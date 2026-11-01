package admin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrBusinessVerificationNotFound = errors.New("business verification not found")
	ErrBusinessAlreadyVerified      = errors.New("business already verified")
	ErrBusinessNotVerified          = errors.New("business not verified")
	ErrBusinessNotPendingReview     = errors.New("business is not awaiting verification")
)

func (a *admin) ListBusinessVerifications(ctx context.Context, query adminDto.ListBusinessVerificationsQuery) ([]adminDto.BusinessVerificationListItem, int, error) {
	query = query.Normalize()
	offset := (query.Page - 1) * query.Limit
	search := strings.TrimSpace(query.Search)

	rows, err := a.db.QueryContext(ctx, listBusinessVerifications, query.Status, search, query.Limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list business verifications: %w", err)
	}
	defer rows.Close()

	var total int
	items := make([]adminDto.BusinessVerificationListItem, 0)
	for rows.Next() {
		var item adminDto.BusinessVerificationListItem
		var industry, city, logo sql.NullString
		var submittedAt sql.NullTime
		if err := rows.Scan(
			&item.UserID,
			&item.BusinessName,
			&item.OwnerFullName,
			&industry,
			&city,
			&logo,
			&item.Email,
			&item.IsVerified,
			&item.VerificationStatus,
			&submittedAt,
			&item.RegisteredAt,
			&total,
		); err != nil {
			return nil, 0, fmt.Errorf("scan business verification: %w", err)
		}
		item.Industry = industry.String
		item.City = city.String
		item.LogoURL = nullStr(logo)
		item.SubmittedAt = nullTime(submittedAt)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (a *admin) GetBusinessVerification(ctx context.Context, userID string) (*adminDto.BusinessVerificationRow, error) {
	var rec adminDto.BusinessVerificationRow
	var legalName, tin, industry, bizType, bizSize, city, region, address, website sql.NullString
	var instagram, tiktok, facebook, linkedin, description, logo, note sql.NullString
	var submittedAt, verifiedAt sql.NullTime

	err := a.db.QueryRowContext(ctx, getBusinessVerification, userID).Scan(
		&rec.UserID,
		&rec.Email,
		&rec.EmailVerified,
		&rec.PhoneNumber,
		&rec.PhoneVerified,
		&rec.Status,
		&rec.BusinessName,
		&legalName,
		&rec.OwnerFullName,
		&tin,
		&industry,
		&bizType,
		&bizSize,
		&city,
		&region,
		&address,
		&website,
		&instagram,
		&tiktok,
		&facebook,
		&linkedin,
		&description,
		&logo,
		&rec.IsVerified,
		&rec.VerificationStatus,
		&note,
		&submittedAt,
		&verifiedAt,
		&rec.RegisteredAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrBusinessVerificationNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get business verification: %w", err)
	}

	rec.LegalBusinessName = nullStr(legalName)
	rec.TINNumber = nullStr(tin)
	rec.Industry = industry.String
	rec.BusinessType = nullStr(bizType)
	rec.BusinessSize = nullStr(bizSize)
	rec.City = city.String
	rec.Region = nullStr(region)
	rec.Address = nullStr(address)
	rec.WebsiteURL = nullStr(website)
	rec.SocialLinks = adminDto.BusinessSocialLinks{
		Instagram: instagram.String,
		TikTok:    tiktok.String,
		Facebook:  facebook.String,
		LinkedIn:  linkedin.String,
	}
	rec.Description = nullStr(description)
	rec.LogoURL = nullStr(logo)
	rec.VerificationNote = nullStr(note)
	rec.SubmittedAt = nullTime(submittedAt)
	rec.VerifiedAt = nullTime(verifiedAt)
	return &rec, nil
}

// VerifyBusiness approves a business that is awaiting review (or was
// rejected before) and returns when it was verified.
func (a *admin) VerifyBusiness(ctx context.Context, userID, adminUserID, note string) (time.Time, error) {
	var verifiedAt time.Time
	err := a.db.QueryRowContext(ctx, verifyBusinessProfile, userID, adminUserID, nullReason(note)).Scan(&verifiedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return time.Time{}, a.verificationConflict(ctx, userID, ErrBusinessAlreadyVerified)
	}
	if err != nil {
		return time.Time{}, fmt.Errorf("verify business: %w", err)
	}
	return verifiedAt, nil
}

// RejectBusinessVerification rejects a pending business.
func (a *admin) RejectBusinessVerification(ctx context.Context, userID, adminUserID, reason string) error {
	return a.setRejected(ctx, userID, adminUserID, reason, "pending", ErrBusinessNotPendingReview)
}

// RevokeBusinessVerification takes verification away from a verified business.
func (a *admin) RevokeBusinessVerification(ctx context.Context, userID, adminUserID, reason string) error {
	return a.setRejected(ctx, userID, adminUserID, reason, "verified", ErrBusinessNotVerified)
}

func (a *admin) setRejected(ctx context.Context, userID, adminUserID, reason, fromStatus string, conflict error) error {
	result, err := a.db.ExecContext(ctx, rejectBusinessVerification, userID, adminUserID, reason, fromStatus)
	if err != nil {
		return fmt.Errorf("reject business verification: %w", err)
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return a.verificationConflict(ctx, userID, conflict)
	}
	return nil
}

// verificationConflict distinguishes "no such business" from "wrong state".
func (a *admin) verificationConflict(ctx context.Context, userID string, conflict error) error {
	var status string
	err := a.db.QueryRowContext(ctx, getBusinessVerificationStatus, userID).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrBusinessVerificationNotFound
	}
	if err != nil {
		return fmt.Errorf("get business verification status: %w", err)
	}
	if status == "verified" && conflict != ErrBusinessNotVerified {
		return ErrBusinessAlreadyVerified
	}
	if status == "unsubmitted" {
		return ErrBusinessNotPendingReview
	}
	return conflict
}

func nullStr(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}

func nullTime(v sql.NullTime) *time.Time {
	if !v.Valid {
		return nil
	}
	t := v.Time
	return &t
}
