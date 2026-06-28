package admin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model/response"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *adminHandler) ListSupportTickets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}

	query := adminDto.ListSupportTicketsQuery{
		Role:     r.URL.Query().Get("role"),
		Status:   r.URL.Query().Get("status"),
		Priority: r.URL.Query().Get("priority"),
		Search:   r.URL.Query().Get("search"),
		Page:     parseIntDefault(r.URL.Query().Get("page"), 1),
		Limit:    parseIntDefault(r.URL.Query().Get("limit"), 20),
	}

	res, err := h.service.ListSupportTickets(ctx, query)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "support tickets retrieved successfully", res, nil)
}

func (h *adminHandler) GetSupportTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}

	ticketID := chi.URLParam(r, "ticket_id")
	res, err := h.service.GetSupportTicket(ctx, ticketID)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "support ticket retrieved successfully", res, nil)
}

func (h *adminHandler) UpdateSupportTicketStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	admin, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}
	adminID := admin.UserID

	r.Body = http.MaxBytesReader(w, r.Body, maxAdminBodySize)
	var req adminDto.UpdateSupportTicketStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}
	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	ticketID := chi.URLParam(r, "ticket_id")
	res, err := h.service.UpdateSupportTicketStatus(ctx, adminID, ticketID, req, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "support ticket status updated successfully", res, nil)
}

func (h *adminHandler) ReplySupportTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	admin, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}
	adminID := admin.UserID

	r.Body = http.MaxBytesReader(w, r.Body, maxAdminBodySize)
	var req adminDto.ReplySupportTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && r.ContentLength > 0 {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}

	ticketID := chi.URLParam(r, "ticket_id")
	res, err := h.service.ReplySupportTicket(ctx, adminID, ticketID, req, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "support ticket reply saved successfully", res, nil)
}
