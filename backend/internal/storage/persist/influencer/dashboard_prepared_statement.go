package influencer

const (
	// One pass per table: COUNT/SUM ... FILTER computes every number together.
	getInfluencerDashboardKPIs = `
	WITH
	k AS (
		SELECT
			COUNT(*) FILTER (WHERE status = 'active')::int AS active,
			COUNT(*) FILTER (WHERE status IN ('pending_influencer', 'accepted'))::int AS pending,
			COUNT(*) FILTER (WHERE status = 'completed')::int AS completed,
			COALESCE(SUM(influencer_receivable_etb) FILTER (WHERE status = 'completed'), 0)::float8 AS earned
		FROM contracts WHERE influencer_user_id = $1
	),
	p AS (
		SELECT
			COALESCE(SUM(amount_etb) FILTER (WHERE status IN ('pending', 'processing')), 0)::float8 AS pending_amount,
			COALESCE(SUM(amount_etb) FILTER (WHERE status = 'paid'), 0)::float8 AS paid_amount,
			MIN(created_at) FILTER (WHERE status IN ('pending', 'processing')) AS oldest_pending
		FROM payouts WHERE influencer_user_id = $1
	)
	SELECT
		(SELECT COUNT(*)::int FROM campaign_invites WHERE influencer_user_id = $1 AND status = 'pending'),
		k.active, k.pending, k.completed, k.earned,
		p.pending_amount, p.paid_amount, p.oldest_pending
	FROM k, p
	`

	listInfluencerDashboardPendingInvites = `
	SELECT
		ci.id,
		bp.business_name,
		c.title,
		ci.status::text,
		c.budget_min_etb,
		c.budget_max_etb,
		c.end_date
	FROM campaign_invites ci
	INNER JOIN campaigns c ON c.id = ci.campaign_id
	INNER JOIN business_profiles bp ON bp.user_id = ci.business_user_id
	WHERE ci.influencer_user_id = $1 AND ci.status = 'pending'
	ORDER BY ci.created_at DESC
	LIMIT 3
	`

	listInfluencerDashboardActiveContracts = `
	SELECT
		ct.id,
		bp.business_name,
		ct.status::text
	FROM contracts ct
	INNER JOIN business_profiles bp ON bp.user_id = ct.business_user_id
	WHERE ct.influencer_user_id = $1 AND ct.status = 'active'
	ORDER BY ct.updated_at DESC
	LIMIT 5
	`
)
