# New Starter

**Score: 58/100**

## What is being assessed
Onboarding effectiveness from the perspective of an engineer joining on day one. Reading order, way of working, development setup speed, ticket workflow, testing expectations, collaboration patterns, and common pitfalls.

## Methodology
Walked through the repository as a new joiner. Read README.md, CLAUDE.md, AGENTS.md. Attempted to reconstruct the path from `git clone` to running tests. Identified undocumented steps and gotchas.

## Findings

### Passing checks
- **README.md** provides installation, all commands, and build instructions — a new user can install and use the tool quickly
- **CLAUDE.md** documents build commands (`make build`, `make test`), coding conventions, and agent workflow — useful for AI-assisted development
- **AGENTS.md** documents the ticket workflow (`tk` commands) — internal process is captured
- **`make test`** runs all tests without any external dependencies — zero-setup test run
- **`make build`** produces a working binary in `dist/` — straightforward build
- **Makefile `help` target** is the entry point for build commands — discoverable

### Issues found
| Finding | Severity | Location | Recommendation |
|---------|----------|----------|----------------|
| No ONBOARDING.md or CONTRIBUTING.md — no documented reading order for new contributors | High | / | Create docs/ONBOARDING.md (template below) |
| `make build` auto-bumps VERSION on every invocation — a new dev running `make build` locally will create unintended version changes | High | Makefile:build | Document this gotcha prominently; separate `make bump` from `make build` |
| CI requires TAP_TOKEN secret with no documentation on how to obtain or configure it | High | .github/workflows/publish.yml | Document in README or docs/SECRETS.md |
| No branching strategy documented — unclear if PRs go to main, if feature branches are used | Medium | / | Add to CONTRIBUTING.md: branch naming convention, PR process, merge strategy |
| Ticket workflow requires `tk` CLI which is not in PATH by default (at /usr/local/bin/tk uses older version with DB schema incompatibility) | High | AGENTS.md | Document correct `tk` binary location: `/Users/simon/go/bin/ticket`; add setup step |
| No `make setup` target to install development dependencies (staticcheck, govulncheck, etc.) | Medium | Makefile | Add `make setup` that installs required tools |
| `go test ./...` is sufficient for running tests — but test conventions (when to write tests, which packages need tests) are undocumented | Medium | / | Add testing expectations to CONTRIBUTING.md |
| No PR template or issue template | Low | .github/ | Add .github/PULL_REQUEST_TEMPLATE.md and .github/ISSUE_TEMPLATE/ |

## Recommended ONBOARDING.md

Create `/Users/simon/code/skills/docs/ONBOARDING.md` with the following structure:

```markdown
# Onboarding Guide

## Day 1: Reading order
1. README.md — understand what the tool does and how to use it
2. CLAUDE.md — coding conventions and agent workflow
3. AGENTS.md — ticket workflow with `tk`
4. This file — development setup and way of working

## Development setup
Prerequisites: Go 1.24+

git clone https://github.com/simonski/skills.git
cd skills
make test          # verify everything works
make build         # builds dist/skills

⚠️ NOTE: `make build` auto-increments the patch version in VERSION.
Do NOT commit VERSION changes from local builds — only the CI pipeline should bump versions.

## Running tests
make test          # or: go test ./...

## Way of working
- Tickets are tracked in `.ticket/` using the `tk` CLI tool (binary: /Users/simon/go/bin/ticket)
- Pick up work: `tk ls`, then `tk claim -id <id>` and `tk state -id <id> active`
- Branch from main; use descriptive branch names
- Commit messages: `type: description` (feat, fix, chore, docs, ci, refactor)

## Adding a new catalog skill
1. Create `internal/catalog/skills/<id>/1.0.0.md` with YAML front matter
2. Create `internal/catalog/skills/<id>/README.md` and `CHANGELOG.md`
3. Add the skill to the catalog table in README.md
4. Run `go test ./...` to verify the skill parses correctly

## Required secrets (for publishing)
- TAP_TOKEN: GitHub PAT with `Contents: write` on simonski/homebrew-tap
  Set in: GitHub repo Settings → Secrets → Actions
```

## Verdict
A new starter can clone, test, and use the tool within 10 minutes thanks to the README and `make test`. However, contributing is poorly documented: there is no CONTRIBUTING.md, branching strategy, or testing guide. The `make build` version-bump gotcha will trip up every new developer. Creating docs/ONBOARDING.md and CONTRIBUTING.md would resolve the majority of friction.

## Changes since last assessment
First assessment.

## Remaining recommendations
| Finding | Severity | Recommendation |
|---------|----------|----------------|
| Create docs/ONBOARDING.md | High | Use template above |
| Document make build version-bump gotcha | High | In ONBOARDING.md and Makefile comment |
| Document TAP_TOKEN setup | High | In ONBOARDING.md / docs/SECRETS.md |
| Create CONTRIBUTING.md | Medium | Branch naming, commit style, PR process |
| Add make setup target | Medium | Install staticcheck, govulncheck |
| Add PR template | Low | .github/PULL_REQUEST_TEMPLATE.md |
