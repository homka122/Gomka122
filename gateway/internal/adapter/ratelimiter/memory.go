package ratelimiter

import (
	"context"
	"sync"
	"time"
)

type clientBucket struct {
	tokens     float64
	lastRefill time.Time
}

type MemoryBucketRateLimiter struct {
	capacity  int
	reqPerSec float64
	clients   map[string]clientBucket
	mu        sync.Mutex
}

func NewMemoryBucketRateLimiter(capacity int, reqPerSec float64) *MemoryBucketRateLimiter {
	return &MemoryBucketRateLimiter{
		capacity:  capacity,
		reqPerSec: reqPerSec,
		clients:   map[string]clientBucket{},
		mu:        sync.Mutex{},
	}
}

func (rl *MemoryBucketRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	key = "rl:ip:" + key
	now := time.Now()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	client, ok := rl.clients[key]

	if !ok {
		client = clientBucket{
			tokens:     float64(rl.capacity),
			lastRefill: now,
		}
	}

	newTokens := now.Sub(client.lastRefill).Seconds() * rl.reqPerSec
	client.tokens = min(float64(rl.capacity), client.tokens+newTokens)
	client.lastRefill = time.Now()

	if client.tokens < 1.0 {
		rl.clients[key] = client
		return false, nil
	}

	client.tokens -= 1.0
	rl.clients[key] = client

	return true, nil
}
