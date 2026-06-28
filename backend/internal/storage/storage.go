package storage

import (
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"context"
	"time"
)

type AuthPersisitence interface {
	SignUp(ctx context.Context, user model.User, name string, isInfluencer bool) (*authDto.SignupResponse, error)
	DeleteUser(ctx context.Context, userID string) error
	GetUserByEmail(ctx context.Context, email string) (*authDto.UserCredential, error)
	UpdateLastLoginAt(ctx context.Context, userID string, at time.Time) error
	CreateSession(ctx context.Context, input authDto.CreateSessionInput) error
	GetSessionByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*authDto.UserSession, error)
	RevokeSession(ctx context.Context, sessionID string) error
	GetUserByID(ctx context.Context, userID string) (*authDto.UserResponse, error)
	GetProfileByUserID(ctx context.Context, userID, role string) (*authDto.ProfileResponse, error)
	CreatePasswordResetToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	GetPasswordResetTokenByHash(ctx context.Context, tokenHash string) (*authDto.PasswordResetToken, error)
	CompletePasswordReset(ctx context.Context, userID, tokenID, passwordHash string) error
	InvalidateEmailVerificationTokens(ctx context.Context, userID string, usedAt time.Time) error
	CreateEmailVerificationToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	CountRecentEmailVerificationTokens(ctx context.Context, userID string, since time.Time) (int, error)
	GetEmailVerificationTokenByHash(ctx context.Context, userID, tokenHash string) (*authDto.EmailVerificationToken, error)
	CompleteEmailVerification(ctx context.Context, userID, tokenID string) error
	InvalidatePhoneVerificationTokens(ctx context.Context, userID string, usedAt time.Time) error
	CreatePhoneVerificationToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	CountRecentPhoneVerificationTokens(ctx context.Context, userID string, since time.Time) (int, error)
	GetPhoneVerificationTokenByHash(ctx context.Context, userID, tokenHash string) (*authDto.PhoneVerificationToken, error)
	CompletePhoneVerification(ctx context.Context, userID, tokenID string) error
	CreateAuditLog(ctx context.Context, input authDto.CreateAuditLogInput) error
	GetLoginLockout(ctx context.Context, email string) (*authDto.LoginLockout, error)
	SaveLoginLockout(ctx context.Context, lockout authDto.LoginLockout) error
	ClearLoginLockout(ctx context.Context, email string) error
	InsertRateLimitEvent(ctx context.Context, scope, rateKey string) error
	CountRateLimitEvents(ctx context.Context, scope, rateKey string, since time.Time) (int, error)
}
