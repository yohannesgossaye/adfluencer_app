package business

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
	"adinfluence_service_app/internal/module/business/helper"
	"context"
	"errors"
	"strings"

	"adinfluence_service_app/platform/utils"
)

func (m *businessModule) CreateContract(ctx context.Context, userID string, req businessDto.CreateContractRequest, meta authDto.SessionMeta) (*commonDto.CreatedIDResponse, error) {
	if err := m.requireVerifiedBusiness(ctx, userID); err != nil {
		return nil, err
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := validateContractDeliverables(req.Deliverables); err != nil {
		return nil, err
	}

	invite, err := m.contracts.GetAcceptedInvite(ctx, userID, req.InviteID)
	if err != nil {
		return nil, mapContractError(err)
	}
	if invite.CampaignStatus != "active" {
		return nil, customErr.ErrInvalidCampaignStatus
	}

	exists, err := m.contracts.ContractExistsForInvite(ctx, req.InviteID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, customErr.ErrContractAlreadyExists
	}

	startDate, err := businessDto.ParseContractDate(req.StartDate)
	if err != nil {
		return nil, customErr.ErrBadRequest
	}
	endDate, err := businessDto.ParseContractDate(req.EndDate)
	if err != nil {
		return nil, customErr.ErrBadRequest
	}

	fees := helper.ComputeContractFees(req.AmountETB)
	contractID := utils.GenerateUUID()

	var description *string
	if v := strings.TrimSpace(req.Description); v != "" {
		description = &v
	}

	if err := m.contracts.CreateContract(ctx, contractDto.SaveContractInput{
		ID:               contractID,
		InviteID:         req.InviteID,
		CampaignID:       invite.CampaignID,
		BusinessUserID:   userID,
		InfluencerUserID: invite.InfluencerUserID,
		Title:            strings.TrimSpace(req.Title),
		Description:      description,
		AmountETB:        req.AmountETB,
		Fees:             fees,
		StartDate:        startDate,
		EndDate:          endDate,
		Deliverables:     req.Deliverables,
	}); err != nil {
		return nil, mapContractError(err)
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &invite.InfluencerUserID,
		Action:       string(model.ProfileUpdateAction),
		EntityType:   "contract",
		EntityID:     &contractID,
		Description:  "business created contract",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	return &commonDto.CreatedIDResponse{ID: contractID}, nil
}

func (m *businessModule) ListContracts(ctx context.Context, userID string, query contractDto.ListQuery) (*businessDto.ListContractsResponse, error) {
	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}

	query = query.Normalize()
	items, total, err := m.contracts.ListContractsForBusiness(ctx, userID, query)
	if err != nil {
		return nil, err
	}

	return &businessDto.ListContractsResponse{
		Items: items,
		Pagination: commonDto.NewListMeta(query.Page, query.Limit, total),
	}, nil
}

func (m *businessModule) GetContract(ctx context.Context, userID, contractID string) (*contractDto.DetailResponse, error) {
	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}

	rec, err := m.contracts.GetContractForBusiness(ctx, userID, contractID)
	if err != nil {
		return nil, mapContractError(err)
	}
	detail := rec.DetailResponse
	return &detail, nil
}

func (m *businessModule) PayContract(ctx context.Context, userID, contractID string, meta authDto.SessionMeta) (*businessDto.PayContractResponse, error) {
	if err := m.requireVerifiedBusiness(ctx, userID); err != nil {
		return nil, err
	}

	paymentID := utils.GenerateUUID()
	payment, err := m.contracts.StubPayContract(ctx, userID, contractID, paymentID)
	if err != nil {
		return nil, mapContractError(err)
	}

	rec, err := m.contracts.GetContractForBusiness(ctx, userID, contractID)
	if err != nil {
		return nil, err
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &rec.InfluencerUserID,
		Action:       string(model.PaymentStatusChangeAction),
		EntityType:   "payment",
		EntityID:     &paymentID,
		Description:  "business paid contract (escrow stub)",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	return &businessDto.PayContractResponse{
		Contract: rec.DetailResponse,
		Payment:  *payment,
	}, nil
}

func (m *businessModule) ApproveWork(ctx context.Context, userID, contractID string, meta authDto.SessionMeta) (*businessDto.ApproveWorkResponse, error) {
	if err := m.requireVerifiedBusiness(ctx, userID); err != nil {
		return nil, err
	}

	payoutID := utils.GenerateUUID()
	payout, err := m.contracts.ApproveWork(ctx, userID, contractID, payoutID)
	if err != nil {
		return nil, mapContractError(err)
	}

	rec, err := m.contracts.GetContractForBusiness(ctx, userID, contractID)
	if err != nil {
		return nil, err
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &rec.InfluencerUserID,
		Action:       string(model.ContractStatusChangeAction),
		EntityType:   "payout",
		EntityID:     &payoutID,
		Description:  "business approved work — payout queued",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	return &businessDto.ApproveWorkResponse{
		Contract: rec.DetailResponse,
		Payout:   *payout,
	}, nil
}

func (m *businessModule) CancelContract(ctx context.Context, userID, contractID string, meta authDto.SessionMeta) (*businessDto.ContractMutationResponse, error) {
	if err := m.requireVerifiedBusiness(ctx, userID); err != nil {
		return nil, err
	}

	if err := m.contracts.CancelContract(ctx, userID, contractID); err != nil {
		return nil, mapContractError(err)
	}

	rec, err := m.contracts.GetContractForBusiness(ctx, userID, contractID)
	if err != nil {
		return nil, err
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.ContractStatusChangeAction),
		EntityType:   "contract",
		EntityID:     &contractID,
		Description:  "business cancelled contract",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	return &businessDto.ContractMutationResponse{Contract: rec.DetailResponse}, nil
}

func (m *businessModule) requireVerifiedBusiness(ctx context.Context, userID string) error {
	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return err
	}

	verified, err := m.business.IsBusinessVerified(ctx, userID)
	if err != nil {
		return err
	}
	if !verified {
		return customErr.ErrAccountNotVerified
	}
	return nil
}

func mapContractError(err error) error {
	if errors.Is(err, contractDto.ErrContractNotFound) {
		return customErr.ErrContractNotFound
	}
	if errors.Is(err, contractDto.ErrInvalidContractStatus) {
		return customErr.ErrInvalidContractStatus
	}
	if errors.Is(err, contractDto.ErrInviteNotAccepted) {
		return customErr.ErrInviteNotAccepted
	}
	if errors.Is(err, contractDto.ErrContractAlreadyExists) {
		return customErr.ErrContractAlreadyExists
	}
	if errors.Is(err, contractDto.ErrPaymentAlreadyExists) {
		return customErr.ErrPaymentAlreadyExists
	}
	if errors.Is(err, contractDto.ErrPaymentRequired) {
		return customErr.ErrPaymentRequired
	}
	return err
}

func validateContractDeliverables(items []contractDto.DeliverableInput) error {
	if len(items) == 0 {
		return customErr.ErrBadRequest
	}
	allowed := map[string]struct{}{
		"tiktok_video": {}, "instagram_reel": {}, "instagram_story": {},
		"youtube_short": {}, "product_review": {}, "post": {}, "other": {},
	}
	for _, item := range items {
		t := strings.ToLower(strings.TrimSpace(item.Type))
		if _, ok := allowed[t]; !ok {
			return customErr.ErrBadRequest
		}
	}
	return nil
}
