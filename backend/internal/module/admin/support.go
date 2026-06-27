package admin

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	supportDto "adinfluence_service_app/internal/constant/model/dto/support"
	supportPersist "adinfluence_service_app/internal/storage/persist/support"
	"context"
	"errors"
	"strings"
)

func (m *adminModule) ListSupportTickets(ctx context.Context, query adminDto.ListSupportTicketsQuery) (*adminDto.ListSupportTicketsResponse, error) {
	query = query.Normalize()
	rows, total, err := m.support.ListAdminTickets(ctx, query)
	if err != nil {
		return nil, err
	}

	stats, err := m.support.GetAdminTicketStats(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]adminDto.SupportTicketListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapAdminSupportListItem(row))
	}

	return &adminDto.ListSupportTicketsResponse{
		Items: items,
		Stats: *stats,
		Pagination: adminDto.PaginationResponse{
			Page:  query.Page,
			Limit: query.Limit,
			Total: total,
		},
	}, nil
}

func (m *adminModule) GetSupportTicket(ctx context.Context, ticketID string) (*adminDto.SupportTicketDetailResponse, error) {
	row, err := m.support.GetAdminTicket(ctx, ticketID)
	if err != nil {
		return nil, mapSupportError(err)
	}
	return mapAdminSupportDetail(*row), nil
}

func (m *adminModule) UpdateSupportTicketStatus(ctx context.Context, adminUserID, ticketID string, req adminDto.UpdateSupportTicketStatusRequest, meta authDto.SessionMeta) (*adminDto.UpdateSupportTicketStatusResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, supportDto.MapValidationError(err)
	}
	status, _ := supportDto.ValidateStatus(req.Status)

	if err := m.support.UpdateTicketStatus(ctx, ticketID, status, adminUserID); err != nil {
		return nil, mapSupportError(err)
	}

	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID: &adminUserID,
		Action:      string(model.StatusChangeAction),
		EntityType:  "support_ticket",
		EntityID:    &ticketID,
		Description: "support ticket status updated to " + status,
		IPAddress:   meta.IPAddress,
		UserAgent:   meta.UserAgent,
	})

	return &adminDto.UpdateSupportTicketStatusResponse{
		ID:     ticketID,
		Status: status,
	}, nil
}

func (m *adminModule) ReplySupportTicket(ctx context.Context, adminUserID, ticketID string, req adminDto.ReplySupportTicketRequest, meta authDto.SessionMeta) (*adminDto.ReplySupportTicketResponse, error) {
	var adminReply, internalNote *string
	if v := strings.TrimSpace(req.AdminReply); v != "" {
		adminReply = &v
	}
	if v := strings.TrimSpace(req.InternalNote); v != "" {
		internalNote = &v
	}
	if adminReply == nil && internalNote == nil {
		return nil, customErr.ErrBadRequest
	}

	if err := m.support.SaveTicketReply(ctx, ticketID, adminUserID, adminReply, internalNote); err != nil {
		return nil, mapSupportError(err)
	}

	desc := "support ticket reply saved"
	if adminReply != nil {
		desc = "support ticket public reply saved"
	}
	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID: &adminUserID,
		Action:      string(model.AdminAction),
		EntityType:  "support_ticket",
		EntityID:    &ticketID,
		Description: desc,
		IPAddress:   meta.IPAddress,
		UserAgent:   meta.UserAgent,
	})

	return &adminDto.ReplySupportTicketResponse{
		ID:         ticketID,
		AdminReply: adminReply,
	}, nil
}

func mapAdminSupportListItem(row adminDto.SupportTicketRow) adminDto.SupportTicketListItem {
	preview := row.Message
	if len(preview) > 100 {
		preview = preview[:100] + "…"
	}
	return adminDto.SupportTicketListItem{
		ID:                   row.ID,
		Subject:              row.Subject,
		MessagePreview:       preview,
		IssueType:            row.IssueType,
		IssueTypeLabel:       supportDto.IssueTypeLabel(row.IssueType),
		Priority:             row.Priority,
		Status:               row.Status,
		SenderUserID:         row.SenderUserID,
		SenderName:           row.SenderName,
		SenderEmail:          row.SenderEmail,
		SenderRole:           row.SenderRole,
		RelatedCampaignTitle: row.RelatedCampaignTitle,
		AdminReply:           row.ResolutionNote,
		CreatedAt:            row.CreatedAt,
	}
}

func mapAdminSupportDetail(row adminDto.SupportTicketRow) *adminDto.SupportTicketDetailResponse {
	return &adminDto.SupportTicketDetailResponse{
		ID:                   row.ID,
		Subject:              row.Subject,
		Message:              row.Message,
		IssueType:            row.IssueType,
		IssueTypeLabel:       supportDto.IssueTypeLabel(row.IssueType),
		Priority:             row.Priority,
		Status:               row.Status,
		SenderUserID:         row.SenderUserID,
		SenderName:           row.SenderName,
		SenderEmail:          row.SenderEmail,
		SenderRole:           row.SenderRole,
		RelatedCampaignID:    row.RelatedCampaignID,
		RelatedCampaignTitle: row.RelatedCampaignTitle,
		AdminReply:           row.ResolutionNote,
		InternalNote:         row.InternalNote,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
	}
}

func mapSupportError(err error) error {
	if errors.Is(err, supportPersist.ErrTicketNotFound) {
		return customErr.ErrSupportTicketNotFound
	}
	return err
}
