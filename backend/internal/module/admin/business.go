package admin

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	adminPersist "adinfluence_service_app/internal/storage/persist/admin"
	"context"
	"errors"
	"strings"
	"time"
)

func (m *adminModule) ListBusinessVerifications(ctx context.Context, query adminDto.ListBusinessVerificationsQuery) (*adminDto.ListBusinessVerificationsResponse, error) {
	query = query.Normalize()
	items, total, err := m.admin.ListBusinessVerifications(ctx, query)
	if err != nil {
		return nil, err
	}

	return &adminDto.ListBusinessVerificationsResponse{
		Items: items,
		Pagination: commonDto.NewListMeta(query.Page, query.Limit, total),
	}, nil
}

func (m *adminModule) GetBusinessVerification(ctx context.Context, userID string) (*adminDto.BusinessVerificationDetailResponse, error) {
	row, err := m.admin.GetBusinessVerification(ctx, userID)
	if err != nil {
		return nil, mapBusinessVerificationError(err)
	}

	return businessVerificationDetailFromRow(row), nil
}

func (m *adminModule) VerifyBusiness(ctx context.Context, adminUserID, targetUserID string, req adminDto.VerifyBusinessRequest, meta authDto.SessionMeta) (*adminDto.VerifyBusinessResponse, error) {
	if err := m.admin.VerifyBusiness(ctx, targetUserID); err != nil {
		return nil, mapBusinessVerificationError(err)
	}

	now := time.Now().UTC()
	desc := "business marketplace verification approved"
	if note := strings.TrimSpace(req.AdminNote); note != "" {
		desc = desc + ": " + note
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &adminUserID,
		TargetUserID: &targetUserID,
		Action:       string(model.ApproveAction),
		EntityType:   "business_profile",
		EntityID:     &targetUserID,
		Description:  desc,
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	return &adminDto.VerifyBusinessResponse{
		UserID:     targetUserID,
		IsVerified: true,
		VerifiedAt: now,
	}, nil
}

func (m *adminModule) RevokeBusinessVerification(ctx context.Context, adminUserID, targetUserID string, req adminDto.RevokeBusinessVerificationRequest, meta authDto.SessionMeta) (*adminDto.RevokeBusinessVerificationResponse, error) {
	if err := m.admin.RevokeBusinessVerification(ctx, targetUserID); err != nil {
		return nil, mapBusinessVerificationError(err)
	}

	desc := "business marketplace verification revoked"
	if reason := strings.TrimSpace(req.Reason); reason != "" {
		desc = desc + ": " + reason
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &adminUserID,
		TargetUserID: &targetUserID,
		Action:       string(model.RejectAction),
		EntityType:   "business_profile",
		EntityID:     &targetUserID,
		Description:  desc,
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	return &adminDto.RevokeBusinessVerificationResponse{
		UserID:     targetUserID,
		IsVerified: false,
	}, nil
}

func businessVerificationDetailFromRow(row *adminDto.BusinessVerificationRow) *adminDto.BusinessVerificationDetailResponse {
	industry := ""
	if row.Industry != nil {
		industry = *row.Industry
	}
	city := ""
	if row.City != nil {
		city = *row.City
	}

	return &adminDto.BusinessVerificationDetailResponse{
		UserID:        row.UserID,
		Email:         row.Email,
		PhoneNumber:   row.PhoneNumber,
		Status:        row.UserStatus,
		BusinessName:  row.BusinessName,
		OwnerFullName: row.OwnerFullName,
		Industry:      industry,
		City:          city,
		Description:   row.Description,
		LogoURL:       row.LogoURL,
		IsVerified:    row.IsVerified,
		RegisteredAt:  row.RegisteredAt,
	}
}

func mapBusinessVerificationError(err error) error {
	if errors.Is(err, adminPersist.ErrBusinessVerificationNotFound) {
		return customErr.ErrBusinessNotFound
	}
	if errors.Is(err, adminPersist.ErrBusinessAlreadyVerified) {
		return customErr.ErrBusinessAlreadyVerified
	}
	if errors.Is(err, adminPersist.ErrBusinessNotVerified) {
		return customErr.ErrBusinessNotVerified
	}
	return err
}
