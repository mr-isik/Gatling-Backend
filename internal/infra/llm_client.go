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
		baseURL: "https://api.anthropic.com/v1/messages",
	}
}

type LLMRequest struct {
	Model     string       `json:"model"`
	MaxTokens int          `json:"max_tokens"`
	Messages  []LLMMessage `json:"messages"`
	System    string       `json:"system,omitempty"`
}

type LLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LLMResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

func (c *LLMClient) SendPrompt(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	reqBody := LLMRequest{
		Model:     "claude-3-5-sonnet-20240620", // or claude-3-haiku for speed
		MaxTokens: 4096,
		System:    systemPrompt,
		Messages: []LLMMessage{
			{Role: "user", Content: userPrompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
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

	if len(llmResp.Content) == 0 {
		return "", fmt.Errorf("empty response from LLM")
	}

	return llmResp.Content[0].Text, nil
}
