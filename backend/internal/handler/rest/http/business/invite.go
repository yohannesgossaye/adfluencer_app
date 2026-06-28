package business

import (
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	"adinfluence_service_app/internal/constant/model/response"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const maxInviteBodySize = 1 << 20

func (h *businessHandler) CreateInvite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, maxInviteBodySize)

	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	campaignID := chi.URLParam(r, "campaign_id")

	var req businessDto.CreateInviteRequest
	if decodeJSON(w, r, &req) {
		return
	}

	res, err := h.service.CreateInvite(ctx, userID, campaignID, req, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusCreated, "invite sent successfully", res, nil)
}

func (h *businessHandler) ListInvites(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	campaignID := chi.URLParam(r, "campaign_id")
	res, err := h.service.ListInvites(ctx, userID, campaignID)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "invites retrieved successfully", res, nil)
}

func (h *businessHandler) CancelInvite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	campaignID := chi.URLParam(r, "campaign_id")
	inviteID := chi.URLParam(r, "invite_id")

	if err := h.service.CancelInvite(ctx, userID, campaignID, inviteID, sessionMetaFromRequest(r)); err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "invite cancelled successfully", nil, nil)
}
