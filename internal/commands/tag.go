package commands

import (
	"fmt"
	"strings"
	"time"

	"githum.com/Murchoid/iwashere/internal/repository"
	"githum.com/Murchoid/iwashere/internal/utils"
)

type TagCommand struct {
	BaseCommand
}

func NewTagCommandFactory() Command {
	return &TagCommand{
		BaseCommand{
			NameStr: "tag",
			DescStr: "add or remove a tag from a note",
			UsageStr: `iwashere tag <subcommand> [arguments]

Subcommands:
  add    <note-id> <tag>     Add tag to note
  remove <note-id> <tag>     Remove tag from note
  list   [tag]              List notes by tag`,
			ExamplesList: []string{
				"iwashere tag add 123 bug",
				"iwashere tag add 456 urgent,frontend",
				"iwashere tag remove 123 bug",
			},
		},
	}
}

func (a *TagCommand) Name() string {
	return a.BaseCommand.Name()
}

func (a *TagCommand) Description() string {
	return a.BaseCommand.Description()
}

func (a *TagCommand) Usage() string {
	return a.BaseCommand.Usage()
}

func (a *TagCommand) Examples() []string {
	return a.BaseCommand.Examples()
}

func (a *TagCommand) Execute(ctx *Context) error {
	if ctx.Repo == nil {
		return fmt.Errorf("not in an iwashere project (run iwashere init first)")
	}

	repo := ctx.Repo
	if len(ctx.Args) == 0 {

		fmt.Println("Tag option must be provided")
		fmt.Println()
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
		return nil
	}
	tag := ctx.Args[0]
	if tag == "" {
		fmt.Println("Cannot use an empty tag option")
		return nil
	}

	if len(ctx.Args) <= 1 {
		fmt.Println("Tag names must be provided")
		fmt.Println()
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
		return nil
	}

	switch tag {
	case "add":
		if err := addNewTag(repo, ctx.Args[1:]); err != nil {
			return err
		}
	case "remove":
		if err := removeTag(repo, ctx.Args[1:]); err != nil {
			return err
		}
	case "list":
		if err := listTag(repo, ctx.Args[1:]); err != nil {
			return err
		}
	}

	return nil
}

func addNewTag(repo repository.Repository, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("An id must be provided")
	}
	id := args[0]

	if len(args[1:]) == 0 {
		return fmt.Errorf("At least one tag to be added must be given")
	}
	tags := args[1:]

	note, err := repo.GetNote(id)

	if err != nil {
		return err
	}

	note.Tags = append(note.Tags, tags...)

	if err := repo.UpdateNote(note); err != nil {
		return err
	}

	fmt.Println("New tag added to note #", id)
	return nil
}

func removeTag(repo repository.Repository, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("An id must be provided")
	}
	id := args[0]

	if len(args[1:]) == 0 {
		return fmt.Errorf("At least one tag to be added must be given")
	}
	tags := args[1:]

	note, err := repo.GetNote(id)

	if err != nil {
		return err
	}

	var newTags []string
	note.UpdatedAt = time.Now()

	for idx := range note.Tags {
		if idx > len(tags)-1 {
			break
		}
		if strings.Compare(note.Tags[idx], tags[idx]) != 0 {
			newTags = append(newTags, note.Tags[idx])
		}
	}

	note.Tags = newTags
	note.UpdatedAt = time.Now()

	if err := repo.UpdateNote(note); err != nil {
		return err
	}

	fmt.Printf("Tags %v removed from note #%v", tags, id)
	return nil
}

func listTag(repo repository.Repository, args []string) error {
	tags := utils.ParseTags(args[0])

	notes, err := repo.ListNotes(&repository.NoteFilter{Tags: tags})

	if err != nil {
		return err
	}

	utils.PrintNotes(notes, nil, "short")

	return nil
}

func init() {
	Register("tag", NewTagCommandFactory)
}
