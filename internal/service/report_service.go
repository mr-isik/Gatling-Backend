package service

import (
	"context"

	"github.com/mr-isik/gatling-backend/internal/domain"
	"github.com/mr-isik/gatling-backend/internal/repository"
)

type ReportService struct {
	reportRepo repository.ReportRepository
	metricRepo repository.MetricRepository
	aiService  *AIService
}

func NewReportService(
	reportRepo repository.ReportRepository,
	metricRepo repository.MetricRepository,
	ai *AIService,
) *ReportService {
	return &ReportService{
		reportRepo: reportRepo,
		metricRepo: metricRepo,
		aiService:  ai,
	}
}

func (s *ReportService) Get(ctx context.Context, runID string) (*domain.Report, error) {
	return s.reportRepo.GetByRunID(ctx, runID)
}

func (s *ReportService) CreateFromRun(ctx context.Context, runID string) (*domain.Report, error) {
	// Query metrics from Influx and build summary
	// For now, this is a stub for the calculation
	summary := domain.Summary{}

	report := &domain.Report{
		RunID:   runID,
		Summary: summary,
	}

	return s.reportRepo.Create(ctx, report)
}

func (s *ReportService) AISummary(ctx context.Context, runID string) (*domain.Report, error) {
	report, err := s.reportRepo.GetByRunID(ctx, runID)
	if err != nil {
		return nil, err
	}

	aiSummary, err := s.aiService.GenerateSummary(ctx, report.Summary)
	if err != nil {
		return nil, err
	}

	report.AISummary = aiSummary
	err = s.reportRepo.Update(ctx, report)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (s *ReportService) Export(ctx context.Context, runID, format string) ([]byte, error) {
	// JSON/PDF Export logic stub
	return nil, nil
}

func (s *ReportService) Compare(ctx context.Context, runID1, runID2 string) (interface{}, error) {
	// Compare logic stub
	return nil, nil
}
