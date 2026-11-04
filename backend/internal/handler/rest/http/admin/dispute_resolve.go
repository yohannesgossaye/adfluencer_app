package admin

import (
	"encoding/json"
	"net/http"

	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
	"adinfluence_service_app/internal/constant/model/response"

	"github.com/go-chi/chi/v5"
)

// ResolveDispute settles a dispute: {"outcome": "continue|complete|terminate", "note": "..."}.
func (h *adminHandler) ResolveDispute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	admin, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxAdminBodySize)
	var req contractDto.ResolveDisputeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}
	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "validation_error", err, err.Error())
		return
	}
	if err := h.service.ResolveDispute(ctx, admin.UserID, chi.URLParam(r, "dispute_id"), req); err != nil {
		writeError(w, err)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "dispute resolved", nil, nil)
}
