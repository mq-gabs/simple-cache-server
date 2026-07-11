package worker

import (
	"context"
	"sync"
)

type Task func(context.Context)

type Pool struct {
	ctx    context.Context
	cancel context.CancelFunc

	queue chan Task
	wg    sync.WaitGroup
}

func New(ctx context.Context, workers, queueSize int) *Pool {
	if workers <= 0 {
		panic("workers must be > 0")
	}
	if queueSize < 0 {
		panic("queueSize must be >= 0")
	}

	ctx, cancel := context.WithCancel(ctx)

	p := &Pool{
		ctx:    ctx,
		cancel: cancel,
		queue:  make(chan Task, queueSize),
	}

	for range workers {
		p.wg.Add(1)

		go func() {
			defer p.wg.Done()

			for {
				select {
				case <-p.ctx.Done():
					return

				case task := <-p.queue:
					if task == nil {
						return
					}

					func() {
						defer func() {
							// Prevent a panic from terminating the worker.
							recover()
						}()

						task(p.ctx)
					}()
				}
			}
		}()
	}

	return p
}

// Submit returns false if:
//   - the context has been canceled
//   - the queue is full
func (p *Pool) Submit(task Task) bool {
	if p.ctx.Err() != nil {
		return false
	}

	select {
	case p.queue <- task:
		return true

	default:
		return false
	}
}

// Stop cancels the pool and waits for workers to exit.
func (p *Pool) Stop() {
	p.cancel()
	p.wg.Wait()
}
