package business

import (
	"adinfluence_service_app/internal/constant/model/response"
	"net/http"
)

func (h *businessHandler) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}

	res, err := h.service.GetDashboardSummary(ctx, user.UserID)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "dashboard summary retrieved successfully", res, nil)
}
