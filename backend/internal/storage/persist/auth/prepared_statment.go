package auth

const (
	createUser = `
	INSERT INTO users (
		id,
		email,
		phone_number,
		password_hash,
		role,
		status,
		email_verified,
		phone_verified,
		created_at,
		updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	insertInfluencer = `
	INSERT INTO influencer_profiles (
		id,
		user_id,
		full_name,
		creator_name,
		is_verified,
		created_at,
		updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	insertBusiness = `
	INSERT INTO business_profiles (
		id,
		user_id,
		owner_full_name,
		business_name,
		is_verified,
		created_at,
		updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	insertSession = `
	INSERT INTO user_sessions (
		id,
		user_id,
		refresh_token_hash,
		ip_address,
		user_agent,
		expires_at,
		created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	getSessionByRefreshHash = `
	SELECT
		id,
		user_id,
		expires_at
	FROM user_sessions
	WHERE refresh_token_hash = $1
		AND revoked_at IS NULL
		AND expires_at > NOW()
	LIMIT 1
	`

	revokeSession = `
	UPDATE user_sessions
	SET revoked_at = $2
	WHERE id = $1
		AND revoked_at IS NULL
	`

	getUserByID = `
	SELECT
		id,
		email,
		COALESCE(phone_number, ''),
		role,
		status,
		email_verified,
		phone_verified,
		last_login_at,
		suspended_until,
		status_reason,
		created_at
	FROM users
	WHERE id = $1
	LIMIT 1
	`

	getBusinessProfileByUserID = `
	SELECT id, business_name, owner_full_name, is_verified
	FROM business_profiles
	WHERE user_id = $1
	LIMIT 1
	`

	getInfluencerProfileByUserID = `
	SELECT id, full_name, creator_name, is_verified
	FROM influencer_profiles
	WHERE user_id = $1
	LIMIT 1
	`

	deleteUser = `
	DELETE FROM users
	WHERE id = $1
	`

	getUserByEmail = `
	SELECT
		id,
		email,
		COALESCE(phone_number, ''),
		password_hash,
		role,
		status,
		email_verified,
		phone_verified,
		last_login_at,
		suspended_until,
		status_reason,
		created_at
	FROM users
	WHERE LOWER(email) = LOWER($1)
	LIMIT 1
	`

	updateLastLoginAt = `
	UPDATE users
	SET last_login_at = $2, updated_at = $2
	WHERE id = $1
	`

	getAdminProfileByUserID = `
	SELECT id, full_name, permission_level
	FROM admin_profiles
	WHERE user_id = $1
	LIMIT 1
	`

	insertPasswordResetToken = `
	INSERT INTO password_reset_tokens (
		id,
		user_id,
		token_hash,
		expires_at
	) VALUES ($1, $2, $3, $4)
	`

	getPasswordResetTokenByHash = `
	SELECT id, user_id, expires_at
	FROM password_reset_tokens
	WHERE token_hash = $1
		AND used_at IS NULL
		AND expires_at > NOW()
	LIMIT 1
	`

	completePasswordReset = `
	UPDATE users
	SET password_hash = $2, updated_at = $3
	WHERE id = $1
	`

	markPasswordResetTokenUsed = `
	UPDATE password_reset_tokens
	SET used_at = $2
	WHERE id = $1
		AND used_at IS NULL
	`

	revokeAllUserSessions = `
	UPDATE user_sessions
	SET revoked_at = $2
	WHERE user_id = $1
		AND revoked_at IS NULL
	`

	invalidateEmailVerificationTokens = `
	UPDATE email_verification_tokens
	SET used_at = $2
	WHERE user_id = $1
		AND used_at IS NULL
	`

	insertEmailVerificationToken = `
	INSERT INTO email_verification_tokens (
		id,
		user_id,
		token_hash,
		expires_at
	) VALUES ($1, $2, $3, $4)
	`

	getEmailVerificationTokenByHash = `
	SELECT id, user_id, expires_at
	FROM email_verification_tokens
	WHERE user_id = $1
		AND token_hash = $2
		AND used_at IS NULL
		AND expires_at > NOW()
	LIMIT 1
	`

	countRecentEmailVerificationTokens = `
	SELECT COUNT(*)
	FROM email_verification_tokens
	WHERE user_id = $1
		AND created_at > $2
	`

	markEmailVerified = `
	UPDATE users
	SET email_verified = TRUE, updated_at = $2
	WHERE id = $1
	`

	markEmailVerificationTokenUsed = `
	UPDATE email_verification_tokens
	SET used_at = $2
	WHERE id = $1
		AND used_at IS NULL
	`

	invalidatePhoneVerificationTokens = `
	UPDATE phone_verification_tokens
	SET used_at = $2
	WHERE user_id = $1
		AND used_at IS NULL
	`

	insertPhoneVerificationToken = `
	INSERT INTO phone_verification_tokens (
		id,
		user_id,
		token_hash,
		expires_at
	) VALUES ($1, $2, $3, $4)
	`

	getPhoneVerificationTokenByHash = `
	SELECT id, user_id, expires_at
	FROM phone_verification_tokens
	WHERE user_id = $1
		AND token_hash = $2
		AND used_at IS NULL
		AND expires_at > NOW()
	LIMIT 1
	`

	countRecentPhoneVerificationTokens = `
	SELECT COUNT(*)
	FROM phone_verification_tokens
	WHERE user_id = $1
		AND created_at > $2
	`

	markPhoneVerified = `
	UPDATE users
	SET phone_verified = TRUE, updated_at = $2
	WHERE id = $1
	`

	markPhoneVerificationTokenUsed = `
	UPDATE phone_verification_tokens
	SET used_at = $2
	WHERE id = $1
		AND used_at IS NULL
	`

	insertAuditLog = `
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
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	getLoginLockout = `
	SELECT email, failed_attempts, lockout_tier, locked_until
	FROM auth_login_lockouts
	WHERE email = $1
	LIMIT 1
	`

	upsertLoginLockout = `
	INSERT INTO auth_login_lockouts (
		email,
		failed_attempts,
		lockout_tier,
		locked_until,
		updated_at
	) VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (email) DO UPDATE SET
		failed_attempts = EXCLUDED.failed_attempts,
		lockout_tier = EXCLUDED.lockout_tier,
		locked_until = EXCLUDED.locked_until,
		updated_at = EXCLUDED.updated_at
	`

	deleteLoginLockout = `
	DELETE FROM auth_login_lockouts
	WHERE email = $1
	`

	insertRateLimitEvent = `
	INSERT INTO auth_rate_limit_events (id, scope, rate_key)
	VALUES ($1, $2, $3)
	`

	countRateLimitEvents = `
	SELECT COUNT(*)
	FROM auth_rate_limit_events
	WHERE scope = $1
		AND rate_key = $2
		AND created_at > $3
	`
)
