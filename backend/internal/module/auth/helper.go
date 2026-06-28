package auth

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"adinfluence_service_app/internal/module/auth/helper"
	"context"
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
)

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
