package auth

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	customErr "adinfluence_service_app/internal/constant/errors"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	authPersist "adinfluence_service_app/internal/storage/persist/auth"
	"adinfluence_service_app/internal/testutil"
	"adinfluence_service_app/platform/jwt"
	"adinfluence_service_app/platform/ratelimit"
	"adinfluence_service_app/platform/utils"

	"go.uber.org/zap"
)

// Wrong passwords are limited per network, per email + network and per email,
// without letting a stranger lock the real owner out.
func TestLoginProtection(t *testing.T) {
	db := testutil.MigratedDB(t)
	log := utils.New(zap.NewNop())
	tokens := jwt.NewTokenManager(jwt.ManagerConfig{SecretKey: "test-secret-0123456789abcdef0123456789", AccessTokenTTL: 5 * time.Minute})
	store := ratelimit.NewMemoryStore()
	m := InitAuthModule(authPersist.InitAuthPersistence(db, log), tokens, nil, nil, nil, store, "http://localhost:3000", "", false, log)
	ctx := context.Background()

	const password = "Test-pass-2026"
	hash, err := jwt.HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO users (id, email, password_hash, role, status, email_verified) VALUES
		('52000000-0000-4000-8000-000000000001', 'owner@example.com', $1, 'business', 'active', TRUE),
		('52000000-0000-4000-8000-000000000002', 'staff@example.com', $1, 'admin', 'active', TRUE),
		('52000000-0000-4000-8000-000000000003', 'popular@example.com', $1, 'business', 'active', TRUE)`, hash); err != nil {
		t.Fatal(err)
	}
	db.Exec(`INSERT INTO business_profiles (user_id, owner_full_name, business_name) VALUES
		('52000000-0000-4000-8000-000000000001', 'Owner', 'Shop'), ('52000000-0000-4000-8000-000000000003', 'Pop', 'Pop Shop')`)
	db.Exec(`INSERT INTO admin_profiles (user_id, full_name) VALUES ('52000000-0000-4000-8000-000000000002', 'Staff')`)

	login := func(email, pass, ip, portal string) error {
		_, err := m.Login(ctx, authDto.LoginRequest{Email: email, Password: pass, Portal: portal}, authDto.SessionMeta{IPAddress: ip, UserAgent: "test"})
		return err
	}
	wantWait := func(t *testing.T, err, kind error, wait time.Duration) {
		t.Helper()
		var w *customErr.WaitError
		if !errors.As(err, &w) || !errors.Is(err, kind) {
			t.Fatalf("got %v, want %v", err, kind)
		}
		if w.Wait > wait || w.Wait < wait-time.Minute {
			t.Fatalf("wait %v, want about %v", w.Wait, wait)
		}
	}

	t.Run("a stranger cannot lock the owner out", func(t *testing.T) {
		for i := 1; i <= 4; i++ {
			if err := login("owner@example.com", "wrong", "198.51.100.1", ""); !errors.Is(err, customErr.ErrInvalidCredentials) {
				t.Fatalf("attempt %d: %v", i, err)
			}
		}
		// The 5th wrong password locks this email from this network for 5 minutes…
		wantWait(t, login("owner@example.com", "wrong", "198.51.100.1", ""), customErr.ErrLoginLocked, 5*time.Minute)
		// …even with the right password now.
		wantWait(t, login("owner@example.com", password, "198.51.100.1", ""), customErr.ErrLoginLocked, 5*time.Minute)
		// The owner on their own network logs in normally.
		if err := login("owner@example.com", password, "203.0.113.7", ""); err != nil {
			t.Fatalf("owner locked out by a stranger: %v", err)
		}
	})

	t.Run("locks get longer", func(t *testing.T) {
		k := keysFor("owner@example.com", "198.51.100.1")
		store.Del(ctx, k.pairLock) // the first lock is over
		for i := 1; i <= 4; i++ {
			login("owner@example.com", "wrong", "198.51.100.1", "")
		}
		wantWait(t, login("owner@example.com", "wrong", "198.51.100.1", ""), customErr.ErrLoginLocked, 15*time.Minute)
	})

	t.Run("a correct password forgets earlier mistakes", func(t *testing.T) {
		for round := 0; round < 2; round++ {
			for i := 1; i <= 4; i++ {
				if err := login("popular@example.com", "wrong", "198.51.100.2", ""); !errors.Is(err, customErr.ErrInvalidCredentials) {
					t.Fatalf("round %d attempt %d: %v", round, i, err)
				}
			}
			if err := login("popular@example.com", password, "198.51.100.2", ""); err != nil {
				t.Fatalf("typo then success locked the owner: %v", err)
			}
		}
	})

	t.Run("unknown emails look exactly like wrong passwords", func(t *testing.T) {
		for i := 1; i <= 4; i++ {
			if err := login("nobody@example.com", "wrong", "198.51.100.3", ""); !errors.Is(err, customErr.ErrInvalidCredentials) {
				t.Fatalf("attempt %d: %v", i, err)
			}
		}
		wantWait(t, login("nobody@example.com", "wrong", "198.51.100.3", ""), customErr.ErrLoginLocked, 5*time.Minute)
	})

	t.Run("one network trying many accounts is paused", func(t *testing.T) {
		ip := "198.51.100.4"
		for i := 1; i < 20; i++ {
			if err := login(fmt.Sprintf("guess%d@example.com", i), "123456", ip, ""); !errors.Is(err, customErr.ErrInvalidCredentials) {
				t.Fatalf("attempt %d: %v", i, err)
			}
		}
		wantWait(t, login("guess20@example.com", "123456", ip, ""), customErr.ErrTooManyAttempts, 10*time.Minute)
		// Paused for every account from that network, right password or not.
		wantWait(t, login("popular@example.com", password, ip, ""), customErr.ErrTooManyAttempts, 10*time.Minute)
		// Other networks are unaffected.
		if err := login("popular@example.com", password, "203.0.113.8", ""); err != nil {
			t.Fatalf("other network refused: %v", err)
		}
	})

	t.Run("an attack from many networks pauses the email", func(t *testing.T) {
		// 30 wrong passwords from 6 networks (5 each, each network gets locked).
		var err error
		for n := 1; n <= 6; n++ {
			for i := 1; i <= 5; i++ {
				err = login("popular@example.com", "wrong", fmt.Sprintf("192.0.2.%d", n), "")
			}
		}
		wantWait(t, err, customErr.ErrLoginLocked, 15*time.Minute)
		// Now paused everywhere, including a network that never failed.
		wantWait(t, login("popular@example.com", password, "203.0.113.9", ""), customErr.ErrLoginLocked, 15*time.Minute)
	})

	t.Run("admin accounts lock after 3", func(t *testing.T) {
		for i := 1; i <= 2; i++ {
			if err := login("staff@example.com", "wrong", "198.51.100.5", "admin"); !errors.Is(err, customErr.ErrInvalidCredentials) {
				t.Fatalf("attempt %d: %v", i, err)
			}
		}
		wantWait(t, login("staff@example.com", "wrong", "198.51.100.5", "admin"), customErr.ErrLoginLocked, 15*time.Minute)
		// Stricter wherever the admin logs in, not just the admin page.
		for i := 1; i <= 2; i++ {
			login("staff@example.com", "wrong", "198.51.100.6", "")
		}
		wantWait(t, login("staff@example.com", "wrong", "198.51.100.6", ""), customErr.ErrLoginLocked, 15*time.Minute)
	})

	t.Run("the message says how long to wait", func(t *testing.T) {
		err := login("nobody@example.com", "wrong", "198.51.100.3", "")
		if got, want := err.Error(), "Too many failed attempts. Try again in 5 minutes, or reset your password."; got != want {
			t.Fatalf("message %q, want %q", got, want)
		}
	})

	t.Run("sign-ups and reset requests are limited per network", func(t *testing.T) {
		var err error
		for i := 1; i <= ratelimit.ForgotPasswordPerIPPerHour+1; i++ {
			err = m.ForgotPassword(ctx, authDto.ForgotPasswordRequest{Email: fmt.Sprintf("r%d@example.com", i)}, authDto.SessionMeta{IPAddress: "198.51.100.7"})
		}
		if !errors.Is(err, customErr.ErrForgotPasswordRateLimit) {
			t.Fatalf("reset request %d from one network: %v", ratelimit.ForgotPasswordPerIPPerHour+1, err)
		}
		store.Set(ctx, "rl:signup:ip:198.51.100.8", int64(ratelimit.SignupPerIPPerHour), time.Hour)
		_, err = m.SignUp(ctx, authDto.SignUpRequest{Name: "X", Email: "new@example.com", Password: password, Role: "business"}, authDto.SessionMeta{IPAddress: "198.51.100.8"})
		if !errors.Is(err, customErr.ErrTooManyAttempts) {
			t.Fatalf("sign-up past the limit: %v", err)
		}
	})
}
