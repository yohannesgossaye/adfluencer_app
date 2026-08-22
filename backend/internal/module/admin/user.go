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
	"golang.org/x/sync/errgroup"
)

func (m *adminModule) ListUsers(ctx context.Context, query adminDto.ListUsersQuery) (*adminDto.ListUsersResponse, error) {
	m.logger.Info(ctx, "admin listing users", zap.String("role", query.Role), zap.String("status", query.Status))

	query = query.Normalize()

	var (
		rows  []adminDto.UserRow
		total int
		stats *adminDto.UserStats
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		rows, total, err = m.admin.ListUsers(ctx, query)
		return err
	})

	g.Go(func() error {
		var err error
		stats, err = m.admin.GetUserStats(ctx)
		return err
	})

	if err := g.Wait(); err != nil {
		m.logger.Error(ctx, "failed to query users or user stats", zap.String("reason", err.Error()))
		return nil, err
	}

	items := make([]adminDto.UserListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapUserListItem(row))
	}

	return &adminDto.ListUsersResponse{
		Items:      items,
		Stats:      *stats,
		Pagination: commonDto.NewListMeta(query.Page, query.Limit, total),
	}, nil
}

func (m *adminModule) GetUser(ctx context.Context, userID string) (*adminDto.UserDetailResponse, error) {
	m.logger.Info(ctx, "admin fetching user details", zap.String("target_user_id", userID))

	row, err := m.admin.GetUser(ctx, userID)
	if err != nil {
		m.logger.Warn(ctx, "failed to get user", zap.String("target_user_id", userID), zap.String("reason", err.Error()))
		return nil, mapUserError(err)
	}
	item := mapUserListItem(*row)
	return &adminDto.UserDetailResponse{
		UserListItem: item,
		StatusReason: row.StatusReason,
	}, nil
}

func (m *adminModule) BlockUser(ctx context.Context, adminUserID, targetUserID string, req adminDto.BlockUserRequest, meta authDto.SessionMeta) (*adminDto.BlockUserResponse, error) {
	m.logger.Info(ctx, "admin blocking user", zap.String("admin_id", adminUserID), zap.String("target_user_id", targetUserID), zap.String("reason", req.Reason))

	if err := m.admin.BlockUser(ctx, targetUserID, strings.TrimSpace(req.Reason)); err != nil {
		m.logger.Error(ctx, "failed to block user in persistence", zap.String("target_user_id", targetUserID), zap.String("reason", err.Error()))
		return nil, mapUserError(err)
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &adminUserID,
		TargetUserID: &targetUserID,
		Action:       string(model.BlockAction),
		EntityType:   "user",
		EntityID:     &targetUserID,
		Description:  "user blocked",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	m.logger.Info(ctx, "user blocked successfully", zap.String("target_user_id", targetUserID))

	return &adminDto.BlockUserResponse{UserID: targetUserID, Status: "blocked"}, nil
}

func (m *adminModule) UnblockUser(ctx context.Context, adminUserID, targetUserID string, meta authDto.SessionMeta) (*adminDto.UnblockUserResponse, error) {
	m.logger.Info(ctx, "admin unblocking user", zap.String("admin_id", adminUserID), zap.String("target_user_id", targetUserID))

	if err := m.admin.UnblockUser(ctx, targetUserID); err != nil {
		m.logger.Error(ctx, "failed to unblock user in persistence", zap.String("target_user_id", targetUserID), zap.String("reason", err.Error()))
		return nil, mapUserError(err)
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &adminUserID,
		TargetUserID: &targetUserID,
		Action:       string(model.UnblockAction),
		EntityType:   "user",
		EntityID:     &targetUserID,
		Description:  "user unblocked",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	m.logger.Info(ctx, "user unblocked successfully", zap.String("target_user_id", targetUserID))

	return &adminDto.UnblockUserResponse{UserID: targetUserID, Status: "active"}, nil
}

func (m *adminModule) RemoveUser(ctx context.Context, adminUserID, targetUserID string, req adminDto.RemoveUserRequest, meta authDto.SessionMeta) (*adminDto.RemoveUserResponse, error) {
	m.logger.Info(ctx, "admin removing user", zap.String("admin_id", adminUserID), zap.String("target_user_id", targetUserID), zap.String("reason", req.Reason))

	if err := m.admin.RemoveUser(ctx, targetUserID, strings.TrimSpace(req.Reason)); err != nil {
		m.logger.Error(ctx, "failed to remove user in persistence", zap.String("target_user_id", targetUserID), zap.String("reason", err.Error()))
		return nil, mapUserError(err)
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID:  &adminUserID,
		TargetUserID: &targetUserID,
		Action:       string(model.DeleteAction),
		EntityType:   "user",
		EntityID:     &targetUserID,
		Description:  "user removed",
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})

	m.logger.Info(ctx, "user removed successfully", zap.String("target_user_id", targetUserID))

	return &adminDto.RemoveUserResponse{UserID: targetUserID, Status: "removed"}, nil
}

func mapUserListItem(row adminDto.UserRow) adminDto.UserListItem {
	return adminDto.UserListItem{
		ID:             row.ID,
		Name:           row.DisplayName,
		Email:          row.Email,
		Role:           row.Role,
		Status:         normalizeAdminUserStatus(row.Status),
		JoinedAt:       row.CreatedAt,
		LastActiveAt:   row.LastLoginAt,
		CampaignsCount: row.CampaignsCount,
		ContractsCount: row.ContractsCount,
	}
}

func normalizeAdminUserStatus(status string) string {
	switch status {
	case "pending_profile", "pending_approval":
		return "pending"
	default:
		return status
	}
}

func mapUserError(err error) error {
	if errors.Is(err, adminPersist.ErrUserNotFound) {
		return customErr.ErrUserNotFound
	}
	if errors.Is(err, adminPersist.ErrCannotModifyAdminUser) {
		return customErr.ErrCannotModifyAdminUser
	}
	if errors.Is(err, adminPersist.ErrInvalidUserStatusChange) {
		return customErr.ErrInvalidUserStatusChange
	}
	return err
}
