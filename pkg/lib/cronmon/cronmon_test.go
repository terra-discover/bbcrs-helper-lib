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

	utils.AssertEqual(t, HealthNeverRun, classify(JobStatus{IntervalSeconds: 60}))

	utils.AssertEqual(t, HealthFailing, classify(JobStatus{
		IntervalSeconds: 60,
		LastStart:       &now,
	}))

	utils.AssertEqual(t, HealthFailing, classify(JobStatus{
		IntervalSeconds:     60,
		LastStart:           &now,
		LastSuccess:         &recent,
		ConsecutiveFailures: 1,
	}))

	utils.AssertEqual(t, HealthOK, classify(JobStatus{
		IntervalSeconds: 60,
		LastStart:       &now,
		LastSuccess:     &recent,
	}))

	utils.AssertEqual(t, HealthStale, classify(JobStatus{
		IntervalSeconds: 60,
		LastStart:       &old,
		LastSuccess:     &old,
	}))

	// Without an interval, staleness cannot be judged and a success stands.
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

// Fail-closed: no lock means no run, for every job. A duplicate enqueue sends
// a customer a second email, while a skipped tick only delays one.
func TestRunWithoutLockSkips(t *testing.T) {
	for _, service := range []string{"integration-app", "queue-app"} {
		m := New(nil, service)

		executed := false
		m.Run(context.Background(), JobSpec{Name: "job", Interval: time.Minute}, func() error {
			executed = true
			return nil
		})

		utils.AssertEqual(t, false, executed)
	}
}

// An unstoppable renewal goroutine leaks one per run.
func TestRenewLockStops(t *testing.T) {
	m := New(nil, "integration-app")

	returned := make(chan struct{})
	go func() {
		defer close(returned)
		m.renewLock(context.Background(), "cron:lock:job", JobSpec{Name: "job", Interval: time.Minute})()
	}()

	select {
	case <-returned:
	case <-time.After(2 * time.Second):
		t.Fatal("renewLock stop function did not return")
	}
}

// No TTL means renewal is disabled rather than ticking forever.
func TestRenewLockDisabledWithoutTTL(t *testing.T) {
	m := New(nil, "integration-app")
	m.renewLock(context.Background(), "cron:lock:job", JobSpec{Name: "job"})()
}

func TestRunProtectedTurnsPanicIntoError(t *testing.T) {
	// Must not propagate, or the scheduler goes down with the job.
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
