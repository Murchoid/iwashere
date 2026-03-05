package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// Duration is a custom type that handles both string and number formats
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

// MarshalJSON implements custom JSON marshaling
func (d Duration) MarshalJSON() ([]byte, error) {
	// Always marshal as string for readability
	return json.Marshal(time.Duration(d).String())
}

// Duration returns the time.Duration value
func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}
