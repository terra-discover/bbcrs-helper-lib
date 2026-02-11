package distlock

import (
	"context"
	"time"
)

// DistributedLock defines the interface for distributed locking
type DistributedLock interface {
	// Acquire tries to acquire a lock with the given key and TTL.
	// Returns true if lock was acquired, false if already locked by another process.
	Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error)

	// Release releases the lock for the given key.
	// Only releases if the lock is owned by this instance.
	Release(ctx context.Context, key string) error

	// WithLock executes the function while holding the lock.
	// Automatically acquires and releases the lock.
	// If lock cannot be acquired, the function is skipped.
	WithLock(ctx context.Context, key string, ttl time.Duration, fn func() error) error
}

// LockKeyPrefix is the prefix for all cron lock keys
const LockKeyPrefix = "cron:lock:"

// BuildLockKey builds a lock key with the standard prefix
func BuildLockKey(jobName string) string {
	return LockKeyPrefix + jobName
}
