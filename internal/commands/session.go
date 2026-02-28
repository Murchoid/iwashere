package commands

import (
	"fmt"
	"time"

	"githum.com/Murchoid/iwashere/internal/domain/models"
	"githum.com/Murchoid/iwashere/internal/repository"
	"githum.com/Murchoid/iwashere/internal/services/git"
	"githum.com/Murchoid/iwashere/internal/utils"
)

type SessionCommand struct{}

func NewSessionCommandFactory() Command {
	return &SessionCommand{}
}

func (a *SessionCommand) Name() string {
	return "add"
}

func (a *SessionCommand) Description() string {
	return "Add a new note"
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

	switch (sessionTags) {
	case "start":
		if ctx.Args[1] == "" {
			return fmt.Errorf("You have to give a session name")
		}

		if err:= startSession(repo, ctx.WorkDir, ctx.Args[1]); err != nil {
			return err
		}
	case "end":
		if err := endSession(repo); err != nil {
			return err
		}

	default:
		return fmt.Errorf("Unknown argument %v", sessionTags)
	}

	return nil
}


func startSession(repo repository.Repository,workDir, sName string) error {
	getRepoService := git.NewService(workDir)

	info, err := getRepoService.GetInfo()

	if err != nil {
		return err
	}

	branch := info.Branch
	session := models.Session {
		ID: utils.GenerateId(),
		Title: sName,
		StartTime: time.Now(),
		Branch: branch,
	}

	if err := repo.SaveSession(&session); err != nil {
		return err
	}

	fmt.Printf("Session %v started at %v", sName, session.StartTime)
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
	fmt.Printf("Session %v ended duration  (%v)", session.Title, howLongAgo)
	return nil
}

func init() {
	Register("session", NewSessionCommandFactory)
}
