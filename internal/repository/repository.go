package repository

import (
	"context"
	"time"

	"github.com/mr-isik/gatling-backend/internal/domain"
)

type ScenarioRepository interface {
	Create(ctx context.Context, s *domain.Scenario) (*domain.Scenario, error)
	GetByID(ctx context.Context, id string) (*domain.Scenario, error)
	List(ctx context.Context, projectID string) ([]*domain.Scenario, error)
	Update(ctx context.Context, s *domain.Scenario) (*domain.Scenario, error)
	Delete(ctx context.Context, id string) error
}

type TestRunRepository interface {
	Create(ctx context.Context, r *domain.TestRun) (*domain.TestRun, error)
	GetByID(ctx context.Context, id string) (*domain.TestRun, error)
	List(ctx context.Context, projectID string) ([]*domain.TestRun, error)
	UpdateStatus(ctx context.Context, id string, status domain.RunStatus) error
	UpdateFinished(ctx context.Context, id string, status domain.RunStatus, finishedAt time.Time, errMsg string) error
}

type ReportRepository interface {
	Create(ctx context.Context, r *domain.Report) (*domain.Report, error)
	GetByRunID(ctx context.Context, runID string) (*domain.Report, error)
	Update(ctx context.Context, r *domain.Report) error
}

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	CreateAPIKey(ctx context.Context, k *domain.APIKey) (*domain.APIKey, error)
	GetAPIKeyByHash(ctx context.Context, hash string) (*domain.APIKey, error)
	DeleteAPIKey(ctx context.Context, id string) error
}

type ProjectRepository interface {
	Create(ctx context.Context, p *domain.Project) (*domain.Project, error)
	GetByID(ctx context.Context, id string) (*domain.Project, error)
	List(ctx context.Context, ownerID string) ([]*domain.Project, error)
	AddMember(ctx context.Context, projectID, userID, role string) error
	GetMembers(ctx context.Context, projectID string) ([]*domain.ProjectMember, error)
}

type MetricRepository interface {
	Write(ctx context.Context, metrics []domain.Metric) error
	QueryTimeSeries(ctx context.Context, runID string, from, to time.Time) ([]domain.AggregatedMetric, error)
	QueryLatest(ctx context.Context, runID string) (*domain.AggregatedMetric, error)
}

type CacheRepository interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

type RunStateRepository interface {
	SetState(ctx context.Context, runID string, state domain.RunStatus, ttl time.Duration) error
	GetState(ctx context.Context, runID string) (domain.RunStatus, error)
	DeleteState(ctx context.Context, runID string) error
}
