package storage

import (
	"context"
	"io"
	"time"
)

type Service interface {
	Enabled() bool
	UploadBytes(ctx context.Context, objectKey string, data io.Reader, size int64, contentType string) (string, error)
	Check(ctx context.Context, objectKey string) (bool, error)
	GetFileURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error)
	Delete(ctx context.Context, objectKey string) error
	PublicURLBase() string
}
