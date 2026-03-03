package models

import (
	"time"
)

type PrivateNote struct {
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

// Team note (sanitized for collaboration)
type TeamNote struct {
    ID          string    `json:"id"`
    Message     string    `json:"message"`
    Author      string    `json:"author"`
    CreatedAt   time.Time `json:"created_at"`
    Tags        []string  `json:"tags"`
    SessionName string    `json:"session_name,omitempty"`
    // NO personal context (files, branches)
}

// Shared reference
type SharedReference struct {
    PrivateNoteID string    `json:"private_note_id"`
    SharedBy      string    `json:"shared_by"`
    SharedAt      time.Time `json:"shared_at"`
    NotePreview   string    `json:"note_preview"` // First 100 chars
}