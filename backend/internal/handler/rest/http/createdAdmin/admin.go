package createdadmin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	"adinfluence_service_app/internal/constant/model/response"
	"adinfluence_service_app/internal/handler/rest/httputil"
	"adinfluence_service_app/internal/module"
	"adinfluence_service_app/platform/utils"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type adminCreateHandler struct {
	handler           module.CreateAdminService
	adminCreateApiKey string
	logger            utils.Logger
}

func InitCreateAdminHandler(s module.CreateAdminService, adminCreateApiKey string, log utils.Logger) *adminCreateHandler {
	return &adminCreateHandler{
		handler:           s,
		adminCreateApiKey: adminCreateApiKey,
		logger:            log,
	}
}

func (h *adminCreateHandler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	var req adminDto.CreateAdminRequest
	ctx := r.Context()
	apiKeyValue := h.adminCreateApiKey
	if apiKeyValue == "" {
		h.logger.Error(ctx, "unable to find adminApiKey from vault user hardcode apikey.....")
		apiKeyValue = adminDto.AdminApikey
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(ctx, "Failed to decode request body", zap.Error(err))
		response.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body", err, "Invalid request body")
		return
	}
	if err := req.Validate(); err != nil {
		h.logger.Error(ctx, "Invalid request data", zap.Error(err))
		response.SendErrorResponse(w, http.StatusBadRequest, "Invalid request data", err, "Invalid request data")
		return
	}
	apiKey := r.Header.Get("X-Api-key")
	if apiKey == "" {
		h.logger.Error(ctx, "Missing API key in request header")
		response.SendErrorResponse(w, http.StatusUnauthorized, "Missing API key", nil, "Missing API key")
		return
	}
	if apiKey != apiKeyValue {
		h.logger.Error(ctx, "Invalid API key provided")
		response.SendErrorResponse(w, http.StatusUnauthorized, "Invalid API key", nil, "Invalid API key")
		return
	}
	if err := h.handler.CreateAdmin(ctx, req); err != nil {
		h.logger.Error(ctx, "Failed to create admin", zap.Error(err))
		httputil.WriteServiceErrorOr(w, err, "created_admin", "unable to create admin")
		return
	}
	response.SendSuccessResponse(w, r, http.StatusCreated, "Admin created successfully", nil, nil)

}
