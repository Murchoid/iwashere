package commands

import (
	"fmt"
	"slices"

	"github.com/Murchoid/iwashere/internal/domain/models"
	"github.com/Murchoid/iwashere/internal/repository"
	"github.com/Murchoid/iwashere/internal/services/git"
	"github.com/Murchoid/iwashere/internal/utils"
)

type BranchCommand struct {
	spec        *CommandSpec
	baseCommand BaseCommand
}

func NewBranchCommandFactory() Command {
	return &BranchCommand{
		spec: BranchCommandSpec,
		baseCommand: BaseCommand{
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
	return a.baseCommand.Name()
}

func (a *BranchCommand) Description() string {
	return a.baseCommand.Description()
}

func (a *BranchCommand) Usage() string {
	return a.baseCommand.Usage()
}

func (a *BranchCommand) Examples() []string {
	return a.baseCommand.Examples()
}

func (a *BranchCommand) Execute(ctx *Context) error {
	if ctx.Repo == nil {
		return fmt.Errorf("not in an iwashere project (run iwashere init first)")
	}

	repo := ctx.Repo
	parsedArgs, err := a.spec.Parse(ctx.Args)

	if err != nil {
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
		return fmt.Errorf("invalid arguments: %w", err)
	}

	branch := ""
	if len(parsedArgs.Positional) > 0 {
		branch = parsedArgs.Positional[0]
	}

	var filters repository.NoteFilter
	filters.ProjectPath = ctx.ProjectPath
	filters.Limit = a.getLimit(parsedArgs)

	tags := parsedArgs.Flags["tags"]
	pTags, err := tags.String()

	if tags.Present && err == nil {
		filters.Tags = utils.ParseTags(pTags)
	}

	if ctx.Config.Git.AutoContext {
		gitService := git.NewService(ctx.WorkDir)
		gitInfo, err := gitService.GetInfo()

		if err == nil && gitInfo != nil {
			filters.Branch = gitInfo.Branch

			fmt.Printf("Git context: %s @ %s\n", gitInfo.Branch, gitInfo.CommitHash)
			if gitInfo.HasChanges {
				fmt.Printf("You have uncommitted changes\n")
			}

		}

		if branch != "" {
			branchName := branch
			branchIsThere := slices.Contains(gitInfo.Allbranches, branchName)
			if branchIsThere {
				filters.Branch = branchName
			} else {
				fmt.Println()
				return fmt.Errorf("Branch '%s' does not exist in your git", branchName)
			}
		}
	}

	notes, err := repo.ListNotes(&filters)
	if err != nil {
		return err
	}

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
	if ok, err := parsedArgs.Flags["short"].Bool(); ok && err == nil {
		format = "short"
	} else if ok, err := parsedArgs.Flags["compact"].Bool(); ok && err == nil {
		format = "compact"
	}

	utils.PrintNotes(notes, sessions, format)
	return nil
}

func (c *BranchCommand) getLimit(parsedArgs *ParsedArgs) int {
	if num, err := parsedArgs.Flags["limit"].Int(); err == nil {
		return num
	}
	return 5
}

func init() {
	Register("branch", NewBranchCommandFactory)
}
