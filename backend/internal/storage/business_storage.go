package storage

import (
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	"context"
)

type BusinessPersistence interface {
	GetProfileByUserID(ctx context.Context, userID string) (*businessDto.ProfileRecord, error)
	UpdateProfile(ctx context.Context, input businessDto.UpdateProfileInput) error
	UpdateLogoURL(ctx context.Context, userID, logoURL string) error
	CompleteProfile(ctx context.Context, userID string) error
	CreateAuditLog(ctx context.Context, input authDto.CreateAuditLogInput) error

	ResolveCategoryID(ctx context.Context, nameOrSlug string) (*businessDto.CampaignCategoryResponse, error)
	CreateCampaign(ctx context.Context, input businessDto.SaveCampaignInput) (*businessDto.CampaignRecord, error)
	ListCampaigns(ctx context.Context, userID string, query businessDto.ListCampaignsQuery) ([]businessDto.CampaignListItem, int, error)
	GetCampaignByID(ctx context.Context, userID, campaignID string) (*businessDto.CampaignRecord, error)
	UpdateCampaign(ctx context.Context, input businessDto.SaveCampaignInput) error
	PublishCampaign(ctx context.Context, userID, campaignID string) error
	CancelCampaign(ctx context.Context, userID, campaignID string) error

	ListMatchCandidates(ctx context.Context) ([]businessDto.MatchCandidate, error)
	ReplaceRecommendations(ctx context.Context, campaignID string, items []businessDto.SaveRecommendationInput) (int, error)
	ListRecommendations(ctx context.Context, campaignID string, query businessDto.ListRecommendationsQuery) ([]businessDto.RecommendationItem, int, error)
	GetRecommendation(ctx context.Context, campaignID, influencerUserID string) (*businessDto.RecommendationRecord, error)
	IsVerifiedActiveInfluencer(ctx context.Context, userID string) (bool, error)

	CreateInvite(ctx context.Context, input businessDto.CreateInviteInput) (*businessDto.InviteResponse, error)
	ListInvites(ctx context.Context, businessUserID, campaignID string) ([]businessDto.InviteResponse, error)
	GetInviteByID(ctx context.Context, businessUserID, campaignID, inviteID string) (*businessDto.InviteResponse, error)
	CancelInvite(ctx context.Context, businessUserID, campaignID, inviteID string) error

	GetDashboardKPIs(ctx context.Context, businessUserID string) (*businessDto.DashboardKPIRow, error)
	ListDashboardRecentCampaigns(ctx context.Context, businessUserID string) ([]businessDto.DashboardCampaignRow, error)
	ListDashboardTopRecommendations(ctx context.Context, businessUserID string) ([]businessDto.DashboardRecommendationRow, error)
}
