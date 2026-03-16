package commands

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// These will be set during build
var (
	Version = "0.3.0"
	Commit  = "none"
)

type VersionCommand struct{}

func NewVersionCommand() Command {
	return &VersionCommand{}
}
func (c *VersionCommand) Name() string        { return "version" }
func (c *VersionCommand) Description() string { return "Show iwashere version information" }
func (c *VersionCommand) Usage() string       { return "iwashere version" }
func (c *VersionCommand) Examples() []string  { return []string{"iwashere version", "iwashere -v"} }

func (c *VersionCommand) Execute(ctx *Context) error {
	Version = GetVersion()
	printVersionInfo()
	return nil
}

func printVersionInfo() {
	banner := `
    ╔══════════════════════════════════════╗
    ║        🔖 iwashere v%s            ║
    ║     Context Preservation Tool        ║
    ╚══════════════════════════════════════╝
    `
	fmt.Printf(banner, GetVersion())

	// Fun facts
	facts := []string{
		"Built with Go",
		"Never lose your train of thought",
		"Your context, your control",
		"Fast. Simple. Effective.",
		"Mkenya daima",
	}

	fmt.Println()
	fmt.Println("  " + facts[getRandomFactIndex(len(facts))])
	fmt.Println()

	// Detailed info
	fmt.Printf("iwashere version v%s\n", GetInfo())
	fmt.Printf("  built with %s\n", runtime.Version())
	fmt.Printf("  on %s/%s\n", runtime.GOOS, runtime.GOARCH)
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

// Get returns the version string, trying multiple strategies
func GetVersion() string {

	if Version != "dev" && Version != "" {
		b, _, _ := strings.Cut(Version, "-")
		Version = b
		return Version
	}

	return "dev"
}

func GetCommit() string {
	return Commit
}

// GetInfo returns formatted version info
func GetInfo() string {
	version := GetVersion()
	commit := GetCommit()

	if commit != "unknown" && commit != "none" {
		return fmt.Sprintf("%s (%s)", version, commit)
	}
	return version
}

// IsReleaseBuild returns true if this is a proper release
func IsReleaseBuild() bool {
	return Version != "dev" && Version != ""
}

// Simple pseudo-random for facts
func getRandomFactIndex(max int) int {
	return int(time.Now().UnixNano() % int64(max))
}

func init() {
	Register("version", NewVersionCommand)
}
