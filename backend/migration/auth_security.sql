-- Auth security: login lockouts + rate-limit event log
-- Safe to re-run (IF NOT EXISTS). For existing DB: make migrate-security

CREATE TABLE IF NOT EXISTS auth_login_lockouts (
    email VARCHAR(255) PRIMARY KEY,

    failed_attempts INT NOT NULL DEFAULT 0,
    lockout_tier INT NOT NULL DEFAULT 0,
    locked_until TIMESTAMPTZ,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS auth_rate_limit_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    scope VARCHAR(50) NOT NULL,
    rate_key VARCHAR(255) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_auth_rate_limit_scope_key_created
    ON auth_rate_limit_events (scope, rate_key, created_at DESC);
