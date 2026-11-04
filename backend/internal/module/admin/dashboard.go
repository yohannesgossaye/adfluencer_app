package admin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	"context"
	"strings"

	"go.uber.org/zap"
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

	// Run one after another: each query takes a few ms, and running them in
	// parallel would hold several pool connections for one page load.
	err := func() (err error) {
		if kpis, err = m.admin.GetDashboardKPIs(ctx); err != nil {
			return err
		}
		if pending, err = m.admin.ListDashboardPendingApprovals(ctx, dashboardPendingLimit); err != nil {
			return err
		}
		if activity, err = m.admin.ListDashboardRecentActivity(ctx, dashboardActivityLimit); err != nil {
			return err
		}
		if trend, err = m.admin.ListDashboardMonthlyRevenue(ctx, dashboardTrendMonths); err != nil {
			return err
		}
		if alerts, err = m.admin.ListDashboardAlerts(ctx, dashboardAlertLimit); err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
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
			PaymentsToConfirm:          kpis.PaymentsToConfirm,
			BusinessesToReview:         kpis.BusinessesToReview,
			UrgentSupportTickets:       kpis.UrgentSupportTickets,
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
