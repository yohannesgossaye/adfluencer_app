package business

import (
	"context"

	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
)

// OpenDispute is "Report a problem" on one of the business's contracts.
func (m *businessModule) OpenDispute(ctx context.Context, userID, contractID string, req contractDto.OpenDisputeRequest) (*contractDto.DisputeSummary, error) {
	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}
	return m.disputes.Open(ctx, userID, "business", contractID, req)
}
