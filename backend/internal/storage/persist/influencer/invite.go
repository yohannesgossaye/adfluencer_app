package influencer

import (
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (p *influencer) ListInvites(ctx context.Context, influencerUserID string, query influencerDto.ListInvitesQuery) ([]influencerDto.InviteListItem, error) {
	query = query.Normalize()
	rows, err := p.db.QueryContext(ctx, listInfluencerInvites, influencerUserID, query.Status)
	if err != nil {
		return nil, fmt.Errorf("list influencer invites: %w", err)
	}
	defer rows.Close()

	items := make([]influencerDto.InviteListItem, 0)
	for rows.Next() {
		item, campaignID, err := scanInviteListRow(rows)
		if err != nil {
			return nil, err
		}
		platforms, err := p.getCampaignPlatforms(ctx, campaignID)
		if err != nil {
			return nil, err
		}
		item.Campaign.Platforms = platforms
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate influencer invites: %w", err)
	}
	return items, nil
}

func (p *influencer) GetInviteByID(ctx context.Context, influencerUserID, inviteID string) (*influencerDto.InviteRecord, error) {
	var rec influencerDto.InviteRecord
	var message sql.NullString
	var respondedAt sql.NullTime
	var budgetMin, budgetMax sql.NullFloat64
	var startDate, endDate sql.NullTime

	err := p.db.QueryRowContext(ctx, getInfluencerInviteByID, influencerUserID, inviteID).Scan(
		&rec.ID,
		&rec.Status,
		&message,
		&rec.CreatedAt,
		&respondedAt,
		&rec.BusinessUserID,
		&rec.Campaign.ID,
		&rec.Campaign.Title,
		&rec.Campaign.Status,
		&budgetMin,
		&budgetMax,
		&startDate,
		&endDate,
		&rec.Campaign.BusinessName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, influencerDto.ErrInviteNotFound
		}
		return nil, fmt.Errorf("get influencer invite: %w", err)
	}

	if message.Valid {
		rec.Message = &message.String
	}
	if respondedAt.Valid {
		t := respondedAt.Time
		rec.RespondedAt = &t
	}
	if budgetMin.Valid {
		v := budgetMin.Float64
		rec.Campaign.BudgetMinETB = &v
	}
	if budgetMax.Valid {
		v := budgetMax.Float64
		rec.Campaign.BudgetMaxETB = &v
	}
	if startDate.Valid {
		t := startDate.Time
		rec.Campaign.StartDate = &t
	}
	if endDate.Valid {
		t := endDate.Time
		rec.Campaign.EndDate = &t
	}

	rec.CampaignID = rec.Campaign.ID
	rec.InfluencerUserID = influencerUserID

	platforms, err := p.getCampaignPlatforms(ctx, rec.Campaign.ID)
	if err != nil {
		return nil, err
	}
	rec.Campaign.Platforms = platforms

	return &rec, nil
}

func (p *influencer) AcceptInvite(ctx context.Context, influencerUserID, inviteID string) (*influencerDto.RespondInviteResponse, error) {
	return p.respondInvite(ctx, influencerUserID, inviteID, acceptInfluencerInvite, "accepted")
}

func (p *influencer) RejectInvite(ctx context.Context, influencerUserID, inviteID string) (*influencerDto.RespondInviteResponse, error) {
	return p.respondInvite(ctx, influencerUserID, inviteID, rejectInfluencerInvite, "rejected")
}

func (p *influencer) respondInvite(ctx context.Context, influencerUserID, inviteID, query, status string) (*influencerDto.RespondInviteResponse, error) {
	result, err := p.db.ExecContext(ctx, query, inviteID, influencerUserID)
	if err != nil {
		return nil, fmt.Errorf("respond invite: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("respond invite rows: %w", err)
	}
	if rows == 0 {
		return nil, influencerDto.ErrInvalidInviteStatus
	}

	rec, err := p.GetInviteByID(ctx, influencerUserID, inviteID)
	if err != nil {
		return nil, err
	}

	respondedAt := rec.CreatedAt
	if rec.RespondedAt != nil {
		respondedAt = *rec.RespondedAt
	}

	return &influencerDto.RespondInviteResponse{
		ID:          rec.ID,
		Status:      status,
		RespondedAt: respondedAt,
	}, nil
}

func scanInviteListRow(rows *sql.Rows) (influencerDto.InviteListItem, string, error) {
	var item influencerDto.InviteListItem
	var message sql.NullString
	var respondedAt sql.NullTime
	var budgetMin, budgetMax sql.NullFloat64
	var startDate, endDate sql.NullTime
	var campaignID string

	if err := rows.Scan(
		&item.ID,
		&item.Status,
		&message,
		&item.CreatedAt,
		&respondedAt,
		&campaignID,
		&item.Campaign.Title,
		&item.Campaign.Status,
		&budgetMin,
		&budgetMax,
		&startDate,
		&endDate,
		&item.Campaign.BusinessName,
	); err != nil {
		return influencerDto.InviteListItem{}, "", fmt.Errorf("scan invite list item: %w", err)
	}

	item.Campaign.ID = campaignID
	if message.Valid {
		item.Message = &message.String
	}
	if respondedAt.Valid {
		t := respondedAt.Time
		item.RespondedAt = &t
	}
	if budgetMin.Valid {
		v := budgetMin.Float64
		item.Campaign.BudgetMinETB = &v
	}
	if budgetMax.Valid {
		v := budgetMax.Float64
		item.Campaign.BudgetMaxETB = &v
	}
	if startDate.Valid {
		t := startDate.Time
		item.Campaign.StartDate = &t
	}
	if endDate.Valid {
		t := endDate.Time
		item.Campaign.EndDate = &t
	}

	return item, campaignID, nil
}

func (p *influencer) getCampaignPlatforms(ctx context.Context, campaignID string) ([]string, error) {
	rows, err := p.db.QueryContext(ctx, getCampaignPlatformsForInvite, campaignID)
	if err != nil {
		return nil, fmt.Errorf("get campaign platforms: %w", err)
	}
	defer rows.Close()

	platforms := make([]string, 0)
	for rows.Next() {
		var platform string
		if err := rows.Scan(&platform); err != nil {
			return nil, fmt.Errorf("scan campaign platform: %w", err)
		}
		platforms = append(platforms, platform)
	}
	if platforms == nil {
		platforms = []string{}
	}
	return platforms, rows.Err()
}
