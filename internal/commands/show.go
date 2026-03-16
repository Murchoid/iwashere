package commands

import (
	"fmt"

	"github.com/Murchoid/iwashere/internal/domain/models"
	"github.com/Murchoid/iwashere/internal/utils"
)

type ShowCommand struct {
	spec        *CommandSpec
	baseCommand BaseCommand
}

func NewShowCommandFactory() Command {
	return &ShowCommand{
		spec: ShowCommandSpec,
		baseCommand: BaseCommand{
			NameStr:  "show",
			DescStr:  "Shows the note specified by the id",
			UsageStr: "iwashere show <id>",
			ExamplesList: []string{
				"iwashere show 123",
				"iwashere s 123",
			},
		},
	}
}

func (a *ShowCommand) Name() string {
	return a.baseCommand.Name()
}

func (a *ShowCommand) Description() string {
	return a.baseCommand.Description()
}

func (a *ShowCommand) Usage() string {
	return a.baseCommand.Usage()
}

func (a *ShowCommand) Examples() []string {
	return a.baseCommand.Examples()
}

func (a *ShowCommand) Execute(ctx *Context) error {

	repo := ctx.Repo

	parsedArgs, err := a.spec.Parse(ctx.Args)

	if err != nil {
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
		return fmt.Errorf("invalid arguments: %w", err)
	}

	if len(parsedArgs.Positional) > 1 {
		fmt.Println("show only accepts one argument")
		fmt.Println()
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
		return nil
	}

	id := parsedArgs.Positional[0]

	note, err := repo.GetNote(id)
	if err != nil {
		return err
	}

	format := "detailed"
	for f := range parsedArgs.Flags {
		switch f {
		case "short":
			format = "short"
		case "compact":
			format = "compact"
		default:
			format = "detailed"
		}
	}

	session, err := repo.GetSession(note.SessionID)
	if err != nil && note.SessionID != ""{
		return err
	}
	sessionMap := map[string]*models.Session{}
	sessionMap[note.SessionID] = session
	utils.PrintNotes([]*models.PrivateNote{note}, sessionMap, format)

	return nil
}

func init() {
	Register("show", NewShowCommandFactory)
}
