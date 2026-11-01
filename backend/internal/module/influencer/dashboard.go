package influencer

import (
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
	"adinfluence_service_app/internal/module/influencer/helper"
	"context"

	"go.uber.org/zap"
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

	// Run one after another: each query takes a few ms, and running them in
	// parallel would hold several pool connections for one page load.
	err := func() (err error) {
		if rec, err = m.influencer.GetProfileByUserID(ctx, userID); err != nil {
			return err
		}
		if kpis, err = m.influencer.GetDashboardKPIs(ctx, userID); err != nil {
			return err
		}
		if invites, err = m.influencer.ListDashboardPendingInvites(ctx, userID); err != nil {
			return err
		}
		if contracts, err = m.influencer.ListDashboardActiveContracts(ctx, userID); err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		m.logger.Error(ctx, "failed to load influencer dashboard data", zap.String("user_id", userID), zap.String("reason", err.Error()))
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
