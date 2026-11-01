package storage

import (
	dto "adinfluence_service_app/internal/constant/model/dto/marketplace"
	"context"
)

type ChatPersistence interface {
	GetOrCreateConversation(ctx context.Context, userOneID, userTwoID string, campaignID *string) (*dto.ConversationItem, bool, error)
	GetCampaignOwner(ctx context.Context, campaignID string) (string, error)
	GetConversationSummary(ctx context.Context, conversationID, currentUserID string) (*dto.ConversationItem, error)
	ListConversationsByUserID(ctx context.Context, userID string, limit, offset int) ([]dto.ConversationItem, error)
	IsConversationParticipant(ctx context.Context, conversationID, userID string) (bool, error)
	GetOtherParticipantID(ctx context.Context, conversationID, currentUserID string) (string, error)
	InsertMessage(ctx context.Context, req dto.SendMessageRequest, senderUserID string) (*dto.Message, bool, error)
	GetMessages(ctx context.Context, conversationID, currentUserID, cursor string, limit int) ([]dto.Message, *dto.CursorPaginationMeta, error)
	MarkMessagesAsRead(ctx context.Context, conversationID, currentUserID string) (int, error)
	DeleteMessage(ctx context.Context, messageID, currentUserID, deleteType string) (*dto.Message, error)
	GetMessageByID(ctx context.Context, messageID string) (*dto.Message, error)
}
