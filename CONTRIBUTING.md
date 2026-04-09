# Contributing to skills

Thank you for contributing! This document covers branching, commit style, PR process, and how to add new skills to the catalog.

## Development setup

```bash
git clone https://github.com/simonski/skills.git
cd skills
make test       # verify everything works — no external dependencies required
```

> ⚠️ **Gotcha:** `make build` auto-increments the patch version in `VERSION`. Do **not** commit `VERSION` changes from local builds. Only the CI pipeline should bump the version via `make release`.
> Use `go build -o dist/skills .` for a plain local build that doesn't touch `VERSION`.

## Branching

- Branch from `main`
- Use descriptive branch names: `feat/add-rust-skill`, `fix/semver-comparison`, `docs/onboarding`
- Keep branches short-lived; merge via PR

## Commit messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>: <short description>

[optional body]
```

Types: `feat`, `fix`, `chore`, `docs`, `ci`, `refactor`, `test`

Examples:
```
feat: add rust skill to catalog
fix: semver comparison for double-digit minor versions
docs: add ONBOARDING guide
ci: add govulncheck to test job
```

## Pull requests

1. Open a PR against `main`
2. Fill in the PR template
3. Ensure `make test` and `make lint` pass locally
4. Request review — one approval required to merge

## Running tests and lint

```bash
make test     # go test ./...
make lint     # go vet + staticcheck
make setup    # install dev tools (staticcheck, govulncheck)
```

## Adding a new skill to the catalog

1. Create the directory: `internal/catalog/skills/<id>/`
2. Create the versioned skill file: `internal/catalog/skills/<id>/1.0.0.md`

   The file must have YAML front matter followed by the skill content:

   ```markdown
   ---
   id: my-skill
   version: 1.0.0
   description: One-line description shown in skills ls
   ---

   # My Skill

   Skill content here...
   ```

3. Create `internal/catalog/skills/<id>/README.md` — overview of the skill
4. Create `internal/catalog/skills/<id>/CHANGELOG.md` — version history
5. Add the skill to the catalog table in `README.md`
6. Run `go test ./...` — the catalog tests will verify the new skill parses correctly

## Required secrets (for maintainers publishing releases)

See [docs/SECRETS.md](docs/SECRETS.md).
