package admin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	"adinfluence_service_app/internal/constant/model/response"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func (h *adminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}

	query := adminDto.ListUsersQuery{
		Role:   r.URL.Query().Get("role"),
		Status: r.URL.Query().Get("status"),
		Search: r.URL.Query().Get("search"),
		Page:   parseIntDefault(r.URL.Query().Get("page"), 1),
		Limit:  parseIntDefault(r.URL.Query().Get("limit"), 20),
	}

	res, err := h.service.ListUsers(ctx, query)
	if err != nil {
		h.handleError(w, ctx, "list users", err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "users retrieved successfully", res, nil)
}

func (h *adminHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}

	userID := chi.URLParam(r, "user_id")
	res, err := h.service.GetUser(ctx, userID)
	if err != nil {
		h.handleError(w, ctx, "get user", err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "user retrieved successfully", res, nil)
}

func (h *adminHandler) BlockUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	adminID, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxAdminBodySize)
	var req adminDto.BlockUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && r.ContentLength > 0 {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}

	targetUserID := chi.URLParam(r, "user_id")
	res, err := h.service.BlockUser(ctx, adminID, targetUserID, req, sessionMetaFromRequest(r))
	if err != nil {
		h.handleError(w, ctx, "block user", err)
		return
	}

	h.logger.Info(ctx, "user blocked",
		zap.String("admin_id", adminID),
		zap.String("user_id", targetUserID),
	)
	response.SendSuccessResponse(w, r, http.StatusOK, "user blocked successfully", res, nil)
}

func (h *adminHandler) UnblockUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	adminID, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}

	targetUserID := chi.URLParam(r, "user_id")
	res, err := h.service.UnblockUser(ctx, adminID, targetUserID, sessionMetaFromRequest(r))
	if err != nil {
		h.handleError(w, ctx, "unblock user", err)
		return
	}

	h.logger.Info(ctx, "user unblocked",
		zap.String("admin_id", adminID),
		zap.String("user_id", targetUserID),
	)
	response.SendSuccessResponse(w, r, http.StatusOK, "user unblocked successfully", res, nil)
}

func (h *adminHandler) RemoveUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	adminID, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxAdminBodySize)
	var req adminDto.RemoveUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && r.ContentLength > 0 {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}

	targetUserID := chi.URLParam(r, "user_id")
	res, err := h.service.RemoveUser(ctx, adminID, targetUserID, req, sessionMetaFromRequest(r))
	if err != nil {
		h.handleError(w, ctx, "remove user", err)
		return
	}

	h.logger.Info(ctx, "user removed",
		zap.String("admin_id", adminID),
		zap.String("user_id", targetUserID),
	)
	response.SendSuccessResponse(w, r, http.StatusOK, "user removed successfully", res, nil)
}
