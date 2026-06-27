package model

import "time"

type User struct {
	ID             string     `json:"id"`
	Email          string     `json:"email"`
	PhoneNumber    string     `json:"phone_number"`
	PasswordHash   string     `json:"password_hash"`
	Role           Role       `json:"role"`
	Status         UserStatus `json:"status"`
	EmailVerified  bool       `json:"email_verified"`
	PhoneVerified  bool       `json:"phone_verified"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`
	SuspendedUntil *time.Time `json:"suspended_until,omitempty"`
	StatusReason   *string    `json:"status_reason,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type BusinessProfile struct {
	ID                     string    `json:"id"`
	UserID                 string    `json:"user_id"`
	OwnerFullName          string    `json:"owner_full_name"`
	BusinessName           string    `json:"business_name"`
	LegalBusinessName      *string   `json:"legal_business_name,omitempty"`
	LogoURL                *string   `json:"logo_url,omitempty"`
	Description            *string   `json:"description,omitempty"`
	Industry               *string   `json:"industry,omitempty"`
	BusinessType           *string   `json:"business_type,omitempty"`
	BusinessSize           *string   `json:"business_size,omitempty"`
	City                   *string   `json:"city,omitempty"`
	Region                 *string   `json:"region,omitempty"`
	Address                *string   `json:"address,omitempty"`
	TINNumber              *string   `json:"tin_number,omitempty"`
	WebsiteURL             *string   `json:"website_url,omitempty"`
	InstagramURL           *string   `json:"instagram_url,omitempty"`
	TikTokURL              *string   `json:"tiktok_url,omitempty"`
	FacebookURL            *string   `json:"facebook_url,omitempty"`
	LinkedInURL            *string   `json:"linkedin_url,omitempty"`
	PreferredPlatforms     []string  `json:"preferred_platforms,omitempty"`
	MonthlyMarketingBudget *float64  `json:"monthly_marketing_budget,omitempty"`
	IsVerified             bool      `json:"is_verified"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type InfluencerProfile struct {
	ID                   string         `json:"id"`
	UserID               string         `json:"user_id"`
	FullName             string         `json:"full_name"`
	CreatorName          string         `json:"creator_name"`
	Username             *string        `json:"username,omitempty"`
	ProfileImageURL      *string        `json:"profile_image_url,omitempty"`
	Bio                  *string        `json:"bio,omitempty"`
	City                 *string        `json:"city,omitempty"`
	Region               *string        `json:"region,omitempty"`
	Languages            []string       `json:"languages,omitempty"`
	SocialLinks          map[string]any `json:"social_links"`
	InstagramFollowers   int            `json:"instagram_followers"`
	TikTokFollowers      int            `json:"tiktok_followers"`
	YouTubeSubscribers   int            `json:"youtube_subscribers"`
	TotalFollowers       int            `json:"total_followers"`
	EngagementRate       *float64       `json:"engagement_rate,omitempty"`
	AudienceDescription  *string        `json:"audience_description,omitempty"`
	ContentStyle         *string        `json:"content_style,omitempty"`
	PortfolioDescription *string        `json:"portfolio_description,omitempty"`
	MinPriceETB          *float64       `json:"min_price_etb,omitempty"`
	MaxPriceETB          *float64       `json:"max_price_etb,omitempty"`
	ApprovalNote         *string        `json:"approval_note,omitempty"`
	IsVerified           bool           `json:"is_verified"`
	ApprovedAt           *time.Time     `json:"approved_at,omitempty"`
	ApprovedBy           *string        `json:"approved_by,omitempty"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
}
