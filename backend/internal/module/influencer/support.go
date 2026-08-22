package influencer

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
	supportDto "adinfluence_service_app/internal/constant/model/dto/support"
	supportPersist "adinfluence_service_app/internal/storage/persist/support"
	"adinfluence_service_app/platform/utils"
	"context"
	"errors"
	"strings"

	"go.uber.org/zap"
)

func (m *influencerModule) ListSupportInviteOptions(ctx context.Context, userID string) (*influencerDto.ListSupportInviteOptionsResponse, error) {
	m.logger.Info(ctx, "listing influencer support invite options", zap.String("user_id", userID))

	if err := m.requireActiveInfluencer(ctx, userID); err != nil {
		return nil, err
	}

	rows, err := m.support.ListInfluencerInviteOptions(ctx, userID)
	if err != nil {
		m.logger.Error(ctx, "failed to query support invite options", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	items := make([]influencerDto.SupportInviteOption, 0, len(rows))
	for _, row := range rows {
		items = append(items, influencerDto.SupportInviteOption{
			ID:            row.ID,
			CampaignTitle: row.CampaignTitle,
			BusinessName:  row.BusinessName,
		})
	}

	return &influencerDto.ListSupportInviteOptionsResponse{Items: items}, nil
}

func (m *influencerModule) CreateSupportTicket(ctx context.Context, userID string, req influencerDto.CreateSupportTicketRequest, meta authDto.SessionMeta) (*commonDto.CreatedIDResponse, error) {
	m.logger.Info(ctx, "creating influencer support ticket", zap.String("user_id", userID), zap.String("subject", req.Subject), zap.String("issue_type", req.IssueType))

	if err := m.requireActiveInfluencer(ctx, userID); err != nil {
		return nil, err
	}
	if err := req.Validate(); err != nil {
		m.logger.Warn(ctx, "influencer support ticket validation failed", zap.String("reason", err.Error()))
		return nil, supportDto.MapValidationError(err)
	}

	priority, _ := supportDto.NormalizePriority(req.Priority)
	issueType, _ := supportDto.ValidateInfluencerIssueType(req.IssueType)

	var relatedCampaignID *string
	if req.RelatedInviteID != nil && strings.TrimSpace(*req.RelatedInviteID) != "" {
		inviteID := strings.TrimSpace(*req.RelatedInviteID)
		campaignID, err := m.support.GetInviteCampaignID(ctx, userID, inviteID)
		if err != nil {
			m.logger.Warn(ctx, "failed to get invite campaign ID for support ticket", zap.String("invite_id", inviteID), zap.String("reason", err.Error()))
			return nil, mapSupportError(err)
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
		m.logger.Error(ctx, "failed to persist influencer support ticket", zap.String("ticket_id", ticketID), zap.String("reason", err.Error()))
		return nil, err
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &userID,
		TargetUserID: &userID,
		Action:       string(model.AdminAction),
		EntityType:   "support_ticket",
		EntityID:     &ticketID,
		Description:  "influencer support ticket created",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	m.logger.Info(ctx, "influencer support ticket created successfully", zap.String("ticket_id", ticketID), zap.String("user_id", userID))

	return &commonDto.CreatedIDResponse{ID: ticketID}, nil
}

func (m *influencerModule) ListSupportTickets(ctx context.Context, userID string) (*influencerDto.ListSupportTicketsResponse, error) {
	m.logger.Info(ctx, "listing influencer support tickets", zap.String("user_id", userID))

	if err := m.requireActiveInfluencer(ctx, userID); err != nil {
		return nil, err
	}

	rows, err := m.support.ListTicketsByUser(ctx, userID)
	if err != nil {
		m.logger.Error(ctx, "failed to list influencer support tickets", zap.String("user_id", userID), zap.String("reason", err.Error()))
		return nil, err
	}

	items := make([]influencerDto.SupportTicketListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapInfluencerSupportListItem(row))
	}

	return &influencerDto.ListSupportTicketsResponse{Items: items}, nil
}

func (m *influencerModule) GetSupportTicket(ctx context.Context, userID, ticketID string) (*influencerDto.SupportTicketDetailResponse, error) {
	m.logger.Info(ctx, "fetching influencer support ticket details", zap.String("user_id", userID), zap.String("ticket_id", ticketID))

	if err := m.requireActiveInfluencer(ctx, userID); err != nil {
		return nil, err
	}

	row, err := m.support.GetTicketByUser(ctx, userID, ticketID)
	if err != nil {
		m.logger.Warn(ctx, "failed to get support ticket", zap.String("ticket_id", ticketID), zap.String("reason", err.Error()))
		return nil, mapSupportError(err)
	}

	item := mapInfluencerSupportListItem(*row)
	return &influencerDto.SupportTicketDetailResponse{
		SupportTicketListItem: item,
		Message:               row.Message,
	}, nil
}

func mapInfluencerSupportListItem(row businessDto.SupportTicketRecord) influencerDto.SupportTicketListItem {
	return influencerDto.SupportTicketListItem{
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
	if errors.Is(err, supportPersist.ErrInviteNotFound) {
		return customErr.ErrInviteNotFound
	}
	return err
}
