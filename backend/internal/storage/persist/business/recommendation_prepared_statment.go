package business

const (
	listEligibleInfluencers = `
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
		COALESCE((
			SELECT array_agg(ipc.category_id::text)
			FROM influencer_profile_categories ipc
			WHERE ipc.influencer_user_id = u.id
		), '{}'),
		COALESCE((
			SELECT array_agg(c.name ORDER BY c.name)
			FROM influencer_profile_categories ipc
			JOIN categories c ON c.id = ipc.category_id
			WHERE ipc.influencer_user_id = u.id
		), '{}'),
		COALESCE((
			SELECT array_agg(ipp.platform::text ORDER BY ipp.platform)
			FROM influencer_profile_platforms ipp
			WHERE ipp.influencer_user_id = u.id
		), '{}')
	FROM users u
	JOIN influencer_profiles ip ON ip.user_id = u.id
	WHERE u.role = 'influencer'
		AND u.status = 'active'
		AND ip.is_verified = TRUE
	`

	getInfluencerCategoryIDs = `
	SELECT category_id::text
	FROM influencer_profile_categories
	WHERE influencer_user_id = $1
	`

	getInfluencerCategoryNames = `
	SELECT c.name
	FROM influencer_profile_categories ipc
	JOIN categories c ON c.id = ipc.category_id
	WHERE ipc.influencer_user_id = $1
	ORDER BY c.name
	`

	deleteRecommendedRows = `
	DELETE FROM campaign_recommendations
	WHERE campaign_id = $1 AND status = 'recommended'
	`

	upsertRecommendation = `
	INSERT INTO campaign_recommendations (
		id, campaign_id, influencer_user_id, match_score, score_breakdown, rank_position, status
	) VALUES ($1, $2, $3, $4, $5::jsonb, $6, 'recommended')
	ON CONFLICT (campaign_id, influencer_user_id) DO UPDATE SET
		match_score = EXCLUDED.match_score,
		score_breakdown = EXCLUDED.score_breakdown,
		rank_position = EXCLUDED.rank_position
	WHERE campaign_recommendations.status = 'recommended'
	`

	listCampaignRecommendationsMatch = `
	SELECT
		cr.influencer_user_id,
		ip.creator_name,
		ip.profile_image_url,
		COALESCE(primary_niche.name, ''),
		ip.total_followers,
		ip.engagement_rate,
		ip.min_price_etb,
		ip.max_price_etb,
		cr.match_score,
		cr.score_breakdown,
		ip.is_verified,
		inv.invite_status,
		cr.rank_position,
		COALESCE((
			SELECT array_agg(ipp.platform::text ORDER BY ipp.platform)
			FROM influencer_profile_platforms ipp
			WHERE ipp.influencer_user_id = cr.influencer_user_id
		), '{}') AS platforms,
		COUNT(*) OVER()::int AS total_count
	FROM campaign_recommendations cr
	JOIN influencer_profiles ip ON ip.user_id = cr.influencer_user_id
	LEFT JOIN LATERAL (
		SELECT c.name
		FROM influencer_profile_categories ipc
		JOIN categories c ON c.id = ipc.category_id
		WHERE ipc.influencer_user_id = cr.influencer_user_id
		ORDER BY c.name
		LIMIT 1
	) primary_niche ON true
	LEFT JOIN LATERAL (
		SELECT ci.status::text AS invite_status
		FROM campaign_invites ci
		WHERE ci.campaign_id = cr.campaign_id
			AND ci.influencer_user_id = cr.influencer_user_id
			AND ci.status IN ('pending', 'accepted')
		ORDER BY ci.created_at DESC
		LIMIT 1
	) inv ON true
	WHERE cr.campaign_id = $1
		AND cr.match_score >= $2
		AND ($3 = '' OR ip.creator_name ILIKE '%' || $3 || '%')
		AND ($4 = '' OR EXISTS (
			SELECT 1
			FROM influencer_profile_categories ipc
			JOIN categories c ON c.id = ipc.category_id
			WHERE ipc.influencer_user_id = cr.influencer_user_id
				AND (LOWER(c.name) = LOWER($4) OR LOWER(c.slug) = LOWER($4))
		))
		AND ($5 = '' OR EXISTS (
			SELECT 1
			FROM influencer_profile_platforms ipp
			WHERE ipp.influencer_user_id = cr.influencer_user_id
				AND ipp.platform::text = LOWER($5)
		))
	ORDER BY cr.rank_position ASC NULLS LAST, cr.match_score DESC
	LIMIT $6 OFFSET $7
	`

	listCampaignRecommendationsFollowers = `
	SELECT
		cr.influencer_user_id,
		ip.creator_name,
		ip.profile_image_url,
		COALESCE(primary_niche.name, ''),
		ip.total_followers,
		ip.engagement_rate,
		ip.min_price_etb,
		ip.max_price_etb,
		cr.match_score,
		cr.score_breakdown,
		ip.is_verified,
		inv.invite_status,
		cr.rank_position,
		COALESCE((
			SELECT array_agg(ipp.platform::text ORDER BY ipp.platform)
			FROM influencer_profile_platforms ipp
			WHERE ipp.influencer_user_id = cr.influencer_user_id
		), '{}') AS platforms,
		COUNT(*) OVER()::int AS total_count
	FROM campaign_recommendations cr
	JOIN influencer_profiles ip ON ip.user_id = cr.influencer_user_id
	LEFT JOIN LATERAL (
		SELECT c.name
		FROM influencer_profile_categories ipc
		JOIN categories c ON c.id = ipc.category_id
		WHERE ipc.influencer_user_id = cr.influencer_user_id
		ORDER BY c.name
		LIMIT 1
	) primary_niche ON true
	LEFT JOIN LATERAL (
		SELECT ci.status::text AS invite_status
		FROM campaign_invites ci
		WHERE ci.campaign_id = cr.campaign_id
			AND ci.influencer_user_id = cr.influencer_user_id
			AND ci.status IN ('pending', 'accepted')
		ORDER BY ci.created_at DESC
		LIMIT 1
	) inv ON true
	WHERE cr.campaign_id = $1
		AND cr.match_score >= $2
		AND ($3 = '' OR ip.creator_name ILIKE '%' || $3 || '%')
		AND ($4 = '' OR EXISTS (
			SELECT 1
			FROM influencer_profile_categories ipc
			JOIN categories c ON c.id = ipc.category_id
			WHERE ipc.influencer_user_id = cr.influencer_user_id
				AND (LOWER(c.name) = LOWER($4) OR LOWER(c.slug) = LOWER($4))
		))
		AND ($5 = '' OR EXISTS (
			SELECT 1
			FROM influencer_profile_platforms ipp
			WHERE ipp.influencer_user_id = cr.influencer_user_id
				AND ipp.platform::text = LOWER($5)
		))
	ORDER BY ip.total_followers DESC, cr.match_score DESC
	LIMIT $6 OFFSET $7
	`

	listCampaignRecommendationsPriceLow = `
	SELECT
		cr.influencer_user_id,
		ip.creator_name,
		ip.profile_image_url,
		COALESCE(primary_niche.name, ''),
		ip.total_followers,
		ip.engagement_rate,
		ip.min_price_etb,
		ip.max_price_etb,
		cr.match_score,
		cr.score_breakdown,
		ip.is_verified,
		inv.invite_status,
		cr.rank_position,
		COALESCE((
			SELECT array_agg(ipp.platform::text ORDER BY ipp.platform)
			FROM influencer_profile_platforms ipp
			WHERE ipp.influencer_user_id = cr.influencer_user_id
		), '{}') AS platforms,
		COUNT(*) OVER()::int AS total_count
	FROM campaign_recommendations cr
	JOIN influencer_profiles ip ON ip.user_id = cr.influencer_user_id
	LEFT JOIN LATERAL (
		SELECT c.name
		FROM influencer_profile_categories ipc
		JOIN categories c ON c.id = ipc.category_id
		WHERE ipc.influencer_user_id = cr.influencer_user_id
		ORDER BY c.name
		LIMIT 1
	) primary_niche ON true
	LEFT JOIN LATERAL (
		SELECT ci.status::text AS invite_status
		FROM campaign_invites ci
		WHERE ci.campaign_id = cr.campaign_id
			AND ci.influencer_user_id = cr.influencer_user_id
			AND ci.status IN ('pending', 'accepted')
		ORDER BY ci.created_at DESC
		LIMIT 1
	) inv ON true
	WHERE cr.campaign_id = $1
		AND cr.match_score >= $2
		AND ($3 = '' OR ip.creator_name ILIKE '%' || $3 || '%')
		AND ($4 = '' OR EXISTS (
			SELECT 1
			FROM influencer_profile_categories ipc
			JOIN categories c ON c.id = ipc.category_id
			WHERE ipc.influencer_user_id = cr.influencer_user_id
				AND (LOWER(c.name) = LOWER($4) OR LOWER(c.slug) = LOWER($4))
		))
		AND ($5 = '' OR EXISTS (
			SELECT 1
			FROM influencer_profile_platforms ipp
			WHERE ipp.influencer_user_id = cr.influencer_user_id
				AND ipp.platform::text = LOWER($5)
		))
	ORDER BY ip.min_price_etb ASC NULLS LAST, cr.match_score DESC
	LIMIT $6 OFFSET $7
	`

	listCampaignRecommendationsPriceHigh = `
	SELECT
		cr.influencer_user_id,
		ip.creator_name,
		ip.profile_image_url,
		COALESCE(primary_niche.name, ''),
		ip.total_followers,
		ip.engagement_rate,
		ip.min_price_etb,
		ip.max_price_etb,
		cr.match_score,
		cr.score_breakdown,
		ip.is_verified,
		inv.invite_status,
		cr.rank_position,
		COALESCE((
			SELECT array_agg(ipp.platform::text ORDER BY ipp.platform)
			FROM influencer_profile_platforms ipp
			WHERE ipp.influencer_user_id = cr.influencer_user_id
		), '{}') AS platforms,
		COUNT(*) OVER()::int AS total_count
	FROM campaign_recommendations cr
	JOIN influencer_profiles ip ON ip.user_id = cr.influencer_user_id
	LEFT JOIN LATERAL (
		SELECT c.name
		FROM influencer_profile_categories ipc
		JOIN categories c ON c.id = ipc.category_id
		WHERE ipc.influencer_user_id = cr.influencer_user_id
		ORDER BY c.name
		LIMIT 1
	) primary_niche ON true
	LEFT JOIN LATERAL (
		SELECT ci.status::text AS invite_status
		FROM campaign_invites ci
		WHERE ci.campaign_id = cr.campaign_id
			AND ci.influencer_user_id = cr.influencer_user_id
			AND ci.status IN ('pending', 'accepted')
		ORDER BY ci.created_at DESC
		LIMIT 1
	) inv ON true
	WHERE cr.campaign_id = $1
		AND cr.match_score >= $2
		AND ($3 = '' OR ip.creator_name ILIKE '%' || $3 || '%')
		AND ($4 = '' OR EXISTS (
			SELECT 1
			FROM influencer_profile_categories ipc
			JOIN categories c ON c.id = ipc.category_id
			WHERE ipc.influencer_user_id = cr.influencer_user_id
				AND (LOWER(c.name) = LOWER($4) OR LOWER(c.slug) = LOWER($4))
		))
		AND ($5 = '' OR EXISTS (
			SELECT 1
			FROM influencer_profile_platforms ipp
			WHERE ipp.influencer_user_id = cr.influencer_user_id
				AND ipp.platform::text = LOWER($5)
		))
	ORDER BY ip.max_price_etb DESC NULLS LAST, cr.match_score DESC
	LIMIT $6 OFFSET $7
	`

	getRecommendationForCampaign = `
	SELECT match_score, score_breakdown, rank_position, status::text, created_at
	FROM campaign_recommendations
	WHERE campaign_id = $1 AND influencer_user_id = $2
	`

	isVerifiedActiveInfluencer = `
	SELECT EXISTS (
		SELECT 1
		FROM users u
		JOIN influencer_profiles ip ON ip.user_id = u.id
		WHERE u.id = $1
			AND u.role = 'influencer'
			AND (u.status = 'active' OR ip.is_verified = TRUE)
	)
	`
)
