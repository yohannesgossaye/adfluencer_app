package admin

import (
	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// FundEscrowRequest is sent by an admin after checking the business's payment arrived.
type FundEscrowRequest struct {
	// Reference identifies the payment that was checked, e.g. a bank transfer number.
	Reference string `json:"reference"`
}

func (r FundEscrowRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Reference,
			validation.Required.Error("payment reference is required"),
			validation.Length(3, 200).Error("payment reference must be between 3 and 200 characters"),
		),
	)
}

type EscrowQueueResponse struct {
	Items []contractDto.EscrowQueueItem `json:"items"`
}

type FundEscrowResponse struct {
	ContractID string                     `json:"contract_id"`
	Status     string                     `json:"status"`
	Payment    contractDto.PaymentSummary `json:"payment"`
}
