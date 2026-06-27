package admin

import (
	"adinfluence_service_app/internal/glue"
	"adinfluence_service_app/internal/handler"
	"adinfluence_service_app/internal/handler/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func AdminRoutes(router chi.Router, rest handler.AdminHandler, authMiddleware middleware.AuthMiddleware) {
	router.Route("/admin", func(r chi.Router) {
		auth := []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken}
		glue.RegisterRoutes(r, []glue.Route{
			{
				Method:      http.MethodGet,
				Path:        "/dashboard/summary",
				Handler:     rest.GetDashboardSummary,
				Middlewares: auth,
			},
		})
		r.Route("/influencers", func(r chi.Router) {
			auth := []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken}
			routes := []glue.Route{
				{
					Method:      http.MethodGet,
					Path:        "/applications",
					Handler:     rest.ListApplications,
					Middlewares: auth,
				},
				{
					Method:      http.MethodGet,
					Path:        "/applications/{user_id}",
					Handler:     rest.GetApplication,
					Middlewares: auth,
				},
				{
					Method:      http.MethodPost,
					Path:        "/applications/{user_id}/approve",
					Handler:     rest.ApproveApplication,
					Middlewares: auth,
				},
				{
					Method:      http.MethodPost,
					Path:        "/applications/{user_id}/reject",
					Handler:     rest.RejectApplication,
					Middlewares: auth,
				},
				{
					Method:      http.MethodPatch,
					Path:        "/applications/{user_id}/notes",
					Handler:     rest.UpdateAdminNotes,
					Middlewares: auth,
				},
			}
			glue.RegisterRoutes(r, routes)
		})
		r.Route("/businesses", func(r chi.Router) {
			auth := []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken}
			routes := []glue.Route{
				{
					Method:      http.MethodGet,
					Path:        "/verifications",
					Handler:     rest.ListBusinessVerifications,
					Middlewares: auth,
				},
				{
					Method:      http.MethodGet,
					Path:        "/verifications/{user_id}",
					Handler:     rest.GetBusinessVerification,
					Middlewares: auth,
				},
				{
					Method:      http.MethodPost,
					Path:        "/verifications/{user_id}/verify",
					Handler:     rest.VerifyBusiness,
					Middlewares: auth,
				},
				{
					Method:      http.MethodPost,
					Path:        "/verifications/{user_id}/revoke",
					Handler:     rest.RevokeBusinessVerification,
					Middlewares: auth,
				},
			}
			glue.RegisterRoutes(r, routes)
		})
		r.Route("/support", func(r chi.Router) {
			auth := []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken}
			routes := []glue.Route{
				{
					Method:      http.MethodGet,
					Path:        "/tickets",
					Handler:     rest.ListSupportTickets,
					Middlewares: auth,
				},
				{
					Method:      http.MethodGet,
					Path:        "/tickets/{ticket_id}",
					Handler:     rest.GetSupportTicket,
					Middlewares: auth,
				},
				{
					Method:      http.MethodPatch,
					Path:        "/tickets/{ticket_id}/status",
					Handler:     rest.UpdateSupportTicketStatus,
					Middlewares: auth,
				},
				{
					Method:      http.MethodPost,
					Path:        "/tickets/{ticket_id}/reply",
					Handler:     rest.ReplySupportTicket,
					Middlewares: auth,
				},
			}
			glue.RegisterRoutes(r, routes)
		})
		r.Route("/categories", func(r chi.Router) {
			auth := []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken}
			routes := []glue.Route{
				{
					Method:      http.MethodGet,
					Path:        "/",
					Handler:     rest.ListCategories,
					Middlewares: auth,
				},
				{
					Method:      http.MethodPost,
					Path:        "/",
					Handler:     rest.CreateCategory,
					Middlewares: auth,
				},
				{
					Method:      http.MethodPatch,
					Path:        "/{category_id}",
					Handler:     rest.UpdateCategory,
					Middlewares: auth,
				},
				{
					Method:      http.MethodPost,
					Path:        "/{category_id}/disable",
					Handler:     rest.DisableCategory,
					Middlewares: auth,
				},
				{
					Method:      http.MethodPost,
					Path:        "/{category_id}/enable",
					Handler:     rest.EnableCategory,
					Middlewares: auth,
				},
			}
			glue.RegisterRoutes(r, routes)
		})
		r.Route("/users", func(r chi.Router) {
			auth := []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken}
			routes := []glue.Route{
				{
					Method:      http.MethodGet,
					Path:        "/",
					Handler:     rest.ListUsers,
					Middlewares: auth,
				},
				{
					Method:      http.MethodGet,
					Path:        "/{user_id}",
					Handler:     rest.GetUser,
					Middlewares: auth,
				},
				{
					Method:      http.MethodPost,
					Path:        "/{user_id}/block",
					Handler:     rest.BlockUser,
					Middlewares: auth,
				},
				{
					Method:      http.MethodPost,
					Path:        "/{user_id}/unblock",
					Handler:     rest.UnblockUser,
					Middlewares: auth,
				},
				{
					Method:      http.MethodPost,
					Path:        "/{user_id}/remove",
					Handler:     rest.RemoveUser,
					Middlewares: auth,
				},
			}
			glue.RegisterRoutes(r, routes)
		})
		r.Route("/campaigns", func(r chi.Router) {
			auth := []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken}
			routes := []glue.Route{
				{
					Method:      http.MethodGet,
					Path:        "/",
					Handler:     rest.ListAdminCampaigns,
					Middlewares: auth,
				},
				{
					Method:      http.MethodGet,
					Path:        "/{campaign_id}",
					Handler:     rest.GetAdminCampaign,
					Middlewares: auth,
				},
				{
					Method:      http.MethodPost,
					Path:        "/{campaign_id}/pause",
					Handler:     rest.PauseAdminCampaign,
					Middlewares: auth,
				},
				{
					Method:      http.MethodPost,
					Path:        "/{campaign_id}/complete",
					Handler:     rest.CompleteAdminCampaign,
					Middlewares: auth,
				},
				{
					Method:      http.MethodPost,
					Path:        "/{campaign_id}/flag",
					Handler:     rest.FlagAdminCampaign,
					Middlewares: auth,
				},
			}
			glue.RegisterRoutes(r, routes)
		})
		r.Route("/disputes", func(r chi.Router) {
			auth := []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken}
			routes := []glue.Route{
				{
					Method:      http.MethodGet,
					Path:        "/",
					Handler:     rest.ListAdminDisputes,
					Middlewares: auth,
				},
				{
					Method:      http.MethodGet,
					Path:        "/{dispute_id}",
					Handler:     rest.GetAdminDispute,
					Middlewares: auth,
				},
				{
					Method:      http.MethodPatch,
					Path:        "/{dispute_id}/status",
					Handler:     rest.UpdateDisputeStatus,
					Middlewares: auth,
				},
				{
					Method:      http.MethodPatch,
					Path:        "/{dispute_id}/notes",
					Handler:     rest.UpdateDisputeNotes,
					Middlewares: auth,
				},
			}
			glue.RegisterRoutes(r, routes)
		})
	})
}
