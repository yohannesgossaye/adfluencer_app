package marketplace

import (
	"adinfluence_service_app/internal/glue"
	"adinfluence_service_app/internal/handler"
	"adinfluence_service_app/internal/handler/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func MarketPlaceRoutes(router chi.Router, chatHandler handler.ChatHandler, authMiddleware middleware.AuthMiddleware) {
	router.Route("/chat", func(r chi.Router) {
		// REST endpoints (Protected with Bearer token)
		routes := []glue.Route{
			{
				Method:  http.MethodPost,
				Path:    "/conversations",
				Handler: chatHandler.GetOrCreateConversation,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/conversations",
				Handler: chatHandler.ListConversations,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/conversations/{id}",
				Handler: chatHandler.GetConversation,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/messages",
				Handler: chatHandler.SendMessage,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/messages",
				Handler: chatHandler.GetMessages,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/messages/read",
				Handler: chatHandler.MarkMessagesAsRead,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodDelete,
				Path:    "/messages/{id}",
				Handler: chatHandler.DeleteMessage,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				// WebSocket endpoint handles token validation internally via query param or header
				Method:  http.MethodGet,
				Path:    "/ws",
				Handler: chatHandler.WebSocket,
			},
		}
		glue.RegisterRoutes(r, routes)
	})
}
