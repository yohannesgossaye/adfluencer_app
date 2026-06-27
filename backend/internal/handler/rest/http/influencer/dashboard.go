package influencer

import (
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model/response"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const maxSupportBodySize = 1 << 20

func (h *influencerHandler) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}

	user, err := h.auth.GetMe(ctx, userID)
	if err != nil {
		h.handleError(w, ctx, "get dashboard summary", userID, err)
		return
	}

	res, err := h.service.GetDashboardSummary(ctx, userID, &user.User)
	if err != nil {
		h.handleError(w, ctx, "get dashboard summary", userID, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "dashboard summary retrieved successfully", res, nil)
}

func (h *influencerHandler) ListSupportInviteOptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}

	res, err := h.service.ListSupportInviteOptions(ctx, userID)
	if err != nil {
		h.handleError(w, ctx, "list support invite options", userID, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "support invite options retrieved successfully", res, nil)
}

func (h *influencerHandler) CreateSupportTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, maxSupportBodySize)

	userID, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}

	var req influencerDto.CreateSupportTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}
	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	res, err := h.service.CreateSupportTicket(ctx, userID, req, sessionMetaFromRequest(r))
	if err != nil {
		h.handleError(w, ctx, "create support ticket", userID, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusCreated, "support ticket created successfully", res, nil)
}

func (h *influencerHandler) ListSupportTickets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}

	res, err := h.service.ListSupportTickets(ctx, userID)
	if err != nil {
		h.handleError(w, ctx, "list support tickets", userID, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "support tickets retrieved successfully", res, nil)
}

func (h *influencerHandler) GetSupportTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}

	ticketID := chi.URLParam(r, "ticket_id")
	res, err := h.service.GetSupportTicket(ctx, userID, ticketID)
	if err != nil {
		h.handleError(w, ctx, "get support ticket", userID, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "support ticket retrieved successfully", res, nil)
}
