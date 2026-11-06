package helper

import (
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"adinfluence_service_app/platform/jwt"
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

// SafeInternalPath returns next only when it is a safe relative in-app path.
func SafeInternalPath(next string) string {
	next = strings.TrimSpace(next)
	if next == "" {
		return ""
	}
	if !strings.HasPrefix(next, "/") || strings.HasPrefix(next, "//") {
		return ""
	}
	if strings.Contains(next, "://") {
		return ""
	}
	return next
}

func ProfileName(profile authDto.ProfileResponse, role string) string {
	if role == string(model.BusinessRole) {
		return profile.BusinessName
	}
	if role == string(model.AdminRole) {
		return profile.FullName
	}
	return profile.CreatorName
}

func AccessTokenInput(user authDto.UserResponse, profile authDto.ProfileResponse) jwt.AccessTokenInput {
	return jwt.AccessTokenInput{
		UserID:      user.ID,
		Email:       user.Email,
		Role:        user.Role,
		Status:      user.Status,
		ProfileID:   profile.ID,
		ProfileName: ProfileName(profile, user.Role),
	}
}

func TokenResponse(pair jwt.TokenPair) authDto.TokenResponse {
	return authDto.TokenResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresIn:    pair.AccessExpiresIn,
	}
}

func RedirectTo(role string) string {
	if role == string(model.InfluencerRole) {
		return "/influencer/dashboard"
	}
	if role == string(model.AdminRole) {
		return "/admin/dashboard"
	}
	return "/business/dashboard"
}
