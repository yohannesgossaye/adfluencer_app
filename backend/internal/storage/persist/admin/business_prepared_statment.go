package admin

const (
	listBusinessVerifications = `
	SELECT
		u.id,
		bp.business_name,
		bp.owner_full_name,
		bp.industry,
		bp.city,
		bp.logo_url,
		u.email,
		bp.is_verified,
		u.created_at
	FROM users u
	INNER JOIN business_profiles bp ON bp.user_id = u.id
	WHERE u.role = 'business'
		AND u.status = 'active'
		AND (
			($1 = 'pending' AND bp.is_verified = FALSE)
			OR ($1 = 'verified' AND bp.is_verified = TRUE)
		)
		AND (
			$2 = ''
			OR bp.business_name ILIKE '%' || $2 || '%'
			OR bp.owner_full_name ILIKE '%' || $2 || '%'
			OR u.email ILIKE '%' || $2 || '%'
		)
	ORDER BY u.created_at DESC
	LIMIT $3 OFFSET $4
	`

	countBusinessVerifications = `
	SELECT COUNT(*)::int
	FROM users u
	INNER JOIN business_profiles bp ON bp.user_id = u.id
	WHERE u.role = 'business'
		AND u.status = 'active'
		AND (
			($1 = 'pending' AND bp.is_verified = FALSE)
			OR ($1 = 'verified' AND bp.is_verified = TRUE)
		)
		AND (
			$2 = ''
			OR bp.business_name ILIKE '%' || $2 || '%'
			OR bp.owner_full_name ILIKE '%' || $2 || '%'
			OR u.email ILIKE '%' || $2 || '%'
		)
	`

	getBusinessVerification = `
	SELECT
		u.id,
		u.email,
		COALESCE(u.phone_number, ''),
		u.status::text,
		bp.business_name,
		bp.owner_full_name,
		bp.industry,
		bp.city,
		bp.description,
		bp.logo_url,
		bp.is_verified,
		u.created_at
	FROM users u
	INNER JOIN business_profiles bp ON bp.user_id = u.id
	WHERE u.id = $1
		AND u.role = 'business'
		AND u.status = 'active'
	`

	verifyBusinessProfile = `
	UPDATE business_profiles bp
	SET is_verified = TRUE, updated_at = NOW()
	FROM users u
	WHERE bp.user_id = u.id
		AND u.id = $1
		AND u.role = 'business'
		AND u.status = 'active'
		AND bp.is_verified = FALSE
	`

	revokeBusinessVerification = `
	UPDATE business_profiles bp
	SET is_verified = FALSE, updated_at = NOW()
	FROM users u
	WHERE bp.user_id = u.id
		AND u.id = $1
		AND u.role = 'business'
		AND u.status = 'active'
		AND bp.is_verified = TRUE
	`
)
