package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mr-isik/gatling-backend/internal/domain"
	"github.com/mr-isik/gatling-backend/internal/repository"
)

type testRunRepository struct {
	db *pgxpool.Pool
}

func NewTestRunRepository(db *pgxpool.Pool) repository.TestRunRepository {
	return &testRunRepository{db: db}
}

func (r *testRunRepository) Create(ctx context.Context, tr *domain.TestRun) (*domain.TestRun, error) {
	query := `INSERT INTO testruns (scenario_id, project_id, status, config, started_by, started_at)
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`
	err := r.db.QueryRow(ctx, query, tr.ScenarioID, tr.ProjectID, tr.Status, tr.Config, tr.StartedBy, tr.StartedAt).
		Scan(&tr.ID, &tr.CreatedAt)
	if err != nil {
		return nil, err
	}
	return tr, nil
}

func (r *testRunRepository) GetByID(ctx context.Context, id string) (*domain.TestRun, error) {
	query := `SELECT id, scenario_id, project_id, status, config, started_by, started_at, finished_at, COALESCE(error_msg, ''), created_at 
	          FROM testruns WHERE id = $1`
	var tr domain.TestRun
	err := r.db.QueryRow(ctx, query, id).Scan(&tr.ID, &tr.ScenarioID, &tr.ProjectID, &tr.Status, &tr.Config, &tr.StartedBy, &tr.StartedAt, &tr.FinishedAt, &tr.ErrorMsg, &tr.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &tr, nil
}

func (r *testRunRepository) List(ctx context.Context, projectID string) ([]*domain.TestRun, error) {
	query := `SELECT id, scenario_id, project_id, status, config, started_by, started_at, finished_at, COALESCE(error_msg, ''), created_at 
	          FROM testruns WHERE project_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var runs []*domain.TestRun
	for rows.Next() {
		var tr domain.TestRun
		if err := rows.Scan(&tr.ID, &tr.ScenarioID, &tr.ProjectID, &tr.Status, &tr.Config, &tr.StartedBy, &tr.StartedAt, &tr.FinishedAt, &tr.ErrorMsg, &tr.CreatedAt); err != nil {
			return nil, err
		}
		runs = append(runs, &tr)
	}
	return runs, nil
}

func (r *testRunRepository) UpdateStatus(ctx context.Context, id string, status domain.RunStatus) error {
	_, err := r.db.Exec(ctx, "UPDATE testruns SET status = $1 WHERE id = $2", status, id)
	return err
}

func (r *testRunRepository) UpdateFinished(ctx context.Context, id string, status domain.RunStatus, finishedAt time.Time, errMsg string) error {
	_, err := r.db.Exec(ctx, "UPDATE testruns SET status = $1, finished_at = $2, error_msg = $3 WHERE id = $4", status, finishedAt, errMsg, id)
	return err
}
