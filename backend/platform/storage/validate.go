package storage

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

const MaxAvatarSize = 5 << 20 // 5 MB

var (
	ErrInvalidFile       = errors.New("invalid file")
	ErrUnsupportedFormat = errors.New("unsupported file format")
)

type ValidatedImage struct {
	ContentType string
	Extension   string
	Data        []byte
	Size        int64
}

func ValidateAvatarFile(filename, contentType string, data []byte) (ValidatedImage, error) {
	if len(data) == 0 {
		return ValidatedImage{}, ErrInvalidFile
	}
	if int64(len(data)) > MaxAvatarSize {
		return ValidatedImage{}, fmt.Errorf("%w: file exceeds 5 MB", ErrInvalidFile)
	}

	detected, ext := detectImage(data)
	if detected == "" {
		return ValidatedImage{}, ErrUnsupportedFormat
	}

	ct := strings.TrimSpace(contentType)
	if ct != "" && !strings.HasPrefix(strings.ToLower(ct), "image/") {
		return ValidatedImage{}, ErrUnsupportedFormat
	}
	if ct == "" {
		ct = detected
	}

	_ = filename

	return ValidatedImage{
		ContentType: detected,
		Extension:   ext,
		Data:        data,
		Size:        int64(len(data)),
	}, nil
}

func detectImage(data []byte) (contentType, ext string) {
	switch {
	case len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF:
		return "image/jpeg", "jpg"
	case len(data) >= 8 && bytes.Equal(data[:8], []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}):
		return "image/png", "png"
	case len(data) >= 12 && string(data[:4]) == "RIFF" && string(data[8:12]) == "WEBP":
		return "image/webp", "webp"
	default:
		return "", ""
	}
}

func ObjectKeyFromURL(publicURL, publicBase string) string {
	publicURL = strings.TrimSpace(publicURL)
	publicBase = strings.TrimRight(strings.TrimSpace(publicBase), "/")
	if publicURL == "" || publicBase == "" {
		return ""
	}
	if !strings.HasPrefix(publicURL, publicBase) {
		return ""
	}
	key := strings.TrimPrefix(publicURL, publicBase)
	return strings.TrimLeft(key, "/")
}
