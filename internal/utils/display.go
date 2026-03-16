package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/Murchoid/iwashere/internal/domain/models"
	"github.com/Murchoid/iwashere/internal/repository"
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

type noteGroup struct {
	Session *models.Session
	Notes   []*models.PrivateNote
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

// Notes display function
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
				ShowSession: false,
				ShowFiles:   true,
				ShowTags:    true,
				ShowID:      true,
				Format:      format,
			}
			fmt.Println(display.String())
		}
	}
}

//Session display function

// PrintSessions displays sessions in a format consistent with notes
func PrintSessions(sessions []*models.Session, showNotes bool, repo repository.Repository) {
	if len(sessions) == 0 {
		fmt.Println("No sessions found")
		fmt.Println("   Start one with: iwashere session start \"session name\"")
		return
	}

	fmt.Println(" Sessions")
	fmt.Println("================")
	fmt.Println()

	for i, session := range sessions {
		printSession(session, i == len(sessions)-1, showNotes, repo)
	}
}

// PrintCurrentSession shows the active session prominently
func PrintCurrentSession(session *models.Session, notes []*models.PrivateNote) {
	fmt.Println(" Active Session")
	fmt.Println("================")
	fmt.Println()

	// Session header with duration
	fmt.Printf("  %s", session.Title)
	if session.EndTime.IsZero() {
		fmt.Printf(" (started %s, ongoing)\n", HowLongAgo(session.StartTime, 0))
	} else {
		sessionDuration := session.EndTime.Sub(session.StartTime).Round(time.Minute)
		fmt.Printf(" (started %s, lasted %s)\n",
			HowLongAgo(session.StartTime, 0),
			sessionDuration)
	}

	// Session summary if exists
	if session.Summary != "" {
		fmt.Printf("   %s\n", session.Summary)
	}

	fmt.Println()

	// Show notes in this session
	if len(notes) > 0 {
		fmt.Println("  Notes in this session:")
		for i, note := range notes {
			prefix := "  ├─ "
			if i == len(notes)-1 {
				prefix = "  └─ "
			}
			fmt.Printf("%s[%s] %s\n",
				prefix,
				HowLongAgo(note.CreatedAt, 0),
				note.Message)

			// Show tags if any (indented)
			if len(note.Tags) > 0 {
				fmt.Printf("  %s    %s\n",
					strings.Repeat(" ", len(prefix)-2),
					strings.Join(note.Tags, ", "))
			}
		}
	}
}

// PrintSessionDetails shows comprehensive session info
func PrintSessionDetails(session *models.Session, notes []*models.PrivateNote) {
	fmt.Printf("Session: %s\n", session.Title)
	fmt.Println(strings.Repeat("=", len(session.Title)+10))
	fmt.Println()

	// Timeline
	fmt.Printf("Started:  %s (%s)\n",
		session.StartTime.Format("Jan 2, 2006 at 15:04"),
		HowLongAgo(session.StartTime, 0))

	if !session.EndTime.IsZero() {
		fmt.Printf("Ended:    %s (%s)\n",
			session.EndTime.Format("Jan 2, 2006 at 15:04"),
			HowLongAgo(session.EndTime, 0))

		duration := session.EndTime.Sub(session.StartTime).Round(time.Minute)
		fmt.Printf("Duration: %s\n", duration)
	} else {
		duration := time.Since(session.StartTime).Round(time.Minute)
		fmt.Printf("Duration: %s (ongoing)\n", duration)
	}

	fmt.Println()

	// Summary
	if session.Summary != "" {
		fmt.Printf("Summary: %s\n", session.Summary)
		fmt.Println()
	}

	// Notes in this session
	if len(notes) > 0 {
		fmt.Printf("Notes (%d):\n", len(notes))
		fmt.Println()

		for i, note := range notes {
			fmt.Printf("  %d. [%s] %s\n",
				i+1,
				HowLongAgo(note.CreatedAt, 0),
				note.Message)

			if len(note.Tags) > 0 {
				fmt.Printf("      %s\n", strings.Join(note.Tags, ", "))
			}

			if note.Branch != "" {
				fmt.Printf("     %s\n", note.Branch)
			}

			if i < len(notes)-1 {
				fmt.Println()
			}
		}
	}
}

// Reminder display
// PrintReminders displays sessions in a format consistent with notes
func PrintReminders(reminders []*models.Reminder, showNotes bool, repo repository.Repository) {
	if len(reminders) == 0 {
		fmt.Println("No reminders found")
		fmt.Println("   create one with: iwashere remind 123 \"reminder message\"")
		return
	}

	fmt.Println(" Reminders")
	fmt.Println("================")
	fmt.Println()

	for i, reminder := range reminders {
		printReminder(reminder, i == len(reminders)-1, showNotes, repo)
	}
}

// Helpers
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
	var duration string
	switch session.State {

	case models.Ended:
		duration = session.EndTime.Sub(session.StartTime).Round(time.Minute).String()
	case models.Paused:
		duration = "paused"
	case models.Continued, models.Ongoing:
		duration = "ongoing"
	}

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
	timeStr := HowLongAgo(d.Note.CreatedAt, 0)
	parts = append(parts, fmt.Sprintf("%s[%s%s]",
		ColorGray, timeStr, ColorReset))

	// Branch with Color
	if d.Note.Branch != "" {
		parts = append(parts, fmt.Sprintf("%s(%s)%s",
			ColorGreen, d.Note.Branch, ColorReset))
	}

	// ID (shortened)
	if d.ShowID {
		shortID := d.Note.ID
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
	timeStr := HowLongAgo(d.Note.CreatedAt, 0)
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
		HowLongAgo(d.Note.CreatedAt, 0),
		d.Note.Message)
}

func printSession(session *models.Session, isLast bool, showNotes bool, repo repository.Repository) {
	// Choose the right tree character
	prefix := "├─ "
	if isLast {
		prefix = "└─ "
	}

	// Format duration
	var durationStr string
	if session.State == models.Ongoing || session.State == models.Continued {
		durationStr = fmt.Sprintf("(started %s, ongoing)", HowLongAgo(session.StartTime, 0))
	} else if session.State == models.Paused {
		durationStr = fmt.Sprintf("(started %s, paused)", HowLongAgo(session.StartTime, 0))
	} else {
		sessionDuration := session.EndTime.Sub(session.StartTime).Round(time.Minute)
		durationStr = fmt.Sprintf("(started %s, lasted %s)",
			HowLongAgo(session.StartTime, 0),
			sessionDuration)
	}

	// Print session header
	fmt.Printf("%s %s (%s) %s\n", prefix, session.ID, session.Title, durationStr)

	// Print session summary if exists (indented)
	if session.Summary != "" {
		fmt.Printf("  %s   %s\n",
			strings.Repeat(" ", len(prefix)-2),
			session.Summary)
	}

	// Print notes in this session if requested
	if showNotes && len(session.Notes) > 0 && repo != nil {
		notes, _ := repo.GetNotesBySession(session.ID)
		if len(notes) > 0 {
			for i, note := range notes {
				notePrefix := "    ├─ "
				if i == len(notes)-1 {
					notePrefix = "    └─ "
				}
				fmt.Printf("%s[%s] %s\n",
					notePrefix,
					HowLongAgo(note.CreatedAt, 0),
					note.Message)
			}
		}
	}

	// Blank line between sessions for readability
	if !isLast {
		fmt.Println()
	}
}

func printReminder(reminder *models.Reminder, isLast bool, showNotes bool, repo repository.Repository) {

	prefix := "├─ "
	if isLast {
		prefix = "└─ "
	}

	// Format duration
	var durationStr string
	if reminder.Active {
		durationStr = fmt.Sprintf("(Created %s, Due at %s)", HowLongAgo(reminder.CreatedAt, 0), reminder.DueAt.Format("Mon Jan 2 at 15:04"))
	} else if !reminder.Active {
		durationStr = fmt.Sprintf("(Created %s, Done)", HowLongAgo(reminder.CreatedAt, 0))
	}

	// Print session header
	var formatedMessage string
	if len([]byte(reminder.Message)) > 20 {
		formatedMessage = reminder.Message[:20] + "..."
	} else {
		formatedMessage = reminder.Message
	}

	fmt.Printf("%s %s (%s) %s\n", prefix, reminder.ID, formatedMessage, durationStr)

	// Print notes in this reminder if requested
	if showNotes && reminder.NoteID != "" && repo != nil {
		note, _ := repo.GetNote(reminder.NoteID)

		notePrefix := "    └─ "
		fmt.Printf("%s[%s] %s\n",
			notePrefix,
			HowLongAgo(note.CreatedAt, 0),
			note.Message)

	}

	// Blank line between sessions for readability
	if !isLast {
		fmt.Println()
	}
}
