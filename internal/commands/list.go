package commands

import (
	"fmt"

	"githum.com/Murchoid/iwashere/internal/repository"
	"githum.com/Murchoid/iwashere/internal/utils"
)

type ListCommand struct{}

func NewListCommandFactory() Command {
	return &ListCommand{}
}

func (a *ListCommand) Name() string {
	return "list"
}

func (a *ListCommand) Description() string {
	return "Show a a list of n-number of notes"
}

func (a *ListCommand) Execute(ctx *Context) error {

	repo := ctx.Repo

	filters := repository.NoteFilter{
		ProjectPath: ctx.ProjectPath,
		Limit:       5,
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
		}
	}
	
	return nil
}

func init() {
	Register("list", NewListCommandFactory)
}
