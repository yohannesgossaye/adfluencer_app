package admin

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	supportDto "adinfluence_service_app/internal/constant/model/dto/support"
	"adinfluence_service_app/internal/module/admin/helper"
	adminPersist "adinfluence_service_app/internal/storage/persist/admin"
	"context"
	"errors"
	"strings"
)

func (m *adminModule) ListAdminDisputes(ctx context.Context, query adminDto.ListAdminDisputesQuery) (*adminDto.ListAdminDisputesResponse, error) {
	query = query.Normalize()
	rows, total, err := m.admin.ListAdminDisputes(ctx, query)
	if err != nil {
		return nil, err
	}

	items := make([]adminDto.AdminDisputeListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapDisputeListItem(row))
	}

	return &adminDto.ListAdminDisputesResponse{
		Items: items,
		Pagination: commonDto.NewListMeta(query.Page, query.Limit, total),
	}, nil
}

func (m *adminModule) GetAdminDispute(ctx context.Context, disputeID string) (*adminDto.AdminDisputeDetailResponse, error) {
	row, err := m.admin.GetAdminDispute(ctx, disputeID)
	if err != nil {
		return nil, mapDisputeError(err)
	}
	res := mapDisputeDetail(*row)
	return &res, nil
}

func (m *adminModule) UpdateDisputeStatus(ctx context.Context, adminUserID, disputeID string, req adminDto.UpdateDisputeStatusRequest, meta authDto.SessionMeta) (*adminDto.UpdateDisputeStatusResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	status := strings.ToLower(strings.TrimSpace(req.Status))

	if err := m.admin.UpdateDisputeStatus(ctx, disputeID, status, adminUserID); err != nil {
		return nil, mapDisputeError(err)
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID: &adminUserID,
		Action:      string(model.StatusChangeAction),
		EntityType:  "dispute",
		EntityID:    &disputeID,
		Description: "dispute status updated to " + status,
		IPAddress:   meta.IPAddress,
		UserAgent:   meta.UserAgent,
	})

	return &adminDto.UpdateDisputeStatusResponse{ID: disputeID, Status: status}, nil
}

func (m *adminModule) UpdateDisputeNotes(ctx context.Context, adminUserID, disputeID string, req adminDto.UpdateDisputeNotesRequest, meta authDto.SessionMeta) (*adminDto.UpdateDisputeNotesResponse, error) {
	if req.AdminNote == nil && req.ResolutionNote == nil {
		return nil, customErr.ErrBadRequest
	}

	if req.ResolutionNote != nil {
		note := strings.TrimSpace(*req.ResolutionNote)
		var notePtr *string
		if note != "" {
			notePtr = &note
		}
		if err := m.admin.UpdateDisputeResolutionNote(ctx, disputeID, notePtr); err != nil {
			return nil, mapDisputeError(err)
		}
	}

	adminNote := ""
	if req.AdminNote != nil {
		adminNote = strings.TrimSpace(*req.AdminNote)
		if adminNote != "" {
			if err := m.admin.SaveDisputeAdminNote(ctx, adminUserID, disputeID, adminNote); err != nil {
				return nil, err
			}
		}
	}

	row, err := m.admin.GetDisputeNotes(ctx, disputeID)
	if err != nil {
		return nil, mapDisputeError(err)
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID: &adminUserID,
		Action:      string(model.AdminAction),
		EntityType:  "dispute",
		EntityID:    &disputeID,
		Description: "dispute notes updated",
		IPAddress:   meta.IPAddress,
		UserAgent:   meta.UserAgent,
	})

	return &adminDto.UpdateDisputeNotesResponse{
		ID:             disputeID,
		AdminNote:      row.AdminNote,
		ResolutionNote: row.ResolutionNote,
	}, nil
}

func mapDisputeListItem(row adminDto.DisputeRow) adminDto.AdminDisputeListItem {
	label := supportDto.IssueTypeLabel(row.IssueType)
	return adminDto.AdminDisputeListItem{
		ID:             row.ID,
		Title:          helper.DisputeTitle(label, row.CampaignTitle),
		BusinessName:   row.BusinessName,
		InfluencerName: row.InfluencerName,
		CampaignTitle:  row.CampaignTitle,
		IssueType:      row.IssueType,
		IssueTypeLabel: label,
		Severity:       row.Severity,
		Status:         row.Status,
		Preview:        helper.DisputePreview(row.Description),
		CreatedAt:      row.CreatedAt,
	}
}

func mapDisputeDetail(row adminDto.DisputeRow) adminDto.AdminDisputeDetailResponse {
	label := supportDto.IssueTypeLabel(row.IssueType)
	return adminDto.AdminDisputeDetailResponse{
		ID:               row.ID,
		Title:            helper.DisputeTitle(label, row.CampaignTitle),
		BusinessName:     row.BusinessName,
		InfluencerName:   row.InfluencerName,
		InfluencerEmail:  row.InfluencerEmail,
		InfluencerUserID: row.InfluencerUserID,
		CampaignTitle:    row.CampaignTitle,
		CampaignID:       row.CampaignID,
		ContractID:       row.ContractID,
		IssueType:        row.IssueType,
		IssueTypeLabel:   label,
		Severity:         row.Severity,
		Status:           row.Status,
		Description:      row.Description,
		Preview:          helper.DisputePreview(row.Description),
		AdminNote:        row.AdminNote,
		ResolutionNote:   row.ResolutionNote,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
		ResolvedAt:       row.ResolvedAt,
	}
}

func mapDisputeError(err error) error {
	if errors.Is(err, adminPersist.ErrAdminDisputeNotFound) {
		return customErr.ErrDisputeNotFound
	}
	return err
}
