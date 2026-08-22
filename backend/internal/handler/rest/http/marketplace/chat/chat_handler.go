package chat

import (
	dto "adinfluence_service_app/internal/constant/model/dto/marketplace"
	"adinfluence_service_app/internal/constant/model/response"
	"adinfluence_service_app/internal/handler"
	"adinfluence_service_app/internal/handler/middleware"
	"adinfluence_service_app/internal/module"
	"adinfluence_service_app/internal/module/marketplace/chat/realtime"
	"adinfluence_service_app/platform/jwt"
	"adinfluence_service_app/platform/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/go-chi/chi/v5"
)

type chatHandler struct {
	service module.ChatService
	tokens  *jwt.TokenManager
	logger  utils.Logger
}

func InitChatHandler(service module.ChatService, tokens *jwt.TokenManager, logger utils.Logger) handler.ChatHandler {
	return &chatHandler{
		service: service,
		tokens:  tokens,
		logger:  logger,
	}
}

func (h *chatHandler) requireUser(w http.ResponseWriter, r *http.Request) (*jwt.AccessClaims, bool) {
	claims, ok := middleware.GetAccessClaims(r.Context())
	if !ok || claims == nil {
		response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", nil, "authentication required")
		return nil, false
	}
	return claims, true
}

func (h *chatHandler) GetOrCreateConversation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}

	var req dto.CreateConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_body", err, "invalid request body")
		return
	}

	conv, err := h.service.GetOrCreateConversation(ctx, user.UserID, req)
	if err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "failed_create_conversation", err, err.Error())
		return
	}

	response.SendSuccessResponse(w, r, http.StatusCreated, "conversation ready", conv, nil)
}

func (h *chatHandler) ListConversations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	conversations, err := h.service.ListConversations(ctx, user.UserID, limit, offset)
	if err != nil {
		response.SendErrorResponse(w, http.StatusInternalServerError, "failed_list_conversations", err, err.Error())
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "conversations fetched successfully", conversations, nil)
}

func (h *chatHandler) GetConversation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}

	convID := chi.URLParam(r, "id")
	if convID == "" {
		convID = r.URL.Query().Get("conversation_id")
	}

	conv, err := h.service.GetConversation(ctx, user.UserID, convID)
	if err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "failed_get_conversation", err, err.Error())
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "conversation fetched successfully", conv, nil)
}

func (h *chatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}

	var req dto.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_body", err, "invalid request body")
		return
	}

	msg, err := h.service.SendMessage(ctx, user.UserID, req)
	if err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "failed_send_message", err, err.Error())
		return
	}

	response.SendSuccessResponse(w, r, http.StatusCreated, "message sent successfully", msg, nil)
}

func (h *chatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}

	conversationID := r.URL.Query().Get("conversation_id")
	cursor := r.URL.Query().Get("cursor")
	if cursor == "" {
		cursor = r.URL.Query().Get("before")
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	messages, meta, err := h.service.GetMessages(ctx, user.UserID, conversationID, cursor, limit)
	if err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "failed_get_messages", err, err.Error())
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "messages fetched successfully", messages, meta)
}

func (h *chatHandler) MarkMessagesAsRead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}

	var req dto.MarkReadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_body", err, "invalid request body")
		return
	}

	if err := h.service.MarkMessagesAsRead(ctx, user.UserID, req.ConversationID); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "failed_mark_read", err, err.Error())
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "messages marked as read", nil, nil)
}

func (h *chatHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}

	messageID := chi.URLParam(r, "id")
	deleteType := r.URL.Query().Get("delete_type")
	if deleteType == "" {
		var req dto.DeleteMessageRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		deleteType = req.DeleteType
	}

	msg, err := h.service.DeleteMessage(ctx, user.UserID, messageID, deleteType)
	if err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "failed_delete_message", err, err.Error())
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "message deleted successfully", msg, nil)
}

func (h *chatHandler) WebSocket(w http.ResponseWriter, r *http.Request) {
	// Extract and validate token from query parameter or header
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if tokenString == "" {
		http.Error(w, "missing access token", http.StatusUnauthorized)
		return
	}

	claims, err := h.tokens.ParseAccessToken(tokenString)
	if err != nil || claims == nil || claims.UserID == "" {
		http.Error(w, "invalid or expired token", http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return
	}

	hub := h.service.GetHub()
	client := realtime.NewClient(claims.UserID, conn, hub)
	hub.Register(client)
	defer hub.Unregister(client)

	ctx := r.Context()
	go client.WriteLoop(ctx)

	for {
		var incoming dto.WSIncomingMessage
		err := wsjson.Read(ctx, conn, &incoming)
		if err != nil {
			return
		}

		switch incoming.Type {
		case "send_message":
			_, err := h.service.SendMessage(ctx, claims.UserID, dto.SendMessageRequest{
				ConversationID: incoming.ConversationID,
				ReceiverUserID: incoming.ReceiverUserID,
				Body:           incoming.Body,
				MediaURL:       incoming.MediaURL,
				MediaType:      incoming.MediaType,
			})
			if err != nil {
				_ = wsjson.Write(ctx, conn, dto.WSOutgoingMessage{
					Type:    "error",
					Payload: err.Error(),
				})
			}

		case "mark_read":
			_ = h.service.MarkMessagesAsRead(ctx, claims.UserID, incoming.ConversationID)

		case "delete_message":
			_, err := h.service.DeleteMessage(ctx, claims.UserID, incoming.MessageID, incoming.DeleteType)
			if err != nil {
				_ = wsjson.Write(ctx, conn, dto.WSOutgoingMessage{
					Type:    "error",
					Payload: err.Error(),
				})
			}
		}
	}
}
