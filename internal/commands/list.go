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
	}
	return nil
}

func init() {
	Register("list", NewListCommandFactory)
}
