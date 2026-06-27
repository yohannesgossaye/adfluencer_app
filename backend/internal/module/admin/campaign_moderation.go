package admin

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	"adinfluence_service_app/internal/module/admin/helper"
	adminPersist "adinfluence_service_app/internal/storage/persist/admin"
	"context"
	"errors"
)

func (m *adminModule) ListAdminCampaigns(ctx context.Context, query adminDto.ListAdminCampaignsQuery) (*adminDto.ListAdminCampaignsResponse, error) {
	query = query.Normalize()
	rows, total, err := m.admin.ListAdminCampaigns(ctx, query)
	if err != nil {
		return nil, err
	}

	stats, err := m.admin.GetAdminCampaignStats(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]adminDto.AdminCampaignListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, helper.MapCampaignListItem(row))
	}

	return &adminDto.ListAdminCampaignsResponse{
		Items: items,
		Stats: *stats,
		Pagination: adminDto.PaginationResponse{
			Page:  query.Page,
			Limit: query.Limit,
			Total: total,
		},
	}, nil
}

func (m *adminModule) GetAdminCampaign(ctx context.Context, campaignID string) (*adminDto.AdminCampaignDetailResponse, error) {
	row, err := m.admin.GetAdminCampaign(ctx, campaignID)
	if err != nil {
		return nil, mapCampaignError(err)
	}
	res := helper.MapCampaignDetail(*row)
	return &res, nil
}

func (m *adminModule) PauseAdminCampaign(ctx context.Context, adminUserID, campaignID string, meta authDto.SessionMeta) (*adminDto.ModerateCampaignResponse, error) {
	return m.moderateCampaign(ctx, adminUserID, campaignID, "pending", "campaign paused by admin", meta)
}

func (m *adminModule) CompleteAdminCampaign(ctx context.Context, adminUserID, campaignID string, meta authDto.SessionMeta) (*adminDto.ModerateCampaignResponse, error) {
	return m.moderateCampaign(ctx, adminUserID, campaignID, "completed", "campaign marked completed by admin", meta)
}

func (m *adminModule) FlagAdminCampaign(ctx context.Context, adminUserID, campaignID string, meta authDto.SessionMeta) (*adminDto.ModerateCampaignResponse, error) {
	return m.moderateCampaign(ctx, adminUserID, campaignID, "paused", "campaign flagged by admin", meta)
}

func (m *adminModule) moderateCampaign(ctx context.Context, adminUserID, campaignID, status, desc string, meta authDto.SessionMeta) (*adminDto.ModerateCampaignResponse, error) {
	if err := m.admin.SetAdminCampaignStatus(ctx, campaignID, status); err != nil {
		return nil, mapCampaignError(err)
	}

	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID: &adminUserID,
		Action:      string(model.StatusChangeAction),
		EntityType:  "campaign",
		EntityID:    &campaignID,
		Description: desc,
		IPAddress:   meta.IPAddress,
		UserAgent:   meta.UserAgent,
	})

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
