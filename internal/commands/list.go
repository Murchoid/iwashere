package commands

import (
	"fmt"

	"githum.com/Murchoid/iwashere/internal/repository"
	"githum.com/Murchoid/iwashere/internal/utils"
)

type ListCommand struct{
	BaseCommand
}

func NewListCommandFactory() Command {
	return &ListCommand{
		BaseCommand{
			NameStr: "list",
			DescStr: "Lists all notes in the project/repo",
			UsageStr: "iwashere list/ls [options]",
			ExamplesList: []string{
				"iwashere list",
				"iwashere list --limit 10",
				"iwashere ls",
				"iwashere ls --limit 10",
			},
		},
	}
}

func (a *ListCommand) Name() string {
	return a.BaseCommand.Name()
}

func (a *ListCommand) Description() string {
	return a.BaseCommand.Description()
}

func (a *ListCommand) Usage() string {
	return a.BaseCommand.Usage()
}

func (a *ListCommand) Examples() []string {
	return a.BaseCommand.Examples()
}


func (a *ListCommand) Execute(ctx *Context) error {

	repo := ctx.Repo

	filters := repository.NoteFilter{
		ProjectPath: ctx.ProjectPath,
		Limit:       5,
		Tags: utils.ParseTags(ctx.Flags["--tags"]),
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
