package domain

import "time"

type Metric struct {
	RunID      string    `json:"run_id"`
	StepName   string    `json:"step_name"`
	Timestamp  time.Time `json:"timestamp"`
	Latency    int64     `json:"latency_ms"`
	StatusCode int       `json:"status_code"`
	Success    bool      `json:"success"`
	BytesSent  int64     `json:"bytes_sent"`
	BytesRecv  int64     `json:"bytes_received"`
}

type AggregatedMetric struct {
	RunID      string    `json:"run_id"`
	Timestamp  time.Time `json:"timestamp"`
	P50        float64   `json:"p50"`
	P95        float64   `json:"p95"`
	P99        float64   `json:"p99"`
	AvgLatency float64   `json:"avg_latency"`
	Throughput float64   `json:"throughput"` // req/s
	ErrorRate  float64   `json:"error_rate"` // %
	ActiveVUs  int       `json:"active_vus"`
}
