package ratelimit

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStore keeps counters in Redis, shared by every API server.
type RedisStore struct {
	client *redis.Client
}

// NewRedisStore connects using a URL such as redis://redis:6379/0. Calls time
// out quickly: a slow Redis must not hold up a login.
func NewRedisStore(url string) (*RedisStore, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	opts.DialTimeout = time.Second
	opts.ReadTimeout = 500 * time.Millisecond
	opts.WriteTimeout = 500 * time.Millisecond
	opts.MaxRetries = 1
	return &RedisStore{client: redis.NewClient(opts)}, nil
}

func (s *RedisStore) Ping(ctx context.Context) error { return s.client.Ping(ctx).Err() }

func (s *RedisStore) Close() error { return s.client.Close() }

// incr adds one and sets the expiry only on a new key, in one atomic step.
var incr = redis.NewScript(`
local v = redis.call('INCR', KEYS[1])
if v == 1 then redis.call('PEXPIRE', KEYS[1], ARGV[1]) end
return v`)

func (s *RedisStore) Incr(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	return incr.Run(ctx, s.client, []string{key}, ttl.Milliseconds()).Int64()
}

func (s *RedisStore) Get(ctx context.Context, key string) (int64, time.Duration, error) {
	pipe := s.client.Pipeline()
	get := pipe.Get(ctx, key)
	ttl := pipe.PTTL(ctx, key)
	if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		return 0, 0, err
	}
	v, err := get.Int64()
	if errors.Is(err, redis.Nil) {
		return 0, 0, nil
	}
	if err != nil {
		return 0, 0, err
	}
	left := ttl.Val()
	if left < 0 {
		left = 0
	}
	return v, left, nil
}

func (s *RedisStore) Set(ctx context.Context, key string, value int64, ttl time.Duration) error {
	return s.client.Set(ctx, key, value, ttl).Err()
}

func (s *RedisStore) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return s.client.Del(ctx, keys...).Err()
}
