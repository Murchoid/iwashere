package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Murchoid/iwashere/internal/domain/models"
	"github.com/Murchoid/iwashere/internal/encryption"
	"github.com/Murchoid/iwashere/internal/services/git"
	"github.com/Murchoid/iwashere/internal/utils"
)

type ShowSharedCommand struct {
	spec        *CommandSpec
	baseCommand BaseCommand
}

func NewShowShareCommand() Command {
	return &ShowSharedCommand{
		spec: ShowSharedCommandSpec,
		baseCommand: BaseCommand{
			NameStr:  "show-shared",
			DescStr:  "Show notes shared with you",
			UsageStr: "iwashere show-shared [note-id]",
			ExamplesList: []string{
				"iwashere show-shared           # List all shared notes",
				"iwashere show-shared 123       # Show specific shared note",
			},
		},
	}
}

func (c *ShowSharedCommand) Name() string {
	return c.baseCommand.Name()
}

func (c *ShowSharedCommand) Description() string {
	return c.baseCommand.Description()
}

func (c *ShowSharedCommand) Usage() string {
	return c.baseCommand.Usage()
}

func (c *ShowSharedCommand) Examples() []string {
	return c.baseCommand.Examples()
}

func (c *ShowSharedCommand) Execute(ctx *Context) error {
	// Get current user's email from git
	gitService := git.NewService(ctx.WorkDir)
	gitInfo, _ := gitService.GetInfo()
	currentEmail := gitInfo.UserEmail
	if currentEmail == "" {
		currentEmail = "unknown@local"
	}

	sharedDir := filepath.Join(ctx.ProjectPath, ".iwashere-shared", currentEmail)
	teamName := ctx.Config.Team.TeamName
	teamDir := filepath.Join(ctx.ProjectPath, ".iwashere-shared", "team", teamName)

	// Check if directory exists
	if _, err := os.Stat(sharedDir); os.IsNotExist(err) {
		fmt.Println("No notes shared with you yet")
		return nil
	}

	parsedArgs, err := c.spec.Parse(ctx.Args)

	if err != nil {
		utils.PrintCommandHelp(c.Name(), c.Description(), c.Usage(), c.Examples())
		return fmt.Errorf("invalid arguments: %w", err)
	}

	if len(ctx.Args) > 1 {
		fmt.Println("show-shared only accpets one argument")
		fmt.Println()
		utils.PrintCommandHelp(c.Name(), c.Description(), c.Usage(), c.Examples())
		return nil
	}

	// If specific note ID provided
	if len(parsedArgs.Positional) > 0 {
		noteID := parsedArgs.Positional[0]
		return c.showSpecificNote(sharedDir, noteID, currentEmail)
	}

	//Ekse jus shhow List all shared notes
	return c.listSharedNotes(sharedDir, teamDir)
}

func (c *ShowSharedCommand) listSharedNotes(sharedDir, teamDir string) error {
	personalFiles, err := filepath.Glob(filepath.Join(sharedDir, "*.share"))
	teamFiles, err := filepath.Glob(filepath.Join(teamDir, "*.team"))
	if err != nil {
		return err
	}

	if len(personalFiles) == 0 && len(teamFiles) == 0 {
		fmt.Println("No shared notes found")
		return nil
	}

	fmt.Println("Notes shared with you")
	fmt.Println("=======================")
	fmt.Println()

	if len(personalFiles) > 0 {
		fmt.Println("Personal notes")
		fmt.Println("================")
	}
	for _, file := range personalFiles {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		var payload models.EncryptedPayload
		if err := json.Unmarshal(data, &payload); err != nil {
			continue
		}

		// Show preview without decrypting full note
		fmt.Printf(" [%s] from %s\n", payload.NoteID, payload.SharedBy)
		fmt.Printf("     %s\n", payload.NotePreview)
		fmt.Printf("     shared %s\n", utils.HowLongAgo(payload.SharedAt, 0))
		fmt.Println()
	}

	if len(teamFiles) > 0 {
		fmt.Println("Team notes")
		fmt.Println("===========")
	}

	for _, file := range teamFiles {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		var payload models.TeamNote
		if err := json.Unmarshal(data, &payload); err != nil {
			continue
		}

		fmt.Printf(" [%s] from %s\n", payload.ID, payload.Author)
		fmt.Printf("     %s\n", payload.Message)
		fmt.Printf("     shared %s\n", utils.HowLongAgo(payload.CreatedAt, 0))
		fmt.Println()
	}

	fmt.Println("Use 'iwashere show-shared <note-id>' to view full note (only works for personal notes)")
	return nil
}

func (c *ShowSharedCommand) showSpecificNote(sharedDir string, noteID string, currentEmail string) error {
	payloadPath := filepath.Join(sharedDir, noteID+".share")

	data, err := os.ReadFile(payloadPath)
	if err != nil {
		return fmt.Errorf("shared note not found: %s", noteID)
	}

	var payload models.EncryptedPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return fmt.Errorf("failed to parse shared note: %w", err)
	}

	// Decrypt the key using recipient's email
	noteKey, err := encryption.DecryptKeyForRecipient(payload.EncryptedKey, currentEmail)
	if err != nil {
		return fmt.Errorf("failed to decrypt key: %w", err)
	}

	// Decrypt the note
	noteJSON, err := encryption.Decrypt(payload.EncryptedNote, payload.IV, noteKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt note: %w", err)
	}

	var sharedNote models.SharedNote
	if err := json.Unmarshal(noteJSON, &sharedNote); err != nil {
		return fmt.Errorf("failed to parse note: %w", err)
	}

	// Display the note
	fmt.Printf("Note from %s%s%s\n",utils.ColorYellow,payload.SharedBy, utils.ColorReset)
	fmt.Printf("   shared %s%s%s\n", utils.ColorCyan,utils.HowLongAgo(payload.SharedAt, 0), utils.ColorReset)
	fmt.Println()
	fmt.Printf("   %s\n", sharedNote.Message)
	fmt.Println()

	if len(sharedNote.Tags) > 0 {
		fmt.Printf("    %s\n", strings.Join(sharedNote.Tags, ", "))
	}

	if sharedNote.Branch != "" {
		fmt.Printf("  Branch: %s%s%s\n", utils.ColorGreen,sharedNote.Branch, utils.ColorReset)
	}

	if sharedNote.SessionName != "" {
		fmt.Printf("   Session: %s%s%s\n", utils.ColorPurple, sharedNote.SessionName, utils.ColorPurple)
	}

	return nil
}

func init() {
	Register("show-shared", NewShowShareCommand)
}
