package admin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var ErrUserNotFound = errors.New("user not found")
var ErrCannotModifyAdminUser = errors.New("cannot modify admin user")
var ErrInvalidUserStatusChange = errors.New("invalid user status change")

func (a *admin) GetUserStats(ctx context.Context) (*adminDto.UserStats, error) {
	var stats adminDto.UserStats
	if err := a.db.QueryRowContext(ctx, getAdminUserStats).Scan(
		&stats.Total,
		&stats.Active,
		&stats.Blocked,
		&stats.Pending,
	); err != nil {
		return nil, fmt.Errorf("get user stats: %w", err)
	}
	return &stats, nil
}

func (a *admin) ListUsers(ctx context.Context, query adminDto.ListUsersQuery) ([]adminDto.UserRow, int, error) {
	query = query.Normalize()
	offset := (query.Page - 1) * query.Limit

	rows, err := a.db.QueryContext(ctx, listAdminUsers, query.Role, query.Status, query.Search, query.Limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var total int
	items := make([]adminDto.UserRow, 0)
	for rows.Next() {
		row, err := scanUserListRow(rows, &total)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *row)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (a *admin) GetUser(ctx context.Context, userID string) (*adminDto.UserRow, error) {
	row, err := scanUserRow(a.db.QueryRowContext(ctx, getAdminUser, userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user: %w", err)
	}
	return row, nil
}

func (a *admin) BlockUser(ctx context.Context, userID, reason string) error {
	role, status, err := a.getUserRoleStatus(ctx, userID)
	if err != nil {
		return err
	}
	if role == "admin" {
		return ErrCannotModifyAdminUser
	}
	if status == "removed" {
		return ErrInvalidUserStatusChange
	}

	return a.setStatusAndRevokeSessions(ctx, blockUser, "block user", userID, reason)
}

// UnblockUser restores the status the user had before being blocked and returns it.
func (a *admin) UnblockUser(ctx context.Context, userID string) (string, error) {
	role, _, err := a.getUserRoleStatus(ctx, userID)
	if err != nil {
		return "", err
	}
	if role == "admin" {
		return "", ErrCannotModifyAdminUser
	}

	var restored string
	err = a.db.QueryRowContext(ctx, unblockUser, userID).Scan(&restored)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrInvalidUserStatusChange
	}
	if err != nil {
		return "", fmt.Errorf("unblock user: %w", err)
	}
	return restored, nil
}

func (a *admin) RemoveUser(ctx context.Context, userID, reason string) error {
	role, status, err := a.getUserRoleStatus(ctx, userID)
	if err != nil {
		return err
	}
	if role == "admin" {
		return ErrCannotModifyAdminUser
	}
	if status == "removed" {
		return ErrInvalidUserStatusChange
	}

	return a.setStatusAndRevokeSessions(ctx, removeUser, "remove user", userID, reason)
}

// setStatusAndRevokeSessions applies a blocking status change and signs the
// user out everywhere in the same transaction, so refresh tokens issued
// before the block can no longer be used.
func (a *admin) setStatusAndRevokeSessions(ctx context.Context, statusQuery, op, userID, reason string) error {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: begin tx: %w", op, err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, statusQuery, userID, nullReason(reason))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return ErrUserNotFound
	}
	if _, err := tx.ExecContext(ctx, revokeUserSessions, userID); err != nil {
		return fmt.Errorf("%s: revoke sessions: %w", op, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: commit: %w", op, err)
	}
	return nil
}

func (a *admin) getUserRoleStatus(ctx context.Context, userID string) (role, status string, err error) {
	err = a.db.QueryRowContext(ctx, getUserRole, userID).Scan(&role, &status)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", ErrUserNotFound
	}
	if err != nil {
		return "", "", fmt.Errorf("get user role: %w", err)
	}
	return role, status, nil
}

func (a *admin) userExists(ctx context.Context, userID string) (bool, error) {
	var exists bool
	err := a.db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check user exists: %w", err)
	}
	return exists, nil
}

func scanUserListRow(s scanner, total *int) (*adminDto.UserRow, error) {
	var row adminDto.UserRow
	var reason sql.NullString
	var lastLogin sql.NullTime
	if err := s.Scan(
		&row.ID,
		&row.Email,
		&row.Role,
		&row.Status,
		&reason,
		&row.DisplayName,
		&row.CreatedAt,
		&lastLogin,
		&row.CampaignsCount,
		&row.ContractsCount,
		total,
	); err != nil {
		return nil, err
	}
	if reason.Valid {
		row.StatusReason = &reason.String
	}
	if lastLogin.Valid {
		t := lastLogin.Time
		row.LastLoginAt = &t
	}
	return &row, nil
}

func scanUserRow(s scanner) (*adminDto.UserRow, error) {
	var row adminDto.UserRow
	var reason sql.NullString
	var lastLogin sql.NullTime
	if err := s.Scan(
		&row.ID,
		&row.Email,
		&row.Role,
		&row.Status,
		&reason,
		&row.DisplayName,
		&row.CreatedAt,
		&lastLogin,
		&row.CampaignsCount,
		&row.ContractsCount,
	); err != nil {
		return nil, err
	}
	if reason.Valid {
		row.StatusReason = &reason.String
	}
	if lastLogin.Valid {
		t := lastLogin.Time
		row.LastLoginAt = &t
	}
	return &row, nil
}

func nullReason(reason string) sql.NullString {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: reason, Valid: true}
}
