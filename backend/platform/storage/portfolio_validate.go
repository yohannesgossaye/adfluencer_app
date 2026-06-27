package storage

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	MaxPortfolioImageSize = 10 << 20  // 10 MB
	MaxPortfolioVideoSize = 100 << 20 // 100 MB
)

type ValidatedPortfolioFile struct {
	ContentType string
	Extension   string
	MediaType   string // image | video
	Data        []byte
	Size        int64
}

func ValidatePortfolioFile(filename, contentType string, data []byte) (ValidatedPortfolioFile, error) {
	if len(data) == 0 {
		return ValidatedPortfolioFile{}, ErrInvalidFile
	}

	detected, ext, mediaType := detectPortfolioMedia(data)
	if detected == "" {
		return ValidatedPortfolioFile{}, ErrUnsupportedFormat
	}

	maxSize := MaxPortfolioImageSize
	if mediaType == "video" {
		maxSize = MaxPortfolioVideoSize
	}
	if int64(len(data)) > int64(maxSize) {
		return ValidatedPortfolioFile{}, fmt.Errorf("%w: file exceeds size limit", ErrInvalidFile)
	}

	ct := strings.TrimSpace(contentType)
	if ct != "" {
		prefix := "image/"
		if mediaType == "video" {
			prefix = "video/"
		}
		if !strings.HasPrefix(strings.ToLower(ct), prefix) {
			return ValidatedPortfolioFile{}, ErrUnsupportedFormat
		}
	}
	if ct == "" {
		ct = detected
	}

	_ = filename

	return ValidatedPortfolioFile{
		ContentType: ct,
		Extension:   ext,
		MediaType:   mediaType,
		Data:        data,
		Size:        int64(len(data)),
	}, nil
}

func detectPortfolioMedia(data []byte) (contentType, ext, mediaType string) {
	if ct, ext := detectImage(data); ct != "" {
		return ct, ext, "image"
	}
	if ct, ext := detectVideo(data); ct != "" {
		return ct, ext, "video"
	}
	return "", "", ""
}

func detectVideo(data []byte) (contentType, ext string) {
	if len(data) >= 12 && string(data[4:8]) == "ftyp" {
		return "video/mp4", "mp4"
	}
	if len(data) >= 4 && bytes.Equal(data[:4], []byte{0x1A, 0x45, 0xDF, 0xA3}) {
		return "video/webm", "webm"
	}
	return "", ""
}
