package admin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrAdminDisputeNotFound = errors.New("dispute not found")

func (a *admin) ListAdminDisputes(ctx context.Context, query adminDto.ListAdminDisputesQuery) ([]adminDto.DisputeRow, int, error) {
	query = query.Normalize()
	offset := (query.Page - 1) * query.Limit

	var total int
	if err := a.db.QueryRowContext(ctx, countAdminDisputes, query.Tab, query.Search).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count admin disputes: %w", err)
	}

	rows, err := a.db.QueryContext(ctx, listAdminDisputes, query.Tab, query.Search, query.Limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list admin disputes: %w", err)
	}
	defer rows.Close()

	items := make([]adminDto.DisputeRow, 0)
	for rows.Next() {
		var row adminDto.DisputeRow
		if err := rows.Scan(
			&row.ID,
			&row.Description,
			&row.IssueType,
			&row.Severity,
			&row.Status,
			&row.CreatedAt,
			&row.CampaignTitle,
			&row.BusinessName,
			&row.InfluencerName,
		); err != nil {
			return nil, 0, fmt.Errorf("scan dispute list row: %w", err)
		}
		items = append(items, row)
	}
	if items == nil {
		items = []adminDto.DisputeRow{}
	}
	return items, total, rows.Err()
}

func (a *admin) GetAdminDispute(ctx context.Context, disputeID string) (*adminDto.DisputeRow, error) {
	var row adminDto.DisputeRow
	var resolutionNote sql.NullString
	var resolvedAt sql.NullTime
	err := a.db.QueryRowContext(ctx, getAdminDispute, disputeID).Scan(
		&row.ID,
		&row.ContractID,
		&row.CampaignID,
		&row.CampaignTitle,
		&row.BusinessName,
		&row.InfluencerName,
		&row.InfluencerEmail,
		&row.InfluencerUserID,
		&row.AgainstUserID,
		&row.IssueType,
		&row.Severity,
		&row.Status,
		&row.Description,
		&resolutionNote,
		&row.CreatedAt,
		&row.UpdatedAt,
		&resolvedAt,
		&row.AdminNote,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAdminDisputeNotFound
		}
		return nil, fmt.Errorf("get admin dispute: %w", err)
	}
	if resolutionNote.Valid {
		row.ResolutionNote = &resolutionNote.String
	}
	if resolvedAt.Valid {
		t := resolvedAt.Time
		row.ResolvedAt = &t
	}
	return &row, nil
}

func (a *admin) UpdateDisputeStatus(ctx context.Context, disputeID, status, adminUserID string) error {
	result, err := a.db.ExecContext(ctx, updateDisputeStatus, disputeID, status, adminUserID)
	if err != nil {
		return fmt.Errorf("update dispute status: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return ErrAdminDisputeNotFound
	}
	return nil
}

func (a *admin) UpdateDisputeResolutionNote(ctx context.Context, disputeID string, note *string) error {
	result, err := a.db.ExecContext(ctx, updateDisputeResolutionNote, disputeID, nullStringPtr(note))
	if err != nil {
		return fmt.Errorf("update dispute resolution note: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return ErrAdminDisputeNotFound
	}
	return nil
}

func (a *admin) SaveDisputeAdminNote(ctx context.Context, adminUserID, disputeID, note string) error {
	return a.CreateAuditLog(ctx, authDto.CreateAuditLogInput{
		ActorUserID: &adminUserID,
		Action:      "admin_action",
		EntityType:  "dispute_note",
		EntityID:    &disputeID,
		Description: note,
	})
}
