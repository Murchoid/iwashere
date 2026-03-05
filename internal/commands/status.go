package commands

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Murchoid/iwashere/internal/domain/models"
	"github.com/Murchoid/iwashere/internal/services/git"
	"github.com/Murchoid/iwashere/internal/utils"
)

type StatusCommand struct {
	BaseCommand
}

func NewStatusCommandFactory() Command {
	return &StatusCommand{
		BaseCommand{
			NameStr:  "status",
			DescStr:  "Quick status of what you were doing",
			UsageStr: "iwashere status ",
			ExamplesList: []string{
				"iwashere status",
			},
		},
	}
}

func (a *StatusCommand) Name() string {
	return a.BaseCommand.Name()
}

func (a *StatusCommand) Description() string {
	return a.BaseCommand.Description()
}

func (a *StatusCommand) Usage() string {
	return a.BaseCommand.Usage()
}

func (a *StatusCommand) Examples() []string {
	return a.BaseCommand.Examples()
}

func (a *StatusCommand) Execute(ctx *Context) error {
	if ctx.Repo == nil {
		return fmt.Errorf("not in an iwashere project (run iwashere init first)")
	}

	// Check for unrecognized arguments
	if len(ctx.Args) > 0 {
		fmt.Println("Unrecognized arguments")
		fmt.Println()
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
		return nil
	}

	repo := ctx.Repo

	// Get active session
	session, err := repo.GetOpenSession()
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	// Get notes for the session
	var notes []*models.PrivateNote
	if session != nil && len(session.Notes) > 0 {
		for _, noteId := range session.Notes {
			note, err := repo.GetNote(noteId)
			if err != nil {
				// Log but don't fail - just skip corrupted notes
				fmt.Fprintf(os.Stderr, "Warning: failed to fetch note %s: %v\n", noteId, err)
				continue
			}
			notes = append(notes, note)
		}
	}

	// Get git info for unstaged files
	gitService := git.NewService(ctx.WorkDir)

	printStatus(session, notes, gitService)
	return nil
}

func printStatus(session *models.Session, notes []*models.PrivateNote, gitService *git.Service) {

	fmt.Println("iwashere status")
	fmt.Println("=================")
	fmt.Println()
	gitInfo, _ := gitService.GetInfo() // Ignore error, git might not be present

	// Session info
	if session == nil || session.ID == "" {
		fmt.Println("No active session")
		fmt.Println("Start one with: iwashere session start \"session name\"")
	} else {
		fmt.Printf("You were working on '%s' (%s)\n",
			session.Title,
			utils.HowLongAgo(session.StartTime))

		if session.EndTime.IsZero() {
			fmt.Printf("Session ongoing")
		} else {
			duration := session.EndTime.Sub(session.StartTime).Round(time.Minute)
			fmt.Printf("Session lasted %s\n", duration)
		}
	}
	fmt.Println()

	// Last note
	if len(notes) > 0 {
		lastNote := notes[len(notes)-1]
		fmt.Printf("Last note: %s\n", lastNote.Message)

		// Show tags if any
		if len(lastNote.Tags) > 0 {
			fmt.Printf("%s\n", strings.Join(lastNote.Tags, ", "))
		}
		fmt.Println()
	}

	// Modified files (combine from notes and git)
	modifiedFiles := make(map[string]bool)

	// Get files from notes
	for _, note := range notes {
		for _, file := range note.ModifiedFiles {
			modifiedFiles[file] = true
		}
	}

	// Get unstaged files from git
	if gitInfo != nil && gitInfo.HasChanges {
		files, _ := gitService.GetModifiedFiles()
		for _, file := range files {
			modifiedFiles[file] = true
		}
	}

	if len(modifiedFiles) > 0 {
		fmt.Println("Modified files:")

		// Sort files for consistent display
		var fileList []string
		for file := range modifiedFiles {
			fileList = append(fileList, file)
		}
		sort.Strings(fileList)

		for _, file := range fileList {
			// Check if unstaged in git
			unstaged := ""
			if gitInfo != nil && gitInfo.HasChanges {
				// You'd need a more sophisticated check here
				unstaged = " (unstaged)"
			}
			fmt.Printf("   • %s%s\n", file, unstaged)
		}
		fmt.Println()
	}

	// Related notes (last 5 notes, excluding current session)
	if len(notes) > 0 {
		fmt.Println("Related notes from this session:")

		// Show notes in reverse order (newest first)
		for i := len(notes) - 1; i >= 0; i-- {
			if i < len(notes)-5 { // Only show last 5
				break
			}
			printRelatedNote(notes[i], i == len(notes)-1)
		}
	}

	// Next steps suggestion
	fmt.Println()
	fmt.Println("What's next?")
	if session == nil {
		fmt.Println("   • Start a session: iwashere session start \"feature name\"")
	} else {
		fmt.Println("   • Add a note: iwashere add \"next task\"")
		fmt.Println("   • End session: iwashere session end")
	}
	fmt.Println("   • View all notes: iwashere list")
}

func printRelatedNote(note *models.PrivateNote, isLast bool) {
	prefix := "   • "
	if isLast {
		prefix = "   └─ "
	}

	// Format time
	timeStr := utils.HowLongAgo(note.CreatedAt)

	// Format tags if any
	tagsStr := ""
	if len(note.Tags) > 0 {
		tagsStr = fmt.Sprintf(" [%s]", strings.Join(note.Tags, ", "))
	}

	fmt.Printf("%s%s - %s%s\n",
		prefix,
		timeStr,
		note.Message,
		tagsStr)
}

func init() {
	Register("status", NewStatusCommandFactory)
}
