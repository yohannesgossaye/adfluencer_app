package influencer

import (
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
	"adinfluence_service_app/internal/module/influencer/helper"
	"context"

	"golang.org/x/sync/errgroup"
)

func (m *influencerModule) GetDashboardSummary(ctx context.Context, userID string) (*influencerDto.DashboardSummaryResponse, error) {
	if err := m.requireActiveInfluencer(ctx, userID); err != nil {
		return nil, err
	}

	var (
		rec       *influencerDto.ProfileRecord
		kpis      *influencerDto.DashboardKPIRow
		invites   []influencerDto.DashboardInviteRow
		contracts []influencerDto.DashboardContractRow
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		rec, err = m.influencer.GetProfileByUserID(ctx, userID)
		return err
	})

	g.Go(func() error {
		var err error
		kpis, err = m.influencer.GetDashboardKPIs(ctx, userID)
		return err
	})

	g.Go(func() error {
		var err error
		invites, err = m.influencer.ListDashboardPendingInvites(ctx, userID)
		return err
	})

	g.Go(func() error {
		var err error
		contracts, err = m.influencer.ListDashboardActiveContracts(ctx, userID)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	detail := rec.ProfileDetailResponse
	completeness := helper.BuildCompleteness(helper.CompletenessInput{
		Profile:   detail,
		Platforms: rec.Platforms,
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
