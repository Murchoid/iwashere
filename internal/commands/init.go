package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Murchoid/iwashere/internal/domain/models"
	"github.com/Murchoid/iwashere/internal/repository/jsonRepo"
	"github.com/Murchoid/iwashere/internal/utils"
)

type InitCommand struct {
	BaseCommand
}

func NewInitCommand() Command {
	return &InitCommand{
		BaseCommand: BaseCommand{
			NameStr:  "init",
			DescStr:  "Initialize iwashere in current directory",
			UsageStr: "iwashere init [--force] [--no-ignore]",
			ExamplesList: []string{
				"iwashere init",
				"iwashere init --force",
				"iwashere init --no-ignore",
			},
		},
	}
}

func (c *InitCommand) Name() string {
	return c.BaseCommand.Name()
}

func (c *InitCommand) Description() string {
	return c.BaseCommand.Description()
}

func (c *InitCommand) Usage() string {
	return c.BaseCommand.Usage()
}

func (c *InitCommand) Examples() []string {
	return c.BaseCommand.Examples()
}

func (c *InitCommand) Execute(ctx *Context) error {

	if len(ctx.Args) > 0 {
		fmt.Println("Unrecognized arguments")
		fmt.Println()
		utils.PrintCommandHelp(c.Name(), c.Description(), c.Usage(), c.Examples())
		return nil
	}

	// Check if already initialized
	// Check for force flag
	force := ctx.Flags["--force"] == "true"
	if ctx.ProjectPath != "" && !force {
		fmt.Println(".iwashere already exists (use --force to reinitialize)")
		fmt.Println("use 'iwashere init --force' if you want to forcefully reinitialize (Use this if you know what you are doing)")
		return nil
	}

	iwasherePath := filepath.Join(ctx.WorkDir, ".iwashere")

	// Check if directory exists
	if _, err := os.Stat(iwasherePath); err == nil {
		// Remove existing directory if force
		if err := os.RemoveAll(iwasherePath); err != nil {
			return fmt.Errorf("failed to remove existing .iwashere: %w", err)
		}
	}

	// Create directory structure
	dirs := []string{
		iwasherePath,
		filepath.Join(iwasherePath, "db"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create default config
	config := models.DefaultConfig()
	config.Project.Name = filepath.Base(ctx.WorkDir)
	config.Project.InitDate = time.Now()

	//Save config to file
	configPath := filepath.Join(iwasherePath, "Config", "config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal note: %w", err)
	}

	os.WriteFile(configPath, data, 0644)

	// Create .gitignore entry
	if ctx.Flags["--no-ignore"] != "true" {
		if err := c.updateGitignore(ctx.WorkDir); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to update .gitignore: %v\n", err)
		}
	}

	// Initialize database
	if err := c.initDatabase(ctx, iwasherePath); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	fmt.Printf("Initialized empty iwashere repository in %s\n", iwasherePath)

	// Show next steps
	fmt.Println("\nNext steps:")
	fmt.Println("  iwashere add \"Your first note\"")
	fmt.Println("  iwashere show")

	return nil
}

func (c *InitCommand) updateGitignore(workDir string) error {
	gitignorePath := filepath.Join(workDir, ".gitignore")

	// Check if .gitignore exists
	content := ""
	if _, err := os.Stat(gitignorePath); err == nil {
		// Read existing
		data, err := os.ReadFile(gitignorePath)
		if err != nil {
			return err
		}
		content = string(data)
	}

	// Check if already has .iwashere entry
	if strings.Contains(content, ".iwashere/") {
		return nil // Already there
	}

	// Append entry
	f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if content != "" && !strings.HasSuffix(content, "\n") {
		if _, err := f.WriteString("\n"); err != nil {
			return err
		}
	}

	_, err = f.WriteString("# iwashere personal notes\n.iwashere/\n")
	return err
}

func (c *InitCommand) initDatabase(ctx *Context, iwasherePath string) error {
	// For now, just create an empty file as placeholder
	// Later, this will initialize SQLite schema
	var dbPath string

	switch ctx.Config.Storage.Type {
	case "sqlite":
		dbPath = filepath.Join(iwasherePath, "db", "notes.db")

	case "json":
		jsonRepo := jsonRepo.NewJSONRepository(iwasherePath)
		ctx.Repo = jsonRepo
		return nil
	default:
		dbPath = filepath.Join(iwasherePath, "db", "notes.db")
	}

	f, err := os.Create(dbPath)
	if err != nil {
		return err
	}
	return f.Close()
}

func init() {
	Register("init", NewInitCommand)
}
