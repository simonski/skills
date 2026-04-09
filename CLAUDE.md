# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

Read AGENTS.md

## Build & Test Commands

```bash
make test              # run all tests (go test ./...)
make lint              # go vet + staticcheck
make setup             # install dev tools (staticcheck, govulncheck)
go test ./internal/catalog/...   # run tests for a single package
go test -run TestName ./...      # run a single test by name
go build -o dist/skills .       # local build WITHOUT bumping VERSION
make build             # build binary (auto-bumps patch version — don't commit VERSION from local builds)
make release           # cross-compile release archives for all platforms
```

## Architecture

Single-binary Go CLI built with **cobra**. Entry point: `main.go` → `cmd.Execute()`.

### Package layout

- **`cmd/`** — Cobra command definitions (one file per subcommand: `add.go`, `ls.go`, `rm.go`, etc.). `root.go` wires them together and runs an auto-update check via `internal/version`.
- **`internal/catalog/`** — Reads the embedded skill catalog. Skills are stored as versioned markdown files under `internal/catalog/skills/<id>/<semver>.md` with YAML front matter (`id`, `version`, `description`). Uses `//go:embed skills/*/*.md` to bundle them into the binary.
- **`internal/project/`** — Manages skills installed in a user's project (the `.skills/` directory). Handles install, remove, list, and reading front matter from installed skill files.
- **`internal/version/`** — GitHub release version checking. Compares current build version against the latest GitHub release tag.

### Key patterns

- **Embedded catalog**: All catalog skills are compiled into the binary via `embed.FS`. Adding a new skill means adding a `<semver>.md` file under `internal/catalog/skills/<id>/` — the catalog tests automatically validate it.
- **Front matter parsing**: Both catalog and project packages parse a simple `---` delimited YAML front matter (id, version, description) from markdown files. This is a custom parser, not a YAML library.
- **Version is set at build time**: `cmd.Version` is injected via `-ldflags "-X .../cmd.Version=..."`. Default is `"dev"`.

## Commit Style

Conventional Commits: `feat:`, `fix:`, `chore:`, `docs:`, `ci:`, `refactor:`, `test:`
