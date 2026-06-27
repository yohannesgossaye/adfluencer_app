package helper

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
	"fmt"
	"strings"
)

func PricingSummary(min, max *float64) string {
	if min == nil || max == nil {
		return ""
	}
	return fmt.Sprintf("ETB %.0f – %.0f", *min, *max)
}

func BioPreview(bio *string) string {
	if bio == nil {
		return ""
	}
	text := strings.TrimSpace(*bio)
	if len(text) <= 120 {
		return text
	}
	return text[:117] + "..."
}

func CompletenessFromRow(row adminDto.ApplicationRow) int {
	profile := influencerDto.ProfileDetailResponse{
		CreatorName:     row.CreatorName,
		Bio:             row.Bio,
		City:            row.City,
		ProfileImageURL: row.ProfileImageURL,
		TotalFollowers:  row.TotalFollowers,
		EngagementRate:  row.EngagementRate,
		MinPriceETB:     row.MinPriceETB,
		MaxPriceETB:     row.MaxPriceETB,
	}
	if row.PrimaryNiche != "" {
		profile.Niches = []influencerDto.CategoryResponse{{Name: row.PrimaryNiche}}
	}
	if row.PrimaryPlatform != "" {
		profile.Platforms = []string{row.PrimaryPlatform}
	}

	// Reuse influencer completeness weights without importing cycle - duplicate minimal calc
	var score int
	if strings.TrimSpace(profile.CreatorName) != "" {
		score += 10
	}
	if profile.Bio != nil && len(strings.TrimSpace(*profile.Bio)) >= 20 {
		score += 15
	}
	if len(profile.Niches) > 0 {
		score += 10
	}
	if len(profile.Platforms) > 0 {
		score += 10
	}
	if profile.TotalFollowers > 0 {
		score += 10
	}
	if profile.EngagementRate != nil {
		score += 5
	}
	if profile.MinPriceETB != nil && profile.MaxPriceETB != nil {
		score += 10
	}
	if profile.City != nil && strings.TrimSpace(*profile.City) != "" {
		score += 5
	}
	if profile.ProfileImageURL != nil && strings.TrimSpace(*profile.ProfileImageURL) != "" {
		score += 10
	}
	if score > 100 {
		score = 100
	}
	return score
}
