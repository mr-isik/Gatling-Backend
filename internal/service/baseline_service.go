package service

import (
	"context"
	"fmt"

	"github.com/mr-isik/gatling-backend/internal/repository"
)

type BaselineComparison struct {
	P95Change        float64 `json:"p95_change_pct"`
	P99Change        float64 `json:"p99_change_pct"`
	ThroughputChange float64 `json:"throughput_change_pct"`
	ErrorRateChange  float64 `json:"error_rate_change_pct"`
	IsRegression     bool    `json:"is_regression"`
	Details          string  `json:"details"`
}

type BaselineService struct {
	runRepo    repository.TestRunRepository
	metricRepo repository.MetricRepository
}

func NewBaselineService(runRepo repository.TestRunRepository, metricRepo repository.MetricRepository) *BaselineService {
	return &BaselineService{
		runRepo:    runRepo,
		metricRepo: metricRepo,
	}
}

func (s *BaselineService) CompareWithBaseline(ctx context.Context, currentRunID, baselineRunID string) (*BaselineComparison, error) {
	// Stub calculation
	// 1. Fetch latest aggregated metrics or full summary for both runs
	// 2. Calculate percentages
	// 3. Determine IsRegression (e.g. P95 > 20% increase)

	// Since metric aggregation isn't fully implemented in metricRepo yet, we return a mock response.
	comp := &BaselineComparison{
		P95Change:        15.5,
		P99Change:        22.1,
		ThroughputChange: -5.0,
		ErrorRateChange:  1.2,
		IsRegression:     true,
		Details:          fmt.Sprintf("Compared current run %s with baseline %s", currentRunID, baselineRunID),
	}

	// Example regression logic
	if comp.P95Change > 20.0 || comp.ErrorRateChange > 5.0 {
		comp.IsRegression = true
	} else {
		comp.IsRegression = false
	}

	return comp, nil
}
