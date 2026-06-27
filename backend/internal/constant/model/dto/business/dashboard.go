package business

import "time"

type DashboardSummaryResponse struct {
	Greeting           DashboardGreeting            `json:"greeting"`
	MarketplaceAccess  MarketplaceAccessResponse    `json:"marketplace_access"`
	KPIs               DashboardKPIs                `json:"kpis"`
	ContractsSummary   DashboardContractsSummary    `json:"contracts_summary"`
	RecentCampaigns    []DashboardCampaignRow       `json:"recent_campaigns"`
	TopRecommendations []DashboardRecommendationRow `json:"top_recommendations"`
}

type DashboardGreeting struct {
	BusinessName  string `json:"business_name"`
	OwnerFullName string `json:"owner_full_name"`
}

type DashboardKPIs struct {
	ActiveCampaigns              int `json:"active_campaigns"`
	TotalCampaigns               int `json:"total_campaigns"`
	RecommendedCreatorsTotal     int `json:"recommended_creators_total"`
	RecommendedCreatorsHighMatch int `json:"recommended_creators_high_match"`
	PendingContracts             int `json:"pending_contracts"`
	ContractsAwaitingInfluencer  int `json:"contracts_awaiting_influencer"`
	ContractsAwaitingPayment     int `json:"contracts_awaiting_payment"`
	PendingInvites               int `json:"pending_invites"`
}

type DashboardContractsSummary struct {
	Pending   int `json:"pending"`
	Active    int `json:"active"`
	Completed int `json:"completed"`
}

type DashboardCampaignRow struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	Niche            string    `json:"niche,omitempty"`
	BudgetSummary    string    `json:"budget_summary,omitempty"`
	Status           string    `json:"status"`
	ProgressPercent  *int      `json:"progress_percent,omitempty"`
	AcceptedCreators int       `json:"accepted_creators"`
	CreatorsNeeded   int       `json:"creators_needed"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type DashboardRecommendationRow struct {
	InfluencerUserID string  `json:"influencer_user_id"`
	CreatorName      string  `json:"creator_name"`
	ProfileImageURL  *string `json:"profile_image_url,omitempty"`
	MatchScore       float64 `json:"match_score"`
	CampaignID       string  `json:"campaign_id"`
	CampaignTitle    string  `json:"campaign_title"`
}

type DashboardKPIRow struct {
	ActiveCampaigns              int
	TotalCampaigns               int
	RecommendedCreatorsTotal     int
	RecommendedCreatorsHighMatch int
	PendingContracts             int
	ContractsAwaitingInfluencer  int
	ContractsAwaitingPayment     int
	PendingInvites               int
	ContractsActive              int
	ContractsCompleted           int
}
