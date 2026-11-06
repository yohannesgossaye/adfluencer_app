package secrets

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

const goodToken = "hvs.test-read-only"

// fakeVault answers like Vault's KV v2 and token endpoints.
func fakeVault(t *testing.T, entry string, sealed bool) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	renewals := &atomic.Int32{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sealed {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		if r.Header.Get("X-Vault-Token") != goodToken {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/secret/data/adfluence/prod":
			if entry == "" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Write([]byte(`{"data":{"data":` + entry + `,"metadata":{"version":3}}}`))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/auth/token/renew-self":
			renewals.Add(1)
			w.Write([]byte(`{"auth":{"lease_duration":2764800}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)
	return srv, renewals
}

func vaultFor(t *testing.T, addr, token string) *Vault {
	t.Helper()
	t.Setenv("VAULT_ADDR", addr)
	t.Setenv("VAULT_TOKEN", token)
	t.Setenv("VAULT_PATH", "secret/adfluence/prod")
	v, err := FromEnv()
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func TestApplyLoadsEveryKeyAndWinsOverEnv(t *testing.T) {
	srv, _ := fakeVault(t, `{"JWT_SECRET":"from-vault-0123456789abcdef0123456789\n","SMTP_PORT":587,"S3_USE_PATH_STYLE":true,"FRONTEND_URL":"https://adfluence.example"}`, false)
	v := vaultFor(t, srv.URL, goodToken)
	t.Setenv("JWT_SECRET", "old-value-from-dotenv")
	t.Setenv("SMTP_PORT", "")
	t.Setenv("S3_USE_PATH_STYLE", "")
	t.Setenv("FRONTEND_URL", "")

	names, err := v.Apply(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(names, ","); got != "FRONTEND_URL,JWT_SECRET,S3_USE_PATH_STYLE,SMTP_PORT" {
		t.Fatalf("applied %s", got)
	}
	for k, want := range map[string]string{
		"JWT_SECRET": "from-vault-0123456789abcdef0123456789", "SMTP_PORT": "587",
		"S3_USE_PATH_STYLE": "true", "FRONTEND_URL": "https://adfluence.example",
	} {
		if got := os.Getenv(k); got != want {
			t.Fatalf("%s = %q, want %q", k, got, want)
		}
	}
}

func TestApplyRefusesClearly(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name, entry, token, want string
		sealed                   bool
	}{
		{"wrong token", `{"A":"1"}`, "hvs.stolen-or-revoked", "refused the token", false},
		{"sealed", `{"A":"1"}`, goodToken, "sealed", true},
		{"missing entry", "", goodToken, "not found", false},
		{"empty entry", `{}`, goodToken, "empty", false},
		{"bad key name", `{"jwt secret":"x"}`, goodToken, "not a valid name", false},
		{"vault settings cannot be overridden", `{"VAULT_TOKEN":"x"}`, goodToken, "not allowed", false},
		{"nested value", `{"A":{"b":1}}`, goodToken, "must be text", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			srv, _ := fakeVault(t, c.entry, c.sealed)
			_, err := vaultFor(t, srv.URL, c.token).Apply(ctx)
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("err = %v, want it to mention %q", err, c.want)
			}
		})
	}

	t.Run("unreachable", func(t *testing.T) {
		_, err := vaultFor(t, "http://127.0.0.1:1", goodToken).Apply(ctx)
		if err == nil || !strings.Contains(err.Error(), "unreachable") {
			t.Fatalf("err = %v", err)
		}
	})
}

func TestErrorsNeverContainValues(t *testing.T) {
	srv, _ := fakeVault(t, `{"JWT_SECRET":"super-secret-value","BAD":{"x":"super-secret-value"}}`, false)
	_, err := vaultFor(t, srv.URL, goodToken).Apply(context.Background())
	if err == nil || strings.Contains(err.Error(), "super-secret-value") || strings.Contains(err.Error(), goodToken) {
		t.Fatalf("error leaks a value: %v", err)
	}
}

func TestFromEnv(t *testing.T) {
	t.Setenv("VAULT_ADDR", "")
	if v, err := FromEnv(); v != nil || err != nil {
		t.Fatalf("without VAULT_ADDR: %v, %v; want nil, nil (use .env)", v, err)
	}
	t.Setenv("VAULT_ADDR", "http://vault:8200")
	t.Setenv("VAULT_TOKEN", goodToken)
	t.Setenv("VAULT_PATH", "adfluence-without-mount")
	if _, err := FromEnv(); err == nil {
		t.Fatal("a path without a mount was accepted")
	}
}

func TestKeepTokenAliveRenews(t *testing.T) {
	srv, renewals := fakeVault(t, `{"A":"1"}`, false)
	v := vaultFor(t, srv.URL, goodToken)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		v.KeepTokenAlive(ctx, 20*time.Millisecond, func(err error) { t.Errorf("renew failed: %v", err) })
		close(done)
	}()
	time.Sleep(70 * time.Millisecond)
	cancel()
	<-done
	if n := renewals.Load(); n < 3 {
		t.Fatalf("renewed %d times, want at start and on every tick", n)
	}
}
