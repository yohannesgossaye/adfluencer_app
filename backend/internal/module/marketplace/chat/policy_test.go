package chat

import (
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	dto "adinfluence_service_app/internal/constant/model/dto/marketplace"
	"errors"
	"strings"
	"testing"
	"time"
)

func user(role, status string) *authDto.UserResponse {
	return &authDto.UserResponse{Role: role, Status: status}
}

func TestCanStartOrSend(t *testing.T) {
	cases := []struct {
		role, status string
		want         bool
	}{
		{"influencer", "active", true},
		{"influencer", "pending_approval", false},
		{"influencer", "pending_profile", false},
		{"business", "active", true},
		{"business", "pending_profile", true},
		{"business", "blocked", false},
		{"influencer", "removed", false},
		{"business", "suspended", false},
		{"admin", "active", false},
	}
	for _, tc := range cases {
		if got := canStartOrSend(user(tc.role, tc.status)); got != tc.want {
			t.Errorf("canStartOrSend(%s,%s) = %v, want %v", tc.role, tc.status, got, tc.want)
		}
	}
}

func TestIsBusinessInfluencerPair(t *testing.T) {
	b, i := user("business", "active"), user("influencer", "active")
	if !isBusinessInfluencerPair(b, i) || !isBusinessInfluencerPair(i, b) {
		t.Fatal("business and influencer must be allowed to chat")
	}
	if isBusinessInfluencerPair(i, user("influencer", "active")) {
		t.Fatal("influencer to influencer must be rejected")
	}
	if isBusinessInfluencerPair(b, user("business", "active")) {
		t.Fatal("business to business must be rejected")
	}
}

func TestValidateMessage(t *testing.T) {
	conv := "7c1c7d3e-8c0e-4c8a-9a3f-2b1f0e4b5a6d"
	link := func(s string) *string { return &s }

	ok := []dto.SendMessageRequest{
		{ConversationID: conv, Body: "hello"},
		{ConversationID: conv, Body: strings.Repeat("ሰ", maxMessageRunes)},
		{ConversationID: conv, MediaURL: link("https://cdn.example.com/a.png"), MediaType: "image"},
	}
	for i, req := range ok {
		r := req
		if err := validateMessage(&r); err != nil {
			t.Errorf("case %d: unexpected error %v", i, err)
		}
	}

	bad := []dto.SendMessageRequest{
		{ConversationID: "not-a-uuid", Body: "hi"},
		{ConversationID: conv, Body: "   "},
		{ConversationID: conv, Body: strings.Repeat("x", maxMessageRunes+1)},
		{ConversationID: conv, Body: "hi", MediaType: "exe"},
		{ConversationID: conv, MediaURL: link("javascript:alert(1)")},
		{ConversationID: conv, Body: "hi", ClientMessageID: link("123")},
	}
	for i, req := range bad {
		r := req
		var v *dto.ValidationError
		if err := validateMessage(&r); !errors.As(err, &v) {
			t.Errorf("case %d: expected validation error, got %v", i, err)
		}
	}
}

func TestSendLimiter(t *testing.T) {
	l := newSendLimiter()
	now := time.Unix(1_000_000, 0)
	for i := 0; i < sendBurst; i++ {
		if !l.allow("u1", now) {
			t.Fatalf("message %d within burst was rejected", i+1)
		}
	}
	if l.allow("u1", now) {
		t.Fatal("message beyond burst should be rejected")
	}
	if !l.allow("u2", now) {
		t.Fatal("limits must be per user")
	}
	if !l.allow("u1", now.Add(time.Second)) {
		t.Fatal("tokens should refill over time")
	}
}
