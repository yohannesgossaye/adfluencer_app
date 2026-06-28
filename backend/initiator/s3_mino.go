package initiator

import (
	platformStorage "adinfluence_service_app/platform/storage"
	"adinfluence_service_app/platform/utils"
	"context"

	"go.uber.org/zap"
)

func initFileStorage(ctx context.Context, cfg S3Config, logger utils.Logger) platformStorage.Service {
	storageCfg := platformStorage.Config{
		Endpoint:     cfg.Endpoint,
		Region:       cfg.Region,
		AccessKeyID:  cfg.AccessKeyID,
		SecretKey:    cfg.SecretKey,
		Bucket:       cfg.Bucket,
		PublicURL:    cfg.PublicURL,
		UsePathStyle: cfg.UsePathStyle,
	}
	if !storageCfg.Enabled() {
		logger.Warn(ctx, "s3 not configured — avatar upload disabled until S3_* env vars are set")
		return platformStorage.NewNoop()
	}

	s3, err := platformStorage.NewS3(storageCfg)
	if err != nil {
		logger.Warn(ctx, "s3 init failed — avatar upload disabled", zap.String("reason", err.Error()))
		return platformStorage.NewNoop()
	}

	logger.Info(ctx, "s3 storage ready", zap.String("bucket", cfg.Bucket), zap.String("endpoint", cfg.Endpoint))
	return s3
}
