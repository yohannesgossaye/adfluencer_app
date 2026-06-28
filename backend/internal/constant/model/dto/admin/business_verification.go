package admin

import (
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	"time"
)

type ListBusinessVerificationsQuery struct {
	Status string
	Search string
	Page   int
	Limit  int
}

func (q ListBusinessVerificationsQuery) Normalize() ListBusinessVerificationsQuery {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}
	switch q.Status {
	case "pending", "verified":
	default:
		q.Status = "pending"
	}
	return q
}

type BusinessVerificationListItem struct {
	UserID        string    `json:"user_id"`
	BusinessName  string    `json:"business_name"`
	OwnerFullName string    `json:"owner_full_name"`
	Industry      string    `json:"industry,omitempty"`
	City          string    `json:"city,omitempty"`
	LogoURL       *string   `json:"logo_url,omitempty"`
	Email         string    `json:"email"`
	IsVerified    bool      `json:"is_verified"`
	RegisteredAt  time.Time `json:"registered_at"`
}

type ListBusinessVerificationsResponse struct {
	Items      []BusinessVerificationListItem `json:"items"`
	Pagination commonDto.ListMeta             `json:"pagination"`
}

type BusinessVerificationDetailResponse struct {
	UserID        string    `json:"user_id"`
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phone_number,omitempty"`
	Status        string    `json:"status"`
	BusinessName  string    `json:"business_name"`
	OwnerFullName string    `json:"owner_full_name"`
	Industry      string    `json:"industry,omitempty"`
	City          string    `json:"city,omitempty"`
	Description   *string   `json:"description,omitempty"`
	LogoURL       *string   `json:"logo_url,omitempty"`
	IsVerified    bool      `json:"is_verified"`
	RegisteredAt  time.Time `json:"registered_at"`
}

type VerifyBusinessRequest struct {
	AdminNote string `json:"admin_note,omitempty"`
}

type VerifyBusinessResponse struct {
	UserID     string    `json:"user_id"`
	IsVerified bool      `json:"is_verified"`
	VerifiedAt time.Time `json:"verified_at"`
}

type RevokeBusinessVerificationRequest struct {
	Reason string `json:"reason,omitempty"`
}

type RevokeBusinessVerificationResponse struct {
	UserID     string `json:"user_id"`
	IsVerified bool   `json:"is_verified"`
}

type BusinessVerificationRow struct {
	UserID        string
	Email         string
	PhoneNumber   string
	UserStatus    string
	BusinessName  string
	OwnerFullName string
	Industry      *string
	City          *string
	Description   *string
	LogoURL       *string
	IsVerified    bool
	RegisteredAt  time.Time
}
