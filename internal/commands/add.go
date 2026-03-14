package commands

import (
	"fmt"

	"github.com/Murchoid/iwashere/internal/domain/models"
	"github.com/Murchoid/iwashere/internal/repository"
	"github.com/Murchoid/iwashere/internal/services/git"
	"github.com/Murchoid/iwashere/internal/utils"
)

type AddCommand struct {
	spec *CommandSpec
}

func NewAddCommandFactory() Command {
	return &AddCommand{
		spec: AddCommandSpec,
	}
}

func (a *AddCommand) Name() string {
	return "add"
}

func (a *AddCommand) Description() string {
	return "Add a new note"
}

func (a *AddCommand) Usage() string {
	return "iwashere add/a <message> [options]"
}

func (a *AddCommand) Examples() []string {
	return []string{
		"iwashere add \"Working on authentication\"",
		"iwashere a \"Fix memory leak\" --tags bug,performance",
		"iwashere add \"Update README\" --branch main",
		"iwashere add \"Add current note in current session\" --session",
	}
}

func (a *AddCommand) Execute(ctx *Context) error {
	if ctx.Repo == nil {
		return fmt.Errorf("not in an iwashere project (run iwashere init first)")
	}

	repo := ctx.Repo
	parsed, err := a.spec.Parse(ctx.Args)
	if err != nil {
		// Show help on parse error
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
		return fmt.Errorf("invalid arguments: %w", err)
	}

	if len(parsed.Positional) == 0 {
		fmt.Println("Note message must be provided")
		fmt.Println()
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())

		return nil
	}

	if len(parsed.Positional) > 1 {
		fmt.Println("add only accepts one argument")
		fmt.Println()
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
		return nil
	}

	message := parsed.Positional[0]
	if message == "" {
		fmt.Println("message Cannot be empty")
		return nil
	}

	note := &models.PrivateNote{
		Message:     message,
		ProjectPath: ctx.ProjectPath,
	}

	tags := parsed.Flags["tags"]
	pTag, err := tags.String()
	if err == nil && tags.Present {
		note.Tags = utils.ParseTags(pTag)
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

			branch, err := parsed.Flags["branch"].String()
			if branch != "" && err == nil && parsed.Flags["branch"].Present {
				branchName := branch
				branchIsThere := false
				fmt.Println(gitInfo.Allbranches)
				for _, branch := range gitInfo.Allbranches {
					if branch != branchName {
						branchIsThere = true
					}
				}

				if branchIsThere {
					note.Branch = branchName
				} else {
					return fmt.Errorf("Branch %s does not exist in your git", branchName)
				}
			}

			fmt.Printf("Git context: %s @ %s\n", gitInfo.Branch, gitInfo.CommitHash)
			if gitInfo.HasChanges {
				fmt.Printf("You have uncommitted changes\n")
			}
		}
	}

	session, err := parsed.Flags["session"].Bool()
	if session && err == nil && parsed.Flags["session"].Present {
		if err := addNoteToCurrentSession(repo, note); err != nil {
			return err
		}
	}

	if err := repo.SaveNote(note); err != nil {
		return fmt.Errorf("failed to save note: %w", err)
	}

	fmt.Printf("Note saved (ID: %s)\n", note.ID)
	return nil
}

func addNoteToCurrentSession(repo repository.Repository, note *models.PrivateNote) error {
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
