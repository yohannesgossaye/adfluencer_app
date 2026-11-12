// Package createdadmin serves POST /api/v1/admin/create_admin: creating an
// Adfluence admin without seeding a fixed password into the database.
//
// It is a powerful endpoint, so it is locked down in layers, checked in order:
//  1. Off unless ADMIN_API_KEY (from Vault) is set and at least 32 characters.
//  2. Only from the server itself (loopback or the private Docker network).
//     Requests through nginx carry the visitor's real IP and are refused.
//  3. The X-Api-Key header must match, compared in constant time.
//
// Refusals at steps 1 and 2 answer 404, so from outside the endpoint looks
// like it doesn't exist.
package createdadmin

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"net/netip"

	customErr "adinfluence_service_app/internal/constant/errors"
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"adinfluence_service_app/internal/constant/model/response"
	"adinfluence_service_app/internal/handler/rest/httputil"
	"adinfluence_service_app/internal/module"
	"adinfluence_service_app/platform/utils"

	"go.uber.org/zap"
)

// MinAPIKeyLength keeps guessable keys from switching the endpoint on.
const MinAPIKeyLength = 32

type adminCreateHandler struct {
	service module.CreateAdminService
	apiKey  []byte
	logger  utils.Logger
}

func InitCreateAdminHandler(s module.CreateAdminService, apiKey string, log utils.Logger) *adminCreateHandler {
	h := &adminCreateHandler{service: s, logger: log}
	if len(apiKey) >= MinAPIKeyLength {
		h.apiKey = []byte(apiKey)
	}
	return h
}

// fromThisServer reports whether the caller is the server itself: loopback,
// or Docker's private network when calling the container's published port.
func fromThisServer(r *http.Request) bool {
	ip, err := netip.ParseAddr(utils.ClientIPFromRequest(r))
	if err != nil {
		return false
	}
	ip = ip.Unmap()
	return ip.IsLoopback() || ip.IsPrivate()
}

func (h *adminCreateHandler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if h.apiKey == nil {
		h.logger.Warn(ctx, "create_admin refused: ADMIN_API_KEY is not set (or shorter than 32 characters)")
		http.NotFound(w, r)
		return
	}
	if !fromThisServer(r) {
		h.logger.Warn(ctx, "create_admin refused: not called from the server itself", zap.String("ip", utils.ClientIPFromRequest(r)))
		http.NotFound(w, r)
		return
	}
	if subtle.ConstantTimeCompare([]byte(r.Header.Get("X-Api-Key")), h.apiKey) != 1 {
		h.logger.Warn(ctx, "create_admin refused: missing or wrong API key")
		response.SendErrorResponse(w, http.StatusUnauthorized, "invalid_api_key", nil, "missing or invalid API key")
		return
	}

	var req adminDto.CreateAdminRequest
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}
	req.Normalize()
	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, customErr.ClientErrorCode(err), err, err.Error())
		return
	}

	meta := authDto.SessionMeta{IPAddress: utils.ClientIPFromRequest(r), UserAgent: r.UserAgent()}
	userID, err := h.service.CreateAdmin(ctx, req, meta)
	if err != nil {
		httputil.WriteServiceErrorOr(w, err, "create_admin_failed", "unable to create admin")
		return
	}
	response.SendSuccessResponse(w, r, http.StatusCreated, "admin created", map[string]string{
		"user_id": userID,
		"email":   req.AdminEmail,
	}, nil)
}
