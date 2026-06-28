package admin

import (
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	"time"
)

type ListApplicationsQuery struct {
	Status string
	Search string
	Page   int
	Limit  int
}

type ApplicationListItem struct {
	UserID              string   `json:"user_id"`
	CreatorName         string   `json:"creator_name"`
	ProfileImageURL     *string  `json:"profile_image_url,omitempty"`
	PrimaryNiche        string   `json:"primary_niche,omitempty"`
	PrimaryPlatform     string   `json:"primary_platform,omitempty"`
	TotalFollowers      int      `json:"total_followers"`
	EngagementRate      *float64 `json:"engagement_rate,omitempty"`
	CompletenessPercent int      `json:"completeness_percent"`
	PricingSummary      string   `json:"pricing_summary,omitempty"`
	Status              string   `json:"status"`
	SubmittedAt         time.Time `json:"submitted_at"`
	BioPreview          string   `json:"bio_preview,omitempty"`
}

type ListApplicationsResponse struct {
	Items      []ApplicationListItem `json:"items"`
	Pagination commonDto.ListMeta    `json:"pagination"`
}

type ApplicationDetailResponse struct {
	UserID              string                   `json:"user_id"`
	Status              string                   `json:"status"`
	StatusReason        *string                  `json:"status_reason,omitempty"`
	Email               string                   `json:"email"`
	PhoneNumber         string                   `json:"phone_number,omitempty"`
	Profile             ApplicationProfileDetail `json:"profile"`
	CompletenessPercent int                      `json:"completeness_percent"`
	AdminNote           string                   `json:"admin_note,omitempty"`
	SubmittedAt         time.Time                `json:"submitted_at"`
}

type ApplicationProfileDetail struct {
	CreatorName         string              `json:"creator_name"`
	Bio                 *string             `json:"bio,omitempty"`
	ProfileImageURL     *string             `json:"profile_image_url,omitempty"`
	City                *string             `json:"city,omitempty"`
	Niches              []NicheResponse     `json:"niches"`
	Platforms           []string            `json:"platforms"`
	SocialLinks         map[string]string   `json:"social_links"`
	TotalFollowers      int                 `json:"total_followers"`
	EngagementRate      *float64            `json:"engagement_rate,omitempty"`
	MinPriceETB         *float64            `json:"min_price_etb,omitempty"`
	MaxPriceETB         *float64            `json:"max_price_etb,omitempty"`
	AudienceDescription *string             `json:"audience_description,omitempty"`
	Languages           []string            `json:"languages"`
	PortfolioItems      []PortfolioItem     `json:"portfolio_items"`
}

type PortfolioItem struct {
	ID            string    `json:"id"`
	FileURL       string    `json:"file_url"`
	MediaType     string    `json:"media_type"`
	MimeType      string    `json:"mime_type"`
	FileSizeBytes int64     `json:"file_size_bytes"`
	SortOrder     int       `json:"sort_order"`
	CreatedAt     time.Time `json:"created_at"`
}

type NicheResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type ApproveApplicationRequest struct {
	AdminNote string `json:"admin_note,omitempty"`
}

type RejectApplicationRequest struct {
	Reason    string `json:"reason"`
	AdminNote string `json:"admin_note,omitempty"`
}

type UpdateAdminNotesRequest struct {
	AdminNote string `json:"admin_note"`
}

type ApproveApplicationResponse struct {
	UserID     string    `json:"user_id"`
	Status     string    `json:"status"`
	ApprovedAt time.Time `json:"approved_at"`
}

type RejectApplicationResponse struct {
	UserID string `json:"user_id"`
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type ApplicationRow struct {
	UserID          string
	Status          string
	StatusReason    *string
	Email           string
	PhoneNumber     string
	UpdatedAt       time.Time
	CreatorName     string
	ProfileImageURL *string
	Bio             *string
	City            *string
	TotalFollowers  int
	EngagementRate  *float64
	MinPriceETB     *float64
	MaxPriceETB     *float64
	AudienceDesc    *string
	ApprovalNote    *string
	PrimaryNiche    string
	PrimaryPlatform string
}
