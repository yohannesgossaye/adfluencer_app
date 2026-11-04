package support

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"errors"
	"strings"
	"time"
)

const (
	IssueCampaignIssue        = "campaign_issue"
	IssueContractIssue        = "contract_issue"
	IssuePaymentQuestion      = "payment_question"
	IssueCreatorCommunication = "creator_communication"
	IssueTechnicalIssue       = "technical_issue"
	IssueGeneralFeedback      = "general_feedback"
	IssueProfileApproval      = "profile_approval"
	IssueComplaint            = "complaint"

	PriorityLow    = "low"
	PriorityNormal = "normal"
	PriorityHigh   = "high"
	PriorityUrgent = "urgent"

	StatusOpen      = "open"
	StatusInReview  = "in_review"
	StatusResolved  = "resolved"
	StatusEscalated = "escalated"
	StatusClosed    = "closed"
)

var businessIssueTypes = map[string]string{
	IssueCampaignIssue:        "Campaign issue",
	IssueContractIssue:        "Contract issue",
	IssuePaymentQuestion:      "Payment question",
	IssueCreatorCommunication: "Creator communication",
	IssueTechnicalIssue:       "Technical issue",
	IssueGeneralFeedback:      "General feedback",
}

var influencerIssueTypes = map[string]string{
	IssueProfileApproval: "Profile approval",
	IssueCampaignIssue:   "Campaign invite",
	IssueContractIssue:   "Contract issue",
	IssuePaymentQuestion: "Payment/payout question",
	IssueComplaint:       "Business communication",
	IssueTechnicalIssue:  "Technical issue",
	IssueGeneralFeedback: "General feedback",
}

var ErrInvalidIssueType = errors.New("invalid issue type")
var ErrInvalidPriority = errors.New("invalid priority")
var ErrInvalidStatus = errors.New("invalid support ticket status")
var ErrSubjectRequired = errors.New("subject is required")
var ErrMessageRequired = errors.New("message is required")

func NormalizePriority(raw string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", PriorityNormal, "medium":
		return PriorityNormal, nil
	case PriorityLow:
		return PriorityLow, nil
	case PriorityHigh:
		return PriorityHigh, nil
	case PriorityUrgent:
		return PriorityUrgent, nil
	default:
		return "", ErrInvalidPriority
	}
}

func ValidateBusinessIssueType(raw string) (string, error) {
	return validateIssueType(raw, businessIssueTypes)
}

func ValidateInfluencerIssueType(raw string) (string, error) {
	return validateIssueType(raw, influencerIssueTypes)
}

func validateIssueType(raw string, allowed map[string]string) (string, error) {
	v := strings.ToLower(strings.TrimSpace(raw))
	if _, ok := allowed[v]; ok {
		return v, nil
	}
	return "", ErrInvalidIssueType
}

func ValidateStatus(raw string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case StatusOpen, StatusInReview, StatusResolved, StatusEscalated, StatusClosed:
		return strings.ToLower(strings.TrimSpace(raw)), nil
	default:
		return "", ErrInvalidStatus
	}
}

func IssueTypeLabel(issueType string) string {
	if label, ok := businessIssueTypes[issueType]; ok {
		return label
	}
	if label, ok := influencerIssueTypes[issueType]; ok {
		return label
	}
	return strings.ReplaceAll(issueType, "_", " ")
}

func MapValidationError(err error) error {
	switch {
	case errors.Is(err, ErrInvalidIssueType):
		return customErr.ErrBadRequest
	case errors.Is(err, ErrInvalidPriority):
		return customErr.ErrBadRequest
	case errors.Is(err, ErrInvalidStatus):
		return customErr.ErrBadRequest
	case errors.Is(err, ErrSubjectRequired), errors.Is(err, ErrMessageRequired), errors.Is(err, ErrMessageTooLong):
		return customErr.ErrBadRequest
	default:
		return err
	}
}

// Who sent a message in a support conversation.
const (
	FromUser    = "user"
	FromSupport = "support"
)

// Message is one public message in a support conversation, after the
// opening message stored on the ticket itself.
type Message struct {
	ID         string    `json:"id"`
	From       string    `json:"from"` // "user" or "support"
	SenderName string    `json:"sender_name"`
	Body       string    `json:"body"`
	CreatedAt  time.Time `json:"created_at"`
}

const MaxMessageLength = 4000

var ErrMessageTooLong = errors.New("message must be at most 4000 characters")

// ReplyRequest is a new message from the user on their own ticket.
type ReplyRequest struct {
	Message string `json:"message"`
}

func (r ReplyRequest) Validate() error {
	n := len([]rune(strings.TrimSpace(r.Message)))
	if n == 0 {
		return ErrMessageRequired
	}
	if n > MaxMessageLength {
		return ErrMessageTooLong
	}
	return nil
}

// ReplyResult tells the user whether their reply reopened a finished ticket.
type ReplyResult struct {
	Message  Message `json:"message"`
	Status   string  `json:"status"`
	Reopened bool    `json:"reopened"`
}
