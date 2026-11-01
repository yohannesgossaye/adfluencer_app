package admin

const (
	// The queue lists businesses that have submitted at least once. Pending
	// reviews come oldest-first so nobody waits forever.
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
		bp.verification_status,
		bp.verification_submitted_at,
		u.created_at,
		COUNT(*) OVER()::int AS total_count
	FROM users u
	INNER JOIN business_profiles bp ON bp.user_id = u.id
	WHERE u.role = 'business'
		AND u.status = 'active'
		AND bp.verification_status <> 'unsubmitted'
		AND ($1 = 'all' OR bp.verification_status = $1)
		AND (
			$2 = ''
			OR bp.business_name ILIKE '%' || $2 || '%'
			OR bp.owner_full_name ILIKE '%' || $2 || '%'
			OR bp.legal_business_name ILIKE '%' || $2 || '%'
			OR bp.tin_number ILIKE '%' || $2 || '%'
			OR u.email ILIKE '%' || $2 || '%'
		)
	ORDER BY
		CASE WHEN $1 = 'pending' THEN bp.verification_submitted_at END ASC NULLS LAST,
		bp.updated_at DESC
	LIMIT $3 OFFSET $4
	`

	getBusinessVerification = `
	SELECT
		u.id,
		u.email,
		u.email_verified,
		COALESCE(u.phone_number, ''),
		u.phone_verified,
		u.status::text,
		bp.business_name,
		bp.legal_business_name,
		bp.owner_full_name,
		bp.tin_number,
		bp.industry,
		bp.business_type,
		bp.business_size,
		bp.city,
		bp.region,
		bp.address,
		bp.website_url,
		bp.instagram_url,
		bp.tiktok_url,
		bp.facebook_url,
		bp.linkedin_url,
		bp.description,
		bp.logo_url,
		bp.is_verified,
		bp.verification_status,
		bp.verification_note,
		bp.verification_submitted_at,
		bp.verified_at,
		u.created_at
	FROM users u
	INNER JOIN business_profiles bp ON bp.user_id = u.id
	WHERE u.id = $1
		AND u.role = 'business'
	`

	getBusinessVerificationStatus = `
	SELECT bp.verification_status
	FROM business_profiles bp
	JOIN users u ON u.id = bp.user_id
	WHERE bp.user_id = $1 AND u.role = 'business'
	`

	// $2 admin id, $3 optional note for the business.
	verifyBusinessProfile = `
	UPDATE business_profiles bp
	SET verification_status = 'verified',
		is_verified = TRUE,
		verified_at = NOW(),
		verified_by = $2,
		verification_note = $3,
		updated_at = NOW()
	FROM users u
	WHERE bp.user_id = u.id
		AND u.id = $1
		AND u.role = 'business'
		AND u.status = 'active'
		AND bp.verification_status IN ('pending', 'rejected')
	RETURNING bp.verified_at
	`

	// Reject a pending review ($4 = 'pending') or revoke a verified business
	// ($4 = 'verified'). $2 admin id, $3 reason shown to the business.
	rejectBusinessVerification = `
	UPDATE business_profiles bp
	SET verification_status = 'rejected',
		is_verified = FALSE,
		verified_at = NULL,
		verified_by = $2,
		verification_note = $3,
		updated_at = NOW()
	FROM users u
	WHERE bp.user_id = u.id
		AND u.id = $1
		AND u.role = 'business'
		AND bp.verification_status = $4
	`
)
