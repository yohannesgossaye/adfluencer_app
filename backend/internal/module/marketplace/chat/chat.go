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
	if req.RecipientUserID == "" {
		return nil, errors.New("recipient_user_id is required")
	}

	if req.RecipientUserID == currentUserID {
		return nil, errors.New("cannot start conversation with yourself")
	}

	// Verify recipient exists
	if _, err := s.authRepo.GetUserByID(ctx, req.RecipientUserID); err != nil {
		return nil, errors.New("recipient user not found")
	}

	return s.chatRepo.GetOrCreateConversation(ctx, currentUserID, req.RecipientUserID, req.CampaignID)
}

func (s *chatService) ListConversations(ctx context.Context, currentUserID string, limit, offset int) ([]dto.ConversationItem, error) {
	return s.chatRepo.ListConversationsByUserID(ctx, currentUserID, limit, offset)
}

func (s *chatService) GetConversation(ctx context.Context, currentUserID, conversationID string) (*dto.ConversationItem, error) {
	if conversationID == "" {
		return nil, errors.New("conversation_id is required")
	}

	isParticipant, err := s.chatRepo.IsConversationParticipant(ctx, conversationID, currentUserID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		return nil, errors.New("unauthorized: you are not a participant in this conversation")
	}

	return s.chatRepo.GetConversationSummary(ctx, conversationID, currentUserID)
}

func (s *chatService) SendMessage(ctx context.Context, currentUserID string, req dto.SendMessageRequest) (*dto.Message, error) {
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
		return nil, errors.New("unauthorized: sender is not a participant in this conversation")
	}

	// Auto-detect receiver if not provided or verify if provided
	otherUserID, err := s.chatRepo.GetOtherParticipantID(ctx, req.ConversationID, currentUserID)
	if err != nil {
		return nil, err
	}
	req.ReceiverUserID = otherUserID

	msg, err := s.chatRepo.InsertMessage(ctx, req, currentUserID)
	if err != nil {
		return nil, err
	}

	// Broadcast realtime event to both participants
	outgoing := dto.WSOutgoingMessage{
		Type:    "message_created",
		Payload: msg,
	}

	s.hub.SendToUser(msg.ReceiverUserID, outgoing)
	s.hub.SendToUser(msg.SenderUserID, outgoing)

	return msg, nil
}

func (s *chatService) GetMessages(ctx context.Context, currentUserID, conversationID, cursor string, limit int) ([]dto.Message, *dto.CursorPaginationMeta, error) {
	if conversationID == "" {
		return nil, nil, errors.New("conversation_id is required")
	}

	isParticipant, err := s.chatRepo.IsConversationParticipant(ctx, conversationID, currentUserID)
	if err != nil {
		return nil, nil, err
	}
	if !isParticipant {
		return nil, nil, errors.New("unauthorized: you are not a participant in this conversation")
	}

	return s.chatRepo.GetMessages(ctx, conversationID, currentUserID, cursor, limit)
}

func (s *chatService) MarkMessagesAsRead(ctx context.Context, currentUserID, conversationID string) error {
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
	if messageID == "" {
		return nil, errors.New("message_id is required")
	}

	if deleteType != "me" && deleteType != "everyone" {
		deleteType = "me"
	}

	deletedMsg, err := s.chatRepo.DeleteMessage(ctx, messageID, currentUserID, deleteType)
	if err != nil {
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

	return deletedMsg, nil
}
