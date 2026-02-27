package commands

import (
	"fmt"
)

type DeleteCommand struct{}

func NewDeleteCommandFactory() Command {
	return &DeleteCommand{}
}

func (a *DeleteCommand) Name() string {
	return "delete"
}

func (a *DeleteCommand) Description() string {
	return "Delete a note with id"
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
