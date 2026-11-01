package business

const identityChanged = `(verification_status = 'verified' AND (
		business_name::text IS DISTINCT FROM $3::text
		OR legal_business_name::text IS DISTINCT FROM $4::text
		OR tin_number::text IS DISTINCT FROM $12::text))`

const (
	getBusinessProfileFull = `
	SELECT
		id,
		user_id,
		owner_full_name,
		business_name,
		legal_business_name,
		logo_url,
		description,
		industry,
		business_type,
		business_size,
		city,
		region,
		address,
		tin_number,
		website_url,
		instagram_url,
		tiktok_url,
		facebook_url,
		linkedin_url,
		preferred_platforms,
		monthly_marketing_budget,
		is_verified,
		verification_status,
		verification_note,
		verification_submitted_at,
		verified_at
	FROM business_profiles
	WHERE user_id = $1
	`

	updateBusinessProfile = `
	UPDATE business_profiles
	SET
		owner_full_name = $2,
		business_name = $3::text,
		legal_business_name = $4::text,
		description = $5,
		industry = $6,
		business_type = $7,
		business_size = $8,
		city = $9,
		region = $10,
		address = $11,
		tin_number = $12::text,
		website_url = $13,
		instagram_url = $14,
		tiktok_url = $15,
		facebook_url = $16,
		linkedin_url = $17,
		preferred_platforms = $18::text[],
		monthly_marketing_budget = $19,
		-- Changing the legal identity of a verified business sends it back to review.
		verification_status = CASE WHEN ` + identityChanged + ` THEN 'pending' ELSE verification_status END,
		is_verified = CASE WHEN ` + identityChanged + ` THEN FALSE ELSE is_verified END,
		verification_submitted_at = CASE WHEN ` + identityChanged + ` THEN NOW() ELSE verification_submitted_at END,
		verification_note = CASE WHEN ` + identityChanged + `
			THEN 'Business name, legal name or TIN changed, so the profile is being re-verified.'
			ELSE verification_note END,
		updated_at = NOW()
	WHERE user_id = $1
	RETURNING
		id,
		user_id,
		owner_full_name,
		business_name,
		legal_business_name,
		logo_url,
		description,
		industry,
		business_type,
		business_size,
		city,
		region,
		address,
		tin_number,
		website_url,
		instagram_url,
		tiktok_url,
		facebook_url,
		linkedin_url,
		preferred_platforms,
		monthly_marketing_budget,
		is_verified,
		verification_status,
		verification_note,
		verification_submitted_at,
		verified_at
	`

	isBusinessVerified = `
	SELECT is_verified
	FROM business_profiles
	WHERE user_id = $1
	`

	updateBusinessLogoURL = `
	UPDATE business_profiles
	SET logo_url = $2, updated_at = NOW()
	WHERE user_id = $1
	`

	activateBusinessUser = `
	UPDATE users
	SET status = 'active', status_reason = NULL, updated_at = NOW()
	WHERE id = $1
		AND role = 'business'
		AND status IN ('pending_profile', 'rejected', 'active')
	`

	// Submitting (or re-submitting after a rejection) puts the business in
	// the admin review queue.
	submitBusinessForVerification = `
	UPDATE business_profiles
	SET verification_status = 'pending',
		verification_submitted_at = NOW(),
		verification_note = NULL,
		is_verified = FALSE,
		updated_at = NOW()
	WHERE user_id = $1
		AND verification_status IN ('unsubmitted', 'rejected')
	`

	insertBusinessAuditLog = `
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
