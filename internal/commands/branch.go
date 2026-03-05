package commands

import (
	"fmt"
	"strconv"

	"github.com/Murchoid/iwashere/internal/domain/models"
	"github.com/Murchoid/iwashere/internal/repository"
	"github.com/Murchoid/iwashere/internal/services/git"
	"github.com/Murchoid/iwashere/internal/utils"
)

type BranchCommand struct {
	BaseCommand
}

func NewBranchCommandFactory() Command {
	return &BranchCommand{
		BaseCommand{
			NameStr:  "branch",
			DescStr:  "Shows all notes of the current branch",
			UsageStr: `iwashere branch  [argument]`,
			ExamplesList: []string{
				"iwashere branch",
				"iwashere branch main",
			},
		},
	}
}

func (a *BranchCommand) Name() string {
	return a.BaseCommand.Name()
}

func (a *BranchCommand) Description() string {
	return a.BaseCommand.Description()
}

func (a *BranchCommand) Usage() string {
	return a.BaseCommand.Usage()
}

func (a *BranchCommand) Examples() []string {
	return a.BaseCommand.Examples()
}

func (a *BranchCommand) Execute(ctx *Context) error {
	if ctx.Repo == nil {
		return fmt.Errorf("not in an iwashere project (run iwashere init first)")
	}

	repo := ctx.Repo
	branch := ""
	if len(ctx.Args) > 0 {
		branch = ctx.Args[0]

	}

	var filters repository.NoteFilter
	filters.ProjectPath = ctx.ProjectPath
	filters.Limit = a.getLimit(ctx)

	if ctx.Flags["--tags"] != "" {
		filters.Tags = utils.ParseTags(ctx.Flags["--tags"])
	}

	if ctx.Config.Git.AutoContext {
		gitService := git.NewService(ctx.WorkDir)
		if gitInfo, err := gitService.GetInfo(); err == nil && gitInfo != nil {
			filters.Branch = gitInfo.Branch

			fmt.Printf("Git context: %s @ %s\n", gitInfo.Branch, gitInfo.CommitHash)
			if gitInfo.HasChanges {
				fmt.Printf("You have uncommitted changes\n")
			}
		}
		if branch != "" {
			filters.Branch = branch
		}
	}

	notes, err := repo.ListNotes(&filters)
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

	format := "detailed"
	if ctx.Flags["--short"] == "true" {
		format = "short"
	}

	utils.PrintNotes(notes, sessions, format)
	return nil
}

func (c *BranchCommand) getLimit(ctx *Context) int {
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
	Register("branch", NewBranchCommandFactory)
}
