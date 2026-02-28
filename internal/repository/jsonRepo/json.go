package jsonRepo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"githum.com/Murchoid/iwashere/internal/domain/errors"
	"githum.com/Murchoid/iwashere/internal/domain/models"
	"githum.com/Murchoid/iwashere/internal/repository"
	"githum.com/Murchoid/iwashere/internal/utils"
)

// Json implementation as a storage option
type JSONRepository struct {
	notesBasePath string // .iwashere/notes/
	sessionBasePath string // .iwashere/sessions/
}

func NewJSONRepository(iwasherePath string) *JSONRepository {
	notesPath := filepath.Join(iwasherePath, "notes")
	sessionPath := filepath.Join(iwasherePath, "sessions")
	// Ensure directory exists
	os.MkdirAll(notesPath, 0755)
	os.MkdirAll(sessionPath, 0755)
	return &JSONRepository{notesBasePath: notesPath, sessionBasePath: sessionPath}
}

func (r *JSONRepository) ListNotes(filter *repository.NoteFilter) ([]*models.Note, error) {
	// Read all note files
	files, err := os.ReadDir(r.notesBasePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*models.Note{}, nil // No notes yet
		}
		return nil, fmt.Errorf("failed to read notes directory: %w", err)
	}

	var notes []*models.Note

	for _, file := range files {
		// Skip directories and non-JSON files
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		// Read and parse each note
		note, err := r.readNoteFile(file.Name())
		if err != nil {
			// Log but don't stop - could be corrupted file
			fmt.Fprintf(os.Stderr, "Warning: failed to read note %s: %v\n", file.Name(), err)
			continue
		}

		// Apply filters
		if r.matchesFilter(note, filter) {
			notes = append(notes, note)
		}
	}

	// Sort by timestamp (newest first)
	r.sortNotesByTime(notes)

	// Apply limit
	if filter != nil && filter.Limit > 0 && len(notes) > filter.Limit {
		notes = notes[:filter.Limit]
	}

	return notes, nil
}

func (r *JSONRepository) readNoteFile(filename string) (*models.Note, error) {
	path := filepath.Join(r.notesBasePath, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var note models.Note
	if err := json.Unmarshal(data, &note); err != nil {
		return nil, err
	}

	return &note, nil
}

func (r *JSONRepository) readSessionFile(filename string) (*models.Session, error) {
	path := filepath.Join(r.sessionBasePath, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var session models.Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *JSONRepository) matchesFilter(note *models.Note, filter *repository.NoteFilter) bool {
	if filter == nil {
		return true
	}

	// Filter by project path
	if filter.ProjectPath != "" && note.ProjectPath != filter.ProjectPath {
		return false
	}

	// Filter by branch
	if filter.Branch != "" && note.Branch != filter.Branch {
		return false
	}

	// Filter by tags (note has ALL specified tags)
	if len(filter.Tags) > 0 {
		tagMap := make(map[string]bool)
		for _, tag := range note.Tags {
			tagMap[tag] = true
		}

		for _, requiredTag := range filter.Tags {
			if !tagMap[requiredTag] {
				return false
			}
		}
	}

	// Filter by time
	if !filter.Since.IsZero() && note.CreatedAt.Before(filter.Since) {
		return false
	}

	return true
}

func (r *JSONRepository) sortNotesByTime(notes []*models.Note) {
	// Simple bubble sort for now (optimize later if needed)
	for i := 0; i < len(notes)-1; i++ {
		for j := i + 1; j < len(notes); j++ {
			if notes[i].CreatedAt.Before(notes[j].CreatedAt) {
				notes[i], notes[j] = notes[j], notes[i]
			}
		}
	}
}

func (r *JSONRepository) SaveNote(note *models.Note) error {
	if note.ID == "" {
		note.ID = utils.GenerateId()
	}

	if note.CreatedAt.IsZero() {
		note.CreatedAt = time.Now()
	}
	note.UpdatedAt = time.Now()

	path := filepath.Join(r.notesBasePath, note.ID+".json")
	data, err := json.MarshalIndent(note, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal note: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

func (r *JSONRepository) GetNote(id string) (*models.Note, error) {
	path := filepath.Join(r.notesBasePath, id+".json")

	// Check if exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.ErrNoteNotFound
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var note models.Note
	if err := json.Unmarshal(data, &note); err != nil {
		return nil, err
	}

	return &note, nil
}

func (r *JSONRepository) UpdateNote(note *models.Note) error {
	// First verify note exists
	existing, err := r.GetNote(note.ID)
	if err != nil {
		return err
	}

	// Preserve creation time
	note.CreatedAt = existing.CreatedAt
	note.UpdatedAt = time.Now()

	// Save (overwrite)
	return r.SaveNote(note)
}

func (r *JSONRepository) DeleteNote(id string) error {
	path := filepath.Join(r.notesBasePath, id+".json")

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return errors.ErrNoteNotFound
		}
		return err
	}

	return nil
}

func (r *JSONRepository) Close() error {
	// Nothing to close for JSON
	return nil
}

func (r *JSONRepository) SaveSession(session *models.Session) error {
	if session.ID == "" {
		session.ID = utils.GenerateId()
	}

	if session.StartTime.IsZero() {
		session.StartTime = time.Now()
	}

	path := filepath.Join(r.sessionBasePath, session.ID+".json")
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal note: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

func (r *JSONRepository) GetSession(id string) (*models.Session, error) {
	path := filepath.Join(r.sessionBasePath, id+".json")

	// Check if exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.ErrSessionNotFound
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var session models.Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	return &session, nil
}


func (r *JSONRepository) GetOpenSession() (*models.Session, error) {
	
	files, err := os.ReadDir(r.sessionBasePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &models.Session{}, nil // No notes yet
		}
		return nil, fmt.Errorf("failed to read notes directory: %w", err)
	}
	
	var session models.Session
	
	for _, file := range files {
		// Skip directories and non-JSON files
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		s, err := r.readSessionFile(file.Name())

		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to read note %s: %v\n", file.Name(), err)
			continue
		}

		if s.ID != "" && s.EndTime.IsZero() {
			session = *s
			break
		}
	}


	return &session, nil
}

func (r *JSONRepository) ListSessions() ([]*models.Session, error) {
// Read all note files
files, err := os.ReadDir(r.sessionBasePath)
if err != nil {
	if os.IsNotExist(err) {
		return []*models.Session{}, nil // No sessions yet
	}
	return nil, fmt.Errorf("failed to read notes directory: %w", err)
}

var sessions []*models.Session

for _, file := range files {
	// Skip directories and non-JSON files
	if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
		continue
	}

	// Read and parse each note
	session, err := r.readSessionFile(file.Name())
	if err != nil {
		// Log but don't stop - could be corrupted file
		fmt.Fprintf(os.Stderr, "Warning: failed to read note %s: %v\n", file.Name(), err)
		continue
	}
	sessions = append(sessions, session)
}
return sessions, nil
}