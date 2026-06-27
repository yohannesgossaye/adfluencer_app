package influencer

import (
	"errors"
	"time"
)

var (
	ErrProfileNotFound          = errors.New("influencer profile not found")
	ErrProfileIncomplete        = errors.New("profile incomplete")
	ErrInvalidStatusForSubmit   = errors.New("invalid status for submit")
	ErrProfileLocked            = errors.New("profile locked")
	ErrForbidden                = errors.New("forbidden")
	ErrUnknownCategory          = errors.New("unknown category")
	ErrInvalidPlatform          = errors.New("invalid platform")
	ErrPortfolioLimitReached    = errors.New("portfolio limit reached")
	ErrPortfolioItemNotFound    = errors.New("portfolio item not found")
)

type SocialLinks struct {
	TikTok    string `json:"tiktok,omitempty"`
	Instagram string `json:"instagram,omitempty"`
	YouTube   string `json:"youtube,omitempty"`
}

type CategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
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
	ID                  string           `json:"id"`
	CreatorName         string           `json:"creator_name"`
	FullName            string           `json:"full_name"`
	Username            *string          `json:"username,omitempty"`
	Bio                 *string          `json:"bio,omitempty"`
	ProfileImageURL     *string          `json:"profile_image_url,omitempty"`
	City                *string          `json:"city,omitempty"`
	Region              *string          `json:"region,omitempty"`
	Languages           []string         `json:"languages"`
	Niches              []CategoryResponse `json:"niches"`
	Platforms           []string         `json:"platforms"`
	SocialLinks         SocialLinks      `json:"social_links"`
	InstagramFollowers  int              `json:"instagram_followers"`
	TikTokFollowers     int              `json:"tiktok_followers"`
	YouTubeSubscribers  int              `json:"youtube_subscribers"`
	TotalFollowers      int              `json:"total_followers"`
	EngagementRate      *float64         `json:"engagement_rate,omitempty"`
	MinPriceETB         *float64         `json:"min_price_etb,omitempty"`
	MaxPriceETB         *float64         `json:"max_price_etb,omitempty"`
	AudienceDescription *string          `json:"audience_description,omitempty"`
	IsVerified          bool             `json:"is_verified"`
	ApprovedAt          *time.Time       `json:"approved_at,omitempty"`
	PortfolioItems      []PortfolioItem  `json:"portfolio_items"`
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

type PortfolioItemRecord struct {
	PortfolioItem
	ObjectKey string
}

type CompletenessResponse struct {
	Percent    int      `json:"percent"`
	Missing    []string `json:"missing"`
	CanSubmit  bool     `json:"can_submit"`
}

type GetProfileResponse struct {
	User         ProfileUserResponse   `json:"user"`
	Profile      ProfileDetailResponse `json:"profile"`
	Completeness CompletenessResponse  `json:"completeness"`
}

type UpdateProfileRequest struct {
	CreatorName         string      `json:"creator_name"`
	Bio                 string      `json:"bio"`
	Niches              []string    `json:"niches"`
	Platforms           []string    `json:"platforms"`
	SocialLinks         SocialLinks `json:"social_links"`
	Followers           *int        `json:"followers,omitempty"`
	InstagramFollowers  *int        `json:"instagram_followers,omitempty"`
	TikTokFollowers     *int        `json:"tiktok_followers,omitempty"`
	YouTubeSubscribers  *int        `json:"youtube_subscribers,omitempty"`
	EngagementRate      *float64    `json:"engagement_rate,omitempty"`
	MinPriceETB         *float64    `json:"min_price_etb,omitempty"`
	MaxPriceETB         *float64    `json:"max_price_etb,omitempty"`
	Languages           []string    `json:"languages"`
	City                string      `json:"city"`
	Region              string      `json:"region,omitempty"`
	AudienceDescription string      `json:"audience_description,omitempty"`
}

type SubmitProfileResponse struct {
	Status      string    `json:"status"`
	Message     string    `json:"message"`
	SubmittedAt time.Time `json:"submitted_at"`
}

type UploadAvatarResponse struct {
	ProfileImageURL string `json:"profile_image_url"`
}

type UploadPortfolioResponse struct {
	Item           PortfolioItem   `json:"item"`
	PortfolioItems []PortfolioItem `json:"portfolio_items"`
	PortfolioCount int             `json:"portfolio_count"`
}

type DeletePortfolioResponse struct {
	PortfolioItems []PortfolioItem `json:"portfolio_items"`
	PortfolioCount int             `json:"portfolio_count"`
}

type CreatePortfolioItemInput struct {
	ID            string
	UserID        string
	FileURL       string
	ObjectKey     string
	MediaType     string
	MimeType      string
	FileSizeBytes int64
}

type UpdateProfileInput struct {
	UserID              string
	CreatorName         string
	FullName            string
	Bio                 string
	City                string
	Region              *string
	Languages           []string
	SocialLinksJSON     []byte
	InstagramFollowers  int
	TikTokFollowers     int
	YouTubeSubscribers  int
	TotalFollowers      int
	EngagementRate      *float64
	MinPriceETB         *float64
	MaxPriceETB         *float64
	AudienceDescription *string
	CategoryIDs         []string
	Platforms           []string
}

type ProfileRecord struct {
	ProfileDetailResponse
	UserID string
}
