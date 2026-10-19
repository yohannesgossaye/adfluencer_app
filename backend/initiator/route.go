package initiator

import (
	"adinfluence_service_app/internal/constant/model/response"
	adminRoutes "adinfluence_service_app/internal/glue/routing/admin"
	auth "adinfluence_service_app/internal/glue/routing/auth"
	businessRoutes "adinfluence_service_app/internal/glue/routing/business"
	influencerRoutes "adinfluence_service_app/internal/glue/routing/influencer"
	marketplaceRoutes "adinfluence_service_app/internal/glue/routing/marketplace"
	"adinfluence_service_app/internal/handler/middleware"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func InirRoute(router chi.Router, rest handlers, authMW middleware.AuthMiddleware) {
	// Health check — used by Docker, Nginx upstream checks, and CI/CD pipeline
	// Always responds 200 OK as long as the process is alive
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.SendSuccessResponse(w, r, http.StatusOK, "ok", map[string]string{
			"status":  "healthy",
			"service": "adfluence-backend",
			"time":    time.Now().UTC().Format(time.RFC3339),
		}, nil)
	})

	router.Route("/api/v1", func(r chi.Router) {
		auth.AuthRoutes(r, rest.auth, authMW)
		influencerRoutes.InfluencerRoutes(r, rest.influencer, authMW)
		businessRoutes.BusinessRoutes(r, rest.business, authMW)
		adminRoutes.AdminRoutes(r, rest.admin, authMW)
		marketplaceRoutes.MarketPlaceRoutes(r, rest.chat, authMW)
	})

	// API root — quick connectivity check
	router.Get("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		response.SendSuccessResponse(w, r, http.StatusOK, "🚀 Adfluence service is running", nil, nil)
	})
}
