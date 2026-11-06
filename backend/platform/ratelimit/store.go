package ratelimit

import (
	"context"
	"sync"
	"time"
)

// Store keeps short-lived counters. Every value expires on its own, so nothing
// needs cleaning up and losing them all (a restart) only resets the limits.
type Store interface {
	// Incr adds one to key and returns the new value. A new key lives for ttl;
	// later increments keep that expiry, so the count covers a fixed window.
	Incr(ctx context.Context, key string, ttl time.Duration) (int64, error)
	// Get returns the value and the time it has left (0, 0 when there is none).
	Get(ctx context.Context, key string) (int64, time.Duration, error)
	// Set stores value for ttl, replacing what was there.
	Set(ctx context.Context, key string, value int64, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
}

// MemoryStore keeps counters in this process. It is the backup when Redis is
// unreachable, and the store for tests and local runs without Redis.
type MemoryStore struct {
	mu        sync.Mutex
	items     map[string]memoryItem
	now       func() time.Time
	lastSweep time.Time
}

type memoryItem struct {
	value   int64
	expires time.Time
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{items: map[string]memoryItem{}, now: time.Now}
}

// live returns the item if it has not expired; the caller holds the lock.
func (s *MemoryStore) live(key string, now time.Time) (memoryItem, bool) {
	it, ok := s.items[key]
	if ok && !now.Before(it.expires) {
		delete(s.items, key)
		return memoryItem{}, false
	}
	return it, ok
}

// sweep drops expired items now and then, so the map cannot grow forever.
func (s *MemoryStore) sweep(now time.Time) {
	if now.Sub(s.lastSweep) < time.Minute {
		return
	}
	s.lastSweep = now
	for k, it := range s.items {
		if !now.Before(it.expires) {
			delete(s.items, k)
		}
	}
}

func (s *MemoryStore) Incr(_ context.Context, key string, ttl time.Duration) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	s.sweep(now)
	it, ok := s.live(key, now)
	if !ok {
		it = memoryItem{expires: now.Add(ttl)}
	}
	it.value++
	s.items[key] = it
	return it.value, nil
}

func (s *MemoryStore) Get(_ context.Context, key string) (int64, time.Duration, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	it, ok := s.live(key, now)
	if !ok {
		return 0, 0, nil
	}
	return it.value, it.expires.Sub(now), nil
}

func (s *MemoryStore) Set(_ context.Context, key string, value int64, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = memoryItem{value: value, expires: s.now().Add(ttl)}
	return nil
}

func (s *MemoryStore) Del(_ context.Context, keys ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, k := range keys {
		delete(s.items, k)
	}
	return nil
}

// Fallback uses primary (Redis) and switches to backup (memory) for any call
// primary fails, so an unreachable Redis never blocks logins and the limits
// keep working on this server. onError is told about each failure.
type Fallback struct {
	Primary, Backup Store
	OnError         func(op string, err error)
}

func (f Fallback) failed(op string, err error) {
	if f.OnError != nil {
		f.OnError(op, err)
	}
}

func (f Fallback) Incr(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	v, err := f.Primary.Incr(ctx, key, ttl)
	if err != nil {
		f.failed("incr", err)
		return f.Backup.Incr(ctx, key, ttl)
	}
	return v, nil
}

func (f Fallback) Get(ctx context.Context, key string) (int64, time.Duration, error) {
	v, ttl, err := f.Primary.Get(ctx, key)
	if err != nil {
		f.failed("get", err)
		return f.Backup.Get(ctx, key)
	}
	return v, ttl, nil
}

func (f Fallback) Set(ctx context.Context, key string, value int64, ttl time.Duration) error {
	if err := f.Primary.Set(ctx, key, value, ttl); err != nil {
		f.failed("set", err)
		return f.Backup.Set(ctx, key, value, ttl)
	}
	return nil
}

func (f Fallback) Del(ctx context.Context, keys ...string) error {
	// Clear both: a value may have been written to the backup during an outage.
	_ = f.Backup.Del(ctx, keys...)
	if err := f.Primary.Del(ctx, keys...); err != nil {
		f.failed("del", err)
	}
	return nil
}
