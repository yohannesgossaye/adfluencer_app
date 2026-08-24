package influencer

import (
	"errors"
	"time"
)

var (
	ErrInviteNotFound      = errors.New("invite not found")
	ErrInvalidInviteStatus = errors.New("invalid invite status")
)

type ListInvitesQuery struct {
	Status string
}

type InviteCampaignSummary struct {
	ID             string     `json:"id"`
	BusinessUserID string     `json:"business_user_id,omitempty"`
	Title          string     `json:"title"`
	BusinessName   string     `json:"business_name"`
	Status         string     `json:"status"`
	BudgetMinETB   *float64   `json:"budget_min_etb,omitempty"`
	BudgetMaxETB   *float64   `json:"budget_max_etb,omitempty"`
	Platforms      []string   `json:"platforms"`
	StartDate      *time.Time `json:"start_date,omitempty"`
	EndDate        *time.Time `json:"end_date,omitempty"`
}

type InviteListItem struct {
	ID             string                `json:"id"`
	BusinessUserID string                `json:"business_user_id,omitempty"`
	Status         string                `json:"status"`
	Message        *string               `json:"message,omitempty"`
	Campaign       InviteCampaignSummary `json:"campaign"`
	CreatedAt      time.Time             `json:"created_at"`
	RespondedAt    *time.Time            `json:"responded_at,omitempty"`
}

type ListInvitesResponse struct {
	Items []InviteListItem `json:"items"`
}

type InviteDetailResponse struct {
	ID          string                `json:"id"`
	Status      string                `json:"status"`
	Message     *string               `json:"message,omitempty"`
	Campaign    InviteCampaignSummary `json:"campaign"`
	CreatedAt   time.Time             `json:"created_at"`
	RespondedAt *time.Time            `json:"responded_at,omitempty"`
}

type RespondInviteResponse struct {
	ID          string    `json:"id"`
	Status      string    `json:"status"`
	RespondedAt time.Time `json:"responded_at"`
}

type InviteRecord struct {
	ID               string
	CampaignID       string
	BusinessUserID   string
	InfluencerUserID string
	Status           string
	Message          *string
	CreatedAt        time.Time
	RespondedAt      *time.Time
	Campaign         InviteCampaignSummary
}

func (q ListInvitesQuery) Normalize() ListInvitesQuery {
	switch q.Status {
	case "pending", "accepted", "rejected", "cancelled":
	default:
		q.Status = ""
	}
	return q
}
