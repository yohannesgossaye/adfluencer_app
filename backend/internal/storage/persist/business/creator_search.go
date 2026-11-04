package business

import (
	"context"
	"database/sql"
	"fmt"

	businessDto "adinfluence_service_app/internal/constant/model/dto/business"
)

// maxCreatorCandidates caps one search. Ranking happens in Go over this set,
// which is plenty for the marketplace's size; revisit with full-text search
// once there are thousands of approved creators.
const maxCreatorCandidates = 1000

const searchCreators = `
SELECT
	u.id,
	ip.creator_name,
	ip.profile_image_url,
	ip.bio,
	ip.total_followers,
	ip.engagement_rate,
	ip.min_price_etb,
	ip.max_price_etb,
	ip.is_verified,
	COALESCE(ip.city, ''),
	ip.tiktok_followers,
	ip.instagram_followers,
	ip.youtube_subscribers,
	COALESCE((SELECT array_agg(c.id::text ORDER BY c.name) FROM influencer_profile_categories ipc JOIN categories c ON c.id = ipc.category_id WHERE ipc.influencer_user_id = u.id), '{}'),
	COALESCE((SELECT array_agg(c.name ORDER BY c.name) FROM influencer_profile_categories ipc JOIN categories c ON c.id = ipc.category_id WHERE ipc.influencer_user_id = u.id), '{}'),
	COALESCE((SELECT array_agg(c.slug ORDER BY c.name) FROM influencer_profile_categories ipc JOIN categories c ON c.id = ipc.category_id WHERE ipc.influencer_user_id = u.id), '{}'),
	COALESCE((SELECT array_agg(ipp.platform::text ORDER BY ipp.platform) FROM influencer_profile_platforms ipp WHERE ipp.influencer_user_id = u.id), '{}')
FROM users u
JOIN influencer_profiles ip ON ip.user_id = u.id
WHERE u.role = 'influencer'
	AND u.status = 'active'
	AND ip.is_verified = TRUE
	AND (
		$1 = ''
		OR ip.creator_name ILIKE '%' || $1 || '%'
		OR ip.username ILIKE '%' || $1 || '%'
		OR ip.city ILIKE '%' || $1 || '%'
		OR ip.bio ILIKE '%' || $1 || '%'
		OR EXISTS (
			SELECT 1 FROM influencer_profile_categories ipc JOIN categories c ON c.id = ipc.category_id
			WHERE ipc.influencer_user_id = u.id AND (c.name ILIKE '%' || $1 || '%' OR c.slug ILIKE '%' || $1 || '%')
		)
	)
	AND ($2 = '' OR EXISTS (
		SELECT 1 FROM influencer_profile_categories ipc JOIN categories c ON c.id = ipc.category_id
		WHERE ipc.influencer_user_id = u.id AND c.slug = $2
	))
	AND ($3 = '' OR EXISTS (
		SELECT 1 FROM influencer_profile_platforms ipp
		WHERE ipp.influencer_user_id = u.id AND ipp.platform::text = $3
	))
	-- Within budget: the creator's starting price fits.
	AND ($4::numeric = 0 OR ip.min_price_etb IS NULL OR ip.min_price_etb <= $4::numeric)
	AND ($5 = '' OR ip.city ILIKE '%' || $5 || '%')
ORDER BY ip.total_followers DESC, u.id
LIMIT $6
`

// SearchCreators returns approved, active creators that pass the filters.
func (p *business) SearchCreators(ctx context.Context, f businessDto.CreatorFilter) ([]businessDto.CreatorRow, error) {
	rows, err := p.db.QueryContext(ctx, searchCreators, f.Q, f.Niche, f.Platform, f.MaxPrice, f.City, maxCreatorCandidates)
	if err != nil {
		return nil, fmt.Errorf("search creators: %w", err)
	}
	defer rows.Close()

	out := make([]businessDto.CreatorRow, 0)
	for rows.Next() {
		var r businessDto.CreatorRow
		var image, bio sql.NullString
		var engagement, minPrice, maxPrice sql.NullFloat64
		var ids, names, slugs, platforms []byte
		if err := rows.Scan(&r.UserID, &r.CreatorName, &image, &bio, &r.TotalFollowers, &engagement, &minPrice, &maxPrice, &r.IsVerified,
			&r.City, &r.TikTok, &r.Instagram, &r.YouTube, &ids, &names, &slugs, &platforms); err != nil {
			return nil, fmt.Errorf("scan creator: %w", err)
		}
		if image.Valid {
			r.ProfileImageURL = &image.String
		}
		if bio.Valid {
			r.Bio = &bio.String
		}
		if engagement.Valid {
			v := engagement.Float64
			r.EngagementRate = &v
		}
		if minPrice.Valid {
			v := minPrice.Float64
			r.MinPriceETB = &v
		}
		if maxPrice.Valid {
			v := maxPrice.Float64
			r.MaxPriceETB = &v
		}
		r.CategoryIDs = parseTextArray(ids)
		r.CategoryNames = parseTextArray(names)
		r.NicheSlugs = parseTextArray(slugs)
		r.Platforms = parseTextArray(platforms)
		out = append(out, r)
	}
	return out, rows.Err()
}
