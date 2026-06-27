package influencer

const (
	getInfluencerDashboardKPIs = `
	SELECT
		(SELECT COUNT(*)::int FROM campaign_invites WHERE influencer_user_id = $1 AND status = 'pending'),
		(SELECT COUNT(*)::int FROM contracts WHERE influencer_user_id = $1 AND status = 'active'),
		(SELECT COUNT(*)::int FROM contracts WHERE influencer_user_id = $1 AND status IN ('pending_influencer', 'accepted')),
		(SELECT COUNT(*)::int FROM contracts WHERE influencer_user_id = $1 AND status = 'completed'),
		(SELECT COALESCE(SUM(influencer_receivable_etb), 0)::float8 FROM contracts WHERE influencer_user_id = $1 AND status = 'completed'),
		(SELECT COALESCE(SUM(amount_etb), 0)::float8 FROM payouts WHERE influencer_user_id = $1 AND status IN ('pending', 'processing')),
		(SELECT COALESCE(SUM(amount_etb), 0)::float8 FROM payouts WHERE influencer_user_id = $1 AND status = 'paid'),
		(SELECT MIN(created_at) FROM payouts WHERE influencer_user_id = $1 AND status IN ('pending', 'processing'))
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
