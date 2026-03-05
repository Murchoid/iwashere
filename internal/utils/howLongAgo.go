package utils

import (
	"fmt"
	"time"
)

func HowLongAgo(since time.Time, duration time.Duration) string {
	duration += time.Since(since)

	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		minutes := int(duration.Minutes())
		return fmt.Sprintf("%dm ago", minutes)
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		return fmt.Sprintf("%dh ago", hours)
	case duration < 7*24*time.Hour:
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	default:
		return since.Add(duration).Format("Jan 2")
	}
}
