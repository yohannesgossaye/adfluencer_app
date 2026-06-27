package business

import (
	"errors"
	"time"
)

const MinRecommendationScore = 78
const MaxStoredRecommendations = 50

var (
	ErrCreatorNotFound = errors.New("creator not found")
)

type ScoreBreakdown struct {
	Niche       int `json:"niche"`
	Platform    int `json:"platform"`
	Budget      int `json:"budget"`
	CreatorSize int `json:"creator_size"`
	Verified    int `json:"verified"`
}

func (b ScoreBreakdown) Total() float64 {
	return float64(b.Niche + b.Platform + b.Budget + b.CreatorSize + b.Verified)
}

type MatchCandidate struct {
	UserID          string
	CreatorName     string
	ProfileImageURL *string
	Bio             *string
	TotalFollowers  int
	EngagementRate  *float64
	MinPriceETB     *float64
	MaxPriceETB     *float64
	IsVerified      bool
	CategoryIDs     []string
	CategoryNames   []string
	Platforms       []string
}

type SaveRecommendationInput struct {
	ID               string
	CampaignID       string
	InfluencerUserID string
	MatchScore       float64
	ScoreBreakdown   ScoreBreakdown
	RankPosition     int
}

type ListRecommendationsQuery struct {
	Search    string
	Niche     string
	Platform  string
	MinScore  float64
	Sort      string
	Page      int
	Limit     int
}

type RecommendationItem struct {
	InfluencerUserID string          `json:"influencer_user_id"`
	CreatorName      string          `json:"creator_name"`
	ProfileImageURL  *string         `json:"profile_image_url,omitempty"`
	PrimaryNiche     string          `json:"primary_niche,omitempty"`
	Platforms        []string        `json:"platforms"`
	TotalFollowers   int             `json:"total_followers"`
	EngagementRate   *float64        `json:"engagement_rate,omitempty"`
	MinPriceETB      *float64        `json:"min_price_etb,omitempty"`
	MaxPriceETB      *float64        `json:"max_price_etb,omitempty"`
	MatchScore       float64         `json:"match_score"`
	ScoreBreakdown   ScoreBreakdown  `json:"score_breakdown"`
	IsVerified       bool            `json:"is_verified"`
	InviteStatus     *string         `json:"invite_status,omitempty"`
	RankPosition     int             `json:"rank_position,omitempty"`
}

type ListRecommendationsResponse struct {
	Items      []RecommendationItem `json:"items"`
	Pagination PaginationResponse   `json:"pagination"`
}

type GenerateRecommendationsResponse struct {
	Generated int `json:"generated"`
	Total     int `json:"total"`
}

func (q ListRecommendationsQuery) Normalize() ListRecommendationsQuery {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}
	if q.MinScore <= 0 {
		q.MinScore = MinRecommendationScore
	}
	switch q.Sort {
	case "followers", "price_low", "price_high", "match":
	default:
		q.Sort = "match"
	}
	return q
}

type RecommendationRecord struct {
	CampaignID       string
	InfluencerUserID string
	MatchScore       float64
	ScoreBreakdown   ScoreBreakdown
	RankPosition     int
	Status           string
	CreatedAt        time.Time
}
