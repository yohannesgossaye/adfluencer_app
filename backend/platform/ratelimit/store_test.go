package ratelimit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

// The memory store and the Redis store must behave the same.
func TestStores(t *testing.T) {
	mr := miniredis.RunT(t)
	rs, err := NewRedisStore("redis://" + mr.Addr() + "/0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { rs.Close() })

	clock := time.Date(2026, 11, 6, 9, 0, 0, 0, time.UTC)
	ms := NewMemoryStore()
	ms.now = func() time.Time { return clock }

	for name, tc := range map[string]struct {
		store   Store
		advance func(time.Duration)
	}{
		"memory": {ms, func(d time.Duration) { clock = clock.Add(d) }},
		"redis":  {rs, mr.FastForward},
	} {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			s := tc.store

			// A fixed window: later hits keep the first expiry.
			for want := int64(1); want <= 3; want++ {
				got, err := s.Incr(ctx, "k", 10*time.Minute)
				if err != nil || got != want {
					t.Fatalf("incr %d: got %d, %v", want, got, err)
				}
				tc.advance(time.Minute)
			}
			v, left, err := s.Get(ctx, "k")
			if err != nil || v != 3 || left > 7*time.Minute+time.Second || left < 6*time.Minute {
				t.Fatalf("get: %d, %v left, %v; want 3 with ~7m left", v, left, err)
			}
			tc.advance(8 * time.Minute)
			if v, _, _ := s.Get(ctx, "k"); v != 0 {
				t.Fatalf("expired key still counts %d", v)
			}
			if got, _ := s.Incr(ctx, "k", time.Minute); got != 1 {
				t.Fatalf("new window starts at %d, want 1", got)
			}

			// Set replaces, Del removes, missing keys read as zero.
			if err := s.Set(ctx, "lock", 2, 5*time.Minute); err != nil {
				t.Fatal(err)
			}
			if v, left, _ := s.Get(ctx, "lock"); v != 2 || left <= 4*time.Minute {
				t.Fatalf("set: %d with %v left", v, left)
			}
			if err := s.Del(ctx, "lock", "k", "missing"); err != nil {
				t.Fatal(err)
			}
			if v, left, err := s.Get(ctx, "lock"); v != 0 || left != 0 || err != nil {
				t.Fatalf("deleted key reads %d, %v, %v", v, left, err)
			}
		})
	}
}

type brokenStore struct{}

var errDown = errors.New("redis down")

func (brokenStore) Incr(context.Context, string, time.Duration) (int64, error) { return 0, errDown }
func (brokenStore) Get(context.Context, string) (int64, time.Duration, error) {
	return 0, 0, errDown
}
func (brokenStore) Set(context.Context, string, int64, time.Duration) error { return errDown }
func (brokenStore) Del(context.Context, ...string) error                    { return errDown }

// When Redis is down, limits keep working from memory and nothing errors.
func TestFallbackKeepsCounting(t *testing.T) {
	ctx := context.Background()
	failures := 0
	f := Fallback{Primary: brokenStore{}, Backup: NewMemoryStore(), OnError: func(string, error) { failures++ }}
	for want := int64(1); want <= 2; want++ {
		if got, err := f.Incr(ctx, "k", time.Minute); err != nil || got != want {
			t.Fatalf("incr: %d, %v", got, err)
		}
	}
	if v, _, err := f.Get(ctx, "k"); err != nil || v != 2 {
		t.Fatalf("get: %d, %v", v, err)
	}
	if err := f.Del(ctx, "k"); err != nil {
		t.Fatal(err)
	}
	if v, _, _ := f.Get(ctx, "k"); v != 0 {
		t.Fatalf("del left %d", v)
	}
	if failures == 0 {
		t.Fatal("redis failures were not reported")
	}
}
