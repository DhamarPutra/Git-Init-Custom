package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsLocal(t *testing.T) {
	fetcher := &Fetcher{}

	// Test URL paths
	if fetcher.IsLocal("https://github.com/user/repo") {
		t.Error("expected https URL not to be local")
	}
	if fetcher.IsLocal("http://github.com/user/repo") {
		t.Error("expected http URL not to be local")
	}
	if fetcher.IsLocal("git@github.com:user/repo.git") {
		t.Error("expected git ssh URL not to be local")
	}

	// Test non-existent path
	if fetcher.IsLocal("./non-existent-template-dir-xyz") {
		t.Error("expected non-existent path not to be local")
	}

	// Test existing path
	tmpDir, err := os.MkdirTemp("", "gitnew-template-test")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	if !fetcher.IsLocal(tmpDir) {
		t.Error("expected existing temp directory to be local")
	}

	// Test relative path
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	relPath, err := filepath.Rel(cwd, tmpDir)
	if err == nil && !fetcher.IsLocal(relPath) {
		t.Error("expected existing relative path to be local")
	}
}
