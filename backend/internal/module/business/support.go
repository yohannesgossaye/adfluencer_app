package business

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	supportDto "adinfluence_service_app/internal/constant/model/dto/support"
	supportPersist "adinfluence_service_app/internal/storage/persist/support"
	"adinfluence_service_app/platform/utils"
	"context"
	"errors"
	"strings"

	"go.uber.org/zap"
)

func (m *businessModule) ListSupportCampaignOptions(ctx context.Context, userID string) (*businessDto.ListSupportCampaignOptionsResponse, error) {
	m.logger.Info(ctx, "listing campaign options for support ticket", zap.String("business_id", userID))

	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}

	items, err := m.support.ListBusinessCampaignOptions(ctx, userID)
	if err != nil {
		m.logger.Error(ctx, "failed to query business campaign options", zap.String("business_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}
	if items == nil {
		items = []businessDto.SupportCampaignOption{}
	}

	return &businessDto.ListSupportCampaignOptionsResponse{Items: items}, nil
}

func (m *businessModule) CreateSupportTicket(ctx context.Context, userID string, req businessDto.CreateSupportTicketRequest, meta authDto.SessionMeta) (*commonDto.CreatedIDResponse, error) {
	m.logger.Info(ctx, "creating business support ticket", zap.String("business_id", userID), zap.String("subject", req.Subject), zap.String("issue_type", req.IssueType))

	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}
	if err := req.Validate(); err != nil {
		m.logger.Warn(ctx, "support ticket validation failed", zap.String("reason", err.Error()))
		return nil, supportDto.MapValidationError(err)
	}

	priority, _ := supportDto.NormalizePriority(req.Priority)
	issueType, _ := supportDto.ValidateBusinessIssueType(req.IssueType)

	var relatedCampaignID *string
	if req.RelatedCampaignID != nil && strings.TrimSpace(*req.RelatedCampaignID) != "" {
		campaignID := strings.TrimSpace(*req.RelatedCampaignID)
		ok, err := m.support.BusinessOwnsCampaign(ctx, userID, campaignID)
		if err != nil {
			return nil, err
		}
		if !ok {
			m.logger.Warn(ctx, "business does not own related campaign for support ticket", zap.String("campaign_id", campaignID))
			return nil, customErr.ErrCampaignNotFound
		}
		relatedCampaignID = &campaignID
	}

	ticketID := utils.GenerateUUID()
	if err := m.support.CreateTicket(ctx, businessDto.CreateSupportTicketInput{
		ID:                ticketID,
		UserID:            userID,
		IssueType:         issueType,
		Subject:           strings.TrimSpace(req.Subject),
		Message:           strings.TrimSpace(req.Message),
		Priority:          priority,
		RelatedCampaignID: relatedCampaignID,
	}); err != nil {
		m.logger.Error(ctx, "failed to persist support ticket", zap.String("ticket_id", ticketID), zap.String("reason", err.Error()))
		return nil, err
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.AdminAction),
		EntityType:   "support_ticket",
		EntityID:     &ticketID,
		Description:  "business support ticket created",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	m.logger.Info(ctx, "support ticket created successfully", zap.String("ticket_id", ticketID), zap.String("business_id", userID))

	return &commonDto.CreatedIDResponse{ID: ticketID}, nil
}

func (m *businessModule) ListSupportTickets(ctx context.Context, userID string) (*businessDto.ListSupportTicketsResponse, error) {
	m.logger.Info(ctx, "listing business support tickets", zap.String("business_id", userID))

	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}

	rows, err := m.support.ListTicketsByUser(ctx, userID)
	if err != nil {
		m.logger.Error(ctx, "failed to list support tickets from database", zap.String("business_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	items := make([]businessDto.SupportTicketListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapBusinessSupportListItem(row))
	}

	return &businessDto.ListSupportTicketsResponse{Items: items}, nil
}

func (m *businessModule) GetSupportTicket(ctx context.Context, userID, ticketID string) (*businessDto.SupportTicketDetailResponse, error) {
	m.logger.Info(ctx, "fetching business support ticket details", zap.String("business_id", userID), zap.String("ticket_id", ticketID))

	if err := m.requireActiveBusiness(ctx, userID); err != nil {
		return nil, err
	}

	row, err := m.support.GetTicketByUser(ctx, userID, ticketID)
	if err != nil {
		m.logger.Warn(ctx, "failed to get support ticket", zap.String("ticket_id", ticketID), zap.String("reason", err.Error()))
		return nil, mapSupportError(err)
	}

	item := mapBusinessSupportListItem(*row)
	return &businessDto.SupportTicketDetailResponse{
		SupportTicketListItem: item,
		Message:               row.Message,
	}, nil
}

func mapBusinessSupportListItem(row businessDto.SupportTicketRecord) businessDto.SupportTicketListItem {
	return businessDto.SupportTicketListItem{
		ID:                   row.ID,
		Subject:              row.Subject,
		IssueType:            row.IssueType,
		IssueTypeLabel:       supportDto.IssueTypeLabel(row.IssueType),
		Priority:             row.Priority,
		Status:               row.Status,
		AdminReply:           row.ResolutionNote,
		RelatedCampaignID:    row.RelatedCampaignID,
		RelatedCampaignTitle: row.RelatedCampaignTitle,
		CreatedAt:            row.CreatedAt,
	}
}

func mapSupportError(err error) error {
	if errors.Is(err, supportPersist.ErrTicketNotFound) {
		return customErr.ErrSupportTicketNotFound
	}
	return err
}
