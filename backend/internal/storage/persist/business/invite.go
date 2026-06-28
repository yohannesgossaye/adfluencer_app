package business

import (
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (p *business) CreateInvite(ctx context.Context, input businessDto.CreateInviteInput) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin create invite: %w", err)
	}
	defer tx.Rollback()

	var exists bool
	if err := tx.QueryRowContext(ctx, inviteExistsForCampaign, input.CampaignID, input.InfluencerUserID).Scan(&exists); err != nil {
		return fmt.Errorf("check invite exists: %w", err)
	}
	if exists {
		return businessDto.ErrInviteAlreadyExists
	}

	if _, err := tx.ExecContext(ctx, insertCampaignInvite,
		input.ID,
		input.CampaignID,
		input.BusinessUserID,
		input.InfluencerUserID,
		nullString(input.Message),
	); err != nil {
		return fmt.Errorf("insert invite: %w", err)
	}

	if _, err := tx.ExecContext(ctx, updateRecommendationStatusInvited, input.CampaignID, input.InfluencerUserID); err != nil {
		return fmt.Errorf("update recommendation invited: %w", err)
	}

	return tx.Commit()
}

func (p *business) ListInvites(ctx context.Context, businessUserID, campaignID string) ([]businessDto.InviteResponse, error) {
	rows, err := p.db.QueryContext(ctx, listCampaignInvites, campaignID, businessUserID)
	if err != nil {
		return nil, fmt.Errorf("list invites: %w", err)
	}
	defer rows.Close()

	items := make([]businessDto.InviteResponse, 0)
	var rowCount int
	for rows.Next() {
		rowCount++
		item, ok, err := scanInviteListRow(rows)
		if err != nil {
			return nil, err
		}
		if ok {
			items = append(items, item)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate invites: %w", err)
	}
	if rowCount == 0 {
		return nil, businessDto.ErrCampaignNotFound
	}
	return items, nil
}

func (p *business) GetInviteByID(ctx context.Context, businessUserID, campaignID, inviteID string) (*businessDto.InviteResponse, error) {
	row := p.db.QueryRowContext(ctx, getCampaignInviteByID, inviteID, campaignID, businessUserID)
	item, err := scanInviteFromRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, businessDto.ErrInviteNotFound
		}
		return nil, fmt.Errorf("get invite: %w", err)
	}
	return &item, nil
}

func (p *business) CancelInvite(ctx context.Context, businessUserID, campaignID, inviteID string) error {
	result, err := p.db.ExecContext(ctx, cancelPendingInvite, inviteID, campaignID, businessUserID)
	if err != nil {
		return fmt.Errorf("cancel invite: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("cancel invite rows: %w", err)
	}
	if rows == 0 {
		return businessDto.ErrInvalidInviteStatus
	}
	return nil
}

func scanInviteListRow(rows *sql.Rows) (businessDto.InviteResponse, bool, error) {
	var item businessDto.InviteResponse
	var inviteID sql.NullString
	var influencerUserID sql.NullString
	var creatorName sql.NullString
	var status sql.NullString
	var message sql.NullString
	var createdAt sql.NullTime
	if err := rows.Scan(
		&inviteID,
		&item.CampaignID,
		&influencerUserID,
		&creatorName,
		&status,
		&message,
		&createdAt,
	); err != nil {
		return businessDto.InviteResponse{}, false, fmt.Errorf("scan invite: %w", err)
	}
	if !inviteID.Valid {
		return businessDto.InviteResponse{}, false, nil
	}
	item.ID = inviteID.String
	item.InfluencerUserID = influencerUserID.String
	item.CreatorName = creatorName.String
	item.Status = status.String
	if message.Valid {
		item.Message = &message.String
	}
	if createdAt.Valid {
		item.CreatedAt = createdAt.Time
	}
	return item, true, nil
}

func scanInviteFromRow(row *sql.Row) (businessDto.InviteResponse, error) {
	var item businessDto.InviteResponse
	var businessUserID string
	var message sql.NullString
	if err := row.Scan(
		&item.ID,
		&item.CampaignID,
		&businessUserID,
		&item.InfluencerUserID,
		&item.CreatorName,
		&item.Status,
		&message,
		&item.CreatedAt,
	); err != nil {
		return businessDto.InviteResponse{}, err
	}
	if message.Valid {
		item.Message = &message.String
	}
	return item, nil
}
