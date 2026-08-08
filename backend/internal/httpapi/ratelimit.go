package httpapi

import (
	"sync"
	"time"
)

type LoginLimiter interface {
	Allow(key string) bool
}

type tokenBucket struct {
	tokens     float64
	capacity   float64
	refillRate float64
	lastRefill time.Time
}

type tokenBucketLimiter struct {
	mu    sync.Mutex
	limit float64
	rate  float64
	now   func() time.Time
	keys  map[string]*tokenBucket
}

func NewTokenBucketLimiter() *tokenBucketLimiter {
	return &tokenBucketLimiter{
		limit: 5,
		rate:  0.1,
		now:   time.Now,
		keys:  make(map[string]*tokenBucket),
	}
}

func NewTokenBucketLimiterWith(capacity, refillPerSecond float64) *tokenBucketLimiter {
	return &tokenBucketLimiter{
		limit: capacity,
		rate:  refillPerSecond,
		now:   time.Now,
		keys:  make(map[string]*tokenBucket),
	}
}

func (l *tokenBucketLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	b, ok := l.keys[key]
	if !ok {
		if len(l.keys) >= 10000 {
			for k, old := range l.keys {
				if now.Sub(old.lastRefill) > time.Hour {
					delete(l.keys, k)
				}
			}
		}
		b = &tokenBucket{tokens: l.limit, capacity: l.limit, refillRate: l.rate, lastRefill: now}
		l.keys[key] = b
	}

	elapsed := now.Sub(b.lastRefill).Seconds()
	if elapsed > 0 {
		b.tokens = min(l.limit, b.tokens+elapsed*b.refillRate)
		b.lastRefill = now
	}
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}
