package commands

import (
	"fmt"

	"githum.com/Murchoid/iwashere/internal/domain/models"
	"githum.com/Murchoid/iwashere/internal/services/git"
)

type AddCommand struct{}

func NewAddCommandFactory() Command {
	return &AddCommand{}
}

func (a *AddCommand) Name() string {
	return "add"
}

func (a *AddCommand) Description() string {
	return "Add a new note"
}

func (a *AddCommand) Execute(ctx *Context) error {
	if ctx.Repo == nil {
		return fmt.Errorf("not in an iwashere project (run iwashere init first)")
	}

	message := ctx.Args[0]
	if message == "" {
		return fmt.Errorf("message required (use -m or provide as argument)")
	}

	note := &models.Note{
		Message:     message,
		ProjectPath: ctx.ProjectPath,
		Tags:        a.parseTags(ctx.Flags["--tags"]),
	}

	if ctx.Config.Git.AutoContext {
		gitService := git.NewService(ctx.WorkDir)
		if gitInfo, err := gitService.GetInfo(); err == nil && gitInfo != nil {
			note.Branch = gitInfo.Branch
			note.CommitHash = gitInfo.CommitHash
			note.CommitMsg = gitInfo.CommitMsg
			note.Remote = gitInfo.Remote

			if mFiles, err := gitService.GetModifiedFiles(); err == nil && mFiles != nil {
				note.ModifiedFiles = mFiles
			}

			fmt.Printf("Git context: %s @ %s\n", gitInfo.Branch, gitInfo.CommitHash)
			if gitInfo.HasChanges {
				fmt.Printf("You have uncommitted changes\n")
			}
		}
	}

	if err := ctx.Repo.SaveNote(note); err != nil {
		return fmt.Errorf("failed to save note: %w", err)
	}

	fmt.Printf("Note saved (ID: %s)\n", note.ID)
	return nil
}

func (a *AddCommand) parseTags(tagFlags string) []string {
	var tags []string
	tags = append(tags, tagFlags)

	return tags
}

func init() {
	Register("add", NewAddCommandFactory)
}
