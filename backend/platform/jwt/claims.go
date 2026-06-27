package jwt

import "github.com/golang-jwt/jwt/v5"

// AccessClaims is embedded in the access JWT.
// Middleware reads these claims after signature validation.
type AccessClaims struct {
	UserID      string `json:"sub"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	Status      string `json:"status"`
	ProfileID   string `json:"profile_id"`
	ProfileName string `json:"profile_name"`
	jwt.RegisteredClaims
}
