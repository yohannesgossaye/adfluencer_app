package jwt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ManagerConfig struct {
	SecretKey      string
	AccessTokenTTL time.Duration
	Sessions       SessionPolicy
}

// SessionPolicy (bank-style): a session ends after Idle without activity
// (each refresh renews it) and after MaxAge in any case. Admins get their
// own, stricter values.
type SessionPolicy struct {
	Idle        time.Duration
	MaxAge      time.Duration
	AdminIdle   time.Duration
	AdminMaxAge time.Duration
}

type TokenManager struct {
	cfg ManagerConfig
}

func NewTokenManager(cfg ManagerConfig) *TokenManager {
	return &TokenManager{cfg: cfg}
}

type AccessTokenInput struct {
	UserID      string
	Email       string
	Role        string
	Status      string
	ProfileID   string
	ProfileName string
}

type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresIn  int64
	RefreshExpiresAt time.Time
}

func (m *TokenManager) GenerateTokenPair(input AccessTokenInput) (TokenPair, error) {
	accessToken, accessExpiresIn, err := m.GenerateAccessToken(input)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresIn:  accessExpiresIn,
		RefreshExpiresAt: time.Now().UTC().Add(m.cfg.Sessions.Idle),
	}, nil
}

func (m *TokenManager) GenerateAccessToken(input AccessTokenInput) (string, int64, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(m.cfg.AccessTokenTTL)

	claims := AccessClaims{
		UserID:      input.UserID,
		Email:       input.Email,
		Role:        input.Role,
		Status:      input.Status,
		ProfileID:   input.ProfileID,
		ProfileName: input.ProfileName,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   input.UserID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(m.cfg.SecretKey))
	if err != nil {
		return "", 0, fmt.Errorf("sign access token: %w", err)
	}

	return signed, int64(m.cfg.AccessTokenTTL.Seconds()), nil
}

func (m *TokenManager) ParseAccessToken(tokenString string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(m.cfg.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid access token")
	}

	return claims, nil
}

// SessionLimits returns how long a session of this role may sit idle and how
// long it may last at all.
func (m *TokenManager) SessionLimits(role string) (idle, maxAge time.Duration) {
	p := m.cfg.Sessions
	idle, maxAge = p.Idle, p.MaxAge
	if role == "admin" {
		idle, maxAge = p.AdminIdle, p.AdminMaxAge
	}
	if idle <= 0 {
		idle = 15 * time.Minute
	}
	if maxAge <= 0 {
		maxAge = 12 * time.Hour
	}
	return idle, maxAge
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
