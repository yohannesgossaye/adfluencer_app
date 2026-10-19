package initiator

import (
	platformStorage "adinfluence_service_app/platform/storage"
	"adinfluence_service_app/platform/utils"
	"context"
	"os"
	"strings"

	"go.uber.org/zap"
)

func initFileStorage(ctx context.Context, cfg S3Config, port string, logger utils.Logger) platformStorage.Service {
	storageCfg := platformStorage.Config{
		Endpoint:     cfg.Endpoint,
		Region:       cfg.Region,
		AccessKeyID:  cfg.AccessKeyID,
		SecretKey:    cfg.SecretKey,
		Bucket:       cfg.Bucket,
		PublicURL:    cfg.PublicURL,
		UsePathStyle: cfg.UsePathStyle,
	}
	if storageCfg.Enabled() {
		s3, err := platformStorage.NewS3(storageCfg)
		if err == nil {
			logger.Info(ctx, "s3 storage ready", zap.String("bucket", cfg.Bucket), zap.String("endpoint", cfg.Endpoint))
			return s3
		}
		logger.Warn(ctx, "s3 init failed — falling back to local disk storage", zap.String("reason", err.Error()))
	}

	// Automatic local disk storage fallback (stores in ./uploads and serves via static route).
	// On AWS EC2: set STORAGE_PUBLIC_URL=https://your-domain.com/uploads
	// On Render:  set STORAGE_PUBLIC_URL or let RENDER_EXTERNAL_URL be auto-detected.
	publicURL := os.Getenv("STORAGE_PUBLIC_URL")
	if publicURL == "" {
		if renderURL := os.Getenv("RENDER_EXTERNAL_URL"); renderURL != "" {
			publicURL = strings.TrimRight(renderURL, "/") + "/uploads"
		} else {
			// No public URL configured — local file URLs will be relative.
			// Set STORAGE_PUBLIC_URL in your environment to fix this.
			logger.Warn(ctx, "STORAGE_PUBLIC_URL not set — uploaded file URLs may not be accessible publicly; set STORAGE_PUBLIC_URL=https://your-domain.com/uploads")
			publicURL = ""
		}
	}


	localStorage, err := platformStorage.NewLocal("./uploads", publicURL)
	if err != nil {
		logger.Warn(ctx, "local storage init failed", zap.String("reason", err.Error()))
		return platformStorage.NewNoop()
	}

	logger.Info(ctx, "local file storage ready", zap.String("public_url", publicURL))
	return localStorage
}
