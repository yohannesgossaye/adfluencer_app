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
	// RefreshTTLSeconds is the session lifetime chosen at login; it is reused
	// on every rotation so "remember me" sessions keep their longer lifetime.
	RefreshTTLSeconds int
}

type UserSession struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
	// RefreshTTL is zero for sessions created before lifetimes were stored.
	RefreshTTL time.Duration
}

type SessionMeta struct {
	IPAddress string
	UserAgent string
}

type RefreshTokenResponse struct {
	Tokens     TokenResponse `json:"tokens"`
	RedirectTo string        `json:"redirect_to,omitempty"`
}
