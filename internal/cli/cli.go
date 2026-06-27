package cli

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/DhamarPutra/git-new/internal/config"
	"github.com/DhamarPutra/git-new/internal/filesystem"
	"github.com/DhamarPutra/git-new/internal/git"
	"github.com/DhamarPutra/git-new/internal/logger"
	"github.com/DhamarPutra/git-new/internal/registry"
	"github.com/DhamarPutra/git-new/internal/template"
)

// Options holds CLI configurations.
type Options struct {
	Template string
	Branch   string
	DryRun   bool
	Verbose  bool
	DestDir  string
}

// Run executes the subcommand logic.
func Run(args []string) int {
	// 1. Load User Configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		return 1
	}

	// 2. Parse Flags
	fs := flag.NewFlagSet("git-new", flag.ContinueOnError)
	var opts Options
	fs.StringVar(&opts.Template, "template", cfg.DefaultTemplate, "Template name, alias, URL, or local path")
	fs.StringVar(&opts.Template, "t", cfg.DefaultTemplate, "Template name, alias, URL, or local path (shorthand)")
	fs.StringVar(&opts.Branch, "branch", cfg.DefaultBranch, "Default branch name for the new repository")
	fs.StringVar(&opts.Branch, "b", cfg.DefaultBranch, "Default branch name for the new repository (shorthand)")
	fs.BoolVar(&opts.DryRun, "dry-run", false, "Show what would be done without modifying the disk")
	fs.BoolVar(&opts.Verbose, "verbose", false, "Enable verbose output")
	fs.BoolVar(&opts.Verbose, "v", false, "Enable verbose output (shorthand)")

	// Parse arguments (reordered so flags can be placed in any order)
	if err := fs.Parse(reorderArgs(args)); err != nil {
		return 1
	}

	// The positional argument is the destination directory
	remainingArgs := fs.Args()
	if len(remainingArgs) < 1 {
		fmt.Println("Usage: git new <directory> [options]")
		fmt.Println("\nOptions:")
		fs.PrintDefaults()
		return 1
	}
	opts.DestDir = remainingArgs[0]

	// 3. Initialize Logger
	l := logger.New(opts.Verbose)

	// 4. Resolve Template Source
	resolver := registry.NewResolver(cfg)
	resolvedSource, err := resolver.Resolve(opts.Template)
	if err != nil {
		l.Error("Failed to resolve template: %v", err)
		return 1
	}

	if opts.DryRun {
		l.Info("[Dry Run] Would bootstrap new repository:")
		l.Info("[Dry Run]   Destination: %s", opts.DestDir)
		l.Info("[Dry Run]   Template source: %s", resolvedSource)
		l.Info("[Dry Run]   Default branch: %s", opts.Branch)
		return 0
	}

	// Verify Destination Directory
	destExists := filesystem.Exists(opts.DestDir)
	if destExists {
		isEmpty, err := filesystem.IsEmpty(opts.DestDir)
		if err != nil {
			l.Error("Failed to read destination directory: %v", err)
			return 1
		}
		if !isEmpty {
			l.Error("Destination directory '%s' is not empty.", opts.DestDir)
			return 1
		}
	}

	// 5. Bootstrap
	gitRunner := git.NewRunner(l)
	fetcher := template.NewFetcher(gitRunner, l)

	isRawDownload := strings.HasPrefix(resolvedSource, "raw-download://")
	isGitIgnoreTemplate := strings.HasPrefix(resolvedSource, "gitignore://")

	if isRawDownload {
		techName := strings.TrimPrefix(resolvedSource, "raw-download://")
		l.Info("Bootstrapping from DhamarPutra/Git-Init-Custom...")

		// Ensure destination exists
		if err := os.MkdirAll(opts.DestDir, 0755); err != nil {
			l.Error("Failed to create destination directory: %v", err)
			return 1
		}

		// Helper to download a URL to a local filepath
		downloadFile := func(url, filePath string) error {
			l.Verbose("Downloading %s...", url)
			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("failed to fetch %s (status: %d)", url, resp.StatusCode)
			}

			// Ensure parent directory exists
			if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
				return err
			}

			out, err := os.Create(filePath)
			if err != nil {
				return err
			}
			defer out.Close()

			_, err = io.Copy(out, resp.Body)
			return err
		}

		// 1. Download templates/log.md
		logURL := "https://raw.githubusercontent.com/DhamarPutra/Git-Init-Custom/main/templates/log.md"
		if err := downloadFile(logURL, filepath.Join(opts.DestDir, "log.md")); err != nil {
			l.Error("Failed to download template files: %v", err)
			if !destExists {
				filesystem.SafeRemoveAll(opts.DestDir)
			}
			return 1
		}

		// 2. Download templates/.github/workflows/config.yml
		workflowURL := "https://raw.githubusercontent.com/DhamarPutra/Git-Init-Custom/main/templates/.github/workflows/config.yml"
		if err := downloadFile(workflowURL, filepath.Join(opts.DestDir, ".github", "workflows", "config.yml")); err != nil {
			l.Error("Failed to download workflow files: %v", err)
			if !destExists {
				filesystem.SafeRemoveAll(opts.DestDir)
			}
			return 1
		}

		// 3. Download gitignore/<Tech>.gitignore if requested (and is not empty/starter)
		if techName != "" && strings.ToLower(techName) != "starter" {
			gitignoreURL := fmt.Sprintf("https://raw.githubusercontent.com/DhamarPutra/Git-Init-Custom/main/gitignore/%s.gitignore", techName)
			if err := downloadFile(gitignoreURL, filepath.Join(opts.DestDir, ".gitignore")); err != nil {
				l.Error("Failed to download .gitignore for %s: %v", techName, err)
				if !destExists {
					filesystem.SafeRemoveAll(opts.DestDir)
				}
				return 1
			}
		}
	} else if isGitIgnoreTemplate {
		// Copy default template first so the user doesn't lose standard files (e.g. .github, log.md)
		defaultTemplateSrc, err := resolver.Resolve(cfg.DefaultTemplate)
		if err == nil && defaultTemplateSrc != "" && !strings.HasPrefix(defaultTemplateSrc, "gitignore://") {
			l.Verbose("Copying default template files from %s...", defaultTemplateSrc)
			if err := fetcher.Fetch(defaultTemplateSrc, opts.DestDir); err != nil {
				l.Error("Failed to copy default template files: %v", err)
				return 1
			}
		}

		gitIgnoreSrc := strings.TrimPrefix(resolvedSource, "gitignore://")
		l.Info("Extracting local .gitignore template from %s...", filepath.Base(gitIgnoreSrc))

		// Create destination directory if it doesn't exist
		if err := os.MkdirAll(opts.DestDir, 0755); err != nil {
			l.Error("Failed to create destination directory: %v", err)
			return 1
		}

		// Copy the file as .gitignore
		destGitIgnore := filepath.Join(opts.DestDir, ".gitignore")
		if err := filesystem.CopyFile(gitIgnoreSrc, destGitIgnore); err != nil {
			l.Error("Failed to copy .gitignore: %v", err)
			if !destExists {
				filesystem.SafeRemoveAll(opts.DestDir)
			}
			return 1
		}
	} else {
		// Fetch template to destination directory
		if err := fetcher.Fetch(resolvedSource, opts.DestDir); err != nil {
			l.Error("Failed to fetch template: %v", err)
			// Clean up partial directory if it was created
			if !destExists {
				l.Verbose("Cleaning up target directory %s due to failure...", opts.DestDir)
				filesystem.SafeRemoveAll(opts.DestDir)
			}
			return 1
		}

		// Remove original .git directory inside the template
		gitDir := filepath.Join(opts.DestDir, ".git")
		if filesystem.Exists(gitDir) {
			l.Verbose("Removing Git history...")
			if err := filesystem.SafeRemoveAll(gitDir); err != nil {
				l.Error("Failed to clean template git history: %v", err)
				return 1
			}
		}
	}

	// Initialize a fresh Git repository
	l.Verbose("Initializing repository...")
	if err := gitRunner.Init(opts.DestDir, opts.Branch); err != nil {
		l.Error("Failed to initialize git repository: %v", err)
		return 1
	}

	// Add files to stage
	if err := gitRunner.AddAll(opts.DestDir); err != nil {
		l.Error("Failed to stage files: %v", err)
		return 1
	}

	// Ensure local user.name and user.email are configured if global aren't
	if err := gitRunner.EnsureUserConfig(opts.DestDir); err != nil {
		l.Verbose("Warning: Failed to ensure git user config: %v", err)
	}

	// Create initial commit
	l.Verbose("Creating initial commit...")
	commitMsg := fmt.Sprintf("Initial commit from template (%s)", opts.Template)
	if err := gitRunner.Commit(opts.DestDir, commitMsg); err != nil {
		l.Error("Failed to create initial commit: %v", err)
		return 1
	}

	l.Info("Done.")
	return 0
}

// reorderArgs moves all flags (starting with -) and their values to the front
// so the standard Go flag package can parse them in any position.
func reorderArgs(args []string) []string {
	var flags []string
	var positionals []string

	i := 0
	for i < len(args) {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
			// Check if this flag requires a value
			name := strings.TrimLeft(arg, "-")
			requiresValue := name == "template" || name == "t" || name == "branch" || name == "b"
			if requiresValue && i+1 < len(args) {
				flags = append(flags, args[i+1])
				i += 2
			} else {
				i++
			}
		} else {
			positionals = append(positionals, arg)
			i++
		}
	}
	return append(flags, positionals...)
}
