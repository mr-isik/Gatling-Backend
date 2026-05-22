package engine

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mr-isik/gatling-backend/internal/domain"
)

type WorkerPool struct {
	executor  *HTTPExecutor
	metricCh  chan<- domain.Metric
	steps     []domain.Step
	runID     string
	activeVUs atomic.Int32
	mu        sync.Mutex
	cancels   []context.CancelFunc
	wg        sync.WaitGroup
}

func NewWorkerPool(
	executor *HTTPExecutor,
	metricCh chan<- domain.Metric,
	steps []domain.Step,
	runID string,
) *WorkerPool {
	return &WorkerPool{
		executor: executor,
		metricCh: metricCh,
		steps:    steps,
		runID:    runID,
		cancels:  make([]context.CancelFunc, 0, 1000), // Preallocate capacity
	}
}

func (p *WorkerPool) SetTargetVUs(target int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	current := int(p.activeVUs.Load())

	if target > current {
		// Spawn new VUs
		toSpawn := target - current
		for i := 0; i < toSpawn; i++ {
			ctx, cancel := context.WithCancel(context.Background())
			p.cancels = append(p.cancels, cancel)
			p.activeVUs.Add(1)
			p.wg.Add(1)
			go p.runVU(ctx)
		}
	} else if target < current {
		// Retire VUs
		toRetire := current - target
		for i := 0; i < toRetire; i++ {
			// Cancel the last one in the list
			idx := len(p.cancels) - 1
			if idx >= 0 {
				p.cancels[idx]()
				p.cancels = p.cancels[:idx]
			}
		}
	}
}

func (p *WorkerPool) runVU(ctx context.Context) {
	defer p.wg.Done()
	defer p.activeVUs.Add(-1)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			for _, step := range p.steps {
				select {
				case <-ctx.Done():
					return
				default:
					metric := p.executor.Execute(ctx, p.runID, step)
					p.metricCh <- metric

					if step.ThinkTime > 0 {
						// Wait for think time, but be responsive to context cancellation
						select {
						case <-ctx.Done():
							return
						case <-time.After(time.Duration(step.ThinkTime) * time.Millisecond):
						}
					}
				}
			}
		}
	}
}

func (p *WorkerPool) Stop() {
	p.mu.Lock()
	for _, cancel := range p.cancels {
		cancel()
	}
	p.cancels = p.cancels[:0]
	p.mu.Unlock()

	// Wait for all goroutines to finish
	p.wg.Wait()
}

func (p *WorkerPool) ActiveVUs() int {
	return int(p.activeVUs.Load())
}
