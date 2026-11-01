package business

import (
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	"adinfluence_service_app/internal/constant/model"
	"adinfluence_service_app/internal/module/business/helper"
	"context"

	"go.uber.org/zap"
)

func (m *businessModule) GetDashboardSummary(ctx context.Context, userID string) (*businessDto.DashboardSummaryResponse, error) {
	m.logger.Info(ctx, "fetching business dashboard summary", zap.String("user_id", userID))

	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		m.logger.Warn(ctx, "business dashboard access denied", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	var (
		profile         *businessDto.ProfileRecord
		kpis            *businessDto.DashboardKPIRow
		campaigns       []businessDto.DashboardCampaignRow
		recommendations []businessDto.DashboardRecommendationRow
	)

	// Run one after another: each query takes a few ms, and running them in
	// parallel would hold several pool connections for one page load.
	err := func() (err error) {
		if profile, err = m.business.GetProfileByUserID(ctx, userID); err != nil {
			return err
		}
		if kpis, err = m.business.GetDashboardKPIs(ctx, userID); err != nil {
			return err
		}
		if campaigns, err = m.business.ListDashboardRecentCampaigns(ctx, userID); err != nil {
			return err
		}
		if recommendations, err = m.business.ListDashboardTopRecommendations(ctx, userID); err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		m.logger.Error(ctx, "failed to load business dashboard summary data", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	detail := profile.ProfileDetailResponse
	marketplace := helper.BuildMarketplaceAccess(detail, string(model.ActiveStatus))

	return &businessDto.DashboardSummaryResponse{
		Greeting: businessDto.DashboardGreeting{
			BusinessName:  detail.BusinessName,
			OwnerFullName: detail.OwnerFullName,
		},
		MarketplaceAccess: marketplace,
		KPIs: businessDto.DashboardKPIs{
			ActiveCampaigns:              kpis.ActiveCampaigns,
			TotalCampaigns:               kpis.TotalCampaigns,
			RecommendedCreatorsTotal:     kpis.RecommendedCreatorsTotal,
			RecommendedCreatorsHighMatch: kpis.RecommendedCreatorsHighMatch,
			PendingContracts:             kpis.PendingContracts,
			ContractsAwaitingInfluencer:  kpis.ContractsAwaitingInfluencer,
			ContractsAwaitingPayment:     kpis.ContractsAwaitingPayment,
			PendingInvites:               kpis.PendingInvites,
		},
		ContractsSummary: businessDto.DashboardContractsSummary{
			Pending:   kpis.PendingContracts,
			Active:    kpis.ContractsActive,
			Completed: kpis.ContractsCompleted,
		},
		RecentCampaigns:    campaigns,
		TopRecommendations: recommendations,
	}, nil
}
