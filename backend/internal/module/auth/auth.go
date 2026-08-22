package auth

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	module "adinfluence_service_app/internal/module"
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

func InitAuthModule(p storage.AuthPersisitence, tokens *jwt.TokenManager, email *mail.Client, smsClient *sms.Client, frontendURL string, log utils.Logger) module.AuthService {
	return &userAuthModule{
		persistence: p,
		tokens:      tokens,
		mail:        email,
		sms:         smsClient,
		frontendURL: strings.TrimRight(frontendURL, "/"),
		logger:      log,
	}
}

func (m *userAuthModule) SignUp(ctx context.Context, req authDto.SignUpRequest, meta authDto.SessionMeta) (*authDto.SignupResponse, error) {
	m.logger.Info(ctx, "processing user signup request",
		zap.String("email", req.Email),
		zap.String("role", req.Role),
		zap.String("name", req.Name),
		zap.String("ip", meta.IPAddress),
	)

	id := utils.GenerateUUID()

	hashedPassword, err := jwt.HashPassword(req.Password)
	if err != nil {
		m.logger.Error(ctx, "unable to hash password during signup", zap.String("email", req.Email), zap.String("reason", err.Error()))
		return nil, err
	}

	isInfluencer := req.Role == string(model.InfluencerRole)
	user := helper.MapperForSignUp(id, req, hashedPassword)

	res, err := m.persistence.SignUp(ctx, user, req.Name, isInfluencer)
	if err != nil {
		m.logger.Error(ctx, "failed to persist signup user", zap.String("email", req.Email), zap.String("reason", err.Error()))
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
		m.logger.Error(ctx, "unable to issue signup tokens", zap.String("user_id", id), zap.String("reason", err.Error()))
		return nil, err
	}

	res.Tokens = helper.TokenResponse(tokenPair)
	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &res.User.ID,
		TargetUserID: &res.User.ID,
		Action:       string(model.SignupAction),
		EntityType:   "user",
		EntityID:     &res.User.ID,
		Description:  "user signed up",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	m.logger.Info(ctx, "user signup completed successfully",
		zap.String("user_id", res.User.ID),
		zap.String("email", res.User.Email),
		zap.String("role", res.User.Role),
	)

	return res, nil
}

func (m *userAuthModule) Login(ctx context.Context, req authDto.LoginRequest, meta authDto.SessionMeta) (*authDto.LoginResponse, error) {
	email := helper.NormalizeEmail(req.Email)
	m.logger.Info(ctx, "processing user login request",
		zap.String("email", email),
		zap.String("portal", req.Portal),
		zap.String("ip", meta.IPAddress),
	)

	if err := m.checkLoginLockout(ctx, email); err != nil {
		m.logger.Warn(ctx, "login blocked due to lockout", zap.String("email", email), zap.String("reason", err.Error()))
		return nil, err
	}

	cred, err := m.persistence.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, customErr.ErrUserNotFound) {
			m.logger.Warn(ctx, "login failed: user not found", zap.String("email", email))
			if lockErr := m.recordLoginFailure(ctx, email); lockErr != nil {
				return nil, lockErr
			}
			return nil, customErr.ErrInvalidCredentials
		}
		m.logger.Error(ctx, "database error querying user by email", zap.String("email", email), zap.String("reason", err.Error()))
		return nil, err
	}

	ok, err := jwt.VerifyPassword(req.Password, cred.PasswordHash)
	if err != nil || !ok {
		m.logger.Warn(ctx, "login failed: invalid password", zap.String("email", email))
		if lockErr := m.recordLoginFailure(ctx, email); lockErr != nil {
			return nil, lockErr
		}
		return nil, customErr.ErrInvalidCredentials
	}

	if err := checkLoginAllowed(cred.UserResponse); err != nil {
		m.logger.Warn(ctx, "login rejected: account status disallowed",
			zap.String("user_id", cred.ID),
			zap.String("status", cred.Status),
			zap.String("reason", err.Error()),
		)
		return nil, err
	}

	if req.Portal == "admin" && cred.Role != string(model.AdminRole) {
		m.logger.Warn(ctx, "non-admin user attempted admin login", zap.String("email", email), zap.String("role", cred.Role))
		return nil, customErr.ErrNotAdmin
	}

	profile, err := m.persistence.GetProfileByUserID(ctx, cred.ID, cred.Role)
	if err != nil {
		m.logger.Error(ctx, "failed to get profile for user", zap.String("user_id", cred.ID), zap.String("reason", err.Error()))
		return nil, err
	}

	if err := m.persistence.ClearLoginLockout(ctx, email); err != nil {
		m.logger.Warn(ctx, "failed to clear login lockout", zap.String("email", email), zap.String("reason", err.Error()))
	}

	now := time.Now().UTC()
	if err := m.persistence.UpdateLastLoginAt(ctx, cred.ID, now); err != nil {
		m.logger.Warn(ctx, "failed to update last login timestamp", zap.String("user_id", cred.ID), zap.String("reason", err.Error()))
	}
	cred.LastLoginAt = &now

	refreshTTL := time.Duration(0)
	if req.RememberMe {
		refreshTTL = m.tokens.RememberMeRefreshTokenTTL()
	}

	user := cred.UserResponse
	tokenPair, err := m.issueSession(ctx, user, *profile, meta, refreshTTL)
	if err != nil {
		m.logger.Error(ctx, "unable to issue login tokens", zap.String("user_id", user.ID), zap.String("reason", err.Error()))
		return nil, err
	}

	redirectTo := helper.RedirectToLogin(user.Role, user.Status)
	if next := helper.SafeInternalPath(req.Next); next != "" {
		redirectTo = next
	}
	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &user.ID,
		TargetUserID: &user.ID,
		Action:       string(model.LoginAction),
		EntityType:   "user",
		EntityID:     &user.ID,
		Description:  "user logged in",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	m.logger.Info(ctx, "user login successful",
		zap.String("user_id", user.ID),
		zap.String("email", user.Email),
		zap.String("role", user.Role),
		zap.String("redirect_to", redirectTo),
	)

	return &authDto.LoginResponse{
		User:       user,
		Profile:    *profile,
		Tokens:     helper.TokenResponse(tokenPair),
		RedirectTo: redirectTo,
	}, nil
}

func (m *userAuthModule) RefreshToken(ctx context.Context, req authDto.RefreshTokenRequest, meta authDto.SessionMeta) (*authDto.RefreshTokenResponse, error) {
	m.logger.Info(ctx, "processing token refresh request", zap.String("ip", meta.IPAddress))

	if req.RefreshToken == "" {
		m.logger.Warn(ctx, "refresh token request missing token")
		return nil, customErr.ErrInvalidRefreshToken
	}

	hash := jwt.HashRefreshToken(req.RefreshToken)
	session, err := m.persistence.GetSessionByRefreshTokenHash(ctx, hash)
	if err != nil {
		m.logger.Warn(ctx, "invalid or expired session for refresh token hash")
		return nil, err
	}

	user, err := m.persistence.GetUserByID(ctx, session.UserID)
	if err != nil {
		m.logger.Error(ctx, "failed to get user for session", zap.String("user_id", session.UserID), zap.String("reason", err.Error()))
		return nil, err
	}

	if user.Status == string(model.BlockedStatus) ||
		user.Status == string(model.RemovedStatus) ||
		user.Status == string(model.SuspendedStatus) {
		_ = m.persistence.RevokeSession(ctx, session.ID)
		m.logger.Warn(ctx, "session revoked due to inactive user status", zap.String("user_id", user.ID), zap.String("status", user.Status))
		return nil, customErr.ErrInvalidRefreshToken
	}

	profile, err := m.persistence.GetProfileByUserID(ctx, user.ID, user.Role)
	if err != nil {
		m.logger.Error(ctx, "failed to get profile during refresh", zap.String("user_id", user.ID), zap.String("reason", err.Error()))
		return nil, err
	}

	if err := m.persistence.RevokeSession(ctx, session.ID); err != nil {
		m.logger.Error(ctx, "failed to rotate session", zap.String("session_id", session.ID), zap.String("reason", err.Error()))
		return nil, err
	}

	tokenPair, err := m.issueSession(ctx, *user, *profile, meta, 0)
	if err != nil {
		m.logger.Error(ctx, "failed to issue rotated session tokens", zap.String("user_id", user.ID), zap.String("reason", err.Error()))
		return nil, err
	}

	m.logger.Info(ctx, "token refreshed successfully", zap.String("user_id", user.ID), zap.String("role", user.Role))

	return &authDto.RefreshTokenResponse{
		Tokens:     helper.TokenResponse(tokenPair),
		RedirectTo: helper.RedirectToLogin(user.Role, user.Status),
	}, nil
}

func (m *userAuthModule) GetMe(ctx context.Context, userID string) (*authDto.MeResponse, error) {
	m.logger.Info(ctx, "fetching current user me profile", zap.String("user_id", userID))

	user, err := m.persistence.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, customErr.ErrUserNotFound) {
			m.logger.Warn(ctx, "get me failed: user not found", zap.String("user_id", userID))
			return nil, customErr.ErrUnauthorized
		}
		m.logger.Error(ctx, "database error in get me", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	if err := checkLoginAllowed(*user); err != nil {
		m.logger.Warn(ctx, "get me rejected: login not allowed", zap.String("user_id", userID), zap.String("status", user.Status))
		return nil, err
	}

	profile, err := m.persistence.GetProfileByUserID(ctx, user.ID, user.Role)
	if err != nil {
		m.logger.Error(ctx, "failed to get profile in get me", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	return &authDto.MeResponse{
		User:    *user,
		Profile: *profile,
	}, nil
}

func (m *userAuthModule) Logout(ctx context.Context, userID string, req authDto.LogoutRequest, meta authDto.SessionMeta) error {
	m.logger.Info(ctx, "processing user logout", zap.String("user_id", userID))

	hash := jwt.HashRefreshToken(req.RefreshToken)

	session, err := m.persistence.GetSessionByRefreshTokenHash(ctx, hash)
	if err != nil {
		m.logger.Warn(ctx, "logout session lookup failed", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return err
	}

	if session.UserID != userID {
		m.logger.Warn(ctx, "logout session user mismatch", zap.String("user_id", userID), zap.String("session_user_id", session.UserID))
		return customErr.ErrUnauthorized
	}

	if err := m.persistence.RevokeSession(ctx, session.ID); err != nil {
		return err
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
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
	m.logger.Info(ctx, "processing forgot password request", zap.String("email", email))

	since := time.Now().UTC().Add(-time.Hour)
	count, err := m.persistence.CountRateLimitEvents(ctx, forgotPasswordRateScope, email, since)
	if err != nil {
		m.logger.Error(ctx, "rate limit check failed for forgot password", zap.String("email", email), zap.String("reason", err.Error()))
		return err
	}
	if count >= ratelimit.ForgotPasswordMaxPerHour {
		m.logger.Warn(ctx, "forgot password rate limit exceeded", zap.String("email", email))
		return customErr.ErrForgotPasswordRateLimit
	}

	if err := m.persistence.InsertRateLimitEvent(ctx, forgotPasswordRateScope, email); err != nil {
		m.logger.Error(ctx, "failed to record forgot password rate limit event", zap.String("email", email), zap.String("reason", err.Error()))
		return err
	}

	cred, err := m.persistence.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, customErr.ErrUserNotFound) {
			m.logger.Info(ctx, "forgot password requested for non-existing user (silent pass)", zap.String("email", email))
			return nil
		}
		m.logger.Error(ctx, "database error querying user for forgot password", zap.String("email", email), zap.String("reason", err.Error()))
		return err
	}

	rawToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		m.logger.Error(ctx, "failed to generate reset token", zap.String("email", email), zap.String("reason", err.Error()))
		return err
	}

	expiresAt := time.Now().UTC().Add(time.Hour)
	if err := m.persistence.CreatePasswordResetToken(ctx, cred.ID, jwt.HashRefreshToken(rawToken), expiresAt); err != nil {
		m.logger.Error(ctx, "failed to persist password reset token", zap.String("user_id", cred.ID), zap.String("reason", err.Error()))
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
		} else {
			m.logger.Info(ctx, "password reset email dispatched", zap.String("email", cred.Email))
		}
		return nil
	}

	m.logger.Info(ctx, "password reset link generated (console fallback)",
		zap.String("email", cred.Email),
		zap.String("link", resetLink),
	)
	return nil
}

func (m *userAuthModule) ResetPassword(ctx context.Context, req authDto.ResetPasswordRequest) (*authDto.ResetPasswordResponse, error) {
	m.logger.Info(ctx, "processing password reset confirmation")

	token, err := m.persistence.GetPasswordResetTokenByHash(ctx, jwt.HashRefreshToken(req.Token))
	if err != nil {
		m.logger.Warn(ctx, "invalid or expired password reset token", zap.String("reason", err.Error()))
		return nil, err
	}

	passwordHash, err := jwt.HashPassword(req.Password)
	if err != nil {
		m.logger.Error(ctx, "unable to hash new password", zap.String("reason", err.Error()))
		return nil, err
	}

	if err := m.persistence.CompletePasswordReset(ctx, token.UserID, token.ID, passwordHash); err != nil {
		m.logger.Error(ctx, "failed to update password in database", zap.String("user_id", token.UserID), zap.String("reason", err.Error()))
		return nil, err
	}

	m.logger.Info(ctx, "password reset completed successfully", zap.String("user_id", token.UserID))

	return &authDto.ResetPasswordResponse{
		RedirectTo: "/auth/login",
	}, nil
}

func (m *userAuthModule) SendVerifyEmail(ctx context.Context, userID string) error {
	m.logger.Info(ctx, "processing send email verification request", zap.String("user_id", userID))

	user, err := m.persistence.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, customErr.ErrUserNotFound) {
			m.logger.Warn(ctx, "user not found for send verify email", zap.String("user_id", userID))
			return customErr.ErrUnauthorized
		}
		m.logger.Error(ctx, "database error querying user for email verify", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return err
	}

	if user.EmailVerified {
		m.logger.Info(ctx, "email already verified", zap.String("user_id", userID))
		return customErr.ErrEmailAlreadyVerified
	}

	since := time.Now().UTC().Add(-verifyEmailRateLimitWindow)
	count, err := m.persistence.CountRecentEmailVerificationTokens(ctx, userID, since)
	if err != nil {
		m.logger.Error(ctx, "rate limit check failed for verify email", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return err
	}
	if count >= verifyEmailRateLimit {
		m.logger.Warn(ctx, "email verification rate limit exceeded", zap.String("user_id", userID))
		return customErr.ErrTooManyVerifyEmailRequests
	}

	code, err := otp.Generate6Digit()
	if err != nil {
		m.logger.Error(ctx, "failed to generate 6-digit OTP", zap.String("reason", err.Error()))
		return err
	}

	now := time.Now().UTC()
	expiresAt := now.Add(verifyEmailOTPTTL)

	if err := m.persistence.InvalidateEmailVerificationTokens(ctx, userID, now); err != nil {
		m.logger.Error(ctx, "failed to invalidate existing verification tokens", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return err
	}

	if err := m.persistence.CreateEmailVerificationToken(ctx, userID, otp.Hash(code), expiresAt); err != nil {
		m.logger.Error(ctx, "failed to store email verification token", zap.String("user_id", userID), zap.String("reason", err.Error()))
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
		} else {
			m.logger.Info(ctx, "verification email dispatched", zap.String("email", user.Email))
		}
		return nil
	}

	m.logger.Info(ctx, "email verification code generated (console fallback)",
		zap.String("email", user.Email),
		zap.String("otp", code),
	)
	return nil
}

func (m *userAuthModule) VerifyEmail(ctx context.Context, userID string, req authDto.VerifyEmailRequest) (*authDto.VerifyEmailResponse, error) {
	m.logger.Info(ctx, "verifying email OTP", zap.String("user_id", userID))

	user, err := m.persistence.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, customErr.ErrUserNotFound) {
			m.logger.Warn(ctx, "user not found during verify email", zap.String("user_id", userID))
			return nil, customErr.ErrUnauthorized
		}
		m.logger.Error(ctx, "database error querying user during email verify", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	if user.EmailVerified {
		m.logger.Info(ctx, "email already verified", zap.String("user_id", userID))
		return nil, customErr.ErrEmailAlreadyVerified
	}

	token, err := m.persistence.GetEmailVerificationTokenByHash(ctx, userID, otp.Hash(req.OTP))
	if err != nil {
		m.logger.Warn(ctx, "invalid or expired email verification OTP", zap.String("user_id", userID))
		return nil, err
	}

	if err := m.persistence.CompleteEmailVerification(ctx, userID, token.ID); err != nil {
		m.logger.Error(ctx, "failed to mark email as verified", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	m.logger.Info(ctx, "email verified successfully", zap.String("user_id", userID), zap.String("email", user.Email))

	return &authDto.VerifyEmailResponse{
		User: authDto.VerifyEmailUserResponse{
			ID:            userID,
			EmailVerified: true,
		},
	}, nil
}

func (m *userAuthModule) SendVerifyPhone(ctx context.Context, userID string) error {
	m.logger.Info(ctx, "processing send phone verification request", zap.String("user_id", userID))

	user, err := m.persistence.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, customErr.ErrUserNotFound) {
			m.logger.Warn(ctx, "user not found for phone verification", zap.String("user_id", userID))
			return customErr.ErrUnauthorized
		}
		m.logger.Error(ctx, "database error querying user for phone verification", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return err
	}

	if user.PhoneVerified {
		m.logger.Info(ctx, "phone already verified", zap.String("user_id", userID))
		return customErr.ErrPhoneAlreadyVerified
	}

	phone := strings.TrimSpace(user.PhoneNumber)
	if phone == "" {
		m.logger.Warn(ctx, "no phone number on user account", zap.String("user_id", userID))
		return customErr.ErrNoPhoneOnAccount
	}

	since := time.Now().UTC().Add(-verifyPhoneRateLimitWindow)
	count, err := m.persistence.CountRecentPhoneVerificationTokens(ctx, userID, since)
	if err != nil {
		m.logger.Error(ctx, "rate limit check failed for phone verify", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return err
	}
	if count >= verifyPhoneRateLimit {
		m.logger.Warn(ctx, "phone verification rate limit exceeded", zap.String("user_id", userID))
		return customErr.ErrTooManyVerifyPhoneRequests
	}

	code, err := otp.Generate6Digit()
	if err != nil {
		m.logger.Error(ctx, "failed to generate phone OTP", zap.String("reason", err.Error()))
		return err
	}

	now := time.Now().UTC()
	expiresAt := now.Add(verifyPhoneOTPTTL)

	if err := m.persistence.InvalidatePhoneVerificationTokens(ctx, userID, now); err != nil {
		m.logger.Error(ctx, "failed to invalidate existing phone verification tokens", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return err
	}

	if err := m.persistence.CreatePhoneVerificationToken(ctx, userID, otp.Hash(code), expiresAt); err != nil {
		m.logger.Error(ctx, "failed to store phone verification token", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return err
	}

	message := fmt.Sprintf("Your Adfluence phone verification code is %s. Valid for 5 minutes.", code)
	if m.sms != nil && m.sms.Enabled() {
		if err := m.sms.Send(phone, message); err != nil {
			m.logger.Warn(ctx, "failed to send SMS OTP",
				zap.String("phone", phone),
				zap.String("reason", err.Error()),
			)
		} else {
			m.logger.Info(ctx, "phone verification SMS sent", zap.String("phone", phone))
		}
		return nil
	}

	m.logger.Info(ctx, "phone verification OTP generated (console fallback)",
		zap.String("phone", phone),
		zap.String("otp", code),
	)
	return nil
}

func (m *userAuthModule) VerifyPhone(ctx context.Context, userID string, req authDto.VerifyPhoneRequest) (*authDto.VerifyPhoneResponse, error) {
	m.logger.Info(ctx, "verifying phone OTP", zap.String("user_id", userID))

	user, err := m.persistence.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, customErr.ErrUserNotFound) {
			m.logger.Warn(ctx, "user not found during phone verify", zap.String("user_id", userID))
			return nil, customErr.ErrUnauthorized
		}
		m.logger.Error(ctx, "database error querying user during phone verify", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	if user.PhoneVerified {
		m.logger.Info(ctx, "phone already verified", zap.String("user_id", userID))
		return nil, customErr.ErrPhoneAlreadyVerified
	}

	token, err := m.persistence.GetPhoneVerificationTokenByHash(ctx, userID, otp.Hash(req.OTP))
	if err != nil {
		m.logger.Warn(ctx, "invalid or expired phone verification OTP", zap.String("user_id", userID))
		return nil, err
	}

	if err := m.persistence.CompletePhoneVerification(ctx, userID, token.ID); err != nil {
		m.logger.Error(ctx, "failed to mark phone as verified", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	m.logger.Info(ctx, "phone verified successfully", zap.String("user_id", userID), zap.String("phone", user.PhoneNumber))

	return &authDto.VerifyPhoneResponse{
		User: authDto.VerifyPhoneUserResponse{
			ID:            userID,
			PhoneVerified: true,
		},
	}, nil
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
