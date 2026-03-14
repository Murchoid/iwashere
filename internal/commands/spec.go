package commands

type FlagType int

const (
	FlagTypeString FlagType = iota
	FlagTypeBool
	FlagTypeInt
	FlagTypeDuration
	FlagTypeTime
)

type FlagSpec struct {
	Name     string
	Type     FlagType
	Usage    string
	Required bool
	Default  any
	Short    string // Short flag like "-m" for "--message"
}

type CommandSpec struct {
	Name        string
	Description string
	Usage       string
	Args        []ArgSpec
	Flags       []FlagSpec
	Subcommands map[string]*CommandSpec
}

type ArgSpec struct {
	Name     string
	Usage    string
	Required bool
	Variadic bool // For commands that take multiple args
}

// Add command spec
var AddCommandSpec = &CommandSpec{
	Name:        "add",
	Description: "Add a new note",
	Usage:       "iwashere add/a <message> [options]",
	Args: []ArgSpec{
		{
			Name:     "message",
			Usage:    "The message of the note",
			Required: true,
		},
	},
	Flags: []FlagSpec{
		{
			Name:     "session",
			Type:     FlagTypeBool,
			Usage:    "adds the note to the current session",
			Required: false,
			Short:    "s",
		},
		{
			Name:     "tags",
			Type:     FlagTypeString,
			Usage:    "tags to attach to the current note",
			Required: false,
			Short:    "t",
		},
		{
			Name:     "branch",
			Type:     FlagTypeString,
			Usage:    "The branch to which to relate the current note with",
			Required: false,
			Short:    "b",
		},
	},
}

// Remind command spec
var RemindCommandSpec = &CommandSpec{
	Name:        "remind",
	Description: "Set reminders for notes",
	Usage:       "iwashere remind [note-id] --at <when> [options]",
	Args: []ArgSpec{
		{
			Name:     "note-id",
			Usage:    "ID of the note to remind about (optional with --message)",
			Required: false,
		},
	},
	Flags: []FlagSpec{
		{
			Name:     "at",
			Type:     FlagTypeTime,
			Usage:    "When to remind (e.g., 'tomorrow 9:00', 'in 2h')",
			Required: true,
			Short:    "a",
		},
		{
			Name:     "message",
			Type:     FlagTypeString,
			Usage:    "Custom reminder message (instead of note message)",
			Required: false,
			Short:    "m",
		},
		{
			Name:     "repeat",
			Type:     FlagTypeString,
			Usage:    "Repeat interval: daily, weekly, monthly",
			Required: false,
			Short:    "r",
		},
	},
	Subcommands: map[string]*CommandSpec{
		"list": {
			Name:        "list",
			Description: "List all reminders",
			Usage:       "iwashere remind list [options]",
			Flags: []FlagSpec{
				{
					Name:     "all",
					Type:     FlagTypeBool,
					Usage:    "Show all reminders (including completed)",
					Required: false,
				},
				{
					Name:     "note",
					Type:     FlagTypeString,
					Usage:    "Filter by note ID",
					Required: false,
				},
			},
		},
		"delete": {
			Name:        "delete",
			Description: "Delete a reminder",
			Usage:       "iwashere remind delete <reminder-id>",
			Args: []ArgSpec{
				{
					Name:     "reminder-id",
					Usage:    "ID of the reminder to delete",
					Required: true,
				},
			},
		},
		"done": {
			Name:        "done",
			Description: "Mark a reminder as done",
			Usage:       "iwashere remind done <reminder-id>",
			Args: []ArgSpec{
				{
					Name:     "reminder-id",
					Usage:    "ID of the reminder to mark as done",
					Required: true,
				},
			},
		},
	},
}

var comandSpecRegestry = map[string]CommandSpec{}

func RegisterSpec(name string, spec *CommandSpec) {
	if _, ok := comandSpecRegestry[name]; ok {
		panic("Command spec already exists")
	}
	comandSpecRegestry[name] = *spec
}

func GetSpec(name string) *CommandSpec {
	if spec, ok := comandSpecRegestry[name]; ok {
		return &spec
	}

	return nil
}

func init() {
	RegisterSpec("remind", RemindCommandSpec)
	RegisterSpec("add", AddCommandSpec)
}
