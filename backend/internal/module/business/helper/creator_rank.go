package helper

import (
	"math"
	"sort"
	"strings"

	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
)

// RankedCreator is a creator with what decided their place in the list.
type RankedCreator struct {
	Row       businessDto.CreatorRow
	Score     *float64 // campaign match score, only when ranked for a campaign
	Breakdown *businessDto.ScoreBreakdown
	textRank  int
	quality   float64
}

// textRank says how well a creator matches the search words: an exact name
// beats a name that starts with them, which beats a name that contains them,
// which beats a niche, city or bio match.
func textRank(row businessDto.CreatorRow, q string) int {
	if q == "" {
		return 0
	}
	q = strings.ToLower(q)
	name := strings.ToLower(row.CreatorName)
	switch {
	case name == q:
		return 4
	case strings.HasPrefix(name, q):
		return 3
	case strings.Contains(name, q):
		return 2
	}
	for _, n := range row.CategoryNames {
		if strings.Contains(strings.ToLower(n), q) {
			return 1
		}
	}
	return 0
}

// quality balances reach and engagement: the square of the following's order
// of magnitude times the square root of the engagement rate, so a 7K account
// with 10% engagement doesn't outrank a 400K account with 4%, and a big but
// inactive account doesn't win on size alone. No engagement rate counts as 2%.
func quality(row businessDto.CreatorRow) float64 {
	e := 2.0
	if row.EngagementRate != nil && *row.EngagementRate > 0 {
		e = math.Min(*row.EngagementRate, 25)
	}
	reach := math.Log10(float64(row.TotalFollowers) + 10)
	q := math.Sqrt(e) * reach * reach
	if row.IsVerified {
		q += 3
	}
	return q
}

func minPrice(row businessDto.CreatorRow) float64 {
	if row.MinPriceETB != nil {
		return *row.MinPriceETB
	}
	return math.MaxFloat64
}

// RankCreators filters by size and orders creators for "Find creators".
// sortBy is best, followers, engagement, price_low or price_high. "best"
// puts the searched name first, then (with a campaign) the strongest
// campaign fit, then engaged audiences.
func RankCreators(rows []businessDto.CreatorRow, q, size, sortBy string, campaign *businessDto.CampaignRecord) []RankedCreator {
	out := make([]RankedCreator, 0, len(rows))
	for _, row := range rows {
		if size != "" && ClassifyCreatorSize(row.TotalFollowers) != size {
			continue
		}
		rc := RankedCreator{Row: row, textRank: textRank(row, q), quality: quality(row)}
		if campaign != nil {
			score, breakdown := ScoreCandidate(*campaign, row.MatchCandidate)
			rc.Score, rc.Breakdown = &score, &breakdown
		}
		out = append(out, rc)
	}

	engagement := func(r RankedCreator) float64 {
		if r.Row.EngagementRate == nil {
			return 0
		}
		return *r.Row.EngagementRate
	}
	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i], out[j]
		switch sortBy {
		case "followers":
			if a.Row.TotalFollowers != b.Row.TotalFollowers {
				return a.Row.TotalFollowers > b.Row.TotalFollowers
			}
		case "engagement":
			if engagement(a) != engagement(b) {
				return engagement(a) > engagement(b)
			}
		case "price_low":
			if minPrice(a.Row) != minPrice(b.Row) {
				return minPrice(a.Row) < minPrice(b.Row)
			}
		case "price_high":
			if minPrice(a.Row) != minPrice(b.Row) {
				// Creators without a price go last here too.
				if minPrice(a.Row) == math.MaxFloat64 || minPrice(b.Row) == math.MaxFloat64 {
					return minPrice(a.Row) < minPrice(b.Row)
				}
				return minPrice(a.Row) > minPrice(b.Row)
			}
		default: // best: the searched name first, then campaign fit, then quality
			if a.textRank != b.textRank {
				return a.textRank > b.textRank
			}
			if a.Score != nil && b.Score != nil && *a.Score != *b.Score {
				return *a.Score > *b.Score
			}
		}
		if a.quality != b.quality {
			return a.quality > b.quality
		}
		return a.Row.TotalFollowers > b.Row.TotalFollowers
	})
	return out
}
