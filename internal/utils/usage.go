package utils

import (
    "fmt"
    "runtime"
)

// Simple color support (disable on Windows if needed)
var (
    Green  = "\033[32m"
    Yellow = "\033[33m"
    Blue   = "\033[34m"
    Reset  = "\033[0m"
)

func init() {
    if runtime.GOOS == "windows" {
        // Windows might not support ANSI colors in cmd
        Green = ""
        Yellow = ""
        Blue = ""
        Reset = ""
    }
}

func PrintCommandHelp(cmdName, description, usage string, examples []string) {
    fmt.Printf("%s%s iwashere %s - %s%s%s\n", Blue, Reset, cmdName, Yellow, description, Reset)
    fmt.Println()
    fmt.Printf("%sUSAGE:%s\n", Green, Reset)
    fmt.Printf("  %s\n", usage)
    fmt.Println()
    if len(examples) > 0 {
        fmt.Printf("%sEXAMPLES:%s\n", Green, Reset)
        for _, ex := range examples {
            fmt.Printf("  %s•%s %s\n", Green, Reset, ex)
        }
    }
}