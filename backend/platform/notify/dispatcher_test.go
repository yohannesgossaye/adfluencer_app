package notify

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"adinfluence_service_app/platform/utils"

	"go.uber.org/zap"
)

func testLogger(t *testing.T) utils.Logger {
	t.Helper()
	return utils.New(zap.NewNop())
}

func TestDeliverRetriesThenSucceeds(t *testing.T) {
	d := newDispatcher(testLogger(t), 1, 4, time.Millisecond)
	var calls atomic.Int32
	d.Submit(Job{Kind: "email", Send: func() error {
		if calls.Add(1) < 3 {
			return errors.New("temporary")
		}
		return nil
	}})
	d.Close(context.Background())
	if got := calls.Load(); got != 3 {
		t.Fatalf("calls = %d, want 3", got)
	}
}

func TestDeliverGivesUpAfterMaxAttempts(t *testing.T) {
	d := newDispatcher(testLogger(t), 1, 4, time.Millisecond)
	var calls atomic.Int32
	d.Submit(Job{Kind: "sms", Send: func() error { calls.Add(1); return errors.New("down") }})
	d.Close(context.Background())
	if got := calls.Load(); got != maxAttempts {
		t.Fatalf("calls = %d, want %d", got, maxAttempts)
	}
}

func TestSubmitDoesNotBlockWhenFull(t *testing.T) {
	block := make(chan struct{})
	d := newDispatcher(testLogger(t), 1, 1, time.Millisecond)
	slow := Job{Kind: "email", Send: func() error { <-block; return nil }}
	d.Submit(slow) // taken by the worker
	time.Sleep(10 * time.Millisecond)
	d.Submit(slow) // fills the queue
	done := make(chan bool, 1)
	go func() { done <- d.Submit(slow) }()
	select {
	case ok := <-done:
		if ok {
			t.Fatal("Submit on a full queue returned true")
		}
	case <-time.After(time.Second):
		t.Fatal("Submit blocked on a full queue")
	}
	close(block)
	d.Close(context.Background())
}

func TestSubmitAfterCloseIsRejected(t *testing.T) {
	d := newDispatcher(testLogger(t), 1, 1, time.Millisecond)
	d.Close(context.Background())
	if d.Submit(Job{Kind: "email", Send: func() error { return nil }}) {
		t.Fatal("Submit after Close returned true")
	}
}
