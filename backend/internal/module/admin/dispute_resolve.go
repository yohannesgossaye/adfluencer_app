package admin

import (
	"context"

	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
)

// ResolveDispute records the admin's decision: continue, complete or terminate the contract.
func (m *adminModule) ResolveDispute(ctx context.Context, adminUserID, disputeID string, req contractDto.ResolveDisputeRequest) error {
	return m.disputes.Resolve(ctx, adminUserID, disputeID, req)
}
