package admin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	"context"
	"strings"
)

const (
	dashboardPendingLimit  = 3
	dashboardActivityLimit = 8
	dashboardAlertLimit    = 6
	dashboardTrendMonths   = 6
)

func (m *adminModule) GetDashboardSummary(ctx context.Context) (*adminDto.DashboardSummaryResponse, error) {
	kpis, err := m.admin.GetDashboardKPIs(ctx)
	if err != nil {
		return nil, err
	}

	pending, err := m.admin.ListDashboardPendingApprovals(ctx, dashboardPendingLimit)
	if err != nil {
		return nil, err
	}

	activity, err := m.admin.ListDashboardRecentActivity(ctx, dashboardActivityLimit)
	if err != nil {
		return nil, err
	}
	for i := range activity {
		activity[i].Label = activityLabel(activity[i].Action)
		if strings.TrimSpace(activity[i].Detail) == "" {
			activity[i].Detail = activity[i].Label
		}
	}

	trend, err := m.admin.ListDashboardMonthlyRevenue(ctx, dashboardTrendMonths)
	if err != nil {
		return nil, err
	}

	alerts, err := m.admin.ListDashboardAlerts(ctx, dashboardAlertLimit)
	if err != nil {
		return nil, err
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
			TotalCampaignValueETB:          kpis.TotalCampaignValueETB,
			TotalPlatformFeesETB:           kpis.TotalPlatformFeesETB,
			PlatformFeePercentBusiness:     5.0,
			PlatformFeePercentInfluencer:   5.0,
			MonthlyTrend:                   trend,
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
