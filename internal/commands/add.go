package commands

import (
	"fmt"
	"time"

	"githum.com/Murchoid/iwashere/internal/domain/models"
	"githum.com/Murchoid/iwashere/internal/utils"
)

type AddCommand struct{}

func NewAddCommandFactory() Command {
	return &AddCommand{}
}

func (a *AddCommand) Name() string {
	return "add"
}

func (a *AddCommand) Description() string {
	return "Add a new note"
}

func (a *AddCommand) Execute(ctx *Context) error {

	repo := ctx.Repo
	msg := ctx.Args[0]

	newNote := models.Note{
		ID:          utils.GenerateId(),
		Message:     msg,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ProjectPath: ctx.ProjectPath,
	}

	if err := repo.SaveNote(&newNote); err != nil {
		return err
	}

	fmt.Println("New note added : ", newNote.Message)
	return nil
}

func init() {
	Register("add", NewAddCommandFactory)
}
