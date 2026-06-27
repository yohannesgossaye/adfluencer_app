package business

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
	"facebook":  {},
	"telegram":  {},
}

func (r UpdateProfileRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.OwnerFullName,
			validation.Required.Error(customErr.ErrNameIsRequired.Error()),
			validation.Length(2, 255).Error(customErr.ErrInvalidNameLength.Error()),
		),
		validation.Field(&r.BusinessName,
			validation.Required.Error("business name is required"),
			validation.Length(2, 255).Error("business name must be between 2 and 255 characters"),
		),
		validation.Field(&r.Description,
			validation.Required.Error("description is required"),
			validation.Length(20, 5000).Error("description must be between 20 and 5000 characters"),
		),
		validation.Field(&r.Industry,
			validation.Required.Error("industry is required"),
			validation.Length(2, 100).Error("industry must be between 2 and 100 characters"),
		),
		validation.Field(&r.City,
			validation.Required.Error("city is required"),
			validation.Length(2, 100).Error("city must be between 2 and 100 characters"),
		),
		validation.Field(&r.PreferredPlatforms,
			validation.Required.Error("at least one preferred platform is required"),
			validation.Length(1, 5).Error("select between 1 and 5 platforms"),
			validation.By(validatePlatforms),
		),
		validation.Field(&r.WebsiteURL, validation.By(validateOptionalURL)),
		validation.Field(&r.SocialLinks, validation.By(validateSocialLinks)),
	)
}

func validatePlatforms(value interface{}) error {
	platforms, _ := value.([]string)
	for _, p := range platforms {
		if _, ok := allowedPlatforms[strings.ToLower(strings.TrimSpace(p))]; !ok {
			return validation.NewError("invalid_platform", "platform must be tiktok, instagram, youtube, facebook, or telegram")
		}
	}
	return nil
}

func validateOptionalURL(value interface{}) error {
	raw, _ := value.(string)
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	u, err := url.ParseRequestURI(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return validation.NewError("invalid_url", "website must be a valid URL")
	}
	return nil
}

func validateSocialLinks(value interface{}) error {
	links, ok := value.(SocialLinks)
	if !ok {
		return nil
	}
	for field, raw := range map[string]string{
		"instagram": links.Instagram,
		"tiktok":    links.TikTok,
		"facebook":  links.Facebook,
		"linkedin":  links.LinkedIn,
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
