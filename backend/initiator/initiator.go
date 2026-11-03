package initiator

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"adinfluence_service_app/internal/constant/model/response"
	"adinfluence_service_app/internal/handler/middleware"
	"adinfluence_service_app/migration"
	"adinfluence_service_app/platform/cleanup"
	"adinfluence_service_app/platform/jwt"
	"adinfluence_service_app/platform/mail"
	"adinfluence_service_app/platform/notify"
	"adinfluence_service_app/platform/sms"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

func Initiator(ctx context.Context) {
	logger := InitLogger()
	defer logger.Sync()

	cfg := Load()

	response.OnHiddenError = func(status int, code string, err error) {
		logger.Error(context.Background(), "error details withheld from client",
			zap.Int("status", status), zap.String("code", code), zap.Error(err))
	}

	logger.Info(ctx, "initiating database connection ......")

	db, err := NewDB(cfg)
	if err != nil {
		logger.Fatal(ctx, "failed to connect database", zap.Error(err))
	}
	defer db.Close()

	logger.Info(ctx, "database connected successfully")

	if cfg.MigrateOnStart {
		migrateCtx, cancelMigrate := context.WithTimeout(ctx, 5*time.Minute)
		version, err := migration.Up(migrateCtx, db)
		cancelMigrate()
		if err != nil {
			logger.Fatal(ctx, "database migration failed", zap.Error(err))
		}
		logger.Info(ctx, "database schema up to date", zap.Int64("version", version))
	}

	bgCtx, stopBackground := context.WithCancel(context.Background())
	defer stopBackground()
	cleanup.Start(bgCtx, db, logger, time.Hour)

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
		logger.Warn(ctx, "sms not configured — phone verification codes are not delivered until SMS_USERNAME + SMS_API_KEY are set")
	}

	fileStorage := initFileStorage(ctx, cfg.S3, cfg.Port, logger)

	r := chi.NewRouter()

	// Only origins listed in CORS_ORIGINS (comma-separated, "*" wildcards allowed
	// within a host, e.g. https://*.example.com) may call the API from a browser.
	logger.Info(ctx, "cors allowed origins", zap.Strings("origins", cfg.CORSOrigins))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin"},
		ExposedHeaders:   []string{"Link", "Content-Length"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(requestTimeout(requestTimeoutDefault, requestTimeoutUpload))

	// Serve uploaded files statically (files only — directory listings are disabled)
	workDir, _ := os.Getwd()
	uploadsDir := filesOnlyFS{http.Dir(filepath.Join(workDir, "uploads"))}
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(uploadsDir)))

	postgresRepo := InitPersistence(db, logger)
	tokens := jwt.NewTokenManager(jwt.ManagerConfig{
		SecretKey:                 cfg.JWT.SecretKey,
		AccessTokenTTL:            cfg.JWT.AccessTokenTTL,
		RefreshTokenTTL:           cfg.JWT.RefreshTokenTTL,
		RememberMeRefreshTokenTTL: cfg.JWT.RememberMeRefreshTokenTTL,
	})
	dispatcher := notify.New(logger)
	service := InitModules(postgresRepo, tokens, emailClient, smsClient, dispatcher, fileStorage, cfg.FrontendURL, cfg.ContactInbox, cfg.IsDevelopment(), logger)
	handler := InitHandler(service, tokens, cfg.CORSOrigins, logger)
	authMW := middleware.InitAuthMiddleware(tokens)
	InitRoute(r, handler, authMW)

	// No server-wide ReadTimeout/WriteTimeout: they would also cut long-lived
	// chat WebSockets and slow uploads. Per-request limits live in requestTimeout.
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	// WebSockets are hijacked connections that Shutdown does not wait for;
	// tell chat clients to reconnect instead of dropping them silently.
	server.RegisterOnShutdown(service.chat.GetHub().CloseAll)

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
	}
	// Let queued verification emails go out before the process exits.
	dispatcher.Close(shutdownCtx)

	logger.Info(ctx, "server stopped gracefully")
}
