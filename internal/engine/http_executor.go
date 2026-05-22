package engine

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/mr-isik/gatling-backend/internal/domain"
)

type HTTPExecutor struct {
	client *http.Client
}

func NewHTTPExecutor() *HTTPExecutor {
	// Create a shared HTTP client optimized for high concurrency
	transport := &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 1000,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
	}

	return &HTTPExecutor{
		client: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second, // Default timeout
		},
	}
}

func (e *HTTPExecutor) Execute(ctx context.Context, runID string, step domain.Step) domain.Metric {
	start := time.Now()
	metric := domain.Metric{
		RunID:     runID,
		StepName:  step.Name,
		Timestamp: start,
		Success:   false,
	}

	var reqBody io.Reader
	if step.Body != "" {
		reqBody = bytes.NewBufferString(step.Body)
		metric.BytesSent = int64(len(step.Body))
	}

	req, err := http.NewRequestWithContext(ctx, step.Method, step.URL, reqBody)
	if err != nil {
		metric.Latency = time.Since(start).Milliseconds()
		metric.StatusCode = 0 // Represents internal error
		return metric
	}

	for k, v := range step.Headers {
		req.Header.Set(k, v)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		metric.Latency = time.Since(start).Milliseconds()
		metric.StatusCode = 0 // Network error or timeout
		return metric
	}
	defer resp.Body.Close()

	// Read body to measure actual completion time and byte count
	bodyBytes, err := io.ReadAll(resp.Body)
	if err == nil {
		metric.BytesRecv = int64(len(bodyBytes))
	}

	metric.Latency = time.Since(start).Milliseconds()
	metric.StatusCode = resp.StatusCode

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		metric.Success = true
	}

	return metric
}
