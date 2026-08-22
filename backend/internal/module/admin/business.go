package admin

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	adminPersist "adinfluence_service_app/internal/storage/persist/admin"
	"context"
	"errors"
	"strings"
	"time"

	"go.uber.org/zap"
)

func (m *adminModule) ListBusinessVerifications(ctx context.Context, query adminDto.ListBusinessVerificationsQuery) (*adminDto.ListBusinessVerificationsResponse, error) {
	m.logger.Info(ctx, "admin listing business verification requests", zap.String("status", query.Status))

	query = query.Normalize()
	items, total, err := m.admin.ListBusinessVerifications(ctx, query)
	if err != nil {
		m.logger.Error(ctx, "failed to query business verifications", zap.String("reason", err.Error()))
		return nil, err
	}

	return &adminDto.ListBusinessVerificationsResponse{
		Items:      items,
		Pagination: commonDto.NewListMeta(query.Page, query.Limit, total),
	}, nil
}

func (m *adminModule) GetBusinessVerification(ctx context.Context, userID string) (*adminDto.BusinessVerificationDetailResponse, error) {
	m.logger.Info(ctx, "admin fetching business verification details", zap.String("target_user_id", userID))

	row, err := m.admin.GetBusinessVerification(ctx, userID)
	if err != nil {
		m.logger.Warn(ctx, "business verification lookup failed", zap.String("target_user_id", userID), zap.String("reason", err.Error()))
		return nil, mapBusinessVerificationError(err)
	}

	return businessVerificationDetailFromRow(row), nil
}

func (m *adminModule) VerifyBusiness(ctx context.Context, adminUserID, targetUserID string, req adminDto.VerifyBusinessRequest, meta authDto.SessionMeta) (*adminDto.VerifyBusinessResponse, error) {
	m.logger.Info(ctx, "admin verifying business profile", zap.String("admin_id", adminUserID), zap.String("target_user_id", targetUserID))

	if err := m.admin.VerifyBusiness(ctx, targetUserID); err != nil {
		m.logger.Error(ctx, "failed to verify business in persistence", zap.String("target_user_id", targetUserID), zap.String("reason", err.Error()))
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

	m.logger.Info(ctx, "business verified successfully", zap.String("target_user_id", targetUserID))

	return &adminDto.VerifyBusinessResponse{
		UserID:     targetUserID,
		IsVerified: true,
		VerifiedAt: now,
	}, nil
}

func (m *adminModule) RevokeBusinessVerification(ctx context.Context, adminUserID, targetUserID string, req adminDto.RevokeBusinessVerificationRequest, meta authDto.SessionMeta) (*adminDto.RevokeBusinessVerificationResponse, error) {
	m.logger.Info(ctx, "admin revoking business verification", zap.String("admin_id", adminUserID), zap.String("target_user_id", targetUserID), zap.String("reason", req.Reason))

	if err := m.admin.RevokeBusinessVerification(ctx, targetUserID); err != nil {
		m.logger.Error(ctx, "failed to revoke business verification", zap.String("target_user_id", targetUserID), zap.String("reason", err.Error()))
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

	m.logger.Info(ctx, "business verification revoked successfully", zap.String("target_user_id", targetUserID))

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
