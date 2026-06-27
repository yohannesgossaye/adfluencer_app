package utils

import (
	"net"
	"net/http"
)

// ClientIPFromRequest returns a PostgreSQL-compatible IP string, or empty if unknown/invalid.
// Invalid or unparseable addresses are ignored (stored as NULL) — never fail auth because of IP.
func ClientIPFromRequest(r *http.Request) string {
	raw := r.RemoteAddr
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		raw = host
	}

	parsed := net.ParseIP(raw)
	if parsed == nil {
		return ""
	}

	return parsed.String()
}
