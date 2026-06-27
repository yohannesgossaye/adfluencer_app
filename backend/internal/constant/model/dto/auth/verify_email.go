package auth

import (
	"errors"
	"time"
)

var ErrInvalidVerifyOTP = errors.New("invalid or expired verification code")

type VerifyEmailRequest struct {
	OTP string `json:"otp"`
}

type VerifyEmailResponse struct {
	User VerifyEmailUserResponse `json:"user"`
}

type VerifyEmailUserResponse struct {
	ID            string `json:"id"`
	EmailVerified bool   `json:"email_verified"`
}

type EmailVerificationToken struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}
