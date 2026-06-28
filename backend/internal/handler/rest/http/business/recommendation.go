package business

import (
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	"adinfluence_service_app/internal/constant/model/response"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *businessHandler) GenerateRecommendations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	campaignID := chi.URLParam(r, "campaign_id")
	res, err := h.service.GenerateRecommendations(ctx, userID, campaignID, sessionMetaFromRequest(r))
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "recommendations generated successfully", res, nil)
}

func (h *businessHandler) ListRecommendations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	campaignID := chi.URLParam(r, "campaign_id")
	query := businessDto.ListRecommendationsQuery{
		Search:   r.URL.Query().Get("search"),
		Niche:    r.URL.Query().Get("niche"),
		Platform: r.URL.Query().Get("platform"),
		MinScore: parseFloatDefault(r.URL.Query().Get("min_score"), businessDto.MinRecommendationScore),
		Sort:     r.URL.Query().Get("sort"),
		Page:     parseIntDefault(r.URL.Query().Get("page"), 1),
		Limit:    parseIntDefault(r.URL.Query().Get("limit"), 20),
	}

	res, err := h.service.ListRecommendations(ctx, userID, campaignID, query)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "recommendations retrieved successfully", res, nil)
}

func (h *businessHandler) GetInfluencerForBusiness(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	userID := user.UserID

	influencerUserID := chi.URLParam(r, "user_id")
	campaignID := r.URL.Query().Get("campaign_id")

	res, err := h.service.GetInfluencerForBusiness(ctx, userID, influencerUserID, campaignID)
	if err != nil {
		writeError(w, err)
		return
	}

	response.SendSuccessResponse(w, r, http.StatusOK, "creator profile retrieved successfully", res, nil)
}

func parseFloatDefault(raw string, fallback float64) float64 {
	if raw == "" {
		return fallback
	}
	n, err := strconv.ParseFloat(raw, 64)
	if err != nil || n < 0 {
		return fallback
	}
	return n
}
