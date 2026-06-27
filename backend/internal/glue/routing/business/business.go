package business

import (
	"adinfluence_service_app/internal/glue"
	"adinfluence_service_app/internal/handler"
	"adinfluence_service_app/internal/handler/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func BusinessRoutes(router chi.Router, rest handler.BusinessHandler, authMiddleware middleware.AuthMiddleware) {
	router.Route("/business", func(r chi.Router) {
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
				Path:    "/profile/logo",
				Handler: rest.UploadLogo,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/profile/complete",
				Handler: rest.CompleteProfile,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/campaigns",
				Handler: rest.CreateCampaign,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/campaigns",
				Handler: rest.ListCampaigns,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/campaigns/{campaign_id}",
				Handler: rest.GetCampaign,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPut,
				Path:    "/campaigns/{campaign_id}",
				Handler: rest.UpdateCampaign,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/campaigns/{campaign_id}/publish",
				Handler: rest.PublishCampaign,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/campaigns/{campaign_id}/cancel",
				Handler: rest.CancelCampaign,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/campaigns/{campaign_id}/recommendations/generate",
				Handler: rest.GenerateRecommendations,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/campaigns/{campaign_id}/recommendations",
				Handler: rest.ListRecommendations,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/influencers/{user_id}",
				Handler: rest.GetInfluencerForBusiness,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/campaigns/{campaign_id}/invites",
				Handler: rest.CreateInvite,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/campaigns/{campaign_id}/invites",
				Handler: rest.ListInvites,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodDelete,
				Path:    "/campaigns/{campaign_id}/invites/{invite_id}",
				Handler: rest.CancelInvite,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/contracts",
				Handler: rest.CreateContract,
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
				Path:    "/contracts/{contract_id}/pay",
				Handler: rest.PayContract,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/contracts/{contract_id}/approve-work",
				Handler: rest.ApproveWork,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/contracts/{contract_id}/cancel",
				Handler: rest.CancelContract,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/support/campaign-options",
				Handler: rest.ListSupportCampaignOptions,
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
