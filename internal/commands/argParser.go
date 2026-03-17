package commands

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type WeekDay int

const (
	Sunday WeekDay = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

type ParsedArgs struct {
	Command    string
	Subcommand string
	Positional []string
	Flags      map[string]FlagValue
}

type FlagValue struct {
	Value   any
	Present bool
	Spec    FlagSpec
}

func (c *CommandSpec) Parse(args []string) (*ParsedArgs, error) {
	result := &ParsedArgs{
		Command: c.Name,
		Flags:   make(map[string]FlagValue),
	}

	for _, flag := range c.Flags {
		if flag.Default != nil {
			result.Flags[flag.Name] = FlagValue{
				Value:   flag.Default,
				Present: false,
				Spec:    flag,
			}
		}
	}

	remaining := args

	if len(c.Subcommands) > 0 && len(remaining) > 0 {
		if subSpec, ok := c.Subcommands[remaining[0]]; ok {
			result.Subcommand = remaining[0]
			remaining = remaining[1:]
			return subSpec.parseWithParent(remaining, result)
		}
	}

	return c.parseWithParent(remaining, result)
}

func (c *CommandSpec) parseWithParent(args []string, result *ParsedArgs) (*ParsedArgs, error) {
	for i := 0; i < len(args); i++ {
		arg := args[i]

		if strings.HasPrefix(arg, "-") {
			flagName := strings.TrimLeft(arg, "-")
			var flagSpec *FlagSpec
			for _, f := range c.Flags {
				if f.Name == flagName || f.Short == flagName {
					flagSpec = &f
					break
				}
			}

			if flagSpec == nil {
				return nil, fmt.Errorf("unknown flag: %s", arg)
			}

			var value any
			var present bool

			switch flagSpec.Type {
			case FlagTypeBool:
				value = true
				present = true

			case FlagTypeString, FlagTypeTime:
				if i+1 >= len(args) || strings.HasPrefix(args[i+1], "-") {
					return nil, fmt.Errorf("flag %s requires a value", arg)
				}
				val := args[i+1]
				i++

				switch flagSpec.Type {

				case FlagTypeString:
					value = val

				case FlagTypeTime:
					t, err := parseNaturalTime(val)
					if err != nil {
						return nil, fmt.Errorf("invalid time format for %s: %w", arg, err)
					}
					value = t
				}
				present = true

			case FlagTypeInt:
				if i+1 >= len(args) || strings.HasPrefix(args[i+1], "-") {
					return nil, fmt.Errorf("flag %s requires a value", arg)
				}

				val := args[i+1]
				i++

				intVal, err := strconv.Atoi(val)
				if err != nil {
					return nil, fmt.Errorf("invalid integer for %s: %w", arg, err)
				}

				value = intVal
				present = true
			}

			result.Flags[flagSpec.Name] = FlagValue{
				Value:   value,
				Present: present,
				Spec:    *flagSpec,
			}

		} else {

			result.Positional = append(result.Positional, arg)
		}
	}

	for _, flag := range c.Flags {
		if flag.Required {
			if val, ok := result.Flags[flag.Name]; !ok || !val.Present {
				return nil, fmt.Errorf("required flag --%s is missing", flag.Name)
			}
		}
	}

	for i, argSpec := range c.Args {
		if argSpec.Required && i >= len(result.Positional) {
			return nil, fmt.Errorf("missing required argument: %s", argSpec.Name)
		}
	}

	return result, nil
}

func (f FlagValue) String() (string, error) {
	if !f.Present && f.Spec.Default == nil {
		return "", fmt.Errorf("flag --%s not provided", f.Spec.Name)
	}
	if str, ok := f.Value.(string); ok {
		return str, nil
	}
	return "", fmt.Errorf("flag --%s is not a string", f.Spec.Name)
}

func (f FlagValue) Bool() (bool, error) {
	if f.Spec.Type != FlagTypeBool {
		return false, fmt.Errorf("flag --%s is not a boolean", f.Spec.Name)
	}
	return f.Present, nil
}

func (f FlagValue) Time() (time.Time, error) {
	if t, ok := f.Value.(time.Time); ok {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("flag --%s is not a time", f.Spec.Name)
}

func (f FlagValue) Int() (int, error) {
	if i, ok := f.Value.(int); ok {
		return i, nil
	}
	return 0, fmt.Errorf("flag --%s is not an integer", f.Spec.Name)
}

func parseNaturalTime(val string) (time.Time, error) {

	if timeStr, ok := strings.CutPrefix(val, "in "); ok {
		d, err := time.ParseDuration(timeStr)
		if err != nil {
			return time.Time{}, err
		}

		dueTime := time.Now().Add(d)
		return dueTime, nil

	} else if timeStr, ok := strings.CutPrefix(val, "tomorrow "); ok {
		t, err := time.Parse("15:04", timeStr)
		dueDate := time.Now().AddDate(0, 0, 1)

		if err != nil {
			return time.Time{}, err
		}

		return time.Date(
			dueDate.Year(), dueDate.Month(), dueDate.Day(),
			t.Hour(), t.Minute(), 0, 0, dueDate.Location(),
		), nil

	} else if timeStr, ok := strings.CutPrefix(val, "today "); ok {
		t, err := time.Parse("15:04", timeStr)
		dueDate := time.Now()

		if err != nil {
			return time.Time{}, err
		}
		return time.Date(
			dueDate.Year(), dueDate.Month(), dueDate.Day(),
			t.Hour(), t.Minute(), 0, 0, dueDate.Location(),
		), nil

	} else {
		return parseDayTime(val)
	}

}

func parseDayTime(val string) (time.Time, error) {
	today := time.Now().Weekday()
	now := time.Now()

	day, t, _ := strings.Cut(val, " ")
	var dueDate time.Time
	var diff float64

	switch day {
	case "Sunday":
		diff = math.Abs(float64((int(Sunday) - int(today))))
	case "Monday":
		diff = math.Abs(float64((int(Monday) - int(today))))
	case "Tuesday":
		diff = math.Abs(float64((int(Tuesday) - int(today))))
	case "Wednesday":
		diff = math.Abs(float64((int(Wednesday) - int(today))))
	case "Thursday":
		diff = math.Abs(float64((int(Thursday) - int(today))))
	case "Friday":
		diff = math.Abs(float64((int(Friday) - int(today))))
	case "Saturday":
		diff = math.Abs(float64((int(Saturday) - int(today))))
	}

	parsedTime, err := time.Parse("15:00", t)
	if err != nil {
		return time.Time{}, nil
	}

	if diff == 0 {
		dueDate = now.AddDate(0, 0, 7)

	} else {
		dueDate = now.Add(time.Duration((24 * diff)))
	}

	return time.Date(
		dueDate.Year(), dueDate.Month(), dueDate.Day(),
		parsedTime.Hour(), parsedTime.Minute(), 0, 0, dueDate.Location(),
	), nil

}
