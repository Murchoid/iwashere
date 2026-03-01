package commands

import (
	"fmt"

	"githum.com/Murchoid/iwashere/internal/utils"
)

type ShowCommand struct {
	BaseCommand
}

func NewShowCommandFactory() Command {
	return &ShowCommand{
		BaseCommand{
			NameStr:  "show",
			DescStr:  "Shows the note specified by the id",
			UsageStr: "iwashere show/s <id>",
			ExamplesList: []string{
				"iwashere show 123",
				"iwashere s 123",
			},
		},
	}
}

func (a *ShowCommand) Name() string {
	return a.BaseCommand.Name()
}

func (a *ShowCommand) Description() string {
	return a.BaseCommand.Description()
}

func (a *ShowCommand) Usage() string {
	return a.BaseCommand.Usage()
}

func (a *ShowCommand) Examples() []string {
	return a.BaseCommand.Examples()
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
