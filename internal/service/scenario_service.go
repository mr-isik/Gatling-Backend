package service

import (
	"context"

	"github.com/mr-isik/gatling-backend/internal/domain"
	"github.com/mr-isik/gatling-backend/internal/repository"
)

type ScenarioService struct {
	scenarioRepo repository.ScenarioRepository
	aiService    *AIService
}

func NewScenarioService(repo repository.ScenarioRepository, ai *AIService) *ScenarioService {
	return &ScenarioService{
		scenarioRepo: repo,
		aiService:    ai,
	}
}

func (s *ScenarioService) Create(ctx context.Context, scenario *domain.Scenario) (*domain.Scenario, error) {
	if scenario.Name == "" || scenario.ProjectID == "" {
		return nil, domain.ErrBadRequest
	}
	return s.scenarioRepo.Create(ctx, scenario)
}

func (s *ScenarioService) GetByID(ctx context.Context, id string) (*domain.Scenario, error) {
	return s.scenarioRepo.GetByID(ctx, id)
}

func (s *ScenarioService) Update(ctx context.Context, scenario *domain.Scenario) (*domain.Scenario, error) {
	return s.scenarioRepo.Update(ctx, scenario)
}

func (s *ScenarioService) Delete(ctx context.Context, id string) error {
	return s.scenarioRepo.Delete(ctx, id)
}

func (s *ScenarioService) List(ctx context.Context, projectID string) ([]*domain.Scenario, error) {
	return s.scenarioRepo.List(ctx, projectID)
}

func (s *ScenarioService) Generate(ctx context.Context, prompt string) (*domain.Scenario, error) {
	// AI integration
	scenario, err := s.aiService.GenerateScenario(ctx, prompt)
	if err != nil {
		return nil, err
	}
	// We return it for the user to review, not saving immediately to DB unless specified.
	// For this phase, we'll just return the generated struct.
	return scenario, nil
}

func (s *ScenarioService) Clone(ctx context.Context, id, newName string) (*domain.Scenario, error) {
	original, err := s.scenarioRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	clone := &domain.Scenario{
		ProjectID:   original.ProjectID,
		Name:        newName,
		Description: original.Description,
		Tags:        append([]string(nil), original.Tags...),       // deep copy
		Steps:       append([]domain.Step(nil), original.Steps...), // deep copy
		IsAIGen:     original.IsAIGen,
		CreatedBy:   original.CreatedBy,
	}

	return s.scenarioRepo.Create(ctx, clone)
}
