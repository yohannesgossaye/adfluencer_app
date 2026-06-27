package storage

import (
	"context"
	"errors"
	"io"
	"time"
)

type NoopStorage struct{}

func NewNoop() *NoopStorage {
	return &NoopStorage{}
}

func (n *NoopStorage) Enabled() bool {
	return false
}

func (n *NoopStorage) PublicURLBase() string {
	return ""
}

func (n *NoopStorage) UploadBytes(ctx context.Context, objectKey string, data io.Reader, size int64, contentType string) (string, error) {
	return "", errors.New("storage is not configured")
}

func (n *NoopStorage) Check(ctx context.Context, objectKey string) (bool, error) {
	return false, errors.New("storage is not configured")
}

func (n *NoopStorage) GetFileURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	return "", errors.New("storage is not configured")
}

func (n *NoopStorage) Delete(ctx context.Context, objectKey string) error {
	return errors.New("storage is not configured")
}
