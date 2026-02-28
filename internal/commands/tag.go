package commands

import (
	"fmt"
	"strings"
	"time"

	"githum.com/Murchoid/iwashere/internal/repository"
)

type TagCommand struct{}

func NewTagCommandFactory() Command {
	return &TagCommand{}
}

func (a *TagCommand) Name() string {
	return "add"
}

func (a *TagCommand) Description() string {
	return "Add a new note"
}

func (a *TagCommand) Execute(ctx *Context) error {
	if ctx.Repo == nil {
		return fmt.Errorf("not in an iwashere project (run iwashere init first)")
	}

	repo := ctx.Repo
	tag := ctx.Args[0]
	if tag == "" {
		return fmt.Errorf("tag name required (use tag or provide as argument)")
	}

	for idx := range ctx.Args {
		fmt.Println("arguments: ", ctx.Args[idx])
	}

	switch (tag) {
	case "add":
		if err:= addNewTag(repo, ctx.Args[1:]); err != nil {
			return err
		}
	case "remove":
		if err := removeTag(repo, ctx.Args[1:]); err != nil {
			return err
		}
	}

	return nil
}

func addNewTag(repo repository.Repository, args []string) error {
	id:=args[0]
	tags := args[1:]

	note, err := repo.GetNote(id)
	
	if err!= nil {
		return err
	}


	note.Tags = append(note.Tags, tags...)

	if err:=repo.UpdateNote(note); err != nil {
		return err
	}


	fmt.Println("New tag added to note #",id)
	return nil
}

func removeTag(repo repository.Repository, args []string) error {
	id:=args[0]
	tags := args[1:]

	note, err := repo.GetNote(id)
	
	if err!= nil {
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

	if err:=repo.UpdateNote(note); err != nil {
		return err
	}

	fmt.Printf("Tags %v removed from note #%v", tags, id)
	return nil
}

func init() {
	Register("tag", NewTagCommandFactory)
}
