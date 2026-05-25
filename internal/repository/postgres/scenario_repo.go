package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mr-isik/gatling-backend/internal/domain"
	"github.com/mr-isik/gatling-backend/internal/repository"
)

type scenarioRepository struct {
	db *pgxpool.Pool
}

func NewScenarioRepository(db *pgxpool.Pool) repository.ScenarioRepository {
	return &scenarioRepository{db: db}
}

func (r *scenarioRepository) Create(ctx context.Context, s *domain.Scenario) (*domain.Scenario, error) {
	query := `INSERT INTO scenarios (project_id, name, description, tags, steps, is_ai_generated, created_by)
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, s.ProjectID, s.Name, s.Description, s.Tags, s.Steps, s.IsAIGen, s.CreatedBy).
		Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *scenarioRepository) GetByID(ctx context.Context, id string) (*domain.Scenario, error) {
	query := `SELECT id, project_id, name, description, tags, steps, is_ai_generated, created_by, created_at, updated_at 
	          FROM scenarios WHERE id = $1`
	var s domain.Scenario
	err := r.db.QueryRow(ctx, query, id).Scan(&s.ID, &s.ProjectID, &s.Name, &s.Description, &s.Tags, &s.Steps, &s.IsAIGen, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *scenarioRepository) List(ctx context.Context, projectID string) ([]*domain.Scenario, error) {
	query := `SELECT id, project_id, name, description, tags, steps, is_ai_generated, created_by, created_at, updated_at 
	          FROM scenarios WHERE project_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scenarios []*domain.Scenario
	for rows.Next() {
		var s domain.Scenario
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.Name, &s.Description, &s.Tags, &s.Steps, &s.IsAIGen, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		scenarios = append(scenarios, &s)
	}
	return scenarios, nil
}

func (r *scenarioRepository) Update(ctx context.Context, s *domain.Scenario) (*domain.Scenario, error) {
	query := `UPDATE scenarios SET name = $1, description = $2, tags = $3, steps = $4, updated_at = NOW() 
	          WHERE id = $5 RETURNING updated_at`
	err := r.db.QueryRow(ctx, query, s.Name, s.Description, s.Tags, s.Steps, s.ID).Scan(&s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *scenarioRepository) Delete(ctx context.Context, id string) error {
	// First, delete any testruns associated with this scenario to prevent foreign key violations.
	// Since testruns cascade to reports, this will safely clean up related data.
	_, err := r.db.Exec(ctx, "DELETE FROM testruns WHERE scenario_id = $1", id)
	if err != nil {
		return err
	}
	
	_, err = r.db.Exec(ctx, "DELETE FROM scenarios WHERE id = $1", id)
	return err
}
