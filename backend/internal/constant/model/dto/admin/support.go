package admin

import (
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	supportDto "adinfluence_service_app/internal/constant/model/dto/support"
	"strings"
	"time"
)

type ListSupportTicketsQuery struct {
	Role     string
	Status   string
	Priority string
	Search   string
	Page     int
	Limit    int
}

func (q ListSupportTicketsQuery) Normalize() ListSupportTicketsQuery {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}
	switch q.Role {
	case "", "business", "influencer":
	default:
		q.Role = ""
	}
	q.Status = strings.TrimSpace(q.Status)
	q.Priority = strings.TrimSpace(q.Priority)
	q.Search = strings.TrimSpace(q.Search)
	return q
}

type SupportTicketStats struct {
	OpenCount      int `json:"open_count"`
	HighPriority   int `json:"high_priority_count"`
	BusinessCount  int `json:"business_count"`
	InfluencerCount int `json:"influencer_count"`
	ResolvedCount  int `json:"resolved_count"`
}

type SupportTicketListItem struct {
	ID                   string    `json:"id"`
	Subject              string    `json:"subject"`
	MessagePreview       string    `json:"message_preview"`
	IssueType            string    `json:"issue_type"`
	IssueTypeLabel       string    `json:"issue_type_label"`
	Priority             string    `json:"priority"`
	Status               string    `json:"status"`
	SenderUserID         string    `json:"sender_user_id"`
	SenderName           string    `json:"sender_name"`
	SenderEmail          string    `json:"sender_email"`
	SenderRole           string    `json:"sender_role"`
	RelatedCampaignTitle *string   `json:"related_campaign_title,omitempty"`
	AdminReply           *string   `json:"admin_reply,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
}

type ListSupportTicketsResponse struct {
	Items      []SupportTicketListItem `json:"items"`
	Stats      SupportTicketStats      `json:"stats"`
	Pagination commonDto.ListMeta      `json:"pagination"`
}

type SupportTicketDetailResponse struct {
	ID                   string    `json:"id"`
	Subject              string    `json:"subject"`
	Message              string    `json:"message"`
	IssueType            string    `json:"issue_type"`
	IssueTypeLabel       string    `json:"issue_type_label"`
	Priority             string    `json:"priority"`
	Status               string    `json:"status"`
	SenderUserID         string    `json:"sender_user_id"`
	SenderName           string    `json:"sender_name"`
	SenderEmail          string    `json:"sender_email"`
	SenderRole           string    `json:"sender_role"`
	RelatedCampaignID    *string   `json:"related_campaign_id,omitempty"`
	RelatedCampaignTitle *string   `json:"related_campaign_title,omitempty"`
	AdminReply           *string   `json:"admin_reply,omitempty"`
	InternalNote         *string   `json:"internal_note,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type UpdateSupportTicketStatusRequest struct {
	Status string `json:"status"`
}

func (r UpdateSupportTicketStatusRequest) Validate() error {
	if _, err := supportDto.ValidateStatus(r.Status); err != nil {
		return err
	}
	return nil
}

type UpdateSupportTicketStatusResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type ReplySupportTicketRequest struct {
	AdminReply   string `json:"admin_reply,omitempty"`
	InternalNote string `json:"internal_note,omitempty"`
}

type ReplySupportTicketResponse struct {
	ID         string  `json:"id"`
	AdminReply *string `json:"admin_reply,omitempty"`
}

type SupportTicketRow struct {
	ID                   string
	Subject              string
	Message              string
	IssueType            string
	Priority             string
	Status               string
	ResolutionNote       *string
	InternalNote         *string
	SenderUserID         string
	SenderName           string
	SenderEmail          string
	SenderRole           string
	RelatedCampaignID    *string
	RelatedCampaignTitle *string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
