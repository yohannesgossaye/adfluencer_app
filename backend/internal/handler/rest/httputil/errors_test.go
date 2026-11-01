package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func decode(t *testing.T, rec *httptest.ResponseRecorder) (string, string) {
	t.Helper()
	var body struct {
		Message string `json:"message"`
		Error   struct {
			Type   string `json:"type"`
			Detail []struct {
				Description string `json:"description"`
			} `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	return body.Error.Type, rec.Body.String()
}

func TestWriteServiceErrorHidesInternalErrors(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteServiceError(rec, fmt.Errorf("get campaign: %w", errors.New("dial tcp 10.0.0.5:5432: connection refused")))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", rec.Code)
	}
	_, raw := decode(t, rec)
	if strings.Contains(raw, "10.0.0.5") || strings.Contains(raw, "get campaign") {
		t.Fatalf("internal error leaked: %s", raw)
	}
}

func TestWriteServiceErrorMapsInvalidUUID(t *testing.T) {
	rec := httptest.NewRecorder()
	pgErr := &pgconn.PgError{Code: "22P02", Message: `invalid input syntax for type uuid: "x"`}
	WriteServiceError(rec, fmt.Errorf("get campaign: %w", pgErr))

	code, raw := decode(t, rec)
	if rec.Code != http.StatusBadRequest || code != "invalid_id" {
		t.Fatalf("got %d %s", rec.Code, code)
	}
	if strings.Contains(raw, "uuid") {
		t.Fatalf("database error leaked: %s", raw)
	}
}

func TestWriteServiceErrorMapsUniqueViolation(t *testing.T) {
	rec := httptest.NewRecorder()
	pgErr := &pgconn.PgError{Code: "23505", ConstraintName: "campaign_invites_campaign_id_influencer_user_id_key"}
	WriteServiceError(rec, fmt.Errorf("insert invite: %w", pgErr))

	code, raw := decode(t, rec)
	if rec.Code != http.StatusConflict || code != "conflict" {
		t.Fatalf("got %d %s", rec.Code, code)
	}
	if strings.Contains(raw, "campaign_invites") {
		t.Fatalf("constraint name leaked: %s", raw)
	}
}
