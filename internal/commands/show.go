package commands

import (
	"fmt"

	"githum.com/Murchoid/iwashere/internal/utils"
)

type ShowCommand struct{}

func NewShowCommandFactory() Command {
	return &ShowCommand{}
}

func (a *ShowCommand) Name() string {
	return "show"
}

func (a *ShowCommand) Description() string {
	return "Show a exisiting notes"
}

func (a *ShowCommand) Execute(ctx *Context) error {

	repo := ctx.Repo
	id := ctx.Args[0]

	note, err := repo.GetNote(id)
	if err != nil {
		return err
	}
	howLongAgo := utils.HowLongAgo(note.UpdatedAt)
	fmt.Printf("[%v](%v) %v: %v\n", howLongAgo, note.Branch, note.ID, note.Message)
	if len(note.Tags) > 0 {
		fmt.Printf("Tags:[")
		for idx := range note.Tags {
			fmt.Printf("%v", note.Tags[idx])
		}
		fmt.Printf("]\n")
	}
	if len(note.ModifiedFiles) > 0 {
		fmt.Println("Modified files")
		for idx := range note.ModifiedFiles {
			fmt.Printf("[%v]\n", note.ModifiedFiles[idx])
		}
	}
	return nil
}

func init() {
	Register("show", NewShowCommandFactory)
}
