package initiator

import (
	"context"
	"net/http"
	"strings"
	"time"
)

const (
	requestTimeoutDefault = 30 * time.Second
	requestTimeoutUpload  = 2 * time.Minute
)

// requestTimeout puts a deadline on each request's context so a slow query
// is cancelled instead of holding a DB connection indefinitely. Multipart
// uploads get longer (mobile networks), and WebSocket upgrades get none
// because the chat connection is meant to stay open.
func requestTimeout(normal, upload time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
				next.ServeHTTP(w, r)
				return
			}
			d := normal
			if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/") {
				d = upload
			}
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
