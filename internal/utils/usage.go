package utils

import (
	"fmt"
)

func PrintCommandHelp(cmdName, description, usage string, examples []string) {
	fmt.Printf("%s%s iwashere %s - %s%s%s\n", ColorBlue, ColorReset, cmdName, ColorYellow, description, ColorReset)
	fmt.Println()
	fmt.Printf("%sUSAGE:%s\n", ColorGreen, ColorReset)
	fmt.Printf("  %s\n", usage)
	fmt.Println()
	if len(examples) > 0 {
		fmt.Printf("%sEXAMPLES:%s\n", ColorGreen, ColorReset)
		for _, ex := range examples {
			fmt.Printf("  %s•%s %s\n", ColorGreen, ColorReset, ex)
		}
	}
}
