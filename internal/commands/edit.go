package commands

import (
	"fmt"

	"githum.com/Murchoid/iwashere/internal/domain/models"
)

type EditCommand struct{}

func NewEditCommandFactory() Command {
	return &EditCommand{}
}

func (a *EditCommand) Name() string {
	return "add"
}

func (a *EditCommand) Description() string {
	return "Add a new note"
}

func (a *EditCommand) Execute(ctx *Context) error {

	repo := ctx.Repo
	id := ctx.Args[0]
	newMsg := ctx.Flags["--message"]

	newNote := models.Note {
		ID: id,
		Message: newMsg,
		ProjectPath: ctx.ProjectPath,
	}
	if err := repo.UpdateNote(&newNote); err != nil {
		return err
	}

	fmt.Printf("Edited #%v note\n", id)
	return nil
}

func init() {
	Register("edit", NewEditCommandFactory)
}
