package influencer

import (
	"adinfluence_service_app/internal/glue"
	"adinfluence_service_app/internal/handler"
	"adinfluence_service_app/internal/handler/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func InfluencerRoutes(router chi.Router, rest handler.InfluencerHandler, authMiddleware middleware.AuthMiddleware) {
	router.Route("/influencer", func(r chi.Router) {
		routes := []glue.Route{
			{
				Method:  http.MethodGet,
				Path:    "/dashboard/summary",
				Handler: rest.GetDashboardSummary,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/profile",
				Handler: rest.GetProfile,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPut,
				Path:    "/profile",
				Handler: rest.UpdateProfile,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/profile/submit",
				Handler: rest.SubmitProfile,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/profile/avatar",
				Handler: rest.UploadAvatar,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/profile/portfolio",
				Handler: rest.UploadPortfolio,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodDelete,
				Path:    "/profile/portfolio/{item_id}",
				Handler: rest.DeletePortfolio,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/invites",
				Handler: rest.ListInvites,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/invites/{id}",
				Handler: rest.GetInvite,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/invites/{id}/accept",
				Handler: rest.AcceptInvite,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/invites/{id}/reject",
				Handler: rest.RejectInvite,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/contracts",
				Handler: rest.ListContracts,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/contracts/{contract_id}",
				Handler: rest.GetContract,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/contracts/{contract_id}/accept",
				Handler: rest.AcceptContract,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/contracts/{contract_id}/reject",
				Handler: rest.RejectContract,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/support/invite-options",
				Handler: rest.ListSupportInviteOptions,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/support/tickets",
				Handler: rest.CreateSupportTicket,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/support/tickets",
				Handler: rest.ListSupportTickets,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/support/tickets/{ticket_id}",
				Handler: rest.GetSupportTicket,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
		}
		glue.RegisterRoutes(r, routes)
	})
}
