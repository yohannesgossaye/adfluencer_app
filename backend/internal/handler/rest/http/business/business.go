package business

import (
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	customErr "adinfluence_service_app/internal/constant/errors"
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	"adinfluence_service_app/internal/constant/model/response"
	"adinfluence_service_app/internal/handler/middleware"
	"adinfluence_service_app/internal/module"
	"adinfluence_service_app/platform/jwt"
	"adinfluence_service_app/platform/utils"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type businessHandler struct {
	service module.BusinessService
	auth    module.AuthService
	logger  utils.Logger
}

func InitBusinessHandler(s module.BusinessService, auth module.AuthService, log utils.Logger) *businessHandler {
	return &businessHandler{
		service: s,
		auth:    auth,
		logger:  log,
	}
}

const maxProfileBodySize = 1 << 20

func (h *businessHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}

	user, err := h.auth.GetMe(ctx, userID)
	if err != nil {
		h.handleError(w, ctx, "get profile", userID, err)
		return
	}

	res, err := h.service.GetProfile(ctx, userID, &user.User)
	if err != nil {
		h.handleError(w, ctx, "get profile", userID, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "profile retrieved successfully", res, nil)
}

func (h *businessHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, maxProfileBodySize)

	userID, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}

	var req businessDto.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	res, err := h.service.UpdateProfile(ctx, userID, req, sessionMetaFromRequest(r))
	if err != nil {
		h.handleError(w, ctx, "update profile", userID, err)
		return
	}

	h.logger.Info(ctx, "business profile updated", zap.String("user_id", userID))
	response.SendSuccessResponse(w, r, http.StatusOK, "profile updated successfully", res, nil)
}

func (h *businessHandler) UploadLogo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, maxLogoBodySize)

	userID, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}

	filename, contentType, data, ok := parseMultipartFile(w, r, maxLogoBodySize)
	if !ok {
		return
	}

	res, err := h.service.UploadLogo(ctx, userID, filename, contentType, data, sessionMetaFromRequest(r))
	if err != nil {
		h.handleError(w, ctx, "upload logo", userID, err)
		return
	}

	h.logger.Info(ctx, "business logo uploaded", zap.String("user_id", userID))
	response.SendSuccessResponse(w, r, http.StatusOK, "logo uploaded successfully", res, nil)
}

func (h *businessHandler) CompleteProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}

	res, err := h.service.CompleteProfile(ctx, userID, sessionMetaFromRequest(r))
	if err != nil {
		h.handleError(w, ctx, "complete profile", userID, err)
		return
	}

	h.logger.Info(ctx, "business profile completed", zap.String("user_id", userID))
	response.SendSuccessResponse(w, r, http.StatusOK, "profile completed successfully", res, nil)
}

const maxLogoBodySize = 6 << 20

func parseMultipartFile(w http.ResponseWriter, r *http.Request, maxSize int64) (filename, contentType string, data []byte, ok bool) {
	if err := r.ParseMultipartForm(maxSize); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_file", customErr.ErrInvalidFile, "invalid multipart form")
		return "", "", nil, false
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_file", customErr.ErrInvalidFile, "file field is required")
		return "", "", nil, false
	}
	defer file.Close()

	data, err = io.ReadAll(file)
	if err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_file", customErr.ErrInvalidFile, "failed to read file")
		return "", "", nil, false
	}

	if header != nil {
		filename = header.Filename
		contentType = header.Header.Get("Content-Type")
	}
	return filename, contentType, data, true
}

func (h *businessHandler) requireBusiness(w http.ResponseWriter, r *http.Request) (string, bool) {
	claims, ok := middleware.GetAccessClaims(r.Context())
	if !ok {
		response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", nil, "unauthorized")
		return "", false
	}
	if claims.Role != string(jwtRoleBusiness()) {
		response.SendErrorResponse(w, http.StatusForbidden, "forbidden", customErr.ErrForbiddenRole, "business access required")
		return "", false
	}
	return userIDFromClaims(claims), true
}

func jwtRoleBusiness() string {
	return "business"
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

func (h *businessHandler) handleError(w http.ResponseWriter, ctx context.Context, action, userID string, err error) {
	if errors.Is(err, businessDto.ErrProfileNotFound) {
		response.SendErrorResponse(w, http.StatusNotFound, "profile_not_found", err, err.Error())
		return
	}
	if status, ok := customErr.HTTPStatusFromError(err); ok {
		h.logger.Warn(ctx, action+" rejected",
			zap.String("user_id", userID),
			zap.String("reason", customErr.ClientErrorCode(err)),
		)
		response.SendErrorResponse(w, status, customErr.ClientErrorCode(err), err, err.Error())
		return
	}
	h.logger.Error(ctx, action+" failed", zap.String("user_id", userID), zap.String("reason", err.Error()))
	response.SendErrorResponse(w, http.StatusInternalServerError, "internal_server_error", err, "internal server error")
}
