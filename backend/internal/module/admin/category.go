package admin

import (
	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	"adinfluence_service_app/internal/module/admin/helper"
	adminPersist "adinfluence_service_app/internal/storage/persist/admin"
	"adinfluence_service_app/platform/utils"
	"context"
	"errors"
	"strings"
)

func (m *adminModule) ListCategories(ctx context.Context) (*adminDto.ListCategoriesResponse, error) {
	rows, err := m.admin.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]adminDto.CategoryListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapCategoryItem(row))
	}

	return &adminDto.ListCategoriesResponse{Items: items}, nil
}

func (m *adminModule) CreateCategory(ctx context.Context, adminUserID string, req adminDto.CreateCategoryRequest, meta authDto.SessionMeta) (*adminDto.CreateCategoryResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, customErr.ErrBadRequest
	}

	name := strings.TrimSpace(req.Name)
	slug := helper.Slugify(name)
	exists, err := m.admin.CategorySlugExists(ctx, slug, "")
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, customErr.ErrCategoryAlreadyExists
	}

	id := utils.GenerateUUID()
	row, err := m.admin.CreateCategory(ctx, id, name, slug, optionalTrimmed(req.Description))
	if err != nil {
		return nil, err
	}

	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID: &adminUserID,
		Action:      string(model.AdminAction),
		EntityType:  "category",
		EntityID:    &id,
		Description: "category created: " + name,
		IPAddress:   meta.IPAddress,
		UserAgent:   meta.UserAgent,
	})

	item := mapCategoryItem(*row)
	return &adminDto.CreateCategoryResponse{CategoryListItem: item}, nil
}

func (m *adminModule) UpdateCategory(ctx context.Context, adminUserID, categoryID string, req adminDto.UpdateCategoryRequest, meta authDto.SessionMeta) (*adminDto.UpdateCategoryResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, customErr.ErrBadRequest
	}

	if _, err := m.admin.GetCategoryByID(ctx, categoryID); err != nil {
		return nil, mapCategoryError(err)
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
			return nil, customErr.ErrCategoryAlreadyExists
		}
		slugPtr = &slug
	}

	if err := m.admin.UpdateCategory(ctx, categoryID, namePtr, slugPtr, optionalTrimmedPtr(req.Description)); err != nil {
		return nil, mapCategoryError(err)
	}

	row, err := m.admin.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return nil, mapCategoryError(err)
	}

	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID: &adminUserID,
		Action:      string(model.AdminAction),
		EntityType:  "category",
		EntityID:    &categoryID,
		Description: "category updated",
		IPAddress:   meta.IPAddress,
		UserAgent:   meta.UserAgent,
	})

	item := mapCategoryItem(*row)
	return &adminDto.UpdateCategoryResponse{CategoryListItem: item}, nil
}

func (m *adminModule) DisableCategory(ctx context.Context, adminUserID, categoryID string, meta authDto.SessionMeta) (*adminDto.ToggleCategoryActiveResponse, error) {
	return m.setCategoryActive(ctx, adminUserID, categoryID, false, meta)
}

func (m *adminModule) EnableCategory(ctx context.Context, adminUserID, categoryID string, meta authDto.SessionMeta) (*adminDto.ToggleCategoryActiveResponse, error) {
	return m.setCategoryActive(ctx, adminUserID, categoryID, true, meta)
}

func (m *adminModule) setCategoryActive(ctx context.Context, adminUserID, categoryID string, active bool, meta authDto.SessionMeta) (*adminDto.ToggleCategoryActiveResponse, error) {
	if err := m.admin.SetCategoryActive(ctx, categoryID, active); err != nil {
		return nil, mapCategoryError(err)
	}

	action := "category disabled"
	if active {
		action = "category enabled"
	}
	m.writeAudit(ctx, authDto.CreateAuditLogInput{
		ActorUserID: &adminUserID,
		Action:      string(model.AdminAction),
		EntityType:  "category",
		EntityID:    &categoryID,
		Description: action,
		IPAddress:   meta.IPAddress,
		UserAgent:   meta.UserAgent,
	})

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
