package utils

import (
	"context"

	"go.uber.org/zap"
)

type ContextKey string

const (
	UserIDKey    ContextKey = "user_id"
	RequestIDKey ContextKey = "request_id"
)

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Fatal(ctx context.Context, msg string, fields ...zap.Field)
	Named(name string) Logger
	With(fields ...zap.Field) Logger
	Sync() error
}

type ZapLogger struct {
	log *zap.Logger
}

func New(log *zap.Logger) Logger {
	return &ZapLogger{
		log: log,
	}
}

func (l *ZapLogger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.log.With(contextFields(ctx)...).Debug(msg, fields...)
}

func (l *ZapLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.log.With(contextFields(ctx)...).Info(msg, fields...)
}

func (l *ZapLogger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.log.With(contextFields(ctx)...).Warn(msg, fields...)
}

func (l *ZapLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.log.With(contextFields(ctx)...).Error(msg, fields...)
}

func (l *ZapLogger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.log.With(contextFields(ctx)...).Fatal(msg, fields...)
}

func (l *ZapLogger) Named(name string) Logger {
	return &ZapLogger{
		log: l.log.Named(name),
	}
}

func (l *ZapLogger) With(fields ...zap.Field) Logger {
	return &ZapLogger{
		log: l.log.With(fields...),
	}
}

func (l *ZapLogger) Sync() error {
	return l.log.Sync()
}

func contextFields(ctx context.Context) []zap.Field {
	if ctx == nil {
		return nil
	}

	fields := make([]zap.Field, 0, 2)

	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		fields = append(fields, zap.String("user_id", userID))
	}

	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		fields = append(fields, zap.String("request_id", requestID))
	}

	return fields
}
