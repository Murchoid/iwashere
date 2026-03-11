package argparser

import "strings"

type Arguments struct {
	Name  string
	Args  []string
	Flags map[string]string
}

func NewArgParser() Arguments {
	return Arguments{
		Flags: make(map[string]string),
	}
}

func (a *Arguments) ParseArguments(args []string) Arguments {
	var arguments Arguments

	arguments.Name = a.parseCmdName(args[0])
	arguments.Args, arguments.Flags = a.parseFlags(args[1:])

	return arguments
}

func (a *Arguments) parseCmdName(argument string) string {
	switch argument {
	case "a":
		return "add"
	case "s":
		return "show"
	case "ls":
		return "list"
	case "rm":
		return "delete"
	case "b":
		return "branch"
	case "-h":
		return "help"
	case "--help":
		return "help"
	case "-v":
		return "version"
	case "--version":
		return "version"
	}

	return argument
}

func (a *Arguments) parseFlags(args []string) ([]string, map[string]string) {

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
