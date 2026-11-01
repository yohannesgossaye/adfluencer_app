-- Session rotation hardening (Phase 1A). Safe to re-run.
--
-- refresh_ttl_seconds: lifetime chosen at login (e.g. 30 days with "remember me"),
--   carried over every time the refresh token is rotated.
-- rotated_at: when this refresh token was exchanged for a new one. A rotated
--   token stays usable for a short grace window so that two browser tabs
--   refreshing at the same moment do not log each other out.

ALTER TABLE user_sessions ADD COLUMN IF NOT EXISTS refresh_ttl_seconds INTEGER;
ALTER TABLE user_sessions ADD COLUMN IF NOT EXISTS rotated_at TIMESTAMPTZ;

-- Every refresh and logout looks sessions up by token hash.
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_sessions_refresh_token_hash
    ON user_sessions (refresh_token_hash);
