package admin

import (
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	customErr "adinfluence_service_app/internal/constant/errors"
	"strings"
	"time"
)

type ListAdminDisputesQuery struct {
	Tab    string
	Search string
	Page   int
	Limit  int
}

func (q ListAdminDisputesQuery) Normalize() ListAdminDisputesQuery {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}
	switch strings.ToLower(strings.TrimSpace(q.Tab)) {
	case "", "open", "resolved":
		q.Tab = strings.ToLower(strings.TrimSpace(q.Tab))
		if q.Tab == "" {
			q.Tab = "open"
		}
	default:
		q.Tab = "open"
	}
	q.Search = strings.TrimSpace(q.Search)
	return q
}

type AdminDisputeListItem struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	BusinessName    string    `json:"business_name"`
	InfluencerName  string    `json:"influencer_name"`
	CampaignTitle   string    `json:"campaign_title"`
	IssueType       string    `json:"issue_type"`
	IssueTypeLabel  string    `json:"issue_type_label"`
	Severity        string    `json:"severity"`
	Status          string    `json:"status"`
	Preview         string    `json:"preview"`
	CreatedAt       time.Time `json:"created_at"`
}

type ListAdminDisputesResponse struct {
	Items      []AdminDisputeListItem `json:"items"`
	Pagination commonDto.ListMeta     `json:"pagination"`
}

type AdminDisputeDetailResponse struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	BusinessName     string    `json:"business_name"`
	InfluencerName   string    `json:"influencer_name"`
	InfluencerEmail  string    `json:"influencer_email"`
	InfluencerUserID string    `json:"influencer_user_id"`
	CampaignTitle    string    `json:"campaign_title"`
	CampaignID       string    `json:"campaign_id"`
	ContractID       string    `json:"contract_id"`
	IssueType        string    `json:"issue_type"`
	IssueTypeLabel   string    `json:"issue_type_label"`
	Severity         string    `json:"severity"`
	Status           string    `json:"status"`
	Description      string    `json:"description"`
	Preview          string    `json:"preview"`
	AdminNote        string    `json:"admin_note,omitempty"`
	ResolutionNote   *string   `json:"resolution_note,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	ResolvedAt       *time.Time `json:"resolved_at,omitempty"`
}

type UpdateDisputeStatusRequest struct {
	Status string `json:"status"`
}

func (r UpdateDisputeStatusRequest) Validate() error {
	switch strings.ToLower(strings.TrimSpace(r.Status)) {
	case "under_review", "resolved", "escalated", "closed", "open":
		return nil
	default:
		return customErr.ErrInvalidDisputeStatus
	}
}

type UpdateDisputeStatusResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type UpdateDisputeNotesRequest struct {
	AdminNote      *string `json:"admin_note,omitempty"`
	ResolutionNote *string `json:"resolution_note,omitempty"`
}

type UpdateDisputeNotesResponse struct {
	ID             string  `json:"id"`
	AdminNote      string  `json:"admin_note,omitempty"`
	ResolutionNote *string `json:"resolution_note,omitempty"`
}

type DisputeNotesSnapshot struct {
	AdminNote      string
	ResolutionNote *string
}

type DisputeRow struct {
	ID               string
	ContractID       string
	CampaignID       string
	CampaignTitle    string
	BusinessName     string
	InfluencerName   string
	InfluencerEmail  string
	InfluencerUserID string
	AgainstUserID    string
	IssueType        string
	Severity         string
	Status           string
	Description      string
	ResolutionNote   *string
	AdminNote        string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ResolvedAt       *time.Time
}
