package business

import (
	"context"

	supportDto "adinfluence_service_app/internal/constant/model/dto/support"
)

// ReplySupportTicket adds the user's message to their own ticket; a finished
// ticket is opened again.
func (m *businessModule) ReplySupportTicket(ctx context.Context, userID, ticketID string, req supportDto.ReplyRequest) (*supportDto.ReplyResult, error) {
	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}
	return m.threads.Reply(ctx, userID, ticketID, req)
}
