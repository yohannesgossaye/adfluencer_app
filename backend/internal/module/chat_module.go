package module

import (
	dto "adinfluence_service_app/internal/constant/model/dto/marketplace"
	"adinfluence_service_app/internal/module/marketplace/chat/realtime"
	"context"
)

type ChatService interface {
	// GetOrCreateConversation reports created=true when a new conversation was made.
	GetOrCreateConversation(ctx context.Context, currentUserID string, req dto.CreateConversationRequest) (*dto.ConversationItem, bool, error)
	ListConversations(ctx context.Context, currentUserID string, limit, offset int) ([]dto.ConversationItem, error)
	GetConversation(ctx context.Context, currentUserID, conversationID string) (*dto.ConversationItem, error)
	SendMessage(ctx context.Context, currentUserID string, req dto.SendMessageRequest) (*dto.Message, error)
	GetMessages(ctx context.Context, currentUserID, conversationID, cursor string, limit int) ([]dto.Message, *dto.CursorPaginationMeta, error)
	MarkMessagesAsRead(ctx context.Context, currentUserID, conversationID string) (*dto.MarkReadResponse, error)
	DeleteMessage(ctx context.Context, currentUserID, messageID, deleteType string) (*dto.Message, error)
	// CheckCanConnect returns an error when the account may not hold a realtime connection.
	CheckCanConnect(ctx context.Context, userID string) error
	GetHub() *realtime.Hub
}
