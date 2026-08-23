package cronmon

import (
	"context"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

// Redis hash fields of a heartbeat. Constants because every service reads the
// hashes written by every other service.
const (
	fieldService      = "service"
	fieldJob          = "job"
	fieldIntervalSec  = "interval_sec"
	fieldDescription  = "description"
	fieldRegisteredAt = "registered_at"
	fieldInstance     = "instance"
	fieldLastStart    = "last_start"
	fieldLastSuccess  = "last_success"
	fieldLastFailure  = "last_failure"
	fieldLastError    = "last_error"
	fieldLastDuration = "last_duration_ms"
	fieldRunCount     = "run_count"
	fieldFailureCount = "failure_count"
	fieldConsecFails  = "consecutive_failures"
)

// Health classifies a job for alerting.
type Health string

const (
	// HealthOK means the last run succeeded and the next one is not overdue.
	HealthOK Health = "ok"
	// HealthFailing means the most recent run returned an error or panicked.
	HealthFailing Health = "failing"
	// HealthStale means nothing has run for far longer than the schedule allows.
	HealthStale Health = "stale"
	// HealthNeverRun means the job is registered but has not run yet.
	HealthNeverRun Health = "never_run"
)

// JobStatus is a single job's observed state.
type JobStatus struct {
	Service             string     `json:"service"`
	Job                 string     `json:"job"`
	Health              Health     `json:"health"`
	Description         string     `json:"description,omitempty"`
	IntervalSeconds     int64      `json:"interval_seconds,omitempty"`
	Instance            string     `json:"instance,omitempty"`
	LastStart           *time.Time `json:"last_start,omitempty"`
	LastSuccess         *time.Time `json:"last_success,omitempty"`
	LastFailure         *time.Time `json:"last_failure,omitempty"`
	LastError           string     `json:"last_error,omitempty"`
	LastDurationMs      int64      `json:"last_duration_ms,omitempty"`
	SecondsSinceSuccess *int64     `json:"seconds_since_success,omitempty"`
	RunCount            int64      `json:"run_count"`
	FailureCount        int64      `json:"failure_count"`
	ConsecutiveFailures int64      `json:"consecutive_failures"`
}

// markStarted increments run_count on start rather than on completion, so a
// job killed mid-run (OOM, pod eviction) still leaves a trace.
func (m *Monitor) markStarted(ctx context.Context, spec JobSpec, start time.Time) {
	if m.client == nil {
		return
	}

	key := m.heartbeatKey(spec.Name)
	pipe := m.client.TxPipeline()
	pipe.HSet(ctx, key, map[string]interface{}{
		fieldService:     m.service,
		fieldJob:         spec.Name,
		fieldIntervalSec: int64(spec.Interval.Seconds()),
		fieldDescription: spec.Description,
		fieldInstance:    m.instance,
		fieldLastStart:   formatTime(start.UTC()),
	})
	pipe.HIncrBy(ctx, key, fieldRunCount, 1)
	pipe.Expire(ctx, key, heartbeatTTL)
	if _, err := pipe.Exec(ctx); err != nil {
		log.Printf("[CronMon] failed writing start heartbeat for %s: %v", spec.Name, err)
	}
}

func (m *Monitor) markFinished(ctx context.Context, spec JobSpec, duration time.Duration, runErr error) {
	if m.client == nil {
		return
	}

	key := m.heartbeatKey(spec.Name)
	now := formatTime(time.Now().UTC())

	pipe := m.client.TxPipeline()
	pipe.HSet(ctx, key, fieldLastDuration, duration.Milliseconds())
	if runErr != nil {
		pipe.HSet(ctx, key, map[string]interface{}{
			fieldLastFailure: now,
			fieldLastError:   truncate(runErr.Error(), 500),
		})
		pipe.HIncrBy(ctx, key, fieldFailureCount, 1)
		pipe.HIncrBy(ctx, key, fieldConsecFails, 1)
	} else {
		pipe.HSet(ctx, key, map[string]interface{}{
			fieldLastSuccess: now,
			fieldConsecFails: 0,
		})
		pipe.HDel(ctx, key, fieldLastError)
	}
	pipe.Expire(ctx, key, heartbeatTTL)
	if _, err := pipe.Exec(ctx); err != nil {
		log.Printf("[CronMon] failed writing finish heartbeat for %s: %v", spec.Name, err)
	}
}

// Status returns every cron heartbeat in Redis. Pass a service name to narrow
// the scan to one service, or "" for all.
func Status(ctx context.Context, client *redis.Client, service string) ([]JobStatus, error) {
	if client == nil {
		return []JobStatus{}, nil
	}

	match := HeartbeatKeyPrefix + "*"
	if service != "" {
		match = HeartbeatKeyPrefix + service + ":*"
	}

	var (
		statuses []JobStatus
		cursor   uint64
	)
	for {
		keys, next, err := client.Scan(ctx, cursor, match, 100).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			fields, err := client.HGetAll(ctx, key).Result()
			if err != nil {
				log.Printf("[CronMon] failed reading heartbeat %s: %v", key, err)
				continue
			}
			if len(fields) == 0 {
				continue
			}
			statuses = append(statuses, buildStatus(key, fields))
		}

		cursor = next
		if cursor == 0 {
			break
		}
	}

	sort.Slice(statuses, func(i, j int) bool {
		if statuses[i].Service != statuses[j].Service {
			return statuses[i].Service < statuses[j].Service
		}
		return statuses[i].Job < statuses[j].Job
	})

	if statuses == nil {
		statuses = []JobStatus{}
	}
	return statuses, nil
}

func buildStatus(key string, fields map[string]string) JobStatus {
	status := JobStatus{
		Service:             fields[fieldService],
		Job:                 fields[fieldJob],
		Description:         fields[fieldDescription],
		Instance:            fields[fieldInstance],
		IntervalSeconds:     parseInt(fields[fieldIntervalSec]),
		LastStart:           parseTime(fields[fieldLastStart]),
		LastSuccess:         parseTime(fields[fieldLastSuccess]),
		LastFailure:         parseTime(fields[fieldLastFailure]),
		LastError:           fields[fieldLastError],
		LastDurationMs:      parseInt(fields[fieldLastDuration]),
		RunCount:            parseInt(fields[fieldRunCount]),
		FailureCount:        parseInt(fields[fieldFailureCount]),
		ConsecutiveFailures: parseInt(fields[fieldConsecFails]),
	}

	// Older heartbeats may predate the service/job fields; recover them from the key.
	if status.Service == "" || status.Job == "" {
		if trimmed := strings.TrimPrefix(key, HeartbeatKeyPrefix); trimmed != key {
			if service, job, ok := strings.Cut(trimmed, ":"); ok {
				if status.Service == "" {
					status.Service = service
				}
				if status.Job == "" {
					status.Job = job
				}
			}
		}
	}

	if status.LastSuccess != nil {
		since := int64(time.Since(*status.LastSuccess).Seconds())
		status.SecondsSinceSuccess = &since
	}
	status.Health = classify(status)
	return status
}

// classify decides a job's health. Failing outranks stale: a job that runs but
// errors is reported as failing rather than merely overdue.
func classify(status JobStatus) Health {
	if status.ConsecutiveFailures > 0 {
		return HealthFailing
	}
	if status.LastSuccess == nil {
		if status.LastStart == nil {
			return HealthNeverRun
		}
		return HealthFailing
	}
	// Two missed intervals, so a single slow or skipped tick is not an alert.
	if status.IntervalSeconds > 0 {
		threshold := time.Duration(status.IntervalSeconds) * time.Second * 2
		if time.Since(*status.LastSuccess) > threshold {
			return HealthStale
		}
	}
	return HealthOK
}

const timeFormat = time.RFC3339

func formatTime(t time.Time) string {
	return t.Format(timeFormat)
}

func parseTime(value string) *time.Time {
	if value == "" {
		return nil
	}
	parsed, err := time.Parse(timeFormat, value)
	if err != nil {
		return nil
	}
	return &parsed
}

func parseInt(value string) int64 {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return parsed
}

func truncate(value string, max int) string {
	if len(value) <= max {
		return value
	}
	return value[:max]
}
