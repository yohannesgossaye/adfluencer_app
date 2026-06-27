package influencer

import (
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	influencerDto "adinfluence_service_app/internal/constant/model/dto/influencer"
	"adinfluence_service_app/platform/utils"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type influencer struct {
	db     *sql.DB
	logger utils.Logger
}

func InitInfluencerPersistence(db *sql.DB, log utils.Logger) *influencer {
	return &influencer{db: db, logger: log}
}

func (p *influencer) GetProfileByUserID(ctx context.Context, userID string) (*influencerDto.ProfileRecord, error) {
	rec, err := p.scanProfileRow(ctx, userID)
	if err != nil {
		return nil, err
	}

	niches, err := p.GetCategoriesForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	platforms, err := p.GetPlatformsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	rec.Niches = niches
	rec.Platforms = platforms
	if rec.Languages == nil {
		rec.Languages = []string{}
	}

	portfolio, err := p.ListPortfolioItems(ctx, userID)
	if err != nil {
		return nil, err
	}
	rec.PortfolioItems = portfolio
	if rec.PortfolioItems == nil {
		rec.PortfolioItems = []influencerDto.PortfolioItem{}
	}

	return rec, nil
}

func (p *influencer) scanProfileRow(ctx context.Context, userID string) (*influencerDto.ProfileRecord, error) {
	var rec influencerDto.ProfileRecord
	var username, profileImage, bio, city, region, audience sql.NullString
	var languages []byte
	var socialRaw []byte
	var engagement, minPrice, maxPrice sql.NullFloat64
	var approvedAt sql.NullTime

	err := p.db.QueryRowContext(ctx, getInfluencerProfileFull, userID).Scan(
		&rec.ID,
		&rec.UserID,
		&rec.FullName,
		&rec.CreatorName,
		&username,
		&profileImage,
		&bio,
		&city,
		&region,
		&languages,
		&socialRaw,
		&rec.InstagramFollowers,
		&rec.TikTokFollowers,
		&rec.YouTubeSubscribers,
		&rec.TotalFollowers,
		&engagement,
		&minPrice,
		&maxPrice,
		&audience,
		&rec.IsVerified,
		&approvedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, influencerDto.ErrProfileNotFound
		}
		return nil, fmt.Errorf("get influencer profile: %w", err)
	}

	if username.Valid {
		rec.Username = &username.String
	}
	if profileImage.Valid {
		rec.ProfileImageURL = &profileImage.String
	}
	if bio.Valid {
		rec.Bio = &bio.String
	}
	if city.Valid {
		rec.City = &city.String
	}
	if region.Valid {
		rec.Region = &region.String
	}
	if audience.Valid {
		rec.AudienceDescription = &audience.String
	}
	if engagement.Valid {
		v := engagement.Float64
		rec.EngagementRate = &v
	}
	if minPrice.Valid {
		v := minPrice.Float64
		rec.MinPriceETB = &v
	}
	if maxPrice.Valid {
		v := maxPrice.Float64
		rec.MaxPriceETB = &v
	}
	if approvedAt.Valid {
		t := approvedAt.Time
		rec.ApprovedAt = &t
	}

	rec.Languages = parseTextArray(languages)
	rec.SocialLinks = parseSocialLinks(socialRaw)

	return &rec, nil
}

func (p *influencer) UpdateProfile(ctx context.Context, input influencerDto.UpdateProfileInput) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin profile update: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, updateInfluencerProfile,
		input.UserID,
		input.FullName,
		input.CreatorName,
		input.Bio,
		input.City,
		nullString(input.Region),
		formatTextArray(input.Languages),
		input.SocialLinksJSON,
		input.InstagramFollowers,
		input.TikTokFollowers,
		input.YouTubeSubscribers,
		input.TotalFollowers,
		nullFloat(input.EngagementRate),
		nullFloat(input.MinPriceETB),
		nullFloat(input.MaxPriceETB),
		nullString(input.AudienceDescription),
	)
	if err != nil {
		return fmt.Errorf("update influencer profile: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update influencer profile rows: %w", err)
	}
	if rows == 0 {
		return influencerDto.ErrProfileNotFound
	}

	if _, err := tx.ExecContext(ctx, deleteProfileCategories, input.UserID); err != nil {
		return fmt.Errorf("clear profile categories: %w", err)
	}
	for _, id := range input.CategoryIDs {
		if _, err := tx.ExecContext(ctx, insertProfileCategory, input.UserID, id); err != nil {
			return fmt.Errorf("insert profile category: %w", err)
		}
	}

	if _, err := tx.ExecContext(ctx, deleteProfilePlatforms, input.UserID); err != nil {
		return fmt.Errorf("clear profile platforms: %w", err)
	}
	for _, platform := range input.Platforms {
		if _, err := tx.ExecContext(ctx, insertProfilePlatform, input.UserID, platform); err != nil {
			return fmt.Errorf("insert profile platform: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit profile update: %w", err)
	}
	return nil
}

func (p *influencer) ResolveCategoryIDs(ctx context.Context, namesOrSlugs []string) ([]influencerDto.CategoryResponse, error) {
	keys := dedupeLowerKeys(namesOrSlugs)
	if len(keys) == 0 {
		return nil, influencerDto.ErrUnknownCategory
	}

	rows, err := p.db.QueryContext(ctx, resolveCategories, formatTextArray(keys))
	if err != nil {
		return nil, fmt.Errorf("resolve categories: %w", err)
	}
	defer rows.Close()

	bySlug := make(map[string]influencerDto.CategoryResponse)
	byName := make(map[string]influencerDto.CategoryResponse)
	for rows.Next() {
		var c influencerDto.CategoryResponse
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}
		bySlug[strings.ToLower(c.Slug)] = c
		byName[strings.ToLower(c.Name)] = c
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	out := make([]influencerDto.CategoryResponse, 0, len(keys))
	for _, key := range keys {
		c, ok := bySlug[key]
		if !ok {
			c, ok = byName[key]
		}
		if !ok {
			return nil, influencerDto.ErrUnknownCategory
		}
		out = append(out, c)
	}
	return out, nil
}

func (p *influencer) UpdateProfileImageURL(ctx context.Context, userID, imageURL string) error {
	res, err := p.db.ExecContext(ctx, updateProfileImageURL, userID, imageURL)
	if err != nil {
		return fmt.Errorf("update profile image url: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update profile image url rows affected: %w", err)
	}
	if n == 0 {
		return influencerDto.ErrProfileNotFound
	}
	return nil
}

func (p *influencer) CountPortfolioItems(ctx context.Context, userID string) (int, error) {
	var count int
	if err := p.db.QueryRowContext(ctx, countPortfolioItems, userID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count portfolio items: %w", err)
	}
	return count, nil
}

func (p *influencer) ListPortfolioItems(ctx context.Context, userID string) ([]influencerDto.PortfolioItem, error) {
	rows, err := p.db.QueryContext(ctx, listPortfolioItems, userID)
	if err != nil {
		return nil, fmt.Errorf("list portfolio items: %w", err)
	}
	defer rows.Close()

	var out []influencerDto.PortfolioItem
	for rows.Next() {
		item, err := scanPortfolioItem(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (p *influencer) GetPortfolioItem(ctx context.Context, userID, itemID string) (*influencerDto.PortfolioItemRecord, error) {
	var rec influencerDto.PortfolioItemRecord
	var objectKey string
	err := p.db.QueryRowContext(ctx, getPortfolioItem, itemID, userID).Scan(
		&rec.ID,
		&rec.FileURL,
		&objectKey,
		&rec.MediaType,
		&rec.MimeType,
		&rec.FileSizeBytes,
		&rec.SortOrder,
		&rec.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, influencerDto.ErrPortfolioItemNotFound
		}
		return nil, fmt.Errorf("get portfolio item: %w", err)
	}
	rec.ObjectKey = objectKey
	return &rec, nil
}

func (p *influencer) CreatePortfolioItem(ctx context.Context, input influencerDto.CreatePortfolioItemInput) (*influencerDto.PortfolioItem, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin portfolio insert: %w", err)
	}
	defer tx.Rollback()

	var count int
	if err := tx.QueryRowContext(ctx, countPortfolioItems, input.UserID).Scan(&count); err != nil {
		return nil, fmt.Errorf("count portfolio items: %w", err)
	}
	if count >= 3 {
		return nil, influencerDto.ErrPortfolioLimitReached
	}

	sortOrder := count

	var createdAt time.Time
	err = tx.QueryRowContext(ctx, insertPortfolioItem,
		input.ID,
		input.UserID,
		input.FileURL,
		input.ObjectKey,
		input.MediaType,
		input.MimeType,
		input.FileSizeBytes,
		sortOrder,
	).Scan(&sortOrder, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("insert portfolio item: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit portfolio insert: %w", err)
	}

	return &influencerDto.PortfolioItem{
		ID:            input.ID,
		FileURL:       input.FileURL,
		MediaType:     input.MediaType,
		MimeType:      input.MimeType,
		FileSizeBytes: input.FileSizeBytes,
		SortOrder:     sortOrder,
		CreatedAt:     createdAt,
	}, nil
}

func (p *influencer) DeletePortfolioItem(ctx context.Context, userID, itemID string) error {
	res, err := p.db.ExecContext(ctx, deletePortfolioItem, itemID, userID)
	if err != nil {
		return fmt.Errorf("delete portfolio item: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete portfolio item rows: %w", err)
	}
	if n == 0 {
		return influencerDto.ErrPortfolioItemNotFound
	}
	return nil
}

type portfolioScanner interface {
	Scan(dest ...any) error
}

func scanPortfolioItem(s portfolioScanner) (influencerDto.PortfolioItem, error) {
	var item influencerDto.PortfolioItem
	if err := s.Scan(
		&item.ID,
		&item.FileURL,
		&item.MediaType,
		&item.MimeType,
		&item.FileSizeBytes,
		&item.SortOrder,
		&item.CreatedAt,
	); err != nil {
		return influencerDto.PortfolioItem{}, fmt.Errorf("scan portfolio item: %w", err)
	}
	return item, nil
}

func dedupeLowerKeys(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, v := range values {
		key := strings.ToLower(strings.TrimSpace(v))
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, key)
	}
	return out
}

func (p *influencer) GetCategoriesForUser(ctx context.Context, userID string) ([]influencerDto.CategoryResponse, error) {
	rows, err := p.db.QueryContext(ctx, getCategoriesForUser, userID)
	if err != nil {
		return nil, fmt.Errorf("get profile categories: %w", err)
	}
	defer rows.Close()

	var out []influencerDto.CategoryResponse
	for rows.Next() {
		var c influencerDto.CategoryResponse
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug); err != nil {
			return nil, fmt.Errorf("scan profile category: %w", err)
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (p *influencer) GetPlatformsForUser(ctx context.Context, userID string) ([]string, error) {
	rows, err := p.db.QueryContext(ctx, getPlatformsForUser, userID)
	if err != nil {
		return nil, fmt.Errorf("get profile platforms: %w", err)
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var platform string
		if err := rows.Scan(&platform); err != nil {
			return nil, fmt.Errorf("scan profile platform: %w", err)
		}
		out = append(out, platform)
	}
	return out, rows.Err()
}

func (p *influencer) SubmitProfile(ctx context.Context, userID string) error {
	result, err := p.db.ExecContext(ctx, submitInfluencerProfile, userID)
	if err != nil {
		return fmt.Errorf("submit profile: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("submit profile rows: %w", err)
	}
	if rows == 0 {
		return influencerDto.ErrInvalidStatusForSubmit
	}
	return nil
}

func (p *influencer) CreateAuditLog(ctx context.Context, input authDto.CreateAuditLogInput) error {
	_, err := p.db.ExecContext(ctx, insertInfluencerAuditLog,
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

func parseSocialLinks(raw []byte) influencerDto.SocialLinks {
	if len(raw) == 0 {
		return influencerDto.SocialLinks{}
	}
	var m map[string]string
	if err := json.Unmarshal(raw, &m); err != nil {
		return influencerDto.SocialLinks{}
	}
	return influencerDto.SocialLinks{
		TikTok:    m["tiktok"],
		Instagram: m["instagram"],
		YouTube:   m["youtube"],
	}
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
