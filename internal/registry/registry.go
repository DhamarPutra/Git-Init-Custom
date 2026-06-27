package registry

import (
	"fmt"
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

	// 5. Default fallback to raw-download for simple technology/template names
	capitalized := capitalize(name)
	return "raw-download://" + capitalized, nil
}

func capitalize(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
