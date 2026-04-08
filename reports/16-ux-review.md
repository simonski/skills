# UX Review

**Score: N/A (CLI)**

## What is being assessed
For a CLI tool: output consistency, colour usage, error display, help text quality, command discoverability, progressive disclosure, and keyboard/terminal ergonomics.

## Methodology
Read all cmd/*.go files. Assessed help text, colour usage, error message patterns, output formatting, and command naming conventions.

## Findings

### Passing checks
- **Consistent colour scheme**: green = success, yellow = warning/up-to-date, red = not installed, cyan = available update — cmd/ls.go, cmd/init.go, cmd/update.go
- **All commands have Long help text** with examples — cobra Long field populated in all command files
- **`skills ls` output is tabular** with aligned columns — cmd/ls.go
- **`skills init` is interactive** with a clear selection loop and explicit instructions — cmd/init.go
- **Dry-run pattern in `skills update`** — prevents accidental changes; shows what would happen first — cmd/update.go
- **Version check notice** appears on stderr (via `fmt.Fprintf(os.Stderr, ...)`) — does not pollute stdout pipeline — cmd/root.go:60
- **`skills get <id>`** pipes cleanly to other tools (no colour, just raw content) — cmd/get.go
- **Error messages are quoted**: `skill "go" is not installed` — consistent use of %q format

### Issues found
| Finding | Severity | Location | Recommendation |
|---------|----------|----------|----------------|
| `skills update` and `skills init` missing from README usage table | Medium | README.md | Add both commands to the documented usage section |
| `skills search` output has no INSTALLED column — inconsistent with `skills ls` | Low | cmd/search.go | Add installed status indicator to search output |
| `skills rm` has no `--yes` confirmation flag — unlike `skills update` which requires `-y` | Low | cmd/rm.go | Add `-y` flag for consistency, or add an interactive "are you sure?" prompt |
| `skills ls` column widths are hardcoded with `%-20s` — long skill IDs or descriptions may misalign | Low | cmd/ls.go:44 | Calculate column widths dynamically from data |
| `skills init` prints `[1]`, `[2]` indices but entering `0` gives an unhelpful "Unknown input" message | Low | cmd/init.go:174 | Specifically handle `0` and out-of-range inputs with a clearer message |
| No `--no-color` flag or `NO_COLOR` env var support | Low | cmd/root.go | Check `os.Getenv("NO_COLOR")` or add `--no-color` flag; pass to fatih/color's `NoColor` var |
| Version update notice is on stderr — correct — but mixed with command output in some terminals | Low | cmd/root.go | Acceptable; no change needed |

## Verdict
CLI UX is polished and consistent. Colour usage is purposeful, help text is complete, and the dry-run pattern in `update` is excellent. Minor gaps: `search` lacks install status, `rm` lacks a safety confirmation, and `NO_COLOR` is not respected. None are blockers.

## Changes since last assessment
First assessment.

## Remaining recommendations
| Finding | Severity | Recommendation |
|---------|----------|----------------|
| Respect NO_COLOR env var | Low | Set color.NoColor = true when NO_COLOR is set |
| Add install status to search output | Low | Reuse skillStatusStr in runSearch |
| Add -y to rm | Low | Match update command pattern |
| Dynamic column widths in ls | Low | Pre-scan data to calculate max widths |
