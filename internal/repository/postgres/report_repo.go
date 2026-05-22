package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mr-isik/gatling-backend/internal/domain"
	"github.com/mr-isik/gatling-backend/internal/repository"
)

type reportRepository struct {
	db *pgxpool.Pool
}

func NewReportRepository(db *pgxpool.Pool) repository.ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) Create(ctx context.Context, rep *domain.Report) (*domain.Report, error) {
	query := `INSERT INTO reports (run_id, summary, ai_summary, ai_recommendations, anomalies) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	err := r.db.QueryRow(ctx, query, rep.RunID, rep.Summary, rep.AISummary, rep.AIRecommendations, rep.Anomalies).
		Scan(&rep.ID, &rep.CreatedAt)
	if err != nil {
		return nil, err
	}
	return rep, nil
}

func (r *reportRepository) GetByRunID(ctx context.Context, runID string) (*domain.Report, error) {
	query := `SELECT id, run_id, summary, ai_summary, ai_recommendations, anomalies, created_at 
	          FROM reports WHERE run_id = $1`
	var rep domain.Report
	err := r.db.QueryRow(ctx, query, runID).Scan(&rep.ID, &rep.RunID, &rep.Summary, &rep.AISummary, &rep.AIRecommendations, &rep.Anomalies, &rep.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &rep, nil
}

func (r *reportRepository) Update(ctx context.Context, rep *domain.Report) error {
	query := `UPDATE reports SET summary = $1, ai_summary = $2, ai_recommendations = $3, anomalies = $4 
	          WHERE id = $5`
	_, err := r.db.Exec(ctx, query, rep.Summary, rep.AISummary, rep.AIRecommendations, rep.Anomalies, rep.ID)
	return err
}
