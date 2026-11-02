package admin

import (
	"encoding/json"
	"net/http"

	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	"adinfluence_service_app/internal/constant/model/response"

	"github.com/go-chi/chi/v5"
)

func (h *adminHandler) ListEscrowQueue(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	res, err := h.service.ListEscrowQueue(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "contracts awaiting payment retrieved successfully", res, nil)
}

func (h *adminHandler) FundContractEscrow(w http.ResponseWriter, r *http.Request) {
	admin, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxAdminBodySize)
	var req adminDto.FundEscrowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && r.ContentLength > 0 {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}
	res, err := h.service.FundContractEscrow(r.Context(), admin.UserID, chi.URLParam(r, "contract_id"), req, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "payment confirmed; escrow funded", res, nil)
}
