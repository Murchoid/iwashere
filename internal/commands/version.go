package commands

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// These will be set during build
var (
	// These can still be set by ldflags as fallback
	Version = "dev"
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
    ║        🔖 iwashere %s            ║
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
	fmt.Printf("iwashere version %s\n", GetInfo())
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
	// 1. If we have a real version from ldflags, use it
	if Version != "dev" && Version != "" {
		return Version
	}

	// 2. Try to get from git describe
	if v := getFromGitDescribe(); v != "" {
		return v
	}

	// 3. Try to get from git tag
	if v := getFromGitTag(); v != "" {
		return v
	}

	// 4. Try to get from git commit
	if v := getFromGitCommit(); v != "" {
		return v + "-dev"
	}

	// 5. Fallback
	return "dev"
}

// GetCommit returns the git commit hash
func GetCommit() string {
	if Commit != "none" {
		return Commit
	}

	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
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

func getFromGitDescribe() string {
	cmd := exec.Command("git", "describe", "--tags", "--always")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func getFromGitTag() string {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func getFromGitCommit() string {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
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
