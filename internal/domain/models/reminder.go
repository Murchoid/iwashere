package models

import "time"

type Reminder struct {
	ID        string
	NoteID    string
	Message   string
	DueAt     time.Time
	CreatedAt time.Time
	Repeats   string
	Active    bool
}
