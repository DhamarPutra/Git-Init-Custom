package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/DhamarPutra/git-new/internal/logger"
)

// Runner runs git commands.
type Runner struct {
	logger *logger.Logger
}

// NewRunner creates a new git runner.
func NewRunner(l *logger.Logger) *Runner {
	return &Runner{logger: l}
}

// runCmd executes a git command in the specified directory.
func (r *Runner) runCmd(dir string, args ...string) (string, error) {
	r.logger.Verbose("Running: git %s", strings.Join(args, " "))

	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg == "" {
			errMsg = err.Error()
		}
		return "", fmt.Errorf("git %s failed: %s", args[0], errMsg)
	}

	return strings.TrimSpace(stdout.String()), nil
}

// Clone clones a repository to a target directory.
func (r *Runner) Clone(url, targetDir string) error {
	_, err := r.runCmd("", "clone", url, targetDir)
	return err
}

// Init initializes a new git repository.
func (r *Runner) Init(dir string, defaultBranch string) error {
	args := []string{"init"}
	if defaultBranch != "" {
		args = append(args, "-b", defaultBranch)
	}
	_, err := r.runCmd(dir, args...)
	return err
}

// AddAll adds all files to the staging area.
func (r *Runner) AddAll(dir string) error {
	_, err := r.runCmd(dir, "add", ".")
	return err
}

// Commit creates a new commit with the specified message.
func (r *Runner) Commit(dir string, message string) error {
	_, err := r.runCmd(dir, "commit", "-m", message)
	return err
}

// SetUserConfig sets local git user details if not set globally, preventing commit failures.
func (r *Runner) EnsureUserConfig(dir string) error {
	// Check if user.name is set
	name, _ := r.runCmd(dir, "config", "user.name")
	if name == "" {
		_, err := r.runCmd(dir, "config", "user.name", "git-new")
		if err != nil {
			return err
		}
	}
	// Check if user.email is set
	email, _ := r.runCmd(dir, "config", "user.email")
	if email == "" {
		_, err := r.runCmd(dir, "config", "user.email", "git-new@localhost")
		if err != nil {
			return err
		}
	}
	return nil
}
