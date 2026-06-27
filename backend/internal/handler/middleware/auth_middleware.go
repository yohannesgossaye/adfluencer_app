package middleware

import (
	"adinfluence_service_app/internal/constant/model/response"
	"adinfluence_service_app/platform/jwt"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const AccessClaimsKey contextKey = "access_claims"

type authMiddleware struct {
	tokens *jwt.TokenManager
}

func InitAuthMiddleware(tokens *jwt.TokenManager) AuthMiddleware {
	return &authMiddleware{tokens: tokens}
}

func (m *authMiddleware) AuthenticateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", nil, "missing authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", nil, "invalid authorization header")
			return
		}

		claims, err := m.tokens.ParseAccessToken(parts[1])
		if err != nil {
			response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", err, "invalid or expired access token")
			return
		}

		ctx := context.WithValue(r.Context(), AccessClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetAccessClaims(ctx context.Context) (*jwt.AccessClaims, bool) {
	claims, ok := ctx.Value(AccessClaimsKey).(*jwt.AccessClaims)
	return claims, ok
}
