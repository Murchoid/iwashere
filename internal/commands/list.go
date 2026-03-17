package commands

import (
	"fmt"

	"github.com/Murchoid/iwashere/internal/domain/models"
	"github.com/Murchoid/iwashere/internal/repository"
	"github.com/Murchoid/iwashere/internal/utils"
)

type ListCommand struct {
	spec        *CommandSpec
	baseCommand BaseCommand
}

func NewListCommandFactory() Command {
	return &ListCommand{
		spec: ListCommandSpec,
		baseCommand: BaseCommand{
			NameStr:  "list",
			DescStr:  "Lists all notes in the project/repo",
			UsageStr: "iwashere list [options]",
			ExamplesList: []string{
				"iwashere list",
				"iwashere list --limit 10",
				"iwashere list --limit 10 --short",
				"iwashere list --short",
				"iwashere list --detailed",
				"iwashere list --compact",
			},
		},
	}
}

func (a *ListCommand) Name() string {
	return a.baseCommand.Name()
}

func (a *ListCommand) Description() string {
	return a.baseCommand.Description()
}

func (a *ListCommand) Usage() string {
	return a.baseCommand.Usage()
}

func (a *ListCommand) Examples() []string {
	return a.baseCommand.Examples()
}

func (c *ListCommand) Execute(ctx *Context) error {

	parsedArgs, err := c.spec.Parse(ctx.Args)

	if err != nil {
		utils.PrintCommandHelp(c.Name(), c.Description(), c.Usage(), c.Examples())
		return fmt.Errorf("invalid arguments: %w", err)
	}

	filter := &repository.NoteFilter{
		ProjectPath: ctx.ProjectPath,
	}

	limit := parsedArgs.Flags["limit"]
	pLimit, err := limit.Int()
	if err != nil && limit.Present {
		return err
	}

	if limit.Present {
		filter.Limit = pLimit
	} else {
		filter.Limit = 5
	}

	notes, err := ctx.Repo.ListNotes(filter)
	if err != nil {
		return err
	}

	// Get sessions for grouping
	sessions := make(map[string]*models.Session)
	for _, note := range notes {
		if note.SessionID != "" {
			session, _ := ctx.Repo.GetSession(note.SessionID)
			if session != nil {
				sessions[note.SessionID] = session
			}
		}
	}

	// set up format
	format := "detailed"
	short := parsedArgs.Flags["short"]
	pShort, err := short.Bool()
	if err != nil && short.Present {
		return err
	}

	compact := parsedArgs.Flags["compact"]
	pCompact, err := compact.Bool()
	if err != nil && compact.Present {
		return err
	}

	if pShort {
		format = "short"
	} else if pCompact {
		format = "compact"
	}

	utils.PrintNotes(notes, sessions, format)
	return nil
}

func init() {
	Register("list", NewListCommandFactory)
}
