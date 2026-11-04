package business

import (
	"net/http"

	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	"adinfluence_service_app/internal/constant/model/response"
)

// SearchCreators is "Find creators": GET /business/creators?q=&niche=&platform=&size=&max_price=&city=&campaign_id=&sort=&page=&limit=
func (h *businessHandler) SearchCreators(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireBusiness(w, r)
	if !ok {
		return
	}
	v := r.URL.Query()
	res, err := h.service.SearchCreators(r.Context(), user.UserID, businessDto.SearchCreatorsQuery{
		Q:          v.Get("q"),
		Niche:      v.Get("niche"),
		Platform:   v.Get("platform"),
		Size:       v.Get("size"),
		MaxPrice:   parseFloatDefault(v.Get("max_price"), 0),
		City:       v.Get("city"),
		CampaignID: v.Get("campaign_id"),
		Sort:       v.Get("sort"),
		Page:       parseIntDefault(v.Get("page"), 1),
		Limit:      parseIntDefault(v.Get("limit"), 24),
	})
	if err != nil {
		writeError(w, err)
		return
	}
	response.SendSuccessResponse(w, r, http.StatusOK, "creators retrieved successfully", res, nil)
}
