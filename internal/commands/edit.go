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
	if len(ctx.Args) <= 1 {
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
		return fmt.Errorf("Missing arguments")
	}
	parseArgs, err := a.spec.Parse(ctx.Args)
	if err != nil {
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
		return fmt.Errorf("invalid arguments: %w", err)
	}

	id := parseArgs.Positional[0]

	newNote := models.PrivateNote{
		ID:          id,
		ProjectPath: ctx.ProjectPath,
	}

	//check for message
	msg := parseArgs.Flags["message"]
	newMsg, err := msg.String()
	if err != nil && msg.Present {
		return err
	}

	if msg.Present {
		newNote.Message = newMsg
		if err := repo.UpdateMessage(&newNote); err != nil {
			return err
		}
	}

	//check for tags
	tags := parseArgs.Flags["tags"]
	pTags, err := tags.String()

	if err != nil && tags.Present {
		return err
	}

	if tags.Present {
		newNote.Tags = utils.ParseTags(pTags)
		if err := repo.UpdateTags(&newNote); err != nil {
			return err
		}
	}

	//check for appended tags
	appendTags := parseArgs.Flags["add-tags"]
	aTags, err := appendTags.String()
	if err != nil && appendTags.Present {
		return err
	}

	if appendTags.Present {
		newNote.Tags = utils.ParseTags(aTags)
		if err := repo.AddTagsToNote(&newNote); err != nil {
			return err
		}
	}

	//check for tags being removed
	removeTags := parseArgs.Flags["remove-tags"]
	rTags, err := removeTags.String()
	if err != nil && removeTags.Present {
		return err
	}

	if removeTags.Present {
		newNote.Tags = utils.ParseTags(rTags)
		if err := repo.RemoveTagsFromNote(&newNote); err != nil {
			return err
		}
	}

	fmt.Printf("Edited #%v note\n", id)
	return nil
}

func init() {
	Register("edit", NewEditCommandFactory)
}
