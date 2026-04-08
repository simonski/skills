# Tech Lead

**Score: 79/100**

## What is being assessed
File sizes, code duplication, error message consistency, magic numbers, cyclomatic complexity, dead code, naming conventions, interface sizes, helper reuse, and refactoring opportunities.

## Methodology
Read all source files. Counted lines per file. Searched for duplicated logic patterns, magic numbers, and inconsistent naming. Assessed function length and complexity.

## Findings

### Passing checks
- No file exceeds 308 lines (cmd/init.go) — well within the 700-line limit
- `buildSkillFile()` is defined once in cmd/add.go:78 and reused by cmd/update.go and cmd/init.go — good helper reuse
- Error messages are consistently lowercase and wrapped with context (e.g. `"reading skill %s@%s: %w"`)
- No magic numbers in business logic — `0o644` and `0o755` file permissions are self-documenting constants
- No dead code identified — all exported functions are used
- `parseSkillArg` is appropriately extracted for reuse between add/update flows
- Package `cmd` is flat — no sub-packages needed at this scale
- `skillStatusStr` helper extracted and reused across `ls` and `init` commands — cmd/init.go:229, cmd/ls.go

### Issues found
| Finding | Severity | Location | Recommendation |
|---------|----------|----------|----------------|
| Duplicate semver logic: `compareVersions`+`parseVersion` in catalog.go vs `semverGT`+`splitSemver` in version.go | High | internal/catalog/catalog.go:131, internal/version/version.go:62 | Extract `internal/semver` package with a single canonical implementation |
| `runInit` is 120+ lines of interleaved UI, business logic, and I/O | Medium | cmd/init.go:93 | Extract interactive loop, agent detection, and installation into separate functions |
| `os.IsNotExist` vs `errors.Is(err, fs.ErrNotExist)` used inconsistently | Low | internal/project/project.go:35, internal/catalog/catalog.go:63 | Standardise on `errors.Is` throughout |
| `updateConfirm` global var for flag state | Low | cmd/update.go:15 | Use `cmd.Flags().GetBool("yes")` in RunE to avoid global state |
| `init.go` mixes `skillEntry` struct definition with large interactive function — struct belongs in a types file or catalog package | Low | cmd/init.go:16-22 | Move `skillEntry` to cmd/ls.go or a shared cmd types file |

## Verdict
Code quality is good. The codebase is small, well-organised, and free of obvious duplication apart from the semver logic. The main refactoring opportunity is extracting a shared semver package and splitting `runInit` into smaller functions. Neither is urgent but both would improve testability and maintainability.

## Changes since last assessment
First assessment.

## Remaining recommendations
| Finding | Severity | Recommendation |
|---------|----------|----------------|
| Deduplicate semver logic | High | Create internal/semver package |
| Split runInit | Medium | Extract sub-functions for agent detection, selection loop, apply step |
| Standardise error sentinel check | Low | Use errors.Is(err, fs.ErrNotExist) everywhere |
