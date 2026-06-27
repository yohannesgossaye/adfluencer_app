-- Phone OTP verification tokens (mirrors email_verification_tokens).
-- Safe to re-run (IF NOT EXISTS). For existing DB: make migrate-phone-verification

CREATE TABLE IF NOT EXISTS phone_verification_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_phone_verification_tokens_user_id ON phone_verification_tokens(user_id);
