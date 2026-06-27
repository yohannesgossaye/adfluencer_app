package helper

import (
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"adinfluence_service_app/platform/utils"
	"strings"
)

func MapperForSignUp(id string, req authDto.SignUpRequest, passwordHash string) model.User {
	return model.User{
		ID:            id,
		Email:         NormalizeEmail(req.Email),
		PhoneNumber:   utils.NormalizePhone(req.Phone),
		PasswordHash:  passwordHash,
		Role:          model.Role(req.Role),
		Status:        model.PendingProfileStatus,
		EmailVerified: false,
		PhoneVerified: false,
	}
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func RedirectToLogin(role, status string) string {
	if status == string(model.RejectedStatus) {
		return "/auth/account-status"
	}

	switch model.Role(role) {
	case model.AdminRole:
		return "/admin/dashboard"
	case model.InfluencerRole:
		if status == string(model.PendingProfileStatus) {
			return "/influencer/profile"
		}
		return "/influencer/dashboard"
	default:
		return "/business/dashboard"
	}
}
