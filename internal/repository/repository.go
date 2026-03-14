package repository

import (
	"time"

	"github.com/Murchoid/iwashere/internal/domain/models"
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
	GetNotesBySession(id string) ([]*models.PrivateNote, error)

	//Reminders
	SaveReminder(reminder *models.Reminder) error
	GetReminder(id string) (*models.Reminder, error)
	ListReminders() ([]*models.Reminder, error)
	ListDueReminders() ([]*models.Reminder, error)
	DeleteReminder(id string) error
	DeactivateOrUpdateReminder(id string) error

	Close() error
}

type NoteFilter struct {
	ProjectPath string
	Branch      string
	Tags        []string
	Limit       int
	Since       time.Time
	SessionID   string
}
