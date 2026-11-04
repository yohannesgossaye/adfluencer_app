package influencer

import (
	"encoding/json"
	"net/http"

	supportDto "adinfluence_service_app/internal/constant/model/dto/support"
	"adinfluence_service_app/internal/constant/model/response"

	"github.com/go-chi/chi/v5"
)

// ReplySupportTicket adds the user's message to their own support ticket.
func (h *influencerHandler) ReplySupportTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
	var req supportDto.ReplyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}
	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "validation_error", err, err.Error())
		return
	}
	res, err := h.service.ReplySupportTicket(ctx, user.UserID, chi.URLParam(r, "ticket_id"), req)
	if err != nil {
		writeError(w, err)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusCreated, "reply sent", res, nil)
}
