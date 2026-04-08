# Product Owner

**Score: 74/100**

## What is being assessed
Feature completeness versus stated goals, user journey quality, error UX, missing user-facing features, and accessibility basics for a CLI tool.

## Methodology
Read README.md to identify stated goals and listed features. Exercised all command help text. Assessed error messages from source. Evaluated user journey for install, add, update, remove, and search.

## Findings

### Passing checks
- Core user journey (install → `skills ls` → `skills add go` → `skills update`) is complete and coherent
- `skills ls` colour-codes status (green/yellow/red) — clear visual affordance — cmd/ls.go
- `skills add go@1.0.0` version-pinning is supported — cmd/add.go
- `skills search <term>` searches ID, description, and content — cmd/search.go
- `skills get <id>` prints full skill content to stdout for piping — cmd/get.go
- `skills versions <id>` lists available versions — cmd/versions.go
- `skills init` wizard auto-detects installed agents and pre-selects appropriate skills — cmd/init.go
- `skills update` dry-run-by-default pattern prevents accidental mass updates — cmd/update.go
- Version update notice on every command — users always know if they are outdated
- Clear error messages: `"skill %q is not installed"`, `"skill %q not found in catalog"` — cmd/rm.go, cmd/add.go

### Issues found
| Finding | Severity | Location | Recommendation |
|---------|----------|----------|----------------|
| `skills update` and `skills init` not documented in README usage table | Medium | README.md | Add both commands to the usage section |
| `skills ls` shows catalog skills but not locally installed skills that are NOT in the catalog | Low | cmd/ls.go | Consider a `skills ls --installed` flag to show only what is actually installed |
| No `skills info <id>` command to show description + all versions without piping through `skills get` | Low | — | Add `skills info <id>` as a user-friendly alternative to `skills get` |
| No confirmation prompt for destructive operations (`skills rm`) | Low | cmd/rm.go | Add a `--yes`/`-y` flag to rm (matching update) or a confirmation prompt |
| `skills search` output does not show installation status | Low | cmd/search.go | Add an INSTALLED column to search output, consistent with `ls` |
| Error message when skill not in catalog says "not found in catalog" — no suggestion to run `skills ls` | Low | internal/catalog/catalog.go:50 | Append: "Run 'skills ls' to see available skills." |
| No `skills rm --all` to uninstall everything | Low | — | Nice-to-have for project cleanup |

## Verdict
The product is feature-complete for its stated goals. The user journey from discovery to installation is smooth. Main gaps are documentation of newer commands (update, init) and small UX polish items like showing install status in search results and adding a suggestion on catalog-not-found errors.

## Changes since last assessment
First assessment.

## Remaining recommendations
| Finding | Severity | Recommendation |
|---------|----------|----------------|
| Document update + init in README | Medium | Update usage table and add examples |
| Show install status in search output | Low | Add INSTALLED column to runSearch |
| Add -y flag to rm | Low | Consistent with update command |
| Helpful error suggestions | Low | Append "Run skills ls" hint on not-found errors |
