package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrustedRealIP(t *testing.T) {
	mw := TrustedRealIP(ParseTrustedProxies(DefaultTrustedProxies))
	cases := []struct {
		name, remote string
		headers      map[string]string
		want         string
	}{
		{"direct visitor", "102.213.71.236:5000", nil, "102.213.71.236"},
		{"direct visitor faking headers is ignored", "102.213.71.236:5000",
			map[string]string{"X-Real-IP": "1.2.3.4", "X-Forwarded-For": "1.2.3.4", "True-Client-IP": "1.2.3.4"}, "102.213.71.236"},
		{"through nginx", "172.17.0.1:40000", map[string]string{"X-Real-IP": "102.213.71.236"}, "102.213.71.236"},
		{"through nginx, X-Forwarded-For only: visitor's own entries ignored", "127.0.0.1:40000",
			map[string]string{"X-Forwarded-For": "6.6.6.6, 102.213.71.236"}, "102.213.71.236"},
		{"proxy without headers keeps its address", "172.17.0.1:40000", nil, "172.17.0.1"},
		{"garbage header is ignored", "172.17.0.1:40000", map[string]string{"X-Real-IP": "not-an-ip"}, "172.17.0.1"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got string
			h := mw(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) { got = ClientIPFromRequest(r) }))
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.RemoteAddr = c.remote
			for k, v := range c.headers {
				req.Header.Set(k, v)
			}
			h.ServeHTTP(httptest.NewRecorder(), req)
			if got != c.want {
				t.Fatalf("client IP = %q, want %q", got, c.want)
			}
		})
	}
}
