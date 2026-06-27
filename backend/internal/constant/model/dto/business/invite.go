package business

import (
	"errors"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrInviteNotFound       = errors.New("invite not found")
	ErrInviteAlreadyExists  = errors.New("invite already exists")
	ErrInvalidInviteStatus  = errors.New("invalid invite status")
)

type CreateInviteRequest struct {
	InfluencerUserID string `json:"influencer_user_id"`
	Message          string `json:"message,omitempty"`
}

func (r CreateInviteRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.InfluencerUserID,
			validation.Required.Error("influencer_user_id is required"),
		),
		validation.Field(&r.Message,
			validation.Length(0, 2000).Error("message must be at most 2000 characters"),
		),
	)
}

type CreateInviteInput struct {
	ID               string
	CampaignID       string
	BusinessUserID   string
	InfluencerUserID string
	Message          *string
}

type InviteResponse struct {
	ID               string    `json:"id"`
	CampaignID       string    `json:"campaign_id"`
	InfluencerUserID string    `json:"influencer_user_id"`
	CreatorName      string    `json:"creator_name,omitempty"`
	Status           string    `json:"status"`
	Message          *string   `json:"message,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

type ListInvitesResponse struct {
	Items []InviteResponse `json:"items"`
}

type CreateInviteResponse struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func OptionalInviteMessage(v string) *string {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	return &v
}
