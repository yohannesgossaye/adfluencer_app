// Package notify delivers emails and SMS in the background so a slow
// provider never holds up the HTTP request that triggered the message.
package notify

import (
	"context"
	"sync"
	"time"

	"adinfluence_service_app/platform/utils"

	"go.uber.org/zap"
)

const (
	defaultWorkers   = 4
	defaultQueueSize = 256
	maxAttempts      = 3
)

// Job is one message to deliver. Send must enforce its own timeout.
type Job struct {
	Kind   string // "email" or "sms", for logs
	Fields []zap.Field
	Send   func() error
}

type Dispatcher struct {
	jobs    chan Job
	wg      sync.WaitGroup
	log     utils.Logger
	backoff time.Duration

	mu     sync.RWMutex
	closed bool
}

func New(log utils.Logger) *Dispatcher {
	return newDispatcher(log, defaultWorkers, defaultQueueSize, 2*time.Second)
}

func newDispatcher(log utils.Logger, workers, queue int, backoff time.Duration) *Dispatcher {
	d := &Dispatcher{jobs: make(chan Job, queue), log: log, backoff: backoff}
	for range workers {
		d.wg.Add(1)
		go d.work()
	}
	return d
}

// Submit queues a job without blocking. It returns false when the queue is
// full or the dispatcher is shutting down; the job is then dropped and logged.
func (d *Dispatcher) Submit(job Job) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.closed {
		d.log.Warn(context.Background(), job.Kind+" dropped: dispatcher closed", job.Fields...)
		return false
	}
	select {
	case d.jobs <- job:
		return true
	default:
		d.log.Error(context.Background(), job.Kind+" dropped: queue full", job.Fields...)
		return false
	}
}

// Close stops accepting jobs and waits for queued ones until ctx expires.
func (d *Dispatcher) Close(ctx context.Context) {
	d.mu.Lock()
	if !d.closed {
		d.closed = true
		close(d.jobs)
	}
	d.mu.Unlock()

	done := make(chan struct{})
	go func() { d.wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-ctx.Done():
		d.log.Warn(context.Background(), "notify: shutdown before all messages were sent", zap.Int("pending", len(d.jobs)))
	}
}

func (d *Dispatcher) work() {
	defer d.wg.Done()
	for job := range d.jobs {
		d.deliver(job)
	}
}

func (d *Dispatcher) deliver(job Job) {
	ctx := context.Background()
	var err error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err = job.Send(); err == nil {
			d.log.Info(ctx, job.Kind+" sent", job.Fields...)
			return
		}
		if attempt < maxAttempts {
			time.Sleep(d.backoff * time.Duration(attempt))
		}
	}
	d.log.Warn(ctx, job.Kind+" failed after retries", append(job.Fields, zap.Int("attempts", maxAttempts), zap.String("reason", err.Error()))...)
}
