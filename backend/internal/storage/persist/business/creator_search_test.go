package business

import (
	"context"
	"os"
	"testing"

	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	"adinfluence_service_app/internal/module/business/helper"
	"adinfluence_service_app/internal/testutil"
)

// The demo seed loads cleanly (twice) and "Find creators" filters it right.
func TestSearchCreatorsOnDemoSeed(t *testing.T) {
	db := testutil.MigratedDB(t)
	seed, err := os.ReadFile("../../../../scripts/seed_demo_creators.sql")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if _, err := db.Exec(string(seed)); err != nil {
			t.Fatalf("seed run %d: %v", i+1, err)
		}
	}
	// A creator who is not approved must never show up.
	if _, err := db.Exec(`UPDATE influencer_profiles SET is_verified = FALSE WHERE creator_name = 'Dawit Unboxing'`); err != nil {
		t.Fatal(err)
	}
	p := &business{db: db}
	ctx := context.Background()
	count := func(f businessDto.CreatorFilter) int {
		t.Helper()
		rows, err := p.SearchCreators(ctx, f)
		if err != nil {
			t.Fatal(err)
		}
		return len(rows)
	}

	cases := []struct {
		name string
		f    businessDto.CreatorFilter
		want int
	}{
		{"all approved", businessDto.CreatorFilter{}, 34},
		{"fashion", businessDto.CreatorFilter{Niche: "fashion"}, 10},
		{"food", businessDto.CreatorFilter{Niche: "food"}, 5},
		{"tech without the unapproved one", businessDto.CreatorFilter{Niche: "tech"}, 4},
		{"sports incl. runners", businessDto.CreatorFilter{Niche: "sports"}, 6},
		{"fashion on tiktok", businessDto.CreatorFilter{Niche: "fashion", Platform: "tiktok"}, 6},
		{"starting price within 5,000", businessDto.CreatorFilter{MaxPrice: 5000}, 15},
		{"city", businessDto.CreatorFilter{City: "hawassa"}, 2},
		{"search by niche word", businessDto.CreatorFilter{Q: "travel"}, 2},
		{"search by bio word", businessDto.CreatorFilter{Q: "injera"}, 1},
	}
	for _, c := range cases {
		if got := count(c.f); got != c.want {
			t.Errorf("%s: %d creators, want %d", c.name, got, c.want)
		}
	}

	rows, err := p.SearchCreators(ctx, businessDto.CreatorFilter{Q: "hana"})
	if err != nil {
		t.Fatal(err)
	}
	ranked := helper.RankCreators(rows, "hana", "", "best", nil)
	if len(ranked) == 0 || ranked[0].Row.CreatorName != "Hana Style Addis" {
		t.Fatalf("searching hana: first = %+v", ranked)
	}
	if r := ranked[0].Row; len(r.NicheSlugs) != 1 || r.NicheSlugs[0] != "fashion" || r.Instagram == 0 || r.City != "Addis Ababa" {
		t.Fatalf("row details: %+v", r)
	}
}
