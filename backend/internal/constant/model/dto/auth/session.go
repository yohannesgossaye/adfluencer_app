package auth

import (
	"time"
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
