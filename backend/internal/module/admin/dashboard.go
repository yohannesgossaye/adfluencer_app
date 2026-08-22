package admin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	"context"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	dashboardPendingLimit  = 3
	dashboardActivityLimit = 8
	dashboardAlertLimit    = 6
	dashboardTrendMonths   = 6
)

func (m *adminModule) GetDashboardSummary(ctx context.Context) (*adminDto.DashboardSummaryResponse, error) {
	m.logger.Info(ctx, "fetching admin dashboard summary")

	var (
		kpis     *adminDto.DashboardKPIRow
		pending  []adminDto.DashboardPendingApprovalRow
		activity []adminDto.DashboardActivityRow
		trend    []adminDto.DashboardMonthlyRevenueRow
		alerts   []adminDto.DashboardAlertRow
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		kpis, err = m.admin.GetDashboardKPIs(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		pending, err = m.admin.ListDashboardPendingApprovals(ctx, dashboardPendingLimit)
		return err
	})

	g.Go(func() error {
		var err error
		activity, err = m.admin.ListDashboardRecentActivity(ctx, dashboardActivityLimit)
		return err
	})

	g.Go(func() error {
		var err error
		trend, err = m.admin.ListDashboardMonthlyRevenue(ctx, dashboardTrendMonths)
		return err
	})

	g.Go(func() error {
		var err error
		alerts, err = m.admin.ListDashboardAlerts(ctx, dashboardAlertLimit)
		return err
	})

	if err := g.Wait(); err != nil {
		m.logger.Error(ctx, "failed to load admin dashboard summary data", zap.String("reason", err.Error()))
		return nil, err
	}

	for i := range activity {
		activity[i].Label = activityLabel(activity[i].Action)
		if strings.TrimSpace(activity[i].Detail) == "" {
			activity[i].Detail = activity[i].Label
		}
	}

	return &adminDto.DashboardSummaryResponse{
		KPIs: adminDto.DashboardKPIs{
			TotalInfluencers:           kpis.TotalInfluencers,
			TotalBusinesses:            kpis.TotalBusinesses,
			PendingInfluencerApprovals: kpis.PendingInfluencerApprovals,
			ActiveCampaigns:            kpis.ActiveCampaigns,
			TotalCampaigns:             kpis.TotalCampaigns,
			TotalUsers:                 kpis.TotalUsers,
			OpenDisputes:               kpis.OpenDisputes,
			OpenSupportTickets:         kpis.OpenSupportTickets,
			BlockedUsers:               kpis.BlockedUsers,
			CompletedContracts:         kpis.CompletedContracts,
			MonthlyPlatformRevenueETB:  kpis.MonthlyPlatformRevenueETB,
		},
		PendingApprovals: pending,
		RecentActivity:   activity,
		RevenueSummary: adminDto.DashboardRevenueSummary{
			TotalCampaignValueETB:        kpis.TotalCampaignValueETB,
			TotalPlatformFeesETB:         kpis.TotalPlatformFeesETB,
			PlatformFeePercentBusiness:   5.0,
			PlatformFeePercentInfluencer: 5.0,
			MonthlyTrend:                 trend,
		},
		Alerts: alerts,
	}, nil
}

func activityLabel(action string) string {
	switch action {
	case "approve":
		return "Application approved"
	case "reject":
		return "Application rejected"
	case "block":
		return "User blocked"
	case "unblock":
		return "User unblocked"
	case "delete":
		return "User removed"
	case "admin":
		return "Admin action"
	case "create":
		return "Record created"
	case "update":
		return "Record updated"
	default:
		if action == "" {
			return "Platform activity"
		}
		return strings.ReplaceAll(action, "_", " ")
	}
}
