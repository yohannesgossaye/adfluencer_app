package helper

import (
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	"sort"
	"strings"
)

func ClassifyCreatorSize(followers int) string {
	switch {
	case followers < 10_000:
		return "nano"
	case followers < 70_000:
		return "micro"
	case followers < 110_000:
		return "mid"
	case followers < 500_000:
		return "macro"
	default:
		return "mega"
	}
}

func CreatorSizeMatches(preference string, followers int) bool {
	preference = strings.ToLower(strings.TrimSpace(preference))
	if preference == "" || preference == "any" {
		return true
	}
	return ClassifyCreatorSize(followers) == preference
}

func budgetsOverlap(campaignMin, campaignMax, influencerMin, influencerMax *float64) bool {
	if influencerMin == nil && influencerMax == nil {
		return false
	}

	if campaignMin == nil && campaignMax == nil {
		return true
	}

	cMin := 0.0
	cMax := 1e18
	if campaignMin != nil {
		cMin = *campaignMin
	}
	if campaignMax != nil {
		cMax = *campaignMax
	}

	iMin := 0.0
	iMax := 1e18
	if influencerMin != nil {
		iMin = *influencerMin
	}
	if influencerMax != nil {
		iMax = *influencerMax
	} else if influencerMin != nil {
		iMax = *influencerMin
	}

	return iMin <= cMax && iMax >= cMin
}

func platformsOverlap(campaignPlatforms, influencerPlatforms []string) bool {
	if len(campaignPlatforms) == 0 || len(influencerPlatforms) == 0 {
		return false
	}
	set := make(map[string]struct{}, len(campaignPlatforms))
	for _, p := range campaignPlatforms {
		set[strings.ToLower(strings.TrimSpace(p))] = struct{}{}
	}
	for _, p := range influencerPlatforms {
		if _, ok := set[strings.ToLower(strings.TrimSpace(p))]; ok {
			return true
		}
	}
	return false
}

func nicheMatches(campaignCategoryID *string, candidateCategoryIDs []string) bool {
	if campaignCategoryID == nil || *campaignCategoryID == "" {
		return false
	}
	for _, id := range candidateCategoryIDs {
		if id == *campaignCategoryID {
			return true
		}
	}
	return false
}

func ScoreCandidate(campaign businessDto.CampaignRecord, candidate businessDto.MatchCandidate) (float64, businessDto.ScoreBreakdown) {
	breakdown := businessDto.ScoreBreakdown{}

	var categoryID *string
	if campaign.Category != nil {
		categoryID = &campaign.Category.ID
	}
	if nicheMatches(categoryID, candidate.CategoryIDs) {
		breakdown.Niche = 40
	}

	if platformsOverlap(campaign.Platforms, candidate.Platforms) {
		breakdown.Platform = 25
	}

	if budgetsOverlap(campaign.BudgetMinETB, campaign.BudgetMaxETB, candidate.MinPriceETB, candidate.MaxPriceETB) {
		breakdown.Budget = 20
	}

	if CreatorSizeMatches(campaign.CreatorSizePreference, candidate.TotalFollowers) {
		breakdown.CreatorSize = 10
	}

	if candidate.IsVerified {
		breakdown.Verified = 5
	}

	return breakdown.Total(), breakdown
}

type scoredCandidate struct {
	candidate businessDto.MatchCandidate
	score     float64
	breakdown businessDto.ScoreBreakdown
}

func RankCandidates(campaign businessDto.CampaignRecord, candidates []businessDto.MatchCandidate) []businessDto.SaveRecommendationInput {
	scoredList := make([]scoredCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		score, breakdown := ScoreCandidate(campaign, candidate)
		if score < businessDto.MinRecommendationScore {
			continue
		}
		scoredList = append(scoredList, scoredCandidate{
			candidate: candidate,
			score:     score,
			breakdown: breakdown,
		})
	}

	sort.Slice(scoredList, func(i, j int) bool {
		if scoredList[i].score != scoredList[j].score {
			return scoredList[i].score > scoredList[j].score
		}
		return scoredList[i].candidate.TotalFollowers > scoredList[j].candidate.TotalFollowers
	})

	limit := businessDto.MaxStoredRecommendations
	if len(scoredList) < limit {
		limit = len(scoredList)
	}

	out := make([]businessDto.SaveRecommendationInput, 0, limit)
	for i := 0; i < limit; i++ {
		item := scoredList[i]
		out = append(out, businessDto.SaveRecommendationInput{
			InfluencerUserID: item.candidate.UserID,
			MatchScore:       item.score,
			ScoreBreakdown:   item.breakdown,
			RankPosition:     i + 1,
		})
	}
	return out
}

func PrimaryNiche(names []string) string {
	if len(names) == 0 {
		return ""
	}
	return names[0]
}
