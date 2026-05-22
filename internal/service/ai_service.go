package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mr-isik/gatling-backend/internal/domain"
	"github.com/mr-isik/gatling-backend/internal/infra"
)

type AIService struct {
	llm *infra.LLMClient
}

func NewAIService(llm *infra.LLMClient) *AIService {
	return &AIService{llm: llm}
}

func (s *AIService) GenerateScenario(ctx context.Context, prompt string) (*domain.Scenario, error) {
	systemPrompt := `Sen bir load test uzmanısın. Kullanıcının tanımına göre HTTP test senaryosu JSON olarak üret. Yanıtın SADECE JSON formatında olmalı.
Beklenen JSON yapısı:
{
  "name": "Senaryo Adı",
  "description": "Senaryo açıklaması",
  "steps": [
    {
      "order": 1,
      "name": "Step 1",
      "method": "GET",
      "url": "https://api.example.com/v1/users",
      "headers": {"Content-Type": "application/json"},
      "body": "",
      "think_time_ms": 1000
    }
  ]
}`

	respText, err := s.llm.SendPrompt(ctx, systemPrompt, prompt)
	if err != nil {
		return nil, err
	}

	// Clean up markdown block if any
	respText = strings.TrimPrefix(respText, "```json\n")
	respText = strings.TrimPrefix(respText, "```\n")
	respText = strings.TrimSuffix(respText, "\n```")

	var scenario domain.Scenario
	if err := json.Unmarshal([]byte(respText), &scenario); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w (Response was: %s)", err, respText)
	}

	scenario.IsAIGen = true
	return &scenario, nil
}

func (s *AIService) AnalyzeMetrics(ctx context.Context, metrics []domain.AggregatedMetric) ([]domain.Anomaly, error) {
	systemPrompt := `Sen bir performans analisti. Verilen metriklerde anomali, spike ve sorunları tespit et. SADECE JSON formatında Array döndür.
Format:
[
  {
    "type": "latency_spike",
    "message": "P95 latency exceeded 500ms",
    "severity": "high"
  }
]`

	metricsJSON, _ := json.Marshal(metrics)
	userPrompt := string(metricsJSON)

	respText, err := s.llm.SendPrompt(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, err
	}

	respText = strings.TrimPrefix(respText, "```json\n")
	respText = strings.TrimPrefix(respText, "```\n")
	respText = strings.TrimSuffix(respText, "\n```")

	var anomalies []domain.Anomaly
	if err := json.Unmarshal([]byte(respText), &anomalies); err != nil {
		return nil, err
	}
	return anomalies, nil
}

func (s *AIService) GenerateSummary(ctx context.Context, summary domain.Summary) (string, error) {
	systemPrompt := `Sen bir teknik rapor yazarısın. Verilen test sonuçlarını Türkçe, yönetici seviyesinde, net ve profesyonel bir dille özetle. (1-2 paragraf)`

	summaryJSON, _ := json.Marshal(summary)
	userPrompt := string(summaryJSON)

	return s.llm.SendPrompt(ctx, systemPrompt, userPrompt)
}
