package influencer

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
	"context"
	"errors"
)

func (m *influencerModule) ListContracts(ctx context.Context, userID string, query contractDto.ListQuery) (*influencerDto.ListContractsResponse, error) {
	if err := m.requireActiveInfluencer(ctx, userID); err != nil {
		return nil, err
	}

	query = query.Normalize()
	items, total, err := m.contracts.ListContractsForInfluencer(ctx, userID, query)
	if err != nil {
		return nil, err
	}

	return &influencerDto.ListContractsResponse{
		Items: items,
		Pagination: contractDto.PaginationResponse{
			Page:  query.Page,
			Limit: query.Limit,
			Total: total,
		},
	}, nil
}

func (m *influencerModule) GetContract(ctx context.Context, userID, contractID string) (*contractDto.DetailResponse, error) {
	if err := m.requireActiveInfluencer(ctx, userID); err != nil {
		return nil, err
	}

	rec, err := m.contracts.GetContractForInfluencer(ctx, userID, contractID)
	if err != nil {
		return nil, mapContractError(err)
	}
	detail := rec.DetailResponse
	return &detail, nil
}

func (m *influencerModule) AcceptContract(ctx context.Context, userID, contractID string, meta authDto.SessionMeta) (*influencerDto.RespondContractResponse, error) {
	if err := m.requireVerifiedInfluencer(ctx, userID); err != nil {
		return nil, err
	}

	rec, err := m.contracts.GetContractForInfluencer(ctx, userID, contractID)
	if err != nil {
		return nil, mapContractError(err)
	}

	if err := m.contracts.AcceptContract(ctx, userID, contractID); err != nil {
		return nil, mapContractError(err)
	}

	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &rec.BusinessUserID,
		Action:       string(model.ContractStatusChangeAction),
		EntityType:   "contract",
		EntityID:     &contractID,
		Description:  "influencer accepted contract",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	rec, err = m.contracts.GetContractForInfluencer(ctx, userID, contractID)
	if err != nil {
		return nil, err
	}

	return &influencerDto.RespondContractResponse{Contract: rec.DetailResponse}, nil
}

func (m *influencerModule) RejectContract(ctx context.Context, userID, contractID string, meta authDto.SessionMeta) (*influencerDto.RespondContractResponse, error) {
	if err := m.requireActiveInfluencer(ctx, userID); err != nil {
		return nil, err
	}

	rec, err := m.contracts.GetContractForInfluencer(ctx, userID, contractID)
	if err != nil {
		return nil, mapContractError(err)
	}

	if err := m.contracts.RejectContract(ctx, userID, contractID); err != nil {
		return nil, mapContractError(err)
	}

	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &rec.BusinessUserID,
		Action:       string(model.ContractStatusChangeAction),
		EntityType:   "contract",
		EntityID:     &contractID,
		Description:  "influencer rejected contract",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	rec, err = m.contracts.GetContractForInfluencer(ctx, userID, contractID)
	if err != nil {
		return nil, err
	}

	return &influencerDto.RespondContractResponse{Contract: rec.DetailResponse}, nil
}

func mapContractError(err error) error {
	if errors.Is(err, contractDto.ErrContractNotFound) {
		return customErr.ErrContractNotFound
	}
	if errors.Is(err, contractDto.ErrInvalidContractStatus) {
		return customErr.ErrInvalidContractStatus
	}
	return err
}
