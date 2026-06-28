package influencer

import (
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
	"adinfluence_service_app/internal/constant/model/response"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *influencerHandler) ListInvites(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	query := influencerDto.ListInvitesQuery{
		Status: r.URL.Query().Get("status"),
	}

	res, err := h.service.ListInvites(ctx, userID, query)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "invites retrieved successfully", res, nil)
}

func (h *influencerHandler) GetInvite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	inviteID := chi.URLParam(r, "id")
	res, err := h.service.GetInvite(ctx, userID, inviteID)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "invite retrieved successfully", res, nil)
}

func (h *influencerHandler) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	inviteID := chi.URLParam(r, "id")
	res, err := h.service.AcceptInvite(ctx, userID, inviteID, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "invite accepted successfully", res, nil)
}

func (h *influencerHandler) RejectInvite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	inviteID := chi.URLParam(r, "id")
	res, err := h.service.RejectInvite(ctx, userID, inviteID, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "invite rejected successfully", res, nil)
}
