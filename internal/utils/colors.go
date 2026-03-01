package utils

// Simple Color support (disable on Windows if needed)
var (
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[37m"
)

func init() {
	if !UseColors() {
		// Windows might not support ANSI Colors in cmd
		ColorGreen = ""
		ColorYellow = ""
		ColorBlue = ""
		ColorReset = ""
		ColorRed = ""
		ColorPurple = ""
		ColorCyan = ""
		ColorGray = ""
	}
}
