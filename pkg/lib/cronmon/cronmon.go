// Package cronmon adds observability to scheduled jobs.
//
// Every service schedules its crons differently (gocron, a vendored fork, or a
// DI-driven server), but they all share one Redis and one Sentry project. cronmon
// wraps a job so that each run publishes a heartbeat to Redis under a common key
// prefix and a check-in to Sentry. The heartbeat answers "what did this job do
// last time?" and the Sentry check-in answers "did this job run at all?" - a
// heartbeat alone cannot tell a job that is late from a job that is dead.
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
	// HeartbeatKeyPrefix prefixes every cron heartbeat key. It is deliberately
	// service-agnostic: scanning this prefix lists every cron of every service.
	HeartbeatKeyPrefix = "cron:heartbeat:"

	// heartbeatTTL keeps a heartbeat readable long after the job stopped running,
	// so a dead cron stays visible instead of silently disappearing from the list.
	heartbeatTTL = 30 * 24 * time.Hour

	// sentryMinInterval is the shortest schedule Sentry cron monitors can express.
	// Jobs faster than this are tracked by heartbeat only.
	sentryMinInterval = time.Minute
)

// JobSpec describes a scheduled job.
type JobSpec struct {
	// Name identifies the job within its service, e.g. "remove_no_show_queue".
	Name string

	// Interval is how often the job is scheduled. It drives both the Sentry
	// monitor schedule and the staleness threshold reported by Status.
	Interval time.Duration

	// LockTTL bounds how long the distributed lock is held. It must exceed the
	// worst-case runtime, otherwise a second replica can start the job while the
	// first is still working. Defaults to twice Interval.
	LockTTL time.Duration

	// Unlocked runs the job without taking the distributed lock. Set it only for
	// a service that runs a single replica and whose job was never locked to
	// begin with: locking such a job would make an unreachable Redis stop it
	// entirely, trading a working cron for an observable one.
	Unlocked bool

	// Description is a human-readable note shown by the status endpoint.
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

// New creates a Monitor for the given service. service is the name the jobs are
// reported under, e.g. "integration-app". A nil client degrades gracefully: jobs
// still run, they are simply not observable.
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

// Service returns the service name this Monitor reports under.
func (m *Monitor) Service() string {
	return m.service
}

// ReleaseAll releases every lock this instance still holds. Call it during
// graceful shutdown: a lock left behind by a terminating pod blocks the job on
// every replica until its TTL expires.
func (m *Monitor) ReleaseAll(ctx context.Context) {
	m.lock.ReleaseAll(ctx)
}

func (m *Monitor) heartbeatKey(job string) string {
	return HeartbeatKeyPrefix + m.service + ":" + job
}

// monitorSlug builds the Sentry monitor slug. Sentry slugs are lowercase and
// capped at 50 characters.
func (m *Monitor) monitorSlug(job string) string {
	slug := strings.ToLower(m.service + "-" + strings.ReplaceAll(job, "_", "-"))
	if len(slug) > 50 {
		slug = slug[:50]
	}
	return slug
}

// Register publishes a job's schedule before it has ever run, so a cron that
// never starts is still listed by Status instead of being invisible.
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

// Wrap adapts a job to the func() signature schedulers expect, running it
// through Run.
func (m *Monitor) Wrap(spec JobSpec, fn func() error) func() {
	return func() {
		m.Run(context.Background(), spec, fn)
	}
}

// Run executes fn while holding the job's distributed lock, recording the
// outcome to Redis and Sentry. When another replica already holds the lock the
// job is skipped without recording anything, so the heartbeat always reflects
// the replica that actually did the work.
//
// A panic inside fn is recovered and recorded as a failure: an unobserved panic
// in a scheduled job is exactly the failure mode this package exists to catch.
func (m *Monitor) Run(ctx context.Context, spec JobSpec, fn func() error) {
	if !spec.Unlocked {
		key := distlock.BuildLockKey(m.service + ":" + spec.Name)
		acquired, err := m.lock.Acquire(ctx, key, spec.lockTTL())
		if err != nil || !acquired {
			return
		}
		defer func() {
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

// runProtected converts a panic into an error so one broken job cannot take the
// whole scheduler down, and so the failure is recorded like any other.
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

// checkInWithID reports to Sentry. CaptureCheckIn is a no-op when Sentry is not
// initialised, so no feature flag is needed here.
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
