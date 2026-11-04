package influencer

import (
	"context"

	supportDto "adinfluence_service_app/internal/constant/model/dto/support"
)

// ReplySupportTicket adds the user's message to their own ticket; a finished
// ticket is opened again.
func (m *influencerModule) ReplySupportTicket(ctx context.Context, userID, ticketID string, req supportDto.ReplyRequest) (*supportDto.ReplyResult, error) {
	if err := m.requireSupportAccess(ctx, userID); err != nil {
		return nil, err
	}
	return m.threads.Reply(ctx, userID, ticketID, req)
}
