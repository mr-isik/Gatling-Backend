package engine

import (
	"context"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/mr-isik/gatling-backend/internal/domain"
	"github.com/mr-isik/gatling-backend/internal/repository"
)

type MetricAggregator struct {
	metricRepo    repository.MetricRepository
	buffer        []domain.Metric
	mu            sync.Mutex
	flushInterval time.Duration
}

func NewMetricAggregator(repo repository.MetricRepository) *MetricAggregator {
	return &MetricAggregator{
		metricRepo:    repo,
		buffer:        make([]domain.Metric, 0, 10000), // preallocate
		flushInterval: 5 * time.Second,
	}
}

func (a *MetricAggregator) Run(ctx context.Context, metricCh <-chan domain.Metric) {
	ticker := time.NewTicker(a.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return // Context cancelled, exit loop. Orchestrator calls Flush.
		case m := <-metricCh:
			a.mu.Lock()
			a.buffer = append(a.buffer, m)
			a.mu.Unlock()
		case <-ticker.C:
			a.Flush(ctx)
		}
	}
}

func (a *MetricAggregator) Flush(ctx context.Context) error {
	a.mu.Lock()
	if len(a.buffer) == 0 {
		a.mu.Unlock()
		return nil
	}
	// Copy buffer to avoid holding lock during network call
	metricsToFlush := make([]domain.Metric, len(a.buffer))
	copy(metricsToFlush, a.buffer)
	// Reset buffer
	a.buffer = a.buffer[:0]
	a.mu.Unlock()

	return a.metricRepo.Write(ctx, metricsToFlush)
}

func (a *MetricAggregator) Snapshot(runID string, activeVUs int) domain.AggregatedMetric {
	a.mu.Lock()
	defer a.mu.Unlock()

	total := len(a.buffer)
	if total == 0 {
		return domain.AggregatedMetric{
			RunID:     runID,
			Timestamp: time.Now(),
			ActiveVUs: activeVUs,
		}
	}

	latencies := make([]float64, total)
	errorCount := 0
	var sumLatency float64

	for i, m := range a.buffer {
		l := float64(m.Latency)
		latencies[i] = l
		sumLatency += l
		if !m.Success {
			errorCount++
		}
	}

	sort.Float64s(latencies)

	p50Idx := int(math.Ceil(float64(total)*0.50)) - 1
	p95Idx := int(math.Ceil(float64(total)*0.95)) - 1
	p99Idx := int(math.Ceil(float64(total)*0.99)) - 1

	if p50Idx < 0 {
		p50Idx = 0
	}
	if p95Idx < 0 {
		p95Idx = 0
	}
	if p99Idx < 0 {
		p99Idx = 0
	}

	return domain.AggregatedMetric{
		RunID:      runID,
		Timestamp:  time.Now(),
		P50:        latencies[p50Idx],
		P95:        latencies[p95Idx],
		P99:        latencies[p99Idx],
		AvgLatency: sumLatency / float64(total),
		Throughput: float64(total) / a.flushInterval.Seconds(),
		ErrorRate:  (float64(errorCount) / float64(total)) * 100,
		ActiveVUs:  activeVUs,
	}
}
