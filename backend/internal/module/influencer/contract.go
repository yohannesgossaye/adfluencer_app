package influencer

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
	"context"
	"errors"

	"go.uber.org/zap"
)

func (m *influencerModule) ListContracts(ctx context.Context, userID string, query contractDto.ListQuery) (*influencerDto.ListContractsResponse, error) {
	m.logger.Info(ctx, "listing influencer contracts", zap.String("user_id", userID), zap.String("status", query.Status))

	if err := m.requireActiveInfluencer(ctx, userID); err != nil {
		return nil, err
	}

	query = query.Normalize()
	items, total, err := m.contracts.ListContractsForInfluencer(ctx, userID, query)
	if err != nil {
		m.logger.Error(ctx, "failed to query influencer contracts", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	return &influencerDto.ListContractsResponse{
		Items:      items,
		Pagination: commonDto.NewListMeta(query.Page, query.Limit, total),
	}, nil
}

func (m *influencerModule) GetContract(ctx context.Context, userID, contractID string) (*contractDto.DetailResponse, error) {
	m.logger.Info(ctx, "fetching influencer contract details", zap.String("user_id", userID), zap.String("contract_id", contractID))

	if err := m.requireActiveInfluencer(ctx, userID); err != nil {
		return nil, err
	}

	rec, err := m.contracts.GetContractForInfluencer(ctx, userID, contractID)
	if err != nil {
		m.logger.Warn(ctx, "contract lookup failed", zap.String("contract_id", contractID), zap.String("reason", err.Error()))
		return nil, mapContractError(err)
	}
	detail := rec.DetailResponse
	return &detail, nil
}

func (m *influencerModule) AcceptContract(ctx context.Context, userID, contractID string, meta authDto.SessionMeta) (*influencerDto.RespondContractResponse, error) {
	m.logger.Info(ctx, "influencer accepting contract", zap.String("user_id", userID), zap.String("contract_id", contractID))

	if err := m.requireVerifiedInfluencer(ctx, userID); err != nil {
		m.logger.Warn(ctx, "cannot accept contract: influencer not verified", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	rec, err := m.contracts.GetContractForInfluencer(ctx, userID, contractID)
	if err != nil {
		return nil, mapContractError(err)
	}

	if err := m.contracts.AcceptContract(ctx, userID, contractID); err != nil {
		m.logger.Error(ctx, "failed to accept contract in persistence", zap.String("contract_id", contractID), zap.String("reason", err.Error()))
		return nil, mapContractError(err)
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
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

	m.logger.Info(ctx, "influencer accepted contract successfully", zap.String("contract_id", contractID))

	return &influencerDto.RespondContractResponse{Contract: rec.DetailResponse}, nil
}

func (m *influencerModule) RejectContract(ctx context.Context, userID, contractID string, meta authDto.SessionMeta) (*influencerDto.RespondContractResponse, error) {
	m.logger.Info(ctx, "influencer rejecting contract", zap.String("user_id", userID), zap.String("contract_id", contractID))

	if err := m.requireActiveInfluencer(ctx, userID); err != nil {
		return nil, err
	}

	rec, err := m.contracts.GetContractForInfluencer(ctx, userID, contractID)
	if err != nil {
		return nil, mapContractError(err)
	}

	if err := m.contracts.RejectContract(ctx, userID, contractID); err != nil {
		m.logger.Error(ctx, "failed to reject contract in persistence", zap.String("contract_id", contractID), zap.String("reason", err.Error()))
		return nil, mapContractError(err)
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
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

	m.logger.Info(ctx, "influencer rejected contract successfully", zap.String("contract_id", contractID))

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
