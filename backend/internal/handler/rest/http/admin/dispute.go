package admin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model/response"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *adminHandler) ListAdminDisputes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}

	query := adminDto.ListAdminDisputesQuery{
		Tab:    r.URL.Query().Get("tab"),
		Search: r.URL.Query().Get("search"),
		Page:   parseIntDefault(r.URL.Query().Get("page"), 1),
		Limit:  parseIntDefault(r.URL.Query().Get("limit"), 20),
	}

	res, err := h.service.ListAdminDisputes(ctx, query)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "disputes retrieved successfully", res, nil)
}

func (h *adminHandler) GetAdminDispute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}

	disputeID := chi.URLParam(r, "dispute_id")
	res, err := h.service.GetAdminDispute(ctx, disputeID)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "dispute retrieved successfully", res, nil)
}

func (h *adminHandler) UpdateDisputeStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	admin, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}
	adminID := admin.UserID

	r.Body = http.MaxBytesReader(w, r.Body, maxAdminBodySize)
	var req adminDto.UpdateDisputeStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}
	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	disputeID := chi.URLParam(r, "dispute_id")
	res, err := h.service.UpdateDisputeStatus(ctx, adminID, disputeID, req, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "dispute status updated successfully", res, nil)
}

func (h *adminHandler) UpdateDisputeNotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	admin, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}
	adminID := admin.UserID

	r.Body = http.MaxBytesReader(w, r.Body, maxAdminBodySize)
	var req adminDto.UpdateDisputeNotesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}

	disputeID := chi.URLParam(r, "dispute_id")
	res, err := h.service.UpdateDisputeNotes(ctx, adminID, disputeID, req, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "dispute notes updated successfully", res, nil)
}
