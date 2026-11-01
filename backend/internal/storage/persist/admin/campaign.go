package admin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrAdminCampaignNotFound      = errors.New("campaign not found")
	ErrAdminCampaignInvalidStatus = errors.New("campaign status does not allow this action")
)

func (a *admin) GetAdminCampaignStats(ctx context.Context) (*adminDto.AdminCampaignStats, error) {
	var stats adminDto.AdminCampaignStats
	if err := a.db.QueryRowContext(ctx, getAdminCampaignStats).Scan(
		&stats.Active,
		&stats.Draft,
		&stats.Paused,
	); err != nil {
		return nil, fmt.Errorf("get admin campaign stats: %w", err)
	}
	return &stats, nil
}

func (a *admin) ListAdminCampaigns(ctx context.Context, query adminDto.ListAdminCampaignsQuery) ([]adminDto.AdminCampaignRow, int, error) {
	query = query.Normalize()
	dbStatus := adminDto.AdminCampaignStatusFilterDB(query.Status)
	offset := (query.Page - 1) * query.Limit

	rows, err := a.db.QueryContext(ctx, listAdminCampaigns, dbStatus, query.Niche, query.Business, query.Limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list admin campaigns: %w", err)
	}
	defer rows.Close()

	var total int
	items := make([]adminDto.AdminCampaignRow, 0)
	for rows.Next() {
		row, err := scanAdminCampaignListRow(rows, &total)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *row)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (a *admin) GetAdminCampaign(ctx context.Context, campaignID string) (*adminDto.AdminCampaignRow, error) {
	row, err := scanAdminCampaignDetailRow(a.db.QueryRowContext(ctx, getAdminCampaign, campaignID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAdminCampaignNotFound
		}
		return nil, fmt.Errorf("get admin campaign: %w", err)
	}
	return row, nil
}

// SetAdminCampaignStatus moves a campaign to status if it is currently in one
// of the from statuses. ErrAdminCampaignInvalidStatus means it exists but is
// in a state that does not allow this change.
func (a *admin) SetAdminCampaignStatus(ctx context.Context, campaignID, status string, from []string) error {
	result, err := a.db.ExecContext(ctx, updateAdminCampaignStatus, campaignID, status, "{"+strings.Join(from, ",")+"}")
	if err != nil {
		return fmt.Errorf("update admin campaign status: %w", err)
	}
	if n, _ := result.RowsAffected(); n > 0 {
		return nil
	}
	var exists bool
	if err := a.db.QueryRowContext(ctx, adminCampaignExists, campaignID).Scan(&exists); err != nil {
		return fmt.Errorf("check campaign exists: %w", err)
	}
	if !exists {
		return ErrAdminCampaignNotFound
	}
	return ErrAdminCampaignInvalidStatus
}

func scanAdminCampaignListRow(s scanner, total *int) (*adminDto.AdminCampaignRow, error) {
	var row adminDto.AdminCampaignRow
	var budgetMin, budgetMax sql.NullFloat64
	if err := s.Scan(
		&row.ID,
		&row.Title,
		&row.BusinessName,
		&row.BusinessUserID,
		&row.Niche,
		&budgetMin,
		&budgetMax,
		&row.Status,
		&row.CreatorsNeeded,
		&row.CreatedAt,
		&row.CreatorsInvited,
		&row.InvitesAccepted,
		&row.ContractsActive,
		&row.ContractsPending,
		&row.PaymentsPaid,
		&row.PaymentsPending,
		total,
	); err != nil {
		return nil, err
	}
	row.BudgetMinETB = nullFloatPtr(budgetMin)
	row.BudgetMaxETB = nullFloatPtr(budgetMax)
	return &row, nil
}

func scanAdminCampaignDetailRow(s scanner) (*adminDto.AdminCampaignRow, error) {
	var row adminDto.AdminCampaignRow
	var budgetMin, budgetMax sql.NullFloat64
	var goal sql.NullString
	var startDate, endDate sql.NullTime
	if err := s.Scan(
		&row.ID,
		&row.Title,
		&row.BusinessName,
		&row.BusinessUserID,
		&row.Niche,
		&budgetMin,
		&budgetMax,
		&row.Status,
		&row.CreatorsNeeded,
		&row.CreatedAt,
		&row.UpdatedAt,
		&goal,
		&row.Description,
		&startDate,
		&endDate,
		&row.CreatorsInvited,
		&row.InvitesAccepted,
		&row.ContractsActive,
		&row.ContractsPending,
		&row.PaymentsPaid,
		&row.PaymentsPending,
	); err != nil {
		return nil, err
	}
	row.BudgetMinETB = nullFloatPtr(budgetMin)
	row.BudgetMaxETB = nullFloatPtr(budgetMax)
	if goal.Valid {
		row.Goal = &goal.String
	}
	row.StartDate = nullTimePtr(startDate)
	row.EndDate = nullTimePtr(endDate)
	return &row, nil
}

func nullFloatPtr(v sql.NullFloat64) *float64 {
	if !v.Valid {
		return nil
	}
	f := v.Float64
	return &f
}

func nullTimePtr(v sql.NullTime) *time.Time {
	if !v.Valid {
		return nil
	}
	t := v.Time
	return &t
}
