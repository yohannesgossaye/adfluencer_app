package storage

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	"context"
)

type SupportInviteOption struct {
	ID            string
	CampaignTitle string
	BusinessName  string
}

type SupportPersistence interface {
	ListBusinessCampaignOptions(ctx context.Context, businessUserID string) ([]businessDto.SupportCampaignOption, error)
	CreateTicket(ctx context.Context, input businessDto.CreateSupportTicketInput) error
	ListTicketsByUser(ctx context.Context, userID string) ([]businessDto.SupportTicketRecord, error)
	GetTicketByUser(ctx context.Context, userID, ticketID string) (*businessDto.SupportTicketRecord, error)
	BusinessOwnsCampaign(ctx context.Context, businessUserID, campaignID string) (bool, error)

	ListInfluencerInviteOptions(ctx context.Context, influencerUserID string) ([]SupportInviteOption, error)
	GetInviteCampaignID(ctx context.Context, influencerUserID, inviteID string) (string, error)

	ListAdminTickets(ctx context.Context, query adminDto.ListSupportTicketsQuery) ([]adminDto.SupportTicketRow, int, error)
	GetAdminTicketStats(ctx context.Context) (*adminDto.SupportTicketStats, error)
	GetAdminTicket(ctx context.Context, ticketID string) (*adminDto.SupportTicketRow, error)
	UpdateTicketStatus(ctx context.Context, ticketID, status, adminUserID string) error
	SaveTicketReply(ctx context.Context, ticketID, adminUserID string, adminReply, internalNote *string) error
}
