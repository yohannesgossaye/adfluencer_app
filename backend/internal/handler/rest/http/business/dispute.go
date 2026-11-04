package business

import (
	"encoding/json"
	"net/http"

	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
	"adinfluence_service_app/internal/constant/model/response"

	"github.com/go-chi/chi/v5"
)

// OpenDispute handles "Report a problem" on a contract.
func (h *businessHandler) OpenDispute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
	var req contractDto.OpenDisputeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}
	if err := req.Validate("business"); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "validation_error", err, err.Error())
		return
	}
	res, err := h.service.OpenDispute(ctx, user.UserID, chi.URLParam(r, "contract_id"), req)
	if err != nil {
		writeError(w, err)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusCreated, "problem reported", res, nil)
}
