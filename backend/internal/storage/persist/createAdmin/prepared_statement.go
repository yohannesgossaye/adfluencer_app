package createadmin

const (
	createAdmin = `
		INSERT INTO users (
			id, email, phone_number, password_hash, role, status,
			email_verified, phone_verified, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		);
	`

	insertAdminProfile = `
		INSERT INTO admin_profiles (
			id, user_id, full_name, permission_level, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6
		);
	`
)
