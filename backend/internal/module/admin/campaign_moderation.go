package admin

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	"adinfluence_service_app/internal/module/admin/helper"
	adminPersist "adinfluence_service_app/internal/storage/persist/admin"
	"context"
	"errors"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (m *adminModule) ListAdminCampaigns(ctx context.Context, query adminDto.ListAdminCampaignsQuery) (*adminDto.ListAdminCampaignsResponse, error) {
	m.logger.Info(ctx, "admin listing campaigns for moderation", zap.String("status", query.Status))

	query = query.Normalize()

	var (
		rows  []adminDto.AdminCampaignRow
		total int
		stats *adminDto.AdminCampaignStats
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		rows, total, err = m.admin.ListAdminCampaigns(ctx, query)
		return err
	})

	g.Go(func() error {
		var err error
		stats, err = m.admin.GetAdminCampaignStats(ctx)
		return err
	})

	if err := g.Wait(); err != nil {
		m.logger.Error(ctx, "failed to query admin campaigns or stats", zap.String("reason", err.Error()))
		return nil, err
	}

	items := make([]adminDto.AdminCampaignListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, helper.MapCampaignListItem(row))
	}

	return &adminDto.ListAdminCampaignsResponse{
		Items:      items,
		Stats:      *stats,
		Pagination: commonDto.NewListMeta(query.Page, query.Limit, total),
	}, nil
}

func (m *adminModule) GetAdminCampaign(ctx context.Context, campaignID string) (*adminDto.AdminCampaignDetailResponse, error) {
	m.logger.Info(ctx, "admin fetching campaign moderation details", zap.String("campaign_id", campaignID))

	row, err := m.admin.GetAdminCampaign(ctx, campaignID)
	if err != nil {
		m.logger.Warn(ctx, "failed to get campaign for admin", zap.String("campaign_id", campaignID), zap.String("reason", err.Error()))
		return nil, mapCampaignError(err)
	}
	res := helper.MapCampaignDetail(*row)
	return &res, nil
}

func (m *adminModule) PauseAdminCampaign(ctx context.Context, adminUserID, campaignID string, meta authDto.SessionMeta) (*adminDto.ModerateCampaignResponse, error) {
	m.logger.Info(ctx, "admin pausing campaign", zap.String("admin_id", adminUserID), zap.String("campaign_id", campaignID))
	return m.moderateCampaign(ctx, adminUserID, campaignID, "pending", "campaign paused by admin", meta)
}

func (m *adminModule) CompleteAdminCampaign(ctx context.Context, adminUserID, campaignID string, meta authDto.SessionMeta) (*adminDto.ModerateCampaignResponse, error) {
	m.logger.Info(ctx, "admin completing campaign", zap.String("admin_id", adminUserID), zap.String("campaign_id", campaignID))
	return m.moderateCampaign(ctx, adminUserID, campaignID, "completed", "campaign marked completed by admin", meta)
}

func (m *adminModule) FlagAdminCampaign(ctx context.Context, adminUserID, campaignID string, meta authDto.SessionMeta) (*adminDto.ModerateCampaignResponse, error) {
	m.logger.Info(ctx, "admin flagging campaign", zap.String("admin_id", adminUserID), zap.String("campaign_id", campaignID))
	return m.moderateCampaign(ctx, adminUserID, campaignID, "paused", "campaign flagged by admin", meta)
}

func (m *adminModule) moderateCampaign(ctx context.Context, adminUserID, campaignID, status, desc string, meta authDto.SessionMeta) (*adminDto.ModerateCampaignResponse, error) {
	if err := m.admin.SetAdminCampaignStatus(ctx, campaignID, status); err != nil {
		m.logger.Error(ctx, "failed to update campaign status in persistence", zap.String("campaign_id", campaignID), zap.String("reason", err.Error()))
		return nil, mapCampaignError(err)
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID: &adminUserID,
		Action:      string(model.StatusChangeAction),
		EntityType:  "campaign",
		EntityID:    &campaignID,
		Description: desc,
		IPAddress:   meta.IPAddress,
		UserAgent:   meta.UserAgent,
	})

	m.logger.Info(ctx, "campaign status updated successfully by admin", zap.String("campaign_id", campaignID), zap.String("new_status", status))

	return &adminDto.ModerateCampaignResponse{
		ID:     campaignID,
		Status: adminDto.MapAdminCampaignStatus(status),
	}, nil
}

func mapCampaignError(err error) error {
	if errors.Is(err, adminPersist.ErrAdminCampaignNotFound) {
		return customErr.ErrCampaignNotFound
	}
	return err
}
