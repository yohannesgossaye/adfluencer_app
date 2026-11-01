-- +goose Up
-- Support the hourly cleanup job (platform/cleanup): it deletes rows by these
-- timestamps, so each delete is an index range scan instead of a full scan.
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires_at ON user_sessions (expires_at);
CREATE INDEX IF NOT EXISTS idx_auth_rate_limit_events_created_at ON auth_rate_limit_events (created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_auth_rate_limit_events_created_at;
DROP INDEX IF EXISTS idx_user_sessions_expires_at;
