package admin

import (
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	"errors"
	"strings"
	"time"
)

var ErrInvalidAdminCampaignAction = errors.New("invalid admin campaign action")

type ListAdminCampaignsQuery struct {
	Status   string
	Niche    string
	Business string
	Page     int
	Limit    int
}

func (q ListAdminCampaignsQuery) Normalize() ListAdminCampaignsQuery {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}
	q.Status = normalizeAdminCampaignStatusFilter(q.Status)
	q.Niche = strings.TrimSpace(q.Niche)
	q.Business = strings.TrimSpace(q.Business)
	return q
}

func normalizeAdminCampaignStatusFilter(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "", "draft", "pending", "active", "completed", "cancelled", "paused", "flagged":
		return strings.ToLower(strings.TrimSpace(status))
	default:
		return ""
	}
}

func AdminCampaignStatusFilterDB(status string) string {
	if status == "flagged" {
		return "paused"
	}
	return status
}

type AdminCampaignStats struct {
	Active  int `json:"active"`
	Pending int `json:"pending"`
	Flagged int `json:"flagged"`
}

type AdminCampaignListItem struct {
	ID                 string    `json:"id"`
	Title              string    `json:"title"`
	BusinessName       string    `json:"business_name"`
	Niche              string    `json:"niche,omitempty"`
	BudgetSummary      string    `json:"budget_summary,omitempty"`
	CreatorsInvited    int       `json:"creators_invited"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	ContractSummary    string    `json:"contract_summary"`
	PaymentSummary     string    `json:"payment_summary"`
	ContractsActive    int       `json:"contracts_active"`
	ContractsPending   int       `json:"contracts_pending"`
}

type ListAdminCampaignsResponse struct {
	Items      []AdminCampaignListItem `json:"items"`
	Stats      AdminCampaignStats      `json:"stats"`
	Pagination commonDto.ListMeta      `json:"pagination"`
}

type AdminCampaignDetailResponse struct {
	AdminCampaignListItem
	Goal               *string   `json:"goal,omitempty"`
	Description        string    `json:"description"`
	BusinessUserID     string    `json:"business_user_id"`
	BudgetMinETB       *float64  `json:"budget_min_etb,omitempty"`
	BudgetMaxETB       *float64  `json:"budget_max_etb,omitempty"`
	CreatorsNeeded     int       `json:"creators_needed"`
	InvitesAccepted    int       `json:"invites_accepted"`
	StartDate          *time.Time `json:"start_date,omitempty"`
	EndDate            *time.Time `json:"end_date,omitempty"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type ModerateCampaignResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type AdminCampaignRow struct {
	ID               string
	Title            string
	BusinessName     string
	BusinessUserID   string
	Niche            string
	BudgetMinETB     *float64
	BudgetMaxETB     *float64
	Status           string
	CreatorsNeeded   int
	CreatorsInvited  int
	InvitesAccepted  int
	ContractsActive  int
	ContractsPending int
	PaymentsPaid     int
	PaymentsPending  int
	Goal             *string
	Description      string
	StartDate          *time.Time
	EndDate            *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func MapAdminCampaignStatus(status string) string {
	if status == "paused" {
		return "flagged"
	}
	return status
}
