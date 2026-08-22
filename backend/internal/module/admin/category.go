package admin

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	"adinfluence_service_app/internal/module/admin/helper"
	adminPersist "adinfluence_service_app/internal/storage/persist/admin"
	"adinfluence_service_app/platform/utils"
	"context"
	"errors"
	"strings"

	"go.uber.org/zap"
)

func (m *adminModule) ListCategories(ctx context.Context) (*adminDto.ListCategoriesResponse, error) {
	m.logger.Info(ctx, "admin listing categories")

	rows, err := m.admin.ListCategories(ctx)
	if err != nil {
		m.logger.Error(ctx, "failed to list categories from database", zap.String("reason", err.Error()))
		return nil, err
	}

	items := make([]adminDto.CategoryListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapCategoryItem(row))
	}

	return &adminDto.ListCategoriesResponse{Items: items}, nil
}

func (m *adminModule) CreateCategory(ctx context.Context, adminUserID string, req adminDto.CreateCategoryRequest, meta authDto.SessionMeta) (*commonDto.CreatedIDResponse, error) {
	m.logger.Info(ctx, "admin creating category", zap.String("admin_id", adminUserID), zap.String("name", req.Name))

	if err := req.Validate(); err != nil {
		m.logger.Warn(ctx, "category validation failed", zap.String("reason", err.Error()))
		return nil, err
	}

	name := strings.TrimSpace(req.Name)
	slug := helper.Slugify(name)
	exists, err := m.admin.CategorySlugExists(ctx, slug, "")
	if err != nil {
		return nil, err
	}
	if exists {
		m.logger.Warn(ctx, "category slug already exists", zap.String("slug", slug))
		return nil, customErr.ErrCategoryAlreadyExists
	}

	id := utils.GenerateUUID()
	if err := m.admin.CreateCategory(ctx, id, name, slug, optionalTrimmed(req.Description)); err != nil {
		m.logger.Error(ctx, "failed to persist category", zap.String("id", id), zap.String("reason", err.Error()))
		return nil, mapCategoryError(err)
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID: &adminUserID,
		Action:      string(model.AdminAction),
		EntityType:  "category",
		EntityID:    &id,
		Description: "category created: " + name,
		IPAddress:   meta.IPAddress,
		UserAgent:   meta.UserAgent,
	})

	m.logger.Info(ctx, "category created successfully", zap.String("id", id), zap.String("name", name), zap.String("slug", slug))

	return &commonDto.CreatedIDResponse{ID: id}, nil
}

func (m *adminModule) UpdateCategory(ctx context.Context, adminUserID, categoryID string, req adminDto.UpdateCategoryRequest, meta authDto.SessionMeta) (*adminDto.UpdateCategoryResponse, error) {
	m.logger.Info(ctx, "admin updating category", zap.String("admin_id", adminUserID), zap.String("category_id", categoryID))

	if err := req.Validate(); err != nil {
		return nil, customErr.ErrBadRequest
	}

	var namePtr, slugPtr *string
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		namePtr = &name
		slug := helper.Slugify(name)
		exists, err := m.admin.CategorySlugExists(ctx, slug, categoryID)
		if err != nil {
			return nil, err
		}
		if exists {
			m.logger.Warn(ctx, "category slug collision on update", zap.String("slug", slug))
			return nil, customErr.ErrCategoryAlreadyExists
		}
		slugPtr = &slug
	}

	row, err := m.admin.UpdateCategory(ctx, categoryID, namePtr, slugPtr, optionalTrimmedPtr(req.Description))
	if err != nil {
		m.logger.Error(ctx, "failed to update category in persistence", zap.String("category_id", categoryID), zap.String("reason", err.Error()))
		return nil, mapCategoryError(err)
	}

	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID: &adminUserID,
		Action:      string(model.AdminAction),
		EntityType:  "category",
		EntityID:    &categoryID,
		Description: "category updated",
		IPAddress:   meta.IPAddress,
		UserAgent:   meta.UserAgent,
	})

	m.logger.Info(ctx, "category updated successfully", zap.String("category_id", categoryID))

	item := mapCategoryItem(*row)
	return &adminDto.UpdateCategoryResponse{CategoryListItem: item}, nil
}

func (m *adminModule) DisableCategory(ctx context.Context, adminUserID, categoryID string, meta authDto.SessionMeta) (*adminDto.ToggleCategoryActiveResponse, error) {
	m.logger.Info(ctx, "admin disabling category", zap.String("category_id", categoryID))
	return m.setCategoryActive(ctx, adminUserID, categoryID, false, meta)
}

func (m *adminModule) EnableCategory(ctx context.Context, adminUserID, categoryID string, meta authDto.SessionMeta) (*adminDto.ToggleCategoryActiveResponse, error) {
	m.logger.Info(ctx, "admin enabling category", zap.String("category_id", categoryID))
	return m.setCategoryActive(ctx, adminUserID, categoryID, true, meta)
}

func (m *adminModule) setCategoryActive(ctx context.Context, adminUserID, categoryID string, active bool, meta authDto.SessionMeta) (*adminDto.ToggleCategoryActiveResponse, error) {
	if err := m.admin.SetCategoryActive(ctx, categoryID, active); err != nil {
		m.logger.Error(ctx, "failed to toggle category active status", zap.String("category_id", categoryID), zap.String("reason", err.Error()))
		return nil, mapCategoryError(err)
	}

	action := "category disabled"
	if active {
		action = "category enabled"
	}
	m.writeAuditAsync(authDto.CreateAuditLogInput{
		ActorUserID: &adminUserID,
		Action:      string(model.AdminAction),
		EntityType:  "category",
		EntityID:    &categoryID,
		Description: action,
		IPAddress:   meta.IPAddress,
		UserAgent:   meta.UserAgent,
	})

	m.logger.Info(ctx, "category active state toggled successfully", zap.String("category_id", categoryID), zap.Bool("active", active))

	return &adminDto.ToggleCategoryActiveResponse{
		ID:       categoryID,
		IsActive: active,
	}, nil
}

func mapCategoryItem(row adminDto.CategoryRow) adminDto.CategoryListItem {
	return adminDto.CategoryListItem{
		ID:              row.ID,
		Name:            row.Name,
		Slug:            row.Slug,
		Description:     row.Description,
		IsActive:        row.IsActive,
		InfluencerCount: row.InfluencerCount,
		CampaignCount:   row.CampaignCount,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}
}

func mapCategoryError(err error) error {
	if errors.Is(err, adminPersist.ErrCategoryNotFound) {
		return customErr.ErrCategoryNotFound
	}
	if errors.Is(err, adminPersist.ErrCategoryAlreadyExists) {
		return customErr.ErrCategoryAlreadyExists
	}
	return err
}

func optionalTrimmed(v *string) *string {
	if v == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*v)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func optionalTrimmedPtr(v *string) *string {
	if v == nil {
		return nil
	}
	return optionalTrimmed(v)
}
