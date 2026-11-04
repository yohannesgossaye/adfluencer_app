package business

import (
	"strings"

	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
)

// SearchCreatorsQuery is "Find creators": approved creators a business can
// browse with or without a campaign. With CampaignID the results are ranked
// by how well each creator fits that campaign.
type SearchCreatorsQuery struct {
	Q          string // name, handle, city, bio or niche
	Niche      string // category slug
	Platform   string // tiktok, instagram, youtube, telegram, facebook
	Size       string // nano, micro, mid, macro, mega
	MaxPrice   float64
	City       string
	CampaignID string
	Sort       string // best, followers, engagement, price_low, price_high
	Page       int
	Limit      int
}

var creatorSorts = map[string]bool{"best": true, "followers": true, "engagement": true, "price_low": true, "price_high": true}
var creatorSizes = map[string]bool{"nano": true, "micro": true, "mid": true, "macro": true, "mega": true}

func (q SearchCreatorsQuery) Normalize() SearchCreatorsQuery {
	q.Q = strings.TrimSpace(q.Q)
	if len([]rune(q.Q)) > 100 {
		q.Q = string([]rune(q.Q)[:100])
	}
	q.Niche = strings.ToLower(strings.TrimSpace(q.Niche))
	q.Platform = strings.ToLower(strings.TrimSpace(q.Platform))
	q.City = strings.TrimSpace(q.City)
	q.CampaignID = strings.TrimSpace(q.CampaignID)
	if q.Size = strings.ToLower(strings.TrimSpace(q.Size)); !creatorSizes[q.Size] {
		q.Size = ""
	}
	if q.Sort = strings.ToLower(strings.TrimSpace(q.Sort)); !creatorSorts[q.Sort] {
		q.Sort = "best"
	}
	if q.MaxPrice < 0 {
		q.MaxPrice = 0
	}
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 48 {
		q.Limit = 24
	}
	return q
}

// CreatorRow is one approved creator as stored, before ranking.
type CreatorRow struct {
	MatchCandidate
	City       string
	NicheSlugs []string
	TikTok     int
	Instagram  int
	YouTube    int
}

type CreatorSearchItem struct {
	InfluencerUserID   string   `json:"influencer_user_id"`
	CreatorName        string   `json:"creator_name"`
	ProfileImageURL    *string  `json:"profile_image_url,omitempty"`
	City               string   `json:"city,omitempty"`
	Bio                string   `json:"bio,omitempty"`
	Niches             []string `json:"niches"`
	Platforms          []string `json:"platforms"`
	TotalFollowers     int      `json:"total_followers"`
	TikTokFollowers    int      `json:"tiktok_followers"`
	InstagramFollowers int      `json:"instagram_followers"`
	YouTubeSubs        int      `json:"youtube_subscribers"`
	Size               string   `json:"size"`
	EngagementRate     *float64 `json:"engagement_rate,omitempty"`
	MinPriceETB        *float64 `json:"min_price_etb,omitempty"`
	MaxPriceETB        *float64 `json:"max_price_etb,omitempty"`
	IsVerified         bool     `json:"is_verified"`
	// Only when ranked for a campaign.
	MatchScore     *float64        `json:"match_score,omitempty"`
	ScoreBreakdown *ScoreBreakdown `json:"score_breakdown,omitempty"`
	InviteStatus   *string         `json:"invite_status,omitempty"`
}

type RankedForCampaign struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type SearchCreatorsResponse struct {
	Items      []CreatorSearchItem `json:"items"`
	Pagination commonDto.ListMeta  `json:"pagination"`
	RankedFor  *RankedForCampaign  `json:"ranked_for,omitempty"`
}

// CreatorFilter is what the database filters on; ranking happens in Go.
type CreatorFilter struct {
	Q        string
	Niche    string
	Platform string
	MaxPrice float64
	City     string
}
