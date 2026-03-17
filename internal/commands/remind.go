package commands

import (
	"fmt"
	"time"

	"github.com/Murchoid/iwashere/internal/domain/models"
	"github.com/Murchoid/iwashere/internal/utils"
)

type RemindCommand struct {
	spec *CommandSpec
}

func NewRemindCommandFactory() Command {
	return &RemindCommand{
		spec: RemindCommandSpec,
	}
}

func (c *RemindCommand) Name() string        { return "remind" }
func (c *RemindCommand) Description() string { return "Set reminders for notes" }
func (c *RemindCommand) Usage() string       { return c.spec.Usage }
func (c *RemindCommand) Examples() []string {
	return []string{
		"iwashere remind 123 --at 'tomorrow 9:00'",
		"iwashere remind list",
		"iwashere remind delete rmd_abc123",
	}
}

func (c *RemindCommand) Execute(ctx *Context) error {
	if ctx.Repo == nil {
		return fmt.Errorf("not in an iwashere project (run iwashere init first)")
	}

	parsed, err := c.spec.Parse(ctx.Args)
	if err != nil {
		// Show help on parse error
		utils.PrintCommandHelp(c.Name(), c.Description(), c.Usage(), c.Examples())
		return fmt.Errorf("invalid arguments: %w", err)
	}

	// Handle subcommands
	switch parsed.Subcommand {
	case "list":
		return c.listReminders(ctx, parsed)
	case "delete":
		return c.deleteReminder(ctx, parsed)
	case "done":
		return c.markAsDone(ctx, parsed)
	default:
		// No subcommand = add reminder
		return c.addReminder(ctx, parsed)
	}
}

func (c *RemindCommand) addReminder(ctx *Context, parsed *ParsedArgs) error {
	// Get note ID from positional args
	var noteID string
	if len(parsed.Positional) > 0 {
		noteID = parsed.Positional[0]
	}

	// Get --at flag (required)
	atFlag := parsed.Flags["at"]
	dueTime, err := atFlag.Time()
	if err != nil {
		return fmt.Errorf("--at flag required")
	}

	// Get message (required)
	var message string
	if msgFlag, ok := parsed.Flags["message"]; ok && msgFlag.Present {
		msg, _ := msgFlag.String()
		message = msg
	} else if noteID != "" {
		note, err := ctx.Repo.GetNote(noteID)
		if err != nil {
			return err
		}
		message = note.Message
	} else {
		return fmt.Errorf("either note-id or --message is required")
	}

	// Get repeat (optional)
	repeats := models.Once
	if repeatFlag, ok := parsed.Flags["repeat"]; ok && repeatFlag.Present {
		rep, err := repeatFlag.String()
		if err != nil {
			return err
		}

		switch rep {
		case models.Daily, models.Once, models.Weekly, models.Monthly, models.Yearly:
			repeats = rep
		default:
			return fmt.Errorf("unrecognize repeat value use (daily, once, weekly, monthly, yearly)")
		}
	}

	reminder := &models.Reminder{
		ID:        utils.GenerateId(),
		NoteID:    noteID,
		Message:   message,
		DueAt:     dueTime,
		Repeats:   repeats,
		CreatedAt: time.Now(),
		Active:    true,
	}

	if err := saveReminder(ctx, reminder); err != nil {
		return err
	}

	fmt.Printf("Reminder set for %s\n", dueTime.Format("Mon Jan 2 at 15:04"))
	fmt.Printf("  Message: %s\n", truncate(message, 50))
	if repeats != models.Once {
		fmt.Printf("   Repeats: %s\n", repeats)
	}

	return nil
}

func (c *RemindCommand) listReminders(ctx *Context, parsed *ParsedArgs) error {
	// Get filter options
	all, _ := parsed.Flags["all"].Bool()
	noteFilter, _ := parsed.Flags["note"].String()

	reminders, _ := getReminders(ctx)

	fmt.Println("Reminders")
	fmt.Println("===========")
	shown := false
	for _, r := range reminders {
		if !all && !r.Active {
			continue
		}
		if noteFilter != "" && r.NoteID != noteFilter {
			continue
		}

		status := "Active"
		if !r.Active {
			status = "Done"
		}

		fmt.Printf("%s%s%s %s[%s]%s %s\n", utils.ColorCyan, status, utils.ColorReset, utils.ColorBlue, r.ID, utils.ColorReset, r.Message)
		fmt.Printf("   Due: %s (%s)\n",
			r.DueAt.Format("Mon Jan 2 15:04"),
			utils.HowLongAgo(r.CreatedAt, 0))
		if r.Repeats != models.Once {
			fmt.Printf("   Repeats: %s\n", r.Repeats)
		}
		fmt.Println()
		shown = true
	}

	if !shown {
		fmt.Println("No reminders you are all good")
	}
	return nil
}

func (c *RemindCommand) deleteReminder(ctx *Context, parsed *ParsedArgs) error {
	if len(parsed.Positional) == 0 {
		return fmt.Errorf("reminder-id required")
	}

	reminderID := parsed.Positional[0]
	if err := deleteReminder(ctx, reminderID); err != nil {
		return err
	}

	fmt.Printf("Reminder %s deleted\n", reminderID[:8])
	return nil
}

func (c *RemindCommand) markAsDone(ctx *Context, parsed *ParsedArgs) error {
	if len(parsed.Positional) == 0 {
		return fmt.Errorf("reminder-id required")
	}

	reminderID := parsed.Positional[0]
	if err := markReminderDone(ctx, reminderID); err != nil {
		return fmt.Errorf("error marking reminder as done: %w", err)
	}

	fmt.Printf("Reminder %s marked as done!\n", reminderID[:8])
	return nil
}

// helper functions
func saveReminder(ctx *Context, reminder *models.Reminder) error {
	repo := ctx.Repo
	return repo.SaveReminder(reminder)
}

func getReminders(ctx *Context) ([]*models.Reminder, error) {
	repo := ctx.Repo

	reminders, err := repo.ListReminders()

	if err != nil {
		return nil, err
	}

	return reminders, nil
}

func deleteReminder(ctx *Context, reminderID string) error {
	repo := ctx.Repo
	return repo.DeleteReminder(reminderID)
}

func markReminderDone(ctx *Context, reminderID string) error {
	repo := ctx.Repo
	return repo.DeactivateOrUpdateReminder(reminderID)
}

func init() {
	Register("remind", NewRemindCommandFactory)
}
