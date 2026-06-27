package admin

const (
	getAdminDashboardKPIs = `
	SELECT
		(SELECT COUNT(*)::int FROM users WHERE role = 'influencer'),
		(SELECT COUNT(*)::int FROM users WHERE role = 'business'),
		(SELECT COUNT(*)::int FROM users WHERE role = 'influencer' AND status = 'pending_approval'),
		(SELECT COUNT(*)::int FROM campaigns WHERE status = 'active'),
		(SELECT COUNT(*)::int FROM campaigns),
		(SELECT COUNT(*)::int FROM users),
		(SELECT COUNT(*)::int FROM disputes WHERE status IN ('open', 'under_review')),
		(SELECT COUNT(*)::int FROM support_tickets WHERE status IN ('open', 'in_review', 'escalated')),
		(SELECT COUNT(*)::int FROM users WHERE status = 'blocked'),
		(SELECT COUNT(*)::int FROM contracts WHERE status = 'completed'),
		(SELECT COALESCE(SUM(platform_total_fee_etb), 0)::float8
		 FROM contracts
		 WHERE status = 'completed'
		   AND completed_at >= date_trunc('month', NOW())),
		(SELECT COALESCE(SUM(amount_etb), 0)::float8
		 FROM contracts
		 WHERE status IN ('active', 'completed')),
		(SELECT COALESCE(SUM(platform_total_fee_etb), 0)::float8
		 FROM contracts
		 WHERE status = 'completed')
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

	listDashboardMonthlyRevenue = `
	SELECT
		date_trunc('month', completed_at)::timestamptz,
		COALESCE(SUM(platform_total_fee_etb), 0)::float8,
		COUNT(*)::int
	FROM contracts
	WHERE status = 'completed'
		AND completed_at >= date_trunc('month', NOW()) - ($1 - 1) * INTERVAL '1 month'
	GROUP BY 1
	ORDER BY 1 ASC
	`

	listDashboardDisputeAlerts = `
	SELECT
		id,
		severity::text,
		LEFT(description, 160)
	FROM disputes
	WHERE status IN ('open', 'under_review')
		AND severity IN ('high', 'medium')
	ORDER BY
		CASE severity WHEN 'high' THEN 1 WHEN 'medium' THEN 2 ELSE 3 END,
		created_at DESC
	LIMIT $1
	`

	listDashboardSupportAlerts = `
	SELECT
		id,
		priority::text,
		subject,
		LEFT(message, 160)
	FROM support_tickets
	WHERE status IN ('open', 'in_review', 'escalated')
		AND priority IN ('urgent', 'high')
	ORDER BY
		CASE priority WHEN 'urgent' THEN 1 WHEN 'high' THEN 2 ELSE 3 END,
		created_at DESC
	LIMIT $1
	`
)
