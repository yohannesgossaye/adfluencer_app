package ratelimit

import "time"

const (
	LoginMaxFailedAttempts = 5
	ForgotPasswordMaxPerHour = 3
)

var LoginLockoutDurations = []time.Duration{
	5 * time.Minute,
	30 * time.Minute,
	2 * time.Hour,
}

func LoginLockoutDuration(tier int) time.Duration {
	if tier < 0 {
		tier = 0
	}
	if tier >= len(LoginLockoutDurations) {
		return LoginLockoutDurations[len(LoginLockoutDurations)-1]
	}
	return LoginLockoutDurations[tier]
}
