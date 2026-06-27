package middleware

import "net/http"

type AuthMiddleware interface {
	AuthenticateToken(next http.Handler) http.Handler
}
