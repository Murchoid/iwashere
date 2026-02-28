package utils

import (
	"fmt"
	"time"
)

func HowLongAgo(date time.Time) string {
	now := time.Now()
	timeDifference := now.Sub(date)

	if timeDifference.Minutes() < 60 {
		return fmt.Sprintf("%.0f minutes ago", timeDifference.Minutes())
	} else if timeDifference.Hours() < 24 {
		hours := int(timeDifference.Hours())
		minutes := int(timeDifference.Minutes()) % 60
		return fmt.Sprintf("%dhr %d mins ago", hours, minutes)
	} else {
		days := int(timeDifference.Hours() / 24)
		return fmt.Sprintf("%d days ago", days)
	}
}
