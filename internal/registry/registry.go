package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/DhamarPutra/git-new/internal/config"
	"github.com/DhamarPutra/git-new/internal/filesystem"
)

// Resolver resolves template names/aliases to actual URLs or paths.
type Resolver struct {
	config *config.Config
}

// NewResolver creates a new Resolver.
func NewResolver(cfg *config.Config) *Resolver {
	return &Resolver{config: cfg}
}

// Resolve translates a template name or alias to a git URL or local path.
func (r *Resolver) Resolve(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("template name cannot be empty")
	}

	// 1. Check if the name exists in config templates mapping
	if val, ok := r.config.Templates[name]; ok {
		// Clean "get from " prefix if user added it in config.yaml
		cleaned := strings.TrimPrefix(val, "get from ")
		cleaned = strings.TrimSpace(cleaned)
		return cleaned, nil
	}

	// 2. Check if it is a GitHub shorthand (e.g., github:owner/repo)
	if strings.HasPrefix(name, "github:") {
		parts := strings.SplitN(name, ":", 2)
		if len(parts) == 2 && strings.Contains(parts[1], "/") {
			return fmt.Sprintf("https://github.com/%s.git", parts[1]), nil
		}
		return "", fmt.Errorf("invalid github template shorthand format: %q (expected github:owner/repository)", name)
	}

	// 3. Check if it is a direct remote URL
	if strings.HasPrefix(name, "http://") ||
		strings.HasPrefix(name, "https://") ||
		strings.HasPrefix(name, "git@") ||
		strings.HasPrefix(name, "ssh://") {
		return name, nil
	}

	// 4. Check if it exists as a local path
	if filesystem.Exists(name) {
		return name, nil
	}

	// 5. Check if it exists as a local folder inside a "templates" directory
	if templatesPath, err := r.findLocalTemplate(name); err == nil && templatesPath != "" {
		return templatesPath, nil
	}

	// 6. Check if it exists as a local .gitignore template inside a "gitignore" folder
	if gitignorePath, err := r.findLocalGitIgnore(name); err == nil && gitignorePath != "" {
		return "gitignore://" + gitignorePath, nil
	}

	// 7. If it's a simple name but not found in config, disk, or templates/gitignore folders
	return "", fmt.Errorf("template %q not found in config.yaml, does not exist as a local path, and is not a valid repository URL", name)
}

// findLocalTemplate looks for <name> folder in the templates directory.
func (r *Resolver) findLocalTemplate(name string) (string, error) {
	var dirsToSearch []string

	// Check relative to executable
	if exePath, err := os.Executable(); err == nil {
		dirsToSearch = append(dirsToSearch, filepath.Join(filepath.Dir(exePath), "templates"))
	}
	// Check relative to current working directory
	dirsToSearch = append(dirsToSearch, filepath.Join(".", "templates"))

	for _, dir := range dirsToSearch {
		targetDir := filepath.Join(dir, name)
		if filesystem.Exists(targetDir) {
			if info, err := os.Stat(targetDir); err == nil && info.IsDir() {
				absPath, err := filepath.Abs(targetDir)
				if err != nil {
					return targetDir, nil
				}
				return absPath, nil
			}
		}
	}

	return "", fmt.Errorf("local template folder not found")
}

// findLocalGitIgnore looks for <name>.gitignore (case-insensitive) inside the gitignore directory.
func (r *Resolver) findLocalGitIgnore(name string) (string, error) {
	var dirsToSearch []string

	// Check relative to executable
	if exePath, err := os.Executable(); err == nil {
		dirsToSearch = append(dirsToSearch, filepath.Join(filepath.Dir(exePath), "gitignore"))
	}
	// Check relative to current working directory
	dirsToSearch = append(dirsToSearch, filepath.Join(".", "gitignore"))

	for _, dir := range dirsToSearch {
		if !filesystem.Exists(dir) {
			continue
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if strings.EqualFold(entry.Name(), name+".gitignore") {
				absPath, err := filepath.Abs(filepath.Join(dir, entry.Name()))
				if err != nil {
					return filepath.Join(dir, entry.Name()), nil
				}
				return absPath, nil
			}
		}
	}

	return "", fmt.Errorf("gitignore file not found")
}
