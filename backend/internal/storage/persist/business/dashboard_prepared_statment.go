package business

const (
	// One pass per table: COUNT(*) FILTER computes every number together.
	getDashboardKPIs = `
	WITH
	c AS (
		SELECT
			COUNT(*) FILTER (WHERE status = 'active')::int AS active,
			COUNT(*)::int AS total
		FROM campaigns WHERE business_user_id = $1
	),
	-- Distinct creators still recommended on open campaigns (not dismissed or
	-- already invited, and not on cancelled/completed campaigns).
	r AS (
		SELECT
			COUNT(DISTINCT cr.influencer_user_id) FILTER (WHERE cr.match_score >= 78)::int AS good,
			COUNT(DISTINCT cr.influencer_user_id) FILTER (WHERE cr.match_score >= 85)::int AS strong
		FROM campaign_recommendations cr
		INNER JOIN campaigns c ON c.id = cr.campaign_id
		WHERE c.business_user_id = $1 AND c.status IN ('draft', 'active')
		  AND cr.status = 'recommended' AND cr.match_score >= 78
	),
	k AS (
		SELECT
			COUNT(*) FILTER (WHERE status IN ('pending_influencer', 'accepted'))::int AS pending,
			COUNT(*) FILTER (WHERE status = 'pending_influencer')::int AS awaiting_creator,
			COUNT(*) FILTER (WHERE status = 'accepted')::int AS awaiting_funding,
			COUNT(*) FILTER (WHERE status = 'active')::int AS active,
			COUNT(*) FILTER (WHERE status = 'completed')::int AS completed
		FROM contracts WHERE business_user_id = $1
	)
	SELECT
		c.active, c.total,
		r.good, r.strong,
		k.pending, k.awaiting_creator, k.awaiting_funding,
		(SELECT COUNT(*)::int
		 FROM campaign_invites ci
		 INNER JOIN campaigns cp ON cp.id = ci.campaign_id
		 WHERE cp.business_user_id = $1 AND ci.status = 'pending'),
		k.active, k.completed
	FROM c, r, k
	`

	listDashboardRecentCampaigns = `
	SELECT
		c.id,
		c.title,
		c.status::text,
		c.budget_min_etb,
		c.budget_max_etb,
		COALESCE(cat.name, ''),
		c.number_of_creators_needed,
		COALESCE(inv.accepted, 0),
		c.updated_at
	FROM campaigns c
	LEFT JOIN categories cat ON cat.id = c.category_id
	LEFT JOIN LATERAL (
		SELECT COUNT(*)::int AS accepted
		FROM campaign_invites ci
		WHERE ci.campaign_id = c.id AND ci.status = 'accepted'
	) inv ON true
	WHERE c.business_user_id = $1
	ORDER BY c.updated_at DESC
	LIMIT 5
	`

	// Best match per creator (each creator appears once), open campaigns only.
	listDashboardTopRecommendations = `
	SELECT influencer_user_id, creator_name, profile_image_url, match_score, campaign_id, campaign_title
	FROM (
		SELECT DISTINCT ON (cr.influencer_user_id)
			cr.influencer_user_id,
			ip.creator_name,
			ip.profile_image_url,
			cr.match_score,
			c.id AS campaign_id,
			c.title AS campaign_title,
			cr.rank_position
		FROM campaign_recommendations cr
		INNER JOIN campaigns c ON c.id = cr.campaign_id
			AND c.business_user_id = $1
			AND c.status IN ('draft', 'active')
		INNER JOIN influencer_profiles ip ON ip.user_id = cr.influencer_user_id
		WHERE cr.match_score >= 78 AND cr.status = 'recommended'
		ORDER BY cr.influencer_user_id, cr.match_score DESC, cr.rank_position ASC NULLS LAST
	) best
	ORDER BY match_score DESC, rank_position ASC NULLS LAST
	LIMIT 5
	`
)
