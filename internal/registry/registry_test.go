package registry

import (
	"os"
	"path/filepath"
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

func TestResolveLocalGitIgnore(t *testing.T) {
	tmpDir := t.TempDir()

	gitignoreDir := filepath.Join(tmpDir, "gitignore")
	if err := os.MkdirAll(gitignoreDir, 0755); err != nil {
		t.Fatalf("failed to create temp gitignore dir: %v", err)
	}

	mockFilePath := filepath.Join(gitignoreDir, "Laravel.gitignore")
	if err := os.WriteFile(mockFilePath, []byte("laravel patterns"), 0644); err != nil {
		t.Fatalf("failed to write mock gitignore file: %v", err)
	}

	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	defer os.Chdir(oldCwd)

	cfg := config.DefaultConfig()
	resolver := NewResolver(cfg)

	res, err := resolver.Resolve("laravel")
	if err != nil {
		t.Fatalf("expected resolution to succeed, got: %v", err)
	}

	absExpectedPath, _ := filepath.Abs(mockFilePath)
	expectedPrefix := "gitignore://" + absExpectedPath
	if res != expectedPrefix {
		t.Errorf("expected resolved path '%s', got '%s'", expectedPrefix, res)
	}
}

func TestResolveLocalTemplate(t *testing.T) {
	tmpDir := t.TempDir()

	templatesDir := filepath.Join(tmpDir, "templates")
	targetTemplateDir := filepath.Join(templatesDir, "my-test-template")
	if err := os.MkdirAll(targetTemplateDir, 0755); err != nil {
		t.Fatalf("failed to create temp template dir: %v", err)
	}

	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	defer os.Chdir(oldCwd)

	cfg := config.DefaultConfig()
	resolver := NewResolver(cfg)

	res, err := resolver.Resolve("my-test-template")
	if err != nil {
		t.Fatalf("expected resolution to succeed, got: %v", err)
	}

	absExpectedPath, _ := filepath.Abs(targetTemplateDir)
	if res != absExpectedPath {
		t.Errorf("expected resolved path '%s', got '%s'", absExpectedPath, res)
	}
}


