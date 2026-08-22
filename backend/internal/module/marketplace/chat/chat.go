package chat

import (
	dto "adinfluence_service_app/internal/constant/model/dto/marketplace"
	"adinfluence_service_app/internal/module"
	"adinfluence_service_app/internal/module/marketplace/chat/realtime"
	"adinfluence_service_app/internal/storage"
	"adinfluence_service_app/platform/utils"
	"context"
	"errors"
	"strings"
	"time"

	"go.uber.org/zap"
)

type chatService struct {
	chatRepo storage.ChatPersistence
	authRepo storage.AuthPersisitence
	hub      *realtime.Hub
	log      utils.Logger
}

func InitChatService(chatRepo storage.ChatPersistence, authRepo storage.AuthPersisitence, log utils.Logger) module.ChatService {
	return &chatService{
		chatRepo: chatRepo,
		authRepo: authRepo,
		hub:      realtime.NewHub(),
		log:      log,
	}
}

func (s *chatService) GetHub() *realtime.Hub {
	return s.hub
}

func (s *chatService) GetOrCreateConversation(ctx context.Context, currentUserID string, req dto.CreateConversationRequest) (*dto.ConversationItem, error) {
	s.log.Info(ctx, "chat get or create conversation", zap.String("current_user_id", currentUserID), zap.String("recipient_user_id", req.RecipientUserID))

	if req.RecipientUserID == "" {
		return nil, errors.New("recipient_user_id is required")
	}

	if req.RecipientUserID == currentUserID {
		return nil, errors.New("cannot start conversation with yourself")
	}

	// Verify recipient exists
	if _, err := s.authRepo.GetUserByID(ctx, req.RecipientUserID); err != nil {
		s.log.Warn(ctx, "recipient user not found for conversation", zap.String("recipient_user_id", req.RecipientUserID))
		return nil, errors.New("recipient user not found")
	}

	res, err := s.chatRepo.GetOrCreateConversation(ctx, currentUserID, req.RecipientUserID, req.CampaignID)
	if err != nil {
		s.log.Error(ctx, "failed to get or create conversation in persistence", zap.String("reason", err.Error()))
		return nil, err
	}

	s.log.Info(ctx, "conversation ready", zap.String("conversation_id", res.ID))
	return res, nil
}

func (s *chatService) ListConversations(ctx context.Context, currentUserID string, limit, offset int) ([]dto.ConversationItem, error) {
	s.log.Info(ctx, "listing chat conversations", zap.String("current_user_id", currentUserID))
	items, err := s.chatRepo.ListConversationsByUserID(ctx, currentUserID, limit, offset)
	if err != nil {
		s.log.Error(ctx, "failed to list conversations from database", zap.String("current_user_id", currentUserID), zap.String("reason", err.Error()))
		return nil, err
	}
	return items, nil
}

func (s *chatService) GetConversation(ctx context.Context, currentUserID, conversationID string) (*dto.ConversationItem, error) {
	s.log.Info(ctx, "fetching chat conversation details", zap.String("current_user_id", currentUserID), zap.String("conversation_id", conversationID))

	if conversationID == "" {
		return nil, errors.New("conversation_id is required")
	}

	isParticipant, err := s.chatRepo.IsConversationParticipant(ctx, conversationID, currentUserID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		s.log.Warn(ctx, "unauthorized conversation access attempt", zap.String("current_user_id", currentUserID), zap.String("conversation_id", conversationID))
		return nil, errors.New("unauthorized: you are not a participant in this conversation")
	}

	return s.chatRepo.GetConversationSummary(ctx, conversationID, currentUserID)
}

func (s *chatService) SendMessage(ctx context.Context, currentUserID string, req dto.SendMessageRequest) (*dto.Message, error) {
	s.log.Info(ctx, "sending chat message", zap.String("sender_user_id", currentUserID), zap.String("conversation_id", req.ConversationID))

	req.Body = strings.TrimSpace(req.Body)

	if req.ConversationID == "" {
		return nil, errors.New("conversation_id is required")
	}

	if req.Body == "" && (req.MediaURL == nil || *req.MediaURL == "") {
		return nil, errors.New("message body or attachment is required")
	}

	// Verify current user is participant
	isSenderParticipant, err := s.chatRepo.IsConversationParticipant(ctx, req.ConversationID, currentUserID)
	if err != nil {
		return nil, err
	}
	if !isSenderParticipant {
		s.log.Warn(ctx, "sender not participant in conversation", zap.String("sender_id", currentUserID), zap.String("conversation_id", req.ConversationID))
		return nil, errors.New("unauthorized: sender is not a participant in this conversation")
	}

	// Auto-detect receiver if not provided or verify if provided
	otherUserID, err := s.chatRepo.GetOtherParticipantID(ctx, req.ConversationID, currentUserID)
	if err != nil {
		s.log.Error(ctx, "failed to get conversation recipient", zap.String("conversation_id", req.ConversationID), zap.String("reason", err.Error()))
		return nil, err
	}
	req.ReceiverUserID = otherUserID

	msg, err := s.chatRepo.InsertMessage(ctx, req, currentUserID)
	if err != nil {
		s.log.Error(ctx, "failed to insert message into persistence", zap.String("conversation_id", req.ConversationID), zap.String("reason", err.Error()))
		return nil, err
	}

	// Broadcast realtime event to both participants
	outgoing := dto.WSOutgoingMessage{
		Type:    "message_created",
		Payload: msg,
	}

	s.hub.SendToUser(msg.ReceiverUserID, outgoing)
	s.hub.SendToUser(msg.SenderUserID, outgoing)

	s.log.Info(ctx, "chat message sent & broadcasted successfully", zap.String("message_id", msg.ID), zap.String("conversation_id", msg.ConversationID))

	return msg, nil
}

func (s *chatService) GetMessages(ctx context.Context, currentUserID, conversationID, cursor string, limit int) ([]dto.Message, *dto.CursorPaginationMeta, error) {
	s.log.Info(ctx, "fetching chat messages", zap.String("current_user_id", currentUserID), zap.String("conversation_id", conversationID))

	if conversationID == "" {
		return nil, nil, errors.New("conversation_id is required")
	}

	isParticipant, err := s.chatRepo.IsConversationParticipant(ctx, conversationID, currentUserID)
	if err != nil {
		return nil, nil, err
	}
	if !isParticipant {
		s.log.Warn(ctx, "unauthorized access to chat message history", zap.String("current_user_id", currentUserID), zap.String("conversation_id", conversationID))
		return nil, nil, errors.New("unauthorized: you are not a participant in this conversation")
	}

	return s.chatRepo.GetMessages(ctx, conversationID, currentUserID, cursor, limit)
}

func (s *chatService) MarkMessagesAsRead(ctx context.Context, currentUserID, conversationID string) error {
	s.log.Info(ctx, "marking conversation messages as read", zap.String("current_user_id", currentUserID), zap.String("conversation_id", conversationID))

	if conversationID == "" {
		return errors.New("conversation_id is required")
	}

	isParticipant, err := s.chatRepo.IsConversationParticipant(ctx, conversationID, currentUserID)
	if err != nil {
		return err
	}
	if !isParticipant {
		return errors.New("unauthorized: you are not a participant in this conversation")
	}

	if err := s.chatRepo.MarkMessagesAsRead(ctx, conversationID, currentUserID); err != nil {
		s.log.Error(ctx, "failed to mark messages read in database", zap.String("conversation_id", conversationID), zap.String("reason", err.Error()))
		return err
	}

	// Notify other participant that messages were read
	otherUserID, err := s.chatRepo.GetOtherParticipantID(ctx, conversationID, currentUserID)
	if err == nil {
		event := dto.WSOutgoingMessage{
			Type: "messages_read",
			Payload: dto.WSReadEvent{
				ConversationID: conversationID,
				ReadByUserID:   currentUserID,
				ReadAt:         time.Now().UTC(),
			},
		}
		s.hub.SendToUser(otherUserID, event)
	}

	return nil
}

func (s *chatService) DeleteMessage(ctx context.Context, currentUserID, messageID, deleteType string) (*dto.Message, error) {
	s.log.Info(ctx, "deleting message", zap.String("current_user_id", currentUserID), zap.String("message_id", messageID), zap.String("delete_type", deleteType))

	if messageID == "" {
		return nil, errors.New("message_id is required")
	}

	if deleteType != "me" && deleteType != "everyone" {
		deleteType = "me"
	}

	deletedMsg, err := s.chatRepo.DeleteMessage(ctx, messageID, currentUserID, deleteType)
	if err != nil {
		s.log.Error(ctx, "failed to delete message from database", zap.String("message_id", messageID), zap.String("reason", err.Error()))
		return nil, err
	}

	// If deleted for everyone, notify both users over WebSocket
	if deleteType == "everyone" {
		event := dto.WSOutgoingMessage{
			Type: "message_deleted",
			Payload: dto.WSDeleteEvent{
				MessageID:      messageID,
				ConversationID: deletedMsg.ConversationID,
				DeleteType:     "everyone",
				Tombstone:      "This message was deleted",
			},
		}
		s.hub.SendToUser(deletedMsg.SenderUserID, event)
		s.hub.SendToUser(deletedMsg.ReceiverUserID, event)
	} else {
		// Deleted for me only: notify current user's other sessions
		event := dto.WSOutgoingMessage{
			Type: "message_deleted",
			Payload: dto.WSDeleteEvent{
				MessageID:      messageID,
				ConversationID: deletedMsg.ConversationID,
				DeleteType:     "me",
			},
		}
		s.hub.SendToUser(currentUserID, event)
	}

	s.log.Info(ctx, "message deleted successfully", zap.String("message_id", messageID), zap.String("delete_type", deleteType))

	return deletedMsg, nil
}
