package models

import (
	"time"
)

type Note struct {
	ID          string    `json:"id"`
	SessionID   string    `json:"session_id"`
	Message     string    `json:"message"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ProjectPath string    `json:"project_path"`
	// Git-related fields
	Branch        string   `json:"branch,omitempty"`
	CommitHash    string   `json:"commit_hash,omitempty"`
	CommitMsg     string   `json:"commit_msg,omitempty"`
	Remote        string   `json:"remote,omitempty"`
	ModifiedFiles []string `json:"modified_files,omitempty"`

	Tags []string `json:"tags,omitempty"`
}
