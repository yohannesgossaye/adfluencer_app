package chat

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	dto "adinfluence_service_app/internal/constant/model/dto/marketplace"
	"adinfluence_service_app/internal/module"
	"adinfluence_service_app/internal/module/marketplace/chat/realtime"
	"adinfluence_service_app/internal/storage"
	"adinfluence_service_app/platform/utils"
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
)

type chatService struct {
	chatRepo storage.ChatPersistence
	authRepo storage.AuthPersisitence
	hub      *realtime.Hub
	limiter  *sendLimiter
	log      utils.Logger
}

func InitChatService(chatRepo storage.ChatPersistence, authRepo storage.AuthPersisitence, log utils.Logger) module.ChatService {
	return &chatService{
		chatRepo: chatRepo,
		authRepo: authRepo,
		hub:      realtime.NewHub(),
		limiter:  newSendLimiter(),
		log:      log,
	}
}

func (s *chatService) GetHub() *realtime.Hub {
	return s.hub
}

// loadUser returns the account, or ErrAccountCannotChat when it no longer exists.
func (s *chatService) loadUser(ctx context.Context, userID string) (*authDto.UserResponse, error) {
	user, err := s.authRepo.GetUserByID(ctx, userID)
	if errors.Is(err, customErr.ErrUserNotFound) {
		return nil, dto.ErrAccountCannotChat
	}
	return user, err
}

// CheckCanConnect is used when opening and while holding a realtime
// connection: blocked, removed or suspended accounts are disconnected.
func (s *chatService) CheckCanConnect(ctx context.Context, userID string) error {
	user, err := s.loadUser(ctx, userID)
	if err != nil {
		return err
	}
	if !canUseChat(user) {
		return dto.ErrAccountCannotChat
	}
	return nil
}

func (s *chatService) requireReader(ctx context.Context, userID string) error {
	return s.CheckCanConnect(ctx, userID)
}

func (s *chatService) requireSender(ctx context.Context, userID string) (*authDto.UserResponse, error) {
	user, err := s.loadUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !canStartOrSend(user) {
		return nil, dto.ErrAccountCannotChat
	}
	return user, nil
}

func (s *chatService) requireParticipant(ctx context.Context, conversationID, userID string) error {
	if err := requireUUID(conversationID, "conversation_id"); err != nil {
		return err
	}
	ok, err := s.chatRepo.IsConversationParticipant(ctx, conversationID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return dto.ErrNotParticipant
	}
	return nil
}

func (s *chatService) GetOrCreateConversation(ctx context.Context, currentUserID string, req dto.CreateConversationRequest) (*dto.ConversationItem, bool, error) {
	if err := requireUUID(req.RecipientUserID, "recipient_user_id"); err != nil {
		return nil, false, err
	}
	if req.RecipientUserID == currentUserID {
		return nil, false, dto.Invalid("you cannot start a conversation with yourself")
	}
	if req.CampaignID != nil && *req.CampaignID == "" {
		req.CampaignID = nil
	}
	if req.CampaignID != nil {
		if err := requireUUID(*req.CampaignID, "campaign_id"); err != nil {
			return nil, false, err
		}
	}

	me, err := s.requireSender(ctx, currentUserID)
	if err != nil {
		return nil, false, err
	}
	other, err := s.authRepo.GetUserByID(ctx, req.RecipientUserID)
	if errors.Is(err, customErr.ErrUserNotFound) {
		return nil, false, dto.ErrRecipientNotFound
	}
	if err != nil {
		return nil, false, err
	}
	if !isBusinessInfluencerPair(me, other) {
		return nil, false, dto.ErrChatNotAllowed
	}
	if !canStartOrSend(other) {
		return nil, false, dto.ErrRecipientUnavailable
	}

	if req.CampaignID != nil {
		ownerID, err := s.chatRepo.GetCampaignOwner(ctx, *req.CampaignID)
		if err != nil {
			return nil, false, err
		}
		businessID := me.ID
		if other.Role == string(model.BusinessRole) {
			businessID = other.ID
		}
		if ownerID != businessID {
			return nil, false, dto.ErrCampaignNotOwned
		}
	}

	conv, created, err := s.chatRepo.GetOrCreateConversation(ctx, currentUserID, req.RecipientUserID, req.CampaignID)
	if err != nil {
		s.log.Error(ctx, "failed to get or create conversation", zap.String("reason", err.Error()))
		return nil, false, err
	}
	return conv, created, nil
}

func (s *chatService) ListConversations(ctx context.Context, currentUserID string, limit, offset int) ([]dto.ConversationItem, error) {
	if err := s.requireReader(ctx, currentUserID); err != nil {
		return nil, err
	}
	return s.chatRepo.ListConversationsByUserID(ctx, currentUserID, limit, offset)
}

func (s *chatService) GetConversation(ctx context.Context, currentUserID, conversationID string) (*dto.ConversationItem, error) {
	if err := s.requireReader(ctx, currentUserID); err != nil {
		return nil, err
	}
	if err := s.requireParticipant(ctx, conversationID, currentUserID); err != nil {
		return nil, err
	}
	return s.chatRepo.GetConversationSummary(ctx, conversationID, currentUserID)
}

func (s *chatService) SendMessage(ctx context.Context, currentUserID string, req dto.SendMessageRequest) (*dto.Message, error) {
	if err := validateMessage(&req); err != nil {
		return nil, err
	}
	if _, err := s.requireSender(ctx, currentUserID); err != nil {
		return nil, err
	}
	if err := s.requireParticipant(ctx, req.ConversationID, currentUserID); err != nil {
		return nil, err
	}
	if !s.limiter.allow(currentUserID, time.Now()) {
		return nil, dto.ErrSendRateLimited
	}

	otherUserID, err := s.chatRepo.GetOtherParticipantID(ctx, req.ConversationID, currentUserID)
	if err != nil {
		return nil, err
	}
	other, err := s.loadUser(ctx, otherUserID)
	if err != nil || !canUseChat(other) {
		return nil, dto.ErrRecipientUnavailable
	}
	req.ReceiverUserID = otherUserID

	msg, created, err := s.chatRepo.InsertMessage(ctx, req, currentUserID)
	if err != nil {
		s.log.Error(ctx, "failed to insert message", zap.String("conversation_id", req.ConversationID), zap.String("reason", err.Error()))
		return nil, err
	}

	if created {
		outgoing := dto.WSOutgoingMessage{Type: "message_created", Payload: msg}
		s.hub.SendToUser(msg.ReceiverUserID, outgoing)
		s.hub.SendToUser(msg.SenderUserID, outgoing)
	}
	return msg, nil
}

func (s *chatService) GetMessages(ctx context.Context, currentUserID, conversationID, cursor string, limit int) ([]dto.Message, *dto.CursorPaginationMeta, error) {
	if err := s.requireReader(ctx, currentUserID); err != nil {
		return nil, nil, err
	}
	if err := s.requireParticipant(ctx, conversationID, currentUserID); err != nil {
		return nil, nil, err
	}
	if cursor != "" {
		if err := requireUUID(cursor, "cursor"); err != nil {
			return nil, nil, err
		}
	}
	return s.chatRepo.GetMessages(ctx, conversationID, currentUserID, cursor, limit)
}

func (s *chatService) MarkMessagesAsRead(ctx context.Context, currentUserID, conversationID string) (*dto.MarkReadResponse, error) {
	if err := s.requireReader(ctx, currentUserID); err != nil {
		return nil, err
	}
	if err := s.requireParticipant(ctx, conversationID, currentUserID); err != nil {
		return nil, err
	}

	count, err := s.chatRepo.MarkMessagesAsRead(ctx, conversationID, currentUserID)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		// Tell the sender (read receipts) and the reader's other tabs/devices
		// (to clear their unread badges).
		event := dto.WSOutgoingMessage{
			Type: "messages_read",
			Payload: dto.WSReadEvent{
				ConversationID: conversationID,
				ReadByUserID:   currentUserID,
				ReadCount:      count,
				ReadAt:         time.Now().UTC(),
			},
		}
		if otherUserID, err := s.chatRepo.GetOtherParticipantID(ctx, conversationID, currentUserID); err == nil {
			s.hub.SendToUser(otherUserID, event)
		}
		s.hub.SendToUser(currentUserID, event)
	}

	return &dto.MarkReadResponse{ConversationID: conversationID, ReadCount: count}, nil
}

func (s *chatService) DeleteMessage(ctx context.Context, currentUserID, messageID, deleteType string) (*dto.Message, error) {
	if err := requireUUID(messageID, "message_id"); err != nil {
		return nil, err
	}
	if deleteType != "everyone" {
		deleteType = "me"
	}
	if err := s.requireReader(ctx, currentUserID); err != nil {
		return nil, err
	}

	deleted, err := s.chatRepo.DeleteMessage(ctx, messageID, currentUserID, deleteType)
	if err != nil {
		return nil, err
	}

	event := dto.WSOutgoingMessage{
		Type: "message_deleted",
		Payload: dto.WSDeleteEvent{
			MessageID:      messageID,
			ConversationID: deleted.ConversationID,
			DeleteType:     deleteType,
		},
	}
	if deleteType == "everyone" {
		event.Payload = dto.WSDeleteEvent{
			MessageID:      messageID,
			ConversationID: deleted.ConversationID,
			DeleteType:     deleteType,
			Tombstone:      "This message was deleted",
		}
		s.hub.SendToUser(deleted.SenderUserID, event)
		s.hub.SendToUser(deleted.ReceiverUserID, event)
	} else {
		s.hub.SendToUser(currentUserID, event)
	}

	return deleted, nil
}
