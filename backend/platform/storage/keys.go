package storage

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func AvatarObjectKey(userID, ext string) string {
	ext = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(ext)), ".")
	return fmt.Sprintf("media/influencer/%s/avatar/%s.%s", userID, uuid.NewString(), ext)
}

func PortfolioObjectKey(userID, ext string) string {
	ext = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(ext)), ".")
	return fmt.Sprintf("media/influencer/%s/portfolio/%s.%s", userID, uuid.NewString(), ext)
}

func BusinessLogoObjectKey(userID, ext string) string {
	ext = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(ext)), ".")
	return fmt.Sprintf("media/business/%s/logo/%s.%s", userID, uuid.NewString(), ext)
}
