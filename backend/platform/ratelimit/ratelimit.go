// Package ratelimit holds the limits that slow down password guessing and
// abuse, and the short-lived counters they are counted in.
package ratelimit

import "time"

// LoginPolicy is how many wrong passwords are tolerated before logins pause.
// Three layers, each counting only failed attempts:
//
//   - Network (IP): many wrong passwords from one IP, for any emails, pause
//     that IP. Stops one machine trying many accounts. Mobile networks put many
//     people behind one IP, so this number is generous.
//   - Email + network: a few wrong passwords for one email from one IP lock
//     that email from that IP only, longer each time. The real owner on another
//     network (their phone, their office) can still log in.
//   - Email, all networks: many wrong passwords for one email from anywhere
//     pause that email and warn its owner by email. Stops spread-out attacks.
type LoginPolicy struct {
	IPMaxFailures int
	IPWindow      time.Duration

	PairMaxFailures int
	PairWindow      time.Duration
	// PairLocks are the lock lengths for the 1st, 2nd, 3rd… lock; the last repeats.
	PairLocks []time.Duration
	// PairLockMemory is how long earlier locks count towards a longer one.
	PairLockMemory time.Duration

	AccountMaxFailures int
	AccountWindow      time.Duration
	AccountLock        time.Duration
}

// UserLogin applies to business and creator logins.
var UserLogin = LoginPolicy{
	IPMaxFailures:      20,
	IPWindow:           10 * time.Minute,
	PairMaxFailures:    5,
	PairWindow:         15 * time.Minute,
	PairLocks:          []time.Duration{5 * time.Minute, 15 * time.Minute, time.Hour},
	PairLockMemory:     24 * time.Hour,
	AccountMaxFailures: 30,
	AccountWindow:      time.Hour,
	AccountLock:        15 * time.Minute,
}

// AdminLogin is stricter: admin accounts can see and change everything.
var AdminLogin = LoginPolicy{
	IPMaxFailures:      20,
	IPWindow:           10 * time.Minute,
	PairMaxFailures:    3,
	PairWindow:         15 * time.Minute,
	PairLocks:          []time.Duration{15 * time.Minute, time.Hour, 4 * time.Hour},
	PairLockMemory:     24 * time.Hour,
	AccountMaxFailures: 10,
	AccountWindow:      time.Hour,
	AccountLock:        30 * time.Minute,
}

// PairLock is the length of the n-th lock (n counts from 1).
func (p LoginPolicy) PairLock(n int64) time.Duration {
	if n < 1 {
		n = 1
	}
	if int(n) > len(p.PairLocks) {
		return p.PairLocks[len(p.PairLocks)-1]
	}
	return p.PairLocks[n-1]
}

const (
	// ForgotPasswordMaxPerHour is per email.
	ForgotPasswordMaxPerHour = 3
	// ForgotPasswordPerIPPerHour stops one machine asking for resets for many emails.
	ForgotPasswordPerIPPerHour = 10
	// SignupPerIPPerHour stops one machine creating many accounts.
	SignupPerIPPerHour = 20
)
