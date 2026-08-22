package initiator

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"adinfluence_service_app/internal/handler/middleware"
	"adinfluence_service_app/platform/jwt"
	"adinfluence_service_app/platform/mail"
	"adinfluence_service_app/platform/sms"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Initiator(ctx context.Context) {
	logger := InitLogger()
	defer logger.Sync()

	cfg := Load()

	logger.Info(ctx, "initiating database connection ......")

	db, err := NewDB(cfg)
	if err != nil {
		logger.Fatal(ctx, "failed to connect database", zap.Error(err))
	}
	defer db.Close()

	logger.Info(ctx, "database connected successfully")

	emailClient := mail.New(mail.Config{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		User:     cfg.SMTP.User,
		Password: cfg.SMTP.Password,
		From:     cfg.SMTP.From,
	})
	if emailClient.Enabled() {
		logger.Info(ctx, "smtp email service ready", zap.String("host", cfg.SMTP.Host))
		if os.Getenv("SMTP_SEND_TEST") == "true" {
			if err := emailClient.Send(cfg.SMTP.User, "Adfluence SMTP test", "If you see this email, SMTP is working."); err != nil {
				logger.Warn(ctx, "smtp test email failed", zap.String("reason", err.Error()))
			} else {
				logger.Info(ctx, "smtp test email sent", zap.String("to", cfg.SMTP.User))
			}
		}
	} else {
		logger.Warn(ctx, "smtp not configured — email features disabled until SMTP_* env vars are set")
	}

	smsClient := sms.New(sms.Config{
		Username: cfg.SMS.Username,
		APIKey:   cfg.SMS.APIKey,
		From:     cfg.SMS.From,
	})
	if smsClient.Enabled() {
		logger.Info(ctx, "sms service ready (Africa's Talking)")
	} else {
		logger.Warn(ctx, "sms not configured — phone OTP logged to console until SMS_USERNAME + SMS_API_KEY are set")
	}

	fileStorage := initFileStorage(ctx, cfg.S3, logger)

	r := chi.NewRouter()

	postgresRepo := InitPersistence(db, logger)
	tokens := jwt.NewTokenManager(jwt.ManagerConfig{
		SecretKey:                 cfg.JWT.SecretKey,
		AccessTokenTTL:            cfg.JWT.AccessTokenTTL,
		RefreshTokenTTL:           cfg.JWT.RefreshTokenTTL,
		RememberMeRefreshTokenTTL: cfg.JWT.RememberMeRefreshTokenTTL,
	})
	service := InitModules(postgresRepo, tokens, emailClient, smsClient, fileStorage, cfg.FrontendURL, logger)
	handler := InitHandler(service, tokens, logger)
	authMW := middleware.InitAuthMiddleware(tokens)
	InirRoute(r, handler, authMW)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		logger.Info(ctx, "🚀 service started", zap.String("port", cfg.Port))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, "server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	sig := <-quit
	logger.Info(ctx, "shutdown signal received", zap.String("signal", sig.String()))

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(ctx, "server shutdown failed", zap.Error(err))
		return
	}

	logger.Info(ctx, "server stopped gracefully")
}
