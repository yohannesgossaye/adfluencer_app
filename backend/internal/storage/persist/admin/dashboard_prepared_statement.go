package admin

const (
	// Each table is read once; COUNT(*) FILTER computes all its numbers in that pass.
	// "This month" is the calendar month in Ethiopia (UTC+3), not UTC.
	getAdminDashboardKPIs = `
	WITH
	u AS (
		SELECT
			COUNT(*) FILTER (WHERE role = 'influencer')::int AS influencers,
			COUNT(*) FILTER (WHERE role = 'business')::int AS businesses,
			COUNT(*) FILTER (WHERE role = 'influencer' AND status = 'pending_approval')::int AS pending,
			COUNT(*)::int AS total,
			COUNT(*) FILTER (WHERE status = 'blocked')::int AS blocked
		FROM users
	),
	c AS (
		SELECT
			COUNT(*) FILTER (WHERE status = 'active')::int AS active,
			COUNT(*)::int AS total
		FROM campaigns
	),
	k AS (
		SELECT
			COUNT(*) FILTER (WHERE status = 'completed')::int AS completed,
			COALESCE(SUM(platform_total_fee_etb) FILTER (
				WHERE status = 'completed'
				  AND completed_at >= date_trunc('month', NOW() AT TIME ZONE '` + businessTimeZone + `') AT TIME ZONE '` + businessTimeZone + `'
			), 0)::float8 AS month_revenue,
			COALESCE(SUM(amount_etb) FILTER (WHERE status IN ('active', 'completed')), 0)::float8 AS total_value,
			COALESCE(SUM(platform_total_fee_etb) FILTER (WHERE status = 'completed'), 0)::float8 AS total_fees
		FROM contracts
	)
	SELECT
		u.influencers, u.businesses, u.pending,
		c.active, c.total,
		u.total,
		(SELECT COUNT(*)::int FROM disputes WHERE status IN ('open', 'under_review')),
		(SELECT COUNT(*)::int FROM support_tickets WHERE status IN ('open', 'in_review', 'escalated')),
		u.blocked,
		k.completed, k.month_revenue, k.total_value, k.total_fees
	FROM u, c, k
	`

	listDashboardPendingApprovals = `
	SELECT
		u.id,
		ip.creator_name,
		ip.total_followers,
		u.updated_at,
		COALESCE((
			SELECT c.name
			FROM influencer_profile_categories ipc
			JOIN categories c ON c.id = ipc.category_id
			WHERE ipc.influencer_user_id = u.id
			ORDER BY c.name
			LIMIT 1
		), ''),
		COALESCE((
			SELECT ipp.platform::text
			FROM influencer_profile_platforms ipp
			WHERE ipp.influencer_user_id = u.id
			ORDER BY ipp.platform
			LIMIT 1
		), '')
	FROM users u
	INNER JOIN influencer_profiles ip ON ip.user_id = u.id
	WHERE u.role = 'influencer'
		AND u.status = 'pending_approval'
	ORDER BY u.updated_at DESC
	LIMIT $1
	`

	listDashboardRecentActivity = `
	SELECT
		id,
		action::text,
		COALESCE(description, ''),
		created_at
	FROM audit_logs
	ORDER BY created_at DESC
	LIMIT $1
	`

	// One row per Ethiopian calendar month, oldest first, including months
	// with no completed contracts (so the trend chart has no gaps).
	listDashboardMonthlyRevenue = `
	WITH months AS (
		SELECT generate_series(
			date_trunc('month', NOW() AT TIME ZONE '` + businessTimeZone + `') - ($1 - 1) * INTERVAL '1 month',
			date_trunc('month', NOW() AT TIME ZONE '` + businessTimeZone + `'),
			INTERVAL '1 month'
		) AS month
	)
	SELECT
		m.month AT TIME ZONE '` + businessTimeZone + `',
		COALESCE(SUM(k.platform_total_fee_etb), 0)::float8,
		COUNT(k.id)::int
	FROM months m
	LEFT JOIN contracts k
		ON k.status = 'completed'
		AND k.completed_at >= m.month AT TIME ZONE '` + businessTimeZone + `'
		AND k.completed_at < (m.month + INTERVAL '1 month') AT TIME ZONE '` + businessTimeZone + `'
	GROUP BY m.month
	ORDER BY m.month ASC
	`

	// Disputes and support tickets ranked together, so an urgent ticket is
	// never pushed off the list by medium-severity disputes.
	listDashboardAlerts = `
	SELECT id, entity_type, level, title, detail FROM (
		SELECT
			id::text, 'dispute' AS entity_type, severity::text AS level,
			'Open dispute' AS title, LEFT(description, 160) AS detail,
			CASE severity WHEN 'high' THEN 1 ELSE 3 END AS sort_rank,
			created_at
		FROM disputes
		WHERE status IN ('open', 'under_review') AND severity IN ('high', 'medium')
		UNION ALL
		SELECT
			id::text, 'support_ticket', priority::text,
			subject, LEFT(message, 160),
			CASE priority WHEN 'urgent' THEN 1 ELSE 2 END,
			created_at
		FROM support_tickets
		WHERE status IN ('open', 'in_review', 'escalated') AND priority IN ('urgent', 'high')
	) a
	ORDER BY sort_rank, created_at DESC
	LIMIT $1
	`
)

// businessTimeZone is where Adfluence operates; month boundaries in reports use it.
const businessTimeZone = "Africa/Addis_Ababa"
