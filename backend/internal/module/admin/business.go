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

	detail := row.BusinessVerificationDetailResponse
	return &detail, nil
}

func (m *adminModule) VerifyBusiness(ctx context.Context, adminUserID, targetUserID string, req adminDto.VerifyBusinessRequest, meta authDto.SessionMeta) (*adminDto.VerifyBusinessResponse, error) {
	m.logger.Info(ctx, "admin verifying business profile", zap.String("admin_id", adminUserID), zap.String("target_user_id", targetUserID))

	note := strings.TrimSpace(req.AdminNote)
	verifiedAt, err := m.admin.VerifyBusiness(ctx, targetUserID, adminUserID, note)
	if err != nil {
		m.logger.Warn(ctx, "business verification not applied", zap.String("target_user_id", targetUserID), zap.String("reason", err.Error()))
		return nil, mapBusinessVerificationError(err)
	}

	desc := "business verification approved"
	if note != "" {
		desc = desc + ": " + note
	}
	m.auditBusinessDecision(adminUserID, targetUserID, model.ApproveAction, desc, meta)

	return &adminDto.VerifyBusinessResponse{
		UserID:             targetUserID,
		IsVerified:         true,
		VerificationStatus: "verified",
		VerifiedAt:         verifiedAt,
	}, nil
}

// RejectBusinessVerification declines a pending business. The reason is shown
// to the business so it can fix its profile and submit again.
func (m *adminModule) RejectBusinessVerification(ctx context.Context, adminUserID, targetUserID string, req adminDto.RejectBusinessVerificationRequest, meta authDto.SessionMeta) (*adminDto.BusinessVerificationDecisionResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	reason := strings.TrimSpace(req.Reason)
	if err := m.admin.RejectBusinessVerification(ctx, targetUserID, adminUserID, reason); err != nil {
		m.logger.Warn(ctx, "business rejection not applied", zap.String("target_user_id", targetUserID), zap.String("reason", err.Error()))
		return nil, mapBusinessVerificationError(err)
	}
	m.auditBusinessDecision(adminUserID, targetUserID, model.RejectAction, "business verification rejected: "+reason, meta)
	return &adminDto.BusinessVerificationDecisionResponse{UserID: targetUserID, VerificationStatus: "rejected", Reason: reason}, nil
}

// RevokeBusinessVerification removes verification from a verified business,
// which immediately locks publishing, invites and contracts.
func (m *adminModule) RevokeBusinessVerification(ctx context.Context, adminUserID, targetUserID string, req adminDto.RevokeBusinessVerificationRequest, meta authDto.SessionMeta) (*adminDto.RevokeBusinessVerificationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	reason := strings.TrimSpace(req.Reason)
	if err := m.admin.RevokeBusinessVerification(ctx, targetUserID, adminUserID, reason); err != nil {
		m.logger.Warn(ctx, "business revocation not applied", zap.String("target_user_id", targetUserID), zap.String("reason", err.Error()))
		return nil, mapBusinessVerificationError(err)
	}
	m.auditBusinessDecision(adminUserID, targetUserID, model.RejectAction, "business verification revoked: "+reason, meta)
	return &adminDto.RevokeBusinessVerificationResponse{UserID: targetUserID, VerificationStatus: "rejected", Reason: reason}, nil
}

func (m *adminModule) auditBusinessDecision(adminUserID, targetUserID string, action model.AuditAction, desc string, meta authDto.SessionMeta) {
	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &adminUserID,
		TargetUserID: &targetUserID,
		Action:       string(action),
		EntityType:   "business_profile",
		EntityID:     &targetUserID,
		Description:  desc,
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})
}

func mapBusinessVerificationError(err error) error {
	switch {
	case errors.Is(err, adminPersist.ErrBusinessVerificationNotFound):
		return customErr.ErrBusinessNotFound
	case errors.Is(err, adminPersist.ErrBusinessAlreadyVerified):
		return customErr.ErrBusinessAlreadyVerified
	case errors.Is(err, adminPersist.ErrBusinessNotVerified):
		return customErr.ErrBusinessNotVerified
	case errors.Is(err, adminPersist.ErrBusinessNotPendingReview):
		return customErr.ErrBusinessNotPendingReview
	}
	return err
}
