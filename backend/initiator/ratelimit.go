package initiator

import (
	"context"
	"sync"
	"time"

	"adinfluence_service_app/platform/ratelimit"
	"adinfluence_service_app/platform/utils"

	"go.uber.org/zap"
)

// initRateLimitStore keeps the login attempt counters in Redis when REDIS_URL
// is set, with memory as the backup for any moment Redis can't be reached.
// Logins never fail because Redis is down.
func initRateLimitStore(ctx context.Context, url string, logger utils.Logger) ratelimit.Store {
	memory := ratelimit.NewMemoryStore()
	if url == "" {
		logger.Warn(ctx, "REDIS_URL not set — login attempt limits are kept in memory (fine for one server)")
		return memory
	}
	rs, err := ratelimit.NewRedisStore(url)
	if err != nil {
		logger.Error(ctx, "invalid REDIS_URL — login attempt limits are kept in memory", zap.Error(err))
		return memory
	}
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := rs.Ping(pingCtx); err != nil {
		logger.Warn(ctx, "redis not reachable yet — using memory until it is", zap.Error(err))
	} else {
		logger.Info(ctx, "redis ready for login attempt limits")
	}

	// Report Redis failures at most once a minute, not once per request.
	var mu sync.Mutex
	var lastReport time.Time
	return ratelimit.Fallback{
		Primary: rs,
		Backup:  memory,
		OnError: func(op string, err error) {
			mu.Lock()
			defer mu.Unlock()
			if time.Since(lastReport) < time.Minute {
				return
			}
			lastReport = time.Now()
			logger.Warn(context.Background(), "redis unavailable — login attempt limits use memory meanwhile", zap.String("op", op), zap.Error(err))
		},
	}
}
