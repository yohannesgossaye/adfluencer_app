package influencer

import (
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
	"adinfluence_service_app/internal/module/influencer/helper"
	"context"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (m *influencerModule) GetDashboardSummary(ctx context.Context, userID string) (*influencerDto.DashboardSummaryResponse, error) {
	m.logger.Info(ctx, "fetching influencer dashboard summary", zap.String("user_id", userID))

	if err := m.requireActiveInfluencer(ctx, userID); err != nil {
		m.logger.Warn(ctx, "influencer dashboard access rejected: account not active", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	var (
		rec       *influencerDto.ProfileRecord
		kpis      *influencerDto.DashboardKPIRow
		invites   []influencerDto.DashboardInviteRow
		contracts []influencerDto.DashboardContractRow
	)

	g, ct := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		rec, err = m.influencer.GetProfileByUserID(ct, userID)
		return err
	})

	g.Go(func() error {
		var err error
		kpis, err = m.influencer.GetDashboardKPIs(ct, userID)
		return err
	})

	g.Go(func() error {
		var err error
		invites, err = m.influencer.ListDashboardPendingInvites(ct, userID)
		return err
	})

	g.Go(func() error {
		var err error
		contracts, err = m.influencer.ListDashboardActiveContracts(ct, userID)
		return err
	})

	if err := g.Wait(); err != nil {
		m.logger.Error(ctx, "failed to load parallel influencer dashboard data", zap.String("user_id", userID), zap.String("reason", err.Error()))
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
