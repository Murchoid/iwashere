// internal/commands/help.go
package commands

import (
	"fmt"
	"sort"

	"githum.com/Murchoid/iwashere/internal/utils"
)

type HelpCommand struct{}

func NewHelpCommandFactory() Command {
	return &HelpCommand{}
}

func (c *HelpCommand) Name() string {
	return "help"
}

func (c *HelpCommand) Description() string {
	return "Show help information for commands"
}

func (c *HelpCommand) Usage() string {
	return "iwashere help [command]"
}

func (c *HelpCommand) Examples() []string {
	return []string{
		"iwashere help",
		"iwashere help add",
		"iwashere help session",
	}
}

func (c *HelpCommand) Execute(ctx *Context) error {
	// If they specified a command, show help for that command
	if len(ctx.Args) > 0 {
		cmdName := ctx.Args[0]
		return showCommandHelp(cmdName)
	}

	// Otherwise show global help
	return ShowGlobalHelp()
}

func ShowGlobalHelp() error {
	fmt.Println("iwashere - Context preservation tool")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("Usage: iwashere <command> [arguments]")
	fmt.Println()

	// Get all commands and sort them alphabetically
	commands := GetAll()
	var names []string
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)

	// Find longest command name for pretty formatting
	maxLen := 0
	for _, name := range names {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	// Print each command with aligned descriptions
	for _, name := range names {
		factory := commands[name]
		cmd := factory()
		fmt.Println("Command:")
		utils.PrintCommandHelp(cmd.Name(), cmd.Description(), cmd.Usage(), cmd.Examples())
		fmt.Println()
	}

	fmt.Println()
	fmt.Println("Use 'iwashere help <command>' for more details about a specific command.")
	fmt.Println("Examples: iwashere help add, iwashere help session")

	return nil
}

func showCommandHelp(cmdName string) error {
	factory, exists := GetFactory(cmdName)
	if !exists {
		return fmt.Errorf("unknown command: %s", cmdName)
	}
	cmd := factory()
	utils.PrintCommandHelp(cmd.Name(), cmd.Description(), cmd.Usage(), cmd.Examples())
	fmt.Println()

	// Could add more sections here like:
	// - Options/Flags
	// - See Also
	// - Notes

	return nil
}

func init() {
	Register("help", NewHelpCommandFactory)
}
