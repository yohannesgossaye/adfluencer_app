package influencer

import (
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	customErr "adinfluence_service_app/internal/constant/errors"
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
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
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type influencerHandler struct {
	service module.InfluencerService
	auth    module.AuthService
	logger  utils.Logger
}

func InitInfluencerHandler(s module.InfluencerService, auth module.AuthService, log utils.Logger) *influencerHandler {
	return &influencerHandler{
		service: s,
		auth:    auth,
		logger:  log,
	}
}

const maxProfileBodySize = 1 << 20

func (h *influencerHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := h.requireInfluencer(w, r)
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

func (h *influencerHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, maxProfileBodySize)

	userID, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}

	var req influencerDto.UpdateProfileRequest
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

	h.logger.Info(ctx, "influencer profile updated", zap.String("user_id", userID))
	response.SendSuccessResponse(w, r, http.StatusOK, "profile updated successfully", res, nil)
}

func (h *influencerHandler) SubmitProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}

	res, err := h.service.SubmitProfile(ctx, userID, sessionMetaFromRequest(r))
	if err != nil {
		h.handleError(w, ctx, "submit profile", userID, err)
		return
	}

	h.logger.Info(ctx, "influencer profile submitted", zap.String("user_id", userID))
	response.SendSuccessResponse(w, r, http.StatusOK, "profile submitted for review", res, nil)
}

const maxAvatarBodySize = 6 << 20 // 6 MB (multipart overhead)

func (h *influencerHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, maxAvatarBodySize)

	userID, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}

	filename, contentType, data, ok := parseMultipartFile(w, r, maxAvatarBodySize)
	if !ok {
		return
	}

	res, err := h.service.UploadAvatar(ctx, userID, filename, contentType, data, sessionMetaFromRequest(r))
	if err != nil {
		h.handleError(w, ctx, "upload avatar", userID, err)
		return
	}

	h.logger.Info(ctx, "influencer avatar uploaded", zap.String("user_id", userID))
	response.SendSuccessResponse(w, r, http.StatusOK, "avatar uploaded successfully", res, nil)
}

const maxPortfolioBodySize = 105 << 20 // 100 MB video + multipart overhead

func (h *influencerHandler) UploadPortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, maxPortfolioBodySize)

	userID, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}

	filename, contentType, data, ok := parseMultipartFile(w, r, maxPortfolioBodySize)
	if !ok {
		return
	}

	res, err := h.service.UploadPortfolio(ctx, userID, filename, contentType, data, sessionMetaFromRequest(r))
	if err != nil {
		h.handleError(w, ctx, "upload portfolio", userID, err)
		return
	}

	h.logger.Info(ctx, "influencer portfolio uploaded", zap.String("user_id", userID))
	response.SendSuccessResponse(w, r, http.StatusOK, "portfolio sample uploaded successfully", res, nil)
}

func (h *influencerHandler) DeletePortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}

	itemID := chi.URLParam(r, "item_id")
	if strings.TrimSpace(itemID) == "" {
		response.SendErrorResponse(w, http.StatusBadRequest, "bad_request", customErr.ErrBadRequest, "item_id is required")
		return
	}

	res, err := h.service.DeletePortfolio(ctx, userID, itemID, sessionMetaFromRequest(r))
	if err != nil {
		h.handleError(w, ctx, "delete portfolio", userID, err)
		return
	}

	h.logger.Info(ctx, "influencer portfolio deleted", zap.String("user_id", userID), zap.String("item_id", itemID))
	response.SendSuccessResponse(w, r, http.StatusOK, "portfolio sample deleted successfully", res, nil)
}

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

func (h *influencerHandler) requireInfluencer(w http.ResponseWriter, r *http.Request) (string, bool) {
	claims, ok := middleware.GetAccessClaims(r.Context())
	if !ok {
		response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", nil, "unauthorized")
		return "", false
	}
	if claims.Role != string(jwtRoleInfluencer()) {
		response.SendErrorResponse(w, http.StatusForbidden, "forbidden", customErr.ErrForbiddenRole, "influencer access required")
		return "", false
	}
	return userIDFromClaims(claims), true
}

func jwtRoleInfluencer() string {
	return "influencer"
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

func (h *influencerHandler) handleError(w http.ResponseWriter, ctx context.Context, action, userID string, err error) {
	if errors.Is(err, influencerDto.ErrProfileNotFound) {
		response.SendErrorResponse(w, http.StatusNotFound, "profile_not_found", err, err.Error())
		return
	}
	if errors.Is(err, influencerDto.ErrUnknownCategory) {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_niche", err, err.Error())
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
