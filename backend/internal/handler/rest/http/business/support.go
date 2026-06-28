package business

import (
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model/response"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const maxSupportBodySize = 1 << 20

func (h *businessHandler) ListSupportCampaignOptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	res, err := h.service.ListSupportCampaignOptions(ctx, userID)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "support campaign options retrieved successfully", res, nil)
}

func (h *businessHandler) CreateSupportTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, maxSupportBodySize)

	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	var req businessDto.CreateSupportTicketRequest
	if decodeJSON(w, r, &req) {
		return
	}
	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	res, err := h.service.CreateSupportTicket(ctx, userID, req, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusCreated, "support ticket created successfully", res, nil)
}

func (h *businessHandler) ListSupportTickets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	res, err := h.service.ListSupportTickets(ctx, userID)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "support tickets retrieved successfully", res, nil)
}

func (h *businessHandler) GetSupportTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	ticketID := chi.URLParam(r, "ticket_id")
	res, err := h.service.GetSupportTicket(ctx, userID, ticketID)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "support ticket retrieved successfully", res, nil)
}
