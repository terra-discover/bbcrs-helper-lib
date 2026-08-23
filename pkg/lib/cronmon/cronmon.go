// Package cronmon adds observability to scheduled jobs: a Redis heartbeat for
// what a job did last time, and a Sentry check-in for whether it ran at all.
package cronmon

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v8"
	"github.com/terra-discover/bbcrs-helper-lib/pkg/lib/distlock"
)

const (
	// Service-agnostic on purpose: scanning it lists every cron of every service.
	HeartbeatKeyPrefix = "cron:heartbeat:"

	// Far longer than any schedule, so a dead cron stays visible.
	heartbeatTTL = 30 * 24 * time.Hour

	// Shortest schedule Sentry cron monitors accept.
	sentryMinInterval = time.Minute
)

// JobSpec describes a scheduled job.
type JobSpec struct {
	Name string

	// Drives the Sentry monitor schedule and the staleness threshold.
	Interval time.Duration

	// Only has to outlast the gap between two renewals, not the whole run.
	// Defaults to 2*Interval.
	LockTTL time.Duration

	// Run anyway when the lock backend is unreachable. Still skipped when
	// another replica holds the lock, so rolling deploys stay protected without
	// making Redis a new point of failure.
	FailOpen bool

	Description string
}

func (s JobSpec) lockTTL() time.Duration {
	if s.LockTTL > 0 {
		return s.LockTTL
	}
	return 2 * s.Interval
}

// Monitor wraps scheduled jobs with distributed locking, Redis heartbeats and
// Sentry check-ins.
type Monitor struct {
	client   *redis.Client
	lock     *distlock.RedisLock
	service  string
	instance string
}

// New creates a Monitor for the given service. A nil client degrades
// gracefully: jobs still run, they are simply not observable.
func New(client *redis.Client, service string) *Monitor {
	instance, err := os.Hostname()
	if err != nil {
		instance = "unknown"
	}
	return &Monitor{
		client:   client,
		lock:     distlock.NewRedisLock(client),
		service:  service,
		instance: instance,
	}
}

func (m *Monitor) Service() string {
	return m.service
}

// ReleaseAll releases every lock this instance holds. Call it during graceful
// shutdown: a lock left by a terminating pod blocks the job on every replica
// until its TTL expires.
func (m *Monitor) ReleaseAll(ctx context.Context) {
	m.lock.ReleaseAll(ctx)
}

func (m *Monitor) heartbeatKey(job string) string {
	return HeartbeatKeyPrefix + m.service + ":" + job
}

// Sentry slugs must be lowercase and at most 50 characters.
func (m *Monitor) monitorSlug(job string) string {
	slug := strings.ToLower(m.service + "-" + strings.ReplaceAll(job, "_", "-"))
	if len(slug) > 50 {
		slug = slug[:50]
	}
	return slug
}

// Register publishes a job's schedule before it has ever run, so a cron that
// never starts is listed by Status instead of being invisible.
func (m *Monitor) Register(ctx context.Context, spec JobSpec) {
	if m.client == nil {
		return
	}

	key := m.heartbeatKey(spec.Name)
	fields := map[string]interface{}{
		fieldService:      m.service,
		fieldJob:          spec.Name,
		fieldIntervalSec:  int64(spec.Interval.Seconds()),
		fieldDescription:  spec.Description,
		fieldRegisteredAt: formatTime(time.Now().UTC()),
	}
	if err := m.client.HSet(ctx, key, fields).Err(); err != nil {
		log.Printf("[CronMon] failed registering job %s: %v", spec.Name, err)
		return
	}
	m.client.Expire(ctx, key, heartbeatTTL)
}

// Wrap adapts a job to the func() signature schedulers expect.
func (m *Monitor) Wrap(spec JobSpec, fn func() error) func() {
	return func() {
		m.Run(context.Background(), spec, fn)
	}
}

// Run executes fn while holding the job's lock, recording the outcome to Redis
// and Sentry. A run skipped because another replica holds the lock records
// nothing, so the heartbeat always reflects the replica that did the work.
func (m *Monitor) Run(ctx context.Context, spec JobSpec, fn func() error) {
	key := distlock.BuildLockKey(m.service + ":" + spec.Name)
	acquired, err := m.lock.Acquire(ctx, key, spec.lockTTL())
	if err != nil || !acquired {
		// Acquire reports "held by another replica" and "backend unreachable"
		// identically, and those need opposite responses.
		if !spec.FailOpen || m.lockReachable(ctx) {
			return
		}
		log.Printf("[CronMon] job %s/%s running unlocked: lock backend unreachable", m.service, spec.Name)
	} else {
		stopRenew := m.renewLock(ctx, key, spec)
		defer func() {
			stopRenew()
			if releaseErr := m.lock.Release(ctx, key); releaseErr != nil {
				log.Printf("[CronMon] failed releasing lock for job %s: %v", spec.Name, releaseErr)
			}
		}()
	}

	start := time.Now()
	checkInID := m.checkIn(spec, sentry.CheckInStatusInProgress, 0)
	m.markStarted(ctx, spec, start)

	runErr := runProtected(fn)

	duration := time.Since(start)
	if runErr != nil {
		log.Printf("[CronMon] job %s/%s failed after %s: %v", m.service, spec.Name, duration, runErr)
		m.checkInWithID(checkInID, spec, sentry.CheckInStatusError, duration)
		m.markFinished(ctx, spec, duration, runErr)
		return
	}

	log.Printf("[CronMon] job %s/%s ok in %s", m.service, spec.Name, duration)
	m.checkInWithID(checkInID, spec, sentry.CheckInStatusOK, duration)
	m.markFinished(ctx, spec, duration, nil)
}

func (m *Monitor) lockReachable(ctx context.Context) bool {
	if m.client == nil {
		return false
	}
	return m.client.Ping(ctx).Err() == nil
}

// renewLock keeps the lock alive while the job runs and returns a stop
// function. Without it a job outliving its TTL loses the lock mid-flight and a
// second replica starts the same work.
func (m *Monitor) renewLock(ctx context.Context, key string, spec JobSpec) func() {
	ttl := spec.lockTTL()
	interval := ttl / 3
	if interval <= 0 {
		return func() {}
	}

	done := make(chan struct{})
	stopped := make(chan struct{})

	go func() {
		defer close(stopped)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				extended, err := m.lock.Extend(ctx, key, ttl)
				if err != nil {
					log.Printf("[CronMon] failed extending lock for job %s: %v", spec.Name, err)
					continue
				}
				if !extended {
					log.Printf("[CronMon] lost lock for job %s while it was still running", spec.Name)
					return
				}
			}
		}
	}()

	return func() {
		close(done)
		<-stopped
	}
}

// runProtected converts a panic into an error so one broken job cannot take the
// whole scheduler down.
func runProtected(fn func() error) (runErr error) {
	defer func() {
		if r := recover(); r != nil {
			runErr = fmt.Errorf("panic: %v", r)
			sentry.CurrentHub().Recover(r)
		}
	}()
	return fn()
}

func (m *Monitor) checkIn(spec JobSpec, status sentry.CheckInStatus, duration time.Duration) sentry.EventID {
	return m.checkInWithID("", spec, status, duration)
}

// CaptureCheckIn is a no-op when Sentry is not initialised, so no feature flag
// is needed here.
func (m *Monitor) checkInWithID(id sentry.EventID, spec JobSpec, status sentry.CheckInStatus, duration time.Duration) sentry.EventID {
	if spec.Interval < sentryMinInterval {
		return ""
	}

	minutes := int64(spec.Interval.Minutes())
	if minutes < 1 {
		minutes = 1
	}
	maxRuntime := int64(spec.lockTTL().Minutes())
	if maxRuntime < 1 {
		maxRuntime = 1
	}

	eventID := sentry.CaptureCheckIn(&sentry.CheckIn{
		ID:          id,
		MonitorSlug: m.monitorSlug(spec.Name),
		Status:      status,
		Duration:    duration,
	}, &sentry.MonitorConfig{
		Schedule:      sentry.IntervalSchedule(minutes, sentry.MonitorScheduleUnitMinute),
		CheckInMargin: minutes,
		MaxRuntime:    maxRuntime,
		Timezone:      "UTC",
	})
	if eventID == nil {
		return ""
	}
	return *eventID
}
