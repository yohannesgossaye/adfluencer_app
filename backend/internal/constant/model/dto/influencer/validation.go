package influencer

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"net/url"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var allowedPlatforms = map[string]struct{}{
	"tiktok":    {},
	"instagram": {},
	"youtube":   {},
}

func (r UpdateProfileRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.CreatorName,
			validation.Required.Error(customErr.ErrNameIsRequired.Error()),
			validation.Length(2, 255).Error(customErr.ErrInvalidNameLength.Error()),
		),
		validation.Field(&r.Bio,
			validation.Required.Error("bio is required"),
			validation.Length(20, 2000).Error("bio must be between 20 and 2000 characters"),
		),
		validation.Field(&r.Niches,
			validation.Required.Error("at least one niche is required"),
			validation.Length(1, 5).Error("select between 1 and 5 niches"),
		),
		validation.Field(&r.Platforms,
			validation.Required.Error("at least one platform is required"),
			validation.Length(1, 5).Error("select between 1 and 5 platforms"),
			validation.By(validatePlatforms),
		),
		validation.Field(&r.City,
			validation.Required.Error("city is required"),
			validation.Length(2, 100).Error("city must be between 2 and 100 characters"),
		),
		validation.Field(&r.MinPriceETB,
			validation.By(func(value interface{}) error {
				if r.MinPriceETB == nil || r.MaxPriceETB == nil {
					return nil
				}
				if *r.MaxPriceETB < *r.MinPriceETB {
					return validation.NewError("invalid_price_range", "max price must be greater than or equal to min price")
				}
				return nil
			}),
		),
		validation.Field(&r.SocialLinks,
			validation.By(validateSocialLinks),
		),
	)
}

func validatePlatforms(value interface{}) error {
	platforms, _ := value.([]string)
	for _, p := range platforms {
		if _, ok := allowedPlatforms[strings.ToLower(strings.TrimSpace(p))]; !ok {
			return validation.NewError("invalid_platform", "platform must be tiktok, instagram, or youtube")
		}
	}
	return nil
}

func validateSocialLinks(value interface{}) error {
	links, ok := value.(SocialLinks)
	if !ok {
		return nil
	}
	for field, raw := range map[string]string{
		"tiktok":    links.TikTok,
		"instagram": links.Instagram,
		"youtube":   links.YouTube,
	} {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		u, err := url.ParseRequestURI(raw)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return validation.NewError("invalid_social_link", field+" link must be a valid URL")
		}
	}
	return nil
}

func NormalizePlatforms(platforms []string) []string {
	out := make([]string, 0, len(platforms))
	seen := make(map[string]struct{}, len(platforms))
	for _, p := range platforms {
		key := strings.ToLower(strings.TrimSpace(p))
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, key)
	}
	return out
}
