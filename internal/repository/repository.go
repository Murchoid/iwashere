package repository

import (
	"time"

	"githum.com/Murchoid/iwashere/internal/domain/models"
)

type Repository interface {
	// Notes
	SaveNote(note *models.PrivateNote) error
	GetNote(id string) (*models.PrivateNote, error)
	ListNotes(filter *NoteFilter) ([]*models.PrivateNote, error)
	UpdateNote(note *models.PrivateNote) error
	DeleteNote(id string) error

	// Sessions
	SaveSession(session *models.Session) error
	ListSessions() ([]*models.Session, error)
	GetSession(id string) (*models.Session, error)
	GetOpenSession() (*models.Session, error)

	Close() error
}

type NoteFilter struct {
	ProjectPath string
	Branch      string
	Tags        []string
	Limit       int
	Since       time.Time
}
