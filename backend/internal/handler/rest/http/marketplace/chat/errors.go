package chat

import (
	dto "adinfluence_service_app/internal/constant/model/dto/marketplace"
	"adinfluence_service_app/internal/constant/model/response"
	"errors"
	"net/http"
)

var chatErrorStatus = []struct {
	err    error
	status int
	code   string
}{
	{dto.ErrConversationNotFound, http.StatusNotFound, "conversation_not_found"},
	{dto.ErrMessageNotFound, http.StatusNotFound, "message_not_found"},
	{dto.ErrRecipientNotFound, http.StatusNotFound, "recipient_not_found"},
	{dto.ErrCampaignNotFound, http.StatusNotFound, "campaign_not_found"},
	{dto.ErrNotParticipant, http.StatusForbidden, "not_participant"},
	{dto.ErrChatNotAllowed, http.StatusForbidden, "chat_not_allowed"},
	{dto.ErrCampaignNotOwned, http.StatusForbidden, "campaign_not_owned"},
	{dto.ErrAccountCannotChat, http.StatusForbidden, "account_cannot_chat"},
	{dto.ErrRecipientUnavailable, http.StatusForbidden, "recipient_unavailable"},
	{dto.ErrDeleteNotAllowed, http.StatusForbidden, "delete_not_allowed"},
	{dto.ErrSendRateLimited, http.StatusTooManyRequests, "rate_limited"},
	{dto.ErrBusinessNotVerified, http.StatusForbidden, "business_not_verified"},
}

// classifyChatError maps a service error to an HTTP status, a stable code
// and a message that is safe to show. Unknown errors are server errors.
func classifyChatError(err error) (int, string, string) {
	var invalid *dto.ValidationError
	if errors.As(err, &invalid) {
		return http.StatusBadRequest, "invalid_request", invalid.Msg
	}
	for _, m := range chatErrorStatus {
		if errors.Is(err, m.err) {
			return m.status, m.code, m.err.Error()
		}
	}
	return http.StatusInternalServerError, "internal_server_error", "internal server error"
}

func writeChatError(w http.ResponseWriter, err error) {
	status, code, msg := classifyChatError(err)
	response.SendErrorResponse(w, status, code, err, msg)
}

func writeBodyError(w http.ResponseWriter, err error) {
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		response.SendErrorResponse(w, http.StatusRequestEntityTooLarge, "body_too_large", nil, "request body is too large")
		return
	}
	response.SendErrorResponse(w, http.StatusBadRequest, "invalid_body", nil, "invalid request body")
}
