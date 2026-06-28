package admin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var ErrCategoryNotFound = errors.New("category not found")
var ErrCategoryAlreadyExists = errors.New("category already exists")

func (a *admin) ListCategories(ctx context.Context) ([]adminDto.CategoryRow, error) {
	rows, err := a.db.QueryContext(ctx, listCategories)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()

	items := make([]adminDto.CategoryRow, 0)
	for rows.Next() {
		row, err := scanCategoryRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *row)
	}
	if items == nil {
		items = []adminDto.CategoryRow{}
	}
	return items, rows.Err()
}

func (a *admin) GetCategoryByID(ctx context.Context, categoryID string) (*adminDto.CategoryRow, error) {
	row, err := scanCategoryRow(a.db.QueryRowContext(ctx, getCategoryByID, categoryID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, fmt.Errorf("get category: %w", err)
	}
	return row, nil
}

func (a *admin) CreateCategory(ctx context.Context, id, name, slug string, description *string) error {
	_, err := a.db.ExecContext(ctx, insertCategory, id, name, slug, nullStringPtr(description))
	if err != nil {
		return fmt.Errorf("create category: %w", err)
	}
	return nil
}

func (a *admin) UpdateCategory(ctx context.Context, categoryID string, name, slug *string, description *string) (*adminDto.CategoryRow, error) {
	row, err := scanCategoryRow(a.db.QueryRowContext(ctx, updateCategory, categoryID, nullStringPtr(name), nullStringPtr(slug), nullStringPtr(description)))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, fmt.Errorf("update category: %w", err)
	}
	return row, nil
}

func (a *admin) SetCategoryActive(ctx context.Context, categoryID string, active bool) error {
	result, err := a.db.ExecContext(ctx, toggleCategoryActive, categoryID, active)
	if err != nil {
		return fmt.Errorf("toggle category: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("toggle category rows: %w", err)
	}
	if n == 0 {
		return ErrCategoryNotFound
	}
	return nil
}

func (a *admin) CategorySlugExists(ctx context.Context, slug, excludeID string) (bool, error) {
	var exists bool
	if err := a.db.QueryRowContext(ctx, categorySlugExists, slug, excludeID).Scan(&exists); err != nil {
		return false, fmt.Errorf("check category slug: %w", err)
	}
	return exists, nil
}

func scanCategoryRow(s scanner) (*adminDto.CategoryRow, error) {
	var row adminDto.CategoryRow
	var desc sql.NullString
	if err := s.Scan(
		&row.ID,
		&row.Name,
		&row.Slug,
		&desc,
		&row.IsActive,
		&row.CreatedAt,
		&row.UpdatedAt,
		&row.InfluencerCount,
		&row.CampaignCount,
	); err != nil {
		return nil, err
	}
	if desc.Valid {
		row.Description = &desc.String
	}
	return &row, nil
}

func nullStringPtr(v *string) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}
	trimmed := strings.TrimSpace(*v)
	if trimmed == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: trimmed, Valid: true}
}
