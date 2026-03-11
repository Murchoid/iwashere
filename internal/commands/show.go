package commands

import (
	"fmt"

	"github.com/Murchoid/iwashere/internal/domain/models"
	"github.com/Murchoid/iwashere/internal/utils"
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

	if len(ctx.Args) == 0 {
		fmt.Println("Id must be provided")
		fmt.Println()
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
		return nil
	}

	if len(ctx.Args) > 1 {
		fmt.Println("show only accepts one argument")
		fmt.Println()
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
		return nil
	}

	id := ctx.Args[0]

	note, err := repo.GetNote(id)
	if err != nil {
		return err
	}

	format := "detailed"

	if ctx.Flags["--short"] != "" {
		if ctx.Flags["--short"] == "true" {
			format = "short"
		} else {
			fmt.Println("Unrecognized argument after short")
			fmt.Println()
			utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
			return nil
		}
	}

	if ctx.Flags["--compact"] != "" {
		if ctx.Flags["--compact"] == "true" {
			format = "compact"
		} else {
			fmt.Println("Unrecognized argument after compact")
			fmt.Println()
			utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
			return nil
		}
	}
	utils.PrintNotes([]*models.PrivateNote{note}, nil, format)

	return nil
}

func init() {
	Register("show", NewShowCommandFactory)
}
