package jsonRepo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/Murchoid/iwashere/internal/domain/errors"
	"github.com/Murchoid/iwashere/internal/domain/models"
	"github.com/Murchoid/iwashere/internal/repository"
	"github.com/Murchoid/iwashere/internal/utils"
)

// Json implementation as a storage option
type JSONRepository struct {
	notesBasePath    string // .iwashere/notes/
	sessionBasePath  string // .iwashere/sessions/
	reminderBasePath string // .iwashere/reminder/
}

func NewJSONRepository(iwasherePath string) *JSONRepository {
	notesPath := filepath.Join(iwasherePath, "notes")
	sessionPath := filepath.Join(iwasherePath, "sessions")
	reminderPath := filepath.Join(iwasherePath, "reminders")
	// Ensure directory exists
	os.MkdirAll(notesPath, 0755)
	os.MkdirAll(sessionPath, 0755)
	os.MkdirAll(reminderPath, 0755)

	return &JSONRepository{notesBasePath: notesPath, sessionBasePath: sessionPath, reminderBasePath: reminderPath}
}

// Notes
func (r *JSONRepository) ListNotes(filter *repository.NoteFilter) ([]*models.PrivateNote, error) {
	// Read all note files
	files, err := os.ReadDir(r.notesBasePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*models.PrivateNote{}, nil // No notes yet
		}
		return nil, fmt.Errorf("failed to read notes directory: %w", err)
	}

	var notes []*models.PrivateNote

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
		// notes = append(notes, note)
	}

	// Sort by timestamp (newest first)
	r.sortNotesByTime(notes)

	// Apply limit
	if filter != nil && filter.Limit > 0 && len(notes) > filter.Limit {
		notes = notes[:filter.Limit]
	}

	return notes, nil
}

func (r *JSONRepository) SaveNote(note *models.PrivateNote) error {
	if note.ID == "" {
		note.ID = utils.GenerateId()
	}

	if note.CreatedAt.IsZero() {
		note.CreatedAt = time.Now()
	}
	note.UpdatedAt = time.Now()
	note.Tags = slices.Compact(note.Tags)

	path := filepath.Join(r.notesBasePath, note.ID+".json")
	data, err := json.MarshalIndent(note, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal note: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

func (r *JSONRepository) SaveTeamNote(note *models.TeamNote) error {
	if note.ID == "" {
		note.ID = utils.GenerateId()
	}

	if note.CreatedAt.IsZero() {
		note.CreatedAt = time.Now()
	}
	note.Tags = slices.Compact(note.Tags)

	path := filepath.Join(r.notesBasePath, note.ID+".json")
	data, err := json.MarshalIndent(note, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal note: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

func (r *JSONRepository) GetNote(id string) (*models.PrivateNote, error) {
	path := filepath.Join(r.notesBasePath, id+".json")

	// Check if exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.ErrNoteNotFound
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var note models.PrivateNote
	if err := json.Unmarshal(data, &note); err != nil {
		return nil, err
	}

	return &note, nil
}

func (r *JSONRepository) UpdateMessage(note *models.PrivateNote) error {
	// First verify note exists
	existing, err := r.GetNote(note.ID)
	if err != nil {
		return err
	}

	note.CreatedAt = existing.CreatedAt
	note.Branch = existing.Branch
	note.SessionID = existing.SessionID
	note.Tags = existing.Tags
	note.CommitHash = existing.CommitHash
	note.CommitMsg = existing.CommitMsg
	note.ModifiedFiles = existing.ModifiedFiles
	note.ProjectPath = existing.ProjectPath
	note.UpdatedAt = time.Now()

	// Save (overwrite)
	return r.SaveNote(note)
}

func (r *JSONRepository) UpdateTags(note *models.PrivateNote) error {
	// First verify note exists
	existing, err := r.GetNote(note.ID)
	if err != nil {
		return err
	}

	note.Tags = slices.Compact(note.Tags)
	note.CreatedAt = existing.CreatedAt
	note.Branch = existing.Branch
	note.Message = existing.Message
	note.SessionID = existing.SessionID
	note.CommitHash = existing.CommitHash
	note.CommitMsg = existing.CommitMsg
	note.ModifiedFiles = existing.ModifiedFiles
	note.ProjectPath = existing.ProjectPath
	note.UpdatedAt = time.Now()

	// Save (overwrite)
	return r.SaveNote(note)
}

func (r *JSONRepository) AddTagsToNote(note *models.PrivateNote) error {
	// First verify note exists
	existing, err := r.GetNote(note.ID)
	if err != nil {
		return err
	}

	note.CreatedAt = existing.CreatedAt
	note.Branch = existing.Branch
	note.Message = existing.Message
	note.Tags = append(note.Tags, existing.Tags...)
	note.Tags = slices.Compact(note.Tags)
	note.SessionID = existing.SessionID
	note.CommitHash = existing.CommitHash
	note.CommitMsg = existing.CommitMsg
	note.ModifiedFiles = existing.ModifiedFiles
	note.ProjectPath = existing.ProjectPath
	note.UpdatedAt = time.Now()

	// Save (overwrite)
	return r.SaveNote(note)
}

func (r *JSONRepository) RemoveTagsFromNote(note *models.PrivateNote) error {
	// First verify note exists
	existing, err := r.GetNote(note.ID)
	if err != nil {
		return err
	}

	removeSet := make(map[string]struct{})
	for _, rTag := range note.Tags {
		removeSet[rTag] = struct{}{}
	}

	var newTags []string
	for _, tag := range existing.Tags {
		if _, found := removeSet[tag]; !found {
			newTags = append(newTags, tag)
		}
	}

	note.Tags = newTags
	note.CreatedAt = existing.CreatedAt
	note.Branch = existing.Branch
	note.Message = existing.Message
	note.SessionID = existing.SessionID
	note.CommitHash = existing.CommitHash
	note.CommitMsg = existing.CommitMsg
	note.ModifiedFiles = existing.ModifiedFiles
	note.ProjectPath = existing.ProjectPath
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

// Sessions
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

func (r *JSONRepository) GetNotesBySession(id string) ([]*models.PrivateNote, error) {
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

	filter := repository.NoteFilter{
		SessionID: session.ID,
		Branch:    session.Branch,
	}

	notes, err := r.ListNotes(&filter)
	if err == nil {
		return nil, err
	}

	return notes, nil
}

func (r *JSONRepository) GetOpenSession() (*models.Session, error) {

	files, err := os.ReadDir(r.sessionBasePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &models.Session{}, nil // No notes yet
		}
		return nil, fmt.Errorf("failed to read notes directory: %w", err)
	}

	var session *models.Session

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

		if s.ID != "" && s.EndTime.IsZero() || s.State == models.Paused || s.State == models.Continued || s.State == models.Ongoing {
			session = s
			break
		} else {
			continue
		}
	}

	return session, nil
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

// Reminders
func (r *JSONRepository) SaveReminder(reminder *models.Reminder) error {
	if reminder.ID == "" {
		reminder.ID = utils.GenerateId()
	}

	if reminder.CreatedAt.IsZero() {
		reminder.CreatedAt = time.Now()
	}

	path := filepath.Join(r.reminderBasePath, reminder.ID+".json")
	data, err := json.MarshalIndent(reminder, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal note: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

func (r *JSONRepository) GetReminder(id string) (*models.Reminder, error) {
	path := filepath.Join(r.reminderBasePath, id+".json")

	// Check if exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.ErrNoteNotFound
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var reminder models.Reminder
	if err := json.Unmarshal(data, &reminder); err != nil {
		return nil, err
	}

	return &reminder, nil
}

func (r *JSONRepository) DeactivateOrUpdateReminder(id string) error {
	existing, err := r.GetReminder(id)
	if err != nil {
		return err
	}

	if existing.Repeats == "none" || existing.Repeats == "" {
		existing.Active = false
	} else {
		existing.DueAt = getNextTimeDependingOnRepeat(existing.Repeats, existing.DueAt)
	}

	return r.SaveReminder(existing)
}

func (r *JSONRepository) ListReminders() ([]*models.Reminder, error) {
	// Read all reminder files
	files, err := os.ReadDir(r.reminderBasePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*models.Reminder{}, nil // No reminders yet
		}
		return nil, fmt.Errorf("failed to read reminder directory: %w", err)
	}

	var reminders []*models.Reminder

	for _, file := range files {
		// Skip directories and non-JSON files
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		// Read and parse each note
		reminder, err := r.readReminderFile(file.Name())
		if err != nil {
			// Log but don't stop - could be corrupted file
			fmt.Fprintf(os.Stderr, "Warning: failed to read note %s: %v\n", file.Name(), err)
			continue
		}
		reminders = append(reminders, reminder)
	}
	return reminders, nil
}

func (r *JSONRepository) ListDueReminders() ([]*models.Reminder, error) {
	// Read all reminder files
	files, err := os.ReadDir(r.reminderBasePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*models.Reminder{}, nil // No reminders yet
		}
		return nil, fmt.Errorf("failed to read reminder directory: %w", err)
	}

	var reminders []*models.Reminder

	for _, file := range files {
		// Skip directories and non-JSON files
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		// Read and parse each note
		reminder, err := r.readReminderFile(file.Name())
		if err != nil {
			// Log but don't stop - could be corrupted file
			fmt.Fprintf(os.Stderr, "Warning: failed to read note %s: %v\n", file.Name(), err)
			continue
		}

		//only append due reminders
		if reminder.Active && (reminder.DueAt.Compare(time.Now()) == -1 || reminder.DueAt.Compare(time.Now()) == 0) {
			reminders = append(reminders, reminder)
		}
	}
	return reminders, nil
}

func (r *JSONRepository) DeleteReminder(id string) error {
	path := filepath.Join(r.reminderBasePath, id+".json")

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return errors.ErrNoteNotFound
		}
		return err
	}

	return nil
}

// Helpers
func (r *JSONRepository) readNoteFile(filename string) (*models.PrivateNote, error) {
	path := filepath.Join(r.notesBasePath, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var note models.PrivateNote
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

func (r *JSONRepository) readReminderFile(filename string) (*models.Reminder, error) {
	path := filepath.Join(r.reminderBasePath, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var reminder models.Reminder
	if err := json.Unmarshal(data, &reminder); err != nil {
		return nil, err
	}

	return &reminder, nil
}

func (r *JSONRepository) matchesFilter(note *models.PrivateNote, filter *repository.NoteFilter) bool {
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

	//Filter by session

	if note.SessionID == filter.SessionID {
		return true
	}

	return true
}

func (r *JSONRepository) sortNotesByTime(notes []*models.PrivateNote) {
	// Simple bubble sort for now (optimize later if needed)
	for i := 0; i < len(notes)-1; i++ {
		for j := i + 1; j < len(notes); j++ {
			if notes[i].CreatedAt.Before(notes[j].CreatedAt) {
				notes[i], notes[j] = notes[j], notes[i]
			}
		}
	}
}

func getNextTimeDependingOnRepeat(repeat string, currentTime time.Time) time.Time {
	switch repeat {
	case "daily":
		return currentTime.Add(24 * time.Hour)
	case "weekly":
		return currentTime.Add(7 * 24 * time.Hour)
	case "monthly":
		return currentTime.AddDate(0, 1, 0)
	case "yearly":
		return currentTime.AddDate(1, 0, 0)
	default:
		return currentTime
	}
}

// Close
func (r *JSONRepository) Close() error {
	// Nothing to close for JSON
	return nil
}
