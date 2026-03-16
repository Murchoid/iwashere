package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/Murchoid/iwashere/internal/domain/models"
	"github.com/Murchoid/iwashere/internal/repository"
	"github.com/Murchoid/iwashere/internal/utils"
)

type  listFlags struct {
			cloud bool
			format string
			limit int
		}
	
type TagCommand struct {
	spec        *CommandSpec
	baseCommand BaseCommand
}

func NewTagCommandFactory() Command {
	return &TagCommand{
		spec: TagCommandSpec,
		baseCommand: BaseCommand{
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
	return a.baseCommand.Name()
}

func (a *TagCommand) Description() string {
	return a.baseCommand.Description()
}

func (a *TagCommand) Usage() string {
	return a.baseCommand.Usage()
}

func (a *TagCommand) Examples() []string {
	return a.baseCommand.Examples()
}

func (a *TagCommand) Execute(ctx *Context) error {
	if ctx.Repo == nil {
		return fmt.Errorf("not in an iwashere project (run iwashere init first)")
	}

	repo := ctx.Repo

	parsedArgs, err := a.spec.Parse(ctx.Args)

	if err != nil {
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
		return fmt.Errorf("invalid arguments: %w", err)
	}

	if parsedArgs.Subcommand == "" {

		fmt.Println("Tag option must be provided")
		fmt.Println()
		utils.PrintCommandHelp(a.Name(), a.Description(), a.Usage(), a.Examples())
		return nil
	}

	tagCommand := parsedArgs.Subcommand

	var tagsInfo []string
	if len(parsedArgs.Positional) > 0 {
		tagsInfo = parsedArgs.Positional //note-id at index 0, and tag at index 1

	}

	switch tagCommand {
	case "add":
		if err := addNewTag(repo, tagsInfo); err != nil {
			return err
		}
	case "remove":
		if err := removeTag(repo, tagsInfo); err != nil {
			return err
		}
	case "list":
		cloud := parsedArgs.Flags["cloud"]
		pCloud, err := cloud.Bool()

		if err != nil && cloud.Present {
			return err
		}

		format := ""
		for f := range parsedArgs.Flags {
			switch f {
			case "short":
				format = "short"
			case "compact":
				format = "compact"
			case "detailed":
				format = "detailed"
			default:
				format = "short"
			}
		}
		limit := parsedArgs.Flags["limit"]
		pLimit, err := limit.Int()
		if err!= nil && limit.Present {
			return err
		}
		if !limit.Present {
			pLimit = 5
		}

		flags := listFlags{
			cloud: pCloud,
			format: format,
			limit: pLimit,
		}

		if err := listTag(repo, tagsInfo,flags); err != nil {
			return err
		}

	}

	return nil
}

func addNewTag(repo repository.Repository, tagInfo []string) error {
	if len(tagInfo) == 0 {
		return fmt.Errorf("An id must be provided")
	}
	id := tagInfo[0]

	if len(tagInfo) == 1 {
		return fmt.Errorf("At least one tag to be added must be given")
	}
	tags := tagInfo[1:]

	note, err := repo.GetNote(id)

	if err != nil {
		return err
	}

	note.Tags = append(note.Tags, tags...)

	if err := repo.AddTagsToNote(note); err != nil {
		return err
	}

	fmt.Println("New tag added to note #", id)
	return nil
}

func removeTag(repo repository.Repository, tagInfo []string) error {
	if len(tagInfo) == 0 {
		return fmt.Errorf("An id must be provided")
	}
	id := tagInfo[0]

	if len(tagInfo[1:]) == 0 {
		return fmt.Errorf("At least one tag to be added must be given")
	}
	tags := tagInfo[1:]

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

	if err := repo.RemoveTagsFromNote(note); err != nil {
		return err
	}

	fmt.Printf("Tags '%v' removed from note #%v", tags, id)
	return nil
}

func listTag(repo repository.Repository, tagInfo []string, flags listFlags) error {
	var tags []string
	filter := repository.NoteFilter{}
	var notes []*models.PrivateNote
	var err error

	if len(tagInfo) > 0 && !flags.cloud {
		tags = utils.ParseTags(tagInfo[0])
		filter.Tags = tags
		filter.Limit = flags.limit
		notes, err = repo.ListNotes(&filter)
		
		if err != nil {
			return err
		}

		fmt.Printf("Notes containing %v\n", tags)
		fmt.Println("===============================")
		fmt.Println()
		for _, note := range notes {
			if err := oneTagResult(repo, note, flags.format); err != nil {
				fmt.Println("An error occured: ", err)
			}
		}

	} else {
		notes, err = repo.ListNotes(nil)
		if err != nil {
			return err
		}

		if !flags.cloud {
			AllTagResults(notes)
		} else if flags.cloud {
			AllTagCloudResults(notes)
		}

	}

	return nil
}

func AllTagResults(notes []*models.PrivateNote) {
	fmt.Println("All Tags")
	fmt.Println("==========")

	var fTags = map[string][]*models.PrivateNote{}

	for _, note := range notes {
		for _, tag := range note.Tags {
			fTags[tag] = append(fTags[tag], note)
		}
	}
	longestTag := 0
	for tag := range fTags {
		if len([]byte(tag)) > longestTag {
			longestTag = len([]byte(tag))
		}
	}

	for tag, tagNote := range fTags {
		fmt.Printf("%-*s (%d notes)\n",longestTag, tag, len(tagNote))
	}
}

func oneTagResult(repo repository.Repository, note *models.PrivateNote, format string) error {

	session, err := repo.GetSession(note.SessionID)
	if err != nil && note.SessionID != ""{
		return err
	}
	sessionMap := map[string]*models.Session{}
	sessionMap[note.SessionID] = session

	utils.PrintNotes([]*models.PrivateNote{note}, sessionMap, format)

	return nil
}

func AllTagCloudResults(notes []*models.PrivateNote) {
	fmt.Println("Tags Cloud")
	fmt.Println("==========")

	var fTags = map[string][]*models.PrivateNote{}
	longestTag := 0

	for _, note := range notes {
		for _, tag := range note.Tags {
			if len([]byte(tag)) > longestTag {
				longestTag = len([]byte(tag))
			}
			fTags[tag] = append(fTags[tag], note)
		}
	}

	for tag, tagNote := range fTags {
		fmt.Printf("%-*s %s\n", longestTag, tag, printStars(len(tagNote)))
	}
}

func printStars(num int) string {
	var stars []byte
	for i := 0; i < num; i++ {
		stars = append(stars, '*')
	}

	return string(stars)
}

func init() {
	Register("tag", NewTagCommandFactory)
}
