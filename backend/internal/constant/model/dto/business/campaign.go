package business

import (
	"errors"
	"time"
)

var (
	ErrCampaignNotFound       = errors.New("campaign not found")
	ErrInvalidCampaignStatus  = errors.New("invalid campaign status")
	ErrAccountNotActive       = errors.New("account not active")
	ErrUnknownCategory        = errors.New("unknown category")
)

type CampaignDeliverableInput struct {
	Type        string `json:"type"`
	Quantity    int    `json:"quantity"`
	Description string `json:"description,omitempty"`
}

type CampaignDeliverableResponse struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Quantity    int     `json:"quantity"`
	Description *string `json:"description,omitempty"`
}

type CreateCampaignRequest struct {
	Title                  string                     `json:"title"`
	Goal                   string                     `json:"goal"`
	Description            string                     `json:"description"`
	Niche                  string                     `json:"niche"`
	TargetAudience         string                     `json:"target_audience,omitempty"`
	City                   string                     `json:"city"`
	Region                 string                     `json:"region,omitempty"`
	Platforms              []string                   `json:"platforms"`
	BudgetMinETB             *float64                   `json:"budget_min_etb,omitempty"`
	BudgetMaxETB             *float64                   `json:"budget_max_etb,omitempty"`
	CreatorSizePreference    string                     `json:"creator_size_preference"`
	NumberOfCreatorsNeeded   int                        `json:"number_of_creators_needed"`
	Deliverables             []CampaignDeliverableInput `json:"deliverables"`
	StartDate                string                     `json:"start_date"`
	EndDate                  string                     `json:"end_date"`
	AdditionalNotes          string                     `json:"additional_notes,omitempty"`
}

type UpdateCampaignRequest struct {
	CreateCampaignRequest
}

type ListCampaignsQuery struct {
	Status string
	Page   int
	Limit  int
}

type PaginationResponse struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type CampaignListItem struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	BusinessName    string     `json:"business_name"`
	PrimaryPlatform string     `json:"primary_platform,omitempty"`
	BudgetSummary   string     `json:"budget_summary,omitempty"`
	Status          string     `json:"status"`
	StartDate       *time.Time `json:"start_date,omitempty"`
	EndDate         *time.Time `json:"end_date,omitempty"`
	CreatorsNeeded  int        `json:"creators_needed"`
	InvitesSent     int        `json:"invites_sent"`
	InvitesAccepted int        `json:"invites_accepted"`
	CreatedAt       time.Time  `json:"created_at"`
}

type ListCampaignsResponse struct {
	Items      []CampaignListItem `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}

type CampaignCategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type CampaignDetailResponse struct {
	ID                     string                        `json:"id"`
	Title                  string                        `json:"title"`
	Goal                   *string                       `json:"goal,omitempty"`
	Description            string                        `json:"description"`
	Category               *CampaignCategoryResponse     `json:"category,omitempty"`
	TargetAudience         *string                       `json:"target_audience,omitempty"`
	City                   *string                       `json:"city,omitempty"`
	Region                 *string                       `json:"region,omitempty"`
	Platforms              []string                      `json:"platforms"`
	BudgetMinETB           *float64                      `json:"budget_min_etb,omitempty"`
	BudgetMaxETB           *float64                      `json:"budget_max_etb,omitempty"`
	CreatorSizePreference  string                        `json:"creator_size_preference"`
	NumberOfCreatorsNeeded int                           `json:"number_of_creators_needed"`
	Deliverables           []CampaignDeliverableResponse `json:"deliverables"`
	StartDate              *time.Time                    `json:"start_date,omitempty"`
	EndDate                *time.Time                    `json:"end_date,omitempty"`
	Status                 string                        `json:"status"`
	BusinessName           string                        `json:"business_name"`
	InvitesSent            int                           `json:"invites_sent"`
	InvitesAccepted        int                           `json:"invites_accepted"`
	CreatedAt              time.Time                     `json:"created_at"`
	UpdatedAt              time.Time                     `json:"updated_at"`
}

type CampaignMutationResponse struct {
	Campaign CampaignDetailResponse `json:"campaign"`
}

type SaveCampaignInput struct {
	ID                     string
	BusinessUserID         string
	Title                  string
	Goal                   *string
	Description            string
	CategoryID             *string
	TargetAudience         *string
	City                   *string
	Region                 *string
	Platforms              []string
	BudgetMinETB           *float64
	BudgetMaxETB           *float64
	CreatorSizePreference  string
	NumberOfCreatorsNeeded int
	Deliverables           []CampaignDeliverableInput
	StartDate              *time.Time
	EndDate                *time.Time
}

type CampaignRecord struct {
	CampaignDetailResponse
	BusinessUserID string
}

func (q ListCampaignsQuery) Normalize() ListCampaignsQuery {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}
	q.Status = normalizeStatusFilter(q.Status)
	return q
}

func normalizeStatusFilter(status string) string {
	switch status {
	case "draft", "pending", "active", "completed", "cancelled", "paused":
		return status
	default:
		return ""
	}
}
