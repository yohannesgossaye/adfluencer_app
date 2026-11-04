package business

import (
	"context"

	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	"adinfluence_service_app/internal/module/business/helper"

	"go.uber.org/zap"
)

// SearchCreators is "Find creators". Any active business can browse approved
// creators; inviting still needs a verified business and a live campaign.
func (m *businessModule) SearchCreators(ctx context.Context, userID string, query businessDto.SearchCreatorsQuery) (*businessDto.SearchCreatorsResponse, error) {
	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}
	query = query.Normalize()

	var campaign *businessDto.CampaignRecord
	if query.CampaignID != "" {
		rec, err := m.business.GetCampaignByID(ctx, userID, query.CampaignID)
		if err != nil {
			return nil, mapCampaignError(err)
		}
		campaign = rec
	}

	rows, err := m.business.SearchCreators(ctx, businessDto.CreatorFilter{
		Q: query.Q, Niche: query.Niche, Platform: query.Platform, MaxPrice: query.MaxPrice, City: query.City,
	})
	if err != nil {
		m.logger.Error(ctx, "failed to search creators", zap.String("reason", err.Error()))
		return nil, err
	}
	ranked := helper.RankCreators(rows, query.Q, query.Size, query.Sort, campaign)

	// Invite status per creator for the chosen campaign.
	invited := map[string]string{}
	if campaign != nil {
		if invites, err := m.business.ListInvites(ctx, userID, campaign.ID); err == nil {
			for _, inv := range invites {
				invited[inv.InfluencerUserID] = inv.Status
			}
		}
	}

	start := (query.Page - 1) * query.Limit
	end := start + query.Limit
	if start > len(ranked) {
		start = len(ranked)
	}
	if end > len(ranked) {
		end = len(ranked)
	}
	items := make([]businessDto.CreatorSearchItem, 0, end-start)
	for _, rc := range ranked[start:end] {
		r := rc.Row
		item := businessDto.CreatorSearchItem{
			InfluencerUserID:   r.UserID,
			CreatorName:        r.CreatorName,
			ProfileImageURL:    r.ProfileImageURL,
			City:               r.City,
			Niches:             r.CategoryNames,
			Platforms:          r.Platforms,
			TotalFollowers:     r.TotalFollowers,
			TikTokFollowers:    r.TikTok,
			InstagramFollowers: r.Instagram,
			YouTubeSubs:        r.YouTube,
			Size:               helper.ClassifyCreatorSize(r.TotalFollowers),
			EngagementRate:     r.EngagementRate,
			MinPriceETB:        r.MinPriceETB,
			MaxPriceETB:        r.MaxPriceETB,
			IsVerified:         r.IsVerified,
			MatchScore:         rc.Score,
			ScoreBreakdown:     rc.Breakdown,
		}
		if r.Bio != nil {
			item.Bio = *r.Bio
		}
		if st, ok := invited[r.UserID]; ok {
			item.InviteStatus = &st
		}
		items = append(items, item)
	}

	res := &businessDto.SearchCreatorsResponse{
		Items:      items,
		Pagination: commonDto.NewListMeta(query.Page, query.Limit, len(ranked)),
	}
	if campaign != nil {
		res.RankedFor = &businessDto.RankedForCampaign{ID: campaign.ID, Title: campaign.Title, Status: campaign.Status}
	}
	return res, nil
}
