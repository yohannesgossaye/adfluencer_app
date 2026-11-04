// Package support is the conversation on a support ticket, shared by the
// business, creator and admin portals: the messages, the user's replies and
// the email a user gets when our team answers.
package support

import (
	"context"
	"errors"
	"fmt"
	"strings"

	customErr "adinfluence_service_app/internal/constant/errors"
	supportDto "adinfluence_service_app/internal/constant/model/dto/support"
	"adinfluence_service_app/internal/storage"
	supportPersist "adinfluence_service_app/internal/storage/persist/support"
	"adinfluence_service_app/platform/notify"
	"adinfluence_service_app/platform/utils"

	"go.uber.org/zap"
)

type Service struct {
	repo  storage.SupportPersistence
	email *notify.Emailer
	log   utils.Logger
}

func New(repo storage.SupportPersistence, email *notify.Emailer, log utils.Logger) *Service {
	return &Service{repo: repo, email: email, log: log}
}

// Messages is the conversation after the opening message, oldest first.
func (s *Service) Messages(ctx context.Context, ticketID string) ([]supportDto.Message, error) {
	return s.repo.ListMessages(ctx, ticketID)
}

// Reply adds the ticket owner's message. Replying to a resolved or closed
// ticket opens it again, so our team sees it in their queue.
func (s *Service) Reply(ctx context.Context, userID, ticketID string, req supportDto.ReplyRequest) (*supportDto.ReplyResult, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	msg, status, reopened, err := s.repo.AddUserMessage(ctx, userID, ticketID, req.Message)
	if err != nil {
		if errors.Is(err, supportPersist.ErrTicketNotFound) {
			return nil, customErr.ErrSupportTicketNotFound
		}
		return nil, err
	}
	s.log.Info(ctx, "support ticket user reply", zap.String("ticket_id", ticketID), zap.Bool("reopened", reopened))
	return &supportDto.ReplyResult{Message: msg, Status: status, Reopened: reopened}, nil
}

// NotifyReply emails the ticket owner that our team answered. It never
// fails the reply itself: a missing ticket or email only skips the email.
func (s *Service) NotifyReply(ctx context.Context, ticketID, reply string) {
	if s.email == nil {
		return
	}
	t, err := s.repo.GetAdminTicket(ctx, ticketID)
	if err != nil {
		s.log.Warn(ctx, "support reply email skipped", zap.String("ticket_id", ticketID), zap.String("reason", err.Error()))
		return
	}
	portal := "influencer"
	if t.SenderRole == "business" {
		portal = "business"
	}
	link := fmt.Sprintf("%s/%s/support?ticket_id=%s", strings.TrimRight(s.email.FrontendURL, "/"), portal, ticketID)
	body := fmt.Sprintf("Hello %s,\n\nOur support team replied to your request \"%s\":\n\n%s\n\nPlease answer in the conversation, so everything stays with your request:\n%s\n\nThe Adfluence team\n",
		t.SenderName, t.Subject, indent(strings.TrimSpace(reply)), link)
	s.email.Send(t.SenderEmail, "Adfluence support replied: "+t.Subject, body, "support_reply")
}

func indent(text string) string {
	lines := strings.Split(text, "\n")
	for i, l := range lines {
		lines[i] = "  " + l
	}
	return strings.Join(lines, "\n")
}
