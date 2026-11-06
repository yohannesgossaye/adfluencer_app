package errors

import (
	"fmt"
	"time"
)

// WaitError is a rate-limit error that says how long until trying again is
// allowed. errors.Is still matches the wrapped error, so status codes and
// error types are unchanged; the message tells the person how long to wait.
type WaitError struct {
	Err  error
	Wait time.Duration
}

func (e *WaitError) Error() string {
	switch e.Err {
	case ErrLoginLocked:
		return fmt.Sprintf("Too many failed attempts. Try again in %s, or reset your password.", humanWait(e.Wait))
	case ErrTooManyAttempts:
		return fmt.Sprintf("Too many attempts from your network. Try again in %s.", humanWait(e.Wait))
	case ErrForgotPasswordRateLimit:
		return fmt.Sprintf("Too many reset requests. Try again in %s.", humanWait(e.Wait))
	}
	return fmt.Sprintf("%s; try again in %s", e.Err, humanWait(e.Wait))
}

func (e *WaitError) Unwrap() error { return e.Err }

// RetryAfterSeconds is the value for the Retry-After header (at least 1).
func (e *WaitError) RetryAfterSeconds() int {
	s := int((e.Wait + time.Second - 1) / time.Second)
	if s < 1 {
		return 1
	}
	return s
}

// Wait wraps err with how long to wait.
func Wait(err error, d time.Duration) error { return &WaitError{Err: err, Wait: d} }

// humanWait rounds up: "1 minute", "14 minutes", "4 hours".
func humanWait(d time.Duration) string {
	if d > time.Hour {
		h := int((d + time.Hour - 1) / time.Hour)
		return plural(h, "hour")
	}
	m := int((d + time.Minute - 1) / time.Minute)
	if m < 1 {
		m = 1
	}
	return plural(m, "minute")
}

func plural(n int, unit string) string {
	if n == 1 {
		return "1 " + unit
	}
	return fmt.Sprintf("%d %ss", n, unit)
}
