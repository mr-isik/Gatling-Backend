package domain

import "time"

type RunStatus string

const (
	StatusPending  RunStatus = "pending"
	StatusRunning  RunStatus = "running"
	StatusStopping RunStatus = "stopping"
	StatusFinished RunStatus = "finished"
	StatusFailed   RunStatus = "failed"
	StatusStopped  RunStatus = "stopped"
)

type TestRun struct {
	ID         string     `json:"id" db:"id"`
	ScenarioID string     `json:"scenario_id" db:"scenario_id"`
	ProjectID  string     `json:"project_id" db:"project_id"`
	Status     RunStatus  `json:"status" db:"status"`
	Config     RunConfig  `json:"config" db:"config"`
	StartedBy  string     `json:"started_by" db:"started_by"`
	StartedAt  *time.Time `json:"started_at,omitempty" db:"started_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty" db:"finished_at"`
	ErrorMsg   string     `json:"error_msg,omitempty" db:"error_msg"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}

type RunConfig struct {
	VUs              int           `json:"vus"`
	Duration         time.Duration `json:"duration"`
	RampUpDuration   time.Duration `json:"ramp_up_duration"`
	RampDownDuration time.Duration `json:"ramp_down_duration"`
	BaselineRunID    string        `json:"baseline_run_id,omitempty"`
}
