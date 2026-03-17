package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// Duration will be our custom type that handles both string and number formats
type Duration time.Duration

// UnmarshalJSON implements custom JSON unmarshaling
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	switch value := v.(type) {
	case float64:
		// Number format (nanoseconds)
		*d = Duration(time.Duration(value))
		return nil
	case string:
		// String format like "1h2m3s"
		dur, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(dur)
		return nil
	default:
		return fmt.Errorf("invalid duration format: %v", value)
	}
}

func (d Duration) MarshalJSON() ([]byte, error) {

	return json.Marshal(time.Duration(d).String())
}


func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}
