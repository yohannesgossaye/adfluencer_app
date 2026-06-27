package helper

import (
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"adinfluence_service_app/platform/ratelimit"
	"time"
)

func ApplyFailedLogin(lockout *authDto.LoginLockout, email string, now time.Time) authDto.LoginLockout {
	if lockout == nil {
		lockout = &authDto.LoginLockout{Email: email}
	}

	if lockout.LockedUntil != nil && lockout.LockedUntil.After(now) {
		return *lockout
	}

	lockout.Email = email
	lockout.FailedAttempts++

	if lockout.FailedAttempts >= ratelimit.LoginMaxFailedAttempts {
		duration := ratelimit.LoginLockoutDuration(lockout.LockoutTier)
		lockedUntil := now.Add(duration)
		lockout.LockedUntil = &lockedUntil
		lockout.LockoutTier++
		lockout.FailedAttempts = 0
	}

	return *lockout
}

func IsLoginLocked(lockout *authDto.LoginLockout, now time.Time) bool {
	return lockout != nil && lockout.LockedUntil != nil && lockout.LockedUntil.After(now)
}
