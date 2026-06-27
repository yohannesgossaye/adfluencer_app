package business

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"fmt"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var allowedCampaignPlatforms = map[string]struct{}{
	"tiktok":    {},
	"instagram": {},
	"youtube":   {},
	"facebook":  {},
	"telegram":  {},
}

var allowedDeliverableTypes = map[string]struct{}{
	"tiktok_video":     {},
	"instagram_reel":   {},
	"instagram_story":  {},
	"youtube_short":    {},
	"product_review":   {},
	"post":             {},
	"other":            {},
}

var allowedCreatorSizes = map[string]struct{}{
	"nano":  {},
	"micro": {},
	"mid":   {},
	"macro": {},
	"mega":  {},
	"any":   {},
}

func (r CreateCampaignRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Title,
			validation.Required.Error("title is required"),
			validation.Length(3, 255).Error("title must be between 3 and 255 characters"),
		),
		validation.Field(&r.Description,
			validation.Required.Error("description is required"),
			validation.Length(20, 10000).Error("description must be between 20 and 10000 characters"),
		),
		validation.Field(&r.Niche,
			validation.Required.Error("niche is required"),
		),
		validation.Field(&r.City,
			validation.Required.Error("city is required"),
			validation.Length(2, 100).Error("city must be between 2 and 100 characters"),
		),
		validation.Field(&r.Platforms,
			validation.Required.Error("at least one platform is required"),
			validation.Length(1, 5).Error("select between 1 and 5 platforms"),
			validation.By(validateCampaignPlatforms),
		),
		validation.Field(&r.CreatorSizePreference,
			validation.Required.Error("creator size preference is required"),
			validation.By(validateCreatorSize),
		),
		validation.Field(&r.StartDate, validation.By(validateDateField("start_date"))),
		validation.Field(&r.EndDate, validation.By(validateDateField("end_date"))),
		validation.Field(&r, validation.By(validateCampaignDateRange)),
		validation.Field(&r.BudgetMinETB, validation.By(validateBudgetRange(&r))),
		validation.Field(&r.Deliverables, validation.By(validateDeliverables)),
		validation.Field(&r.NumberOfCreatorsNeeded,
			validation.Min(1).Error("number of creators must be at least 1"),
			validation.Max(50).Error("number of creators cannot exceed 50"),
		),
	)
}

func (r CreateCampaignRequest) ValidateForPublish() error {
	if err := r.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(r.Goal) == "" {
		return validation.NewError("goal_required", "goal is required")
	}
	if r.BudgetMinETB == nil || r.BudgetMaxETB == nil {
		return validation.NewError("budget_required", "budget range is required")
	}
	return nil
}

func validateCampaignPlatforms(value interface{}) error {
	platforms, _ := value.([]string)
	for _, p := range platforms {
		if _, ok := allowedCampaignPlatforms[strings.ToLower(strings.TrimSpace(p))]; !ok {
			return validation.NewError("invalid_platform", "invalid campaign platform")
		}
	}
	return nil
}

func validateCreatorSize(value interface{}) error {
	raw, _ := value.(string)
	raw = strings.ToLower(strings.TrimSpace(raw))
	if raw == "" {
		raw = "any"
	}
	if _, ok := allowedCreatorSizes[raw]; !ok {
		return validation.NewError("invalid_creator_size", "invalid creator size preference")
	}
	return nil
}

func validateDeliverables(value interface{}) error {
	items, _ := value.([]CampaignDeliverableInput)
	for _, item := range items {
		t := strings.ToLower(strings.TrimSpace(item.Type))
		if t == "" {
			return validation.NewError("invalid_deliverable", "deliverable type is required")
		}
		if _, ok := allowedDeliverableTypes[t]; !ok {
			return validation.NewError("invalid_deliverable", "invalid deliverable type")
		}
		if item.Quantity < 1 {
			return validation.NewError("invalid_deliverable", "deliverable quantity must be at least 1")
		}
	}
	return nil
}

func validateDateField(field string) validation.RuleFunc {
	return func(value interface{}) error {
		raw, _ := value.(string)
		raw = strings.TrimSpace(raw)
		if raw == "" {
			return validation.NewError(field+"_required", field+" is required")
		}
		if _, err := time.Parse("2006-01-02", raw); err != nil {
			return validation.NewError("invalid_date", field+" must be YYYY-MM-DD")
		}
		return nil
	}
}

func validateCampaignDateRange(value interface{}) error {
	req, ok := value.(CreateCampaignRequest)
	if !ok {
		return nil
	}
	start, err1 := time.Parse("2006-01-02", strings.TrimSpace(req.StartDate))
	end, err2 := time.Parse("2006-01-02", strings.TrimSpace(req.EndDate))
	if err1 != nil || err2 != nil {
		return nil
	}
	if end.Before(start) {
		return validation.NewError("invalid_date_range", "end date must be on or after start date")
	}
	return nil
}

func validateBudgetRange(req *CreateCampaignRequest) validation.RuleFunc {
	return func(value interface{}) error {
		if req.BudgetMinETB == nil || req.BudgetMaxETB == nil {
			return nil
		}
		if *req.BudgetMaxETB < *req.BudgetMinETB {
			return validation.NewError("invalid_budget_range", "max budget must be greater than or equal to min budget")
		}
		return nil
	}
}

func NormalizeCampaignPlatforms(platforms []string) []string {
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

func NormalizeCreatorSize(raw string) string {
	raw = strings.ToLower(strings.TrimSpace(raw))
	if raw == "" {
		return "any"
	}
	if _, ok := allowedCreatorSizes[raw]; ok {
		return raw
	}
	return "any"
}

func ParseCampaignDate(raw string) (*time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("%w: date is required", customErr.ErrBadRequest)
	}
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid date format", customErr.ErrBadRequest)
	}
	utc := t.UTC()
	return &utc, nil
}

func MergeDescription(description, notes string) string {
	description = strings.TrimSpace(description)
	notes = strings.TrimSpace(notes)
	if notes == "" {
		return description
	}
	if description == "" {
		return notes
	}
	return description + "\n\nAdditional notes:\n" + notes
}
