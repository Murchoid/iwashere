package models

import (
	"time"
)

type Note struct {
	ID          string    `json:"id"`
	Message     string    `json:"message"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ProjectPath string    `json:"project_path"`
	Branch      string    `json:"branch"`
	GitHash     string    `json:"git_hash"`
	Tags        []string  `json:"tags"`
}
