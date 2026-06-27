package influencer

import (
	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
	"adinfluence_service_app/internal/constant/model/response"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *influencerHandler) ListContracts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}

	query := contractDto.ListQuery{
		Status: r.URL.Query().Get("status"),
		Page:   parseIntDefault(r.URL.Query().Get("page"), 1),
		Limit:  parseIntDefault(r.URL.Query().Get("limit"), 20),
	}

	res, err := h.service.ListContracts(ctx, userID, query)
	if err != nil {
		h.handleError(w, ctx, "list contracts", userID, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "contracts retrieved successfully", res, nil)
}

func (h *influencerHandler) GetContract(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}

	contractID := chi.URLParam(r, "contract_id")
	res, err := h.service.GetContract(ctx, userID, contractID)
	if err != nil {
		h.handleError(w, ctx, "get contract", userID, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "contract retrieved successfully", res, nil)
}

func (h *influencerHandler) AcceptContract(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}

	contractID := chi.URLParam(r, "contract_id")
	res, err := h.service.AcceptContract(ctx, userID, contractID, sessionMetaFromRequest(r))
	if err != nil {
		h.handleError(w, ctx, "accept contract", userID, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "contract accepted successfully", res, nil)
}

func (h *influencerHandler) RejectContract(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := h.requireInfluencer(w, r)
	if !ok {
		return
	}

	contractID := chi.URLParam(r, "contract_id")
	res, err := h.service.RejectContract(ctx, userID, contractID, sessionMetaFromRequest(r))
	if err != nil {
		h.handleError(w, ctx, "reject contract", userID, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "contract rejected successfully", res, nil)
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
