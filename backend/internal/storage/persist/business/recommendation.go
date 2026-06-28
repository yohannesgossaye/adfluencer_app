package business

import (
	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
	"adinfluence_service_app/platform/utils"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func (p *business) ListMatchCandidates(ctx context.Context) ([]businessDto.MatchCandidate, error) {
	rows, err := p.db.QueryContext(ctx, listEligibleInfluencers)
	if err != nil {
		return nil, fmt.Errorf("list eligible influencers: %w", err)
	}
	defer rows.Close()

	candidates := make([]businessDto.MatchCandidate, 0)
	for rows.Next() {
		var candidate businessDto.MatchCandidate
		var profileImage, bio sql.NullString
		var engagement, minPrice, maxPrice sql.NullFloat64
		var categoryIDsRaw, categoryNamesRaw, platformsRaw []byte

		if err := rows.Scan(
			&candidate.UserID,
			&candidate.CreatorName,
			&profileImage,
			&bio,
			&candidate.TotalFollowers,
			&engagement,
			&minPrice,
			&maxPrice,
			&candidate.IsVerified,
			&categoryIDsRaw,
			&categoryNamesRaw,
			&platformsRaw,
		); err != nil {
			return nil, fmt.Errorf("scan eligible influencer: %w", err)
		}

		if profileImage.Valid {
			candidate.ProfileImageURL = &profileImage.String
		}
		if bio.Valid {
			candidate.Bio = &bio.String
		}
		if engagement.Valid {
			v := engagement.Float64
			candidate.EngagementRate = &v
		}
		if minPrice.Valid {
			v := minPrice.Float64
			candidate.MinPriceETB = &v
		}
		if maxPrice.Valid {
			v := maxPrice.Float64
			candidate.MaxPriceETB = &v
		}

		candidate.CategoryIDs = parseTextArray(categoryIDsRaw)
		candidate.CategoryNames = parseTextArray(categoryNamesRaw)
		candidate.Platforms = parseTextArray(platformsRaw)
		candidates = append(candidates, candidate)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate eligible influencers: %w", err)
	}

	return candidates, nil
}

func (p *business) ReplaceRecommendations(ctx context.Context, campaignID string, items []businessDto.SaveRecommendationInput) (int, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin replace recommendations: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, deleteRecommendedRows, campaignID); err != nil {
		return 0, fmt.Errorf("delete recommended rows: %w", err)
	}

	inserted := 0
	for _, item := range items {
		breakdownJSON, err := json.Marshal(item.ScoreBreakdown)
		if err != nil {
			return 0, fmt.Errorf("marshal score breakdown: %w", err)
		}

		id := item.ID
		if id == "" {
			id = utils.GenerateUUID()
		}

		result, err := tx.ExecContext(ctx, upsertRecommendation,
			id,
			campaignID,
			item.InfluencerUserID,
			item.MatchScore,
			breakdownJSON,
			item.RankPosition,
		)
		if err != nil {
			return 0, fmt.Errorf("upsert recommendation: %w", err)
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return 0, fmt.Errorf("upsert recommendation rows: %w", err)
		}
		if rows > 0 {
			inserted++
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit replace recommendations: %w", err)
	}
	return inserted, nil
}

func (p *business) ListRecommendations(ctx context.Context, campaignID string, query businessDto.ListRecommendationsQuery) ([]businessDto.RecommendationItem, int, error) {
	query = query.Normalize()
	offset := (query.Page - 1) * query.Limit
	search := strings.TrimSpace(query.Search)
	niche := strings.TrimSpace(query.Niche)
	platform := strings.ToLower(strings.TrimSpace(query.Platform))

	listSQL := listCampaignRecommendationsMatch
	switch query.Sort {
	case "followers":
		listSQL = listCampaignRecommendationsFollowers
	case "price_low":
		listSQL = listCampaignRecommendationsPriceLow
	case "price_high":
		listSQL = listCampaignRecommendationsPriceHigh
	}

	rows, err := p.db.QueryContext(ctx, listSQL,
		campaignID, query.MinScore, search, niche, platform, query.Limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list recommendations: %w", err)
	}
	defer rows.Close()

	var total int
	items, err := p.scanRecommendationRows(rows, &total)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (p *business) GetRecommendation(ctx context.Context, campaignID, influencerUserID string) (*businessDto.RecommendationRecord, error) {
	var rec businessDto.RecommendationRecord
	var breakdownRaw []byte
	var rank sql.NullInt64

	err := p.db.QueryRowContext(ctx, getRecommendationForCampaign, campaignID, influencerUserID).Scan(
		&rec.MatchScore,
		&breakdownRaw,
		&rank,
		&rec.Status,
		&rec.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get recommendation: %w", err)
	}

	if err := json.Unmarshal(breakdownRaw, &rec.ScoreBreakdown); err != nil {
		return nil, fmt.Errorf("unmarshal score breakdown: %w", err)
	}
	rec.CampaignID = campaignID
	rec.InfluencerUserID = influencerUserID
	if rank.Valid {
		rec.RankPosition = int(rank.Int64)
	}
	return &rec, nil
}

func (p *business) IsVerifiedActiveInfluencer(ctx context.Context, userID string) (bool, error) {
	var exists bool
	if err := p.db.QueryRowContext(ctx, isVerifiedActiveInfluencer, userID).Scan(&exists); err != nil {
		return false, fmt.Errorf("check verified influencer: %w", err)
	}
	return exists, nil
}

func (p *business) getInfluencerCategoryIDs(ctx context.Context, userID string) ([]string, error) {
	rows, err := p.db.QueryContext(ctx, getInfluencerCategoryIDs, userID)
	if err != nil {
		return nil, fmt.Errorf("get influencer category ids: %w", err)
	}
	defer rows.Close()

	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan category id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (p *business) getInfluencerCategoryNames(ctx context.Context, userID string) ([]string, error) {
	rows, err := p.db.QueryContext(ctx, getInfluencerCategoryNames, userID)
	if err != nil {
		return nil, fmt.Errorf("get influencer category names: %w", err)
	}
	defer rows.Close()

	names := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scan category name: %w", err)
		}
		names = append(names, name)
	}
	return names, rows.Err()
}

func (p *business) scanRecommendationRows(rows *sql.Rows, total *int) ([]businessDto.RecommendationItem, error) {
	items := make([]businessDto.RecommendationItem, 0)
	for rows.Next() {
		var item businessDto.RecommendationItem
		var profileImage sql.NullString
		var engagement, minPrice, maxPrice sql.NullFloat64
		var inviteStatus sql.NullString
		var breakdownRaw, platformsRaw []byte
		var rank sql.NullInt64

		if err := rows.Scan(
			&item.InfluencerUserID,
			&item.CreatorName,
			&profileImage,
			&item.PrimaryNiche,
			&item.TotalFollowers,
			&engagement,
			&minPrice,
			&maxPrice,
			&item.MatchScore,
			&breakdownRaw,
			&item.IsVerified,
			&inviteStatus,
			&rank,
			&platformsRaw,
			total,
		); err != nil {
			return nil, fmt.Errorf("scan recommendation item: %w", err)
		}

		if profileImage.Valid {
			item.ProfileImageURL = &profileImage.String
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
		if inviteStatus.Valid {
			v := inviteStatus.String
			item.InviteStatus = &v
		}
		if rank.Valid {
			item.RankPosition = int(rank.Int64)
		}
		if err := json.Unmarshal(breakdownRaw, &item.ScoreBreakdown); err != nil {
			return nil, fmt.Errorf("unmarshal recommendation breakdown: %w", err)
		}

		item.Platforms = parseTextArray(platformsRaw)
		if item.Platforms == nil {
			item.Platforms = []string{}
		}

		items = append(items, item)
	}
	return items, rows.Err()
}
