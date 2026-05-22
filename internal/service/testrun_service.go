package service

import (
	"context"
	"time"

	"github.com/mr-isik/gatling-backend/internal/domain"
	"github.com/mr-isik/gatling-backend/internal/repository"
)

type TestRunService struct {
	runRepo      repository.TestRunRepository
	scenarioRepo repository.ScenarioRepository
	metricRepo   repository.MetricRepository
	stateRepo    repository.RunStateRepository
	cacheRepo    repository.CacheRepository
	orchestrator OrchestratorIface // interface to break dependency cycle
}

type OrchestratorIface interface {
	StartRun(ctx context.Context, runID string, config domain.RunConfig, steps []domain.Step)
	StopRun(runID string)
}

func NewTestRunService(
	runRepo repository.TestRunRepository,
	scenarioRepo repository.ScenarioRepository,
	metricRepo repository.MetricRepository,
	stateRepo repository.RunStateRepository,
	cacheRepo repository.CacheRepository,
) *TestRunService {
	return &TestRunService{
		runRepo:      runRepo,
		scenarioRepo: scenarioRepo,
		metricRepo:   metricRepo,
		stateRepo:    stateRepo,
		cacheRepo:    cacheRepo,
	}
}

func (s *TestRunService) SetOrchestrator(o OrchestratorIface) {
	s.orchestrator = o
}

func (s *TestRunService) Start(ctx context.Context, scenarioID, projectID, userID string, config domain.RunConfig) (string, error) {
	_, err := s.scenarioRepo.GetByID(ctx, scenarioID)
	if err != nil {
		return "", err
	}

	now := time.Now()
	run := &domain.TestRun{
		ScenarioID: scenarioID,
		ProjectID:  projectID,
		Status:     domain.StatusPending,
		Config:     config,
		StartedBy:  userID,
		StartedAt:  &now,
	}

	created, err := s.runRepo.Create(ctx, run)
	if err != nil {
		return "", err
	}

	err = s.stateRepo.SetState(ctx, created.ID, domain.StatusPending, 24*time.Hour)
	if err != nil {
		// Log error but continue
	}

	// Orchestrator trigger
	if s.orchestrator != nil {
		scenario, _ := s.scenarioRepo.GetByID(ctx, scenarioID)
		s.orchestrator.StartRun(context.Background(), created.ID, config, scenario.Steps)
	}

	return created.ID, nil
}

func (s *TestRunService) Stop(ctx context.Context, runID string) error {
	status, err := s.stateRepo.GetState(ctx, runID)
	if err != nil && err != domain.ErrNotFound {
		return err
	}

	if status != domain.StatusRunning && status != domain.StatusPending {
		return domain.ErrBadRequest
	}

	err = s.stateRepo.SetState(ctx, runID, domain.StatusStopping, 24*time.Hour)
	if err != nil {
		return err
	}

	err = s.runRepo.UpdateStatus(ctx, runID, domain.StatusStopping)
	if err == nil && s.orchestrator != nil {
		s.orchestrator.StopRun(runID)
	}
	return err
}

func (s *TestRunService) GetByID(ctx context.Context, runID string) (*domain.TestRun, error) {
	// Attempt cache first? Omitting cache logic for simplicity now.
	return s.runRepo.GetByID(ctx, runID)
}

func (s *TestRunService) List(ctx context.Context, projectID string) ([]*domain.TestRun, error) {
	return s.runRepo.List(ctx, projectID)
}

func (s *TestRunService) GetMetrics(ctx context.Context, runID string, from, to time.Time) ([]domain.AggregatedMetric, error) {
	return s.metricRepo.QueryTimeSeries(ctx, runID, from, to)
}

func (s *TestRunService) UpdateStatus(ctx context.Context, runID string, status domain.RunStatus) error {
	if err := s.stateRepo.SetState(ctx, runID, status, 24*time.Hour); err != nil {
		return err
	}
	return s.runRepo.UpdateStatus(ctx, runID, status)
}

func (s *TestRunService) MarkFinished(ctx context.Context, runID string, status domain.RunStatus, errMsg string) error {
	now := time.Now()
	if err := s.runRepo.UpdateFinished(ctx, runID, status, now, errMsg); err != nil {
		return err
	}
	return s.stateRepo.DeleteState(ctx, runID)
}
