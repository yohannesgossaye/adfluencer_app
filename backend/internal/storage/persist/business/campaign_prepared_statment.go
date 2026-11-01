package business

const (
	resolveCategoryByNameOrSlug = `
	SELECT id, name, slug
	FROM categories
	WHERE is_active = TRUE
		AND (LOWER(name) = LOWER($1) OR LOWER(slug) = LOWER($1))
	LIMIT 1
	`

	insertCampaign = `
	INSERT INTO campaigns (
		id, business_user_id, title, goal, description, category_id,
		target_audience, city, region, budget_min_etb, budget_max_etb,
		creator_size_preference, number_of_creators_needed,
		start_date, end_date, status
	) VALUES (
		$1, $2, $3, $4, $5, $6,
		$7, $8, $9, $10, $11,
		$12::creator_size, $13,
		$14, $15, 'draft'
	)
	`

	getCampaignByID = `
	SELECT
		c.id,
		c.business_user_id,
		c.title,
		c.goal,
		c.description,
		c.category_id,
		cat.name,
		cat.slug,
		c.target_audience,
		c.city,
		c.region,
		c.budget_min_etb,
		c.budget_max_etb,
		c.creator_size_preference::text,
		c.number_of_creators_needed,
		c.start_date,
		c.end_date,
		c.status::text,
		c.created_at,
		c.updated_at,
		bp.business_name,
		COALESCE(inv.total, 0),
		COALESCE(inv.accepted, 0)
	FROM campaigns c
	JOIN business_profiles bp ON bp.user_id = c.business_user_id
	LEFT JOIN categories cat ON cat.id = c.category_id
	LEFT JOIN LATERAL (
		SELECT
			COUNT(*)::int AS total,
			COUNT(*) FILTER (WHERE ci.status = 'accepted')::int AS accepted
		FROM campaign_invites ci
		WHERE ci.campaign_id = c.id
	) inv ON true
	WHERE c.id = $1 AND c.business_user_id = $2
	`

	listBusinessCampaigns = `
	SELECT
		c.id,
		c.title,
		bp.business_name,
		c.budget_min_etb,
		c.budget_max_etb,
		c.status::text,
		c.start_date,
		c.end_date,
		c.number_of_creators_needed,
		c.created_at,
		COALESCE(inv.total, 0),
		COALESCE(inv.accepted, 0),
		(
			SELECT cp.platform::text
			FROM campaign_platforms cp
			WHERE cp.campaign_id = c.id
			ORDER BY cp.platform
			LIMIT 1
		),
		COUNT(*) OVER()::int AS total_count
	FROM campaigns c
	JOIN business_profiles bp ON bp.user_id = c.business_user_id
	LEFT JOIN LATERAL (
		SELECT
			COUNT(*)::int AS total,
			COUNT(*) FILTER (WHERE ci.status = 'accepted')::int AS accepted
		FROM campaign_invites ci
		WHERE ci.campaign_id = c.id
	) inv ON true
	WHERE c.business_user_id = $1
		AND ($2 = '' OR c.status::text = $2)
	ORDER BY c.created_at DESC
	LIMIT $3 OFFSET $4
	`

	updateCampaign = `
	UPDATE campaigns
	SET
		title = $3,
		goal = $4,
		description = $5,
		category_id = $6,
		target_audience = $7,
		city = $8,
		region = $9,
		budget_min_etb = $10,
		budget_max_etb = $11,
		creator_size_preference = $12::creator_size,
		number_of_creators_needed = $13,
		start_date = $14,
		end_date = $15,
		updated_at = NOW()
	WHERE id = $1 AND business_user_id = $2
		AND status IN ('draft', 'active', 'paused')
	`

	publishCampaign = `
	UPDATE campaigns
	SET status = 'active', updated_at = NOW()
	WHERE id = $1 AND business_user_id = $2 AND status = 'draft'
	`

	cancelCampaign = `
	UPDATE campaigns
	SET status = 'cancelled', updated_at = NOW()
	WHERE id = $1 AND business_user_id = $2
		AND status IN ('draft', 'active', 'paused')
	`

	cancelPendingCampaignInvites = `
	UPDATE campaign_invites
	SET status = 'cancelled', updated_at = NOW()
	WHERE campaign_id = $1 AND status = 'pending'
	`

	deleteCampaignPlatforms = `
	DELETE FROM campaign_platforms WHERE campaign_id = $1
	`

	insertCampaignPlatform = `
	INSERT INTO campaign_platforms (id, campaign_id, platform)
	VALUES ($1, $2, $3::social_platform)
	ON CONFLICT DO NOTHING
	`

	deleteCampaignDeliverables = `
	DELETE FROM campaign_deliverables WHERE campaign_id = $1
	`

	insertCampaignDeliverable = `
	INSERT INTO campaign_deliverables (id, campaign_id, deliverable_type, quantity, description)
	VALUES ($1, $2, $3::deliverable_type, $4, $5)
	`

	getCampaignPlatforms = `
	SELECT platform::text
	FROM campaign_platforms
	WHERE campaign_id = $1
	ORDER BY platform
	`

	getCampaignDeliverables = `
	SELECT id, deliverable_type::text, quantity, description
	FROM campaign_deliverables
	WHERE campaign_id = $1
	ORDER BY created_at
	`
)
