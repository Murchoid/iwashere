package commands

import (
	"fmt"

	"githum.com/Murchoid/iwashere/internal/domain/models"
	"githum.com/Murchoid/iwashere/internal/repository"
	"githum.com/Murchoid/iwashere/internal/services/git"
	"githum.com/Murchoid/iwashere/internal/utils"
)

type AddCommand struct {
	BaseCommand
}

func NewAddCommandFactory() Command {
	return &AddCommand{
		BaseCommand{
			NameStr:  "add",
			DescStr:  "Add a new note",
			UsageStr: "iwashere add/a <message> [options]",
			ExamplesList: []string{
				"iwashere add \"Working on authentication\"",
				"iwashere a \"Fix memory leak\" --tags bug,performance",
				"iwashere add \"Update README\" --branch main",
				"iwashere add \"Add current note in current session\" --session",
			},
		},
	}
}

func (a *AddCommand) Name() string {
	return a.BaseCommand.Name()
}

func (a *AddCommand) Description() string {
	return a.BaseCommand.Description()
}

func (a *AddCommand) Usage() string {
	return a.BaseCommand.Usage()
}

func (a *AddCommand) Examples() []string {
	return a.BaseCommand.Examples()
}

func (a *AddCommand) Execute(ctx *Context) error {
	if ctx.Repo == nil {
		return fmt.Errorf("not in an iwashere project (run iwashere init first)")
	}

	repo := ctx.Repo
	message := ctx.Args[0]
	if message == "" {
		return fmt.Errorf("message required (use -m or provide as argument)")
	}

	note := &models.Note{
		Message:     message,
		ProjectPath: ctx.ProjectPath,
	}

	if ctx.Flags["--tags"] != "" {
		note.Tags = utils.ParseTags(ctx.Flags["--tags"])
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

	if err := repo.SaveNote(note); err != nil {
		return fmt.Errorf("failed to save note: %w", err)
	}

	if ctx.Flags["--session"] != "" {
		if err := addNoteToCurrentSession(repo, note); err != nil {
			return err
		}
	}

	fmt.Printf("Note saved (ID: %s)\n", note.ID)
	return nil
}

func addNoteToCurrentSession(repo repository.Repository, note *models.Note) error {
	session, err := repo.GetOpenSession()

	if err != nil {
		return err
	}

	if session.ID == "" {
		fmt.Println("No session open create one to add a note to it")
		return nil
	}

	if session.Branch != note.Branch {
		fmt.Println("There is no open session in this branch")
		return nil
	}

	session.Notes = append(session.Notes, note.ID)
	note.SessionID = session.ID
	if err := repo.SaveSession(session); err != nil {
		return err
	}

	return nil
}

func init() {
	Register("add", NewAddCommandFactory)
}
