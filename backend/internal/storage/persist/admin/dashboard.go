package admin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

func (a *admin) GetDashboardKPIs(ctx context.Context) (*adminDto.DashboardKPIRow, error) {
	var row adminDto.DashboardKPIRow
	err := a.db.QueryRowContext(ctx, getAdminDashboardKPIs).Scan(
		&row.TotalInfluencers,
		&row.TotalBusinesses,
		&row.PendingInfluencerApprovals,
		&row.ActiveCampaigns,
		&row.TotalCampaigns,
		&row.TotalUsers,
		&row.OpenDisputes,
		&row.OpenSupportTickets,
		&row.BlockedUsers,
		&row.CompletedContracts,
		&row.MonthlyPlatformRevenueETB,
		&row.TotalCampaignValueETB,
		&row.TotalPlatformFeesETB,
	)
	if err != nil {
		return nil, fmt.Errorf("get admin dashboard kpis: %w", err)
	}
	return &row, nil
}

func (a *admin) ListDashboardPendingApprovals(ctx context.Context, limit int) ([]adminDto.DashboardPendingApprovalRow, error) {
	rows, err := a.db.QueryContext(ctx, listDashboardPendingApprovals, limit)
	if err != nil {
		return nil, fmt.Errorf("list dashboard pending approvals: %w", err)
	}
	defer rows.Close()

	items := make([]adminDto.DashboardPendingApprovalRow, 0, limit)
	for rows.Next() {
		var item adminDto.DashboardPendingApprovalRow
		if err := rows.Scan(
			&item.UserID,
			&item.Name,
			&item.Followers,
			&item.SubmittedAt,
			&item.Niche,
			&item.Platform,
		); err != nil {
			return nil, fmt.Errorf("scan pending approval: %w", err)
		}
		items = append(items, item)
	}
	if items == nil {
		items = []adminDto.DashboardPendingApprovalRow{}
	}
	return items, rows.Err()
}

func (a *admin) ListDashboardRecentActivity(ctx context.Context, limit int) ([]adminDto.DashboardActivityRow, error) {
	rows, err := a.db.QueryContext(ctx, listDashboardRecentActivity, limit)
	if err != nil {
		return nil, fmt.Errorf("list dashboard activity: %w", err)
	}
	defer rows.Close()

	items := make([]adminDto.DashboardActivityRow, 0, limit)
	for rows.Next() {
		var item adminDto.DashboardActivityRow
		if err := rows.Scan(&item.ID, &item.Action, &item.Detail, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan activity row: %w", err)
		}
		items = append(items, item)
	}
	if items == nil {
		items = []adminDto.DashboardActivityRow{}
	}
	return items, rows.Err()
}

func (a *admin) ListDashboardMonthlyRevenue(ctx context.Context, months int) ([]adminDto.DashboardMonthlyRevenueRow, error) {
	if months < 1 {
		months = 6
	}
	rows, err := a.db.QueryContext(ctx, listDashboardMonthlyRevenue, months)
	if err != nil {
		return nil, fmt.Errorf("list dashboard monthly revenue: %w", err)
	}
	defer rows.Close()

	items := make([]adminDto.DashboardMonthlyRevenueRow, 0, months)
	for rows.Next() {
		var item adminDto.DashboardMonthlyRevenueRow
		if err := rows.Scan(&item.Month, &item.PlatformRevenueETB, &item.CompletedContracts); err != nil {
			return nil, fmt.Errorf("scan monthly revenue: %w", err)
		}
		items = append(items, item)
	}
	if items == nil {
		items = []adminDto.DashboardMonthlyRevenueRow{}
	}
	return items, rows.Err()
}

func (a *admin) ListDashboardAlerts(ctx context.Context, limit int) ([]adminDto.DashboardAlertRow, error) {
	var disputes, tickets []adminDto.DashboardAlertRow

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		disputes, err = a.listDisputeAlerts(ctx, limit)
		return err
	})

	g.Go(func() error {
		var err error
		tickets, err = a.listSupportAlerts(ctx, limit)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	items := append(disputes, tickets...)
	if items == nil {
		items = []adminDto.DashboardAlertRow{}
	}
	if len(items) > limit {
		items = items[:limit]
	}
	return items, nil
}

func (a *admin) listDisputeAlerts(ctx context.Context, limit int) ([]adminDto.DashboardAlertRow, error) {
	rows, err := a.db.QueryContext(ctx, listDashboardDisputeAlerts, limit)
	if err != nil {
		return nil, fmt.Errorf("list dispute alerts: %w", err)
	}
	defer rows.Close()

	items := make([]adminDto.DashboardAlertRow, 0, limit)
	for rows.Next() {
		var id, severity, detail string
		if err := rows.Scan(&id, &severity, &detail); err != nil {
			return nil, fmt.Errorf("scan dispute alert: %w", err)
		}
		items = append(items, adminDto.DashboardAlertRow{
			ID:         id,
			Title:      "Open dispute",
			Detail:     detail,
			Severity:   severity,
			EntityType: "dispute",
		})
	}
	return items, rows.Err()
}

func (a *admin) listSupportAlerts(ctx context.Context, limit int) ([]adminDto.DashboardAlertRow, error) {
	rows, err := a.db.QueryContext(ctx, listDashboardSupportAlerts, limit)
	if err != nil {
		return nil, fmt.Errorf("list support alerts: %w", err)
	}
	defer rows.Close()

	items := make([]adminDto.DashboardAlertRow, 0, limit)
	for rows.Next() {
		var id, priority, subject, detail string
		if err := rows.Scan(&id, &priority, &subject, &detail); err != nil {
			return nil, fmt.Errorf("scan support alert: %w", err)
		}
		items = append(items, adminDto.DashboardAlertRow{
			ID:         id,
			Title:      subject,
			Detail:     detail,
			Severity:   supportPriorityToSeverity(priority),
			EntityType: "support_ticket",
		})
	}
	return items, rows.Err()
}

func supportPriorityToSeverity(priority string) string {
	switch priority {
	case "urgent":
		return "high"
	case "high":
		return "medium"
	default:
		return "low"
	}
}
