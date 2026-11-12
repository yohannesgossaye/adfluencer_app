package createadmin

const (
	insertAdminUser = `
	INSERT INTO users (id, email, password_hash, role, status, email_verified, phone_verified, created_at, updated_at)
	VALUES ($1, $2, $3, 'admin', 'active', TRUE, FALSE, NOW(), NOW())`

	insertAdminProfile = `
	INSERT INTO admin_profiles (user_id, full_name, permission_level, created_at, updated_at)
	VALUES ($1, $2, $3, NOW(), NOW())`

	insertAdminCreatedAudit = `
	INSERT INTO audit_logs (id, actor_user_id, target_user_id, action, entity_type, entity_id, description, ip_address, user_agent, created_at)
	VALUES ($1, NULL, $2, 'admin_action', 'user', $2, 'admin account created with the admin API key', NULLIF($3, '')::inet, NULLIF($4, ''), NOW())`
)
