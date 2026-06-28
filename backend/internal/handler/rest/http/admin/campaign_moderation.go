package admin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"adinfluence_service_app/internal/constant/model/response"
	"context"
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
	h.moderateCampaign(w, r, func(ctx context.Context, adminID, campaignID string, meta authDto.SessionMeta) (*adminDto.ModerateCampaignResponse, error) {
		return h.service.PauseAdminCampaign(ctx, adminID, campaignID, meta)
	})
}

func (h *adminHandler) CompleteAdminCampaign(w http.ResponseWriter, r *http.Request) {
	h.moderateCampaign(w, r, func(ctx context.Context, adminID, campaignID string, meta authDto.SessionMeta) (*adminDto.ModerateCampaignResponse, error) {
		return h.service.CompleteAdminCampaign(ctx, adminID, campaignID, meta)
	})
}

func (h *adminHandler) FlagAdminCampaign(w http.ResponseWriter, r *http.Request) {
	h.moderateCampaign(w, r, func(ctx context.Context, adminID, campaignID string, meta authDto.SessionMeta) (*adminDto.ModerateCampaignResponse, error) {
		return h.service.FlagAdminCampaign(ctx, adminID, campaignID, meta)
	})
}

type moderateCampaignFunc func(ctx context.Context, adminID, campaignID string, meta authDto.SessionMeta) (*adminDto.ModerateCampaignResponse, error)

func (h *adminHandler) moderateCampaign(w http.ResponseWriter, r *http.Request, fn moderateCampaignFunc) {
	ctx := r.Context()
	admin, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}
	adminID := admin.UserID

	campaignID := chi.URLParam(r, "campaign_id")
	res, err := fn(ctx, adminID, campaignID, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "campaign updated successfully", res, nil)
}
