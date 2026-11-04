package contract

import (
	"errors"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrDisputeAlreadyOpen     = errors.New("this contract already has an open dispute")
	ErrDisputeNotFound        = errors.New("dispute not found")
	ErrDisputeAlreadyResolved = errors.New("dispute already resolved")
)

// Reasons a business or a creator can pick when reporting a problem.
var disputeReasons = map[string]map[string]string{
	"business": {
		"not_delivered":        "Work not delivered",
		"delivered_late":       "Delivered late",
		"not_as_agreed":        "Work not as agreed",
		"creator_unresponsive": "Creator not responding",
		"other":                "Other",
	},
	"influencer": {
		"business_unresponsive": "Business not responding",
		"extra_work_requested":  "Asked for extra work not in the contract",
		"approval_refused":      "Approval unfairly refused",
		"other":                 "Other",
	},
}

// DisputeReasonLabel returns the words shown for a reason code.
func DisputeReasonLabel(reason string) string {
	for _, m := range disputeReasons {
		if l, ok := m[reason]; ok {
			return l
		}
	}
	return strings.ReplaceAll(reason, "_", " ")
}

// DisputeSeverity: undelivered work is the most urgent for the admin.
func DisputeSeverity(reason string) string {
	if reason == "not_delivered" {
		return "high"
	}
	return "medium"
}

// OpenDisputeRequest is "Report a problem" on a contract.
type OpenDisputeRequest struct {
	Reason      string `json:"reason"`
	Description string `json:"description"`
}

func (r OpenDisputeRequest) Validate(role string) error {
	allowed := make([]interface{}, 0, len(disputeReasons[role]))
	for k := range disputeReasons[role] {
		allowed = append(allowed, k)
	}
	return validation.ValidateStruct(&r,
		validation.Field(&r.Reason, validation.Required.Error("please choose what went wrong"), validation.In(allowed...).Error("please choose one of the listed reasons")),
		validation.Field(&r.Description, validation.Required.Error("please describe what happened"),
			validation.By(func(v interface{}) error {
				n := len([]rune(strings.TrimSpace(v.(string))))
				if n < 20 || n > 2000 {
					return errors.New("please describe what happened in 20 to 2000 characters")
				}
				return nil
			})),
	)
}

// Dispute outcomes an admin can choose.
const (
	DisputeOutcomeContinue  = "continue"  // back to in progress
	DisputeOutcomeComplete  = "complete"  // work accepted, payout queued
	DisputeOutcomeTerminate = "terminate" // contract cancelled, business refunded by the team
)

type ResolveDisputeRequest struct {
	Outcome string `json:"outcome"`
	Note    string `json:"note"`
}

func (r ResolveDisputeRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Outcome, validation.Required, validation.In(DisputeOutcomeContinue, DisputeOutcomeComplete, DisputeOutcomeTerminate).Error("choose continue, complete or terminate")),
		validation.Field(&r.Note, validation.Required.Error("please explain the decision to both sides"),
			validation.By(func(v interface{}) error {
				n := len([]rune(strings.TrimSpace(v.(string))))
				if n < 10 || n > 2000 {
					return errors.New("the decision note must be 10 to 2000 characters")
				}
				return nil
			})),
	)
}

// DisputeSummary is the dispute shown on a contract (the latest one).
type DisputeSummary struct {
	ID             string     `json:"id"`
	Status         string     `json:"status"`
	Reason         string     `json:"reason"`
	ReasonLabel    string     `json:"reason_label"`
	Description    string     `json:"description"`
	OpenedByRole   string     `json:"opened_by_role"`
	Outcome        *string    `json:"outcome,omitempty"`
	ResolutionNote *string    `json:"resolution_note,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	ResolvedAt     *time.Time `json:"resolved_at,omitempty"`
}

// DisputeParties is who to email about a dispute.
type DisputeParties struct {
	ContractID    string
	ContractTitle string
	BusinessName  string
	BusinessEmail string
	CreatorName   string
	CreatorEmail  string
}
