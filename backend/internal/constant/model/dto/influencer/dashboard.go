package influencer

import "time"

type DashboardSummaryResponse struct {
	Greeting          DashboardGreeting          `json:"greeting"`
	Completeness      CompletenessResponse       `json:"completeness"`
	KPIs              DashboardKPIs              `json:"kpis"`
	ContractsSummary  DashboardContractsSummary  `json:"contracts_summary"`
	EarningsSummary   DashboardEarningsSummary   `json:"earnings_summary"`
	PendingInvites    []DashboardInviteRow       `json:"pending_invites"`
	ActiveContracts   []DashboardContractRow     `json:"active_contracts"`
}

type DashboardGreeting struct {
	CreatorName string `json:"creator_name"`
}

type DashboardKPIs struct {
	PendingInvites           int     `json:"pending_invites"`
	ActiveCollaborations     int     `json:"active_collaborations"`
	PendingContracts         int     `json:"pending_contracts"`
	CompletedContracts       int     `json:"completed_contracts"`
	ProfileScorePercent      int     `json:"profile_score_percent"`
	TotalCompletedEarningsETB float64 `json:"total_completed_earnings_etb"`
}

type DashboardContractsSummary struct {
	Pending   int `json:"pending"`
	Active    int `json:"active"`
	Completed int `json:"completed"`
}

type DashboardEarningsSummary struct {
	PendingPayoutETB     float64    `json:"pending_payout_etb"`
	CompletedEarningsETB float64    `json:"completed_earnings_etb"`
	NextPayoutAt         *time.Time `json:"next_payout_at,omitempty"`
}

type DashboardInviteRow struct {
	ID            string     `json:"id"`
	BusinessName  string     `json:"business_name"`
	CampaignTitle string     `json:"campaign_title"`
	Status        string     `json:"status"`
	BudgetSummary string     `json:"budget_summary,omitempty"`
	EndDate       *time.Time `json:"end_date,omitempty"`
}

type DashboardContractRow struct {
	ID           string `json:"id"`
	BusinessName string `json:"business_name"`
	Status       string `json:"status"`
}

type DashboardKPIRow struct {
	PendingInvites            int
	ActiveCollaborations      int
	PendingContracts          int
	CompletedContracts        int
	TotalCompletedEarningsETB float64
	PendingPayoutETB          float64
	PaidOutEarningsETB        float64
	NextPayoutAt              *time.Time
}
