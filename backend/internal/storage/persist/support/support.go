package support

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	"adinfluence_service_app/internal/storage"
	"adinfluence_service_app/platform/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var ErrTicketNotFound = errors.New("support ticket not found")
var ErrInviteNotFound = errors.New("support invite not found")

type support struct {
	db *sql.DB
}

func InitSupportPersistence(db *sql.DB) *support {
	return &support{db: db}
}

func (s *support) ListBusinessCampaignOptions(ctx context.Context, businessUserID string) ([]businessDto.SupportCampaignOption, error) {
	rows, err := s.db.QueryContext(ctx, listBusinessCampaignOptions, businessUserID)
	if err != nil {
		return nil, fmt.Errorf("list support campaign options: %w", err)
	}
	defer rows.Close()

	items := make([]businessDto.SupportCampaignOption, 0)
	for rows.Next() {
		var item businessDto.SupportCampaignOption
		if err := rows.Scan(&item.ID, &item.Title); err != nil {
			return nil, fmt.Errorf("scan campaign option: %w", err)
		}
		items = append(items, item)
	}
	if items == nil {
		items = []businessDto.SupportCampaignOption{}
	}
	return items, rows.Err()
}

func (s *support) CreateTicket(ctx context.Context, input businessDto.CreateSupportTicketInput) error {
	_, err := s.db.ExecContext(ctx, insertSupportTicket,
		input.ID,
		input.UserID,
		input.Subject,
		input.Message,
		input.IssueType,
		input.Priority,
		nullUUID(input.RelatedCampaignID),
	)
	if err != nil {
		return fmt.Errorf("create support ticket: %w", err)
	}
	return nil
}

func (s *support) ListTicketsByUser(ctx context.Context, userID string) ([]businessDto.SupportTicketRecord, error) {
	rows, err := s.db.QueryContext(ctx, listSupportTicketsByUser, userID)
	if err != nil {
		return nil, fmt.Errorf("list support tickets: %w", err)
	}
	defer rows.Close()
	return scanBusinessTicketRows(rows)
}

func (s *support) GetTicketByUser(ctx context.Context, userID, ticketID string) (*businessDto.SupportTicketRecord, error) {
	row := s.db.QueryRowContext(ctx, getSupportTicketByUser, ticketID, userID)
	rec, err := scanBusinessTicketRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTicketNotFound
		}
		return nil, fmt.Errorf("get support ticket: %w", err)
	}
	return rec, nil
}

func (s *support) BusinessOwnsCampaign(ctx context.Context, businessUserID, campaignID string) (bool, error) {
	var exists bool
	if err := s.db.QueryRowContext(ctx, businessOwnsCampaign, campaignID, businessUserID).Scan(&exists); err != nil {
		return false, fmt.Errorf("check campaign ownership: %w", err)
	}
	return exists, nil
}

func (s *support) ListAdminTickets(ctx context.Context, query adminDto.ListSupportTicketsQuery) ([]adminDto.SupportTicketRow, int, error) {
	query = query.Normalize()
	offset := (query.Page - 1) * query.Limit

	var total int
	if err := s.db.QueryRowContext(ctx, countAdminSupportTickets, query.Role, query.Status, query.Priority, query.Search).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count support tickets: %w", err)
	}

	rows, err := s.db.QueryContext(ctx, listAdminSupportTickets, query.Role, query.Status, query.Priority, query.Search, query.Limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list admin support tickets: %w", err)
	}
	defer rows.Close()

	items := make([]adminDto.SupportTicketRow, 0)
	for rows.Next() {
		row, err := scanAdminListRow(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *row)
	}
	if items == nil {
		items = []adminDto.SupportTicketRow{}
	}
	return items, total, rows.Err()
}

func (s *support) GetAdminTicketStats(ctx context.Context) (*adminDto.SupportTicketStats, error) {
	var stats adminDto.SupportTicketStats
	if err := s.db.QueryRowContext(ctx, getAdminSupportTicketStats).Scan(
		&stats.OpenCount,
		&stats.HighPriority,
		&stats.BusinessCount,
		&stats.InfluencerCount,
		&stats.ResolvedCount,
	); err != nil {
		return nil, fmt.Errorf("get support ticket stats: %w", err)
	}
	return &stats, nil
}

func (s *support) GetAdminTicket(ctx context.Context, ticketID string) (*adminDto.SupportTicketRow, error) {
	row, err := scanAdminTicketRow(s.db.QueryRowContext(ctx, getAdminSupportTicket, ticketID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTicketNotFound
		}
		return nil, fmt.Errorf("get admin support ticket: %w", err)
	}
	return row, nil
}

func (s *support) UpdateTicketStatus(ctx context.Context, ticketID, status, adminUserID string) error {
	result, err := s.db.ExecContext(ctx, updateSupportTicketStatus, ticketID, status, adminUserID)
	if err != nil {
		return fmt.Errorf("update support ticket status: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update support ticket status rows: %w", err)
	}
	if n == 0 {
		return ErrTicketNotFound
	}
	return nil
}

func (s *support) SaveTicketReply(ctx context.Context, ticketID, adminUserID string, adminReply, internalNote *string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin support reply tx: %w", err)
	}
	defer tx.Rollback()

	if adminReply != nil {
		if _, err := tx.ExecContext(ctx, updateSupportTicketReply, ticketID, strings.TrimSpace(*adminReply), adminUserID); err != nil {
			return fmt.Errorf("update support ticket reply: %w", err)
		}
	} else {
		if _, err := tx.ExecContext(ctx, `
			UPDATE support_tickets
			SET assigned_admin_user_id = $2, updated_at = NOW()
			WHERE id = $1`, ticketID, adminUserID); err != nil {
			return fmt.Errorf("assign support ticket: %w", err)
		}
	}

	if internalNote != nil && strings.TrimSpace(*internalNote) != "" {
		if _, err := tx.ExecContext(ctx, deleteInternalNotes, ticketID); err != nil {
			return fmt.Errorf("delete internal notes: %w", err)
		}
		if _, err := tx.ExecContext(ctx, insertInternalNote, utils.GenerateUUID(), ticketID, adminUserID, strings.TrimSpace(*internalNote)); err != nil {
			return fmt.Errorf("insert internal note: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit support reply tx: %w", err)
	}
	return nil
}

func scanBusinessTicketRows(rows *sql.Rows) ([]businessDto.SupportTicketRecord, error) {
	items := make([]businessDto.SupportTicketRecord, 0)
	for rows.Next() {
		rec, err := scanBusinessTicketRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *rec)
	}
	if items == nil {
		items = []businessDto.SupportTicketRecord{}
	}
	return items, rows.Err()
}

func scanBusinessTicketRow(row scanner) (*businessDto.SupportTicketRecord, error) {
	var rec businessDto.SupportTicketRecord
	var resolution, campaignID, campaignTitle sql.NullString
	if err := row.Scan(
		&rec.ID,
		&rec.Subject,
		&rec.Message,
		&rec.IssueType,
		&rec.Priority,
		&rec.Status,
		&resolution,
		&campaignID,
		&campaignTitle,
		&rec.CreatedAt,
	); err != nil {
		return nil, err
	}
	if resolution.Valid {
		rec.ResolutionNote = &resolution.String
	}
	if campaignID.Valid {
		v := campaignID.String
		rec.RelatedCampaignID = &v
	}
	if campaignTitle.Valid {
		rec.RelatedCampaignTitle = &campaignTitle.String
	}
	return &rec, nil
}

func scanAdminListRow(row scanner) (*adminDto.SupportTicketRow, error) {
	var rec adminDto.SupportTicketRow
	var resolution, campaignID, campaignTitle sql.NullString
	if err := row.Scan(
		&rec.ID,
		&rec.Subject,
		&rec.Message,
		&rec.IssueType,
		&rec.Priority,
		&rec.Status,
		&resolution,
		&rec.SenderUserID,
		&rec.SenderName,
		&rec.SenderEmail,
		&rec.SenderRole,
		&campaignID,
		&campaignTitle,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	); err != nil {
		return nil, err
	}
	if resolution.Valid {
		rec.ResolutionNote = &resolution.String
	}
	if campaignID.Valid {
		v := campaignID.String
		rec.RelatedCampaignID = &v
	}
	if campaignTitle.Valid {
		rec.RelatedCampaignTitle = &campaignTitle.String
	}
	return &rec, nil
}

func scanAdminTicketRow(row scanner) (*adminDto.SupportTicketRow, error) {
	var rec adminDto.SupportTicketRow
	var resolution, internalNote, campaignID, campaignTitle sql.NullString
	if err := row.Scan(
		&rec.ID,
		&rec.Subject,
		&rec.Message,
		&rec.IssueType,
		&rec.Priority,
		&rec.Status,
		&resolution,
		&internalNote,
		&rec.SenderUserID,
		&rec.SenderName,
		&rec.SenderEmail,
		&rec.SenderRole,
		&campaignID,
		&campaignTitle,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	); err != nil {
		return nil, err
	}
	if resolution.Valid {
		rec.ResolutionNote = &resolution.String
	}
	if internalNote.Valid {
		rec.InternalNote = &internalNote.String
	}
	if campaignID.Valid {
		v := campaignID.String
		rec.RelatedCampaignID = &v
	}
	if campaignTitle.Valid {
		rec.RelatedCampaignTitle = &campaignTitle.String
	}
	return &rec, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func nullUUID(id *string) sql.NullString {
	if id == nil || strings.TrimSpace(*id) == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: strings.TrimSpace(*id), Valid: true}
}

func (s *support) ListInfluencerInviteOptions(ctx context.Context, influencerUserID string) ([]storage.SupportInviteOption, error) {
	rows, err := s.db.QueryContext(ctx, listInfluencerInviteOptions, influencerUserID)
	if err != nil {
		return nil, fmt.Errorf("list influencer invite options: %w", err)
	}
	defer rows.Close()

	items := make([]storage.SupportInviteOption, 0)
	for rows.Next() {
		var item storage.SupportInviteOption
		if err := rows.Scan(&item.ID, &item.CampaignTitle, &item.BusinessName); err != nil {
			return nil, fmt.Errorf("scan invite option: %w", err)
		}
		items = append(items, item)
	}
	if items == nil {
		items = []storage.SupportInviteOption{}
	}
	return items, rows.Err()
}

func (s *support) GetInviteCampaignID(ctx context.Context, influencerUserID, inviteID string) (string, error) {
	var campaignID string
	err := s.db.QueryRowContext(ctx, getInviteCampaignID, inviteID, influencerUserID).Scan(&campaignID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrInviteNotFound
		}
		return "", fmt.Errorf("get invite campaign id: %w", err)
	}
	return campaignID, nil
}
