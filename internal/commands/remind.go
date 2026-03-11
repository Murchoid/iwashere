package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/Murchoid/iwashere/internal/domain/models"
	"github.com/Murchoid/iwashere/internal/utils"
)

type RemindCommand struct {
	BaseCommand
}

func NewRemindCommand() Command {
	return &RemindCommand{
		BaseCommand{
			NameStr: "remind",
			DescStr: "Set reminders for notes",
			UsageStr: `iwashere remind <note-id> --at <when>
iwashere remind list
iwashere remind delete <reminder-id>`,
			ExamplesList: []string{
				"iwashere remind 123 --at 'tomorrow 9am'",
				"iwashere remind 123 --at 'in 2h'",
				"iwashere remind 123 --at 'in 2h45m'",
				"iwashere remind list",
				"iwashere remind delete <reminder-id>",
			},
		},
	}
}

func (c *RemindCommand) Name() string {
	return c.BaseCommand.Name()
}

func (c *RemindCommand) Description() string {
	return c.BaseCommand.Description()
}

func (c *RemindCommand) Usage() string {
	return c.BaseCommand.Usage()
}

func (c *RemindCommand) Examples() []string {
	return c.BaseCommand.Examples()
}

func (c *RemindCommand) Execute(ctx *Context) error {
	if ctx.Repo == nil {
		return fmt.Errorf("not in an iwashere project (run iwashere init first)")
	}

	if len(ctx.Args) == 0 {
		return c.listReminders(ctx)
	}

	switch ctx.Args[0] {
	case "list":
		return c.listReminders(ctx)
	case "delete":
		return c.deleteReminder(ctx)
	default:
		// Assume it's a note ID
		return c.addReminder(ctx)
	}
}

func (c *RemindCommand) addReminder(ctx *Context) error {
	noteID := ctx.Args[0]
	when := ctx.Flags["--at"]
	if when == "" {
		return fmt.Errorf("--at is required (e.g., --at 'tomorrow 9am')")
	}

	// Get the note
	note, err := ctx.Repo.GetNote(noteID)
	if err != nil {
		return err
	}

	// Parse when
	dueTime, err := parseNaturalTime(when)
	if err != nil {
		return fmt.Errorf("could not parse time: %s", when)
	}

	reminder := &models.Reminder{
		ID:        utils.GenerateId(),
		NoteID:    noteID,
		Message:   note.Message,
		DueAt:     dueTime,
		CreatedAt: time.Now(),
		Active:    true,
	}

	if err := saveReminder(ctx, reminder); err != nil {
		return err
	}

	fmt.Printf("Reminder set for %s\n", dueTime.Format("Mon Jan 2 at 15:04"))
	fmt.Printf("   Note: %s\n", truncate(note.Message, 50))

	return nil
}

func parseNaturalTime(input string) (time.Time, error) {
	now := time.Now()

	// Handle "in X minutes/hours/days"
	if strings.HasPrefix(input, "in ") {
		duration, err := time.ParseDuration(strings.TrimPrefix(input, "in "))
		if err != nil {
			return time.Time{}, err
		}
		return now.Add(duration), nil
	}

	// Handle "tomorrow 9am"
	if strings.HasPrefix(input, "tomorrow") {
		t, err := time.Parse("15:04", strings.TrimPrefix(input, "tomorrow "))
		if err != nil {
			return time.Time{}, err
		}
		tomorrow := now.Add(24 * time.Hour)
		return time.Date(
			tomorrow.Year(), tomorrow.Month(), tomorrow.Day(),
			t.Hour(), t.Minute(), 0, 0, now.Location(),
		), nil
	}

	if strings.HasPrefix(input, "today ") {
		t, err := time.Parse("15:04", strings.TrimPrefix(input, "today "))
		if err != nil {
			return time.Time{}, err
		}
		return time.Date(
			now.Year(), now.Month(), now.Day(),
			t.Hour(), t.Minute(), 0, 0, now.Location(),
		), nil
	}

	// Handle "Friday 5pm"
	// This is complex - for v1, use a library or simplify
	return time.Parse("2006-01-02 15:04", input)
}

func (c *RemindCommand) listReminders(ctx *Context) error {
	repo := ctx.Repo

	if len(ctx.Args) > 1 {
		fmt.Println("Unknown arguments")
		utils.PrintCommandHelp(c.Name(), c.Description(), c.Usage(), c.Examples())
		return nil
	}

	reminders, err := repo.ListReminders()

	if err != nil {
		return err
	}

	utils.PrintReminders(reminders, true, repo)

	return nil
}

func (c *RemindCommand) deleteReminder(ctx *Context) error {
	repo := ctx.Repo

	if len(ctx.Args) < 1 {
		fmt.Println("Too few arguments")
		utils.PrintCommandHelp(c.Name(), c.Description(), c.Usage(), c.Examples())
		return nil
	}

	if len(ctx.Args) > 2 {
		fmt.Println("Unknown arguments")
		utils.PrintCommandHelp(c.Name(), c.Description(), c.Usage(), c.Examples())
		return nil
	}

	id := ctx.Args[1]
	err := repo.DeleteReminder(id)

	if err != nil {
		return err
	}

	return nil
}

func saveReminder(ctx *Context, reminder *models.Reminder) error {
	repo := ctx.Repo

	if len(ctx.Args) < 1 {
		fmt.Println("Too few arguments")
		return nil
	}

	err := repo.SaveReminder(reminder)

	if err != nil {
		return err
	}

	return nil
}

func init() {
	Register("remind", NewRemindCommand)
}
