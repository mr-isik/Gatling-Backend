package domain

import "time"

type Project struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	OwnerID     string    `json:"owner_id" db:"owner_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type ProjectMember struct {
	ProjectID string `json:"project_id" db:"project_id"`
	UserID    string `json:"user_id" db:"user_id"`
	Role      string `json:"role" db:"role"` // owner, admin, member
}
