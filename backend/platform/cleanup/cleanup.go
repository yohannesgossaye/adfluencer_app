// Package cleanup periodically deletes auth rows that can never be used
// again, so these tables stay small. Every rule keeps at least a day of
// history beyond the longest window the auth code looks back over (2h).
package cleanup

import (
	"context"
	"database/sql"
	"time"

	"adinfluence_service_app/platform/utils"

	"go.uber.org/zap"
)

// Rules are run in order; each is one DELETE.
var Rules = []struct{ Name, SQL string }{
	{"expired or revoked sessions", `DELETE FROM user_sessions
		WHERE expires_at < now() - INTERVAL '1 day' OR revoked_at < now() - INTERVAL '1 day'`},
	{"password reset tokens", `DELETE FROM password_reset_tokens WHERE expires_at < now() - INTERVAL '1 day'`},
	{"email verification codes", `DELETE FROM email_verification_tokens WHERE expires_at < now() - INTERVAL '1 day'`},
	{"phone verification codes", `DELETE FROM phone_verification_tokens WHERE expires_at < now() - INTERVAL '1 day'`},
	{"rate limit events", `DELETE FROM auth_rate_limit_events WHERE created_at < now() - INTERVAL '1 day'`},
	{"stale login lockouts", `DELETE FROM auth_login_lockouts
		WHERE (locked_until IS NULL OR locked_until < now()) AND updated_at < now() - INTERVAL '1 day'`},
}

// RunOnce applies every rule and logs how many rows each removed.
func RunOnce(ctx context.Context, db *sql.DB, log utils.Logger) {
	for _, r := range Rules {
		qctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		res, err := db.ExecContext(qctx, r.SQL)
		cancel()
		if err != nil {
			log.Warn(ctx, "cleanup failed", zap.String("rule", r.Name), zap.Error(err))
			continue
		}
		if n, _ := res.RowsAffected(); n > 0 {
			log.Info(ctx, "cleanup removed rows", zap.String("rule", r.Name), zap.Int64("rows", n))
		}
	}
}

// Start runs the cleanup shortly after startup and then every interval
// until ctx is cancelled.
func Start(ctx context.Context, db *sql.DB, log utils.Logger, interval time.Duration) {
	go func() {
		timer := time.NewTimer(time.Minute)
		defer timer.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				RunOnce(ctx, db, log)
				timer.Reset(interval)
			}
		}
	}()
}
