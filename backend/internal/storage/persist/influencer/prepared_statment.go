package influencer

const (
	getInfluencerProfileFull = `
	SELECT
		id,
		user_id,
		full_name,
		creator_name,
		username,
		profile_image_url,
		bio,
		city,
		region,
		languages,
		social_links,
		instagram_followers,
		tiktok_followers,
		youtube_subscribers,
		total_followers,
		engagement_rate,
		min_price_etb,
		max_price_etb,
		audience_description,
		is_verified,
		approved_at
	FROM influencer_profiles
	WHERE user_id = $1
	`

	updateInfluencerProfile = `
	UPDATE influencer_profiles
	SET
		full_name = $2,
		creator_name = $3,
		bio = $4,
		city = $5,
		region = $6,
		languages = $7::text[],
		social_links = $8,
		instagram_followers = $9,
		tiktok_followers = $10,
		youtube_subscribers = $11,
		total_followers = $12,
		engagement_rate = $13,
		min_price_etb = $14,
		max_price_etb = $15,
		audience_description = $16,
		updated_at = NOW()
	WHERE user_id = $1
	`

	lockInfluencerProfileForUpdate = `
	SELECT id
	FROM influencer_profiles
	WHERE user_id = $1
	FOR UPDATE
	`

	getInfluencerProfileIsVerified = `
	SELECT is_verified
	FROM influencer_profiles
	WHERE user_id = $1
	`

	resolveCategories = `
	SELECT id, name, slug
	FROM categories
	WHERE is_active = TRUE
		AND (
			LOWER(slug) = ANY($1::text[])
			OR LOWER(name) = ANY($1::text[])
		)
	`

	deleteProfileCategories = `
	DELETE FROM influencer_profile_categories
	WHERE influencer_user_id = $1
	`

	insertProfileCategory = `
	INSERT INTO influencer_profile_categories (influencer_user_id, category_id)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING
	`

	deleteProfilePlatforms = `
	DELETE FROM influencer_profile_platforms
	WHERE influencer_user_id = $1
	`

	insertProfilePlatform = `
	INSERT INTO influencer_profile_platforms (influencer_user_id, platform)
	VALUES ($1, $2::social_platform)
	ON CONFLICT DO NOTHING
	`

	getCategoriesForUser = `
	SELECT c.id, c.name, c.slug
	FROM influencer_profile_categories ipc
	JOIN categories c ON c.id = ipc.category_id
	WHERE ipc.influencer_user_id = $1
	ORDER BY c.name
	`

	getPlatformsForUser = `
	SELECT platform::text
	FROM influencer_profile_platforms
	WHERE influencer_user_id = $1
	ORDER BY platform
	`

	submitInfluencerProfile = `
	UPDATE users
	SET status = 'pending_approval', status_reason = NULL, updated_at = NOW()
	WHERE id = $1
		AND role = 'influencer'
		AND status IN ('pending_profile', 'rejected')
	`

	updateProfileImageURL = `
	UPDATE influencer_profiles
	SET profile_image_url = $2, updated_at = NOW()
	WHERE user_id = $1
	`

	insertInfluencerAuditLog = `
	INSERT INTO audit_logs (
		id,
		actor_user_id,
		target_user_id,
		action,
		entity_type,
		entity_id,
		description,
		ip_address,
		user_agent,
		created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
	`

	countPortfolioItems = `
	SELECT COUNT(*)::int
	FROM influencer_portfolio_items
	WHERE user_id = $1
	`

	listPortfolioItems = `
	SELECT id, file_url, media_type, mime_type, file_size_bytes, sort_order, created_at
	FROM influencer_portfolio_items
	WHERE user_id = $1
	ORDER BY sort_order, created_at
	`

	getPortfolioItem = `
	SELECT id, file_url, object_key, media_type, mime_type, file_size_bytes, sort_order, created_at
	FROM influencer_portfolio_items
	WHERE id = $1 AND user_id = $2
	`

	insertPortfolioItem = `
	INSERT INTO influencer_portfolio_items (
		id, user_id, file_url, object_key, media_type, mime_type, file_size_bytes, sort_order
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING sort_order, created_at
	`

	deletePortfolioItem = `
	DELETE FROM influencer_portfolio_items
	WHERE id = $1 AND user_id = $2
	`
)
