package admin

const (
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
		 WHERE ct.business_user_id = u.id OR ct.influencer_user_id = u.id),
		COUNT(*) OVER()::int AS total_count
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

	// Blocking remembers the previous status (re-blocking keeps the original).
	blockUser = `
	UPDATE users
	SET previous_status = CASE WHEN status = 'blocked' THEN previous_status ELSE status END,
		status = 'blocked', status_reason = $2, updated_at = NOW()
	WHERE id = $1 AND role <> 'admin'
	`

	// Only blocked accounts can be unblocked; they return to the status they
	// had before (e.g. a creator still awaiting approval stays awaiting it).
	unblockUser = `
	UPDATE users
	SET status = COALESCE(previous_status, 'active'), previous_status = NULL,
		status_reason = NULL, updated_at = NOW()
	WHERE id = $1 AND role <> 'admin' AND status = 'blocked'
	RETURNING status::text
	`

	revokeUserSessions = `
	UPDATE user_sessions
	SET revoked_at = NOW()
	WHERE user_id = $1 AND revoked_at IS NULL
	`

	removeUser = `
	UPDATE users
	SET status = 'removed', status_reason = $2, updated_at = NOW()
	WHERE id = $1 AND role <> 'admin'
	`

	getUserRole = `
	SELECT role::text, status::text FROM users WHERE id = $1
	`
)
