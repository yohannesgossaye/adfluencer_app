package auth

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"adinfluence_service_app/internal/module/auth/helper"
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
)

const (
	verifyEmailOTPTTL          = 5 * time.Minute
	verifyEmailRateLimit       = 3
	verifyEmailRateLimitWindow = time.Hour
	verifyPhoneOTPTTL          = 5 * time.Minute
	verifyPhoneRateLimit       = 3
	verifyPhoneRateLimitWindow = time.Hour
	forgotPasswordRateScope    = "forgot_password_email"

	// A 6-digit code has 10^6 values; capping wrong guesses per user makes
	// brute force impractical within the code's lifetime.
	verifyEmailAttemptScope = "verify_email_attempt"
	verifyPhoneAttemptScope = "verify_phone_attempt"
	maxOTPFailures          = 5
	otpFailureWindow        = 15 * time.Minute
)

// checkOTPAttempts rejects verification once a user has made too many wrong
// guesses inside the failure window.
func (m *userAuthModule) checkOTPAttempts(ctx context.Context, scope, userID string) error {
	count, err := m.persistence.CountRateLimitEvents(ctx, scope, userID, time.Now().UTC().Add(-otpFailureWindow))
	if err != nil {
		m.logger.Error(ctx, "otp attempt check failed", zap.String("scope", scope), zap.String("user_id", userID), zap.String("reason", err.Error()))
		return err
	}
	if count >= maxOTPFailures {
		m.logger.Warn(ctx, "too many incorrect otp attempts", zap.String("scope", scope), zap.String("user_id", userID))
		return customErr.ErrTooManyOTPAttempts
	}
	return nil
}

// recordOTPFailure counts a wrong code. Lookup errors other than a wrong or
// expired code are not counted against the user.
func (m *userAuthModule) recordOTPFailure(ctx context.Context, scope, userID string, verifyErr error) {
	if !errors.Is(verifyErr, authDto.ErrInvalidVerifyOTP) {
		return
	}
	if err := m.persistence.InsertRateLimitEvent(ctx, scope, userID); err != nil {
		m.logger.Error(ctx, "failed to record otp failure", zap.String("scope", scope), zap.String("user_id", userID), zap.String("reason", err.Error()))
	}
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

func (m *userAuthModule) writeAuditAsync(input authDto.CreateAuditLogInput) {
	go func(input authDto.CreateAuditLogInput) {
		defer func() {
			if r := recover(); r != nil {
				m.logger.Error(context.Background(), "audit log panic recovered",
					zap.Any("panic", r),
					zap.String("action", input.Action),
				)
			}
		}()

		auditCtx, cancel := context.WithTimeout(context.Background(), model.AuditAsyncTimeout)
		defer cancel()

		if err := m.persistence.CreateAuditLog(auditCtx, input); err != nil {
			m.logger.Warn(auditCtx, "audit log failed", zap.String("reason", err.Error()))
		}
	}(input)
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
