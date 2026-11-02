package admin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"adinfluence_service_app/internal/constant/model/response"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *adminHandler) ListAdminCampaigns(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}

	query := adminDto.ListAdminCampaignsQuery{
		Status:   r.URL.Query().Get("status"),
		Niche:    r.URL.Query().Get("niche"),
		Business: r.URL.Query().Get("business"),
		Page:     parseIntDefault(r.URL.Query().Get("page"), 1),
		Limit:    parseIntDefault(r.URL.Query().Get("limit"), 20),
	}

	res, err := h.service.ListAdminCampaigns(ctx, query)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "campaigns retrieved successfully", res, nil)
}

func (h *adminHandler) GetAdminCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}

	campaignID := chi.URLParam(r, "campaign_id")
	res, err := h.service.GetAdminCampaign(ctx, campaignID)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "campaign retrieved successfully", res, nil)
}

func (h *adminHandler) PauseAdminCampaign(w http.ResponseWriter, r *http.Request) {
	h.moderateCampaign(w, r, func(ctx context.Context, adminID, campaignID, reason string, meta authDto.SessionMeta) (*adminDto.ModerateCampaignResponse, error) {
		return h.service.PauseAdminCampaign(ctx, adminID, campaignID, reason, meta)
	})
}

func (h *adminHandler) CompleteAdminCampaign(w http.ResponseWriter, r *http.Request) {
	h.moderateCampaign(w, r, func(ctx context.Context, adminID, campaignID, reason string, meta authDto.SessionMeta) (*adminDto.ModerateCampaignResponse, error) {
		return h.service.CompleteAdminCampaign(ctx, adminID, campaignID, reason, meta)
	})
}

func (h *adminHandler) FlagAdminCampaign(w http.ResponseWriter, r *http.Request) {
	h.moderateCampaign(w, r, func(ctx context.Context, adminID, campaignID, reason string, meta authDto.SessionMeta) (*adminDto.ModerateCampaignResponse, error) {
		return h.service.FlagAdminCampaign(ctx, adminID, campaignID, reason, meta)
	})
}

func (h *adminHandler) ResumeAdminCampaign(w http.ResponseWriter, r *http.Request) {
	h.moderateCampaign(w, r, func(ctx context.Context, adminID, campaignID, reason string, meta authDto.SessionMeta) (*adminDto.ModerateCampaignResponse, error) {
		return h.service.ResumeAdminCampaign(ctx, adminID, campaignID, reason, meta)
	})
}

type moderateCampaignFunc func(ctx context.Context, adminID, campaignID, reason string, meta authDto.SessionMeta) (*adminDto.ModerateCampaignResponse, error)

func (h *adminHandler) moderateCampaign(w http.ResponseWriter, r *http.Request, fn moderateCampaignFunc) {
	ctx := r.Context()
	admin, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}
	adminID := admin.UserID

	// The body is optional: {"reason": "..."} is recorded in the audit log.
	r.Body = http.MaxBytesReader(w, r.Body, maxAdminBodySize)
	var req adminDto.ModerateCampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && r.ContentLength > 0 {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid_request", err, "invalid request body")
		return
	}

	campaignID := chi.URLParam(r, "campaign_id")
	res, err := fn(ctx, adminID, campaignID, req.Reason, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "campaign updated successfully", res, nil)
}
