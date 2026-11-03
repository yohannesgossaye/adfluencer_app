package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	customErr "adinfluence_service_app/internal/constant/errors"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"

	"go.uber.org/zap"
)

const (
	contactRateScopeIP     = "contact_ip"
	contactRateScopeEmail  = "contact_email"
	contactMaxPerHourIP    = 5
	contactMaxPerHourEmail = 3
)

// SubmitContact emails a Contact us message to the team inbox. Messages with
// the hidden honeypot field filled are dropped silently, and senders are
// limited per IP and per email address.
func (m *userAuthModule) SubmitContact(ctx context.Context, req authDto.ContactRequest, meta authDto.SessionMeta) error {
	if strings.TrimSpace(req.Website) != "" {
		m.logger.Info(ctx, "contact message dropped by honeypot", zap.String("ip", meta.IPAddress))
		return nil
	}
	req = req.Normalize()
	if m.contactInbox == "" || !m.mail.Enabled() {
		m.logger.Error(ctx, "contact message not sent: no inbox or email service configured")
		return customErr.ErrContactUnavailable
	}

	since := time.Now().UTC().Add(-time.Hour)
	for _, c := range []struct {
		scope, key string
		max        int
	}{{contactRateScopeIP, meta.IPAddress, contactMaxPerHourIP}, {contactRateScopeEmail, req.Email, contactMaxPerHourEmail}} {
		if c.key == "" {
			continue
		}
		n, err := m.persistence.CountRateLimitEvents(ctx, c.scope, c.key, since)
		if err != nil {
			return err
		}
		if n >= c.max {
			m.logger.Warn(ctx, "contact rate limit exceeded", zap.String("scope", c.scope))
			return customErr.ErrContactRateLimit
		}
	}
	for _, c := range [][2]string{{contactRateScopeIP, meta.IPAddress}, {contactRateScopeEmail, req.Email}} {
		if c[1] == "" {
			continue
		}
		if err := m.persistence.InsertRateLimitEvent(ctx, c[0], c[1]); err != nil {
			return err
		}
	}

	subject := "Contact: " + req.Subject
	if req.Subject == "" {
		subject = "Contact message from " + req.Name
	}
	body := fmt.Sprintf("New message from the Adfluence Contact us page.\n\nName: %s\nEmail: %s\nI am: %s\n\n%s\n\n—\nReply to this person at %s.\n",
		req.Name, req.Email, roleLabel(req.Role), req.Message, req.Email)
	m.sendEmail(m.contactInbox, subject, body, "contact_message")
	return nil
}

func roleLabel(role string) string {
	switch role {
	case "business":
		return "a business"
	case "creator":
		return "a creator"
	default:
		return "other"
	}
}
