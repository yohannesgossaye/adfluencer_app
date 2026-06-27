package marketplace

import (
	"adinfluence_service_app/internal/glue"
	"adinfluence_service_app/internal/handler"
	"adinfluence_service_app/internal/handler/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func MarketPlaceRoutes(router chi.Router, rest handler.MarketplaceHandler, authMiddleware middleware.AuthMiddleware) {
	router.Route("/market", func(r chi.Router) {
		routes := []glue.Route{
			{
				Method:  http.MethodPost,
				Path:    "/get_conversation",
				Handler: rest.GetConversation,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
		}
		glue.RegisterRoutes(r, routes)

	})

}
