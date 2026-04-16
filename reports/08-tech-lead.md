# Tech Lead

**Score: 78/100** (was 79)

## What is being assessed
This category looks at maintainability from a lead engineer's perspective: duplication, file sizes, refactoring pressure, consistency, and how easy the code is to evolve safely. Good looks like low duplication, obvious invariants, and small testable units.

## Methodology
Reviewed the largest files, helper reuse, and repeated logic across packages. Cross-checked line counts and command/package organization.

## Findings

### Passing checks
- No Go source file exceeds the 700-line threshold; the largest is `cmd/init.go` at 308 lines.
- Shared skill-file rendering is centralized in `buildSkillFile` and reused across commands (`cmd/add.go:77-80`, `cmd/update.go:102-105,146-148`, `cmd/init.go:288-290`).
- Package boundaries are still small and understandable (`main.go:1-7`, `cmd/root.go:42-53`, `docs/ARCHITECTURE.md:23-48`).

### Issues found
| Finding | Severity | Location | Recommendation |
|---|---|---|---|
| Semver logic is duplicated across catalog and version packages | High | `internal/catalog/catalog.go:129-150`, `internal/version/version.go:71-105` | Extract one shared internal package and delete the duplicates. |
| `runInit` mixes detection, rendering, selection state, stdin parsing, and install/update behavior in one path | Medium | `cmd/init.go:93-212` | Split the wizard into smaller pure helpers plus a thin interactive shell. |
| Project path invariants are enforced nowhere central, which leaves every caller responsible for safety | Medium | `internal/project/project.go:25-28,81-103` | Make `internal/project` own skill-ID validation instead of trusting its callers. |
| The `ls` table uses mismatched status column widths between header and rows | Low | `cmd/ls.go:44-45,63` | Normalize the format string once to avoid slow layout drift. |

## Verdict
The codebase is still compact and easy to reason about, but a few seams are starting to show. The main maintainability win would be centralizing invariants and removing duplicated version logic before more commands accrete around them.

## Changes since last assessment
- File count and feature count grew, but overall file size discipline remains good.
- The major lead-level concerns are still duplication and safety invariants, not broad architectural sprawl.

## Remaining recommendations
| Finding | Severity | Recommendation |
|---|---|---|
| Duplicated semver code | High | Create one semver helper package and test it thoroughly. |
| Monolithic init flow | Medium | Extract pure selection and detection helpers from the interactive loop. |
| Unowned ID invariant | Medium | Move skill-ID validation into the project package. |
