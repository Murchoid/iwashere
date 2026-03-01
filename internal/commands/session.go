package commands

import (
	"fmt"
	"time"

	"githum.com/Murchoid/iwashere/internal/domain/models"
	"githum.com/Murchoid/iwashere/internal/repository"
	"githum.com/Murchoid/iwashere/internal/services/git"
	"githum.com/Murchoid/iwashere/internal/utils"
)

type SessionCommand struct {
	BaseCommand
}

func NewSessionCommandFactory() Command {
	return &SessionCommand{
		BaseCommand{
			NameStr:  "session",
			DescStr:  "create or end a session",
			UsageStr: "iwashere session [option] [title]",
			ExamplesList: []string{
				"iwashere session start \"start debbuging\"",
				"iwashere session end",
			},
		},
	}
}

func (a *SessionCommand) Name() string {
	return a.BaseCommand.Name()
}

func (a *SessionCommand) Description() string {
	return a.BaseCommand.Description()
}

func (a *SessionCommand) Usage() string {
	return a.BaseCommand.Usage()
}

func (a *SessionCommand) Examples() []string {
	return a.BaseCommand.Examples()
}

func (a *SessionCommand) Execute(ctx *Context) error {
	if ctx.Repo == nil {
		return fmt.Errorf("not in an iwashere project (run iwashere init first)")
	}

	repo := ctx.Repo
	sessionTags := ctx.Args[0]
	if sessionTags == "" {
		return fmt.Errorf("tag name required (use tag or provide as argument)")
	}

	switch sessionTags {
	case "start":
		if ctx.Args[1] == "" {
			return fmt.Errorf("You have to give a session name")
		}

		if err := startSession(repo, ctx.WorkDir, ctx.Args[1]); err != nil {
			return err
		}
	case "end":
		if err := endSession(repo); err != nil {
			return err
		}

	case "list":
		if err := listSessions(repo); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown argument %v\n", sessionTags)
	}

	return nil
}

func startSession(repo repository.Repository, workDir, sName string) error {
	getRepoService := git.NewService(workDir)

	info, err := getRepoService.GetInfo()

	if err != nil {
		return err
	}

	branch := info.Branch
	session := models.Session{
		ID:        utils.GenerateId(),
		Title:     sName,
		StartTime: time.Now(),
		Branch:    branch,
	}

	if err := repo.SaveSession(&session); err != nil {
		return err
	}

	fmt.Printf("Session %v started at %v\n", sName, session.StartTime)
	return nil
}

func endSession(repo repository.Repository) error {
	session, err := repo.GetOpenSession()

	if err != nil {
		return err
	}

	if session.ID == "" {
		fmt.Println("No open sessions")
		return nil
	}

	session.EndTime = time.Now()

	if err := repo.SaveSession(session); err != nil {
		return err
	}

	howLongAgo := utils.HowLongAgo(session.StartTime)
	fmt.Printf("Session %v ended duration  (%v)\n", session.Title, howLongAgo)
	return nil
}

func listSessions(repo repository.Repository) error {
	sessions, err := repo.ListSessions()

	if err != nil {
		return err
	}

	fmt.Println("All sessions")
	for _, session := range sessions {
		howLongAgo := utils.HowLongAgo(session.StartTime)
		fmt.Printf("'%v' started %v\n", session.Title, howLongAgo)
	}

	return nil
}

func init() {
	Register("session", NewSessionCommandFactory)
}
