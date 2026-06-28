package admin

import (
	"adinfluence_service_app/internal/constant/model/response"
	"net/http"
)

func (h *adminHandler) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}

	res, err := h.service.GetDashboardSummary(ctx)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "dashboard summary retrieved successfully", res, nil)
}
