package commands

import (
	"fmt"

	"githum.com/Murchoid/iwashere/internal/domain/models"
)

type EditCommand struct{
	BaseCommand
}

func NewEditCommandFactory() Command {
	return &EditCommand{
		BaseCommand{
			NameStr: "edit",
			DescStr: "Edits a note",
			UsageStr: "iwashere edit <id> --message <message>",
			ExamplesList: []string{
				"iwashere edit 123 \"Edited this note to something new\"",
			},
		},
	}
}

func (a *EditCommand) Name() string {
	return a.BaseCommand.Name()
}

func (a *EditCommand) Description() string {
	return a.BaseCommand.Description()
}


func (a *EditCommand) Usage() string {
	return a.BaseCommand.Usage()
}

func (a *EditCommand) Examples() []string {
	return a.BaseCommand.Examples()
}

func (a *EditCommand) Execute(ctx *Context) error {

	repo := ctx.Repo
	id := ctx.Args[0]
	newMsg := ctx.Flags["--message"]

	newNote := models.Note{
		ID:          id,
		Message:     newMsg,
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
