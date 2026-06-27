package business

import (
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	"adinfluence_service_app/platform/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type business struct {
	db     *sql.DB
	logger utils.Logger
}

func InitBusinessPersistence(db *sql.DB, log utils.Logger) *business {
	return &business{db: db, logger: log}
}

func (p *business) GetProfileByUserID(ctx context.Context, userID string) (*businessDto.ProfileRecord, error) {
	var rec businessDto.ProfileRecord
	var legalName, logo, description, industry, bizType, bizSize sql.NullString
	var city, region, address, tin, website sql.NullString
	var instagram, tiktok, facebook, linkedin sql.NullString
	var platforms []byte
	var budget sql.NullFloat64

	err := p.db.QueryRowContext(ctx, getBusinessProfileFull, userID).Scan(
		&rec.ID,
		&rec.UserID,
		&rec.OwnerFullName,
		&rec.BusinessName,
		&legalName,
		&logo,
		&description,
		&industry,
		&bizType,
		&bizSize,
		&city,
		&region,
		&address,
		&tin,
		&website,
		&instagram,
		&tiktok,
		&facebook,
		&linkedin,
		&platforms,
		&budget,
		&rec.IsVerified,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, businessDto.ErrProfileNotFound
		}
		return nil, fmt.Errorf("get business profile: %w", err)
	}

	rec.LegalBusinessName = nullStringPtr(legalName)
	rec.LogoURL = nullStringPtr(logo)
	rec.Description = nullStringPtr(description)
	rec.Industry = nullStringPtr(industry)
	rec.BusinessType = nullStringPtr(bizType)
	rec.BusinessSize = nullStringPtr(bizSize)
	rec.City = nullStringPtr(city)
	rec.Region = nullStringPtr(region)
	rec.Address = nullStringPtr(address)
	rec.TINNumber = nullStringPtr(tin)
	rec.WebsiteURL = nullStringPtr(website)
	rec.SocialLinks = businessDto.SocialLinks{
		Instagram: instagram.String,
		TikTok:    tiktok.String,
		Facebook:  facebook.String,
		LinkedIn:  linkedin.String,
	}
	rec.PreferredPlatforms = parseTextArray(platforms)
	if rec.PreferredPlatforms == nil {
		rec.PreferredPlatforms = []string{}
	}
	if budget.Valid {
		v := budget.Float64
		rec.MonthlyMarketingBudget = &v
	}

	return &rec, nil
}

func (p *business) UpdateProfile(ctx context.Context, input businessDto.UpdateProfileInput) error {
	res, err := p.db.ExecContext(ctx, updateBusinessProfile,
		input.UserID,
		input.OwnerFullName,
		input.BusinessName,
		nullString(input.LegalBusinessName),
		nullString(input.Description),
		nullString(input.Industry),
		nullString(input.BusinessType),
		nullString(input.BusinessSize),
		nullString(input.City),
		nullString(input.Region),
		nullString(input.Address),
		nullString(input.TINNumber),
		nullString(input.WebsiteURL),
		nullString(input.InstagramURL),
		nullString(input.TikTokURL),
		nullString(input.FacebookURL),
		nullString(input.LinkedInURL),
		formatTextArray(input.PreferredPlatforms),
		nullFloat(input.MonthlyMarketingBudget),
	)
	if err != nil {
		return fmt.Errorf("update business profile: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update business profile rows: %w", err)
	}
	if n == 0 {
		return businessDto.ErrProfileNotFound
	}
	return nil
}

func (p *business) UpdateLogoURL(ctx context.Context, userID, logoURL string) error {
	res, err := p.db.ExecContext(ctx, updateBusinessLogoURL, userID, logoURL)
	if err != nil {
		return fmt.Errorf("update business logo url: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update business logo url rows: %w", err)
	}
	if n == 0 {
		return businessDto.ErrProfileNotFound
	}
	return nil
}

func (p *business) CompleteProfile(ctx context.Context, userID string) error {
	res, err := p.db.ExecContext(ctx, completeBusinessProfile, userID)
	if err != nil {
		return fmt.Errorf("complete business profile: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("complete business profile rows: %w", err)
	}
	if n == 0 {
		return businessDto.ErrInvalidStatusForComplete
	}
	return nil
}

func (p *business) CreateAuditLog(ctx context.Context, input authDto.CreateAuditLogInput) error {
	_, err := p.db.ExecContext(ctx, insertBusinessAuditLog,
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

func parseTextArray(raw []byte) []string {
	s := strings.TrimSpace(string(raw))
	if s == "" || s == "{}" {
		return []string{}
	}
	s = strings.TrimPrefix(s, "{")
	s = strings.TrimSuffix(s, "}")
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		part = strings.Trim(part, `"`)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func formatTextArray(vals []string) string {
	if len(vals) == 0 {
		return "{}"
	}
	var b strings.Builder
	b.WriteByte('{')
	for i, v := range vals {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(`"` + strings.ReplaceAll(v, `"`, `\"`) + `"`)
	}
	b.WriteByte('}')
	return b.String()
}

func nullString(v *string) sql.NullString {
	if v == nil || strings.TrimSpace(*v) == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: strings.TrimSpace(*v), Valid: true}
}

func nullStringPtr(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}

func nullFloat(v *float64) sql.NullFloat64 {
	if v == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: *v, Valid: true}
}

func nullUUID(id *string) sql.NullString {
	if id == nil || strings.TrimSpace(*id) == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: *id, Valid: true}
}
