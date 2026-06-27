package business

import (
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	"adinfluence_service_app/platform/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

func (p *business) ResolveCategoryID(ctx context.Context, nameOrSlug string) (*businessDto.CampaignCategoryResponse, error) {
	nameOrSlug = strings.TrimSpace(nameOrSlug)
	if nameOrSlug == "" {
		return nil, businessDto.ErrUnknownCategory
	}

	var cat businessDto.CampaignCategoryResponse
	err := p.db.QueryRowContext(ctx, resolveCategoryByNameOrSlug, nameOrSlug).Scan(&cat.ID, &cat.Name, &cat.Slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, businessDto.ErrUnknownCategory
		}
		return nil, fmt.Errorf("resolve category: %w", err)
	}
	return &cat, nil
}

func (p *business) CreateCampaign(ctx context.Context, input businessDto.SaveCampaignInput) (*businessDto.CampaignRecord, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin create campaign: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, insertCampaign,
		input.ID,
		input.BusinessUserID,
		input.Title,
		nullString(input.Goal),
		input.Description,
		nullUUID(input.CategoryID),
		nullString(input.TargetAudience),
		nullString(input.City),
		nullString(input.Region),
		nullFloat(input.BudgetMinETB),
		nullFloat(input.BudgetMaxETB),
		input.CreatorSizePreference,
		input.NumberOfCreatorsNeeded,
		nullDate(input.StartDate),
		nullDate(input.EndDate),
	); err != nil {
		return nil, fmt.Errorf("insert campaign: %w", err)
	}

	if err := replaceCampaignPlatforms(ctx, tx, input.ID, input.Platforms); err != nil {
		return nil, err
	}
	if err := replaceCampaignDeliverables(ctx, tx, input.ID, input.Deliverables); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit create campaign: %w", err)
	}

	return p.GetCampaignByID(ctx, input.BusinessUserID, input.ID)
}

func (p *business) ListCampaigns(ctx context.Context, userID string, query businessDto.ListCampaignsQuery) ([]businessDto.CampaignListItem, int, error) {
	query = query.Normalize()
	offset := (query.Page - 1) * query.Limit

	var total int
	if err := p.db.QueryRowContext(ctx, countBusinessCampaigns, userID, query.Status).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count campaigns: %w", err)
	}

	rows, err := p.db.QueryContext(ctx, listBusinessCampaigns, userID, query.Status, query.Limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list campaigns: %w", err)
	}
	defer rows.Close()

	items := make([]businessDto.CampaignListItem, 0)
	for rows.Next() {
		var item businessDto.CampaignListItem
		var budgetMin, budgetMax sql.NullFloat64
		var startDate, endDate sql.NullTime
		var primaryPlatform sql.NullString

		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.BusinessName,
			&budgetMin,
			&budgetMax,
			&item.Status,
			&startDate,
			&endDate,
			&item.CreatorsNeeded,
			&item.CreatedAt,
			&item.InvitesSent,
			&item.InvitesAccepted,
			&primaryPlatform,
		); err != nil {
			return nil, 0, fmt.Errorf("scan campaign list item: %w", err)
		}

		var minPtr, maxPtr *float64
		if budgetMin.Valid {
			v := budgetMin.Float64
			minPtr = &v
		}
		if budgetMax.Valid {
			v := budgetMax.Float64
			maxPtr = &v
		}
		item.BudgetSummary = formatBudgetSummary(minPtr, maxPtr)
		if startDate.Valid {
			t := startDate.Time
			item.StartDate = &t
		}
		if endDate.Valid {
			t := endDate.Time
			item.EndDate = &t
		}
		if primaryPlatform.Valid {
			item.PrimaryPlatform = primaryPlatform.String
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (p *business) GetCampaignByID(ctx context.Context, userID, campaignID string) (*businessDto.CampaignRecord, error) {
	rec, err := scanCampaignRow(p.db.QueryRowContext(ctx, getCampaignByID, campaignID, userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, businessDto.ErrCampaignNotFound
		}
		return nil, fmt.Errorf("get campaign: %w", err)
	}

	platforms, err := p.getCampaignPlatforms(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	deliverables, err := p.getCampaignDeliverables(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	rec.Platforms = platforms
	rec.Deliverables = deliverables
	return rec, nil
}

func (p *business) UpdateCampaign(ctx context.Context, input businessDto.SaveCampaignInput) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin update campaign: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, updateCampaign,
		input.ID,
		input.BusinessUserID,
		input.Title,
		nullString(input.Goal),
		input.Description,
		nullUUID(input.CategoryID),
		nullString(input.TargetAudience),
		nullString(input.City),
		nullString(input.Region),
		nullFloat(input.BudgetMinETB),
		nullFloat(input.BudgetMaxETB),
		input.CreatorSizePreference,
		input.NumberOfCreatorsNeeded,
		nullDate(input.StartDate),
		nullDate(input.EndDate),
	)
	if err != nil {
		return fmt.Errorf("update campaign: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update campaign rows: %w", err)
	}
	if n == 0 {
		return businessDto.ErrCampaignNotFound
	}

	if err := replaceCampaignPlatforms(ctx, tx, input.ID, input.Platforms); err != nil {
		return err
	}
	if err := replaceCampaignDeliverables(ctx, tx, input.ID, input.Deliverables); err != nil {
		return err
	}

	return tx.Commit()
}

func (p *business) PublishCampaign(ctx context.Context, userID, campaignID string) error {
	res, err := p.db.ExecContext(ctx, publishCampaign, campaignID, userID)
	if err != nil {
		return fmt.Errorf("publish campaign: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("publish campaign rows: %w", err)
	}
	if n == 0 {
		return businessDto.ErrInvalidCampaignStatus
	}
	return nil
}

func (p *business) CancelCampaign(ctx context.Context, userID, campaignID string) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin cancel campaign: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, cancelCampaign, campaignID, userID)
	if err != nil {
		return fmt.Errorf("cancel campaign: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("cancel campaign rows: %w", err)
	}
	if n == 0 {
		return businessDto.ErrInvalidCampaignStatus
	}

	if _, err := tx.ExecContext(ctx, cancelPendingCampaignInvites, campaignID); err != nil {
		return fmt.Errorf("cancel pending invites: %w", err)
	}

	return tx.Commit()
}

func (p *business) getCampaignPlatforms(ctx context.Context, campaignID string) ([]string, error) {
	rows, err := p.db.QueryContext(ctx, getCampaignPlatforms, campaignID)
	if err != nil {
		return nil, fmt.Errorf("get campaign platforms: %w", err)
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var platform string
		if err := rows.Scan(&platform); err != nil {
			return nil, fmt.Errorf("scan campaign platform: %w", err)
		}
		out = append(out, platform)
	}
	return out, rows.Err()
}

func (p *business) getCampaignDeliverables(ctx context.Context, campaignID string) ([]businessDto.CampaignDeliverableResponse, error) {
	rows, err := p.db.QueryContext(ctx, getCampaignDeliverables, campaignID)
	if err != nil {
		return nil, fmt.Errorf("get campaign deliverables: %w", err)
	}
	defer rows.Close()

	var out []businessDto.CampaignDeliverableResponse
	for rows.Next() {
		var item businessDto.CampaignDeliverableResponse
		var description sql.NullString
		if err := rows.Scan(&item.ID, &item.Type, &item.Quantity, &description); err != nil {
			return nil, fmt.Errorf("scan campaign deliverable: %w", err)
		}
		if description.Valid {
			item.Description = &description.String
		}
		out = append(out, item)
	}
	if out == nil {
		out = []businessDto.CampaignDeliverableResponse{}
	}
	return out, rows.Err()
}

func replaceCampaignPlatforms(ctx context.Context, tx *sql.Tx, campaignID string, platforms []string) error {
	if _, err := tx.ExecContext(ctx, deleteCampaignPlatforms, campaignID); err != nil {
		return fmt.Errorf("clear campaign platforms: %w", err)
	}
	for _, platform := range platforms {
		if _, err := tx.ExecContext(ctx, insertCampaignPlatform, utils.GenerateUUID(), campaignID, platform); err != nil {
			return fmt.Errorf("insert campaign platform: %w", err)
		}
	}
	return nil
}

func replaceCampaignDeliverables(ctx context.Context, tx *sql.Tx, campaignID string, deliverables []businessDto.CampaignDeliverableInput) error {
	if _, err := tx.ExecContext(ctx, deleteCampaignDeliverables, campaignID); err != nil {
		return fmt.Errorf("clear campaign deliverables: %w", err)
	}
	for _, d := range deliverables {
		var desc *string
		if strings.TrimSpace(d.Description) != "" {
			v := strings.TrimSpace(d.Description)
			desc = &v
		}
		if _, err := tx.ExecContext(ctx, insertCampaignDeliverable,
			utils.GenerateUUID(),
			campaignID,
			strings.ToLower(strings.TrimSpace(d.Type)),
			d.Quantity,
			nullString(desc),
		); err != nil {
			return fmt.Errorf("insert campaign deliverable: %w", err)
		}
	}
	return nil
}

func scanCampaignRow(row *sql.Row) (*businessDto.CampaignRecord, error) {
	var rec businessDto.CampaignRecord
	var goal, targetAudience, city, region sql.NullString
	var categoryID, categoryName, categorySlug sql.NullString
	var budgetMin, budgetMax sql.NullFloat64
	var startDate, endDate sql.NullTime

	err := row.Scan(
		&rec.ID,
		&rec.BusinessUserID,
		&rec.Title,
		&goal,
		&rec.Description,
		&categoryID,
		&categoryName,
		&categorySlug,
		&targetAudience,
		&city,
		&region,
		&budgetMin,
		&budgetMax,
		&rec.CreatorSizePreference,
		&rec.NumberOfCreatorsNeeded,
		&startDate,
		&endDate,
		&rec.Status,
		&rec.CreatedAt,
		&rec.UpdatedAt,
		&rec.BusinessName,
		&rec.InvitesSent,
		&rec.InvitesAccepted,
	)
	if err != nil {
		return nil, err
	}

	if goal.Valid {
		rec.Goal = &goal.String
	}
	if targetAudience.Valid {
		rec.TargetAudience = &targetAudience.String
	}
	if city.Valid {
		rec.City = &city.String
	}
	if region.Valid {
		rec.Region = &region.String
	}
	if budgetMin.Valid {
		v := budgetMin.Float64
		rec.BudgetMinETB = &v
	}
	if budgetMax.Valid {
		v := budgetMax.Float64
		rec.BudgetMaxETB = &v
	}
	if startDate.Valid {
		t := startDate.Time
		rec.StartDate = &t
	}
	if endDate.Valid {
		t := endDate.Time
		rec.EndDate = &t
	}
	if categoryID.Valid {
		rec.Category = &businessDto.CampaignCategoryResponse{
			ID:   categoryID.String,
			Name: categoryName.String,
			Slug: categorySlug.String,
		}
	}
	rec.Platforms = []string{}
	rec.Deliverables = []businessDto.CampaignDeliverableResponse{}

	return &rec, nil
}

func nullDate(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func formatBudgetSummary(min, max *float64) string {
	if min == nil && max == nil {
		return ""
	}
	if min != nil && max != nil {
		return fmt.Sprintf("ETB %.0f – %.0f", *min, *max)
	}
	if min != nil {
		return fmt.Sprintf("ETB %.0f+", *min)
	}
	return fmt.Sprintf("ETB up to %.0f", *max)
}
