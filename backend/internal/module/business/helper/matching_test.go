package helper

import (
	"testing"

	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
)

// The SQL prefilter in ListMatchCandidates keeps only creators who share the
// campaign's niche and at least one platform. That is safe only while a
// creator missing either one can never reach MinRecommendationScore.
func TestPrefilterDropsNoRealMatch(t *testing.T) {
	budget := 10000.0
	var campaign businessDto.CampaignRecord
	campaign.Category = &businessDto.CampaignCategoryResponse{ID: "niche-1"}
	campaign.Platforms = []string{"TikTok"}
	campaign.BudgetMinETB = &budget
	campaign.BudgetMaxETB = &budget
	campaign.CreatorSizePreference = "any"
	best := businessDto.MatchCandidate{
		CategoryIDs:    []string{"niche-1"},
		Platforms:      []string{"tiktok"},
		MinPriceETB:    &budget,
		MaxPriceETB:    &budget,
		TotalFollowers: 5000,
		IsVerified:     true,
	}

	if score, _ := ScoreCandidate(campaign, best); score < businessDto.MinRecommendationScore {
		t.Fatalf("a full match scored %v, below the minimum", score)
	}

	noNiche := best
	noNiche.CategoryIDs = []string{"other"}
	if score, _ := ScoreCandidate(campaign, noNiche); score >= businessDto.MinRecommendationScore {
		t.Fatalf("without the niche a creator scored %v; the SQL prefilter would wrongly drop them", score)
	}

	noPlatform := best
	noPlatform.Platforms = []string{"youtube"}
	if score, _ := ScoreCandidate(campaign, noPlatform); score >= businessDto.MinRecommendationScore {
		t.Fatalf("without a platform a creator scored %v; the SQL prefilter would wrongly drop them", score)
	}
}

func TestNormalizePlatforms(t *testing.T) {
	got := NormalizePlatforms([]string{" TikTok ", "", "INSTAGRAM"})
	if len(got) != 2 || got[0] != "tiktok" || got[1] != "instagram" {
		t.Fatalf("got %v", got)
	}
}
