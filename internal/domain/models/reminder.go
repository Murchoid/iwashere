package models

import "time"

var (
	Once    = "once"
	Daily   = "daily"
	Weekly  = "weekly"
	Monthly = "monthly"
	Yearly  = "yearly"
)

type Reminder struct {
	ID        string
	NoteID    string
	Message   string
	DueAt     time.Time
	CreatedAt time.Time
	Repeats   string
	Active    bool
}
