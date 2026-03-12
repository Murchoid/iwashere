package commands

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/Murchoid/iwashere/internal/domain/models"
	"github.com/Murchoid/iwashere/internal/utils"
)

type WeekDays int

const (
	Sunday WeekDays = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

type Repeatition string

const (
	Daily   Repeatition = "daily"
	Weekly  Repeatition = "weekly"
	Monthly Repeatition = "monthly"
	None    Repeatition = "none"
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

	if len(ctx.Args) == 0 && ctx.Flags["--message"] == "" {
		return c.listReminders(ctx)
	}

	if len(ctx.Args) > 0 {
		switch ctx.Args[0] {
		case "list":
			return c.listReminders(ctx)
		case "delete":
			return c.deleteReminder(ctx)
		}
	} else {
		return c.addReminder(ctx)
	}

	return nil
}

func (c *RemindCommand) addReminder(ctx *Context) error {
	var noteID string
	if len(ctx.Args) > 0 {
		noteID = ctx.Args[0]
	}

	when := ctx.Flags["--at"]
	if when == "" {
		return fmt.Errorf("--at is required (e.g., --at 'tomorrow 9am')")
	}

	var note *models.PrivateNote
	if noteID != "" {
		var err error
		note, err = ctx.Repo.GetNote(noteID)
		if err != nil {
			return err
		}
	}

	dueTime, err := parseNaturalTime(when)
	if err != nil {
		return fmt.Errorf("could not parse time: %s %w", when, err)
	}

	var message string
	if noteID == "" && note == nil {
		if ctx.Flags["--message"] != "" {
			message = ctx.Flags["--message"]
		} else {
			return fmt.Errorf("A message or note id must be provided")
		}
	} else if note != nil {
		message = note.Message
	}

	repeats := None

	if ctx.Flags["--repeat"] != "" {
		repeats = Repeatition(strings.ToLower(ctx.Flags["--repeat"]))
	}

	reminder := &models.Reminder{
		ID:        utils.GenerateId(),
		NoteID:    noteID,
		Message:   message,
		DueAt:     dueTime,
		Repeats:   string(repeats),
		CreatedAt: time.Now(),
		Active:    true,
	}

	if err := saveReminder(ctx, reminder); err != nil {
		return err
	}

	fmt.Printf("Reminder set for %s\n", dueTime.Format("Mon Jan 2 at 15:04"))
	fmt.Printf("   Note: %s\n", truncate(message, 50))
	fmt.Printf("   Repeats: %s\n", reminder.Repeats)

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
	return parseDayTime(input)
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

	if len(ctx.Args) < 1 && ctx.Flags["--message"] == "" {
		fmt.Println("Too few arguments")
		return nil
	}

	err := repo.SaveReminder(reminder)

	if err != nil {
		return err
	}

	return nil
}

func parseDayTime(dayTime string) (time.Time, error) {
	day, t, _ := strings.Cut(dayTime, " ")
	today := int(time.Now().Weekday())

	switch day {
	case "Sunday":
		diff := int(Sunday) - today
		return getReminderFullDateTime(diff, t)
	case "Monday":
		diff := int(Monday) - today
		return getReminderFullDateTime(diff, t)
	case "Tuesday":
		diff := int(Tuesday) - today
		return getReminderFullDateTime(diff, t)
	case "Wednesday":
		diff := int(Wednesday) - today
		return getReminderFullDateTime(diff, t)
	case "Thursday":
		diff := int(Thursday) - today
		return getReminderFullDateTime(diff, t)
	case "Friday":
		diff := int(Friday) - today
		return getReminderFullDateTime(diff, t)
	case "Saturday":
		diff := int(Saturday) - today
		return getReminderFullDateTime(diff, t)
	}

	return time.Now(), nil
}

func getReminderFullDateTime(diff int, t string) (time.Time, error) {

	rTime, err := time.Parse("15:00", t)
	if diff != 0 {
		daysAhead := math.Abs(float64(diff))
		now := time.Now()

		if err != nil {
			return time.Time{}, err
		}

		specificDay := 24 * daysAhead
		remindDate := now.Add(time.Duration(specificDay) * time.Hour)
		setDate := time.Date(
			remindDate.Year(), remindDate.Month(), remindDate.Day(),
			rTime.Hour(), rTime.Minute(), 0, 0, now.Location(),
		)

		return setDate, nil
	}

	//if the difference is 0, means the days are 7days apart, i.e Monday = 1, 1 - 1 is 0
	now := time.Now()
	specificDay := 24 * 7
	remindDate := now.Add(time.Duration(specificDay) * time.Hour)
	setDate := time.Date(
		remindDate.Year(), remindDate.Month(), remindDate.Day(),
		rTime.Hour(), rTime.Minute(), 0, 0, now.Location(),
	)

	return setDate, nil
}

func init() {
	Register("remind", NewRemindCommand)
}
