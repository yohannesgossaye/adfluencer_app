package helper

import (
	"testing"

	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
)

func creator(name string, followers int, engagement float64, niche string, platforms ...string) businessDto.CreatorRow {
	e := engagement
	price := float64(followers / 10)
	return businessDto.CreatorRow{MatchCandidate: businessDto.MatchCandidate{
		UserID: name, CreatorName: name, TotalFollowers: followers, EngagementRate: &e,
		MinPriceETB: &price, MaxPriceETB: &price, IsVerified: true,
		CategoryIDs: []string{niche}, CategoryNames: []string{niche}, Platforms: platforms,
	}}
}

func names(r []RankedCreator) []string {
	out := make([]string, len(r))
	for i, c := range r {
		out[i] = c.Row.CreatorName
	}
	return out
}

func TestRankCreatorsSearchedNameComesFirst(t *testing.T) {
	rows := []businessDto.CreatorRow{
		creator("Abel Sport Talk", 300_000, 6, "Sports", "tiktok"),
		creator("Sport Hana", 20_000, 3, "Sports", "tiktok"),
		creator("Hana", 5_000, 2, "Fashion", "instagram"),
		creator("Hana Style", 50_000, 8, "Fashion", "instagram"),
	}
	got := names(RankCreators(rows, "hana", "", "best", nil))
	want := []string{"Hana", "Hana Style", "Sport Hana", "Abel Sport Talk"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("order = %v, want %v", got, want)
		}
	}
}

func TestRankCreatorsCampaignFitFirst(t *testing.T) {
	fashion := "fashion"
	min, max := 2_000.0, 10_000.0
	var campaign businessDto.CampaignRecord
	campaign.Category = &businessDto.CampaignCategoryResponse{ID: fashion}
	campaign.Platforms = []string{"instagram"}
	campaign.BudgetMinETB, campaign.BudgetMaxETB = &min, &max
	campaign.CreatorSizePreference = "micro"
	rows := []businessDto.CreatorRow{
		creator("Big tech", 900_000, 9, "tech", "youtube"),        // no niche, no platform
		creator("Fashion TikTok", 40_000, 5, "fashion", "tiktok"), // niche only
		creator("Perfect fit", 40_000, 4, "fashion", "instagram"), // niche, platform, budget, size
	}
	ranked := RankCreators(rows, "", "", "best", &campaign)
	if got := names(ranked); got[0] != "Perfect fit" || got[2] != "Big tech" {
		t.Fatalf("order = %v", got)
	}
	if ranked[0].Score == nil || *ranked[0].Score != 100 {
		t.Fatalf("perfect fit score = %v, want 100", ranked[0].Score)
	}
}

func TestRankCreatorsSizeAndSorts(t *testing.T) {
	rows := []businessDto.CreatorRow{
		creator("Nano", 5_000, 9, "food", "tiktok"),
		creator("Micro", 30_000, 3, "food", "tiktok"),
		creator("Mega", 2_000_000, 1, "food", "tiktok"),
	}
	if got := names(RankCreators(rows, "", "micro", "best", nil)); len(got) != 1 || got[0] != "Micro" {
		t.Fatalf("size filter = %v", got)
	}
	if got := names(RankCreators(rows, "", "", "followers", nil)); got[0] != "Mega" {
		t.Fatalf("followers sort = %v", got)
	}
	if got := names(RankCreators(rows, "", "", "engagement", nil)); got[0] != "Nano" {
		t.Fatalf("engagement sort = %v", got)
	}
	if got := names(RankCreators(rows, "", "", "price_low", nil)); got[0] != "Nano" || got[2] != "Mega" {
		t.Fatalf("price_low sort = %v", got)
	}
	if got := names(RankCreators(rows, "", "", "price_high", nil)); got[0] != "Mega" {
		t.Fatalf("price_high sort = %v", got)
	}
}

// Browsing without a search balances reach and engagement.
func TestRankCreatorsDefaultBalancesReachAndEngagement(t *testing.T) {
	rows := []businessDto.CreatorRow{
		creator("Tiny but very engaged", 7_000, 10.4, "fashion", "tiktok"),
		creator("Big football channel", 620_000, 4.2, "sports", "youtube"),
		creator("Huge but inactive", 900_000, 0.3, "lifestyle", "facebook"),
	}
	got := names(RankCreators(rows, "", "", "best", nil))
	if got[0] != "Big football channel" || got[2] == "Big football channel" {
		t.Fatalf("order = %v", got)
	}
}
