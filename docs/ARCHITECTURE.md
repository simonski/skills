# Architecture

## Overview

`skills` is a single-binary Go CLI that manages AI-agent skill definition files (`.skills/<id>.md`) in a project directory. Skills are Markdown files with YAML front matter, embedded directly into the binary at compile time.

## Package structure

```
skills/
├── main.go                     # entry point — calls cmd.Execute()
├── cmd/                        # cobra command implementations
│   ├── root.go                 # root command, checkForUpdates, command registration
│   ├── add.go                  # skills add
│   ├── rm.go                   # skills rm
│   ├── ls.go                   # skills ls
│   ├── get.go                  # skills get
│   ├── search.go               # skills search
│   ├── versions.go             # skills versions
│   ├── update.go               # skills update
│   ├── init.go                 # skills init (interactive wizard)
│   └── version.go              # skills version
├── internal/
│   ├── catalog/                # embedded skill catalog (read-only)
│   │   ├── catalog.go          # Get, All, Versions, GetVersion API
│   │   └── skills/             # embedded skill files
│   │       └── <id>/
│   │           ├── <ver>.md    # versioned skill content (YAML front matter + body)
│   │           ├── README.md   # human-readable overview
│   │           └── CHANGELOG.md
│   ├── project/                # local project skill management (read/write)
│   │   └── project.go          # Install, Remove, List, SkillPath
│   └── version/                # GitHub release checking
│       └── version.go          # LatestRelease, IsOutdated, IsNewerThan
└── packaging/
    └── homebrew/
        └── skills.rb.template  # Homebrew formula template
```

## Dependency DAG

```
cmd → internal/catalog    (read skill definitions)
cmd → internal/project    (read/write .skills/ directory)
cmd → internal/version    (check for newer releases)

internal/* packages do NOT import cmd or each other.
```

No circular dependencies. The internal packages are fully independent.

## Catalog

Skills are embedded into the binary at compile time using `//go:embed skills/*/*.md` in `internal/catalog/catalog.go`. This means:

- Zero disk reads for catalog access at runtime
- The catalog cannot be tampered with by unprivileged users
- Adding a skill requires recompiling the binary

Each skill file uses YAML front matter delimited by `---`:

```markdown
---
id: go
version: 1.0.0
description: Go programming best practices
---

# Go

Skill content...
```

`catalog.Get(id)` returns the latest version. `catalog.GetVersion(id, version)` returns a specific version. `catalog.Versions(id)` lists all available versions sorted ascending.

## Project (local skill storage)

Installed skills live in `.skills/<id>.md` within the working directory. `internal/project` handles all filesystem I/O:

- `project.Install(root, skill)` — writes the skill file
- `project.Remove(root, id)` — deletes the skill file
- `project.List(root)` — reads all installed skills

File paths are constructed with `filepath.Join(root, ".skills", id+".md")` to prevent path traversal.

## Version checking

`internal/version.LatestRelease()` makes a single GET request to the GitHub releases API on every command invocation (via `checkForUpdates()` in `cmd/root.go`). If `GITHUB_TOKEN` is set in the environment, it is sent as a Bearer token to avoid the 60 req/hr unauthenticated rate limit.

> ⚠️ Known issue: this call is synchronous and can add up to 5s latency on slow networks (tracked in SKL-T-12 / SKL-T-36).

## Build and release

- `make build` — compiles the binary, **bumps the patch version** in `VERSION`
- `make release` — cross-compiles for all platforms, does **not** bump version
- `make publish` — calls `release`, tags, creates GitHub release, updates Homebrew tap
- CI uses `make release` (not `make build`) to avoid double-bumping

See [CONTRIBUTING.md](../CONTRIBUTING.md) for the full development workflow.
