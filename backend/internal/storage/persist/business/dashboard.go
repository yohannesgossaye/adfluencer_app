package business

import (
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	"context"
	"database/sql"
	"fmt"
)

func (p *business) GetDashboardKPIs(ctx context.Context, businessUserID string) (*businessDto.DashboardKPIRow, error) {
	var row businessDto.DashboardKPIRow
	err := p.db.QueryRowContext(ctx, getDashboardKPIs, businessUserID).Scan(
		&row.ActiveCampaigns,
		&row.TotalCampaigns,
		&row.RecommendedCreatorsTotal,
		&row.RecommendedCreatorsHighMatch,
		&row.PendingContracts,
		&row.ContractsAwaitingInfluencer,
		&row.ContractsAwaitingPayment,
		&row.PendingInvites,
		&row.ContractsActive,
		&row.ContractsCompleted,
	)
	if err != nil {
		return nil, fmt.Errorf("get dashboard kpis: %w", err)
	}
	return &row, nil
}

func (p *business) ListDashboardRecentCampaigns(ctx context.Context, businessUserID string) ([]businessDto.DashboardCampaignRow, error) {
	rows, err := p.db.QueryContext(ctx, listDashboardRecentCampaigns, businessUserID)
	if err != nil {
		return nil, fmt.Errorf("list dashboard campaigns: %w", err)
	}
	defer rows.Close()

	items := make([]businessDto.DashboardCampaignRow, 0, 5)
	for rows.Next() {
		var item businessDto.DashboardCampaignRow
		var budgetMin, budgetMax sql.NullFloat64
		var niche string

		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Status,
			&budgetMin,
			&budgetMax,
			&niche,
			&item.CreatorsNeeded,
			&item.AcceptedCreators,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan dashboard campaign: %w", err)
		}

		item.Niche = niche
		var minPtr, maxPtr *float64
		if budgetMin.Valid {
			v := budgetMin.Float64
			minPtr = &v
		}
		if budgetMax.Valid {
			v := budgetMax.Float64
			maxPtr = &v
		}
		item.BudgetSummary = formatBudgetSummary(minPtr, maxPtr)
		item.ProgressPercent = campaignProgressPercent(item.AcceptedCreators, item.CreatorsNeeded)
		items = append(items, item)
	}
	if items == nil {
		items = []businessDto.DashboardCampaignRow{}
	}
	return items, rows.Err()
}

func (p *business) ListDashboardTopRecommendations(ctx context.Context, businessUserID string) ([]businessDto.DashboardRecommendationRow, error) {
	rows, err := p.db.QueryContext(ctx, listDashboardTopRecommendations, businessUserID)
	if err != nil {
		return nil, fmt.Errorf("list dashboard recommendations: %w", err)
	}
	defer rows.Close()

	items := make([]businessDto.DashboardRecommendationRow, 0, 5)
	for rows.Next() {
		var item businessDto.DashboardRecommendationRow
		var profileImage sql.NullString
		if err := rows.Scan(
			&item.InfluencerUserID,
			&item.CreatorName,
			&profileImage,
			&item.MatchScore,
			&item.CampaignID,
			&item.CampaignTitle,
		); err != nil {
			return nil, fmt.Errorf("scan dashboard recommendation: %w", err)
		}
		if profileImage.Valid {
			item.ProfileImageURL = &profileImage.String
		}
		items = append(items, item)
	}
	if items == nil {
		items = []businessDto.DashboardRecommendationRow{}
	}
	return items, rows.Err()
}

func campaignProgressPercent(accepted, needed int) *int {
	if needed <= 0 {
		return nil
	}
	pct := accepted * 100 / needed
	if pct > 100 {
		pct = 100
	}
	return &pct
}
