package initiator

import (
	adminRoutes "adinfluence_service_app/internal/glue/routing/admin"
	auth "adinfluence_service_app/internal/glue/routing/auth"
	businessRoutes "adinfluence_service_app/internal/glue/routing/business"
	influencerRoutes "adinfluence_service_app/internal/glue/routing/influencer"
	marketplaceRoutes "adinfluence_service_app/internal/glue/routing/marketplace"
	"adinfluence_service_app/internal/handler/middleware"

	"github.com/go-chi/chi/v5"
)

func InirRoute(router chi.Router, rest handlers, authMW middleware.AuthMiddleware) {
	router.Route("/api/v1", func(r chi.Router) {
		auth.AuthRoutes(r, rest.auth, authMW)
		influencerRoutes.InfluencerRoutes(r, rest.influencer, authMW)
		businessRoutes.BusinessRoutes(r, rest.business, authMW)
		adminRoutes.AdminRoutes(r, rest.admin, authMW)
		marketplaceRoutes.MarketPlaceRoutes(r, rest.chat, authMW)
	})
}

