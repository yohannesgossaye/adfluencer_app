package business

import (
	supportDto "adinfluence_service_app/internal/constant/model/dto/support"
	"strings"
	"time"
)

type CreateSupportTicketRequest struct {
	IssueType         string  `json:"issue_type"`
	Subject           string  `json:"subject"`
	Message           string  `json:"message"`
	Priority          string  `json:"priority"`
	RelatedCampaignID *string `json:"related_campaign_id,omitempty"`
}

func (r CreateSupportTicketRequest) Validate() error {
	if _, err := supportDto.ValidateBusinessIssueType(r.IssueType); err != nil {
		return err
	}
	if strings.TrimSpace(r.Subject) == "" {
		return supportDto.ErrSubjectRequired
	}
	if strings.TrimSpace(r.Message) == "" {
		return supportDto.ErrMessageRequired
	}
	if _, err := supportDto.NormalizePriority(r.Priority); err != nil {
		return err
	}
	return nil
}

type SupportCampaignOption struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type ListSupportCampaignOptionsResponse struct {
	Items []SupportCampaignOption `json:"items"`
}

type SupportTicketListItem struct {
	ID                   string    `json:"id"`
	Subject              string    `json:"subject"`
	IssueType            string    `json:"issue_type"`
	IssueTypeLabel       string    `json:"issue_type_label"`
	Priority             string    `json:"priority"`
	Status               string    `json:"status"`
	AdminReply           *string   `json:"admin_reply,omitempty"`
	RelatedCampaignID    *string   `json:"related_campaign_id,omitempty"`
	RelatedCampaignTitle *string   `json:"related_campaign_title,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
}

type ListSupportTicketsResponse struct {
	Items []SupportTicketListItem `json:"items"`
}

type SupportTicketDetailResponse struct {
	SupportTicketListItem
	Message string `json:"message"`
}

type CreateSupportTicketResponse struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateSupportTicketInput struct {
	ID                string
	UserID            string
	IssueType         string
	Subject           string
	Message           string
	Priority          string
	RelatedCampaignID *string
}

type SupportTicketRecord struct {
	ID                   string
	Subject              string
	Message              string
	IssueType            string
	Priority             string
	Status               string
	ResolutionNote       *string
	RelatedCampaignID    *string
	RelatedCampaignTitle *string
	CreatedAt            time.Time
}
