package template

import (
	"path/filepath"
	"strings"

	"github.com/DhamarPutra/git-new/internal/filesystem"
	"github.com/DhamarPutra/git-new/internal/git"
	"github.com/DhamarPutra/git-new/internal/logger"
)

// Fetcher handles obtaining template files.
type Fetcher struct {
	gitRunner *git.Runner
	logger    *logger.Logger
}

// NewFetcher creates a new Fetcher.
func NewFetcher(gitRunner *git.Runner, l *logger.Logger) *Fetcher {
	return &Fetcher{
		gitRunner: gitRunner,
		logger:    l,
	}
}

// Fetch template files from source to destination directory.
func (f *Fetcher) Fetch(source string, dest string) error {
	f.logger.Verbose("Resolving template source: %s", source)

	if f.IsLocal(source) {
		f.logger.Info("Copying local template from %s...", source)
		return filesystem.CopyDir(source, dest)
	}

	f.logger.Info("Cloning remote template from %s...", source)
	return f.gitRunner.Clone(source, dest)
}

// IsLocal checks if the source path represents a local directory.
func (f *Fetcher) IsLocal(source string) bool {
	// If it contains git/ssh scheme or http scheme, it's not local
	if strings.HasPrefix(source, "http://") ||
		strings.HasPrefix(source, "https://") ||
		strings.HasPrefix(source, "git@") ||
		strings.HasPrefix(source, "ssh://") {
		return false
	}

	// Clean path and check if it exists
	cleanSource := filepath.Clean(source)
	return filesystem.Exists(cleanSource)
}
