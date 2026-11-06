package utils

import (
	"net"
	"net/http"
	"net/netip"
	"strings"
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

// DefaultTrustedProxies are the addresses a forwarding proxy can connect
// from: this machine (nginx) and private networks (Docker, the cloud network).
const DefaultTrustedProxies = "127.0.0.0/8,::1/128,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16,fc00::/7"

// ParseTrustedProxies reads a comma-separated list of CIDRs, skipping bad entries.
func ParseTrustedProxies(list string) []netip.Prefix {
	var out []netip.Prefix
	for _, s := range strings.Split(list, ",") {
		if p, err := netip.ParsePrefix(strings.TrimSpace(s)); err == nil {
			out = append(out, p.Masked())
		}
	}
	return out
}

// TrustedRealIP sets the request's address to the visitor's real IP, but
// believes X-Real-IP / X-Forwarded-For only when the request comes from a
// trusted proxy. Anyone can send those headers; trusting them from anywhere
// would let an attacker pick a new IP for every request and dodge the
// per-IP limits.
func TrustedRealIP(trusted []netip.Prefix) func(http.Handler) http.Handler {
	isTrusted := func(a netip.Addr) bool {
		a = a.Unmap()
		for _, p := range trusted {
			if p.Contains(a) {
				return true
			}
		}
		return false
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if ip := forwardedClientIP(r, isTrusted); ip.IsValid() {
				r.RemoteAddr = net.JoinHostPort(ip.String(), "0")
			}
			next.ServeHTTP(w, r)
		})
	}
}

func forwardedClientIP(r *http.Request, isTrusted func(netip.Addr) bool) netip.Addr {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	peer, err := netip.ParseAddr(host)
	if err != nil || !isTrusted(peer) {
		return netip.Addr{}
	}
	// nginx sets X-Real-IP to the address it saw, replacing whatever the visitor sent.
	if ip, err := netip.ParseAddr(strings.TrimSpace(r.Header.Get("X-Real-IP"))); err == nil {
		return ip.Unmap()
	}
	// Otherwise take the last X-Forwarded-For entry not added by our own proxies;
	// entries further left were written by the visitor and can be fake.
	parts := strings.Split(strings.Join(r.Header.Values("X-Forwarded-For"), ","), ",")
	for i := len(parts) - 1; i >= 0; i-- {
		ip, err := netip.ParseAddr(strings.TrimSpace(parts[i]))
		if err != nil {
			return netip.Addr{}
		}
		if !isTrusted(ip) {
			return ip.Unmap()
		}
	}
	return netip.Addr{}
}
