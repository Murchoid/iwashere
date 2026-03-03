package commands

import (
	"strconv"

	"githum.com/Murchoid/iwashere/internal/domain/models"
	"githum.com/Murchoid/iwashere/internal/repository"
	"githum.com/Murchoid/iwashere/internal/utils"
)

type ListCommand struct {
	BaseCommand
}

func NewListCommandFactory() Command {
	return &ListCommand{
		BaseCommand{
			NameStr:  "list",
			DescStr:  "Lists all notes in the project/repo",
			UsageStr: "iwashere list/ls [options]",
			ExamplesList: []string{
				"iwashere list",
				"iwashere list --limit 10",
				"iwashere ls",
				"iwashere ls --limit 10",
			},
		},
	}
}

func (a *ListCommand) Name() string {
	return a.BaseCommand.Name()
}

func (a *ListCommand) Description() string {
	return a.BaseCommand.Description()
}

func (a *ListCommand) Usage() string {
	return a.BaseCommand.Usage()
}

func (a *ListCommand) Examples() []string {
	return a.BaseCommand.Examples()
}

func (c *ListCommand) Execute(ctx *Context) error {
	filter := &repository.NoteFilter{
		ProjectPath: ctx.ProjectPath,
		Limit:       c.getLimit(ctx),
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

	// Use display package
	format := "detailed"
	if ctx.Flags["--short"] == "true" {
		format = "short"
	}

	utils.PrintNotes(notes, sessions, format)
	return nil
}


func (c *ListCommand) getLimit(ctx *Context) int {
	if ctx.Flags["--limit"] != "" {
		limit, err := strconv.Atoi(ctx.Flags["--limit"])
		if err != nil {
			return 5
		}
		return limit
	}
	return 5
}

func init() {
	Register("list", NewListCommandFactory)
}
