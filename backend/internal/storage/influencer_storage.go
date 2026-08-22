package storage

import (
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
	"context"
)

type InfluencerPersistence interface {
	GetProfileByUserID(ctx context.Context, userID string) (*influencerDto.ProfileRecord, error)
	UpdateProfile(ctx context.Context, input influencerDto.UpdateProfileInput) error
	UpdateProfileImageURL(ctx context.Context, userID string, imageURL *string) error
	CountPortfolioItems(ctx context.Context, userID string) (int, error)
	ListPortfolioItems(ctx context.Context, userID string) ([]influencerDto.PortfolioItem, error)
	GetPortfolioItem(ctx context.Context, userID, itemID string) (*influencerDto.PortfolioItemRecord, error)
	CreatePortfolioItem(ctx context.Context, input influencerDto.CreatePortfolioItemInput) (*influencerDto.PortfolioItem, error)
	DeletePortfolioItem(ctx context.Context, userID, itemID string) error
	ResolveCategoryIDs(ctx context.Context, namesOrSlugs []string) ([]influencerDto.CategoryResponse, error)
	ListActiveCategories(ctx context.Context) ([]influencerDto.CategoryResponse, error)
	GetCategoriesForUser(ctx context.Context, userID string) ([]influencerDto.CategoryResponse, error)
	GetPlatformsForUser(ctx context.Context, userID string) ([]string, error)
	IsProfileVerified(ctx context.Context, userID string) (bool, error)
	SubmitProfile(ctx context.Context, userID string) error
	CreateAuditLog(ctx context.Context, input authDto.CreateAuditLogInput) error

	ListInvites(ctx context.Context, influencerUserID string, query influencerDto.ListInvitesQuery) ([]influencerDto.InviteListItem, error)
	GetInviteByID(ctx context.Context, influencerUserID, inviteID string) (*influencerDto.InviteRecord, error)
	AcceptInvite(ctx context.Context, influencerUserID, inviteID string) (*influencerDto.RespondInviteResponse, error)
	RejectInvite(ctx context.Context, influencerUserID, inviteID string) (*influencerDto.RespondInviteResponse, error)

	GetDashboardKPIs(ctx context.Context, influencerUserID string) (*influencerDto.DashboardKPIRow, error)
	ListDashboardPendingInvites(ctx context.Context, influencerUserID string) ([]influencerDto.DashboardInviteRow, error)
	ListDashboardActiveContracts(ctx context.Context, influencerUserID string) ([]influencerDto.DashboardContractRow, error)
}
