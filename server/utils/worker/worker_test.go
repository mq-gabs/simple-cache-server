package worker

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSubmitExecutesTask(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool := New(ctx, 2, 10)
	defer pool.Stop()

	done := make(chan struct{})

	ok := pool.Submit(func(ctx context.Context) {
		close(done)
	})

	if !ok {
		t.Fatal("expected task to be accepted")
	}

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("task was not executed")
	}
}

func TestQueueFull(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool := New(ctx, 1, 1)
	defer pool.Stop()

	started := make(chan struct{})
	block := make(chan struct{})

	// Occupy the only worker.
	if !pool.Submit(func(ctx context.Context) {
		close(started)
		<-block
	}) {
		t.Fatal("expected first submit to succeed")
	}

	<-started

	// Fill the queue.
	if !pool.Submit(func(ctx context.Context) {}) {
		t.Fatal("expected second submit to succeed")
	}

	// Queue is full.
	if pool.Submit(func(ctx context.Context) {}) {
		t.Fatal("expected third submit to be rejected")
	}

	close(block)
}

func TestContextCancellationRejectsSubmit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	pool := New(ctx, 1, 1)

	cancel()

	time.Sleep(20 * time.Millisecond)

	if pool.Submit(func(ctx context.Context) {}) {
		t.Fatal("expected submit to fail after cancellation")
	}

	pool.Stop()
}

func TestStopWaitsForWorkers(t *testing.T) {
	ctx := context.Background()

	pool := New(ctx, 1, 10)

	done := make(chan struct{})

	pool.Submit(func(ctx context.Context) {
		time.Sleep(100 * time.Millisecond)
		close(done)
	})

	pool.Stop()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Stop returned before task finished")
	}
}

func TestWorkerRecoversFromPanic(t *testing.T) {
	ctx := context.Background()

	pool := New(ctx, 1, 10)
	defer pool.Stop()

	done := make(chan struct{})

	pool.Submit(func(ctx context.Context) {
		panic("boom")
	})

	pool.Submit(func(ctx context.Context) {
		close(done)
	})

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("worker died after panic")
	}
}

func TestConcurrentSubmit(t *testing.T) {
	ctx := context.Background()

	pool := New(ctx, 8, 1000)
	defer pool.Stop()

	var executed atomic.Int64
	var wg sync.WaitGroup

	for range 500 {
		wg.Add(1)

		go func() {
			defer wg.Done()

			pool.Submit(func(ctx context.Context) {
				executed.Add(1)
			})
		}()
	}

	deadline := time.Now().Add(time.Second)

	wg.Wait()

	for time.Now().Before(deadline) {
		if executed.Load() == 500 {
			return
		}
		time.Sleep(time.Millisecond)
	}

	t.Fatalf("expected 500 executed tasks, got %d", executed.Load())
}

func TestTaskReceivesCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	pool := New(ctx, 1, 10)
	defer pool.Stop()

	done := make(chan struct{})

	pool.Submit(func(ctx context.Context) {
		<-ctx.Done()
		close(done)
	})

	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("task did not observe canceled context")
	}
}
