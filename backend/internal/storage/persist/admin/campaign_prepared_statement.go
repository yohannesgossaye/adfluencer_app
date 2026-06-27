package admin

const (
	getAdminCampaignStats = `
	SELECT
		(SELECT COUNT(*)::int FROM campaigns WHERE status = 'active'),
		(SELECT COUNT(*)::int FROM campaigns WHERE status = 'pending'),
		(SELECT COUNT(*)::int FROM campaigns WHERE status = 'paused')
	`

	countAdminCampaigns = `
	SELECT COUNT(*)::int
	FROM campaigns c
	INNER JOIN business_profiles bp ON bp.user_id = c.business_user_id
	LEFT JOIN categories cat ON cat.id = c.category_id
	WHERE ($1 = '' OR c.status::text = $1)
		AND ($2 = '' OR cat.name ILIKE '%' || $2 || '%' OR cat.slug ILIKE '%' || $2 || '%')
		AND ($3 = '' OR bp.business_name ILIKE '%' || $3 || '%')
	`

	listAdminCampaigns = `
	SELECT
		c.id,
		c.title,
		bp.business_name,
		c.business_user_id,
		COALESCE(cat.name, ''),
		c.budget_min_etb,
		c.budget_max_etb,
		c.status::text,
		c.number_of_creators_needed,
		c.created_at,
		COALESCE(inv.total, 0),
		COALESCE(inv.accepted, 0),
		COALESCE(ct.active_count, 0),
		COALESCE(ct.pending_count, 0),
		COALESCE(pay.paid_count, 0),
		COALESCE(pay.pending_count, 0)
	FROM campaigns c
	INNER JOIN business_profiles bp ON bp.user_id = c.business_user_id
	LEFT JOIN categories cat ON cat.id = c.category_id
	LEFT JOIN LATERAL (
		SELECT COUNT(*)::int AS total
		FROM campaign_invites ci
		WHERE ci.campaign_id = c.id
	) inv ON true
	LEFT JOIN LATERAL (
		SELECT
			COUNT(*) FILTER (WHERE ct.status = 'active')::int AS active_count,
			COUNT(*) FILTER (WHERE ct.status IN ('pending_influencer', 'accepted', 'draft'))::int AS pending_count
		FROM contracts ct
		WHERE ct.campaign_id = c.id
	) ct ON true
	LEFT JOIN LATERAL (
		SELECT
			COUNT(*) FILTER (WHERE p.status = 'paid')::int AS paid_count,
			COUNT(*) FILTER (WHERE p.status = 'pending')::int AS pending_count
		FROM payments p
		INNER JOIN contracts ct ON ct.id = p.contract_id
		WHERE ct.campaign_id = c.id
	) pay ON true
	WHERE ($1 = '' OR c.status::text = $1)
		AND ($2 = '' OR cat.name ILIKE '%' || $2 || '%' OR cat.slug ILIKE '%' || $2 || '%')
		AND ($3 = '' OR bp.business_name ILIKE '%' || $3 || '%')
	ORDER BY c.created_at DESC
	LIMIT $4 OFFSET $5
	`

	getAdminCampaign = `
	SELECT
		c.id,
		c.title,
		bp.business_name,
		c.business_user_id,
		COALESCE(cat.name, ''),
		c.budget_min_etb,
		c.budget_max_etb,
		c.status::text,
		c.number_of_creators_needed,
		c.created_at,
		c.updated_at,
		c.goal,
		c.description,
		c.start_date,
		c.end_date,
		COALESCE(inv.total, 0),
		COALESCE(inv.accepted, 0),
		COALESCE(ct.active_count, 0),
		COALESCE(ct.pending_count, 0),
		COALESCE(pay.paid_count, 0),
		COALESCE(pay.pending_count, 0)
	FROM campaigns c
	INNER JOIN business_profiles bp ON bp.user_id = c.business_user_id
	LEFT JOIN categories cat ON cat.id = c.category_id
	LEFT JOIN LATERAL (
		SELECT COUNT(*)::int AS total
		FROM campaign_invites ci
		WHERE ci.campaign_id = c.id
	) inv ON true
	LEFT JOIN LATERAL (
		SELECT
			COUNT(*) FILTER (WHERE ct.status = 'active')::int AS active_count,
			COUNT(*) FILTER (WHERE ct.status IN ('pending_influencer', 'accepted', 'draft'))::int AS pending_count
		FROM contracts ct
		WHERE ct.campaign_id = c.id
	) ct ON true
	LEFT JOIN LATERAL (
		SELECT
			COUNT(*) FILTER (WHERE p.status = 'paid')::int AS paid_count,
			COUNT(*) FILTER (WHERE p.status = 'pending')::int AS pending_count
		FROM payments p
		INNER JOIN contracts ct ON ct.id = p.contract_id
		WHERE ct.campaign_id = c.id
	) pay ON true
	WHERE c.id = $1
	`

	updateAdminCampaignStatus = `
	UPDATE campaigns
	SET status = $2::campaign_status, updated_at = NOW()
	WHERE id = $1
	`
)
