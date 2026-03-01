// // internal/commands/config.go
package commands

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"path/filepath"

// 	"iwashere/internal/domain/models"
// 	"iwashere/internal/utils"
// )

// type ConfigCommand struct{}

// func (c *ConfigCommand) Name() string        { return "config" }
// func (c *ConfigCommand) Description() string { return "Manage global iwashere settings" }

// func (c *ConfigCommand) Execute(ctx *Context) error {
// 	configDir := utils.GetConfigDir()
// 	configPath := filepath.Join(configDir, "config.json")

// 	// Ensure config directory exists
// 	if err := os.MkdirAll(configDir, 0755); err != nil {
// 		return fmt.Errorf("failed to create config directory: %w", err)
// 	}

// 	// Handle subcommands
// 	if len(ctx.Args) == 0 {
// 		return c.showConfig(configPath)
// 	}

// 	switch ctx.Args[0] {
// 	case "set":
// 		return c.setConfig(configPath, ctx.Args[1:])
// 	case "get":
// 		return c.getConfig(configPath, ctx.Args[1:])
// 	case "reset":
// 		return c.resetConfig(configPath)
// 	default:
// 		return fmt.Errorf("unknown config subcommand: %s", ctx.Args[0])
// 	}
// }

// func (c *ConfigCommand) showConfig(path string) error {
// 	config, err := loadUserConfig(path)
// 	if err != nil {
// 		config = models.DefaultConfig()
// 	}

// 	data, _ := json.MarshalIndent(config, "", "  ")
// 	fmt.Println(string(data))
// 	return nil
// }

// func loadUserConfig(path string) (*models.Config, error) {
// 	data, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var config models.Config
// 	if err := json.Unmarshal(data, &config); err != nil {
// 		return nil, err
// 	}
// 	return &config, nil
// }
