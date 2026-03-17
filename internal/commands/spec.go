package commands

type FlagType int

const (
	FlagTypeString FlagType = iota
	FlagTypeBool
	FlagTypeInt
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
	Usage:       "iwashere add <message> [options]",
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

// config spec
var ConfigCommandSpec = &CommandSpec{
	Name:        "config",
	Description: "Show, set, get or reset your configs",
	Usage:       "iwashere config [options] [path]",
	Subcommands: map[string]*CommandSpec{
		"set": {
			Name:        "set",
			Usage:       "set <proptery.attribute>",
			Description: "set a value for a particullar attribute",
			Args: []ArgSpec{
				{
					Name:     "property_attribute",
					Required: true,
				},
			},
		},
		"get": {
			Name:        "get",
			Usage:       "get <proptery.attribute>",
			Description: "get a value for a particullar attribute",
			Args: []ArgSpec{
				{
					Name:     "property_attribute",
					Required: true,
				},
			},
		},
	},
}

// delete spec
var DeleteCommandSpec = &CommandSpec{
	Name:        "delete",
	Description: "deletes a note",
	Usage:       "iwashere delete/rm <id>",
	Args: []ArgSpec{
		{
			Name:     "note-id",
			Usage:    "iwashere delete <note-id>",
			Required: true,
		},
	},
}

// branch command spec
var BranchCommandSpec = &CommandSpec{
	Name:        "branch",
	Description: "list notes from a branch",
	Usage:       "iwashere branch [options]",
	Args: []ArgSpec{
		{
			Name:     "branchName",
			Usage:    "The specific branch to list notes from",
			Required: false,
		},
	},
	Flags: []FlagSpec{
		{
			Name:     "limit",
			Type:     FlagTypeInt,
			Usage:    "Number of notes to show",
			Required: false,
			Short:    "l",
		},
		{
			Name:     "tags",
			Type:     FlagTypeString,
			Usage:    "tags to attach to the current note",
			Required: false,
			Short:    "t",
		},
		{
			Name:     "short",
			Type:     FlagTypeBool,
			Usage:    "Displays notes in 'short' format",
			Required: false,
			Short:    "sh",
		},
		{
			Name:     "detailed",
			Type:     FlagTypeBool,
			Usage:    "Displays notes in 'detailed' format",
			Required: false,
			Short:    "dt",
		},
		{
			Name:     "compact",
			Type:     FlagTypeBool,
			Usage:    "Displays notes in 'compact' format",
			Required: false,
			Short:    "cm",
		},
	},
}

// Edit command spec
var EditCommandSpec = &CommandSpec{
	Name:        "edit",
	Description: "Edits a note",
	Usage:       "iwashere edit <id> --message <message>",
	Args: []ArgSpec{
		{
			Name:     "note-id",
			Usage:    "ID of the note to edit",
			Required: true,
		},
	},
	Flags: []FlagSpec{
		{
			Name:     "message",
			Type:     FlagTypeString,
			Usage:    "the new message to replace a note with",
			Required: false,
			Short:    "m",
		},
		{
			Name:     "tags",
			Type:     FlagTypeString,
			Usage:    "tags to overwrite with",
			Required: false,
			Short:    "t",
		},
		{
			Name:     "add-tags",
			Type:     FlagTypeString,
			Usage:    "tags to add to existing ones",
			Required: false,
			Short:    "at",
		},
		{
			Name:     "remove-tags",
			Type:     FlagTypeString,
			Usage:    "tags to overemove from existing ones",
			Required: false,
			Short:    "rt",
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

// init spec
var InitCommandSpec = &CommandSpec{
	Name:        "init",
	Description: "initialize a iwashere in a project",
	Usage:       "iwashere init [options]",
	Flags: []FlagSpec{
		{
			Name:     "force",
			Type:     FlagTypeBool,
			Usage:    "iwashere init --force",
			Required: false,
			Short:    "f",
		},
		{
			Name:     "no-ignore",
			Type:     FlagTypeBool,
			Usage:    "iwashere init --no-ignore",
			Required: false,
			Short:    "ni",
		},
	},
}

// List command spec
var ListCommandSpec = &CommandSpec{
	Name:        "list",
	Description: "Lists all notes in the project/repo",
	Usage:       "iwashere list [options]",
	Flags: []FlagSpec{
		{
			Name:     "limit",
			Type:     FlagTypeInt,
			Usage:    "iwashere list --limit 10",
			Required: false,
			Short:    "l",
		},
		{
			Name:     "short",
			Type:     FlagTypeBool,
			Usage:    "Display the notes in 'short' format",
			Required: false,
			Short:    "s",
		},
		{
			Name:     "detailed",
			Type:     FlagTypeBool,
			Usage:    "Display the notes in 'detailed' format",
			Required: false,
			Short:    "d",
		},
		{
			Name:     "compact",
			Type:     FlagTypeBool,
			Usage:    "Display the notes in 'compact' format",
			Required: false,
			Short:    "c",
		},
	},
}

// Session command spec
var SessionCommandSpec = &CommandSpec{
	Name:        "session",
	Description: "start, pause, continue or end a session",
	Usage:       "iwashere session [option] <title>/<id>",
	Subcommands: map[string]*CommandSpec{
		"start": {
			Name:        "start",
			Usage:       "iwashere session start <title>",
			Description: "Start session with the title 'title'",
			Args: []ArgSpec{
				{
					Name:     "session title",
					Usage:    "iwashere session start 'example session'",
					Required: true,
				},
			},
		},
		"pause": {
			Name:        "pause",
			Usage:       "iwashere session pause <id>",
			Description: "Pause session with the id 'id'",
		},
		"continue": {
			Name:        "continue",
			Usage:       "iwashere session continue",
			Description: "Start session with the title 'title'",
		},

		"list": {
			Name:        "list",
			Usage:       "iwashere session list [options]",
			Description: "list all sessions",
			Args: []ArgSpec{
				{
					Name:     "session id",
					Usage:    "iwashere session list 123",
					Required: false,
				},
			},
		},

		"end": {
			Name:        "end",
			Usage:       "iwashere end",
			Description: "End an ongoing/continued",
		},
	},
}

// Share command spec
var ShareCommandSpec = &CommandSpec{
	Name:        "share",
	Description: "Share notes with teammates",
	Usage:       "iwashere share [note-id] --with <recipient>",
	Args: []ArgSpec{
		{
			Name:     "note-id",
			Usage:    "ID of the note share",
			Required: false,
		},
	},
	Flags: []FlagSpec{
		{
			Name:     "with",
			Type:     FlagTypeString,
			Usage:    "Email or team to share note with",
			Required: true,
			Short:    "w",
		},
	},
}

// show-shared spec
var ShowSharedCommandSpec = &CommandSpec{
	Name:        "show-shared",
	Description: "Show notes shared with you",
	Usage:       "iwashere show-shared [note-id]",
	Args: []ArgSpec{
		{
			Name:     "note id",
			Usage:    "iwashere show-shared 123",
			Required: false,
		},
	},
}

// show spec
var ShowCommandSpec = &CommandSpec{
	Name:        "show",
	Description: "Shows the note specified by the id",
	Usage:       "iwashere show <id>",
	Args: []ArgSpec{
		{
			Name:     "note id",
			Usage:    "iwashere show 123",
			Required: true,
		},
	},
	Flags: []FlagSpec{
		{
			Name:     "limit",
			Type:     FlagTypeInt,
			Usage:    "iwashere list --limit 10",
			Required: false,
			Short:    "l",
		},
		{
			Name:     "short",
			Type:     FlagTypeBool,
			Usage:    "Display the notes in 'short' format",
			Required: false,
			Short:    "s",
		},
		{
			Name:     "detailed",
			Type:     FlagTypeBool,
			Usage:    "Display the notes in 'detailed' format",
			Required: false,
			Short:    "d",
		},
		{
			Name:     "compact",
			Type:     FlagTypeBool,
			Usage:    "Display the notes in 'compact' format",
			Required: false,
			Short:    "c",
		},
	},
}

// Tag command spec
var TagCommandSpec = &CommandSpec{
	Name:        "tag",
	Description: "add or remove a tag from a note",
	Usage:       "iwashere tag <subcommand> [arguments]",
	Subcommands: map[string]*CommandSpec{
		"add": {
			Name:        "add",
			Usage:       "iwashere tag add <id> <tag>",
			Description: "Add a tag to an existing note",
			Args: []ArgSpec{
				{
					Name:     "note id",
					Usage:    "iwashere tag add 123 <tag>",
					Required: true,
				},
				{
					Name:     "tag",
					Usage:    "iwashere tag add 123 feature",
					Required: true,
				},
			},
		},
		"remove": {
			Name:        "remove",
			Usage:       "iwashere tag remove <id> <tag>",
			Description: "Add a tag to an existing note",
			Args: []ArgSpec{
				{
					Name:     "note id",
					Usage:    "iwashere tag remove 123 <tag>",
					Required: true,
				},
				{
					Name:     "tag",
					Usage:    "iwashere tag remove 123 feature",
					Required: true,
				},
			},
		},
		"list": {
			Name:        "list",
			Usage:       "iwashere tag list <tag>",
			Description: "list all occurences of indicated tag ",
			Args: []ArgSpec{
				{
					Name:     "tag",
					Usage:    "iwashere tag list feature",
					Required: false,
				},
			},
			Flags: []FlagSpec{
				{
					Name:     "cloud",
					Type:     FlagTypeBool,
					Usage:    "iwashere tag list --cloud",
					Required: false,
					Short:    "cl",
				},
				{
					Name:     "limit",
					Type:     FlagTypeInt,
					Usage:    "iwashere list --limit 10",
					Required: false,
					Short:    "l",
				},
				{
					Name:     "short",
					Type:     FlagTypeBool,
					Usage:    "Display the notes in 'short' format",
					Required: false,
					Short:    "s",
				},
				{
					Name:     "detailed",
					Type:     FlagTypeBool,
					Usage:    "Display the notes in 'detailed' format",
					Required: false,
					Short:    "d",
				},
				{
					Name:     "compact",
					Type:     FlagTypeBool,
					Usage:    "Display the notes in 'compact' format",
					Required: false,
					Short:    "c",
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
	RegisterSpec("branch", BranchCommandSpec)
	RegisterSpec("config", ConfigCommandSpec)
	RegisterSpec("delete", DeleteCommandSpec)
	RegisterSpec("edit", EditCommandSpec)
	RegisterSpec("init", InitCommandSpec)
	RegisterSpec("list", ListCommandSpec)
	RegisterSpec("session", SessionCommandSpec)
	RegisterSpec("share", ShareCommandSpec)
	RegisterSpec("share-shared", ShowSharedCommandSpec)
	RegisterSpec("show", ShowCommandSpec)
	RegisterSpec("tag", TagCommandSpec)
}
