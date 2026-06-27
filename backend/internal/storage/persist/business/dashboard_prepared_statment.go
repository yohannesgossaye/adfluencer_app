package business

const (
	getDashboardKPIs = `
	SELECT
		(SELECT COUNT(*)::int FROM campaigns WHERE business_user_id = $1 AND status = 'active'),
		(SELECT COUNT(*)::int FROM campaigns WHERE business_user_id = $1),
		(SELECT COUNT(*)::int
		 FROM campaign_recommendations cr
		 INNER JOIN campaigns c ON c.id = cr.campaign_id
		 WHERE c.business_user_id = $1 AND cr.match_score >= 78),
		(SELECT COUNT(*)::int
		 FROM campaign_recommendations cr
		 INNER JOIN campaigns c ON c.id = cr.campaign_id
		 WHERE c.business_user_id = $1 AND cr.match_score >= 85),
		(SELECT COUNT(*)::int FROM contracts WHERE business_user_id = $1 AND status IN ('pending_influencer', 'accepted')),
		(SELECT COUNT(*)::int FROM contracts WHERE business_user_id = $1 AND status = 'pending_influencer'),
		(SELECT COUNT(*)::int FROM contracts WHERE business_user_id = $1 AND status = 'accepted'),
		(SELECT COUNT(*)::int
		 FROM campaign_invites ci
		 INNER JOIN campaigns c ON c.id = ci.campaign_id
		 WHERE c.business_user_id = $1 AND ci.status = 'pending'),
		(SELECT COUNT(*)::int FROM contracts WHERE business_user_id = $1 AND status = 'active'),
		(SELECT COUNT(*)::int FROM contracts WHERE business_user_id = $1 AND status = 'completed')
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

	listDashboardTopRecommendations = `
	SELECT
		cr.influencer_user_id,
		ip.creator_name,
		ip.profile_image_url,
		cr.match_score,
		c.id,
		c.title
	FROM campaign_recommendations cr
	INNER JOIN campaigns c ON c.id = cr.campaign_id AND c.business_user_id = $1
	INNER JOIN influencer_profiles ip ON ip.user_id = cr.influencer_user_id
	WHERE cr.match_score >= 78
	ORDER BY cr.match_score DESC, cr.rank_position ASC NULLS LAST
	LIMIT 5
	`
)
