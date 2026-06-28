package influencer

import (
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
	"adinfluence_service_app/internal/module/influencer/helper"
	"context"
)

func (m *influencerModule) GetDashboardSummary(ctx context.Context, userID string) (*influencerDto.DashboardSummaryResponse, error) {
	if err := m.requireActiveInfluencer(ctx, userID); err != nil {
		return nil, err
	}

	rec, err := m.influencer.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	kpis, err := m.influencer.GetDashboardKPIs(ctx, userID)
	if err != nil {
		return nil, err
	}

	invites, err := m.influencer.ListDashboardPendingInvites(ctx, userID)
	if err != nil {
		return nil, err
	}

	contracts, err := m.influencer.ListDashboardActiveContracts(ctx, userID)
	if err != nil {
		return nil, err
	}

	detail := rec.ProfileDetailResponse
	platforms, err := m.influencer.GetPlatformsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	completeness := helper.BuildCompleteness(helper.CompletenessInput{
		Profile:   detail,
		Platforms: platforms,
	})

	completedEarnings := kpis.PaidOutEarningsETB
	if completedEarnings == 0 {
		completedEarnings = kpis.TotalCompletedEarningsETB
	}

	return &influencerDto.DashboardSummaryResponse{
		Greeting: influencerDto.DashboardGreeting{
			CreatorName: detail.CreatorName,
		},
		Completeness: completeness,
		KPIs: influencerDto.DashboardKPIs{
			PendingInvites:            kpis.PendingInvites,
			ActiveCollaborations:      kpis.ActiveCollaborations,
			PendingContracts:          kpis.PendingContracts,
			CompletedContracts:        kpis.CompletedContracts,
			ProfileScorePercent:       completeness.Percent,
			TotalCompletedEarningsETB: kpis.TotalCompletedEarningsETB,
		},
		ContractsSummary: influencerDto.DashboardContractsSummary{
			Pending:   kpis.PendingContracts,
			Active:    kpis.ActiveCollaborations,
			Completed: kpis.CompletedContracts,
		},
		EarningsSummary: influencerDto.DashboardEarningsSummary{
			PendingPayoutETB:     kpis.PendingPayoutETB,
			CompletedEarningsETB: completedEarnings,
			NextPayoutAt:         kpis.NextPayoutAt,
		},
		PendingInvites:  invites,
		ActiveContracts: contracts,
	}, nil
}
