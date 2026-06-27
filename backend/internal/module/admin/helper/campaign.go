package helper

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	"fmt"
)

func FormatContractSummary(active, pending int) string {
	if active == 0 && pending == 0 {
		return "None"
	}
	if active > 0 && pending > 0 {
		return fmt.Sprintf("%d active, %d pending", active, pending)
	}
	if active > 0 {
		return fmt.Sprintf("%d active", active)
	}
	return fmt.Sprintf("%d pending", pending)
}

func FormatPaymentSummary(paid, pending int) string {
	switch {
	case paid > 0 && pending > 0:
		return "Partially funded"
	case paid > 0:
		return "Funded"
	case pending > 0:
		return "Pending deposit"
	default:
		return "Not started"
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

func DisputeTitle(issueTypeLabel, campaignTitle string) string {
	if campaignTitle == "" {
		return issueTypeLabel
	}
	return issueTypeLabel + " — " + campaignTitle
}

func DisputePreview(description string) string {
	if len(description) <= 160 {
		return description
	}
	return description[:160] + "…"
}

func MapCampaignListItem(row adminDto.AdminCampaignRow) adminDto.AdminCampaignListItem {
	return adminDto.AdminCampaignListItem{
		ID:               row.ID,
		Title:            row.Title,
		BusinessName:     row.BusinessName,
		Niche:            row.Niche,
		BudgetSummary:    FormatBudgetSummary(row.BudgetMinETB, row.BudgetMaxETB),
		CreatorsInvited:  row.CreatorsInvited,
		Status:           adminDto.MapAdminCampaignStatus(row.Status),
		CreatedAt:        row.CreatedAt,
		ContractSummary:  FormatContractSummary(row.ContractsActive, row.ContractsPending),
		PaymentSummary:   FormatPaymentSummary(row.PaymentsPaid, row.PaymentsPending),
		ContractsActive:  row.ContractsActive,
		ContractsPending: row.ContractsPending,
	}
}

func MapCampaignDetail(row adminDto.AdminCampaignRow) adminDto.AdminCampaignDetailResponse {
	item := MapCampaignListItem(row)
	return adminDto.AdminCampaignDetailResponse{
		AdminCampaignListItem: item,
		Goal:                  row.Goal,
		Description:           row.Description,
		BusinessUserID:        row.BusinessUserID,
		BudgetMinETB:          row.BudgetMinETB,
		BudgetMaxETB:          row.BudgetMaxETB,
		CreatorsNeeded:        row.CreatorsNeeded,
		InvitesAccepted:       row.InvitesAccepted,
		StartDate:             row.StartDate,
		EndDate:               row.EndDate,
		UpdatedAt:             row.UpdatedAt,
	}
}
