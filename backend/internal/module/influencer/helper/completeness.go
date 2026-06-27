package helper

import (
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
	"strings"
)

type FollowerInput struct {
	Platforms            []string
	Followers            *int
	InstagramFollowers   *int
	TikTokFollowers      *int
	YouTubeSubscribers   *int
}

func ResolveFollowers(in FollowerInput) (instagram, tiktok, youtube, total int) {
	if in.InstagramFollowers != nil || in.TikTokFollowers != nil || in.YouTubeSubscribers != nil {
		if in.InstagramFollowers != nil {
			instagram = *in.InstagramFollowers
		}
		if in.TikTokFollowers != nil {
			tiktok = *in.TikTokFollowers
		}
		if in.YouTubeSubscribers != nil {
			youtube = *in.YouTubeSubscribers
		}
		return instagram, tiktok, youtube, instagram + tiktok + youtube
	}

	if in.Followers == nil || *in.Followers <= 0 {
		return 0, 0, 0, 0
	}

	count := *in.Followers
	if len(in.Platforms) == 1 {
		switch in.Platforms[0] {
		case "instagram":
			return count, 0, 0, count
		case "tiktok":
			return 0, count, 0, count
		case "youtube":
			return 0, 0, count, count
		}
	}

	return 0, 0, 0, count
}

func SocialLinksToMap(links influencerDto.SocialLinks) map[string]string {
	out := make(map[string]string)
	if v := strings.TrimSpace(links.TikTok); v != "" {
		out["tiktok"] = v
	}
	if v := strings.TrimSpace(links.Instagram); v != "" {
		out["instagram"] = v
	}
	if v := strings.TrimSpace(links.YouTube); v != "" {
		out["youtube"] = v
	}
	return out
}

func MapToSocialLinks(raw map[string]string) influencerDto.SocialLinks {
	return influencerDto.SocialLinks{
		TikTok:    raw["tiktok"],
		Instagram: raw["instagram"],
		YouTube:   raw["youtube"],
	}
}

func HasSocialLinkForPlatform(links influencerDto.SocialLinks, platform string) bool {
	switch platform {
	case "tiktok":
		return strings.TrimSpace(links.TikTok) != ""
	case "instagram":
		return strings.TrimSpace(links.Instagram) != ""
	case "youtube":
		return strings.TrimSpace(links.YouTube) != ""
	default:
		return false
	}
}

type CompletenessInput struct {
	Profile   influencerDto.ProfileDetailResponse
	Platforms []string
}

func portfolioCount(p influencerDto.ProfileDetailResponse) int {
	return len(p.PortfolioItems)
}

func BuildCompleteness(in CompletenessInput) influencerDto.CompletenessResponse {
	var score int
	missing := make([]string, 0, 8)

	add := func(points int, field string, ok bool) {
		if ok {
			score += points
			return
		}
		missing = append(missing, field)
	}

	p := in.Profile
	add(10, "creator_name", strings.TrimSpace(p.CreatorName) != "")
	add(15, "bio", p.Bio != nil && len(strings.TrimSpace(*p.Bio)) >= 20)
	add(10, "niche", len(p.Niches) > 0)
	add(10, "platform", len(p.Platforms) > 0)

	hasLink := false
	for _, platform := range in.Platforms {
		if HasSocialLinkForPlatform(p.SocialLinks, platform) {
			hasLink = true
			break
		}
	}
	add(10, "social_link", hasLink)
	add(10, "followers", p.TotalFollowers > 0)
	add(5, "engagement_rate", p.EngagementRate != nil)
	add(10, "price_range", p.MinPriceETB != nil && p.MaxPriceETB != nil && *p.MaxPriceETB >= *p.MinPriceETB)
	add(5, "city", p.City != nil && strings.TrimSpace(*p.City) != "")
	add(5, "languages", len(p.Languages) > 0)
	add(10, "profile_image", p.ProfileImageURL != nil && strings.TrimSpace(*p.ProfileImageURL) != "")

	portfolioItems := portfolioCount(p)
	if portfolioItems > 0 {
		pts := portfolioItems * 10
		if pts > 30 {
			pts = 30
		}
		score += pts
	} else {
		missing = append(missing, "portfolio_sample")
	}

	if score > 100 {
		score = 100
	}

	canSubmit := len(missing) == 0 || (score >= 70 && !contains(missing, "bio") && !contains(missing, "niche") && !contains(missing, "platform") && !contains(missing, "city") && !contains(missing, "price_range") && !contains(missing, "followers"))

	// Phase A submit rules (no avatar/portfolio required yet)
	canSubmit = strings.TrimSpace(p.CreatorName) != "" &&
		p.Bio != nil && len(strings.TrimSpace(*p.Bio)) >= 20 &&
		len(p.Niches) > 0 &&
		len(p.Platforms) > 0 &&
		hasLink &&
		p.TotalFollowers > 0 &&
		p.EngagementRate != nil &&
		p.MinPriceETB != nil && p.MaxPriceETB != nil &&
		p.City != nil && strings.TrimSpace(*p.City) != ""

	return influencerDto.CompletenessResponse{
		Percent:   score,
		Missing:   missing,
		CanSubmit: canSubmit,
	}
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

func SubmitMissing(in CompletenessInput) []string {
	c := BuildCompleteness(in)
	if c.CanSubmit {
		return nil
	}
	return c.Missing
}
