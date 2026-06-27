package business

import (
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"adinfluence_service_app/internal/module/business/helper"
	"context"
)

func (m *businessModule) GetDashboardSummary(ctx context.Context, userID string, user *authDto.UserResponse) (*businessDto.DashboardSummaryResponse, error) {
	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}

	profile, err := m.business.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	kpis, err := m.business.GetDashboardKPIs(ctx, userID)
	if err != nil {
		return nil, err
	}

	campaigns, err := m.business.ListDashboardRecentCampaigns(ctx, userID)
	if err != nil {
		return nil, err
	}

	recommendations, err := m.business.ListDashboardTopRecommendations(ctx, userID)
	if err != nil {
		return nil, err
	}

	detail := profile.ProfileDetailResponse
	marketplace := helper.BuildMarketplaceAccess(detail, user.Status)

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
