package influencer

import (
	"context"

	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
)

// OpenDispute is "Report a problem" on one of the creator's contracts.
func (m *influencerModule) OpenDispute(ctx context.Context, userID, contractID string, req contractDto.OpenDisputeRequest) (*contractDto.DisputeSummary, error) {
	if err := m.requireActiveInfluencer(ctx, userID); err != nil {
		return nil, err
	}
	return m.disputes.Open(ctx, userID, "influencer", contractID, req)
}
