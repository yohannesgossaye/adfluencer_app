package chat

import (
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	dto "adinfluence_service_app/internal/constant/model/dto/marketplace"
	"net/url"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	maxMessageRunes  = 4000
	maxMediaURLBytes = 2048

	// Per-user send budget: bursts of up to sendBurst messages, refilled at
	// sendRefillPerSecond. Generous for people, a wall for scripts.
	sendBurst           = 20
	sendRefillPerSecond = 2
)

var allowedMediaTypes = map[string]bool{"": true, "text": true, "image": true, "video": true, "file": true}

// canUseChat reports whether an account may read its conversations at all.
func canUseChat(u *authDto.UserResponse) bool {
	switch model.UserStatus(u.Status) {
	case model.BlockedStatus, model.RemovedStatus, model.SuspendedStatus, model.RejectedStatus:
		return false
	}
	return u.Role == string(model.BusinessRole) || u.Role == string(model.InfluencerRole)
}

// canStartOrSend reports whether an account may open conversations and send
// messages: approved influencers, and businesses that are active or still
// completing their profile (the same rule that lets them run campaigns).
func canStartOrSend(u *authDto.UserResponse) bool {
	if !canUseChat(u) {
		return false
	}
	switch model.Role(u.Role) {
	case model.InfluencerRole:
		return u.Status == string(model.ActiveStatus)
	case model.BusinessRole:
		return u.Status == string(model.ActiveStatus) || u.Status == string(model.PendingProfileStatus)
	}
	return false
}

// isBusinessInfluencerPair enforces that chat only connects a business with a
// creator — never two businesses or two creators.
func isBusinessInfluencerPair(a, b *authDto.UserResponse) bool {
	return (a.Role == string(model.BusinessRole) && b.Role == string(model.InfluencerRole)) ||
		(a.Role == string(model.InfluencerRole) && b.Role == string(model.BusinessRole))
}

func requireUUID(value, field string) error {
	if _, err := uuid.Parse(value); err != nil {
		return dto.Invalid(field + " must be a valid id")
	}
	return nil
}

func validateMessage(req *dto.SendMessageRequest) error {
	req.Body = strings.TrimSpace(req.Body)
	if err := requireUUID(req.ConversationID, "conversation_id"); err != nil {
		return err
	}
	if req.ClientMessageID != nil {
		if err := requireUUID(*req.ClientMessageID, "client_message_id"); err != nil {
			return err
		}
	}
	hasMedia := req.MediaURL != nil && strings.TrimSpace(*req.MediaURL) != ""
	if req.Body == "" && !hasMedia {
		return dto.Invalid("message body or attachment is required")
	}
	if utf8.RuneCountInString(req.Body) > maxMessageRunes {
		return dto.Invalid("message is too long (maximum 4000 characters)")
	}
	if !allowedMediaTypes[req.MediaType] {
		return dto.Invalid("media_type must be one of: text, image, video, file")
	}
	if hasMedia {
		link := strings.TrimSpace(*req.MediaURL)
		u, err := url.Parse(link)
		if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" || len(link) > maxMediaURLBytes {
			return dto.Invalid("media_url must be a valid http(s) link")
		}
		req.MediaURL = &link
	} else {
		req.MediaURL = nil
	}
	return nil
}

// sendLimiter is an in-memory token bucket per user. It protects a single
// instance; a shared (Redis) limiter replaces it when running several.
type sendLimiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
}

type bucket struct {
	tokens float64
	last   time.Time
}

func newSendLimiter() *sendLimiter {
	return &sendLimiter{buckets: make(map[string]*bucket)}
}

func (l *sendLimiter) allow(userID string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	b, ok := l.buckets[userID]
	if !ok {
		b = &bucket{tokens: sendBurst, last: now}
		l.buckets[userID] = b
	}
	b.tokens += now.Sub(b.last).Seconds() * sendRefillPerSecond
	if b.tokens > sendBurst {
		b.tokens = sendBurst
	}
	b.last = now
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}
