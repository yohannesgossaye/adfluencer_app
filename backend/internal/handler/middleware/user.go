package middleware

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"context"
)

type UserPayload struct {
	UserID      string
	Email       string
	Role        string
	Status      string
	ProfileID   string
	ProfileName string
}

func ExtractUserInfo(ctx context.Context) (*UserPayload, error) {
	claims, ok := GetAccessClaims(ctx)
	if !ok {
		return nil, customErr.ErrUnauthorized
	}

	userID := claims.UserID
	if userID == "" {
		userID = claims.Subject
	}

	return &UserPayload{
		UserID:      userID,
		Email:       claims.Email,
		Role:        claims.Role,
		Status:      claims.Status,
		ProfileID:   claims.ProfileID,
		ProfileName: claims.ProfileName,
	}, nil
}

func (p *UserPayload) ToUserResponse() authDto.UserResponse {
	return authDto.UserResponse{
		ID:     p.UserID,
		Email:  p.Email,
		Role:   p.Role,
		Status: p.Status,
	}
}
