package business

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	"adinfluence_service_app/internal/module/business/helper"
	"context"

	"adinfluence_service_app/platform/utils"
)

func (m *businessModule) GenerateRecommendations(ctx context.Context, userID, campaignID string, meta authDto.SessionMeta) (*businessDto.GenerateRecommendationsResponse, error) {
	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}

	campaign, err := m.requireActiveCampaign(ctx, userID, campaignID)
	if err != nil {
		return nil, err
	}

	candidates, err := m.business.ListMatchCandidates(ctx)
	if err != nil {
		return nil, err
	}

	ranked := helper.RankCandidates(*campaign, candidates)
	for i := range ranked {
		ranked[i].ID = utils.GenerateUUID()
		ranked[i].CampaignID = campaignID
	}

	generated, err := m.business.ReplaceRecommendations(ctx, campaignID, ranked)
	if err != nil {
		return nil, err
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.ProfileUpdateAction),
		EntityType:   "campaign_recommendations",
		EntityID:     &campaignID,
		Description:  "campaign recommendations generated",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	return &businessDto.GenerateRecommendationsResponse{
		Generated: generated,
		Total:     generated,
	}, nil
}

func (m *businessModule) ListRecommendations(ctx context.Context, userID, campaignID string, query businessDto.ListRecommendationsQuery) (*businessDto.ListRecommendationsResponse, error) {
	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}

	if _, err := m.business.GetCampaignByID(ctx, userID, campaignID); err != nil {
		return nil, mapCampaignError(err)
	}

	query = query.Normalize()
	items, total, err := m.business.ListRecommendations(ctx, campaignID, query)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []businessDto.RecommendationItem{}
	}

	return &businessDto.ListRecommendationsResponse{
		Items: items,
		Pagination: commonDto.NewListMeta(query.Page, query.Limit, total),
	}, nil
}

func (m *businessModule) requireActiveCampaign(ctx context.Context, userID, campaignID string) (*businessDto.CampaignRecord, error) {
	rec, err := m.business.GetCampaignByID(ctx, userID, campaignID)
	if err != nil {
		return nil, mapCampaignError(err)
	}
	if rec.Status != "active" {
		return nil, customErr.ErrInvalidCampaignStatus
	}
	return rec, nil
}
