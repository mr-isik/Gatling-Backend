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

type LLMRequest struct {
	Model     string       `json:"model"`
	Messages  []LLMMessage `json:"messages"`
	MaxTokens int          `json:"max_tokens,omitempty"`
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

func (c *LLMClient) SendPrompt(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	messages := make([]LLMMessage, 0, 2)
	if systemPrompt != "" {
		messages = append(messages, LLMMessage{
			Role:    "system",
			Content: systemPrompt,
		})
	}
	messages = append(messages, LLMMessage{
		Role:    "user",
		Content: userPrompt,
	})

	reqBody := LLMRequest{
		Model:     "gpt-4o-mini", // "gpt-4o-mini" for speed/efficiency
		Messages:  messages,
		MaxTokens: 4096,
	}

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
