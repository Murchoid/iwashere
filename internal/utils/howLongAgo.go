package utils

import (
	"fmt"
	"time"
)

func HowLongAgo(date time.Time) string {
	now := time.Now()
	timeDifference := now.Sub(date)
	howLongAgo := ""
	if timeDifference.Minutes() < 60 {
		howLongAgo = fmt.Sprintf("%.f minutes ago", timeDifference.Round(60000000000).Minutes())
	} else if timeDifference.Hours() < 24 {
		howLongAgo = fmt.Sprintf("%.f hours ago", timeDifference.Round(60).Hours())
	} else {
		days := timeDifference.Round(24).Hours() / 24
		howLongAgo = fmt.Sprintf("%.f days ago", days)
	}

	return howLongAgo
}
