package git

import (
	"os/exec"
	"strings"
)

type Info struct {
	Branch     string
	CommitHash string
	CommitMsg  string
	Remote     string
	HasChanges bool   // Unstaged changes
	HasStaged  bool   // Staged changes
	RootPath   string // Root of git repo
	UserEmail  string
	UserName   string
}

type Service struct {
	workDir string
}

func NewService(workDir string) *Service {
	return &Service{workDir: workDir}
}

// IsRepo checks if current directory is in a git repository
func (s *Service) IsRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = s.workDir
	err := cmd.Run()
	return err == nil
}

// GetInfo gathers all git information
func (s *Service) GetInfo() (*Info, error) {
	if !s.IsRepo() {
		return nil, nil // Not a git repo
	}

	info := &Info{}

	// Get current branch
	if branch, err := s.getBranch(); err == nil {
		info.Branch = branch
	}

	if name, _ := s.getCurrentUser(); name != "" {
		info.UserName = name
	}

	if _, email := s.getCurrentUser(); email != "" {
		info.UserEmail = email
	}

	// Get commit hash
	if hash, err := s.getCommitHash(); err == nil {
		info.CommitHash = hash
	}

	// Get commit message
	if msg, err := s.getCommitMsg(); err == nil {
		info.CommitMsg = msg
	}

	// Get remote
	if remote, err := s.getRemote(); err == nil {
		info.Remote = remote
	}

	// Check for changes
	info.HasChanges = s.hasChanges()
	info.HasStaged = s.hasStaged()

	// Get git root
	if root, err := s.getRoot(); err == nil {
		info.RootPath = root
	}

	return info, nil
}

// Helper methods
func (s *Service) getBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = s.workDir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (s *Service) getCommitHash() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	cmd.Dir = s.workDir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (s *Service) getCommitMsg() (string, error) {
	cmd := exec.Command("git", "log", "-1", "--pretty=%B")
	cmd.Dir = s.workDir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (s *Service) getRemote() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = s.workDir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (s *Service) getRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = s.workDir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (s *Service) hasChanges() bool {
	cmd := exec.Command("git", "diff", "--quiet")
	cmd.Dir = s.workDir
	err := cmd.Run()
	// Exit code 1 means there are changes
	return err != nil
}

func (s *Service) hasStaged() bool {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	cmd.Dir = s.workDir
	err := cmd.Run()
	// Exit code 1 means there are staged changes
	return err != nil
}

func (s *Service) getCurrentUser() (string, string) {
	// Get from git config
	nameCmd := exec.Command("git", "config", "user.name")
	emailCmd := exec.Command("git", "config", "user.email")

	name, _ := nameCmd.Output()
	email, _ := emailCmd.Output()

	return strings.TrimSpace(string(name)), strings.TrimSpace(string(email))
}

// GetModifiedFiles returns list of files changed
func (s *Service) GetModifiedFiles() ([]string, error) {
	cmd := exec.Command("git", "ls-files", "--modified", "--others", "--exclude-standard")
	cmd.Dir = s.workDir
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(files) == 1 && files[0] == "" {
		return []string{}, nil
	}
	return files, nil
}
