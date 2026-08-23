package cronmon

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2/utils"
)

func TestClassify(t *testing.T) {
	now := time.Now()
	recent := now.Add(-30 * time.Second)
	old := now.Add(-10 * time.Minute)

	// A job that has never started is registered but not yet observed.
	utils.AssertEqual(t, HealthNeverRun, classify(JobStatus{IntervalSeconds: 60}))

	// Started but never succeeded means the very first run broke.
	utils.AssertEqual(t, HealthFailing, classify(JobStatus{
		IntervalSeconds: 60,
		LastStart:       &now,
	}))

	// A recent failure outranks an older success.
	utils.AssertEqual(t, HealthFailing, classify(JobStatus{
		IntervalSeconds:     60,
		LastStart:           &now,
		LastSuccess:         &recent,
		ConsecutiveFailures: 1,
	}))

	// Succeeding within the interval is healthy.
	utils.AssertEqual(t, HealthOK, classify(JobStatus{
		IntervalSeconds: 60,
		LastStart:       &now,
		LastSuccess:     &recent,
	}))

	// No success for more than two intervals means the job is likely dead.
	utils.AssertEqual(t, HealthStale, classify(JobStatus{
		IntervalSeconds: 60,
		LastStart:       &old,
		LastSuccess:     &old,
	}))

	// Without a known interval staleness cannot be judged, so a success stands.
	utils.AssertEqual(t, HealthOK, classify(JobStatus{LastSuccess: &old}))
}

func TestMonitorSlug(t *testing.T) {
	m := New(nil, "integration-app")

	utils.AssertEqual(t, "integration-app-remove-no-show-queue", m.monitorSlug("remove_no_show_queue"))
	utils.AssertEqual(t, true, len(m.monitorSlug(strings.Repeat("long_job_name", 10))) <= 50)
}

func TestHeartbeatKey(t *testing.T) {
	m := New(nil, "integration-app")
	utils.AssertEqual(t, "cron:heartbeat:integration-app:remove_no_show_queue", m.heartbeatKey("remove_no_show_queue"))
}

func TestLockTTLDefaultsToTwiceInterval(t *testing.T) {
	utils.AssertEqual(t, 10*time.Minute, JobSpec{Interval: 5 * time.Minute}.lockTTL())
	utils.AssertEqual(t, time.Minute, JobSpec{Interval: 5 * time.Minute, LockTTL: time.Minute}.lockTTL())
}

// Without Redis the lock cannot be acquired, so the job is skipped rather than
// run unlocked on every replica at once.
func TestRunWithoutRedisSkips(t *testing.T) {
	m := New(nil, "integration-app")

	executed := false
	m.Run(context.Background(), JobSpec{Name: "job", Interval: time.Minute}, func() error {
		executed = true
		return nil
	})

	utils.AssertEqual(t, false, executed)
}

// An unlocked job must keep running even with Redis unavailable: monitoring a
// job must never be the reason it stops.
func TestUnlockedRunsWithoutRedis(t *testing.T) {
	m := New(nil, "queue-app")

	executed := false
	m.Run(context.Background(), JobSpec{Name: "job", Interval: time.Minute, Unlocked: true}, func() error {
		executed = true
		return nil
	})

	utils.AssertEqual(t, true, executed)
}

func TestRunProtectedTurnsPanicIntoError(t *testing.T) {
	// Must not propagate: a panicking job would otherwise take down the scheduler.
	err := runProtected(func() error {
		panic("boom")
	})

	utils.AssertEqual(t, true, err != nil)
	utils.AssertEqual(t, true, strings.Contains(err.Error(), "boom"))

	utils.AssertEqual(t, nil, runProtected(func() error { return nil }))
}

func TestStatusWithoutRedisIsEmpty(t *testing.T) {
	statuses, err := Status(context.Background(), nil, "")
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 0, len(statuses))
}

func TestBuildStatusRecoversNamesFromKey(t *testing.T) {
	status := buildStatus("cron:heartbeat:payment-app:check_transaction", map[string]string{
		fieldRunCount: "3",
	})

	utils.AssertEqual(t, "payment-app", status.Service)
	utils.AssertEqual(t, "check_transaction", status.Job)
	utils.AssertEqual(t, int64(3), status.RunCount)
}
