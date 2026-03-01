package models

import "time"

type Session struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Notes     []string  `json:"note_ids"` // References to notes
	Summary   string    `json:"summary"`  // Auto-generated or manual
	Branch    string    `json:"branch"`
}
