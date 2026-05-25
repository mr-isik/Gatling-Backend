package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mr-isik/gatling-backend/internal/domain"
	"github.com/mr-isik/gatling-backend/internal/infra"
)

type AIService struct {
	llm *infra.LLMClient
}

func NewAIService(llm *infra.LLMClient) *AIService {
	return &AIService{llm: llm}
}

func (s *AIService) GenerateScenario(ctx context.Context, prompt string, targetURL string, apiCtx *domain.ApiContext) (*domain.Scenario, error) {
	systemPrompt := `You are an expert load testing engineer. Generate an HTTP load test scenario as a JSON object.
You MUST respond with ONLY a valid JSON object — no prose, no markdown, no explanation.
The JSON MUST follow this exact schema:
{
  "name": "string",
  "description": "string",
  "steps": [
    {
      "order": 1,
      "name": "string",
      "method": "GET|POST|PUT|DELETE|PATCH",
      "url": "https://...",
      "headers": {"key": "value"},
      "body": "string or empty",
      "think_time_ms": 1000
    }
  ]
}`

	if targetURL != "" {
		systemPrompt += fmt.Sprintf("\n\nIMPORTANT: The target application's base URL is: %s\nYou MUST use this exact base URL for ALL step URLs. Do NOT invent or substitute any other domain.", targetURL)
	}

	if apiCtx != nil && len(apiCtx.Endpoints) > 0 {
		// If no explicit targetURL but apiCtx has a base_url, remind the AI
		if targetURL == "" && apiCtx.BaseURL != "" {
			systemPrompt += fmt.Sprintf("\n\nThe API base URL from the documentation is: %s — use it for all step URLs.", apiCtx.BaseURL)
		}
		ctxBytes, _ := json.MarshalIndent(apiCtx, "", "  ")
		systemPrompt += "\n\nUse ONLY the following API documentation to build the scenario steps. Do NOT invent endpoints, paths, methods, or body schemas:\n" + string(ctxBytes)
	}

	respText, err := s.llm.SendJSONPrompt(ctx, systemPrompt, prompt)
	if err != nil {
		return nil, err
	}

	var scenario domain.Scenario
	if err := json.Unmarshal([]byte(respText), &scenario); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w (Response was: %s)", err, respText)
	}

	scenario.IsAIGen = true
	return &scenario, nil
}

func (s *AIService) AnalyzeMetrics(ctx context.Context, metrics []domain.AggregatedMetric) ([]domain.Anomaly, error) {
	systemPrompt := `You are a performance analyst. Detect anomalies, spikes, and issues in the given load test metrics.
You MUST respond with ONLY a valid JSON object — no prose, no markdown.
Use this exact schema:
{
  "anomalies": [
    {
      "type": "latency_spike|error_surge|throughput_drop",
      "message": "string",
      "severity": "low|medium|high"
    }
  ]
}
If no anomalies are found, return: {"anomalies": []}`

	metricsJSON, _ := json.Marshal(metrics)

	respText, err := s.llm.SendJSONPrompt(ctx, systemPrompt, string(metricsJSON))
	if err != nil {
		return nil, err
	}

	var result struct {
		Anomalies []domain.Anomaly `json:"anomalies"`
	}
	if err := json.Unmarshal([]byte(respText), &result); err != nil {
		return nil, err
	}
	return result.Anomalies, nil
}

func (s *AIService) GenerateSummary(ctx context.Context, summary domain.Summary) (string, error) {
	systemPrompt := `You are a technical report writer. Summarize the given test results in English, at an executive level, in a clear and professional language. (1-2 paragraphs)`

	summaryJSON, _ := json.Marshal(summary)
	userPrompt := string(summaryJSON)

	return s.llm.SendPrompt(ctx, systemPrompt, userPrompt)
}
