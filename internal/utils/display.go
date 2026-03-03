package utils

import (
	"fmt"
	"strings"
	"time"

	"githum.com/Murchoid/iwashere/internal/domain/models"
)

type NoteDisplay struct {
	Note        *models.PrivateNote
	Session     *models.Session
	ShowSession bool
	ShowFiles   bool
	ShowTags    bool
	ShowID      bool
	Format      string // "short", "detailed", "compact"
}

func (d *NoteDisplay) String() string {
	switch d.Format {
	case "compact":
		return d.compactFormat()
	case "short":
		return d.shortFormat()
	default:
		return d.detailedFormat()
	}
}

// Main display function
func PrintNotes(notes []*models.PrivateNote, sessions map[string]*models.Session, format string) {
	if len(notes) == 0 {
		fmt.Println("No notes found")
		return
	}

	// Group notes by session
	sessionGroups := groupNotesBySession(notes, sessions)

	for i, group := range sessionGroups {
		if i > 0 {
			fmt.Println() // Blank line between sessions
		}

		// Print session header if this group has a session
		if group.Session != nil {
			printSessionHeader(group.Session)
		}

		// Print notes in this session
		for _, note := range group.Notes {
			display := &NoteDisplay{
				Note:        note,
				Session:     group.Session,
				ShowSession: false, // Already showing session header
				ShowFiles:   true,
				ShowTags:    true,
				ShowID:      true,
				Format:      format,
			}
			fmt.Println(display.String())
		}
	}
}

type noteGroup struct {
	Session *models.Session
	Notes   []*models.PrivateNote
}

func groupNotesBySession(notes []*models.PrivateNote, sessions map[string]*models.Session) []noteGroup {
	var groups []noteGroup

	// First, group by session
	sessionMap := make(map[string][]*models.PrivateNote)
	var standalone []*models.PrivateNote

	for _, note := range notes {
		if note.SessionID != "" {
			sessionMap[note.SessionID] = append(sessionMap[note.SessionID], note)
		} else {
			standalone = append(standalone, note)
		}
	}

	// Create groups for sessions
	for sessionID, sessionNotes := range sessionMap {
		if session, exists := sessions[sessionID]; exists {
			groups = append(groups, noteGroup{
				Session: session,
				Notes:   sessionNotes,
			})
		}
	}

	// Add standalone notes as a pseudo-group
	if len(standalone) > 0 {
		groups = append(groups, noteGroup{
			Session: nil,
			Notes:   standalone,
		})
	}

	return groups
}

func printSessionHeader(session *models.Session) {
	duration := session.EndTime.Sub(session.StartTime).Round(time.Minute)
	fmt.Printf("%s %s%s ", ColorPurple, session.Title, ColorReset)
	fmt.Printf("(%s - %s, %s)\n",
		session.StartTime.Format("15:04"),
		session.EndTime.Format("15:04"),
		duration)
	if session.Summary != "" {
		fmt.Printf(" %s\n", session.Summary)
	}
}

func (d *NoteDisplay) detailedFormat() string {
	var parts []string

	// Time with relative format
	timeStr := HowLongAgo(d.Note.CreatedAt)
	parts = append(parts, fmt.Sprintf("%s[%s%s]",
		ColorGray, timeStr, ColorReset))

	// Branch with Color
	if d.Note.Branch != "" {
		parts = append(parts, fmt.Sprintf("%s(%s)%s",
			ColorGreen, d.Note.Branch, ColorReset))
	}

	// ID (shortened)
	if d.ShowID {
		shortID := d.Note.ID[:8]
		parts = append(parts, fmt.Sprintf("%s%s:%s",
			ColorCyan, shortID, ColorReset))
	}

	// Message (main content)
	parts = append(parts, d.Note.Message)

	// Tags on new line with indentation
	if d.ShowTags && len(d.Note.Tags) > 0 {
		tags := make([]string, len(d.Note.Tags))
		for i, tag := range d.Note.Tags {
			tags[i] = fmt.Sprintf("%s#%s%s", ColorYellow, tag, ColorReset)
		}
		parts = append(parts, fmt.Sprintf("\n  %s", strings.Join(tags, " ")))
	}

	// Modified files
	if d.ShowFiles && len(d.Note.ModifiedFiles) > 0 {
		files := make([]string, len(d.Note.ModifiedFiles))
		for i, f := range d.Note.ModifiedFiles {
			// Show just filename, not full path
			parts := strings.Split(f, "/")
			files[i] = parts[len(parts)-1]
		}
		parts = append(parts, fmt.Sprintf("\nModified files:  %s",
			strings.Join(files, ", ")))
	}

	return strings.Join(parts, " ")
}

func (d *NoteDisplay) shortFormat() string {
	// Compact, one-line format
	timeStr := HowLongAgo(d.Note.CreatedAt)
	branch := ""
	if d.Note.Branch != "" {
		branch = fmt.Sprintf("(%s) ", d.Note.Branch)
	}

	msg := d.Note.Message
	if len(msg) > 50 {
		msg = msg[:47] + "..."
	}

	return fmt.Sprintf("%s %s- %s",
		timeStr, branch, msg)
}

func (d *NoteDisplay) compactFormat() string {
	// Super compact for lists
	return fmt.Sprintf("%s %s",
		HowLongAgo(d.Note.CreatedAt),
		d.Note.Message)
}
