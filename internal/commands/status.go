package commands

import (
	"fmt"

	"githum.com/Murchoid/iwashere/internal/domain/models"
	"githum.com/Murchoid/iwashere/internal/utils"
)

type StatusCommand struct {
	BaseCommand
}

func NewStatusCommandFactory() Command {
	return &StatusCommand{
		BaseCommand{
			NameStr:  "status",
			DescStr:  "Quick status of what you were doing",
			UsageStr: "iwashere status ",
			ExamplesList: []string{
				"iwashere status",
			},
		},
	}
}

func (a *StatusCommand) Name() string {
	return a.BaseCommand.Name()
}

func (a *StatusCommand) Description() string {
	return a.BaseCommand.Description()
}

func (a *StatusCommand) Usage() string {
	return a.BaseCommand.Usage()
}

func (a *StatusCommand) Examples() []string {
	return a.BaseCommand.Examples()
}

func (a *StatusCommand) Execute(ctx *Context) error {
	if ctx.Repo == nil {
		return fmt.Errorf("not in an iwashere project (run iwashere init first)")
	}

	repo := ctx.Repo
	if len(ctx.Args) > 0 {
		fmt.Println("Unrecognized arguments")
		fmt.Println()
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
		return nil
	}

	session, err := repo.GetOpenSession()

	if err != nil {
		return err
	}

	var notes []models.Note

	for _,noteId := range session.Notes {
		fNotes, err := repo.GetNote(noteId)

		if err != nil {
			fmt.Println("Failed to fetc note of id: ", noteId)
			continue
		}

		notes = append(notes, *fNotes)
	}

	printResults(session, notes)
	return nil
}

func printResults(session *models.Session, notes []models.Note) {

	if session.ID == "" {
		fmt.Println("You dont ave any active sessions, create one using \niwashere session start")
	} else {
		fmt.Printf("You were working on session: %v (%v)\n", session.Title, utils.HowLongAgo(session.StartTime))
	}

	
	if len(notes)>0 {
		fmt.Println("Last note: ", notes[len(notes)-1])
		fmt.Println("Modified files: ")
		var modifiedFiles map[string]int
		modifiedFiles = make(map[string]int)
	
		for _, note := range notes {
			if len(note.ModifiedFiles) > 0 {
				for _, file := range note.ModifiedFiles {
				modifiedFiles[file]++
				}
			}
		}
	
		for files := range modifiedFiles {
			fmt.Println(files)
		}
	
		fmt.Println("Relted notes:")
		for _, note := range notes {
			if notes[len(notes)-1].ID != note.ID {
				fmt.Println(note.Message)
			}
		}
	}


}

func init() {
	Register("status", NewStatusCommandFactory)
}
