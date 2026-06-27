package auth

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"adinfluence_service_app/internal/module/auth/helper"
	"adinfluence_service_app/internal/storage"
	"adinfluence_service_app/platform/jwt"
	"adinfluence_service_app/platform/mail"
	"adinfluence_service_app/platform/otp"
	"adinfluence_service_app/platform/ratelimit"
	"adinfluence_service_app/platform/sms"
	"adinfluence_service_app/platform/utils"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

type userAuthModule struct {
	persistence storage.AuthPersisitence
	tokens      *jwt.TokenManager
	mail        *mail.Client
	sms         *sms.Client
	frontendURL string
	logger      utils.Logger
}

func InitAuthModule(p storage.AuthPersisitence, tokens *jwt.TokenManager, email *mail.Client, smsClient *sms.Client, frontendURL string, log utils.Logger) *userAuthModule {
	return &userAuthModule{
		persistence: p,
		tokens:      tokens,
		mail:        email,
		sms:         smsClient,
		frontendURL: strings.TrimRight(frontendURL, "/"),
		logger:      log,
	}
}

const forgotPasswordRateScope = "forgot_password_email"

func (m *userAuthModule) SignUp(ctx context.Context, req authDto.SignUpRequest, meta authDto.SessionMeta) (*authDto.SignupResponse, error) {
	id := utils.GenerateUUID()

	hashedPassword, err := jwt.HashPassword(req.Password)
	if err != nil {
		m.logger.Error(ctx, "unable to hash password", zap.String("reason", err.Error()))
		return nil, err
	}

	isInfluencer := req.Role == string(model.InfluencerRole)
	user := helper.MapperForSignUp(id, req, hashedPassword)

	res, err := m.persistence.SignUp(ctx, user, req.Name, isInfluencer)
	if err != nil {
		return nil, err
	}

	tokenPair, err := m.issueSession(ctx, res.User, res.Profile, meta, 0)
	if err != nil {
		if delErr := m.persistence.DeleteUser(ctx, id); delErr != nil {
			m.logger.Error(ctx, "signup rollback failed",
				zap.String("user_id", id),
				zap.String("reason", delErr.Error()),
			)
		}
		m.logger.Error(ctx, "unable to issue signup tokens", zap.String("reason", err.Error()))
		return nil, err
	}

	res.Tokens = helper.TokenResponse(tokenPair)
	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID:  &res.User.ID,
		TargetUserID: &res.User.ID,
		Action:       string(model.SignupAction),
		EntityType:   "user",
		EntityID:     &res.User.ID,
		Description:  "user signed up",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})
	return res, nil
}

func (m *userAuthModule) Login(ctx context.Context, req authDto.LoginRequest, meta authDto.SessionMeta) (*authDto.LoginResponse, error) {
	email := helper.NormalizeEmail(req.Email)

	if err := m.checkLoginLockout(ctx, email); err != nil {
		return nil, err
	}

	cred, err := m.persistence.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, authDto.ErrUserNotFound) {
			if lockErr := m.recordLoginFailure(ctx, email); lockErr != nil {
				return nil, lockErr
			}
			return nil, customErr.ErrInvalidCredentials
		}
		return nil, err
	}

	ok, err := jwt.VerifyPassword(req.Password, cred.PasswordHash)
	if err != nil || !ok {
		if lockErr := m.recordLoginFailure(ctx, email); lockErr != nil {
			return nil, lockErr
		}
		return nil, customErr.ErrInvalidCredentials
	}

	if err := checkLoginAllowed(cred.UserResponse); err != nil {
		return nil, err
	}

	if req.Portal == "admin" && cred.Role != string(model.AdminRole) {
		return nil, customErr.ErrNotAdmin
	}

	profile, err := m.persistence.GetProfileByUserID(ctx, cred.ID, cred.Role)
	if err != nil {
		return nil, err
	}

	if err := m.persistence.ClearLoginLockout(ctx, email); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if err := m.persistence.UpdateLastLoginAt(ctx, cred.ID, now); err != nil {
		return nil, err
	}
	cred.LastLoginAt = &now

	refreshTTL := time.Duration(0)
	if req.RememberMe {
		refreshTTL = m.tokens.RememberMeRefreshTokenTTL()
	}

	user := cred.UserResponse
	tokenPair, err := m.issueSession(ctx, user, *profile, meta, refreshTTL)
	if err != nil {
		m.logger.Error(ctx, "unable to issue login tokens", zap.String("reason", err.Error()))
		return nil, err
	}

	redirectTo := helper.RedirectToLogin(user.Role, user.Status)
	if next := helper.SafeInternalPath(req.Next); next != "" {
		redirectTo = next
	}

	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID:  &user.ID,
		TargetUserID: &user.ID,
		Action:       string(model.LoginAction),
		EntityType:   "user",
		EntityID:     &user.ID,
		Description:  "user logged in",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	return &authDto.LoginResponse{
		User:       user,
		Profile:    *profile,
		Tokens:     helper.TokenResponse(tokenPair),
		RedirectTo: redirectTo,
	}, nil
}

func (m *userAuthModule) RefreshToken(ctx context.Context, req authDto.RefreshTokenRequest, meta authDto.SessionMeta) (*authDto.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, authDto.ErrInvalidRefreshToken
	}

	hash := jwt.HashRefreshToken(req.RefreshToken)
	session, err := m.persistence.GetSessionByRefreshTokenHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	user, err := m.persistence.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	if user.Status == string(model.BlockedStatus) ||
		user.Status == string(model.RemovedStatus) ||
		user.Status == string(model.SuspendedStatus) {
		_ = m.persistence.RevokeSession(ctx, session.ID)
		return nil, authDto.ErrInvalidRefreshToken
	}

	profile, err := m.persistence.GetProfileByUserID(ctx, user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	if err := m.persistence.RevokeSession(ctx, session.ID); err != nil {
		return nil, err
	}

	tokenPair, err := m.issueSession(ctx, *user, *profile, meta, 0)
	if err != nil {
		return nil, err
	}

	return &authDto.RefreshTokenResponse{
		Tokens:     helper.TokenResponse(tokenPair),
		RedirectTo: helper.RedirectToLogin(user.Role, user.Status),
	}, nil
}

func (m *userAuthModule) GetMe(ctx context.Context, userID string) (*authDto.MeResponse, error) {
	user, err := m.persistence.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, authDto.ErrUserNotFound) {
			return nil, customErr.ErrUnauthorized
		}
		return nil, err
	}

	if err := checkLoginAllowed(*user); err != nil {
		return nil, err
	}

	profile, err := m.persistence.GetProfileByUserID(ctx, user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &authDto.MeResponse{
		User:    *user,
		Profile: *profile,
	}, nil
}

func (m *userAuthModule) Logout(ctx context.Context, userID string, req authDto.LogoutRequest, meta authDto.SessionMeta) error {
	hash := jwt.HashRefreshToken(req.RefreshToken)

	session, err := m.persistence.GetSessionByRefreshTokenHash(ctx, hash)
	if err != nil {
		return err
	}

	if session.UserID != userID {
		return customErr.ErrUnauthorized
	}

	if err := m.persistence.RevokeSession(ctx, session.ID); err != nil {
		return err
	}

	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.LogoutAction),
		EntityType:   "user",
		EntityID:     &userID,
		Description:  "user logged out",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	return nil
}

func (m *userAuthModule) ForgotPassword(ctx context.Context, req authDto.ForgotPasswordRequest) error {
	email := helper.NormalizeEmail(req.Email)

	since := time.Now().UTC().Add(-time.Hour)
	count, err := m.persistence.CountRateLimitEvents(ctx, forgotPasswordRateScope, email, since)
	if err != nil {
		return err
	}
	if count >= ratelimit.ForgotPasswordMaxPerHour {
		return customErr.ErrForgotPasswordRateLimit
	}

	if err := m.persistence.InsertRateLimitEvent(ctx, forgotPasswordRateScope, email); err != nil {
		return err
	}

	cred, err := m.persistence.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, authDto.ErrUserNotFound) {
			return nil
		}
		return err
	}

	rawToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		return err
	}

	expiresAt := time.Now().UTC().Add(time.Hour)
	if err := m.persistence.CreatePasswordResetToken(ctx, cred.ID, jwt.HashRefreshToken(rawToken), expiresAt); err != nil {
		return err
	}

	resetLink := fmt.Sprintf("%s/auth/reset-password?token=%s", m.frontendURL, rawToken)
	body := fmt.Sprintf("Reset your Adfluence password using this link:\n\n%s\n\nThis link expires in 1 hour.", resetLink)

	if m.mail != nil && m.mail.Enabled() {
		if err := m.mail.Send(cred.Email, "Reset your Adfluence password", body); err != nil {
			m.logger.Warn(ctx, "failed to send password reset email",
				zap.String("email", cred.Email),
				zap.String("reason", err.Error()),
			)
		}
		return nil
	}

	m.logger.Info(ctx, "password reset link generated",
		zap.String("email", cred.Email),
		zap.String("link", resetLink),
	)
	return nil
}

func (m *userAuthModule) ResetPassword(ctx context.Context, req authDto.ResetPasswordRequest) (*authDto.ResetPasswordResponse, error) {
	token, err := m.persistence.GetPasswordResetTokenByHash(ctx, jwt.HashRefreshToken(req.Token))
	if err != nil {
		return nil, err
	}

	passwordHash, err := jwt.HashPassword(req.Password)
	if err != nil {
		m.logger.Error(ctx, "unable to hash password", zap.String("reason", err.Error()))
		return nil, err
	}

	if err := m.persistence.CompletePasswordReset(ctx, token.UserID, token.ID, passwordHash); err != nil {
		return nil, err
	}

	return &authDto.ResetPasswordResponse{
		RedirectTo: "/auth/login",
	}, nil
}

const (
	verifyEmailOTPTTL          = 5 * time.Minute
	verifyEmailRateLimit       = 3
	verifyEmailRateLimitWindow = time.Hour
)

func (m *userAuthModule) SendVerifyEmail(ctx context.Context, userID string) error {
	user, err := m.persistence.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, authDto.ErrUserNotFound) {
			return customErr.ErrUnauthorized
		}
		return err
	}

	if user.EmailVerified {
		return customErr.ErrEmailAlreadyVerified
	}

	since := time.Now().UTC().Add(-verifyEmailRateLimitWindow)
	count, err := m.persistence.CountRecentEmailVerificationTokens(ctx, userID, since)
	if err != nil {
		return err
	}
	if count >= verifyEmailRateLimit {
		return customErr.ErrTooManyVerifyEmailRequests
	}

	code, err := otp.Generate6Digit()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	expiresAt := now.Add(verifyEmailOTPTTL)

	if err := m.persistence.InvalidateEmailVerificationTokens(ctx, userID, now); err != nil {
		return err
	}

	if err := m.persistence.CreateEmailVerificationToken(ctx, userID, otp.Hash(code), expiresAt); err != nil {
		return err
	}

	body := fmt.Sprintf("Your Adfluence email verification code is: %s\n\nThis code expires in 5 minutes.", code)
	subject := "Verify your Adfluence email"

	if m.mail != nil && m.mail.Enabled() {
		if err := m.mail.Send(user.Email, subject, body); err != nil {
			m.logger.Warn(ctx, "failed to send verification email",
				zap.String("email", user.Email),
				zap.String("reason", err.Error()),
			)
		}
		return nil
	}

	m.logger.Info(ctx, "email verification code generated",
		zap.String("email", user.Email),
		zap.String("otp", code),
	)
	return nil
}

func (m *userAuthModule) VerifyEmail(ctx context.Context, userID string, req authDto.VerifyEmailRequest) (*authDto.VerifyEmailResponse, error) {
	user, err := m.persistence.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, authDto.ErrUserNotFound) {
			return nil, customErr.ErrUnauthorized
		}
		return nil, err
	}

	if user.EmailVerified {
		return nil, customErr.ErrEmailAlreadyVerified
	}

	token, err := m.persistence.GetEmailVerificationTokenByHash(ctx, userID, otp.Hash(req.OTP))
	if err != nil {
		return nil, err
	}

	if err := m.persistence.CompleteEmailVerification(ctx, userID, token.ID); err != nil {
		return nil, err
	}

	return &authDto.VerifyEmailResponse{
		User: authDto.VerifyEmailUserResponse{
			ID:            userID,
			EmailVerified: true,
		},
	}, nil
}

const (
	verifyPhoneOTPTTL          = 5 * time.Minute
	verifyPhoneRateLimit       = 3
	verifyPhoneRateLimitWindow = time.Hour
)

func (m *userAuthModule) SendVerifyPhone(ctx context.Context, userID string) error {
	user, err := m.persistence.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, authDto.ErrUserNotFound) {
			return customErr.ErrUnauthorized
		}
		return err
	}

	if user.PhoneVerified {
		return customErr.ErrPhoneAlreadyVerified
	}

	phone := strings.TrimSpace(user.PhoneNumber)
	if phone == "" {
		return customErr.ErrNoPhoneOnAccount
	}

	since := time.Now().UTC().Add(-verifyPhoneRateLimitWindow)
	count, err := m.persistence.CountRecentPhoneVerificationTokens(ctx, userID, since)
	if err != nil {
		return err
	}
	if count >= verifyPhoneRateLimit {
		return customErr.ErrTooManyVerifyPhoneRequests
	}

	code, err := otp.Generate6Digit()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	expiresAt := now.Add(verifyPhoneOTPTTL)

	if err := m.persistence.InvalidatePhoneVerificationTokens(ctx, userID, now); err != nil {
		return err
	}

	if err := m.persistence.CreatePhoneVerificationToken(ctx, userID, otp.Hash(code), expiresAt); err != nil {
		return err
	}

	message := fmt.Sprintf("Your Adfluence phone verification code is: %s. Expires in 5 minutes.", code)
	smsTo := "+" + utils.NormalizePhone(phone)

	if m.sms != nil && m.sms.Enabled() {
		if err := m.sms.Send(smsTo, message); err != nil {
			m.logger.Warn(ctx, "failed to send verification sms",
				zap.String("phone", smsTo),
				zap.String("reason", err.Error()),
			)
		}
		return nil
	}

	m.logger.Info(ctx, "phone verification code generated",
		zap.String("phone", smsTo),
		zap.String("otp", code),
	)
	return nil
}

func (m *userAuthModule) VerifyPhone(ctx context.Context, userID string, req authDto.VerifyPhoneRequest) (*authDto.VerifyPhoneResponse, error) {
	user, err := m.persistence.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, authDto.ErrUserNotFound) {
			return nil, customErr.ErrUnauthorized
		}
		return nil, err
	}

	if user.PhoneVerified {
		return nil, customErr.ErrPhoneAlreadyVerified
	}

	token, err := m.persistence.GetPhoneVerificationTokenByHash(ctx, userID, otp.Hash(req.OTP))
	if err != nil {
		return nil, err
	}

	if err := m.persistence.CompletePhoneVerification(ctx, userID, token.ID); err != nil {
		return nil, err
	}

	return &authDto.VerifyPhoneResponse{
		User: authDto.VerifyPhoneUserResponse{
			ID:            userID,
			PhoneVerified: true,
		},
	}, nil
}

func checkLoginAllowed(user authDto.UserResponse) error {
	switch model.UserStatus(user.Status) {
	case model.BlockedStatus:
		return customErr.ErrAccountBlocked
	case model.RemovedStatus:
		return customErr.ErrAccountRemoved
	case model.SuspendedStatus:
		if user.SuspendedUntil != nil && user.SuspendedUntil.After(time.Now().UTC()) {
			return customErr.ErrAccountSuspended
		}
	}
	return nil
}

func (m *userAuthModule) issueSession(ctx context.Context, user authDto.UserResponse, profile authDto.ProfileResponse, meta authDto.SessionMeta, refreshTTL time.Duration) (jwt.TokenPair, error) {
	tokenPair, err := m.tokens.GenerateTokenPair(helper.AccessTokenInput(user, profile))
	if err != nil {
		return jwt.TokenPair{}, fmt.Errorf("generate token pair: %w", err)
	}

	if refreshTTL > 0 {
		tokenPair.RefreshExpiresAt = time.Now().UTC().Add(refreshTTL)
	}

	err = m.persistence.CreateSession(ctx, authDto.CreateSessionInput{
		UserID:           user.ID,
		RefreshTokenHash: jwt.HashRefreshToken(tokenPair.RefreshToken),
		IPAddress:        meta.IPAddress,
		UserAgent:        meta.UserAgent,
		ExpiresAt:        tokenPair.RefreshExpiresAt,
	})
	if err != nil {
		return jwt.TokenPair{}, err
	}

	return tokenPair, nil
}

func (m *userAuthModule) writeAudit(ctx context.Context, input authDto.CreateAuditLogInput) {
	if err := m.persistence.CreateAuditLog(ctx, input); err != nil {
		m.logger.Warn(ctx, "audit log failed", zap.String("reason", err.Error()))
	}
}

func (m *userAuthModule) checkLoginLockout(ctx context.Context, email string) error {
	lockout, err := m.persistence.GetLoginLockout(ctx, email)
	if err != nil {
		return err
	}
	if helper.IsLoginLocked(lockout, time.Now().UTC()) {
		return customErr.ErrLoginLocked
	}
	return nil
}

func (m *userAuthModule) recordLoginFailure(ctx context.Context, email string) error {
	now := time.Now().UTC()

	lockout, err := m.persistence.GetLoginLockout(ctx, email)
	if err != nil {
		return err
	}

	updated := helper.ApplyFailedLogin(lockout, email, now)
	if err := m.persistence.SaveLoginLockout(ctx, updated); err != nil {
		return err
	}

	if helper.IsLoginLocked(&updated, now) {
		return customErr.ErrLoginLocked
	}

	return nil
}
