package admin

import "time"

type DashboardSummaryResponse struct {
	KPIs              DashboardKPIs                  `json:"kpis"`
	PendingApprovals  []DashboardPendingApprovalRow  `json:"pending_approvals"`
	RecentActivity    []DashboardActivityRow         `json:"recent_activity"`
	RevenueSummary    DashboardRevenueSummary        `json:"revenue_summary"`
	Alerts            []DashboardAlertRow            `json:"alerts"`
}

type DashboardKPIs struct {
	TotalInfluencers           int     `json:"total_influencers"`
	TotalBusinesses            int     `json:"total_businesses"`
	PendingInfluencerApprovals int     `json:"pending_influencer_approvals"`
	ActiveCampaigns            int     `json:"active_campaigns"`
	TotalCampaigns             int     `json:"total_campaigns"`
	TotalUsers                 int     `json:"total_users"`
	OpenDisputes               int     `json:"open_disputes"`
	OpenSupportTickets         int     `json:"open_support_tickets"`
	BlockedUsers               int     `json:"blocked_users"`
	CompletedContracts         int     `json:"completed_contracts"`
	MonthlyPlatformRevenueETB  float64 `json:"monthly_platform_revenue_etb"`
}

type DashboardRevenueSummary struct {
	TotalCampaignValueETB        float64                    `json:"total_campaign_value_etb"`
	TotalPlatformFeesETB         float64                    `json:"total_platform_fees_etb"`
	PlatformFeePercentBusiness   float64                    `json:"platform_fee_percent_business"`
	PlatformFeePercentInfluencer float64                    `json:"platform_fee_percent_influencer"`
	MonthlyTrend                 []DashboardMonthlyRevenueRow `json:"monthly_trend"`
}

type DashboardMonthlyRevenueRow struct {
	Month              time.Time `json:"month"`
	PlatformRevenueETB float64   `json:"platform_revenue_etb"`
	CompletedContracts int       `json:"completed_contracts"`
}

type DashboardPendingApprovalRow struct {
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Niche       string    `json:"niche,omitempty"`
	Platform    string    `json:"platform,omitempty"`
	Followers   int       `json:"followers"`
	SubmittedAt time.Time `json:"submitted_at"`
}

type DashboardActivityRow struct {
	ID        string    `json:"id"`
	Action    string    `json:"action"`
	Label     string    `json:"label"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

type DashboardAlertRow struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Detail     string `json:"detail"`
	Severity   string `json:"severity"`
	EntityType string `json:"entity_type,omitempty"`
}

type DashboardKPIRow struct {
	TotalInfluencers           int
	TotalBusinesses            int
	PendingInfluencerApprovals int
	ActiveCampaigns            int
	TotalCampaigns             int
	TotalUsers                 int
	OpenDisputes               int
	OpenSupportTickets         int
	BlockedUsers               int
	CompletedContracts         int
	MonthlyPlatformRevenueETB  float64
	TotalCampaignValueETB      float64
	TotalPlatformFeesETB       float64
}
