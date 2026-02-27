package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"githum.com/Murchoid/iwashere/internal/commands"
	"githum.com/Murchoid/iwashere/internal/domain/models"
	"githum.com/Murchoid/iwashere/internal/repository"
	"githum.com/Murchoid/iwashere/internal/repository/jsonRepo"
)

func main() {

	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	// Parse command
	cmdName := normalizeCmdName(args[0])

	cmdArgs, flags := parseFlags(args[1:])

	// Get working directory
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
		os.Exit(1)
	}

	// Check if we're in an iwashere project (look for .iwashere/)
	projectPath := findProjectRoot(workDir)

	// Load config (if exists)
	config, err := loadConfig(projectPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load config: %v\n", err)
		fmt.Fprintf(os.Stderr, "Using default configuration\n")
		config = models.DefaultConfig()
	}

	repo, err := createRepository(projectPath, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating repository: %v\n", err)
		os.Exit(1)
	}
	// Create command context
	ctx := &commands.Context{
		Args:        cmdArgs,
		Flags:       flags,
		WorkDir:     workDir,
		ProjectPath: projectPath,
		Config:      config,
		Repo:        repo,
	}

	// Find and execute command
	if cmd, exists := commands.GetFactory(cmdName); exists {
		command := cmd()

		if err := command.Execute(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Unknown command: %s\n", cmdName)
		printUsage()
		os.Exit(1)
	}
}

func normalizeCmdName(cmdName string) string {
	switch cmdName {
	case "a":
		return "add"
	case "s":
		return "show"
	case "ls":
		return "list"
	case "rm":
		return "delete"
	}

	return cmdName
}
func parseFlags(args []string) ([]string, map[string]string) {
	var positional []string
	flags := make(map[string]string)

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			// Handle flags
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				flags[arg] = args[i+1]
				i++ // Skip next arg
			} else {
				flags[arg] = "true" // boolean flag
			}
		} else {
			positional = append(positional, arg)
		}
	}

	return positional, flags
}

func findProjectRoot(path string) string {
	for {
		if _, err := os.Stat(filepath.Join(path, ".iwashere")); err == nil {
			return path
		}
		parent := filepath.Dir(path)
		if parent == path { // Reached root
			break
		}
		path = parent
	}
	return "" // No .iwashere found
}

func loadConfig(projectPath string) (*models.Config, error) {
	// If no project path, return default config (no file)
	if projectPath == "" {
		return models.DefaultConfig(), nil
	}

	// Define config path (using lowercase 'config' for consistency)
	configDir := filepath.Join(projectPath, ".iwashere", "config")
	configPath := filepath.Join(configDir, "config.json")

	// Check if config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Config doesn't exist - create it
		return createDefaultConfig(configDir, configPath)
	} else if err != nil {
		// Some other error (permission, etc)
		return nil, fmt.Errorf("error checking config: %w", err)
	}

	// Config exists - read it
	return readConfig(configPath)
}

func createDefaultConfig(configDir, configPath string) (*models.Config, error) {
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Get default config
	config := models.DefaultConfig()

	// Marshal with nice formatting
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write config file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write config: %w", err)
	}

	return config, nil
}

func readConfig(configPath string) (*models.Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config models.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

func createRepository(projectPath string, cfg *models.Config) (repository.Repository, error) {
	if projectPath == "" {
		return nil, nil // No project, can't use repo
	}

	iwasherePath := filepath.Join(projectPath, ".iwashere")

	switch cfg.Storage.Type {
	case "json":
		return jsonRepo.NewJSONRepository(iwasherePath), nil
	case "sqlite":
		// return sqlite.NewSQLiteRepository(iwasherePath), nil
		return nil, fmt.Errorf("sqlite not implemented yet")
	default:
		return jsonRepo.NewJSONRepository(iwasherePath), nil // default to json
	}
}

func printUsage() {
	fmt.Println(`iwashere - Context preservation tool

Usage:
  iwashere <command> [arguments]

Commands:
  init    Initialize iwashere in current directory
  add     Add a note (coming soon)
  show    Show notes (coming soon)
  
Use "iwashere help <command>" for more info about a command.`)
}
