package admin

import (
	"errors"
	"strings"
	"time"
)

var ErrCategoryNameRequired = errors.New("category name is required")

type CategoryListItem struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Slug             string    `json:"slug"`
	Description      *string   `json:"description,omitempty"`
	IsActive         bool      `json:"is_active"`
	InfluencerCount  int       `json:"influencer_count"`
	CampaignCount    int       `json:"campaign_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ListCategoriesResponse struct {
	Items []CategoryListItem `json:"items"`
}

type CreateCategoryRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

func (r CreateCategoryRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return ErrCategoryNameRequired
	}
	return nil
}

type CreateCategoryResponse struct {
	CategoryListItem
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (r UpdateCategoryRequest) Validate() error {
	if r.Name != nil && strings.TrimSpace(*r.Name) == "" {
		return ErrCategoryNameRequired
	}
	if r.Name == nil && r.Description == nil {
		return errors.New("no fields to update")
	}
	return nil
}

type UpdateCategoryResponse struct {
	CategoryListItem
}

type ToggleCategoryActiveResponse struct {
	ID       string `json:"id"`
	IsActive bool   `json:"is_active"`
}

type CategoryRow struct {
	ID              string
	Name            string
	Slug            string
	Description     *string
	IsActive        bool
	InfluencerCount int
	CampaignCount   int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
