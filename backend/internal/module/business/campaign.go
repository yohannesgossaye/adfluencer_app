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

	"adinfluence_service_app/platform/utils"
)

func (m *businessModule) CreateCampaign(ctx context.Context, userID string, req businessDto.CreateCampaignRequest, meta authDto.SessionMeta) (*commonDto.CreatedIDResponse, error) {
	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}

	if req.NumberOfCreatorsNeeded <= 0 {
		req.NumberOfCreatorsNeeded = 1
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}

	category, err := m.business.ResolveCategoryID(ctx, req.Niche)
	if err != nil {
		if errors.Is(err, businessDto.ErrUnknownCategory) {
			return nil, customErr.ErrUnknownCategory
		}
		return nil, err
	}

	campaignID := utils.GenerateUUID()
	input, err := helper.BuildSaveCampaignInput(userID, campaignID, req, &category.ID)
	if err != nil {
		return nil, err
	}

	if err := m.business.CreateCampaign(ctx, input); err != nil {
		return nil, err
	}

	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.ProfileUpdateAction),
		EntityType:   "campaign",
		EntityID:     &campaignID,
		Description:  "business campaign created",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	return &commonDto.CreatedIDResponse{ID: campaignID}, nil
}

func (m *businessModule) ListCampaigns(ctx context.Context, userID string, query businessDto.ListCampaignsQuery) (*businessDto.ListCampaignsResponse, error) {
	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}

	query = query.Normalize()
	items, total, err := m.business.ListCampaigns(ctx, userID, query)
	if err != nil {
		return nil, err
	}

	res := helper.ListCampaignsResponse(items, query, total)
	return &res, nil
}

func (m *businessModule) GetCampaign(ctx context.Context, userID, campaignID string) (*businessDto.CampaignDetailResponse, error) {
	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}

	rec, err := m.business.GetCampaignByID(ctx, userID, campaignID)
	if err != nil {
		return nil, mapCampaignError(err)
	}

	detail := rec.CampaignDetailResponse
	return &detail, nil
}

func (m *businessModule) UpdateCampaign(ctx context.Context, userID, campaignID string, req businessDto.UpdateCampaignRequest, meta authDto.SessionMeta) (*businessDto.CampaignMutationResponse, error) {
	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}

	if req.NumberOfCreatorsNeeded <= 0 {
		req.NumberOfCreatorsNeeded = 1
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if _, err := m.business.GetCampaignByID(ctx, userID, campaignID); err != nil {
		return nil, mapCampaignError(err)
	}

	category, err := m.business.ResolveCategoryID(ctx, req.Niche)
	if err != nil {
		if errors.Is(err, businessDto.ErrUnknownCategory) {
			return nil, customErr.ErrUnknownCategory
		}
		return nil, err
	}

	input, err := helper.BuildSaveCampaignInput(userID, campaignID, req.CreateCampaignRequest, &category.ID)
	if err != nil {
		return nil, err
	}

	if err := m.business.UpdateCampaign(ctx, input); err != nil {
		return nil, mapCampaignError(err)
	}

	rec, err := m.business.GetCampaignByID(ctx, userID, campaignID)
	if err != nil {
		return nil, err
	}

	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.ProfileUpdateAction),
		EntityType:   "campaign",
		EntityID:     &campaignID,
		Description:  "business campaign updated",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	res := helper.CampaignMutation(*rec)
	return &res, nil
}

func (m *businessModule) PublishCampaign(ctx context.Context, userID, campaignID string, meta authDto.SessionMeta) (*businessDto.CampaignMutationResponse, error) {
	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}

	rec, err := m.business.GetCampaignByID(ctx, userID, campaignID)
	if err != nil {
		return nil, mapCampaignError(err)
	}

	if rec.Status != "draft" {
		return nil, customErr.ErrInvalidCampaignStatus
	}

	if err := helper.ValidateCampaignRecordForPublish(*rec); err != nil {
		return nil, customErr.ErrProfileIncomplete
	}

	if err := m.business.PublishCampaign(ctx, userID, campaignID); err != nil {
		return nil, mapCampaignError(err)
	}

	rec, err = m.business.GetCampaignByID(ctx, userID, campaignID)
	if err != nil {
		return nil, err
	}

	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.StatusChangeAction),
		EntityType:   "campaign",
		EntityID:     &campaignID,
		Description:  "business campaign published",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	res := helper.CampaignMutation(*rec)
	return &res, nil
}

func (m *businessModule) CancelCampaign(ctx context.Context, userID, campaignID string, meta authDto.SessionMeta) (*businessDto.CampaignMutationResponse, error) {
	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}

	if _, err := m.business.GetCampaignByID(ctx, userID, campaignID); err != nil {
		return nil, mapCampaignError(err)
	}

	if err := m.business.CancelCampaign(ctx, userID, campaignID); err != nil {
		return nil, mapCampaignError(err)
	}

	rec, err := m.business.GetCampaignByID(ctx, userID, campaignID)
	if err != nil {
		return nil, err
	}

	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.StatusChangeAction),
		EntityType:   "campaign",
		EntityID:     &campaignID,
		Description:  "business campaign cancelled",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

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
	if model.UserStatus(user.Status) != model.ActiveStatus {
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
