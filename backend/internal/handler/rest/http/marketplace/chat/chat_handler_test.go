package chat

import (
	dto "adinfluence_service_app/internal/constant/model/dto/marketplace"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWSAccessToken(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/chat/ws", nil)
	r.Header.Set("Sec-WebSocket-Protocol", "adfluence.v1, eyJhbGciOi.payload.sig")
	if tok, sub := wsAccessToken(r); tok != "eyJhbGciOi.payload.sig" || !sub {
		t.Fatalf("subprotocol token = %q, %v", tok, sub)
	}

	r = httptest.NewRequest(http.MethodGet, "/api/v1/chat/ws", nil)
	r.Header.Set("Authorization", "Bearer abc.def.ghi")
	if tok, sub := wsAccessToken(r); tok != "abc.def.ghi" || sub {
		t.Fatalf("header token = %q, %v", tok, sub)
	}

	r = httptest.NewRequest(http.MethodGet, "/api/v1/chat/ws?token=abc.def.ghi", nil)
	if tok, _ := wsAccessToken(r); tok != "" {
		t.Fatal("tokens in the URL must not be accepted (they end up in logs)")
	}

	r = httptest.NewRequest(http.MethodGet, "/api/v1/chat/ws", nil)
	r.Header.Set("Sec-WebSocket-Protocol", "abc.def.ghi")
	if tok, _ := wsAccessToken(r); tok != "" {
		t.Fatal("a token without the adfluence subprotocol must be ignored")
	}
}

func TestClassifyChatError(t *testing.T) {
	cases := []struct {
		err    error
		status int
		code   string
	}{
		{dto.Invalid("bad"), http.StatusBadRequest, "invalid_request"},
		{fmt.Errorf("wrapped: %w", dto.ErrNotParticipant), http.StatusForbidden, "not_participant"},
		{dto.ErrConversationNotFound, http.StatusNotFound, "conversation_not_found"},
		{dto.ErrChatNotAllowed, http.StatusForbidden, "chat_not_allowed"},
		{dto.ErrSendRateLimited, http.StatusTooManyRequests, "rate_limited"},
		{errors.New("pq: connection refused"), http.StatusInternalServerError, "internal_server_error"},
	}
	for _, tc := range cases {
		status, code, msg := classifyChatError(tc.err)
		if status != tc.status || code != tc.code {
			t.Errorf("%v -> %d %s, want %d %s", tc.err, status, code, tc.status, tc.code)
		}
		if status == http.StatusInternalServerError && msg != "internal server error" {
			t.Errorf("internal error text leaked: %q", msg)
		}
	}
}
