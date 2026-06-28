package business

import (
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	"adinfluence_service_app/internal/constant/model/response"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const maxCampaignBodySize = 1 << 20

func (h *businessHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, maxCampaignBodySize)

	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	var req businessDto.CreateCampaignRequest
	if decodeJSON(w, r, &req) {
		return
	}

	res, err := h.service.CreateCampaign(ctx, userID, req, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusCreated, "campaign created successfully", res, nil)
}

func (h *businessHandler) ListCampaigns(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	query := businessDto.ListCampaignsQuery{
		Status: r.URL.Query().Get("status"),
		Page:   parseIntDefault(r.URL.Query().Get("page"), 1),
		Limit:  parseIntDefault(r.URL.Query().Get("limit"), 20),
	}

	res, err := h.service.ListCampaigns(ctx, userID, query)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "campaigns retrieved successfully", res, nil)
}

func (h *businessHandler) GetCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	campaignID := chi.URLParam(r, "campaign_id")
	res, err := h.service.GetCampaign(ctx, userID, campaignID)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "campaign retrieved successfully", res, nil)
}

func (h *businessHandler) UpdateCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, maxCampaignBodySize)

	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	campaignID := chi.URLParam(r, "campaign_id")

	var req businessDto.UpdateCampaignRequest
	if decodeJSON(w, r, &req) {
		return
	}

	res, err := h.service.UpdateCampaign(ctx, userID, campaignID, req, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "campaign updated successfully", res, nil)
}

func (h *businessHandler) PublishCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	campaignID := chi.URLParam(r, "campaign_id")
	res, err := h.service.PublishCampaign(ctx, userID, campaignID, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "campaign published successfully", res, nil)
}

func (h *businessHandler) CancelCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	campaignID := chi.URLParam(r, "campaign_id")
	res, err := h.service.CancelCampaign(ctx, userID, campaignID, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "campaign cancelled successfully", res, nil)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return true
	}
	return false
}

func parseIntDefault(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return fallback
	}
	return n
}
