package distlock

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// RedisLock implements DistributedLock using Redis SET NX
type RedisLock struct {
	client     *redis.Client
	lockValues sync.Map // stores lock values for safe release
}

// NewRedisLock creates a new Redis-based distributed lock
func NewRedisLock(client *redis.Client) *RedisLock {
	return &RedisLock{
		client: client,
	}
}

// Lua script for safe release - only delete if value matches
const releaseLuaScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else
    return 0
end
`

// Acquire tries to acquire a lock with SET NX EX
func (r *RedisLock) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	if r.client == nil {
		log.Println("[DistLock] Redis client is nil, skipping job execution")
		return false, nil
	}

	// Check if Redis is reachable
	if err := r.client.Ping(ctx).Err(); err != nil {
		log.Printf("[DistLock] Redis is unreachable, skipping job execution: %v", err)
		return false, nil
	}

	lockValue := uuid.New().String()

	ok, err := r.client.SetNX(ctx, key, lockValue, ttl).Result()
	if err != nil {
		log.Printf("[DistLock] Error acquiring lock %s, skipping job execution: %v", key, err)
		return false, nil
	}

	if ok {
		r.lockValues.Store(key, lockValue)
		log.Printf("[DistLock] Acquired lock: %s", key)
	} else {
		log.Printf("[DistLock] Lock already held by another instance: %s", key)
	}

	return ok, nil
}

// Release releases the lock using Lua script for atomic check-and-delete
func (r *RedisLock) Release(ctx context.Context, key string) error {
	if r.client == nil {
		return nil
	}

	value, ok := r.lockValues.Load(key)
	if !ok {
		log.Printf("[DistLock] No lock value found for key: %s", key)
		return nil
	}

	result, err := r.client.Eval(ctx, releaseLuaScript, []string{key}, value).Result()
	if err != nil && err != redis.Nil {
		log.Printf("[DistLock] Error releasing lock %s: %v", key, err)
		return err
	}

	r.lockValues.Delete(key)

	if result == int64(1) {
		log.Printf("[DistLock] Released lock: %s", key)
	} else {
		log.Printf("[DistLock] Lock was not held or already expired: %s", key)
	}

	return nil
}

// WithLock executes the function while holding the lock
func (r *RedisLock) WithLock(ctx context.Context, key string, ttl time.Duration, fn func() error) error {
	acquired, err := r.Acquire(ctx, key, ttl)
	if err != nil {
		return err
	}

	if !acquired {
		return nil // another instance is running, skip
	}

	defer func() {
		if releaseErr := r.Release(ctx, key); releaseErr != nil {
			log.Printf("[DistLock] Error releasing lock after execution: %v", releaseErr)
		}
	}()

	return fn()
}

// Extend extends the TTL of an existing lock
func (r *RedisLock) Extend(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	if r.client == nil {
		return false, nil
	}

	value, ok := r.lockValues.Load(key)
	if !ok {
		return false, nil
	}

	// Lua script to extend TTL only if we own the lock
	const extendLuaScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("pexpire", KEYS[1], ARGV[2])
else
    return 0
end
`

	result, err := r.client.Eval(ctx, extendLuaScript, []string{key}, value, ttl.Milliseconds()).Result()
	if err != nil && err != redis.Nil {
		return false, err
	}

	return result == int64(1), nil
}

// ReleaseAll releases all locks held by this instance (for graceful shutdown)
func (r *RedisLock) ReleaseAll(ctx context.Context) {
	r.lockValues.Range(func(key, value interface{}) bool {
		keyStr := key.(string)
		if err := r.Release(ctx, keyStr); err != nil {
			log.Printf("[DistLock] Error releasing lock %s during shutdown: %v", keyStr, err)
		}
		return true
	})
}
