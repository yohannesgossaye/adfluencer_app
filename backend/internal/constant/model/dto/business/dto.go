package business

import (
	"errors"
	"time"
)

var (
	ErrProfileNotFound           = errors.New("business profile not found")
	ErrProfileIncomplete         = errors.New("profile incomplete")
	ErrInvalidStatusForComplete  = errors.New("invalid status for complete")
	ErrProfileLocked             = errors.New("profile locked")
	ErrNotVerified               = errors.New("account not verified")
)

type SocialLinks struct {
	Instagram string `json:"instagram,omitempty"`
	TikTok    string `json:"tiktok,omitempty"`
	Facebook  string `json:"facebook,omitempty"`
	LinkedIn  string `json:"linkedin,omitempty"`
}

type ProfileUserResponse struct {
	ID            string  `json:"id"`
	Email         string  `json:"email"`
	PhoneNumber   string  `json:"phone_number,omitempty"`
	Status        string  `json:"status"`
	StatusReason  *string `json:"status_reason,omitempty"`
	EmailVerified bool    `json:"email_verified"`
	PhoneVerified bool    `json:"phone_verified"`
}

type ProfileDetailResponse struct {
	ID                     string      `json:"id"`
	OwnerFullName          string      `json:"owner_full_name"`
	BusinessName           string      `json:"business_name"`
	LegalBusinessName      *string     `json:"legal_business_name,omitempty"`
	LogoURL                *string     `json:"logo_url,omitempty"`
	Description            *string     `json:"description,omitempty"`
	Industry               *string     `json:"industry,omitempty"`
	BusinessType           *string     `json:"business_type,omitempty"`
	BusinessSize           *string     `json:"business_size,omitempty"`
	City                   *string     `json:"city,omitempty"`
	Region                 *string     `json:"region,omitempty"`
	Address                *string     `json:"address,omitempty"`
	TINNumber              *string     `json:"tin_number,omitempty"`
	WebsiteURL             *string     `json:"website_url,omitempty"`
	SocialLinks            SocialLinks `json:"social_links"`
	PreferredPlatforms     []string    `json:"preferred_platforms"`
	MonthlyMarketingBudget *float64    `json:"monthly_marketing_budget,omitempty"`
	IsVerified             bool        `json:"is_verified"`
}

type CompletenessResponse struct {
	Percent     int      `json:"percent"`
	Missing     []string `json:"missing"`
	CanComplete bool     `json:"can_complete"`
}

type MarketplaceAccessResponse struct {
	VisibleInMarketplace bool   `json:"visible_in_marketplace"`
	CanCreateContracts   bool   `json:"can_create_contracts"`
	Message              string `json:"message,omitempty"`
}

type GetProfileResponse struct {
	User               ProfileUserResponse       `json:"user"`
	Profile            ProfileDetailResponse     `json:"profile"`
	Completeness       CompletenessResponse      `json:"completeness"`
	MarketplaceAccess  MarketplaceAccessResponse `json:"marketplace_access"`
}

type UpdateProfileRequest struct {
	OwnerFullName          string      `json:"owner_full_name"`
	BusinessName           string      `json:"business_name"`
	LegalBusinessName      string      `json:"legal_business_name,omitempty"`
	Description            string      `json:"description"`
	Industry               string      `json:"industry"`
	BusinessType           string      `json:"business_type,omitempty"`
	BusinessSize           string      `json:"business_size,omitempty"`
	City                   string      `json:"city"`
	Region                 string      `json:"region,omitempty"`
	Address                string      `json:"address,omitempty"`
	TINNumber              string      `json:"tin_number,omitempty"`
	WebsiteURL             string      `json:"website_url,omitempty"`
	SocialLinks            SocialLinks `json:"social_links"`
	PreferredPlatforms     []string    `json:"preferred_platforms"`
	MonthlyMarketingBudget *float64    `json:"monthly_marketing_budget,omitempty"`
}

type UploadLogoResponse struct {
	LogoURL string `json:"logo_url"`
}

type CompleteProfileResponse struct {
	Status      string    `json:"status"`
	Message     string    `json:"message"`
	IsVerified  bool      `json:"is_verified"`
	CompletedAt time.Time `json:"completed_at"`
}

type UpdateProfileInput struct {
	UserID                 string
	OwnerFullName          string
	BusinessName           string
	LegalBusinessName      *string
	Description            *string
	Industry               *string
	BusinessType           *string
	BusinessSize           *string
	City                   *string
	Region                 *string
	Address                *string
	TINNumber              *string
	WebsiteURL             *string
	InstagramURL           *string
	TikTokURL              *string
	FacebookURL            *string
	LinkedInURL            *string
	PreferredPlatforms     []string
	MonthlyMarketingBudget *float64
}

type ProfileRecord struct {
	ProfileDetailResponse
	UserID string
}
