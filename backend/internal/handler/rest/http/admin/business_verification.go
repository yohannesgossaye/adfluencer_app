package admin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	"adinfluence_service_app/internal/constant/model/response"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *adminHandler) ListBusinessVerifications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}

	query := adminDto.ListBusinessVerificationsQuery{
		Status: r.URL.Query().Get("status"),
		Search: r.URL.Query().Get("search"),
		Page:   parseIntDefault(r.URL.Query().Get("page"), 1),
		Limit:  parseIntDefault(r.URL.Query().Get("limit"), 20),
	}

	res, err := h.service.ListBusinessVerifications(ctx, query)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "business verifications retrieved successfully", res, nil)
}

func (h *adminHandler) GetBusinessVerification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}

	userID := chi.URLParam(r, "user_id")
	res, err := h.service.GetBusinessVerification(ctx, userID)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "business verification retrieved successfully", res, nil)
}

func (h *adminHandler) VerifyBusiness(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	admin, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}
	adminID := admin.UserID

	r.Body = http.MaxBytesReader(w, r.Body, maxAdminBodySize)
	var req adminDto.VerifyBusinessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && r.ContentLength > 0 {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}

	targetUserID := chi.URLParam(r, "user_id")
	res, err := h.service.VerifyBusiness(ctx, adminID, targetUserID, req, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "business verified successfully", res, nil)
}

func (h *adminHandler) RevokeBusinessVerification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	admin, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}
	adminID := admin.UserID

	r.Body = http.MaxBytesReader(w, r.Body, maxAdminBodySize)
	var req adminDto.RevokeBusinessVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && r.ContentLength > 0 {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}

	targetUserID := chi.URLParam(r, "user_id")
	res, err := h.service.RevokeBusinessVerification(ctx, adminID, targetUserID, req, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "business verification revoked successfully", res, nil)
}
