package contract

import (
	"errors"
	"time"
)

var (
	ErrContractNotFound      = errors.New("contract not found")
	ErrInvalidContractStatus = errors.New("invalid contract status")
	ErrInviteNotAccepted     = errors.New("invite not accepted")
	ErrContractAlreadyExists = errors.New("contract already exists for invite")
	ErrPaymentAlreadyExists  = errors.New("payment already exists")
	ErrPaymentRequired       = errors.New("payment required")
)

const (
	BusinessFeePercent   = 5.0
	InfluencerFeePercent = 5.0
	StubPaymentMethod    = "manual_stub"
)

type DeliverableInput struct {
	Type        string `json:"type"`
	Quantity    int    `json:"quantity"`
	Description string `json:"description,omitempty"`
}

type DeliverableResponse struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Quantity    int     `json:"quantity"`
	Description *string `json:"description,omitempty"`
	Status      string  `json:"status"`
}

type FeeBreakdown struct {
	AmountETB              float64 `json:"amount_etb"`
	BusinessFeePercent       float64 `json:"business_fee_percent"`
	InfluencerFeePercent     float64 `json:"influencer_fee_percent"`
	BusinessFeeAmountETB     float64 `json:"business_fee_amount_etb"`
	InfluencerFeeAmountETB   float64 `json:"influencer_fee_amount_etb"`
	TotalBusinessPayableETB  float64 `json:"total_business_payable_etb"`
	InfluencerReceivableETB  float64 `json:"influencer_receivable_etb"`
	PlatformTotalFeeETB      float64 `json:"platform_total_fee_etb"`
}

type PaymentSummary struct {
	ID                   string     `json:"id"`
	Status               string     `json:"status"`
	TotalPaidByBusiness  float64    `json:"total_paid_by_business_etb"`
	PaymentMethod        string     `json:"payment_method"`
	PaidAt               *time.Time `json:"paid_at,omitempty"`
}

type PayoutSummary struct {
	ID        string     `json:"id"`
	Status    string     `json:"status"`
	AmountETB float64    `json:"amount_etb"`
	CreatedAt time.Time  `json:"created_at"`
	PaidAt    *time.Time `json:"paid_at,omitempty"`
}

type CampaignSummary struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type PartySummary struct {
	UserID       string `json:"user_id"`
	BusinessName string `json:"business_name,omitempty"`
	CreatorName  string `json:"creator_name,omitempty"`
}

type DetailResponse struct {
	ID                     string                `json:"id"`
	Title                  string                `json:"title"`
	Description            *string               `json:"description,omitempty"`
	Status                 string                `json:"status"`
	Fees                   FeeBreakdown          `json:"fees"`
	Deliverables           []DeliverableResponse `json:"deliverables"`
	Campaign               CampaignSummary       `json:"campaign"`
	Business               PartySummary          `json:"business"`
	Influencer             PartySummary          `json:"influencer"`
	StartDate              *time.Time            `json:"start_date,omitempty"`
	EndDate                *time.Time            `json:"end_date,omitempty"`
	Payment                *PaymentSummary       `json:"payment,omitempty"`
	Payout                 *PayoutSummary        `json:"payout,omitempty"`
	AcceptedByInfluencerAt *time.Time            `json:"accepted_by_influencer_at,omitempty"`
	CompletedAt            *time.Time            `json:"completed_at,omitempty"`
	CreatedAt              time.Time             `json:"created_at"`
	UpdatedAt              time.Time             `json:"updated_at"`
}

type ListItem struct {
	ID                    string    `json:"id"`
	Title                 string    `json:"title"`
	Status                string    `json:"status"`
	AmountETB             float64   `json:"amount_etb"`
	CampaignTitle         string    `json:"campaign_title"`
	CounterpartyName      string    `json:"counterparty_name"`
	TotalBusinessPayable  float64   `json:"total_business_payable_etb,omitempty"`
	InfluencerReceivable  float64   `json:"influencer_receivable_etb,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
}

type SaveContractInput struct {
	ID                   string
	InviteID             string
	CampaignID           string
	BusinessUserID       string
	InfluencerUserID     string
	Title                string
	Description          *string
	AmountETB            float64
	Fees                 FeeBreakdown
	StartDate            *time.Time
	EndDate              *time.Time
	Deliverables         []DeliverableInput
}

type Record struct {
	DetailResponse
	InviteID         string
	BusinessUserID   string
	InfluencerUserID string
}

type AcceptedInviteRecord struct {
	ID               string
	CampaignID       string
	BusinessUserID   string
	InfluencerUserID string
	CampaignStatus   string
}

type ListQuery struct {
	Status string
	Page   int
	Limit  int
}

type PaginationResponse struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

func (q ListQuery) Normalize() ListQuery {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}
	switch q.Status {
	case "pending_influencer", "accepted", "active", "completed", "cancelled", "rejected", "disputed":
	default:
		q.Status = ""
	}
	return q
}
