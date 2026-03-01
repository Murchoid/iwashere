package commands

import (
	"fmt"
)

type DeleteCommand struct{
	BaseCommand
}

func NewDeleteCommandFactory() Command {
	return &DeleteCommand{
		BaseCommand{
			NameStr: "delete",
			DescStr: "deletes a note",
			UsageStr: "iwashere delete/rm <id>",
			ExamplesList: []string{
				"iwashere delete 123",
				"iwashere rm 123",
			},
		},
	}
}

func (a *DeleteCommand) Name() string {
	return a.BaseCommand.Name()
}

func (a *DeleteCommand) Description() string {
	return a.BaseCommand.Description()
}

func (a *DeleteCommand) Usage() string {
	return a.BaseCommand.Usage()
}

func (a *DeleteCommand) Examples() []string {
	return a.BaseCommand.Examples()
}

func (a *DeleteCommand) Execute(ctx *Context) error {

	repo := ctx.Repo
	id := ctx.Args[0]

	if err := repo.DeleteNote(id); err != nil {
		return err
	}

	fmt.Printf("Note #%v deleted\n", id)
	return nil
}

func init() {
	Register("delete", NewDeleteCommandFactory)
}
