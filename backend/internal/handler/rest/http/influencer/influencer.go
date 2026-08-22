package influencer

import (
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	customErr "adinfluence_service_app/internal/constant/errors"
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
	"adinfluence_service_app/internal/constant/model/response"
	"adinfluence_service_app/internal/handler/middleware"
	"adinfluence_service_app/internal/handler/rest/httputil"
	"adinfluence_service_app/internal/module"
	"adinfluence_service_app/platform/utils"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type influencerHandler struct {
	service module.InfluencerService
	logger utils.Logger
}

func InitInfluencerHandler(s module.InfluencerService,log utils.Logger) *influencerHandler {
	return &influencerHandler{
		service: s,
		logger: log,
	}
}

const maxProfileBodySize = 1 << 20

func (h *influencerHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}

	res, err := h.service.GetProfile(ctx, user.UserID)
	if err != nil {
		writeError(w, err)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "profile retrieved successfully", res, nil)
}

func (h *influencerHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, maxProfileBodySize)

	user, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}
	userID := user.UserID

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
		writeError(w, err)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "profile updated successfully", res, nil)
}

func (h *influencerHandler) SubmitProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	res, err := h.service.SubmitProfile(ctx, userID, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "profile submitted for review", res, nil)
}

const maxAvatarBodySize = 6 << 20

func (h *influencerHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, maxAvatarBodySize)

	user, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	filename, contentType, data, ok := parseMultipartFile(w, r, maxAvatarBodySize)
	if !ok {
		return
	}

	res, err := h.service.UploadAvatar(ctx, userID, filename, contentType, data, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "avatar uploaded successfully", res, nil)
}

func (h *influencerHandler) DeleteAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	if err := h.service.DeleteAvatar(ctx, userID, sessionMetaFromRequest(r)); err != nil {
		writeError(w, err)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "avatar deleted successfully", nil, nil)
}

const maxPortfolioBodySize = 105 << 20

func (h *influencerHandler) UploadPortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, maxPortfolioBodySize)

	user, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	filename, contentType, data, ok := parseMultipartFile(w, r, maxPortfolioBodySize)
	if !ok {
		return
	}

	res, err := h.service.UploadPortfolio(ctx, userID, filename, contentType, data, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "portfolio sample uploaded successfully", res, nil)
}

func (h *influencerHandler) DeletePortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	itemID := chi.URLParam(r, "item_id")
	if strings.TrimSpace(itemID) == "" {
		response.SendErrorResponse(w, http.StatusBadRequest, "bad_request", customErr.ErrBadRequest, "item_id is required")
		return
	}

	res, err := h.service.DeletePortfolio(ctx, userID, itemID, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "portfolio sample deleted successfully", res, nil)
}

func parseMultipartFile(w http.ResponseWriter, r *http.Request, maxSize int64) (filename, contentType string, data []byte, ok bool) {
	if err := r.ParseMultipartForm(maxSize); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_file", customErr.ErrInvalidFile, "invalid multipart form")
		return "", "", nil, false
	}

	keys := []string{"file", "avatar", "media", "image", "logo"}
	var file multipart.File
	var header *multipart.FileHeader
	var err error

	for _, k := range keys {
		file, header, err = r.FormFile(k)
		if err == nil {
			break
		}
	}

	if file == nil || err != nil {
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

func (h *influencerHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	items, err := h.service.ListActiveCategories(ctx)
	if err != nil {
		writeError(w, err)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "categories listed successfully", map[string]interface{}{"items": items}, nil)
}

func (h *influencerHandler) requireInfluencer(w http.ResponseWriter, r *http.Request) (*middleware.UserPayload, bool) {
	user, err := middleware.ExtractUserInfo(r.Context())
	if err != nil {
		response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", err, "unauthorized")
		return nil, false
	}
	if user.Role != jwtRoleInfluencer() {
		response.SendErrorResponse(w, http.StatusForbidden, "forbidden", customErr.ErrForbiddenRole, "influencer access required")
		return nil, false
	}
	return user, true
}

func jwtRoleInfluencer() string {
	return "influencer"
}

func sessionMetaFromRequest(r *http.Request) authDto.SessionMeta {
	return authDto.SessionMeta{
		IPAddress: utils.ClientIPFromRequest(r),
		UserAgent: r.UserAgent(),
	}
}

func writeError(w http.ResponseWriter, err error) {
	httputil.WriteServiceError(w, err,
		httputil.ErrorMapping{Err: influencerDto.ErrProfileNotFound, Status: http.StatusNotFound, Code: "profile_not_found"},
		httputil.ErrorMapping{Err: influencerDto.ErrUnknownCategory, Status: http.StatusBadRequest, Code: "invalid_niche"},
	)
}
