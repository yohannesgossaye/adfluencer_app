package admin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"adinfluence_service_app/platform/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var ErrApplicationNotFound = errors.New("application not found")
var ErrInvalidStatusForReview = errors.New("application is not pending approval")

type admin struct {
	db     *sql.DB
	logger utils.Logger
}

func InitAdminPersistence(db *sql.DB, log utils.Logger) *admin {
	return &admin{db: db, logger: log}
}

func (a *admin) ListInfluencerApplications(ctx context.Context, query adminDto.ListApplicationsQuery) ([]adminDto.ApplicationRow, int, error) {
	query = query.Normalize()
	offset := (query.Page - 1) * query.Limit

	var total int
	if err := a.db.QueryRowContext(ctx, countInfluencerApplications, query.Status, query.Search).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count applications: %w", err)
	}

	rows, err := a.db.QueryContext(ctx, listInfluencerApplications, query.Status, query.Search, query.Limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list applications: %w", err)
	}
	defer rows.Close()

	items, err := scanApplicationRows(rows)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (a *admin) GetInfluencerApplication(ctx context.Context, userID string) (*adminDto.ApplicationRow, error) {
	row, err := scanApplicationRowFromScanner(a.db.QueryRowContext(ctx, getInfluencerApplication, userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrApplicationNotFound
		}
		return nil, fmt.Errorf("get application: %w", err)
	}
	return row, nil
}

func (a *admin) ApproveInfluencerApplication(ctx context.Context, userID, adminUserID string, adminNote *string) error {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin approve transaction: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, approveInfluencerApplication, userID)
	if err != nil {
		return fmt.Errorf("approve user: %w", err)
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return ErrInvalidStatusForReview
	}

	if _, err := tx.ExecContext(ctx, setInfluencerApproved, userID, adminUserID, nullString(adminNote)); err != nil {
		return fmt.Errorf("approve profile: %w", err)
	}

	return tx.Commit()
}

func (a *admin) RejectInfluencerApplication(ctx context.Context, userID, reason string, adminNote *string) error {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin reject transaction: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, rejectInfluencerApplication, userID, reason)
	if err != nil {
		return fmt.Errorf("reject user: %w", err)
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return ErrInvalidStatusForReview
	}

	if _, err := tx.ExecContext(ctx, setInfluencerRejected, userID, nullString(adminNote)); err != nil {
		return fmt.Errorf("reject profile: %w", err)
	}

	return tx.Commit()
}

func (a *admin) UpdateInfluencerAdminNote(ctx context.Context, userID string, adminNote string) error {
	result, err := a.db.ExecContext(ctx, updateInfluencerAdminNote, userID, adminNote)
	if err != nil {
		return fmt.Errorf("update admin note: %w", err)
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return ErrApplicationNotFound
	}
	return nil
}

func (a *admin) CreateAuditLog(ctx context.Context, input authDto.CreateAuditLogInput) error {
	_, err := a.db.ExecContext(ctx, insertAdminAuditLog,
		utils.GenerateUUID(),
		nullUUID(input.ActorUserID),
		nullUUID(input.TargetUserID),
		input.Action,
		input.EntityType,
		nullUUID(input.EntityID),
		input.Description,
		input.IPAddress,
		input.UserAgent,
	)
	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}
	return nil
}

func scanApplicationRows(rows *sql.Rows) ([]adminDto.ApplicationRow, error) {
	var out []adminDto.ApplicationRow
	for rows.Next() {
		row, err := scanApplicationRowFromRows(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *row)
	}
	return out, rows.Err()
}

func scanApplicationRowFromRows(rows *sql.Rows) (*adminDto.ApplicationRow, error) {
	return scanApplicationRowFromScanner(rows)
}

type scanner interface {
	Scan(dest ...any) error
}

func scanApplicationRowFromScanner(s scanner) (*adminDto.ApplicationRow, error) {
	var item adminDto.ApplicationRow
	var statusReason, profileImage, bio, city, audience, approvalNote sql.NullString
	var engagement, minPrice, maxPrice sql.NullFloat64

	err := s.Scan(
		&item.UserID,
		&item.Status,
		&statusReason,
		&item.Email,
		&item.PhoneNumber,
		&item.UpdatedAt,
		&item.CreatorName,
		&profileImage,
		&bio,
		&city,
		&item.TotalFollowers,
		&engagement,
		&minPrice,
		&maxPrice,
		&audience,
		&approvalNote,
		&item.PrimaryNiche,
		&item.PrimaryPlatform,
	)
	if err != nil {
		return nil, err
	}

	if statusReason.Valid {
		item.StatusReason = &statusReason.String
	}
	if profileImage.Valid {
		item.ProfileImageURL = &profileImage.String
	}
	if bio.Valid {
		item.Bio = &bio.String
	}
	if city.Valid {
		item.City = &city.String
	}
	if audience.Valid {
		item.AudienceDesc = &audience.String
	}
	if approvalNote.Valid {
		item.ApprovalNote = &approvalNote.String
	}
	if engagement.Valid {
		v := engagement.Float64
		item.EngagementRate = &v
	}
	if minPrice.Valid {
		v := minPrice.Float64
		item.MinPriceETB = &v
	}
	if maxPrice.Valid {
		v := maxPrice.Float64
		item.MaxPriceETB = &v
	}

	return &item, nil
}

func nullString(v *string) sql.NullString {
	if v == nil || strings.TrimSpace(*v) == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: strings.TrimSpace(*v), Valid: true}
}

func nullUUID(id *string) sql.NullString {
	if id == nil || strings.TrimSpace(*id) == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: *id, Valid: true}
}
