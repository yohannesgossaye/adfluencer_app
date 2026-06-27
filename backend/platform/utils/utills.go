package utils

import (
	"strings"

	"github.com/google/uuid"
)

func NormalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	phone = strings.TrimPrefix(phone, "+")
	phone = strings.ReplaceAll(phone, "-", "")

	switch {
	case strings.HasPrefix(phone, "251") && len(phone) >= 12:
		return phone[:12]

	case strings.HasPrefix(phone, "0") && len(phone) >= 10:
		return "251" + phone[1:10]

	case len(phone) == 9:
		return "251" + phone
	}

	return phone
}

func GenerateUUID() string {
	id := uuid.NewString()
	return id
}
