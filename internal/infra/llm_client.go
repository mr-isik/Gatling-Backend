package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LLMClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

func NewLLMClient(apiKey string) *LLMClient {
	return &LLMClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		baseURL: "https://api.openai.com/v1/chat/completions",
	}
}

type LLMResponseFormat struct {
	Type string `json:"type"` // "text" or "json_object"
}

type LLMRequest struct {
	Model          string            `json:"model"`
	Messages       []LLMMessage      `json:"messages"`
	MaxTokens      int               `json:"max_tokens,omitempty"`
	ResponseFormat *LLMResponseFormat `json:"response_format,omitempty"`
}

type LLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LLMResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// sendRequest is the internal helper that calls the OpenAI API.
func (c *LLMClient) sendRequest(ctx context.Context, reqBody LLMRequest) (string, error) {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("LLM API returned status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	var llmResp LLMResponse
	if err := json.NewDecoder(resp.Body).Decode(&llmResp); err != nil {
		return "", err
	}

	if len(llmResp.Choices) == 0 {
		return "", fmt.Errorf("empty response from LLM")
	}

	return llmResp.Choices[0].Message.Content, nil
}

// SendPrompt sends a free-text prompt and returns the model's raw text response.
func (c *LLMClient) SendPrompt(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	messages := make([]LLMMessage, 0, 2)
	if systemPrompt != "" {
		messages = append(messages, LLMMessage{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, LLMMessage{Role: "user", Content: userPrompt})

	return c.sendRequest(ctx, LLMRequest{
		Model:     "gpt-4o-mini",
		Messages:  messages,
		MaxTokens: 4096,
	})
}

// SendJSONPrompt forces the model to return a valid JSON object (no prose, no markdown).
// Uses the OpenAI `response_format: { type: "json_object" }` feature.
// The system prompt MUST mention "JSON" or the API will return an error.
func (c *LLMClient) SendJSONPrompt(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	messages := make([]LLMMessage, 0, 2)
	if systemPrompt != "" {
		messages = append(messages, LLMMessage{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, LLMMessage{Role: "user", Content: userPrompt})

	return c.sendRequest(ctx, LLMRequest{
		Model:     "gpt-4o-mini",
		Messages:  messages,
		MaxTokens: 4096,
		ResponseFormat: &LLMResponseFormat{Type: "json_object"},
	})
}
