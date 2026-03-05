package commands

import (
	"fmt"
	"time"

	"github.com/Murchoid/iwashere/internal/domain/models"
	"github.com/Murchoid/iwashere/internal/repository"
	"github.com/Murchoid/iwashere/internal/services/git"
	"github.com/Murchoid/iwashere/internal/utils"
)

type SessionCommand struct {
	BaseCommand
}

func NewSessionCommandFactory() Command {
	return &SessionCommand{
		BaseCommand{
			NameStr:  "session",
			DescStr:  "create or end a session",
			UsageStr: "iwashere session [option] <title>/<id>",
			ExamplesList: []string{
				"iwashere session start \"start debbuging\" #starts a session with the title given",
				"iwashere session list #lists all sessions in the current project",
				"iwashere session list 123 #lists all info about sessoin with id 123",
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
	if len(ctx.Args) == 0 {
		fmt.Println("Option must be provided")
		fmt.Println()
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
		return nil
	}

	sessionTags := ctx.Args[0]
	if sessionTags == "" {
		fmt.Println("Option must be provided")
		fmt.Println()
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
		return nil
	}

	switch sessionTags {
	case "start":

		if len(ctx.Args) <= 1 {
			fmt.Println("Title of session must be provided")
			fmt.Println()
			utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
			return nil
		}

		if ctx.Args[1] == "" {
			fmt.Println("You have to give a session name")
			return nil
		}

		if err := startSession(repo, ctx.WorkDir, ctx.Args[1]); err != nil {
			return err
		}
	case "end":
		if err := endSession(repo); err != nil {
			return err
		}

	case "pause":
		if err := pauseSession(repo); err != nil {
			return err
		}
	case "continue":
		if err := continueSession(repo); err != nil {
			return err
		}
	case "list":
		if len(ctx.Args) > 1 {
			if err := showSession(ctx); err != nil {
				return err
			}
		} else {
			if err := listSessions(repo); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("Unknown argument %v\n", sessionTags)
	}

	return nil
}

func startSession(repo repository.Repository, workDir, sName string) error {

	isThereOngoinSession, err := repo.GetOpenSession()

	if isThereOngoinSession.State == "ongoing" || isThereOngoinSession.ID != "" {
		fmt.Println("There is an ongoing session, end or pause it to start another")
		fmt.Println()
		utils.PrintSessionDetails(isThereOngoinSession, nil)
		return nil
	}


	getRepoService := git.NewService(workDir)

	info, err := getRepoService.GetInfo()

	if err != nil {
		return err
	}

	branch := info.Branch
	session := models.Session{
		ID:        utils.GenerateId(),
		State: "ongoing",
		Title:     sName,
		StartTime: time.Now(),
		Branch:    branch,
	}

	if err := repo.SaveSession(&session); err != nil {
		return err
	}

	fmt.Printf("Session %v started at just now\n", sName)
	return nil
}

func pauseSession(repo repository.Repository) error {
    session, err := repo.GetOpenSession()
    if err != nil {
        return err
    }

    if session.ID == "" {
        fmt.Println("No active session to pause")
        return nil
    }

    // Can only pause ongoing or continued sessions
    if session.State != "ongoing" && session.State != "continued" {
        fmt.Printf("Cannot pause session in state: %s\n", session.State)
        return nil
    }

    // Calculate duration since last start/continue
    now := time.Now()
    session.TotalTime += models.Duration(now.Sub(session.StartTime))
    session.EndTime = now
    session.State = "paused"

    if err := repo.SaveSession(session); err != nil {
        return err
    }

    fmt.Printf("Session '%s' paused (total: %s)\n", 
        session.Title, 
        session.TotalTime.Duration().Round(time.Second))
    return nil
}

func continueSession(repo repository.Repository) error {
    session, err := repo.GetOpenSession()
    if err != nil {
        return err
    }

    if session.ID == "" {
        fmt.Println("No paused session to continue")
        return nil
    }

    // Can only continue paused sessions
    if session.State != "paused" {
        fmt.Printf("Cannot continue session in state: %s\n", session.State)
        return nil
    }

    // Restart the session
    session.StartTime = time.Now()
    session.State = "continued"

    if err := repo.SaveSession(session); err != nil {
        return err
    }

    fmt.Printf("Session '%s' continued\n", session.Title)
    return nil
}

func endSession(repo repository.Repository) error {
    session, err := repo.GetOpenSession()
    if err != nil {
        return err
    }

    if session.ID == "" {
        fmt.Println("No active session to end")
        return nil
    }

    now := time.Now()
    
    // Add final segment if session was active
    if session.State == "ongoing" || session.State == "continued" {
        session.TotalTime += models.Duration(now.Sub(session.StartTime))
    }
    
    session.EndTime = now
    session.State = "ended"

    if err := repo.SaveSession(session); err != nil {
        return err
    }

    fmt.Printf("Session '%s' ended (total: %s)\n", 
        session.Title, 
        session.TotalTime.Duration().Round(time.Second))
    return nil
}

func listSessions(repo repository.Repository) error {
	sessions, err := repo.ListSessions()
	if err != nil {
		return err
	}

	// Use the new display function
	utils.PrintSessions(sessions, false, nil)
	return nil
}

func showSession(ctx *Context) error {
	if len(ctx.Args) < 2 {
		// Show current active session
		session, err := ctx.Repo.GetOpenSession()
		if err != nil {
			return fmt.Errorf("no active session and no session ID provided")
		}

		notes, err := ctx.Repo.GetNotesBySession(session.ID)
		if err != nil {
			return err
		}

		utils.PrintCurrentSession(session, notes)
		return nil
	}

	// Show specific session by ID
	sessionID := ctx.Args[1]
	session, err := ctx.Repo.GetSession(sessionID)
	if err != nil {
		return err
	}

	notes, err := ctx.Repo.GetNotesBySession(session.ID)
	if err != nil {
		return err
	}

	utils.PrintSessionDetails(session, notes)
	return nil
}

func init() {
	Register("session", NewSessionCommandFactory)
}
