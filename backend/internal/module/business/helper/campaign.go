package helper

import (
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	"fmt"
	"strings"
	"time"
)

func ValidateCampaignRecordForPublish(rec businessDto.CampaignRecord) error {
	req := businessDto.CreateCampaignRequest{
		Title:                  rec.Title,
		Goal:                   stringValue(rec.Goal),
		Description:            rec.Description,
		Niche:                  categoryName(rec.Category),
		TargetAudience:         stringValue(rec.TargetAudience),
		City:                   stringValue(rec.City),
		Region:                 stringValue(rec.Region),
		Platforms:              rec.Platforms,
		BudgetMinETB:           rec.BudgetMinETB,
		BudgetMaxETB:           rec.BudgetMaxETB,
		CreatorSizePreference:  rec.CreatorSizePreference,
		NumberOfCreatorsNeeded: rec.NumberOfCreatorsNeeded,
		StartDate:              formatDate(rec.StartDate),
		EndDate:                formatDate(rec.EndDate),
	}
	return req.ValidateForPublish()
}

func stringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func categoryName(cat *businessDto.CampaignCategoryResponse) string {
	if cat == nil {
		return ""
	}
	return cat.Name
}

func formatDate(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.UTC().Format("2006-01-02")
}

func BuildSaveCampaignInput(userID string, campaignID string, req businessDto.CreateCampaignRequest, categoryID *string) (businessDto.SaveCampaignInput, error) {
	if req.NumberOfCreatorsNeeded <= 0 {
		req.NumberOfCreatorsNeeded = 1
	}

	startDate, err := businessDto.ParseCampaignDate(req.StartDate)
	if err != nil {
		return businessDto.SaveCampaignInput{}, err
	}
	endDate, err := businessDto.ParseCampaignDate(req.EndDate)
	if err != nil {
		return businessDto.SaveCampaignInput{}, err
	}

	description := businessDto.MergeDescription(req.Description, req.AdditionalNotes)

	var goal *string
	if v := strings.TrimSpace(req.Goal); v != "" {
		goal = &v
	}
	var targetAudience *string
	if v := strings.TrimSpace(req.TargetAudience); v != "" {
		targetAudience = &v
	}
	var city *string
	if v := strings.TrimSpace(req.City); v != "" {
		city = &v
	}
	var region *string
	if v := strings.TrimSpace(req.Region); v != "" {
		region = &v
	}

	return businessDto.SaveCampaignInput{
		ID:                     campaignID,
		BusinessUserID:         userID,
		Title:                  strings.TrimSpace(req.Title),
		Goal:                   goal,
		Description:            description,
		CategoryID:             categoryID,
		TargetAudience:         targetAudience,
		City:                   city,
		Region:                 region,
		Platforms:              businessDto.NormalizeCampaignPlatforms(req.Platforms),
		BudgetMinETB:           req.BudgetMinETB,
		BudgetMaxETB:           req.BudgetMaxETB,
		CreatorSizePreference:  businessDto.NormalizeCreatorSize(req.CreatorSizePreference),
		NumberOfCreatorsNeeded: req.NumberOfCreatorsNeeded,
		Deliverables:           normalizeDeliverables(req.Deliverables),
		StartDate:              startDate,
		EndDate:                endDate,
	}, nil
}

func normalizeDeliverables(items []businessDto.CampaignDeliverableInput) []businessDto.CampaignDeliverableInput {
	out := make([]businessDto.CampaignDeliverableInput, 0, len(items))
	for _, item := range items {
		qty := item.Quantity
		if qty <= 0 {
			qty = 1
		}
		out = append(out, businessDto.CampaignDeliverableInput{
			Type:        strings.ToLower(strings.TrimSpace(item.Type)),
			Quantity:    qty,
			Description: strings.TrimSpace(item.Description),
		})
	}
	return out
}

func CampaignMutation(rec businessDto.CampaignRecord) businessDto.CampaignMutationResponse {
	if rec.Platforms == nil {
		rec.Platforms = []string{}
	}
	if rec.Deliverables == nil {
		rec.Deliverables = []businessDto.CampaignDeliverableResponse{}
	}
	return businessDto.CampaignMutationResponse{Campaign: rec.CampaignDetailResponse}
}

func ListCampaignsResponse(items []businessDto.CampaignListItem, query businessDto.ListCampaignsQuery, total int) businessDto.ListCampaignsResponse {
	if items == nil {
		items = []businessDto.CampaignListItem{}
	}
	return businessDto.ListCampaignsResponse{
		Items: items,
		Pagination: commonDto.NewListMeta(query.Page, query.Limit, total),
	}
}

func FormatBudgetSummary(min, max *float64) string {
	if min == nil && max == nil {
		return ""
	}
	if min != nil && max != nil {
		return fmt.Sprintf("ETB %.0f – %.0f", *min, *max)
	}
	if min != nil {
		return fmt.Sprintf("ETB %.0f+", *min)
	}
	return fmt.Sprintf("ETB up to %.0f", *max)
}
