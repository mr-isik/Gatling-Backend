package domain

import "time"

type Report struct {
	ID                string    `json:"id" db:"id"`
	RunID             string    `json:"run_id" db:"run_id"`
	Summary           Summary   `json:"summary" db:"summary"`
	AISummary         string    `json:"ai_summary,omitempty" db:"ai_summary"`
	AIRecommendations []string  `json:"ai_recommendations,omitempty"`
	Anomalies         []Anomaly `json:"anomalies,omitempty"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

type Summary struct {
	TotalRequests int     `json:"total_requests"`
	SuccessCount  int     `json:"success_count"`
	ErrorCount    int     `json:"error_count"`
	AvgLatency    float64 `json:"avg_latency"`
	P95Latency    float64 `json:"p95_latency"`
	P99Latency    float64 `json:"p99_latency"`
	MaxVUs        int     `json:"max_vus"`
	Duration      int64   `json:"duration_ms"`
}

type Anomaly struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // latency_spike, error_surge, throughput_drop
	Message   string    `json:"message"`
	Severity  string    `json:"severity"` // low, medium, high
}
