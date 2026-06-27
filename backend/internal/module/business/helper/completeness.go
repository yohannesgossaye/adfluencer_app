package helper

import (
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	"strings"
)

func BuildCompleteness(profile businessDto.ProfileDetailResponse) businessDto.CompletenessResponse {
	var score int
	missing := make([]string, 0, 10)

	add := func(points int, field string, ok bool) {
		if ok {
			score += points
			return
		}
		missing = append(missing, field)
	}

	add(10, "owner_full_name", strings.TrimSpace(profile.OwnerFullName) != "")
	add(10, "business_name", strings.TrimSpace(profile.BusinessName) != "")
	add(15, "description", profile.Description != nil && len(strings.TrimSpace(*profile.Description)) >= 20)
	add(10, "industry", profile.Industry != nil && strings.TrimSpace(*profile.Industry) != "")
	add(10, "city", profile.City != nil && strings.TrimSpace(*profile.City) != "")
	add(10, "preferred_platforms", len(profile.PreferredPlatforms) > 0)
	add(10, "logo_url", profile.LogoURL != nil && strings.TrimSpace(*profile.LogoURL) != "")

	hasWebOrSocial := (profile.WebsiteURL != nil && strings.TrimSpace(*profile.WebsiteURL) != "") ||
		strings.TrimSpace(profile.SocialLinks.Instagram) != "" ||
		strings.TrimSpace(profile.SocialLinks.TikTok) != "" ||
		strings.TrimSpace(profile.SocialLinks.Facebook) != "" ||
		strings.TrimSpace(profile.SocialLinks.LinkedIn) != ""
	add(5, "website_or_social", hasWebOrSocial)
	add(5, "business_type", profile.BusinessType != nil && strings.TrimSpace(*profile.BusinessType) != "")
	add(5, "monthly_marketing_budget", profile.MonthlyMarketingBudget != nil)

	if score > 100 {
		score = 100
	}

	canComplete := strings.TrimSpace(profile.OwnerFullName) != "" &&
		strings.TrimSpace(profile.BusinessName) != "" &&
		profile.Description != nil && len(strings.TrimSpace(*profile.Description)) >= 20 &&
		profile.Industry != nil && strings.TrimSpace(*profile.Industry) != "" &&
		profile.City != nil && strings.TrimSpace(*profile.City) != "" &&
		len(profile.PreferredPlatforms) > 0

	return businessDto.CompletenessResponse{
		Percent:     score,
		Missing:     missing,
		CanComplete: canComplete,
	}
}

func BuildMarketplaceAccess(profile businessDto.ProfileDetailResponse, userStatus string) businessDto.MarketplaceAccessResponse {
	active := userStatus == "active"
	verified := profile.IsVerified

	access := businessDto.MarketplaceAccessResponse{
		VisibleInMarketplace: active && verified,
		CanCreateContracts:   active && verified,
	}

	switch {
	case !active:
		access.Message = "Complete your business profile to activate your account."
	case !verified:
		access.Message = "Your profile is active. Admin verification is required before you appear to creators and can create contracts."
	default:
		access.Message = "Your business is verified and visible in the marketplace."
	}

	return access
}
