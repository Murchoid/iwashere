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
	spec        *CommandSpec
	baseCommand BaseCommand
}

func NewSessionCommandFactory() Command {
	return &SessionCommand{
		spec: SessionCommandSpec,
		baseCommand: BaseCommand{
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
	return a.baseCommand.Name()
}

func (a *SessionCommand) Description() string {
	return a.baseCommand.Description()
}

func (a *SessionCommand) Usage() string {
	return a.baseCommand.Usage()
}

func (a *SessionCommand) Examples() []string {
	return a.baseCommand.Examples()
}

func (a *SessionCommand) Execute(ctx *Context) error {
	if ctx.Repo == nil {
		return fmt.Errorf("not in an iwashere project (run iwashere init first)")
	}

	repo := ctx.Repo

	parsedArgs, err := a.spec.Parse(ctx.Args)

	if err != nil {
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
		return fmt.Errorf("invalid arguments: %w", err)
	}

	if parsedArgs.Subcommand == "" {
		fmt.Println("Option must be provided")
		fmt.Println()
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
		return nil
	}

	sessionTags := parsedArgs.Subcommand

	switch sessionTags {
	case "start":
		if len(parsedArgs.Positional) == 0 {
			fmt.Println("Session title must be provided")
			fmt.Println()
			utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
			return nil
		}
		title := parsedArgs.Positional[0]
		if err := startSession(repo, ctx.WorkDir, title); err != nil {
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
		var id string
		if len(parsedArgs.Positional) > 0 {
			id = parsedArgs.Positional[0]
		}

		if id != "" {
			if err := showSession(repo, id); err != nil {
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

	if isThereOngoinSession != nil {
		if isThereOngoinSession.State == models.Ongoing || isThereOngoinSession.ID != "" && isThereOngoinSession.EndTime.IsZero() {
			fmt.Println("There is an ongoing session, end it to start another")
			fmt.Println()
			utils.PrintSessionDetails(isThereOngoinSession, nil)
			return nil
		}
	}

	getRepoService := git.NewService(workDir)

	info, err := getRepoService.GetInfo()

	if err != nil {
		return err
	}
	branch := ""

	if info != nil {
		branch = info.Branch
	}
	session := models.Session{
		ID:        utils.GenerateId(),
		State:     models.Ongoing,
		Title:     sName,
		StartTime: time.Now(),
		Branch:    branch,
	}

	if err := repo.SaveSession(&session); err != nil {
		return err
	}

	fmt.Printf("Session '%v' started has been started now\n", sName)
	return nil
}

func pauseSession(repo repository.Repository) error {
	session, err := repo.GetOpenSession()
	if err != nil {
		return err
	}

	if session == nil {
		fmt.Println("No active session to pause")
		return nil
	}

	// Can only pause ongoing or continued sessions
	if session.State != models.Ongoing && session.State != models.Continued {
		fmt.Printf("Cannot pause session in state: %s\n", session.State)
		return nil
	}

	// Calculate duration since last start/continue
	now := time.Now()
	session.TotalTime += models.Duration(now.Sub(session.StartTime))
	session.EndTime = now
	session.State = models.Paused

	if err := repo.SaveSession(session); err != nil {
		return err
	}

	fmt.Printf("Session %s'%s'%s paused (total: %s)\n",
		utils.ColorPurple,
		session.Title,
		utils.ColorReset,
		session.TotalTime.Duration().Round(time.Second))
	return nil
}

func continueSession(repo repository.Repository) error {
	session, err := repo.GetOpenSession()
	if err != nil {
		return err
	}

	if session == nil {
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
	session.State = models.Continued

	if err := repo.SaveSession(session); err != nil {
		return err
	}

	fmt.Printf("Session %s'%s'%s continued\n", utils.ColorPurple, session.Title, utils.ColorReset)
	return nil
}

func endSession(repo repository.Repository) error {
	session, err := repo.GetOpenSession()
	if err != nil {
		return err
	}

	if session == nil {
		fmt.Println("No active session to end")
		return nil
	}

	now := time.Now()

	// Add final segment if session was active
	if session.State == models.Ongoing || session.State == models.Continued {
		session.TotalTime += models.Duration(now.Sub(session.StartTime))
	}

	session.EndTime = now
	session.State = models.Ended

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

func showSession(repo repository.Repository, sessionID string) error {

	// Show specific session by ID
	session, err := repo.GetSession(sessionID)
	if err != nil {
		return err
	}

	notes, err := repo.GetNotesBySession(session.ID)
	if err != nil {
		return err
	}

	utils.PrintSessionDetails(session, notes)
	return nil
}

func init() {
	Register("session", NewSessionCommandFactory)
}
