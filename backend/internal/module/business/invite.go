package business

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
	"context"
	"errors"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"adinfluence_service_app/platform/utils"
)

func (m *businessModule) GetInfluencerForBusiness(ctx context.Context, userID, influencerUserID, campaignID string) (*businessDto.BusinessInfluencerViewResponse, error) {
	m.logger.Info(ctx, "business viewing influencer profile", zap.String("business_id", userID), zap.String("influencer_id", influencerUserID))

	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}

	var (
		verified bool
		rec      *influencerDto.ProfileRecord
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		verified, err = m.business.IsVerifiedActiveInfluencer(ctx, influencerUserID)
		return err
	})

	g.Go(func() error {
		var err error
		rec, err = m.influencer.GetProfileByUserID(ctx, influencerUserID)
		if err != nil {
			if errors.Is(err, influencerDto.ErrProfileNotFound) {
				return customErr.ErrCreatorNotFound
			}
		}
		return err
	})

	if err := g.Wait(); err != nil {
		m.logger.Error(ctx, "failed to query influencer profile details", zap.String("influencer_id", influencerUserID), zap.String("reason", err.Error()))
		return nil, err
	}
	if !verified {
		m.logger.Warn(ctx, "influencer is not verified or active", zap.String("influencer_id", influencerUserID))
		return nil, customErr.ErrCreatorNotFound
	}

	res := buildInfluencerView(rec)
	if strings.TrimSpace(campaignID) != "" {
		var match *businessDto.RecommendationRecord

		g, ctx := errgroup.WithContext(ctx)

		g.Go(func() error {
			_, err := m.business.GetCampaignByID(ctx, userID, campaignID)
			return mapCampaignError(err)
		})

		g.Go(func() error {
			var err error
			match, err = m.business.GetRecommendation(ctx, campaignID, influencerUserID)
			return err
		})

		if err := g.Wait(); err != nil {
			return nil, err
		}

		if match != nil {
			score := match.MatchScore
			breakdown := match.ScoreBreakdown
			res.MatchScore = &score
			res.ScoreBreakdown = &breakdown
		}
	}

	return res, nil
}

func buildInfluencerView(rec *influencerDto.ProfileRecord) *businessDto.BusinessInfluencerViewResponse {
	detail := rec.ProfileDetailResponse
	if detail.Languages == nil {
		detail.Languages = []string{}
	}
	if detail.Niches == nil {
		detail.Niches = []influencerDto.CategoryResponse{}
	}
	if detail.Platforms == nil {
		detail.Platforms = []string{}
	}
	if detail.PortfolioItems == nil {
		detail.PortfolioItems = []influencerDto.PortfolioItem{}
	}

	return &businessDto.BusinessInfluencerViewResponse{
		InfluencerUserID: rec.UserID,
		CreatorName:      detail.CreatorName,
		FullName:         detail.FullName,
		Bio:              detail.Bio,
		ProfileImageURL:  detail.ProfileImageURL,
		City:             detail.City,
		Region:           detail.Region,
		Languages:        detail.Languages,
		Niches:           detail.Niches,
		Platforms:        detail.Platforms,
		SocialLinks:      detail.SocialLinks,
		TotalFollowers:   detail.TotalFollowers,
		EngagementRate:   detail.EngagementRate,
		MinPriceETB:      detail.MinPriceETB,
		MaxPriceETB:      detail.MaxPriceETB,
		IsVerified:       detail.IsVerified,
		PortfolioItems:   detail.PortfolioItems,
	}
}

func (m *businessModule) CreateInvite(ctx context.Context, userID, campaignID string, req businessDto.CreateInviteRequest, meta authDto.SessionMeta) (*commonDto.CreatedIDResponse, error) {
	m.logger.Info(ctx, "creating campaign invite", zap.String("business_id", userID), zap.String("campaign_id", campaignID), zap.String("influencer_id", req.InfluencerUserID))

	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}
	if err := req.Validate(); err != nil {
		m.logger.Warn(ctx, "create invite validation failed", zap.String("reason", err.Error()))
		return nil, err
	}

	if _, err := m.requireActiveCampaign(ctx, userID, campaignID); err != nil {
		return nil, err
	}

	ok, err := m.business.IsVerifiedActiveInfluencer(ctx, req.InfluencerUserID)
	if err != nil {
		return nil, err
	}
	if !ok {
		m.logger.Warn(ctx, "cannot invite non-verified creator", zap.String("influencer_id", req.InfluencerUserID))
		return nil, customErr.ErrCreatorNotFound
	}

	inviteID := utils.GenerateUUID()
	if err := m.business.CreateInvite(ctx, businessDto.CreateInviteInput{
		ID:               inviteID,
		CampaignID:       campaignID,
		BusinessUserID:   userID,
		InfluencerUserID: req.InfluencerUserID,
		Message:          businessDto.OptionalInviteMessage(req.Message),
	}); err != nil {
		m.logger.Error(ctx, "failed to persist campaign invite", zap.String("invite_id", inviteID), zap.String("reason", err.Error()))
		return nil, mapInviteError(err)
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &req.InfluencerUserID,
		Action:       string(model.StatusChangeAction),
		EntityType:   "campaign_invite",
		EntityID:     &inviteID,
		Description:  "business sent campaign invite",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	m.logger.Info(ctx, "campaign invite created successfully", zap.String("invite_id", inviteID))

	return &commonDto.CreatedIDResponse{ID: inviteID}, nil
}

func (m *businessModule) ListInvites(ctx context.Context, userID, campaignID string) (*businessDto.ListInvitesResponse, error) {
	m.logger.Info(ctx, "listing campaign invites", zap.String("business_id", userID), zap.String("campaign_id", campaignID))

	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}

	items, err := m.business.ListInvites(ctx, userID, campaignID)
	if err != nil {
		m.logger.Error(ctx, "failed to query campaign invites", zap.String("campaign_id", campaignID), zap.String("reason", err.Error()))
		return nil, mapCampaignError(err)
	}
	if items == nil {
		items = []businessDto.InviteResponse{}
	}

	return &businessDto.ListInvitesResponse{Items: items}, nil
}

func (m *businessModule) CancelInvite(ctx context.Context, userID, campaignID, inviteID string, meta authDto.SessionMeta) error {
	m.logger.Info(ctx, "cancelling campaign invite", zap.String("business_id", userID), zap.String("campaign_id", campaignID), zap.String("invite_id", inviteID))

	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return err
	}

	if _, err := m.business.GetCampaignByID(ctx, userID, campaignID); err != nil {
		return mapCampaignError(err)
	}

	if err := m.business.CancelInvite(ctx, userID, campaignID, inviteID); err != nil {
		m.logger.Error(ctx, "failed to cancel campaign invite", zap.String("invite_id", inviteID), zap.String("reason", err.Error()))
		return mapInviteError(err)
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.StatusChangeAction),
		EntityType:   "campaign_invite",
		EntityID:     &inviteID,
		Description:  "business cancelled campaign invite",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	m.logger.Info(ctx, "campaign invite cancelled successfully", zap.String("invite_id", inviteID))

	return nil
}

func mapInviteError(err error) error {
	if errors.Is(err, businessDto.ErrInviteNotFound) {
		return customErr.ErrInviteNotFound
	}
	if errors.Is(err, businessDto.ErrInviteAlreadyExists) {
		return customErr.ErrInviteAlreadyExists
	}
	if errors.Is(err, businessDto.ErrInvalidInviteStatus) {
		return customErr.ErrInvalidInviteStatus
	}
	return err
}
