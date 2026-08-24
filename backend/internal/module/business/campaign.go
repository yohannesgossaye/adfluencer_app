package business

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	"adinfluence_service_app/internal/module/business/helper"
	"context"
	"errors"
	"time"

	"adinfluence_service_app/platform/utils"

	"go.uber.org/zap"
)

func (m *businessModule) CreateCampaign(ctx context.Context, userID string, req businessDto.CreateCampaignRequest, meta authDto.SessionMeta) (*commonDto.CreatedIDResponse, error) {
	m.logger.Info(ctx, "creating new campaign", zap.String("user_id", userID), zap.String("title", req.Title), zap.String("niche", req.Niche))

	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		m.logger.Warn(ctx, "campaign creation rejected: business not active", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	if req.NumberOfCreatorsNeeded <= 0 {
		req.NumberOfCreatorsNeeded = 1
	}
	if err := req.Validate(); err != nil {
		m.logger.Warn(ctx, "campaign validation failed", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	category, err := m.business.ResolveCategoryID(ctx, req.Niche)
	if err != nil {
		if errors.Is(err, businessDto.ErrUnknownCategory) {
			m.logger.Warn(ctx, "unknown niche category for campaign", zap.String("niche", req.Niche))
			return nil, customErr.ErrUnknownCategory
		}
		m.logger.Error(ctx, "failed to resolve category ID", zap.String("niche", req.Niche), zap.String("reason", err.Error()))
		return nil, err
	}

	campaignID := utils.GenerateUUID()
	input, err := helper.BuildSaveCampaignInput(userID, campaignID, req, &category.ID)
	if err != nil {
		m.logger.Error(ctx, "failed to build campaign input", zap.String("campaign_id", campaignID), zap.String("reason", err.Error()))
		return nil, err
	}

	if err := m.business.CreateCampaign(ctx, input); err != nil {
		m.logger.Error(ctx, "failed to persist campaign", zap.String("campaign_id", campaignID), zap.String("reason", err.Error()))
		return nil, err
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.ProfileUpdateAction),
		EntityType:   "campaign",
		EntityID:     &campaignID,
		Description:  "business campaign created",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	m.logger.Info(ctx, "campaign created successfully", zap.String("user_id", userID), zap.String("campaign_id", campaignID))

	return &commonDto.CreatedIDResponse{ID: campaignID}, nil
}

func (m *businessModule) ListCampaigns(ctx context.Context, userID string, query businessDto.ListCampaignsQuery) (*businessDto.ListCampaignsResponse, error) {
	m.logger.Info(ctx, "listing campaigns", zap.String("user_id", userID), zap.String("status", query.Status))

	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}

	query = query.Normalize()
	items, total, err := m.business.ListCampaigns(ctx, userID, query)
	if err != nil {
		m.logger.Error(ctx, "failed to list campaigns from database", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	res := helper.ListCampaignsResponse(items, query, total)
	return &res, nil
}

func (m *businessModule) GetCampaign(ctx context.Context, userID, campaignID string) (*businessDto.CampaignDetailResponse, error) {
	m.logger.Info(ctx, "fetching campaign details", zap.String("user_id", userID), zap.String("campaign_id", campaignID))

	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}

	rec, err := m.business.GetCampaignByID(ctx, userID, campaignID)
	if err != nil {
		m.logger.Warn(ctx, "failed to get campaign by ID", zap.String("campaign_id", campaignID), zap.String("reason", err.Error()))
		return nil, mapCampaignError(err)
	}

	detail := rec.CampaignDetailResponse
	return &detail, nil
}

func (m *businessModule) UpdateCampaign(ctx context.Context, userID, campaignID string, req businessDto.UpdateCampaignRequest, meta authDto.SessionMeta) (*businessDto.CampaignMutationResponse, error) {
	m.logger.Info(ctx, "updating campaign", zap.String("user_id", userID), zap.String("campaign_id", campaignID))

	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}

	if req.NumberOfCreatorsNeeded <= 0 {
		req.NumberOfCreatorsNeeded = 1
	}
	if err := req.Validate(); err != nil {
		m.logger.Warn(ctx, "update campaign validation failed", zap.String("campaign_id", campaignID), zap.String("reason", err.Error()))
		return nil, err
	}

	category, err := m.business.ResolveCategoryID(ctx, req.Niche)
	if err != nil {
		if errors.Is(err, businessDto.ErrUnknownCategory) {
			return nil, customErr.ErrUnknownCategory
		}
		m.logger.Error(ctx, "failed to resolve category ID for campaign update", zap.String("niche", req.Niche), zap.String("reason", err.Error()))
		return nil, err
	}

	input, err := helper.BuildSaveCampaignInput(userID, campaignID, req.CreateCampaignRequest, &category.ID)
	if err != nil {
		return nil, err
	}

	if err := m.business.UpdateCampaign(ctx, input); err != nil {
		m.logger.Error(ctx, "failed to update campaign in persistence", zap.String("campaign_id", campaignID), zap.String("reason", err.Error()))
		return nil, mapCampaignError(err)
	}

	rec, err := m.business.GetCampaignByID(ctx, userID, campaignID)
	if err != nil {
		return nil, err
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.ProfileUpdateAction),
		EntityType:   "campaign",
		EntityID:     &campaignID,
		Description:  "business campaign updated",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	m.logger.Info(ctx, "campaign updated successfully", zap.String("user_id", userID), zap.String("campaign_id", campaignID))

	res := helper.CampaignMutation(*rec)
	return &res, nil
}

func (m *businessModule) PublishCampaign(ctx context.Context, userID, campaignID string, meta authDto.SessionMeta) (*businessDto.CampaignMutationResponse, error) {
	m.logger.Info(ctx, "publishing campaign", zap.String("user_id", userID), zap.String("campaign_id", campaignID))

	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}

	rec, err := m.business.GetCampaignByID(ctx, userID, campaignID)
	if err != nil {
		return nil, mapCampaignError(err)
	}

	if rec.Status != "draft" {
		m.logger.Warn(ctx, "cannot publish non-draft campaign", zap.String("campaign_id", campaignID), zap.String("status", rec.Status))
		return nil, customErr.ErrInvalidCampaignStatus
	}

	if err := helper.ValidateCampaignRecordForPublish(*rec); err != nil {
		m.logger.Warn(ctx, "campaign record incomplete for publish", zap.String("campaign_id", campaignID), zap.String("reason", err.Error()))
		return nil, customErr.ErrProfileIncomplete
	}

	if err := m.business.PublishCampaign(ctx, userID, campaignID); err != nil {
		m.logger.Error(ctx, "failed to publish campaign in persistence", zap.String("campaign_id", campaignID), zap.String("reason", err.Error()))
		return nil, mapCampaignError(err)
	}

	rec.Status = "active"
	rec.UpdatedAt = time.Now().UTC()

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.StatusChangeAction),
		EntityType:   "campaign",
		EntityID:     &campaignID,
		Description:  "business campaign published",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	m.logger.Info(ctx, "campaign published successfully", zap.String("user_id", userID), zap.String("campaign_id", campaignID))

	res := helper.CampaignMutation(*rec)
	return &res, nil
}

func (m *businessModule) CancelCampaign(ctx context.Context, userID, campaignID string, meta authDto.SessionMeta) (*businessDto.CampaignMutationResponse, error) {
	m.logger.Info(ctx, "cancelling campaign", zap.String("user_id", userID), zap.String("campaign_id", campaignID))

	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}

	if err := m.business.CancelCampaign(ctx, userID, campaignID); err != nil {
		m.logger.Error(ctx, "failed to cancel campaign", zap.String("campaign_id", campaignID), zap.String("reason", err.Error()))
		return nil, mapCampaignError(err)
	}

	rec, err := m.business.GetCampaignByID(ctx, userID, campaignID)
	if err != nil {
		return nil, err
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.StatusChangeAction),
		EntityType:   "campaign",
		EntityID:     &campaignID,
		Description:  "business campaign cancelled",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	m.logger.Info(ctx, "campaign cancelled successfully", zap.String("user_id", userID), zap.String("campaign_id", campaignID))

	res := helper.CampaignMutation(*rec)
	return &res, nil
}

func (m *businessModule) requireActiveBusiness(ctx context.Context, userID string) error {
	user, err := m.auth.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, customErr.ErrUserNotFound) {
			return customErr.ErrUnauthorized
		}
		return err
	}
	if err := ensureEditableStatus(user.Status); err != nil {
		return err
	}
	switch model.UserStatus(user.Status) {
	case model.ActiveStatus, model.PendingProfileStatus, model.PendingApprovalStatus:
		return nil
	default:
		return customErr.ErrAccountNotActive
	}
	return nil
}

func mapCampaignError(err error) error {
	if errors.Is(err, businessDto.ErrCampaignNotFound) {
		return customErr.ErrCampaignNotFound
	}
	if errors.Is(err, businessDto.ErrInvalidCampaignStatus) {
		return customErr.ErrInvalidCampaignStatus
	}
	return err
}
