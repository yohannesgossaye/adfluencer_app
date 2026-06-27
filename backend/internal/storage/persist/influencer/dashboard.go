package influencer

import (
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
	"context"
	"database/sql"
	"fmt"
)

func (p *influencer) GetDashboardKPIs(ctx context.Context, influencerUserID string) (*influencerDto.DashboardKPIRow, error) {
	var row influencerDto.DashboardKPIRow
	var nextPayout sql.NullTime
	err := p.db.QueryRowContext(ctx, getInfluencerDashboardKPIs, influencerUserID).Scan(
		&row.PendingInvites,
		&row.ActiveCollaborations,
		&row.PendingContracts,
		&row.CompletedContracts,
		&row.TotalCompletedEarningsETB,
		&row.PendingPayoutETB,
		&row.PaidOutEarningsETB,
		&nextPayout,
	)
	if err != nil {
		return nil, fmt.Errorf("get influencer dashboard kpis: %w", err)
	}
	if nextPayout.Valid {
		t := nextPayout.Time
		row.NextPayoutAt = &t
	}
	return &row, nil
}

func (p *influencer) ListDashboardPendingInvites(ctx context.Context, influencerUserID string) ([]influencerDto.DashboardInviteRow, error) {
	rows, err := p.db.QueryContext(ctx, listInfluencerDashboardPendingInvites, influencerUserID)
	if err != nil {
		return nil, fmt.Errorf("list dashboard pending invites: %w", err)
	}
	defer rows.Close()

	items := make([]influencerDto.DashboardInviteRow, 0, 3)
	for rows.Next() {
		var item influencerDto.DashboardInviteRow
		var budgetMin, budgetMax sql.NullFloat64
		var endDate sql.NullTime
		if err := rows.Scan(
			&item.ID,
			&item.BusinessName,
			&item.CampaignTitle,
			&item.Status,
			&budgetMin,
			&budgetMax,
			&endDate,
		); err != nil {
			return nil, fmt.Errorf("scan dashboard invite: %w", err)
		}
		item.BudgetSummary = formatBudgetSummary(nullFloatPtr(budgetMin), nullFloatPtr(budgetMax))
		if endDate.Valid {
			t := endDate.Time
			item.EndDate = &t
		}
		items = append(items, item)
	}
	if items == nil {
		items = []influencerDto.DashboardInviteRow{}
	}
	return items, rows.Err()
}

func (p *influencer) ListDashboardActiveContracts(ctx context.Context, influencerUserID string) ([]influencerDto.DashboardContractRow, error) {
	rows, err := p.db.QueryContext(ctx, listInfluencerDashboardActiveContracts, influencerUserID)
	if err != nil {
		return nil, fmt.Errorf("list dashboard active contracts: %w", err)
	}
	defer rows.Close()

	items := make([]influencerDto.DashboardContractRow, 0, 5)
	for rows.Next() {
		var item influencerDto.DashboardContractRow
		if err := rows.Scan(&item.ID, &item.BusinessName, &item.Status); err != nil {
			return nil, fmt.Errorf("scan dashboard contract: %w", err)
		}
		items = append(items, item)
	}
	if items == nil {
		items = []influencerDto.DashboardContractRow{}
	}
	return items, rows.Err()
}

func formatBudgetSummary(min, max *float64) string {
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

func nullFloatPtr(v sql.NullFloat64) *float64 {
	if !v.Valid {
		return nil
	}
	f := v.Float64
	return &f
}
