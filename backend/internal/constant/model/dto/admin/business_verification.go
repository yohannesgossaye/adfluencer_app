package admin

import (
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ListBusinessVerificationsQuery struct {
	Status string
	Search string
	Page   int
	Limit  int
}

// Normalize defaults to the review queue ("pending"); "all" lists every
// business that has started verification.
func (q ListBusinessVerificationsQuery) Normalize() ListBusinessVerificationsQuery {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}
	switch q.Status {
	case "pending", "verified", "rejected", "all":
	default:
		q.Status = "pending"
	}
	return q
}

type BusinessVerificationListItem struct {
	UserID             string     `json:"user_id"`
	BusinessName       string     `json:"business_name"`
	OwnerFullName      string     `json:"owner_full_name"`
	Industry           string     `json:"industry,omitempty"`
	City               string     `json:"city,omitempty"`
	LogoURL            *string    `json:"logo_url,omitempty"`
	Email              string     `json:"email"`
	IsVerified         bool       `json:"is_verified"`
	VerificationStatus string     `json:"verification_status"`
	SubmittedAt        *time.Time `json:"submitted_at,omitempty"`
	RegisteredAt       time.Time  `json:"registered_at"`
}

type ListBusinessVerificationsResponse struct {
	Items      []BusinessVerificationListItem `json:"items"`
	Pagination commonDto.ListMeta             `json:"pagination"`
}

type BusinessSocialLinks struct {
	Instagram string `json:"instagram,omitempty"`
	TikTok    string `json:"tiktok,omitempty"`
	Facebook  string `json:"facebook,omitempty"`
	LinkedIn  string `json:"linkedin,omitempty"`
}

// BusinessVerificationDetailResponse carries everything an admin needs to
// check a business: identity (legal name, TIN), contact, location and
// online presence, plus the verification history.
type BusinessVerificationDetailResponse struct {
	UserID             string              `json:"user_id"`
	Email              string              `json:"email"`
	EmailVerified      bool                `json:"email_verified"`
	PhoneNumber        string              `json:"phone_number,omitempty"`
	PhoneVerified      bool                `json:"phone_verified"`
	Status             string              `json:"status"`
	BusinessName       string              `json:"business_name"`
	LegalBusinessName  *string             `json:"legal_business_name,omitempty"`
	OwnerFullName      string              `json:"owner_full_name"`
	TINNumber          *string             `json:"tin_number,omitempty"`
	Industry           string              `json:"industry,omitempty"`
	BusinessType       *string             `json:"business_type,omitempty"`
	BusinessSize       *string             `json:"business_size,omitempty"`
	City               string              `json:"city,omitempty"`
	Region             *string             `json:"region,omitempty"`
	Address            *string             `json:"address,omitempty"`
	WebsiteURL         *string             `json:"website_url,omitempty"`
	SocialLinks        BusinessSocialLinks `json:"social_links"`
	Description        *string             `json:"description,omitempty"`
	LogoURL            *string             `json:"logo_url,omitempty"`
	IsVerified         bool                `json:"is_verified"`
	VerificationStatus string              `json:"verification_status"`
	VerificationNote   *string             `json:"verification_note,omitempty"`
	SubmittedAt        *time.Time          `json:"submitted_at,omitempty"`
	VerifiedAt         *time.Time          `json:"verified_at,omitempty"`
	RegisteredAt       time.Time           `json:"registered_at"`
}

type VerifyBusinessRequest struct {
	AdminNote string `json:"admin_note,omitempty"`
}

type VerifyBusinessResponse struct {
	UserID             string    `json:"user_id"`
	IsVerified         bool      `json:"is_verified"`
	VerificationStatus string    `json:"verification_status"`
	VerifiedAt         time.Time `json:"verified_at"`
}

// RejectBusinessVerificationRequest is used to reject a pending business or
// revoke a verified one. The reason is shown to the business.
type RejectBusinessVerificationRequest struct {
	Reason string `json:"reason"`
}

func (r RejectBusinessVerificationRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Reason,
			validation.Required.Error("reason is required"),
			validation.Length(5, 1000).Error("reason must be between 5 and 1000 characters"),
		),
	)
}

// RevokeBusinessVerificationRequest is kept as a separate name for the revoke route.
type RevokeBusinessVerificationRequest = RejectBusinessVerificationRequest

type BusinessVerificationDecisionResponse struct {
	UserID             string `json:"user_id"`
	IsVerified         bool   `json:"is_verified"`
	VerificationStatus string `json:"verification_status"`
	Reason             string `json:"reason"`
}

type RevokeBusinessVerificationResponse = BusinessVerificationDecisionResponse

type BusinessVerificationRow struct {
	BusinessVerificationDetailResponse
}
