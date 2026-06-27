# Project git-new (VsCode Ext)

A Git subcommand that bootstraps a new project from a Git template repository.

Instead of:

```bash
git clone template
cd template
git remote remove origin
```

Users can simply run:

```bash
git new my-app
```

or

```bash
git new my-app --template laravel
```

The command clones a template repository, removes the original Git history, initializes a fresh repository, and leaves the project ready to connect to a new remote.

---

# Goal

Provide a developer experience similar to:

* create-next-app
* npm create vite
* cargo new
* bun create

but built around Git templates.

The project should feel like a native Git command.

---

# Core Flow

Example:

```bash
git new my-app
```

Steps:

1. Resolve template.
2. Clone template into destination.
3. Remove original `.git` directory.
4. Initialize a fresh Git repository.
5. Create initial commit.
6. Remove all template remotes.
7. Ready for:

```bash
git remote add origin <repo>
git push -u origin main
```

---

# Features

## MVP

* git new <directory>
* template registry
* local template support
* GitHub template support
* fresh Git history
* initial commit
* verbose mode
* dry-run mode

---

## Future

* interactive template selector
* GitHub authentication
* private repositories
* template variables
* README generation
* license generation
* .env generation
* CI preset
* VSCode extension
* update template registry
* plugin system

---

# CLI

Examples:

```bash
git new my-app

git new my-app --template react

git new my-app --template laravel

git new my-app --template github:user/template

git new my-app --template ./local-template

git new my-app --dry-run

git new my-app --verbose
```

---

# Architecture

```
cmd/
    root

internal/
    cli/
    git/
    template/
    config/
    registry/
    logger/
    filesystem/

templates/

configs/
```

Responsibilities:

## cli

Argument parsing.

## git

Wrapper around Git commands.

Avoid shell-specific behavior.

## template

Template discovery.

Support:

* local path
* GitHub
* registry

## registry

Resolve template aliases.

Example:

```
react

↓

https://github.com/company/react-template.git
```

## filesystem

Safe file operations.

Never delete outside target directory.

---

# Design Principles

* Small executable.
* No external Git implementation.
* Always use installed Git.
* Cross-platform.
* No shell-specific features.
* Predictable behavior.
* Safe defaults.

---

# Error Handling

Never leave partially initialized repositories.

If any step fails:

* clean temporary directory
* restore working state whenever possible
* provide actionable error messages

---

# Logging

Support:

```
--verbose
```

Example:

```
Resolving template...
Cloning...
Removing Git history...
Initializing repository...
Creating initial commit...
Done.
```

---

# Configuration

User config:

```
~/.gitnew/config.yaml
```

Example:

```yaml
defaultTemplate: starter
defaultBranch: main

templates:
  starter: https://github.com/DhamarPutra/template-project.git
  react: get from https://github.com/github/gitignore.git
  nest: get from https://github.com/github/gitignore.git
  List-All: get from https://github.com/github/gitignore.git
```

---

# Non Goals

This project is NOT:

* a Git replacement
* a package manager
* a deployment tool
* a project generator

It only bootstraps repositories from templates.

---

# Coding Guidelines

* Keep functions small.
* Prefer composition.
* Avoid global state.
* Wrap all Git calls.
* Write unit tests for business logic.
* Write integration tests for CLI behavior.
* Never panic on user errors.
* Return structured errors.

---

# Testing

Unit tests:

* registry
* config
* template resolver

Integration tests:

* clone local template
* clone GitHub template
* initialize repository
* remove remote
* create first commit

---

# Future VSCode Integration

A VSCode extension may invoke:

```
git new
```

instead of calling `git init`.

The extension should not duplicate business logic.

All logic must remain inside the CLI.

The extension acts only as a UI layer.

---

# Success Criteria

Running:

```bash
git new awesome-project --template react
```

should produce:

```
awesome-project/
    .git
    src/
    package.json
    README.md
```

with:

* no template remote
* fresh Git history
* initial commit
* ready for a new remote
