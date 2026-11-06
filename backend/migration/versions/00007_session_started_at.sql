-- +goose Up
-- Bank-style sessions: every session ends after inactivity and after a
-- maximum age counted from login. started_at is that login time; it is
-- copied to every refreshed session row.
ALTER TABLE user_sessions ADD COLUMN IF NOT EXISTS started_at TIMESTAMPTZ;
UPDATE user_sessions SET started_at = created_at WHERE started_at IS NULL;
ALTER TABLE user_sessions ALTER COLUMN started_at SET DEFAULT now();
ALTER TABLE user_sessions ALTER COLUMN started_at SET NOT NULL;

-- Everyone signs in once more, so the new rules apply to every session at
-- once (older sessions were issued for 7 or 30 days).
UPDATE user_sessions SET revoked_at = NOW() WHERE revoked_at IS NULL;

-- +goose Down
ALTER TABLE user_sessions DROP COLUMN IF EXISTS started_at;
