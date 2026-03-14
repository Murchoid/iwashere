// internal/commands/help.go
package commands

import (
	"fmt"
	"sort"

	"github.com/Murchoid/iwashere/internal/utils"
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
		showCommandHelp(cmdName)
	}

	// Otherwise show global help
	ShowGlobalHelp()

	return nil
}

func ShowGlobalHelp() {
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

}

func showCommandHelp(cmdName string) {
	spec := getCommandSpec(cmdName)
	if spec == nil {
		return
	}

	fmt.Printf("Usage: %s\n", spec.Usage)
	fmt.Printf("\nDescription:\n  %s\n", spec.Description)

	if len(spec.Args) > 0 {
		fmt.Printf("\nArguments:\n")
		for _, arg := range spec.Args {
			req := ""
			if arg.Required {
				req = " (required)"
			}
			fmt.Printf("  %s%s\n    %s\n", arg.Name, req, arg.Usage)
		}
	}

	if len(spec.Flags) > 0 {
		fmt.Printf("\nFlags:\n")
		for _, flag := range spec.Flags {
			short := ""
			if flag.Short != "" {
				short = fmt.Sprintf(" -%s", flag.Short)
			}
			req := ""
			if flag.Required {
				req = " (required)"
			}
			def := ""
			if flag.Default != nil {
				def = fmt.Sprintf(" (default: %v)", flag.Default)
			}
			fmt.Printf("  --%s%s%s\n    %s%s\n",
				flag.Name, short, req, flag.Usage, def)
		}
	}

	if len(spec.Subcommands) > 0 {
		fmt.Printf("\nSubcommands:\n")
		for name, sub := range spec.Subcommands {
			fmt.Printf("  %s\n    %s\n", name, sub.Description)
		}
	}
}

func getCommandSpec(name string) *CommandSpec {
	return GetSpec(name)
}

func init() {
	Register("help", NewHelpCommandFactory)
}
