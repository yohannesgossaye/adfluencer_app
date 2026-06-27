package auth

import (
	"adinfluence_service_app/internal/glue"
	"adinfluence_service_app/internal/handler"
	"adinfluence_service_app/internal/handler/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(router chi.Router, rest handler.AuthHandler, authMiddleware middleware.AuthMiddleware) {
	router.Route("/auth", func(r chi.Router) {
		routes := []glue.Route{
			{
				Method:  http.MethodPost,
				Path:    "/signup",
				Handler: rest.SignUp,
			},
			{
				Method:  http.MethodPost,
				Path:    "/login",
				Handler: rest.Login,
			},
			{
				Method:  http.MethodPost,
				Path:    "/refresh",
				Handler: rest.RefreshToken,
			},
			{
				Method:  http.MethodPost,
				Path:    "/forgot-password",
				Handler: rest.ForgotPassword,
			},
			{
				Method:  http.MethodPost,
				Path:    "/reset-password",
				Handler: rest.ResetPassword,
			},
			{
				Method:  http.MethodGet,
				Path:    "/me",
				Handler: rest.GetMe,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/logout",
				Handler: rest.Logout,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/send-verify-email",
				Handler: rest.SendVerifyEmail,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/verify-email",
				Handler: rest.VerifyEmail,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/send-verify-phone",
				Handler: rest.SendVerifyPhone,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/verify-phone",
				Handler: rest.VerifyPhone,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
		}
		glue.RegisterRoutes(r, routes)

	})

}
