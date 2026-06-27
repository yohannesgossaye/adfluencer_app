package business

import influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"

type BusinessInfluencerViewResponse struct {
	InfluencerUserID string                          `json:"influencer_user_id"`
	CreatorName      string                          `json:"creator_name"`
	FullName         string                          `json:"full_name"`
	Bio              *string                         `json:"bio,omitempty"`
	ProfileImageURL  *string                         `json:"profile_image_url,omitempty"`
	City             *string                         `json:"city,omitempty"`
	Region           *string                         `json:"region,omitempty"`
	Languages        []string                        `json:"languages"`
	Niches           []influencerDto.CategoryResponse `json:"niches"`
	Platforms        []string                        `json:"platforms"`
	SocialLinks      influencerDto.SocialLinks       `json:"social_links"`
	TotalFollowers   int                             `json:"total_followers"`
	EngagementRate   *float64                        `json:"engagement_rate,omitempty"`
	MinPriceETB      *float64                        `json:"min_price_etb,omitempty"`
	MaxPriceETB      *float64                        `json:"max_price_etb,omitempty"`
	IsVerified       bool                            `json:"is_verified"`
	PortfolioItems   []influencerDto.PortfolioItem   `json:"portfolio_items"`
	MatchScore       *float64                        `json:"match_score,omitempty"`
	ScoreBreakdown   *ScoreBreakdown                 `json:"score_breakdown,omitempty"`
}
