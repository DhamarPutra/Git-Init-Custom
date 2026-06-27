package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.DefaultTemplate != "starter" {
		t.Errorf("expected default template to be 'starter', got '%s'", cfg.DefaultTemplate)
	}
	if cfg.DefaultBranch != "main" {
		t.Errorf("expected default branch to be 'main', got '%s'", cfg.DefaultBranch)
	}
	if len(cfg.Templates) != 1 {
		t.Errorf("expected default templates map to have size 1, got size %d", len(cfg.Templates))
	}
	if cfg.Templates["starter"] != "templates" {
		t.Errorf("expected starter template to map to templates, got '%s'", cfg.Templates["starter"])
	}
}

func TestLoadFromFileNonExistent(t *testing.T) {
	cfg, err := LoadFromFile("non-existent-file-path-xyz.yaml")
	if err != nil {
		t.Fatalf("expected no error for non-existent file, got: %v", err)
	}
	if cfg.DefaultTemplate != "starter" {
		t.Errorf("expected default template to be 'starter', got '%s'", cfg.DefaultTemplate)
	}
}

func TestLoadFromFileValid(t *testing.T) {
	yamlContent := `
defaultTemplate: custom-starter
defaultBranch: develop
templates:
  custom-starter: https://github.com/DhamarPutra/template-project.git
  react: get from https://github.com/github/gitignore.git
`
	tmpDir, err := os.MkdirTemp("", "gitnew-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(tmpFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write temp config file: %v", err)
	}

	cfg, err := LoadFromFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to load config file: %v", err)
	}

	if cfg.DefaultTemplate != "custom-starter" {
		t.Errorf("expected defaultTemplate 'custom-starter', got '%s'", cfg.DefaultTemplate)
	}
	if cfg.DefaultBranch != "develop" {
		t.Errorf("expected defaultBranch 'develop', got '%s'", cfg.DefaultBranch)
	}
	if cfg.Templates["custom-starter"] != "https://github.com/DhamarPutra/template-project.git" {
		t.Errorf("expected template custom-starter URL match, got '%s'", cfg.Templates["custom-starter"])
	}
}
