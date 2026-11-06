package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	customErr "adinfluence_service_app/internal/constant/errors"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	authPersist "adinfluence_service_app/internal/storage/persist/auth"
	"adinfluence_service_app/internal/testutil"
	"adinfluence_service_app/platform/jwt"
	"adinfluence_service_app/platform/utils"

	"go.uber.org/zap"
)

// Bank-style sessions: they end after inactivity and after a maximum age,
// admins get stricter limits, and "remember me" no longer extends anything.
func TestSessionLimits(t *testing.T) {
	db := testutil.MigratedDB(t)
	log := utils.New(zap.NewNop())
	tokens := jwt.NewTokenManager(jwt.ManagerConfig{
		SecretKey:      "test-secret-0123456789abcdef0123456789",
		AccessTokenTTL: 5 * time.Minute,
		Sessions: jwt.SessionPolicy{
			Idle: 15 * time.Minute, MaxAge: 12 * time.Hour,
			AdminIdle: 10 * time.Minute, AdminMaxAge: 8 * time.Hour,
		},
	})
	m := InitAuthModule(authPersist.InitAuthPersistence(db, log), tokens, nil, nil, nil, nil, "http://localhost:3000", "", false, log)
	ctx := context.Background()
	meta := authDto.SessionMeta{IPAddress: "127.0.0.1", UserAgent: "test"}

	hash, err := jwt.HashPassword("Test-pass-2026")
	if err != nil {
		t.Fatal(err)
	}
	exec := func(q string, args ...any) {
		t.Helper()
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatal(err)
		}
	}
	exec(`INSERT INTO users (id, email, password_hash, role, status, email_verified) VALUES
		('51000000-0000-4000-8000-000000000001', 'biz@example.com', $1, 'business', 'active', TRUE),
		('51000000-0000-4000-8000-000000000002', 'admin@example.com', $1, 'admin', 'active', TRUE)`, hash)
	exec(`INSERT INTO business_profiles (user_id, owner_full_name, business_name) VALUES ('51000000-0000-4000-8000-000000000001', 'Owner', 'Shop')`)
	exec(`INSERT INTO admin_profiles (user_id, full_name) VALUES ('51000000-0000-4000-8000-000000000002', 'Admin')`)

	login := func(email, portal string, remember bool) authDto.TokenResponse {
		t.Helper()
		res, err := m.Login(ctx, authDto.LoginRequest{Email: email, Password: "Test-pass-2026", Portal: portal, RememberMe: remember}, meta)
		if err != nil {
			t.Fatalf("login %s: %v", email, err)
		}
		return res.Tokens
	}
	refresh := func(refreshToken string) (authDto.TokenResponse, error) {
		res, err := m.RefreshToken(ctx, authDto.RefreshTokenRequest{RefreshToken: refreshToken}, meta)
		if err != nil {
			return authDto.TokenResponse{}, err
		}
		return res.Tokens, nil
	}
	latestSession := `(SELECT id FROM user_sessions WHERE revoked_at IS NULL ORDER BY created_at DESC LIMIT 1)`

	// "Remember me" no longer extends the session.
	biz := login("biz@example.com", "", true)
	if biz.IdleTimeoutSeconds != 15*60 || biz.ExpiresIn != 5*60 {
		t.Fatalf("business: idle %ds, access %ds; want 900 and 300", biz.IdleTimeoutSeconds, biz.ExpiresIn)
	}
	if d := time.Until(biz.SessionEndsAt); d < 11*time.Hour+59*time.Minute || d > 12*time.Hour {
		t.Fatalf("business session ends in %v, want 12h", d)
	}
	var expires time.Time
	db.QueryRow(`SELECT expires_at FROM user_sessions WHERE id = ` + latestSession).Scan(&expires)
	if d := time.Until(expires); d > 15*time.Minute || d < 14*time.Minute {
		t.Fatalf("refresh token valid for %v, want 15m despite remember me", d)
	}

	admin := login("admin@example.com", "admin", false)
	if admin.IdleTimeoutSeconds != 10*60 || time.Until(admin.SessionEndsAt) > 8*time.Hour {
		t.Fatalf("admin: idle %ds, ends in %v; want 600 and 8h", admin.IdleTimeoutSeconds, time.Until(admin.SessionEndsAt))
	}

	// Activity renews the idle period but keeps the same end.
	next, err := refresh(biz.RefreshToken)
	if err != nil {
		t.Fatalf("refresh while active: %v", err)
	}
	if !next.SessionEndsAt.Equal(biz.SessionEndsAt) {
		t.Fatalf("refresh moved the session end from %v to %v", biz.SessionEndsAt, next.SessionEndsAt)
	}

	// Idle too long: the server refuses to continue.
	exec(`UPDATE user_sessions SET expires_at = NOW() - interval '1 second' WHERE id = ` + latestSession)
	if _, err := refresh(next.RefreshToken); !errors.Is(err, customErr.ErrInvalidRefreshToken) {
		t.Fatalf("idle session refreshed: err = %v", err)
	}

	// Active, but past the maximum age: refused as well.
	again := login("biz@example.com", "", false)
	exec(`UPDATE user_sessions SET started_at = NOW() - interval '13 hours' WHERE id = ` + latestSession)
	if _, err := refresh(again.RefreshToken); !errors.Is(err, customErr.ErrInvalidRefreshToken) {
		t.Fatalf("session older than 12h refreshed: err = %v", err)
	}

	// A session about to hit its maximum age gets a refresh token that ends with it.
	near := login("biz@example.com", "", false)
	exec(`UPDATE user_sessions SET started_at = NOW() - interval '11 hours 55 minutes' WHERE id = ` + latestSession)
	capped, err := refresh(near.RefreshToken)
	if err != nil {
		t.Fatalf("refresh near the maximum age: %v", err)
	}
	db.QueryRow(`SELECT expires_at FROM user_sessions WHERE id = ` + latestSession).Scan(&expires)
	if d := time.Until(expires); d > 5*time.Minute+5*time.Second {
		t.Fatalf("refresh token outlives the session: valid for %v, session ends in %v", d, time.Until(capped.SessionEndsAt))
	}
}
