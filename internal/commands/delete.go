package commands

import (
	"fmt"

	"github.com/Murchoid/iwashere/internal/utils"
)

type DeleteCommand struct {
	spec *CommandSpec
	baseCommand BaseCommand
}

func NewDeleteCommandFactory() Command {
	return &DeleteCommand{
		spec: DeleteCommandSpec,
		baseCommand: BaseCommand{
			NameStr:  "delete",
			DescStr:  "deletes a note",
			UsageStr: "iwashere delete/rm <id>",
			ExamplesList: []string{
				"iwashere delete 123",
				"iwashere rm 123",
			},
		},
	}
}

func (a *DeleteCommand) Name() string {
	return a.baseCommand.Name()
}

func (a *DeleteCommand) Description() string {
	return a.baseCommand.Description()
}

func (a *DeleteCommand) Usage() string {
	return a.baseCommand.Usage()
}

func (a *DeleteCommand) Examples() []string {
	return a.baseCommand.Examples()
}

func (a *DeleteCommand) Execute(ctx *Context) error {

	repo := ctx.Repo
	parseArgs, err := a.spec.Parse(ctx.Args)
	if err!= nil {
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
		return fmt.Errorf("invalid arguments: %w", err)
	}

	var id string
	if len(parseArgs.Positional) > 0 {
		id = parseArgs.Positional[0]
	}

	if err := repo.DeleteNote(id); err != nil {
		return err
	}

	fmt.Printf("Note #%v deleted\n", id)
	return nil
}

func init() {
	Register("delete", NewDeleteCommandFactory)
}
