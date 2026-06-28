package business

import (
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateContractRequest struct {
	InviteID     string                        `json:"invite_id"`
	Title        string                        `json:"title"`
	AmountETB    float64                       `json:"amount_etb"`
	Description  string                        `json:"description,omitempty"`
	Deliverables []contractDto.DeliverableInput `json:"deliverables"`
	StartDate    string                        `json:"start_date,omitempty"`
	EndDate      string                        `json:"end_date,omitempty"`
}

func (r CreateContractRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.InviteID, validation.Required.Error("invite_id is required")),
		validation.Field(&r.Title,
			validation.Required.Error("title is required"),
			validation.Length(3, 255).Error("title must be between 3 and 255 characters"),
		),
		validation.Field(&r.AmountETB,
			validation.Required.Error("amount_etb is required"),
			validation.Min(1.0).Error("amount_etb must be at least 1"),
		),
		validation.Field(&r.Deliverables,
			validation.Required.Error("at least one deliverable is required"),
			validation.Length(1, 20).Error("select between 1 and 20 deliverables"),
		),
	)
}

type CreateContractResponse struct {
	Contract contractDto.DetailResponse `json:"contract"`
}

type ListContractsResponse struct {
	Items      []contractDto.ListItem       `json:"items"`
	Pagination commonDto.ListMeta         `json:"pagination"`
}

type PayContractResponse struct {
	Contract contractDto.DetailResponse `json:"contract"`
	Payment  contractDto.PaymentSummary `json:"payment"`
}

type ApproveWorkResponse struct {
	Contract contractDto.DetailResponse `json:"contract"`
	Payout   contractDto.PayoutSummary  `json:"payout"`
}

type ContractMutationResponse struct {
	Contract contractDto.DetailResponse `json:"contract"`
}

func ParseContractDate(raw string) (*time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return nil, err
	}
	utc := t.UTC()
	return &utc, nil
}
