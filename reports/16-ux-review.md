# UX Review

**Score: 74/100** (was N/A)

## What is being assessed
UX for a CLI means readable output, safe defaults, discoverable workflows, and accessibility-friendly behavior in terminal environments. Good looks like consistent status signals, accurate help, and command flows that avoid accidental destructive actions.

## Methodology
Reviewed command help text, table output, status coloring, and interactive flows in `ls`, `search`, `update`, `init`, and `rm`. Compared command output patterns for consistency.

## Findings

### Passing checks
- `ls` uses clear status labels and color distinctions for installed, update-available, and missing skills (`cmd/ls.go:17-23,40-63`).
- `update` defaults to dry-run mode, which is a strong UX choice for bulk changes (`cmd/update.go:16-33,95-97`).
- `init` provides an interactive selection loop with explicit affordances for all/none/install/quit (`cmd/init.go:163-210`).
- `version` prints update guidance without polluting command stdout for non-version commands (`cmd/root.go:62-66`, `cmd/version.go:21-39`).

### Issues found
| Finding | Severity | Location | Recommendation |
|---|---|---|---|
| The CLI emits colors unconditionally and does not respect `NO_COLOR` or a `--no-color` option | Medium | `cmd/ls.go:21-23,40-42`, `cmd/update.go:61-63,129-131`, `cmd/init.go:99-102,230-248` | Add color suppression support for accessibility and scripting environments. |
| `search` lacks the install-state context that `ls` already provides | Low | `cmd/search.go:49-53`, `cmd/ls.go:44-63` | Add an installed/status indicator to search results. |
| `rm` has no confirmation step despite being destructive | Low | `cmd/rm.go:13-37` | Add `--yes` or an interactive confirmation mode. |
| `ls` uses mismatched format widths between header and row rendering | Low | `cmd/ls.go:44-45,63` | Reuse one format string for the table to keep alignment stable. |

## Verdict
The CLI is already usable and fairly consistent, especially around `ls`, `update`, and `init`. The biggest UX gaps are accessibility and deletion safety, not fundamental workflow design.

## Changes since last assessment
- This category is newly scored instead of being left N/A.
- No major UX regressions were found, but the repo still lacks color opt-out and safer delete ergonomics.

## Remaining recommendations
| Finding | Severity | Recommendation |
|---|---|---|
| Missing color opt-out | Medium | Respect `NO_COLOR` and/or add `--no-color`. |
| Search without status context | Low | Show installed/update state in search output. |
| Destructive remove flow | Low | Add a confirmation or explicit apply flag. |
