package registry

import (
	"testing"

	"github.com/DhamarPutra/git-new/internal/config"
)

func TestResolve(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Templates["react"] = "get from https://github.com/github/gitignore.git"
	cfg.Templates["nest"] = "https://github.com/github/gitignore.git"

	tempDir := t.TempDir()
	resolver := NewResolver(cfg)

	tests := []struct {
		name     string
		input    string
		expected string
		err      bool
	}{
		{
			name:     "Config alias with prefix clean",
			input:    "react",
			expected: "https://github.com/github/gitignore.git",
			err:      false,
		},
		{
			name:     "Config alias without prefix",
			input:    "nest",
			expected: "https://github.com/github/gitignore.git",
			err:      false,
		},
		{
			name:     "GitHub shorthand",
			input:    "github:DhamarPutra/git-new",
			expected: "https://github.com/DhamarPutra/git-new.git",
			err:      false,
		},
		{
			name:     "Direct HTTPS URL",
			input:    "https://github.com/some/repo.git",
			expected: "https://github.com/some/repo.git",
			err:      false,
		},
		{
			name:     "Local relative path",
			input:    tempDir,
			expected: tempDir,
			err:      false,
		},
		{
			name:     "Simple name falls back to raw-download",
			input:    "laravel",
			expected: "raw-download://Laravel",
			err:      false,
		},
		{
			name:  "Empty template name",
			input: "",
			err:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := resolver.Resolve(tt.input)
			if (err != nil) != tt.err {
				t.Fatalf("expected error presence %v, got %v", tt.err, err)
			}
			if !tt.err && res != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, res)
			}
		})
	}
}
