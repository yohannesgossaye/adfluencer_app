package auth

import (
	"errors"
	"time"
)

var (
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
	ErrUserNotFound        = errors.New("user not found")
	ErrProfileNotFound     = errors.New("profile not found")
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type CreateSessionInput struct {
	UserID           string
	RefreshTokenHash string
	IPAddress        string
	UserAgent        string
	ExpiresAt        time.Time
}

type UserSession struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}

type SessionMeta struct {
	IPAddress string
	UserAgent string
}

type RefreshTokenResponse struct {
	Tokens     TokenResponse `json:"tokens"`
	RedirectTo string        `json:"redirect_to,omitempty"`
}
