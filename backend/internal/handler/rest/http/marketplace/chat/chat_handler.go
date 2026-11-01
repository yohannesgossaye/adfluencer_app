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
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	// maxChatBodyBytes caps every chat request and WebSocket frame.
	maxChatBodyBytes = 64 << 10

	// wsSubprotocol is offered by browser clients together with the access
	// token: new WebSocket(url, ["adfluence.v1", accessToken]). Browsers
	// cannot set an Authorization header on WebSockets, and putting the token
	// in the URL would write it to access logs.
	wsSubprotocol = "adfluence.v1"

	// How often an open socket re-checks that the account may still chat.
	wsStatusRecheck = 60 * time.Second
)

type chatHandler struct {
	service        module.ChatService
	tokens         *jwt.TokenManager
	originPatterns []string
	logger         utils.Logger
}

// allowedOrigins are the same origins allowed by CORS (scheme://host,
// "*" wildcards allowed). Clients without an Origin header (mobile apps) are
// not affected.
func InitChatHandler(service module.ChatService, tokens *jwt.TokenManager, allowedOrigins []string, logger utils.Logger) handler.ChatHandler {
	return &chatHandler{
		service:        service,
		tokens:         tokens,
		originPatterns: allowedOrigins,
		logger:         logger,
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

func decodeBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxChatBodyBytes)
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		writeBodyError(w, err)
		return false
	}
	return true
}

func (h *chatHandler) GetOrCreateConversation(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}

	var req dto.CreateConversationRequest
	if !decodeBody(w, r, &req) {
		return
	}

	conv, created, err := h.service.GetOrCreateConversation(r.Context(), user.UserID, req)
	if err != nil {
		writeChatError(w, err)
		return
	}

	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}
	response.SendSuccessResponse(w, r, status, "conversation ready", conv, nil)
}

func (h *chatHandler) ListConversations(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	conversations, err := h.service.ListConversations(r.Context(), user.UserID, limit, offset)
	if err != nil {
		writeChatError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "conversations fetched successfully", conversations, nil)
}

func (h *chatHandler) GetConversation(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}

	conv, err := h.service.GetConversation(r.Context(), user.UserID, chi.URLParam(r, "id"))
	if err != nil {
		writeChatError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "conversation fetched successfully", conv, nil)
}

func (h *chatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}

	var req dto.SendMessageRequest
	if !decodeBody(w, r, &req) {
		return
	}

	msg, err := h.service.SendMessage(r.Context(), user.UserID, req)
	if err != nil {
		writeChatError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusCreated, "message sent successfully", msg, nil)
}

func (h *chatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}

	q := r.URL.Query()
	cursor := q.Get("cursor")
	if cursor == "" {
		cursor = q.Get("before")
	}
	limit, _ := strconv.Atoi(q.Get("limit"))

	messages, meta, err := h.service.GetMessages(r.Context(), user.UserID, q.Get("conversation_id"), cursor, limit)
	if err != nil {
		writeChatError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "messages fetched successfully", messages, meta)
}

func (h *chatHandler) MarkMessagesAsRead(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}

	var req dto.MarkReadRequest
	if !decodeBody(w, r, &req) {
		return
	}

	res, err := h.service.MarkMessagesAsRead(r.Context(), user.UserID, req.ConversationID)
	if err != nil {
		writeChatError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "messages marked as read", res, nil)
}

func (h *chatHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireUser(w, r)
	if !ok {
		return
	}

	deleteType := r.URL.Query().Get("delete_type")
	if deleteType == "" && r.ContentLength > 0 {
		var req dto.DeleteMessageRequest
		if !decodeBody(w, r, &req) {
			return
		}
		deleteType = req.DeleteType
	}

	msg, err := h.service.DeleteMessage(r.Context(), user.UserID, chi.URLParam(r, "id"), deleteType)
	if err != nil {
		writeChatError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "message deleted successfully", msg, nil)
}

// wsAccessToken reads the token from the Authorization header (mobile/server
// clients) or from the WebSocket subprotocol list (browsers). It reports
// whether the subprotocol form was used, so the handshake can echo it.
func wsAccessToken(r *http.Request) (token string, viaSubprotocol bool) {
	if auth := r.Header.Get("Authorization"); strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer "), false
	}
	var offered bool
	for _, p := range strings.Split(r.Header.Get("Sec-WebSocket-Protocol"), ",") {
		p = strings.TrimSpace(p)
		switch {
		case p == wsSubprotocol:
			offered = true
		case p != "":
			token = p
		}
	}
	if !offered {
		return "", false
	}
	return token, true
}

func (h *chatHandler) WebSocket(w http.ResponseWriter, r *http.Request) {
	tokenString, viaSubprotocol := wsAccessToken(r)
	if tokenString == "" {
		response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", nil, "missing access token")
		return
	}

	claims, err := h.tokens.ParseAccessToken(tokenString)
	if err != nil || claims == nil || claims.UserID == "" {
		response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", nil, "invalid or expired token")
		return
	}

	if err := h.service.CheckCanConnect(r.Context(), claims.UserID); err != nil {
		writeChatError(w, err)
		return
	}

	opts := &websocket.AcceptOptions{OriginPatterns: h.originPatterns}
	if viaSubprotocol {
		opts.Subprotocols = []string{wsSubprotocol}
	}
	conn, err := websocket.Accept(w, r, opts)
	if err != nil {
		return
	}
	conn.SetReadLimit(maxChatBodyBytes)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	hub := h.service.GetHub()
	client := realtime.NewClient(claims.UserID, conn)
	hub.Register(client)
	defer hub.Unregister(client)

	go client.WriteLoop(ctx)
	go h.watchSession(ctx, client, claims)

	for {
		var incoming dto.WSIncomingMessage
		if err := wsjson.Read(ctx, conn, &incoming); err != nil {
			return
		}
		h.handleFrame(ctx, client, claims.UserID, incoming)
	}
}

func (h *chatHandler) handleFrame(ctx context.Context, client *realtime.Client, userID string, in dto.WSIncomingMessage) {
	var err error
	switch in.Type {
	case "send_message":
		_, err = h.service.SendMessage(ctx, userID, dto.SendMessageRequest{
			ConversationID:  in.ConversationID,
			Body:            in.Body,
			MediaURL:        in.MediaURL,
			MediaType:       in.MediaType,
			ClientMessageID: in.ClientMessageID,
		})
	case "mark_read":
		_, err = h.service.MarkMessagesAsRead(ctx, userID, in.ConversationID)
	case "delete_message":
		_, err = h.service.DeleteMessage(ctx, userID, in.MessageID, in.DeleteType)
	default:
		err = dto.Invalid("unknown message type")
	}
	if err == nil {
		return
	}

	_, code, msg := classifyChatError(err)
	if code == "internal_server_error" {
		h.logger.Error(ctx, "websocket action failed", zap.String("type", in.Type), zap.String("user_id", userID), zap.Error(err))
	}
	client.TrySend(dto.WSOutgoingMessage{
		Type:    "error",
		Payload: dto.WSErrorEvent{Code: code, Message: msg, ClientMessageID: in.ClientMessageID},
	})
}

// watchSession closes the socket when the access token expires (the client
// reconnects with a fresh token) or when the account is blocked or removed.
func (h *chatHandler) watchSession(ctx context.Context, client *realtime.Client, claims *jwt.AccessClaims) {
	var expiry <-chan time.Time
	if claims.ExpiresAt != nil {
		timer := time.NewTimer(time.Until(claims.ExpiresAt.Time))
		defer timer.Stop()
		expiry = timer.C
	}
	recheck := time.NewTicker(wsStatusRecheck)
	defer recheck.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-expiry:
			client.Close(realtime.CloseTokenExpired, "access token expired")
			return
		case <-recheck.C:
			err := h.service.CheckCanConnect(ctx, claims.UserID)
			if errors.Is(err, dto.ErrAccountCannotChat) {
				client.Close(realtime.CloseAccountUnavailable, "account unavailable")
				return
			}
		}
	}
}
