package auth

import (
	"errors"
	"time"
)

var ErrInvalidResetToken = errors.New("invalid or expired reset token")

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token           string `json:"token"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type ResetPasswordResponse struct {
	RedirectTo string `json:"redirect_to"`
}

type PasswordResetToken struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}
