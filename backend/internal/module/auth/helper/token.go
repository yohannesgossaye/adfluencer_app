package helper

import (
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"adinfluence_service_app/platform/jwt"
)

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
