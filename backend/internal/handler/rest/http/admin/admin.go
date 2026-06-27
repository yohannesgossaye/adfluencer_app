package admin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model/response"
	"adinfluence_service_app/internal/handler/middleware"
	"adinfluence_service_app/internal/module"
	"adinfluence_service_app/platform/jwt"
	"adinfluence_service_app/platform/utils"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type adminHandler struct {
	service module.AdminService
	logger  utils.Logger
}

func InitAdminHandler(s module.AdminService, log utils.Logger) *adminHandler {
	return &adminHandler{service: s, logger: log}
}

const maxAdminBodySize = 1 << 20

func (h *adminHandler) ListApplications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}

	query := adminDto.ListApplicationsQuery{
		Status: r.URL.Query().Get("status"),
		Search: r.URL.Query().Get("search"),
		Page:   parseIntDefault(r.URL.Query().Get("page"), 1),
		Limit:  parseIntDefault(r.URL.Query().Get("limit"), 20),
	}

	res, err := h.service.ListApplications(ctx, query)
	if err != nil {
		h.handleError(w, ctx, "list applications", err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "applications retrieved successfully", res, nil)
}

func (h *adminHandler) GetApplication(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}

	userID := chi.URLParam(r, "user_id")
	res, err := h.service.GetApplication(ctx, userID)
	if err != nil {
		h.handleError(w, ctx, "get application", err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "application retrieved successfully", res, nil)
}

func (h *adminHandler) ApproveApplication(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	adminID, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxAdminBodySize)
	var req adminDto.ApproveApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && r.ContentLength > 0 {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}

	targetUserID := chi.URLParam(r, "user_id")
	res, err := h.service.ApproveApplication(ctx, adminID, targetUserID, req, sessionMetaFromRequest(r))
	if err != nil {
		h.handleError(w, ctx, "approve application", err)
		return
	}

	h.logger.Info(ctx, "influencer application approved",
		zap.String("admin_id", adminID),
		zap.String("user_id", targetUserID),
	)
	response.SendSuccessResponse(w, r, http.StatusOK, "application approved successfully", res, nil)
}

func (h *adminHandler) RejectApplication(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	adminID, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxAdminBodySize)
	var req adminDto.RejectApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}
	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	targetUserID := chi.URLParam(r, "user_id")
	res, err := h.service.RejectApplication(ctx, adminID, targetUserID, req, sessionMetaFromRequest(r))
	if err != nil {
		h.handleError(w, ctx, "reject application", err)
		return
	}

	h.logger.Info(ctx, "influencer application rejected",
		zap.String("admin_id", adminID),
		zap.String("user_id", targetUserID),
	)
	response.SendSuccessResponse(w, r, http.StatusOK, "application rejected successfully", res, nil)
}

func (h *adminHandler) UpdateAdminNotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxAdminBodySize)
	var req adminDto.UpdateAdminNotesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}
	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	targetUserID := chi.URLParam(r, "user_id")
	if err := h.service.UpdateAdminNotes(ctx, targetUserID, req); err != nil {
		h.handleError(w, ctx, "update admin notes", err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "admin notes updated successfully", nil, nil)
}

func (h *adminHandler) requireAdmin(w http.ResponseWriter, r *http.Request) (string, bool) {
	claims, ok := middleware.GetAccessClaims(r.Context())
	if !ok {
		response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", nil, "unauthorized")
		return "", false
	}
	if claims.Role != "admin" {
		response.SendErrorResponse(w, http.StatusForbidden, "forbidden", customErr.ErrForbiddenRole, "admin access required")
		return "", false
	}
	return userIDFromClaims(claims), true
}

func userIDFromClaims(claims *jwt.AccessClaims) string {
	if claims.UserID != "" {
		return claims.UserID
	}
	return claims.Subject
}

func sessionMetaFromRequest(r *http.Request) authDto.SessionMeta {
	return authDto.SessionMeta{
		IPAddress: utils.ClientIPFromRequest(r),
		UserAgent: r.UserAgent(),
	}
}

func parseIntDefault(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v < 1 {
		return fallback
	}
	return v
}

func (h *adminHandler) handleError(w http.ResponseWriter, ctx context.Context, action string, err error) {
	if status, ok := customErr.HTTPStatusFromError(err); ok {
		h.logger.Warn(ctx, action+" rejected", zap.String("reason", customErr.ClientErrorCode(err)))
		response.SendErrorResponse(w, status, customErr.ClientErrorCode(err), err, err.Error())
		return
	}
	h.logger.Error(ctx, action+" failed", zap.String("reason", err.Error()))
	response.SendErrorResponse(w, http.StatusInternalServerError, "internal_server_error", err, "internal server error")
}
