package commands

import (
	"fmt"

	"githum.com/Murchoid/iwashere/internal/repository"
	"githum.com/Murchoid/iwashere/internal/services/git"
	"githum.com/Murchoid/iwashere/internal/utils"
)

type BranchCommand struct{}

func NewBranchCommandFactory() Command {
	return &BranchCommand{}
}

func (a *BranchCommand) Name() string {
	return "add"
}

func (a *BranchCommand) Description() string {
	return "Add a new note"
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
	filters.Tags = utils.ParseTags(ctx.Flags["--tags"])
	
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

	for idx := range notes {
		howLongAgo := utils.HowLongAgo(notes[idx].UpdatedAt)
		fmt.Printf("[%v](%v) %v: %v\n", howLongAgo, notes[idx].Branch, notes[idx].ID, notes[idx].Message)
				if len(notes[idx].ModifiedFiles) > 0 {
			fmt.Println("Modified files")
			for mIdx := range notes[idx].ModifiedFiles {
				fmt.Printf("[%v]\n", notes[idx].ModifiedFiles[mIdx])
			}
			fmt.Println()
			fmt.Println()
		}


	}
	return nil
}

func (a *BranchCommand) parseTags(tagFlags string) []string {
	var tags []string
	tags = append(tags, tagFlags)

	return tags
}

func init() {
	Register("branch", NewBranchCommandFactory)
}
