package auth

import (
	"adinfluence_service_app/internal/constant/model"
	customErr "adinfluence_service_app/internal/constant/errors"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"adinfluence_service_app/platform/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

type auth struct {
	db     *sql.DB
	logger utils.Logger
}

func InitAuthPersistence(db *sql.DB, log utils.Logger) *auth {
	return &auth{
		db:     db,
		logger: log,
	}
}

func (a *auth) SignUp(ctx context.Context, user model.User, name string, isInfluencer bool) (*authDto.SignupResponse, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	now := time.Now().UTC()
	profileID := utils.GenerateUUID()

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin signup transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, createUser,
		user.ID,
		user.Email,
		user.PhoneNumber,
		user.PasswordHash,
		user.Role,
		user.Status,
		user.EmailVerified,
		user.PhoneVerified,
		now,
		now,
	)
	if err != nil {
		if mapped := mapUniqueViolation(err); mapped != nil {
			return nil, mapped
		}
		return nil, fmt.Errorf("insert user: %w", err)
	}

	profile := authDto.ProfileResponse{
		ID:         profileID,
		IsVerified: false,
	}

	if isInfluencer {
		_, err = tx.ExecContext(ctx, insertInfluencer,
			profileID,
			user.ID,
			name,
			name,
			false,
			now,
			now,
		)
		profile.FullName = name
		profile.CreatorName = name
	} else {
		_, err = tx.ExecContext(ctx, insertBusiness,
			profileID,
			user.ID,
			name,
			name,
			false,
			now,
			now,
		)
		profile.OwnerFullName = name
		profile.BusinessName = name
	}

	if err != nil {
		return nil, fmt.Errorf("insert profile: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit signup transaction: %w", err)
	}

	return &authDto.SignupResponse{
		User: authDto.UserResponse{
			ID:            user.ID,
			Email:         user.Email,
			PhoneNumber:   user.PhoneNumber,
			Role:          string(user.Role),
			Status:        string(user.Status),
			EmailVerified: user.EmailVerified,
			PhoneVerified: user.PhoneVerified,
			CreatedAt:     now,
		},
		Profile:    profile,
		RedirectTo: signupRedirectTo(user.Role),
	}, nil
}

func (a *auth) DeleteUser(ctx context.Context, userID string) error {
	_, err := a.db.ExecContext(ctx, deleteUser, userID)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

func (a *auth) GetUserByEmail(ctx context.Context, email string) (*authDto.UserCredential, error) {
	var cred authDto.UserCredential
	var lastLogin sql.NullTime
	var suspendedUntil sql.NullTime
	var statusReason sql.NullString

	err := a.db.QueryRowContext(ctx, getUserByEmail, email).Scan(
		&cred.ID,
		&cred.Email,
		&cred.PhoneNumber,
		&cred.PasswordHash,
		&cred.Role,
		&cred.Status,
		&cred.EmailVerified,
		&cred.PhoneVerified,
		&lastLogin,
		&suspendedUntil,
		&statusReason,
		&cred.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErr.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if lastLogin.Valid {
		cred.LastLoginAt = &lastLogin.Time
	}
	if suspendedUntil.Valid {
		cred.SuspendedUntil = &suspendedUntil.Time
	}
	if statusReason.Valid {
		cred.StatusReason = &statusReason.String
	}

	return &cred, nil
}

func (a *auth) UpdateLastLoginAt(ctx context.Context, userID string, at time.Time) error {
	_, err := a.db.ExecContext(ctx, updateLastLoginAt, userID, at)
	if err != nil {
		return fmt.Errorf("update last login: %w", err)
	}
	return nil
}

func (a *auth) CreateSession(ctx context.Context, input authDto.CreateSessionInput) error {
	_, err := a.db.ExecContext(ctx, insertSession,
		utils.GenerateUUID(),
		input.UserID,
		input.RefreshTokenHash,
		nullString(input.IPAddress),
		nullString(input.UserAgent),
		input.ExpiresAt,
		time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert session: %w", err)
	}
	return nil
}

func (a *auth) GetSessionByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*authDto.UserSession, error) {
	var session authDto.UserSession
	err := a.db.QueryRowContext(ctx, getSessionByRefreshHash, refreshTokenHash).Scan(
		&session.ID,
		&session.UserID,
		&session.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErr.ErrInvalidRefreshToken
		}
		return nil, fmt.Errorf("get session: %w", err)
	}
	return &session, nil
}

func (a *auth) RevokeSession(ctx context.Context, sessionID string) error {
	_, err := a.db.ExecContext(ctx, revokeSession, sessionID, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	return nil
}

func (a *auth) GetUserByID(ctx context.Context, userID string) (*authDto.UserResponse, error) {
	var user authDto.UserResponse
	var lastLogin sql.NullTime
	var suspendedUntil sql.NullTime
	var statusReason sql.NullString

	err := a.db.QueryRowContext(ctx, getUserByID, userID).Scan(
		&user.ID,
		&user.Email,
		&user.PhoneNumber,
		&user.Role,
		&user.Status,
		&user.EmailVerified,
		&user.PhoneVerified,
		&lastLogin,
		&suspendedUntil,
		&statusReason,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErr.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	if lastLogin.Valid {
		user.LastLoginAt = &lastLogin.Time
	}
	if suspendedUntil.Valid {
		user.SuspendedUntil = &suspendedUntil.Time
	}
	if statusReason.Valid {
		user.StatusReason = &statusReason.String
	}

	return &user, nil
}

func (a *auth) GetProfileByUserID(ctx context.Context, userID, role string) (*authDto.ProfileResponse, error) {
	profile := &authDto.ProfileResponse{}

	switch model.Role(role) {
	case model.BusinessRole:
		err := a.db.QueryRowContext(ctx, getBusinessProfileByUserID, userID).Scan(
			&profile.ID,
			&profile.BusinessName,
			&profile.OwnerFullName,
			&profile.IsVerified,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, customErr.ErrProfileNotFound
			}
			return nil, fmt.Errorf("get business profile: %w", err)
		}
	case model.InfluencerRole:
		err := a.db.QueryRowContext(ctx, getInfluencerProfileByUserID, userID).Scan(
			&profile.ID,
			&profile.FullName,
			&profile.CreatorName,
			&profile.IsVerified,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, customErr.ErrProfileNotFound
			}
			return nil, fmt.Errorf("get influencer profile: %w", err)
		}
	case model.AdminRole:
		err := a.db.QueryRowContext(ctx, getAdminProfileByUserID, userID).Scan(
			&profile.ID,
			&profile.FullName,
			&profile.PermissionLevel,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, customErr.ErrProfileNotFound
			}
			return nil, fmt.Errorf("get admin profile: %w", err)
		}
	default:
		return nil, customErr.ErrProfileNotFound
	}

	return profile, nil
}

func (a *auth) CreatePasswordResetToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	_, err := a.db.ExecContext(ctx, insertPasswordResetToken,
		utils.GenerateUUID(),
		userID,
		tokenHash,
		expiresAt,
	)
	if err != nil {
		return fmt.Errorf("insert password reset token: %w", err)
	}
	return nil
}

func (a *auth) GetPasswordResetTokenByHash(ctx context.Context, tokenHash string) (*authDto.PasswordResetToken, error) {
	var token authDto.PasswordResetToken
	err := a.db.QueryRowContext(ctx, getPasswordResetTokenByHash, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, authDto.ErrInvalidResetToken
		}
		return nil, fmt.Errorf("get password reset token: %w", err)
	}
	return &token, nil
}

func (a *auth) CompletePasswordReset(ctx context.Context, userID, tokenID, passwordHash string) error {
	now := time.Now().UTC()

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin password reset transaction: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, completePasswordReset, userID, passwordHash, now)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update password rows affected: %w", err)
	}
	if rows == 0 {
		return customErr.ErrUserNotFound
	}

	result, err = tx.ExecContext(ctx, markPasswordResetTokenUsed, tokenID, now)
	if err != nil {
		return fmt.Errorf("mark reset token used: %w", err)
	}

	rows, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("mark reset token rows affected: %w", err)
	}
	if rows == 0 {
		return authDto.ErrInvalidResetToken
	}

	if _, err := tx.ExecContext(ctx, revokeAllUserSessions, userID, now); err != nil {
		return fmt.Errorf("revoke user sessions: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit password reset transaction: %w", err)
	}

	return nil
}

func (a *auth) InvalidateEmailVerificationTokens(ctx context.Context, userID string, usedAt time.Time) error {
	_, err := a.db.ExecContext(ctx, invalidateEmailVerificationTokens, userID, usedAt)
	if err != nil {
		return fmt.Errorf("invalidate email verification tokens: %w", err)
	}
	return nil
}

func (a *auth) CreateEmailVerificationToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	_, err := a.db.ExecContext(ctx, insertEmailVerificationToken,
		utils.GenerateUUID(),
		userID,
		tokenHash,
		expiresAt,
	)
	if err != nil {
		return fmt.Errorf("insert email verification token: %w", err)
	}
	return nil
}

func (a *auth) CountRecentEmailVerificationTokens(ctx context.Context, userID string, since time.Time) (int, error) {
	var count int
	err := a.db.QueryRowContext(ctx, countRecentEmailVerificationTokens, userID, since).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count email verification tokens: %w", err)
	}
	return count, nil
}

func (a *auth) GetEmailVerificationTokenByHash(ctx context.Context, userID, tokenHash string) (*authDto.EmailVerificationToken, error) {
	var token authDto.EmailVerificationToken
	err := a.db.QueryRowContext(ctx, getEmailVerificationTokenByHash, userID, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, authDto.ErrInvalidVerifyOTP
		}
		return nil, fmt.Errorf("get email verification token: %w", err)
	}
	return &token, nil
}

func (a *auth) CompleteEmailVerification(ctx context.Context, userID, tokenID string) error {
	now := time.Now().UTC()

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin email verification transaction: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, markEmailVerified, userID, now)
	if err != nil {
		return fmt.Errorf("mark email verified: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("mark email verified rows affected: %w", err)
	}
	if rows == 0 {
		return customErr.ErrUserNotFound
	}

	result, err = tx.ExecContext(ctx, markEmailVerificationTokenUsed, tokenID, now)
	if err != nil {
		return fmt.Errorf("mark email verification token used: %w", err)
	}

	rows, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("mark email verification token rows affected: %w", err)
	}
	if rows == 0 {
		return authDto.ErrInvalidVerifyOTP
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit email verification transaction: %w", err)
	}

	return nil
}

func (a *auth) InvalidatePhoneVerificationTokens(ctx context.Context, userID string, usedAt time.Time) error {
	_, err := a.db.ExecContext(ctx, invalidatePhoneVerificationTokens, userID, usedAt)
	if err != nil {
		return fmt.Errorf("invalidate phone verification tokens: %w", err)
	}
	return nil
}

func (a *auth) CreatePhoneVerificationToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	_, err := a.db.ExecContext(ctx, insertPhoneVerificationToken,
		utils.GenerateUUID(),
		userID,
		tokenHash,
		expiresAt,
	)
	if err != nil {
		return fmt.Errorf("insert phone verification token: %w", err)
	}
	return nil
}

func (a *auth) CountRecentPhoneVerificationTokens(ctx context.Context, userID string, since time.Time) (int, error) {
	var count int
	err := a.db.QueryRowContext(ctx, countRecentPhoneVerificationTokens, userID, since).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count phone verification tokens: %w", err)
	}
	return count, nil
}

func (a *auth) GetPhoneVerificationTokenByHash(ctx context.Context, userID, tokenHash string) (*authDto.PhoneVerificationToken, error) {
	var token authDto.PhoneVerificationToken
	err := a.db.QueryRowContext(ctx, getPhoneVerificationTokenByHash, userID, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, authDto.ErrInvalidVerifyOTP
		}
		return nil, fmt.Errorf("get phone verification token: %w", err)
	}
	return &token, nil
}

func (a *auth) CompletePhoneVerification(ctx context.Context, userID, tokenID string) error {
	now := time.Now().UTC()

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin phone verification transaction: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, markPhoneVerified, userID, now)
	if err != nil {
		return fmt.Errorf("mark phone verified: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("mark phone verified rows affected: %w", err)
	}
	if rows == 0 {
		return customErr.ErrUserNotFound
	}

	result, err = tx.ExecContext(ctx, markPhoneVerificationTokenUsed, tokenID, now)
	if err != nil {
		return fmt.Errorf("mark phone verification token used: %w", err)
	}

	rows, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("mark phone verification token rows affected: %w", err)
	}
	if rows == 0 {
		return authDto.ErrInvalidVerifyOTP
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit phone verification transaction: %w", err)
	}

	return nil
}

func (a *auth) CreateAuditLog(ctx context.Context, input authDto.CreateAuditLogInput) error {
	_, err := a.db.ExecContext(ctx, insertAuditLog,
		utils.GenerateUUID(),
		nullUUID(input.ActorUserID),
		nullUUID(input.TargetUserID),
		input.Action,
		nullString(input.EntityType),
		nullUUID(input.EntityID),
		nullString(input.Description),
		nullString(input.IPAddress),
		nullString(input.UserAgent),
		time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}
	return nil
}

func (a *auth) GetLoginLockout(ctx context.Context, email string) (*authDto.LoginLockout, error) {
	var lockout authDto.LoginLockout
	var lockedUntil sql.NullTime

	err := a.db.QueryRowContext(ctx, getLoginLockout, email).Scan(
		&lockout.Email,
		&lockout.FailedAttempts,
		&lockout.LockoutTier,
		&lockedUntil,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get login lockout: %w", err)
	}

	if lockedUntil.Valid {
		lockout.LockedUntil = &lockedUntil.Time
	}

	return &lockout, nil
}

func (a *auth) SaveLoginLockout(ctx context.Context, lockout authDto.LoginLockout) error {
	now := time.Now().UTC()
	_, err := a.db.ExecContext(ctx, upsertLoginLockout,
		lockout.Email,
		lockout.FailedAttempts,
		lockout.LockoutTier,
		nullTime(lockout.LockedUntil),
		now,
	)
	if err != nil {
		return fmt.Errorf("save login lockout: %w", err)
	}
	return nil
}

func (a *auth) ClearLoginLockout(ctx context.Context, email string) error {
	_, err := a.db.ExecContext(ctx, deleteLoginLockout, email)
	if err != nil {
		return fmt.Errorf("clear login lockout: %w", err)
	}
	return nil
}

func (a *auth) InsertRateLimitEvent(ctx context.Context, scope, rateKey string) error {
	_, err := a.db.ExecContext(ctx, insertRateLimitEvent,
		utils.GenerateUUID(),
		scope,
		rateKey,
	)
	if err != nil {
		return fmt.Errorf("insert rate limit event: %w", err)
	}
	return nil
}

func (a *auth) CountRateLimitEvents(ctx context.Context, scope, rateKey string, since time.Time) (int, error) {
	var count int
	err := a.db.QueryRowContext(ctx, countRateLimitEvents, scope, rateKey, since).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count rate limit events: %w", err)
	}
	return count, nil
}

func signupRedirectTo(role model.Role) string {
	if role == model.InfluencerRole {
		return "/influencer/dashboard"
	}
	return "/business/dashboard"
}

func mapUniqueViolation(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		return nil
	}

	switch pgErr.ConstraintName {
	case "users_email_key":
		return customErr.ErrEmailAlreadyExists
	case "users_phone_number_key":
		return customErr.ErrPhoneAlreadyExists
	default:
		return customErr.ErrBadRequest
	}
}

func nullString(v string) sql.NullString {
	if v == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: v, Valid: true}
}

func nullUUID(v *string) sql.NullString {
	if v == nil || *v == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: *v, Valid: true}
}

func nullTime(v *time.Time) sql.NullTime {
	if v == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *v, Valid: true}
}
