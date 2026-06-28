package storage

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"context"
)

type AdminPersistence interface {
	ListInfluencerApplications(ctx context.Context, query adminDto.ListApplicationsQuery) ([]adminDto.ApplicationRow, int, error)
	GetInfluencerApplication(ctx context.Context, userID string) (*adminDto.ApplicationRow, error)
	ApproveInfluencerApplication(ctx context.Context, userID, adminUserID string, adminNote *string) error
	RejectInfluencerApplication(ctx context.Context, userID, reason string, adminNote *string) error
	UpdateInfluencerAdminNote(ctx context.Context, userID string, adminNote string) error
	CreateAuditLog(ctx context.Context, input authDto.CreateAuditLogInput) error

	ListBusinessVerifications(ctx context.Context, query adminDto.ListBusinessVerificationsQuery) ([]adminDto.BusinessVerificationListItem, int, error)
	GetBusinessVerification(ctx context.Context, userID string) (*adminDto.BusinessVerificationRow, error)
	VerifyBusiness(ctx context.Context, userID string) error
	RevokeBusinessVerification(ctx context.Context, userID string) error

	ListCategories(ctx context.Context) ([]adminDto.CategoryRow, error)
	GetCategoryByID(ctx context.Context, categoryID string) (*adminDto.CategoryRow, error)
	CreateCategory(ctx context.Context, id, name, slug string, description *string) error
	UpdateCategory(ctx context.Context, categoryID string, name, slug *string, description *string) error
	SetCategoryActive(ctx context.Context, categoryID string, active bool) error
	CategorySlugExists(ctx context.Context, slug, excludeID string) (bool, error)

	GetUserStats(ctx context.Context) (*adminDto.UserStats, error)
	ListUsers(ctx context.Context, query adminDto.ListUsersQuery) ([]adminDto.UserRow, int, error)
	GetUser(ctx context.Context, userID string) (*adminDto.UserRow, error)
	BlockUser(ctx context.Context, userID, reason string) error
	UnblockUser(ctx context.Context, userID string) error
	RemoveUser(ctx context.Context, userID, reason string) error

	GetDashboardKPIs(ctx context.Context) (*adminDto.DashboardKPIRow, error)
	ListDashboardPendingApprovals(ctx context.Context, limit int) ([]adminDto.DashboardPendingApprovalRow, error)
	ListDashboardRecentActivity(ctx context.Context, limit int) ([]adminDto.DashboardActivityRow, error)
	ListDashboardMonthlyRevenue(ctx context.Context, months int) ([]adminDto.DashboardMonthlyRevenueRow, error)
	ListDashboardAlerts(ctx context.Context, limit int) ([]adminDto.DashboardAlertRow, error)

	GetAdminCampaignStats(ctx context.Context) (*adminDto.AdminCampaignStats, error)
	ListAdminCampaigns(ctx context.Context, query adminDto.ListAdminCampaignsQuery) ([]adminDto.AdminCampaignRow, int, error)
	GetAdminCampaign(ctx context.Context, campaignID string) (*adminDto.AdminCampaignRow, error)
	SetAdminCampaignStatus(ctx context.Context, campaignID, status string) error

	ListAdminDisputes(ctx context.Context, query adminDto.ListAdminDisputesQuery) ([]adminDto.DisputeRow, int, error)
	GetAdminDispute(ctx context.Context, disputeID string) (*adminDto.DisputeRow, error)
	UpdateDisputeStatus(ctx context.Context, disputeID, status, adminUserID string) error
	UpdateDisputeResolutionNote(ctx context.Context, disputeID string, note *string) error
	SaveDisputeAdminNote(ctx context.Context, adminUserID, disputeID, note string) error
}
