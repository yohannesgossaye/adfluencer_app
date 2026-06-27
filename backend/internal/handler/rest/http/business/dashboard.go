package business

import (
	"adinfluence_service_app/internal/constant/model/response"
	"net/http"
)

func (h *businessHandler) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := h.requireBusiness(w, r)
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
