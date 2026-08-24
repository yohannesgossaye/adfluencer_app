package admin

const (
	listInfluencerApplications = `
	SELECT
		u.id,
		u.status::text,
		u.status_reason,
		u.email,
		COALESCE(u.phone_number, ''),
		u.updated_at,
		ip.creator_name,
		ip.profile_image_url,
		ip.bio,
		ip.city,
		ip.total_followers,
		ip.engagement_rate,
		ip.min_price_etb,
		ip.max_price_etb,
		ip.audience_description,
		ip.approval_note,
		COALESCE((
			SELECT c.name
			FROM influencer_profile_categories ipc
			JOIN categories c ON c.id = ipc.category_id
			WHERE ipc.influencer_user_id = u.id
			ORDER BY c.name
			LIMIT 1
		), '') AS primary_niche,
		COALESCE((
			SELECT ipp.platform::text
			FROM influencer_profile_platforms ipp
			WHERE ipp.influencer_user_id = u.id
			ORDER BY ipp.platform
			LIMIT 1
		), '') AS primary_platform,
		COUNT(*) OVER()::int AS total_count
	FROM users u
	JOIN influencer_profiles ip ON ip.user_id = u.id
	WHERE u.role = 'influencer'
		AND ($1 = '' OR u.status::text = $1)
		AND (
			$2 = ''
			OR ip.creator_name ILIKE '%' || $2 || '%'
			OR ip.bio ILIKE '%' || $2 || '%'
			OR EXISTS (
				SELECT 1
				FROM influencer_profile_categories ipc
				JOIN categories c ON c.id = ipc.category_id
				WHERE ipc.influencer_user_id = u.id
					AND c.name ILIKE '%' || $2 || '%'
			)
			OR EXISTS (
				SELECT 1
				FROM influencer_profile_platforms ipp
				WHERE ipp.influencer_user_id = u.id
					AND ipp.platform::text ILIKE '%' || $2 || '%'
			)
		)
	ORDER BY u.updated_at DESC
	LIMIT $3 OFFSET $4
	`

	getInfluencerApplication = `
	SELECT
		u.id,
		u.status::text,
		u.status_reason,
		u.email,
		COALESCE(u.phone_number, ''),
		u.updated_at,
		ip.creator_name,
		ip.profile_image_url,
		ip.bio,
		ip.city,
		ip.total_followers,
		ip.engagement_rate,
		ip.min_price_etb,
		ip.max_price_etb,
		ip.audience_description,
		ip.approval_note,
		COALESCE((
			SELECT c.name
			FROM influencer_profile_categories ipc
			JOIN categories c ON c.id = ipc.category_id
			WHERE ipc.influencer_user_id = u.id
			ORDER BY c.name
			LIMIT 1
		), '') AS primary_niche,
		COALESCE((
			SELECT ipp.platform::text
			FROM influencer_profile_platforms ipp
			WHERE ipp.influencer_user_id = u.id
			ORDER BY ipp.platform
			LIMIT 1
		), '') AS primary_platform
	FROM users u
	JOIN influencer_profiles ip ON ip.user_id = u.id
	WHERE u.role = 'influencer'
		AND u.id = $1
	`

	approveInfluencerApplication = `
	UPDATE users
	SET status = 'active', status_reason = NULL, updated_at = NOW()
	WHERE id = $1
		AND role = 'influencer'
		AND (status = 'pending_approval' OR status = 'pending_profile' OR status = 'active')
	`

	setInfluencerApproved = `
	UPDATE influencer_profiles
	SET
		is_verified = TRUE,
		approved_at = NOW(),
		approved_by = $2,
		approval_note = COALESCE($3, approval_note),
		updated_at = NOW()
	WHERE user_id = $1
	`

	rejectInfluencerApplication = `
	UPDATE users
	SET status = 'rejected', status_reason = $2, updated_at = NOW()
	WHERE id = $1
		AND role = 'influencer'
		AND (status = 'pending_approval' OR status = 'pending_profile' OR status = 'active')
	`

	setInfluencerRejected = `
	UPDATE influencer_profiles
	SET
		is_verified = FALSE,
		approved_at = NULL,
		approved_by = NULL,
		approval_note = COALESCE($2, approval_note),
		updated_at = NOW()
	WHERE user_id = $1
	`

	updateInfluencerAdminNote = `
	UPDATE influencer_profiles
	SET approval_note = $2, updated_at = NOW()
	WHERE user_id = $1
	`

	insertAdminAuditLog = `
	INSERT INTO audit_logs (
		id,
		actor_user_id,
		target_user_id,
		action,
		entity_type,
		entity_id,
		description,
		ip_address,
		user_agent,
		created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
	`
)
