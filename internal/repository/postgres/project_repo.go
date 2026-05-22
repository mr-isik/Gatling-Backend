package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mr-isik/gatling-backend/internal/domain"
	"github.com/mr-isik/gatling-backend/internal/repository"
)

type projectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) repository.ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(ctx context.Context, p *domain.Project) (*domain.Project, error) {
	query := `INSERT INTO projects (name, description, owner_id) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := r.db.QueryRow(ctx, query, p.Name, p.Description, p.OwnerID).Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *projectRepository) GetByID(ctx context.Context, id string) (*domain.Project, error) {
	query := `SELECT id, name, description, owner_id, created_at FROM projects WHERE id = $1`
	var p domain.Project
	err := r.db.QueryRow(ctx, query, id).Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *projectRepository) List(ctx context.Context, ownerID string) ([]*domain.Project, error) {
	query := `SELECT id, name, description, owner_id, created_at FROM projects WHERE owner_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*domain.Project
	for rows.Next() {
		var p domain.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, &p)
	}
	return projects, nil
}

func (r *projectRepository) AddMember(ctx context.Context, projectID, userID, role string) error {
	query := `INSERT INTO project_members (project_id, user_id, role) VALUES ($1, $2, $3) ON CONFLICT (project_id, user_id) DO UPDATE SET role = $3`
	_, err := r.db.Exec(ctx, query, projectID, userID, role)
	return err
}

func (r *projectRepository) GetMembers(ctx context.Context, projectID string) ([]*domain.ProjectMember, error) {
	query := `SELECT project_id, user_id, role FROM project_members WHERE project_id = $1`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*domain.ProjectMember
	for rows.Next() {
		var m domain.ProjectMember
		if err := rows.Scan(&m.ProjectID, &m.UserID, &m.Role); err != nil {
			return nil, err
		}
		members = append(members, &m)
	}
	return members, nil
}
