// Package secrets loads the API's configuration from HashiCorp Vault.
//
// The server keeps only three values locally: VAULT_ADDR, VAULT_TOKEN (a
// read-only token) and VAULT_PATH. Everything else (secrets and settings) lives
// in one Vault KV v2 entry whose keys are the usual variable names
// (JWT_SECRET, DATABASE_URL, …). At startup those keys are applied to the
// process environment, so every existing setting works unchanged.
//
// Values are never logged; only key names are.
package secrets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Vault reads one KV v2 entry and keeps its token alive.
type Vault struct {
	addr, token, mount, path string
	http                     *http.Client
}

// FromEnv returns nil when VAULT_ADDR is unset (local runs and CI use .env).
func FromEnv() (*Vault, error) {
	addr := strings.TrimRight(strings.TrimSpace(os.Getenv("VAULT_ADDR")), "/")
	if addr == "" {
		return nil, nil
	}
	token := strings.TrimSpace(os.Getenv("VAULT_TOKEN"))
	full := strings.Trim(strings.TrimSpace(os.Getenv("VAULT_PATH")), "/")
	mount, path, ok := strings.Cut(full, "/")
	if token == "" || !ok || mount == "" || path == "" {
		return nil, errors.New("VAULT_ADDR is set, so VAULT_TOKEN and VAULT_PATH (e.g. secret/adfluence/prod) are required")
	}
	return &Vault{addr: addr, token: token, mount: mount, path: path, http: &http.Client{Timeout: 10 * time.Second}}, nil
}

// Path is the entry being read, for log messages.
func (v *Vault) Path() string { return v.mount + "/" + v.path }

var keyPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// Read returns every key of the entry as text.
func (v *Vault) Read(ctx context.Context) (map[string]string, error) {
	var body struct {
		Data struct {
			Data map[string]any `json:"data"`
		} `json:"data"`
	}
	if err := v.call(ctx, http.MethodGet, "/v1/"+v.mount+"/data/"+v.path, &body); err != nil {
		return nil, err
	}
	if len(body.Data.Data) == 0 {
		return nil, fmt.Errorf("vault entry %s is empty or missing", v.Path())
	}
	out := make(map[string]string, len(body.Data.Data))
	for k, raw := range body.Data.Data {
		switch val := raw.(type) {
		case string:
			// A newline pasted at the end of a secret is never meant.
			out[k] = strings.TrimRight(val, "\r\n")
		case bool, float64, json.Number:
			out[k] = fmt.Sprint(val)
		case nil:
			out[k] = ""
		default:
			return nil, fmt.Errorf("vault key %s must be text, a number or true/false", k)
		}
	}
	return out, nil
}

// Apply reads the entry and sets each key in the process environment,
// replacing .env values: Vault is the source of truth. It returns the key
// names applied (never the values).
func (v *Vault) Apply(ctx context.Context) ([]string, error) {
	values, err := v.Read(ctx)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(values))
	for k, val := range values {
		if !keyPattern.MatchString(k) {
			return nil, fmt.Errorf("vault key %q is not a valid name (use CAPITALS_AND_UNDERSCORES)", k)
		}
		if strings.HasPrefix(k, "VAULT_") {
			return nil, fmt.Errorf("vault key %s is not allowed: Vault's own settings stay on the server", k)
		}
		if err := os.Setenv(k, val); err != nil {
			return nil, err
		}
		names = append(names, k)
	}
	sort.Strings(names)
	return names, nil
}

// RenewToken extends the token's lifetime (a periodic token gets its full
// period again).
func (v *Vault) RenewToken(ctx context.Context) error {
	return v.call(ctx, http.MethodPost, "/v1/auth/token/renew-self", nil)
}

// KeepTokenAlive renews the token now and then every day until ctx ends, so a
// running API never loses access. Failures are reported, not fatal: the
// configuration is already loaded.
func (v *Vault) KeepTokenAlive(ctx context.Context, every time.Duration, report func(error)) {
	renew := func() {
		if err := v.RenewToken(ctx); err != nil && ctx.Err() == nil {
			report(err)
		}
	}
	renew()
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			renew()
		}
	}
}

func (v *Vault) call(ctx context.Context, method, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, method, v.addr+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Vault-Token", v.token)
	res, err := v.http.Do(req)
	if err != nil {
		return fmt.Errorf("vault unreachable: %w", err)
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	switch {
	case res.StatusCode == http.StatusForbidden:
		return errors.New("vault refused the token (wrong, expired or revoked token, or no permission for this path)")
	case res.StatusCode == http.StatusNotFound:
		return fmt.Errorf("vault entry %s not found", v.Path())
	case res.StatusCode == http.StatusServiceUnavailable:
		return errors.New("vault is sealed: unseal it, then the API starts by itself")
	case res.StatusCode >= 300:
		return fmt.Errorf("vault answered %d for %s", res.StatusCode, path)
	}
	if out == nil {
		return nil
	}
	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.UseNumber()
	return dec.Decode(out)
}
