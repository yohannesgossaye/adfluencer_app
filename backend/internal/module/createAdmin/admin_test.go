package createadmin

import (
	"context"
	"errors"
	"testing"
	"time"

	customErr "adinfluence_service_app/internal/constant/errors"
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	authModule "adinfluence_service_app/internal/module/auth"
	authPersist "adinfluence_service_app/internal/storage/persist/auth"
	createPersist "adinfluence_service_app/internal/storage/persist/createAdmin"
	"adinfluence_service_app/internal/testutil"
	"adinfluence_service_app/platform/jwt"
	"adinfluence_service_app/platform/utils"

	"go.uber.org/zap"
)

// An admin created through the API can log in to the admin portal, the
// password is stored hashed, the creation is audited, and emails stay unique.
func TestCreatedAdminCanLogIn(t *testing.T) {
	db := testutil.MigratedDB(t)
	log := utils.New(zap.NewNop())
	ctx := context.Background()
	m := InitCreateAdminModule(createPersist.InitCreateAdminPersistence(db, log), log)
	tokens := jwt.NewTokenManager(jwt.ManagerConfig{SecretKey: "test-secret-0123456789abcdef0123456789", AccessTokenTTL: 5 * time.Minute})
	auth := authModule.InitAuthModule(authPersist.InitAuthPersistence(db, log), tokens, nil, nil, nil, nil, "http://localhost:3000", "", false, log)

	req := adminDto.CreateAdminRequest{AdminName: " Hana ", AdminEmail: " Hana@Adfluence.ET ", Password: "Strong-pass-1", AdminRole: "super_admin"}
	req.Normalize()
	if err := req.Validate(); err != nil {
		t.Fatal(err)
	}
	id, err := m.CreateAdmin(ctx, req, authDto.SessionMeta{IPAddress: "127.0.0.1", UserAgent: "curl"})
	if err != nil {
		t.Fatalf("create admin: %v", err)
	}

	var hash, role, level, fullName string
	db.QueryRow(`SELECT u.password_hash, u.role, a.permission_level, a.full_name FROM users u JOIN admin_profiles a ON a.user_id = u.id WHERE u.id = $1`, id).
		Scan(&hash, &role, &level, &fullName)
	if hash == "" || hash == "Strong-pass-1" {
		t.Fatal("password not stored as a hash")
	}
	if role != "admin" || level != "super_admin" || fullName != "Hana" {
		t.Fatalf("stored role %q, level %q, name %q", role, level, fullName)
	}
	var audits int
	db.QueryRow(`SELECT count(*) FROM audit_logs WHERE target_user_id = $1 AND action = 'admin_action'`, id).Scan(&audits)
	if audits != 1 {
		t.Fatalf("%d audit entries, want 1", audits)
	}

	// The admin logs in, typing the email in any case.
	res, err := auth.Login(ctx, authDto.LoginRequest{Email: "HANA@adfluence.et", Password: "Strong-pass-1", Portal: "admin"}, authDto.SessionMeta{IPAddress: "198.51.100.9", UserAgent: "test"})
	if err != nil {
		t.Fatalf("new admin cannot log in: %v", err)
	}
	if res.User.Role != "admin" {
		t.Fatalf("logged in as %q", res.User.Role)
	}

	// The same email again is refused clearly.
	if _, err := m.CreateAdmin(ctx, req, authDto.SessionMeta{}); !errors.Is(err, customErr.ErrEmailAlreadyExists) {
		t.Fatalf("duplicate email: %v, want ErrEmailAlreadyExists", err)
	}
}
