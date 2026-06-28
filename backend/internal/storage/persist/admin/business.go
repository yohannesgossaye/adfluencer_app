package admin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var ErrBusinessVerificationNotFound = errors.New("business verification not found")
var ErrBusinessAlreadyVerified = errors.New("business already verified")
var ErrBusinessNotVerified = errors.New("business not verified")

func (a *admin) ListBusinessVerifications(ctx context.Context, query adminDto.ListBusinessVerificationsQuery) ([]adminDto.BusinessVerificationListItem, int, error) {
	query = query.Normalize()
	offset := (query.Page - 1) * query.Limit
	search := strings.TrimSpace(query.Search)

	rows, err := a.db.QueryContext(ctx, listBusinessVerifications, query.Status, search, query.Limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list business verifications: %w", err)
	}
	defer rows.Close()

	var total int
	items := make([]adminDto.BusinessVerificationListItem, 0)
	for rows.Next() {
		var item adminDto.BusinessVerificationListItem
		var industry, city sql.NullString
		var logo sql.NullString
		if err := rows.Scan(
			&item.UserID,
			&item.BusinessName,
			&item.OwnerFullName,
			&industry,
			&city,
			&logo,
			&item.Email,
			&item.IsVerified,
			&item.RegisteredAt,
			&total,
		); err != nil {
			return nil, 0, fmt.Errorf("scan business verification: %w", err)
		}
		if industry.Valid {
			item.Industry = industry.String
		}
		if city.Valid {
			item.City = city.String
		}
		if logo.Valid {
			item.LogoURL = &logo.String
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (a *admin) GetBusinessVerification(ctx context.Context, userID string) (*adminDto.BusinessVerificationRow, error) {
	row, err := scanBusinessVerificationRow(a.db.QueryRowContext(ctx, getBusinessVerification, userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBusinessVerificationNotFound
		}
		return nil, fmt.Errorf("get business verification: %w", err)
	}
	return row, nil
}

func (a *admin) VerifyBusiness(ctx context.Context, userID string) error {
	result, err := a.db.ExecContext(ctx, verifyBusinessProfile, userID)
	if err != nil {
		return fmt.Errorf("verify business: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("verify business rows: %w", err)
	}
	if n == 0 {
		exists, err := a.businessVerificationExists(ctx, userID)
		if err != nil {
			return err
		}
		if !exists {
			return ErrBusinessVerificationNotFound
		}
		return ErrBusinessAlreadyVerified
	}
	return nil
}

func (a *admin) RevokeBusinessVerification(ctx context.Context, userID string) error {
	result, err := a.db.ExecContext(ctx, revokeBusinessVerification, userID)
	if err != nil {
		return fmt.Errorf("revoke business verification: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("revoke business rows: %w", err)
	}
	if n == 0 {
		exists, err := a.businessVerificationExists(ctx, userID)
		if err != nil {
			return err
		}
		if !exists {
			return ErrBusinessVerificationNotFound
		}
		return ErrBusinessNotVerified
	}
	return nil
}

func (a *admin) businessVerificationExists(ctx context.Context, userID string) (bool, error) {
	var exists bool
	err := a.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM users u
			INNER JOIN business_profiles bp ON bp.user_id = u.id
			WHERE u.id = $1 AND u.role = 'business' AND u.status = 'active'
		)`, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check business verification exists: %w", err)
	}
	return exists, nil
}

func scanBusinessVerificationRow(row *sql.Row) (*adminDto.BusinessVerificationRow, error) {
	var rec adminDto.BusinessVerificationRow
	var industry, city, description, logo sql.NullString
	if err := row.Scan(
		&rec.UserID,
		&rec.Email,
		&rec.PhoneNumber,
		&rec.UserStatus,
		&rec.BusinessName,
		&rec.OwnerFullName,
		&industry,
		&city,
		&description,
		&logo,
		&rec.IsVerified,
		&rec.RegisteredAt,
	); err != nil {
		return nil, err
	}
	if industry.Valid {
		rec.Industry = &industry.String
	}
	if city.Valid {
		rec.City = &city.String
	}
	if description.Valid {
		rec.Description = &description.String
	}
	if logo.Valid {
		rec.LogoURL = &logo.String
	}
	return &rec, nil
}
