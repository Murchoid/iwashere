package commands

import (
	"fmt"

	"github.com/Murchoid/iwashere/internal/domain/models"
	"github.com/Murchoid/iwashere/internal/utils"
)

type EditCommand struct {
	spec        *CommandSpec
	baseCommand BaseCommand
}

func NewEditCommandFactory() Command {
	return &EditCommand{
		spec: EditCommandSpec,
		baseCommand: BaseCommand{
			NameStr:  "edit",
			DescStr:  "Edits a note",
			UsageStr: "iwashere edit <id> --message <message>",
			ExamplesList: []string{
				"iwashere edit 123 \"Edited this note to something new\"",
			},
		},
	}
}

func (a *EditCommand) Name() string {
	return a.baseCommand.Name()
}

func (a *EditCommand) Description() string {
	return a.baseCommand.Description()
}

func (a *EditCommand) Usage() string {
	return a.baseCommand.Usage()
}

func (a *EditCommand) Examples() []string {
	return a.baseCommand.Examples()
}

func (a *EditCommand) Execute(ctx *Context) error {

	repo := ctx.Repo

	parseArgs, err := a.spec.Parse(ctx.Args)
	if err != nil {
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
		return fmt.Errorf("invalid arguments: %w", err)
	}

	id := parseArgs.Positional[0]

	newMsg, err := parseArgs.Flags["--message"].String()
	if err != nil {
		return err
	}

	newNote := models.PrivateNote{
		ID:          id,
		Message:     newMsg,
		ProjectPath: ctx.ProjectPath,
	}
	if err := repo.UpdateMessage(&newNote); err != nil {
		return err
	}

	tags := parseArgs.Flags["tags"]
	pTags, err := tags.String()

	if err != nil && tags.Present {
		return err
	}

	newNote.Tags = utils.ParseTags(pTags)
	if err := repo.UpdateTags(&newNote); err != nil {
		return err
	}

	fmt.Printf("Edited #%v note\n", id)
	return nil
}

func init() {
	Register("edit", NewEditCommandFactory)
}
