package notify

import (
	"adinfluence_service_app/platform/mail"

	"go.uber.org/zap"
)

// Emailer sends emails in the background through the dispatcher. It is
// safe to use when email is not configured: sending is then skipped.
type Emailer struct {
	mail        *mail.Client
	dispatcher  *Dispatcher
	FrontendURL string
}

func NewEmailer(m *mail.Client, d *Dispatcher, frontendURL string) *Emailer {
	return &Emailer{mail: m, dispatcher: d, FrontendURL: frontendURL}
}

// Send queues one email. purpose is only for the logs.
func (e *Emailer) Send(to, subject, body, purpose string) {
	if e == nil || e.mail == nil || e.dispatcher == nil || !e.mail.Enabled() || to == "" {
		return
	}
	e.dispatcher.Submit(Job{
		Kind:   "email",
		Fields: []zap.Field{zap.String("purpose", purpose)},
		Send:   func() error { return e.mail.Send(to, subject, body) },
	})
}
