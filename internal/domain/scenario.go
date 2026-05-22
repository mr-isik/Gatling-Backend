package domain

import "time"

type Scenario struct {
	ID          string    `json:"id" db:"id"`
	ProjectID   string    `json:"project_id" db:"project_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Tags        []string  `json:"tags" db:"tags"`
	Steps       []Step    `json:"steps" db:"steps"`
	IsAIGen     bool      `json:"is_ai_generated" db:"is_ai_generated"`
	CreatedBy   string    `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Step struct {
	Order     int               `json:"order"`
	Name      string            `json:"name"`
	Method    string            `json:"method"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers,omitempty"`
	Body      string            `json:"body,omitempty"`
	ThinkTime int               `json:"think_time_ms,omitempty"`
	Checks    []Check           `json:"checks,omitempty"`
}

type Check struct {
	Type     string `json:"type"`     // status, body_contains, json_path
	Operator string `json:"operator"` // eq, gt, contains
	Value    string `json:"value"`
}
