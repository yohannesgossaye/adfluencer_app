package business

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
		is_verified
	FROM business_profiles
	WHERE user_id = $1
	`

	updateBusinessProfile = `
	UPDATE business_profiles
	SET
		owner_full_name = $2,
		business_name = $3,
		legal_business_name = $4,
		description = $5,
		industry = $6,
		business_type = $7,
		business_size = $8,
		city = $9,
		region = $10,
		address = $11,
		tin_number = $12,
		website_url = $13,
		instagram_url = $14,
		tiktok_url = $15,
		facebook_url = $16,
		linkedin_url = $17,
		preferred_platforms = $18::text[],
		monthly_marketing_budget = $19,
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
		is_verified
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

	completeBusinessProfile = `
	UPDATE users
	SET status = 'active', status_reason = NULL, updated_at = NOW()
	WHERE id = $1
		AND role = 'business'
		AND status IN ('pending_profile', 'rejected')
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
