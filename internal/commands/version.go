package commands

import (
	"fmt"
	"runtime"
	"time"
)

// These will be set during build
var (
    Version   = "dev"
    Commit    = "none"
    BuildDate = "unknown"
    BuiltBy   = "unknown"
)

type VersionCommand struct{}

func NewVersionCommand() Command {
	return &VersionCommand{}
}
func (c *VersionCommand) Name() string { return "version" }
func (c *VersionCommand) Description() string { return "Show iwashere version information" }
func (c *VersionCommand) Usage() string { return "iwashere version" }
func (c *VersionCommand) Examples() []string { return []string{"iwashere version", "iwashere -v"} } 

func (c *VersionCommand) Execute(ctx *Context) error {
    printVersionInfo()
    return nil
}

func printVersionInfo() {
    banner := `
    ╔══════════════════════════════════════╗
    ║         🔖 iwashere %s           ║
    ║     Context Preservation Tool        ║
    ╚══════════════════════════════════════╝
    `
    fmt.Printf(banner, Version)
    
    // Fun facts
    facts := []string{
        "Built with Go",
        "Never lose your train of thought",
        "Your context, your control",
        "Fast. Simple. Effective.",
        "Dogfooding since day one",
        "Made with ❤️  in Kenya",
    }
    
    fmt.Println()
    fmt.Println("  " + facts[getRandomFactIndex(len(facts))])
    fmt.Println()
    
    // Detailed info
    fmt.Println("Build Details:")
    fmt.Printf("   • Version: %s\n", Version)
    fmt.Printf("   • Commit: %s\n", truncate(Commit, 8))
    fmt.Printf("   • Built: %s\n", BuildDate)
    fmt.Printf("   • Built by: %s\n", BuiltBy)
    fmt.Printf("   • Go version: %s\n", runtime.Version())
    fmt.Printf("   • OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
    fmt.Println()
    
    // Quick tips
    fmt.Println("Quick tips:")
    fmt.Println("   • iwashere init     # Start tracking a project")
    fmt.Println("   • iwashere add msg  # Save your context")
    fmt.Println("   • iwashere status   # See where you left off")
    fmt.Println()
    
    // Easter egg if version is dev
    if Version == "dev" {
        fmt.Println("Development build - you're probably hacking on iwashere!")
        fmt.Println("Run 'go build -o iwashere ./cmd/iwashere' to build")
    }
}

func GetVersion() string {
	return Version
}

// Simple pseudo-random for facts
func getRandomFactIndex(max int) int {
    return int(time.Now().UnixNano() % int64(max))
}

func init() {
	Register("version", NewVersionCommand)
}