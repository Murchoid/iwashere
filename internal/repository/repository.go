package repository

import (
	"time"

	"githum.com/Murchoid/iwashere/internal/domain/models"
)

type Repository interface {
	// Notes
	SaveNote(note *models.Note) error
	GetNote(id string) (*models.Note, error)
	ListNotes(filter *NoteFilter) ([]*models.Note, error)
	UpdateNote(note *models.Note) error
	DeleteNote(id string) error

	// Sessions
	SaveSession(session *models.Session) error
	GetSession(id string) (*models.Session, error)

	Close() error
}

type NoteFilter struct {
	ProjectPath string
	Branch      string
	Tags        []string
	Limit       int
	Since       time.Time
}
