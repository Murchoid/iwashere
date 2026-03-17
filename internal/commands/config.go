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
	spec        *CommandSpec
	baseCommand BaseCommand
}

func NewConfigCommand() Command {
	return &ConfigCommand{
		spec: ConfigCommandSpec,
		baseCommand: BaseCommand{
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
	return c.baseCommand.Name()
}

func (c *ConfigCommand) Description() string {
	return c.baseCommand.Description()
}

func (c *ConfigCommand) Usage() string {
	return c.baseCommand.Usage()
}

func (c *ConfigCommand) Examples() []string {
	return c.baseCommand.Examples()
}

func (c *ConfigCommand) Execute(ctx *Context) error {
	configDir := filepath.Join(ctx.ProjectPath, ".iwashere", "config")
	configPath := filepath.Join(configDir, "config.json")

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Handle subcommands
	if len(ctx.Args) == 0 {
		return c.showConfig(configPath)
	}

	parsedArgs, err := c.spec.Parse(ctx.Args)

	if err != nil {
		utils.PrintCommandHelp(c.Name(), c.Description(), c.Usage(), c.Examples())
		return fmt.Errorf("invalid arguments: %w", err)
	}

	switch parsedArgs.Subcommand {
	case "set":
		return c.setConfig(configPath, parsedArgs.Positional[0:])

	case "get":
		return c.getConfig(configPath, parsedArgs.Positional[0:])
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
	case "team.name":
		config.Team.TeamName = value
	case "git.auto_context":
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("git.autocontext must be true/false")
		}
		config.Git.AutoContext = boolVal
	case "git.track_branches":
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("git.trackbranches must be true/false")
		}
		config.Git.TrackBranches = boolVal
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	// Save config
	res := c.saveConfig(configPath, config)
	fmt.Printf("Successfully set value '%v' for '%v'\n", value, key)

	return res
}

func (c *ConfigCommand) getConfig(configPath string, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Unknown arguments")
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
	case "git.auto_context":
		fmt.Println(config.Git.AutoContext)
	case "git.track_branches":
		fmt.Println(config.Git.TrackBranches)
	case "team.name":
		fmt.Println(config.Team.TeamName)
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
