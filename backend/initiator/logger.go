package initiator

import (
	"log"
	"os"
	"strings"

	"adinfluence_service_app/platform/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() utils.Logger {
	format := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_FORMAT")))
	if format == "" {
		format = "console"
	}

	var cfg zap.Config
	switch format {
	case "json":
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "ts"
		cfg.EncoderConfig.MessageKey = "msg"
		cfg.EncoderConfig.LevelKey = "level"
	default:
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	cfg.EncoderConfig.CallerKey = ""
	cfg.EncoderConfig.StacktraceKey = ""

	if level := strings.TrimSpace(os.Getenv("LOG_LEVEL")); level != "" {
		if err := cfg.Level.UnmarshalText([]byte(level)); err != nil {
			log.Printf("invalid LOG_LEVEL %q, using default: %v", level, err)
		}
	}

	logger, err := cfg.Build(zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		log.Fatalf("failed to initialize zap logger: %v", err)
	}

	return utils.New(logger)
}
