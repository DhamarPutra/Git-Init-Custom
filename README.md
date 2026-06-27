# Git Green Screen (`git-new`)

`git-new` is a Git subcommand and matching VSCode Extension designed to bootstrap new projects quickly from template repositories (local or remote), stripping the original template Git history and starting a fresh repository with a clean initial commit.

---

## 1. Product Requirements Document (PRD)

### Core Goals
- Provide a developer onboarding/repository bootstrap experience similar to `create-next-app` or `cargo new`.
- Act as a native Git subcommand (`git new`).
- Provide a companion VSCode Extension for visual repository bootstrapping.
- Keep the tool fast, cross-platform, and lightweight with no dependencies on external Git implementations (always uses the installed Git client).

### Workflow
1. **Resolve**: Resolve the template alias (from `~/.gitnew/config.yaml`), GitHub shorthand (`github:user/repo`), or local folder path.
2. **Clone/Copy**: Clone the remote template repository or recursively copy the local template directory.
3. **Strip History**: Safely remove the original `.git` history folder from the destination.
4. **Reinitialize**: Initialize a fresh Git repository in the target folder with the configured default branch (e.g. `main`).
5. **Initial Commit**: Automatically stage files and create an initial commit, leaving the project ready for its new remote origin.

---

## 2. Technical Stack

- **CLI Engine**: Go 1.26.2 (Native Executable)
- **Configuration Format**: YAML (`gopkg.in/yaml.v3`)
- **VSCode Extension**: TypeScript, Bundled with `esbuild`, utilizing Node's `child_process` to trigger the CLI binary.

---

## 3. Project Structure

```
├── cmd/
│   └── git-new/
│       └── main.go           # CLI Entry Point
├── internal/
│   ├── cli/                  # CLI Orchestrator and Argument Parser
│   ├── config/               # config.yaml Parser
│   ├── filesystem/           # Safe file operations
│   ├── git/                  # Wrapper around system git executable
│   ├── logger/               # Console output and Verbose logging helper
│   ├── registry/             # Template alias resolver
│   └── template/             # Template fetcher (cloning or copying)
├── VscodeExt/                # VSCode Extension source code
└── README.md                 # This documentation
```

---

## 4. Setup & Installation

### Prerequisite
Ensure Go 1.26+ and Git are installed on your machine and configured in your PATH environment variable.

### CLI Configuration Setup
Create a configuration file at `~/.gitnew/config.yaml` (Windows: `C:\Users\<Username>\.gitnew\config.yaml`) to specify aliases:
```yaml
defaultTemplate: starter
defaultBranch: main

templates:
  starter: https://github.com/DhamarPutra/template-project.git
  react: get from https://github.com/github/gitignore.git
```

---

## 5. Build & Run

### A. CLI (Go Subcommand)

#### Run Unit Tests
```bash
$env:GOROOT="PATH_TO_YOUR_GO_INSTALLATION\go"  # Adjust GOROOT if needed on your system
go test ./...
```

#### Compile/Build Binary
To build the native executable:
```bash
go build -o git-new.exe cmd/git-new/main.go
```

#### Make it a Native Git Subcommand
To use it as `git new <project>` from anywhere:
1. Compile the binary to `git-new.exe`.
2. Add the directory containing `git-new.exe` (e.g., `C:\Project\Git-Init-Custom`) to your Windows **Environment Variables PATH**:
   - Press the **Start** / Windows key, search for **"Environment Variables"** (or **"Edit the system environment variables"**), and select it.
   - Click the **"Environment Variables..."** button at the bottom of the window.
   - Under **User variables** (or **System variables**), locate and select the **`Path`** variable, then click **"Edit..."**.
   - Click **"New"** and add the absolute folder path where `git-new.exe` resides (e.g., `C:\Project\Git-Init-Custom`).
   - Click **OK** on all dialog windows to save the changes.
3. **Restart your terminal** or editor for the new PATH settings to take effect.
4. You can now run:
   ```bash
   git new my-app --template react
   ```

---

### B. VSCode Extension

#### Setup & Install Dependencies
Navigate into the `VscodeExt` folder and install dependencies:
```bash
cd VscodeExt
npm install
```

#### Build / Compile Extension
To compile the TypeScript code using `esbuild`:
```bash
npm run compile
```

#### Run in Debug Mode
1. Open the `VscodeExt` folder in VSCode.
2. Press `F5` to start debugging. A new **Extension Development Host** window will open.
3. Open the Command Palette (`Ctrl+Shift+P`) in the new window and run:
   `Git Green Screen: Bootstrap New Project`.

#### Pack into `.vsix` for Production Install
To package the extension into a local installer file:
```bash
npx @vscode/vsce package --no-dependencies
```
This produces `git-green-screen-0.0.1.vsix` inside the `VscodeExt/` folder. You can install it directly via the VSCode Extensions pane under **Install from VSIX...**.
