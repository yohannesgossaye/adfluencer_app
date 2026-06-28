package admin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model/response"
	"adinfluence_service_app/internal/handler/middleware"
	"adinfluence_service_app/internal/handler/rest/httputil"
	"adinfluence_service_app/internal/module"
	"adinfluence_service_app/platform/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type adminHandler struct {
	service module.AdminService
	logger utils.Logger
}

func InitAdminHandler(s module.AdminService,log utils.Logger) *adminHandler {
	return &adminHandler{
		service: s,
		logger: log,

	}
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
		writeError(w, err)
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
		writeError(w, err)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "application retrieved successfully", res, nil)
}

func (h *adminHandler) ApproveApplication(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	admin, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}
	adminID := admin.UserID

	r.Body = http.MaxBytesReader(w, r.Body, maxAdminBodySize)
	var req adminDto.ApproveApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && r.ContentLength > 0 {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}

	targetUserID := chi.URLParam(r, "user_id")
	res, err := h.service.ApproveApplication(ctx, adminID, targetUserID, req, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "application approved successfully", res, nil)
}

func (h *adminHandler) RejectApplication(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	admin, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}
	adminID := admin.UserID

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
		writeError(w, err)
		return
	}
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
		writeError(w, err)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "admin notes updated successfully", nil, nil)
}

func (h *adminHandler) requireAdmin(w http.ResponseWriter, r *http.Request) (*middleware.UserPayload, bool) {
	user, err := middleware.ExtractUserInfo(r.Context())
	if err != nil {
		response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", err, "unauthorized")
		return nil, false
	}
	if user.Role != "admin" {
		response.SendErrorResponse(w, http.StatusForbidden, "forbidden", customErr.ErrForbiddenRole, "admin access required")
		return nil, false
	}
	return user, true
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

func writeError(w http.ResponseWriter, err error) {
	httputil.WriteServiceError(w, err)
}
