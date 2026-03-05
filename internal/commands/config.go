package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Murchoid/iwashere/internal/domain/models"
	"github.com/Murchoid/iwashere/internal/utils"
)

type ConfigCommand struct {
	BaseCommand
}

func NewConfigCommand() Command {
	return &ConfigCommand{
		BaseCommand{
			NameStr:  "config",
			DescStr:  "Show, set, get or reset your configs",
			UsageStr: "iwashere config [options] [path]",
			ExamplesList: []string{
				"iwashere config",
				"iwashere config set <property.attribute>",
				"iwashere config get <property.attribute>",
				"iwashere config reset",
			},
		},
	}
}

func (c *ConfigCommand) Name() string {
	return c.BaseCommand.Name()
}

func (c *ConfigCommand) Description() string {
	return c.BaseCommand.Description()
}

func (c *ConfigCommand) Examples() []string {
	return c.BaseCommand.Examples()
}

func (c *ConfigCommand) Execute(ctx *Context) error {
	configDir := utils.GetConfigDir()
	configPath := filepath.Join(configDir, "config.json")

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Handle subcommands
	if len(ctx.Args) == 0 {
		return c.showConfig(configPath)
	}

	switch ctx.Args[0] {
	case "set":
		if len(ctx.Args) <= 1 {
			fmt.Println("Set requires key value pair")
			utils.PrintCommandHelp(c.Name(), c.Description(), c.Usage(), c.Examples())
			return nil
		}
		return c.setConfig(configPath, ctx.Args[1:])

	case "get":
		if len(ctx.Args) <= 1 {
			fmt.Println("get requires a key")
			utils.PrintCommandHelp(c.Name(), c.Description(), c.Usage(), c.Examples())
			return nil
		}
		return c.getConfig(configPath, ctx.Args[1:])
	case "reset":
		return c.resetConfig(configPath)
	default:
		return fmt.Errorf("unknown config subcommand: %s", ctx.Args[0])
	}
}

func (c *ConfigCommand) showConfig(path string) error {
	config, err := loadUserConfig(path)
	if err != nil {
		config = models.DefaultConfig()
	}

	data, _ := json.MarshalIndent(config, "", "  ")
	fmt.Println(string(data))
	return nil
}

func (c *ConfigCommand) setConfig(configPath string, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("set requires key and value")
	}

	key := args[0]
	value := args[1]

	if value == "" {
		fmt.Println("Cant set an empty attribute, if you meant to reset the config, use iwashere config reset instead")
		return nil
	}

	// Load existing config
	config, err := loadUserConfig(configPath)
	if err != nil {
		config = models.DefaultConfig()
	}

	// Manual mapping of keys to fields
	switch key {
	case "project.name":
		config.Project.Name = value
	case "storage.type":
		config.Storage.Type = value
	case "storage.path":
		config.Storage.Path = value
	case "git.autocontext":
		// Parse string to bool
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("git.autocontext must be true/false")
		}
		config.Git.AutoContext = boolVal
	case "git.trackbranches":
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("git.trackbranches must be true/false")
		}
		config.Git.TrackBranches = boolVal
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	// Save config
	return c.saveConfig(configPath, config)
}

func (c *ConfigCommand) getConfig(configPath string, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("get requires a key")
	}

	key := args[0]

	config, err := loadUserConfig(configPath)
	if err != nil {
		config = models.DefaultConfig()
	}

	switch key {
	case "project.name":
		fmt.Println(config.Project.Name)
	case "storage.type":
		fmt.Println(config.Storage.Type)
	case "storage.path":
		fmt.Println(config.Storage.Path)
	case "git.autocontext":
		fmt.Println(config.Git.AutoContext)
	case "git.trackbranches":
		fmt.Println(config.Git.TrackBranches)
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	return nil
}

func (c *ConfigCommand) resetConfig(path string) error {
	config := models.DefaultConfig()

	data, _ := json.MarshalIndent(config, "", "  ")
	if err := c.saveConfig(path, config); err != nil {
		fmt.Println("Config reset ....")
	}

	fmt.Println(string(data))
	return nil
}

func (c *ConfigCommand) saveConfig(configPath string, config *models.Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal note: %w", err)
	}

	os.WriteFile(configPath, data, 0644)

	return nil
}

func loadUserConfig(path string) (*models.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config models.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func init() {
	Register("config", NewConfigCommand)
}
