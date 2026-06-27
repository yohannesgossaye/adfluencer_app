package admin

const (
	countAdminUsers = `
	SELECT COUNT(*)::int
	FROM users u
	LEFT JOIN business_profiles bp ON bp.user_id = u.id AND u.role = 'business'
	LEFT JOIN influencer_profiles ip ON ip.user_id = u.id AND u.role = 'influencer'
	WHERE ($1 = '' OR u.role::text = $1)
		AND (
			$2 = ''
			OR ($2 = 'pending' AND u.status IN ('pending_profile', 'pending_approval'))
			OR ($2 <> 'pending' AND u.status::text = $2)
		)
		AND (
			$3 = ''
			OR u.email ILIKE '%' || $3 || '%'
			OR COALESCE(bp.business_name, ip.creator_name, '') ILIKE '%' || $3 || '%'
			OR COALESCE(bp.owner_full_name, '') ILIKE '%' || $3 || '%'
		)
	`

	getAdminUserStats = `
	SELECT
		COUNT(*)::int,
		COUNT(*) FILTER (WHERE status = 'active')::int,
		COUNT(*) FILTER (WHERE status = 'blocked')::int,
		COUNT(*) FILTER (WHERE status IN ('pending_profile', 'pending_approval'))::int
	FROM users
	`

	listAdminUsers = `
	SELECT
		u.id,
		u.email,
		u.role::text,
		u.status::text,
		u.status_reason,
		COALESCE(bp.business_name, ip.creator_name, split_part(u.email, '@', 1)),
		u.created_at,
		u.last_login_at,
		(SELECT COUNT(*)::int FROM campaigns c WHERE c.business_user_id = u.id),
		(SELECT COUNT(*)::int FROM contracts ct
		 WHERE ct.business_user_id = u.id OR ct.influencer_user_id = u.id)
	FROM users u
	LEFT JOIN business_profiles bp ON bp.user_id = u.id AND u.role = 'business'
	LEFT JOIN influencer_profiles ip ON ip.user_id = u.id AND u.role = 'influencer'
	WHERE ($1 = '' OR u.role::text = $1)
		AND (
			$2 = ''
			OR ($2 = 'pending' AND u.status IN ('pending_profile', 'pending_approval'))
			OR ($2 <> 'pending' AND u.status::text = $2)
		)
		AND (
			$3 = ''
			OR u.email ILIKE '%' || $3 || '%'
			OR COALESCE(bp.business_name, ip.creator_name, '') ILIKE '%' || $3 || '%'
			OR COALESCE(bp.owner_full_name, '') ILIKE '%' || $3 || '%'
		)
	ORDER BY u.created_at DESC
	LIMIT $4 OFFSET $5
	`

	getAdminUser = `
	SELECT
		u.id,
		u.email,
		u.role::text,
		u.status::text,
		u.status_reason,
		COALESCE(bp.business_name, ip.creator_name, split_part(u.email, '@', 1)),
		u.created_at,
		u.last_login_at,
		(SELECT COUNT(*)::int FROM campaigns c WHERE c.business_user_id = u.id),
		(SELECT COUNT(*)::int FROM contracts ct
		 WHERE ct.business_user_id = u.id OR ct.influencer_user_id = u.id)
	FROM users u
	LEFT JOIN business_profiles bp ON bp.user_id = u.id AND u.role = 'business'
	LEFT JOIN influencer_profiles ip ON ip.user_id = u.id AND u.role = 'influencer'
	WHERE u.id = $1
	`

	blockUser = `
	UPDATE users
	SET status = 'blocked', status_reason = $2, updated_at = NOW()
	WHERE id = $1 AND role <> 'admin' AND status <> 'removed'
	`

	unblockUser = `
	UPDATE users
	SET status = 'active', status_reason = NULL, updated_at = NOW()
	WHERE id = $1 AND role <> 'admin' AND status = 'blocked'
	`

	removeUser = `
	UPDATE users
	SET status = 'removed', status_reason = $2, updated_at = NOW()
	WHERE id = $1 AND role <> 'admin' AND status <> 'removed'
	`

	getUserRole = `
	SELECT role::text, status::text FROM users WHERE id = $1
	`
)
