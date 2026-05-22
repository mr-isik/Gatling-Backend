package engine

import (
	"context"
	"sync"
	"time"

	"github.com/mr-isik/gatling-backend/internal/domain"
	"github.com/mr-isik/gatling-backend/internal/repository"
)

type TestRunServiceIface interface {
	UpdateStatus(ctx context.Context, runID string, status domain.RunStatus) error
	MarkFinished(ctx context.Context, runID string, status domain.RunStatus, errMsg string) error
}

type AIServiceIface interface {
	AnalyzeMetrics(ctx context.Context, metrics []domain.AggregatedMetric) ([]domain.Anomaly, error)
}

// HubIface allows mocking or decoupled usage of WS Hub
type HubIface interface {
	Broadcast(runID, topic string, payload interface{})
}

type Orchestrator struct {
	runService   TestRunServiceIface
	metricRepo   repository.MetricRepository
	scenarioRepo repository.ScenarioRepository
	aiService    AIServiceIface
	hub          HubIface
	mu           sync.Mutex
	cancels      map[string]context.CancelFunc
}

func NewOrchestrator(
	runService TestRunServiceIface,
	metricRepo repository.MetricRepository,
	scenarioRepo repository.ScenarioRepository,
	aiService AIServiceIface,
	hub HubIface,
) *Orchestrator {
	return &Orchestrator{
		runService:   runService,
		metricRepo:   metricRepo,
		scenarioRepo: scenarioRepo,
		aiService:    aiService,
		hub:          hub,
		cancels:      make(map[string]context.CancelFunc),
	}
}

func (o *Orchestrator) StartRun(ctx context.Context, runID string, config domain.RunConfig, steps []domain.Step) {
	// Create a new context for this run that we can cancel
	runCtx, cancel := context.WithCancel(context.Background())

	o.mu.Lock()
	o.cancels[runID] = cancel
	o.mu.Unlock()

	go func() {
		defer func() {
			o.mu.Lock()
			delete(o.cancels, runID)
			o.mu.Unlock()
		}()

		// 1. Status -> Running
		o.runService.UpdateStatus(runCtx, runID, domain.StatusRunning)

		// 2. Init components
		executor := NewHTTPExecutor()
		metricCh := make(chan domain.Metric, 10000)
		aggregator := NewMetricAggregator(o.metricRepo)
		scheduler := NewRampScheduler(config)
		pool := NewWorkerPool(executor, metricCh, steps, runID)

		// 3. Start aggregator
		aggCtx, aggCancel := context.WithCancel(runCtx)
		go aggregator.Run(aggCtx, metricCh)

		// 4. Main loop
		totalDuration := config.RampUpDuration + config.Duration + config.RampDownDuration
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		start := time.Now()
		tickCount := 0

		for {
			select {
			case <-runCtx.Done():
				// Stopped externally via StopRun
				pool.Stop()
				aggCancel()
				aggregator.Flush(context.Background())
				o.runService.MarkFinished(context.Background(), runID, domain.StatusStopped, "")
				return
			case <-ticker.C:
				tickCount++
				elapsed := time.Since(start)

				// 1. Live Metrik Push
				snapshot := aggregator.Snapshot(runID, pool.ActiveVUs())
				if o.hub != nil {
					o.hub.Broadcast(runID, "live", snapshot)
				}

				// 2. Anomali Kontrolü (Her 15sn'de bir)
				if tickCount%15 == 0 && o.aiService != nil {
					// Asynchronous call to prevent blocking the tick loop
					go func(snap domain.AggregatedMetric) {
						// Create an isolated context with timeout for AI call
						aiCtx, aiCancel := context.WithTimeout(context.Background(), 10*time.Second)
						defer aiCancel()

						anomalies, err := o.aiService.AnalyzeMetrics(aiCtx, []domain.AggregatedMetric{snap})
						if err == nil && len(anomalies) > 0 && o.hub != nil {
							o.hub.Broadcast(runID, "anomalies", anomalies)
						}
					}(snapshot)
				}

				if elapsed >= totalDuration {
					// Test completed naturally
					pool.Stop()
					aggCancel()
					aggregator.Flush(context.Background())
					o.runService.MarkFinished(context.Background(), runID, domain.StatusFinished, "")
					return
				}
				target := scheduler.CalculateVUs(elapsed)
				pool.SetTargetVUs(target)
			}
		}
	}()
}

func (o *Orchestrator) StopRun(runID string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if cancel, exists := o.cancels[runID]; exists {
		cancel()
	}
}
