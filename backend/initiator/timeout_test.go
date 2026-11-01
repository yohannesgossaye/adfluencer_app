package initiator

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRequestTimeout(t *testing.T) {
	cases := []struct {
		name    string
		headers map[string]string
		want    time.Duration // 0 = no deadline
	}{
		{"json request", map[string]string{"Content-Type": "application/json"}, time.Second},
		{"upload", map[string]string{"Content-Type": "multipart/form-data; boundary=x"}, time.Minute},
		{"websocket", map[string]string{"Upgrade": "websocket"}, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var deadline time.Time
			var has bool
			h := requestTimeout(time.Second, time.Minute)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				deadline, has = r.Context().Deadline()
			}))
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}
			start := time.Now()
			h.ServeHTTP(httptest.NewRecorder(), req)
			if tc.want == 0 {
				if has {
					t.Fatal("websocket request got a deadline")
				}
				return
			}
			if !has {
				t.Fatal("no deadline set")
			}
			if got := deadline.Sub(start); (got - tc.want).Abs() > 100*time.Millisecond {
				t.Fatalf("deadline in %v, want about %v", got, tc.want)
			}
		})
	}
}
