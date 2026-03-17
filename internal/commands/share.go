package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Murchoid/iwashere/internal/domain/models"
	"github.com/Murchoid/iwashere/internal/encryption"
	"github.com/Murchoid/iwashere/internal/repository"
	"github.com/Murchoid/iwashere/internal/services/git"
	"github.com/Murchoid/iwashere/internal/utils"
)

type ShareCommand struct {
	spec        *CommandSpec
	baseCommand BaseCommand
}

func NewShareCommand() Command {
	return &ShareCommand{
		spec: ShareCommandSpec,
		baseCommand: BaseCommand{
			NameStr:  "share",
			DescStr:  "Share notes with teammates",
			UsageStr: "iwashere share [note-id] --with <recipient>",
			ExamplesList: []string{
				"iwashere share --with alice@example.com              # Share latest note",
				"iwashere share 123 --with alice@example.com          # Share specific note",
				"iwashere share --with @backend-team                   # Share with team",
				"iwashere share --with alice@example.com,bob@example.com # Multiple recipients",
			},
		},
	}
}

func (c *ShareCommand) Name() string {
	return c.baseCommand.Name()
}

func (c *ShareCommand) Description() string {
	return c.baseCommand.Description()
}

func (c *ShareCommand) Usage() string {
	return c.baseCommand.Usage()
}

func (c *ShareCommand) Examples() []string {
	return c.baseCommand.Examples()
}

func (c *ShareCommand) Execute(ctx *Context) error {
	if ctx.Repo == nil {
		return fmt.Errorf("not in an iwashere project (run iwashere init first)")
	}

	parsedArgs, err := c.spec.Parse(ctx.Args)

	if err != nil {
		utils.PrintCommandHelp(c.Name(), c.Description(), c.Usage(), c.Examples())
		return fmt.Errorf("invalid arguments: %w", err)
	}

	with := parsedArgs.Flags["with"]
	pWith, err := with.String()
	if err != nil && with.Present {
		return err
	}

	// Get git info for current user
	gitService := git.NewService(ctx.WorkDir)
	gitInfo, _ := gitService.GetInfo()
	currentEmail := gitInfo.UserEmail
	if currentEmail == "" {
		// Fallback to a default if no git email
		currentEmail = "unknown@local"
	}

	// Get the note to share
	var privateNote *models.PrivateNote

	var noteId string
	if len(parsedArgs.Positional) > 0 {
		noteId = parsedArgs.Positional[0]
	}
	if noteId == "" {
		// Share latest note
		notes, err := ctx.Repo.ListNotes(&repository.NoteFilter{Limit: 1})
		if err != nil {
			return err
		}
		if len(notes) == 0 {
			return fmt.Errorf("no notes found to share")
		}
		privateNote = notes[0]
	} else {
		privateNote, err = ctx.Repo.GetNote(noteId)
		if err != nil {
			return fmt.Errorf("Note not found: %w", err)
		}
	}

	// Parse recipients
	recipients := parseRecipients(pWith)

	// Track success/failure
	successCount := 0
	var errors []string

	// Share with each recipient
	for _, recipient := range recipients {
		if strings.HasPrefix(recipient, "@") {
			// Team sharing - use team directory (git-tracked, no encryption)
			err = c.shareWithTeam(ctx, privateNote, recipient, gitInfo.UserName)
		} else {
			// Individual sharing - use encrypted payloads
			err = c.shareWithIndividual(ctx, privateNote, recipient, currentEmail)
		}

		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", recipient, err))
		} else {
			successCount++
		}
	}

	fmt.Printf("Shared note (%s) with %d recipient(s)\n", noteId[:5]+"...", successCount)
	if len(errors) > 0 {
		fmt.Println("Errors:")
		for _, errMsg := range errors {
			fmt.Printf("   • %s\n", errMsg)
		}
	}

	// Suggest next steps
	if successCount > 0 {
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("   git add .iwashere-shared/")
		fmt.Println("   git commit -m \"Share notes with team\"")
		fmt.Println("   git push")
		fmt.Println()
		fmt.Println("Recipients will need to pull and use 'iwashere show-shared'")
	}

	return nil
}

func (c *ShareCommand) shareWithIndividual(ctx *Context, note *models.PrivateNote, recipient string, sharerEmail string) error {

	noteKey, err := encryption.GenerateKey()
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}


	sharedNote := &models.SharedNote{
		Message:   note.Message,
		CreatedAt: note.CreatedAt,
		Tags:      note.Tags,
		Branch:    note.Branch,
		Author:    sharerEmail,
	}

	// Add session info if available
	if note.SessionID != "" {
		if session, err := ctx.Repo.GetSession(note.SessionID); err == nil {
			sharedNote.SessionName = session.Title
		}
	}


	noteJSON, err := json.Marshal(sharedNote)
	if err != nil {
		return fmt.Errorf("failed to marshal note: %w", err)
	}


	encryptedNote, iv, err := encryption.Encrypt(noteJSON, noteKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt note: %w", err)
	}


	encryptedKey, err := encryption.EncryptKeyForRecipient(noteKey, recipient)
	if err != nil {
		return fmt.Errorf("failed to encrypt key: %w", err)
	}


	payload := &models.EncryptedPayload{
		NoteID:        note.ID,
		EncryptedNote: encryptedNote,
		EncryptedKey:  encryptedKey,
		IV:            iv,
		SharedBy:      sharerEmail,
		SharedAt:      time.Now(),
		NotePreview:   truncate(note.Message, 50), //inly first 50 chars
	}


	// Use .iwashere-shared/ (not .iwashere/) so it can be in git
	sharedDir := filepath.Join(ctx.ProjectPath, ".iwashere-shared", recipient)
	if err := os.MkdirAll(sharedDir, 0755); err != nil {
		return fmt.Errorf("failed to create shared directory: %w", err)
	}

	payloadPath := filepath.Join(sharedDir, note.ID+".share")
	payloadData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	if err := os.WriteFile(payloadPath, payloadData, 0644); err != nil {
		return fmt.Errorf("failed to save payload: %w", err)
	}

	c.updateShareIndex(ctx, recipient, note.ID)

	return nil
}

func (c *ShareCommand) shareWithTeam(ctx *Context, note *models.PrivateNote, teamFlag string, author string) error {
	teamName := strings.TrimPrefix(teamFlag, "@")

	// Create team note (sanitized, no encryption for team)
	teamNote := &models.TeamNote{
		ID:        utils.GenerateId(),
		Message:   note.Message,
		Author:    author,
		CreatedAt: time.Now(),
		Tags:      note.Tags,
	}

	if note.SessionID != "" {
		if session, err := ctx.Repo.GetSession(note.SessionID); err == nil {
			teamNote.SessionName = session.Title
		}
	}

	// Save to team directory (git-tracked)
	teamDir := filepath.Join(ctx.ProjectPath, ".iwashere-shared", "team", teamName)
	if err := os.MkdirAll(teamDir, 0755); err != nil {
		return fmt.Errorf("failed to create team directory: %w", err)
	}

	payloadData, err := json.MarshalIndent(teamNote, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	teamPath := filepath.Join(teamDir, note.ID+".team")
	fmt.Println(teamPath)
	if err := os.WriteFile(teamPath, payloadData, 0644); err != nil {
		return fmt.Errorf("failed to save payload: %w", err)
	}

	return nil
}

func (c *ShareCommand) updateShareIndex(ctx *Context, recipient, noteID string) {
	indexPath := filepath.Join(ctx.ProjectPath, ".iwashere-shared", "index.json")

	var index models.EncryptedPayloadIndex

	// Read existing index if it exists
	if data, err := os.ReadFile(indexPath); err == nil {
		json.Unmarshal(data, &index)
	}

	if index.Shares == nil {
		index.Shares = make(map[string][]string)
	}

	index.Shares[recipient] = append(index.Shares[recipient], noteID)

	if data, err := json.MarshalIndent(index, "", "  "); err == nil {
		os.WriteFile(indexPath, data, 0644)
	}
}

func parseRecipients(recipients string) []string {
	return strings.Split(recipients, ",")
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

func init() {
	Register("share", NewShareCommand)
}
