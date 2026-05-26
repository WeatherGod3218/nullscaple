package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const MILLISECONDS = 1000

type TokenBucket struct {
	client   *redis.Client
	rate     float64
	capacity float64
}

func NewTokenBucket(rate float64, capacity float64) *TokenBucket {
	return &TokenBucket{client: client, rate: rate, capacity: capacity}
}

func (tb *TokenBucket) Allow(ctx context.Context, key string) (bool, int64, error) {
	now := time.Now().UnixMilli()
	ipKey := fmt.Sprintf("ratelimit:%s", key)

	script := redis.NewScript(`
		local key = KEYS[1]
		local now = tonumber(ARGV[1])
		local rate = tonumber(ARGV[2])
		local capacity = tonumber(ARGV[3])
		local ttl = tonumber(ARGV[4])

		local bucketData = redis.call("HMGET", key, "tokens", "last_refill")
		local tokens = tonumber(bucketData[1]) or capacity
		local lastRefill = tonumber(bucketData[2]) or now

		local elapsed = math.max(0, (now - lastRefill) / 1000)
		tokens = math.min(capacity, tokens + (elapsed * rate))

		if tokens < 1 then
			redis.call("HMSET", key, "tokens", tokens, "last_refill", now)
			redis.call("PEXPIRE", key, ttl)
			return 0
		end

		tokens = tokens - 1
		redis.call("HMSET", key, "tokens", tokens, "last_refill", now)
		redis.call("PEXPIRE", key, ttl)
		return {1, tokens, }
	`)

	ttlMs := int64((tb.capacity / tb.rate) * MILLISECONDS * 2)

	result, err := script.Run(ctx, tb.client, []string{ipKey}, now, tb.rate, tb.capacity, ttlMs).Int64Slice()
	if err != nil {
		return false, 0, fmt.Errorf("redis script: %w", err)
	}

	return (result[0] == 1), result[1], nil
}
