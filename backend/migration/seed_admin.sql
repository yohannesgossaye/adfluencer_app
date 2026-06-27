-- Dev admin seed (run once: make seed-admin)
-- Login: admin@adfluence.com / admin1111
-- Use portal: "admin" on POST /api/v1/auth/login

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
) VALUES (
    'a1000000-0000-4000-8000-000000000001',
    'admin@adfluence.com',
    NULL,
    '$argon2id$v=19$m=65536,t=4,p=2$SDDP2Vq45jYwGPUG3auV3g$B/d7aBWuRin0+R9uU0E4TmTxvlY4oY2Xqa+/03fyltA',
    'admin',
    'active',
    TRUE,
    FALSE,
    NOW(),
    NOW()
) ON CONFLICT (email) DO NOTHING;

INSERT INTO admin_profiles (
    id,
    user_id,
    full_name,
    permission_level,
    created_at,
    updated_at
) VALUES (
    'a2000000-0000-4000-8000-000000000001',
    'a1000000-0000-4000-8000-000000000001',
    'Adfluence Admin',
    'super_admin',
    NOW(),
    NOW()
) ON CONFLICT (user_id) DO NOTHING;
