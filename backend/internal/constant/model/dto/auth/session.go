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
	// RefreshTTLSeconds is the idle timeout this session was issued with.
	RefreshTTLSeconds int
	// StartedAt is when the user logged in; it carries over every refresh so
	// the session can't outlive its maximum age.
	StartedAt time.Time
}

type UserSession struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
	// RefreshTTL is zero for sessions created before lifetimes were stored.
	RefreshTTL time.Duration
	StartedAt  time.Time
}

type SessionMeta struct {
	IPAddress string
	UserAgent string
}

type RefreshTokenResponse struct {
	Tokens     TokenResponse `json:"tokens"`
	RedirectTo string        `json:"redirect_to,omitempty"`
}
