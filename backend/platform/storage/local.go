package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LocalStorage struct {
	baseDir   string
	publicURL string
}

func NewLocal(baseDir, publicURL string) (*LocalStorage, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("create storage directory: %w", err)
	}
	return &LocalStorage{
		baseDir:   baseDir,
		publicURL: strings.TrimRight(publicURL, "/"),
	}, nil
}

func (l *LocalStorage) Enabled() bool {
	return true
}

func (l *LocalStorage) PublicURLBase() string {
	return l.publicURL
}

func (l *LocalStorage) UploadBytes(ctx context.Context, objectKey string, data io.Reader, size int64, contentType string) (string, error) {
	objectKey = strings.TrimLeft(strings.TrimSpace(objectKey), "/")
	fullPath := filepath.Join(l.baseDir, filepath.FromSlash(objectKey))

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", fmt.Errorf("create parent directory: %w", err)
	}

	out, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("create destination file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, data); err != nil {
		return "", fmt.Errorf("save file data: %w", err)
	}

	return objectKey, nil
}

func (l *LocalStorage) Check(ctx context.Context, objectKey string) (bool, error) {
	objectKey = strings.TrimLeft(strings.TrimSpace(objectKey), "/")
	fullPath := filepath.Join(l.baseDir, filepath.FromSlash(objectKey))
	_, err := os.Stat(fullPath)
	return err == nil, nil
}

func (l *LocalStorage) GetFileURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	objectKey = strings.TrimLeft(strings.TrimSpace(objectKey), "/")
	escaped := strings.ReplaceAll(objectKey, " ", "%20")
	return url.JoinPath(l.publicURL, escaped)
}

func (l *LocalStorage) Delete(ctx context.Context, objectKey string) error {
	objectKey = strings.TrimLeft(strings.TrimSpace(objectKey), "/")
	fullPath := filepath.Join(l.baseDir, filepath.FromSlash(objectKey))
	_ = os.Remove(fullPath)
	return nil
}
