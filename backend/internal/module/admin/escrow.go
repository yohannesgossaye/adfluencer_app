package admin

import (
	"context"
	"errors"
	"strings"

	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
	"adinfluence_service_app/platform/utils"

	"go.uber.org/zap"
)

// ListEscrowQueue lists accepted contracts whose payment is not confirmed yet.
func (m *adminModule) ListEscrowQueue(ctx context.Context) (*adminDto.EscrowQueueResponse, error) {
	items, err := m.contracts.ListContractsAwaitingFunding(ctx)
	if err != nil {
		m.logger.Error(ctx, "failed to list contracts awaiting funding", zap.String("reason", err.Error()))
		return nil, err
	}
	return &adminDto.EscrowQueueResponse{Items: items}, nil
}

// FundContractEscrow is how escrow is funded until real payments exist: an
// admin confirms the business's payment arrived, and the contract starts.
func (m *adminModule) FundContractEscrow(ctx context.Context, adminUserID, contractID string, req adminDto.FundEscrowRequest, meta authDto.SessionMeta) (*adminDto.FundEscrowResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	reference := strings.TrimSpace(req.Reference)
	paymentID := utils.GenerateUUID()
	payment, err := m.contracts.FundEscrow(ctx, contractID, paymentID, adminUserID, reference)
	if err != nil {
		m.logger.Warn(ctx, "escrow funding not applied", zap.String("contract_id", contractID), zap.String("reason", err.Error()))
		switch {
		case errors.Is(err, contractDto.ErrContractNotFound):
			return nil, customErr.ErrContractNotFound
		case errors.Is(err, contractDto.ErrInvalidContractStatus):
			return nil, customErr.ErrInvalidContractStatus
		case errors.Is(err, contractDto.ErrPaymentAlreadyExists):
			return nil, customErr.ErrPaymentAlreadyExists
		}
		return nil, err
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID: &adminUserID,
		Action:      string(model.PaymentStatusChangeAction),
		EntityType:  "payment",
		EntityID:    &paymentID,
		Description: "admin confirmed contract payment (escrow held): " + reference,
		IPAddress:   meta.IPAddress,
		UserAgent:   meta.UserAgent,
	})
	m.logger.Info(ctx, "escrow funded by admin", zap.String("contract_id", contractID), zap.String("admin_id", adminUserID))

	return &adminDto.FundEscrowResponse{ContractID: contractID, Status: "active", Payment: *payment}, nil
}
