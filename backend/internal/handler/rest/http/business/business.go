package business

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	"adinfluence_service_app/internal/constant/model/response"
	"adinfluence_service_app/internal/handler/middleware"
	"adinfluence_service_app/internal/handler/rest/httputil"
	"adinfluence_service_app/internal/module"
	"adinfluence_service_app/platform/utils"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
)

type businessHandler struct {
	service module.BusinessService
	logger  utils.Logger
}

func InitBusinessHandler(s module.BusinessService, log utils.Logger) *businessHandler {
	return &businessHandler{
		service: s,
		logger:  log,
	}
}

const maxProfileBodySize = 1 << 20

func (h *businessHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireBusiness(w, r)
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

func (h *businessHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, maxProfileBodySize)

	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

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
		writeError(w, err)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "profile updated successfully", res, nil)
}

func (h *businessHandler) UploadLogo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, maxLogoBodySize)

	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	filename, contentType, data, ok := parseMultipartFile(w, r, maxLogoBodySize)
	if !ok {
		return
	}

	res, err := h.service.UploadLogo(ctx, userID, filename, contentType, data, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "logo uploaded successfully", res, nil)
}

func (h *businessHandler) CompleteProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	res, err := h.service.CompleteProfile(ctx, userID, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "profile completed successfully", res, nil)
}

const maxLogoBodySize = 6 << 20

func parseMultipartFile(w http.ResponseWriter, r *http.Request, maxSize int64) (filename, contentType string, data []byte, ok bool) {
	if err := r.ParseMultipartForm(maxSize); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_file", customErr.ErrInvalidFile, "invalid multipart form")
		return "", "", nil, false
	}

	keys := []string{"file", "logo", "avatar", "media", "image"}
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

func (h *businessHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	items, err := h.service.ListActiveCategories(ctx)
	if err != nil {
		writeError(w, err)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "categories listed successfully", map[string]interface{}{"items": items}, nil)
}

func (h *businessHandler) requireBusiness(w http.ResponseWriter, r *http.Request) (*middleware.UserPayload, bool) {
	user, err := middleware.ExtractUserInfo(r.Context())
	if err != nil {
		response.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized", err, "unauthorized")
		return nil, false
	}
	if user.Role != jwtRoleBusiness() {
		response.SendErrorResponse(w, http.StatusForbidden, "forbidden", customErr.ErrForbiddenRole, "business access required")
		return nil, false
	}
	return user, true
}

func jwtRoleBusiness() string {
	return "business"
}

func sessionMetaFromRequest(r *http.Request) authDto.SessionMeta {
	return authDto.SessionMeta{
		IPAddress: utils.ClientIPFromRequest(r),
		UserAgent: r.UserAgent(),
	}
}

func writeError(w http.ResponseWriter, err error) {
	httputil.WriteServiceError(w, err,
		httputil.ErrorMapping{Err: businessDto.ErrProfileNotFound, Status: http.StatusNotFound, Code: "profile_not_found"},
	)
}
