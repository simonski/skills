# Idiomatic Go

**Score: 82/100** (was 81)

## What is being assessed
Idiomatic Go covers package structure, error handling, naming, test style, and whether the code follows common Go conventions without unnecessary abstraction. Good looks like small focused packages, wrapped errors, predictable exported APIs, and simple control flow.

## Methodology
Read all Go packages under `cmd/` and `internal/`, with extra attention to error wrapping, reusable helpers, flag handling, and duplicated logic. Cross-checked the build workflow through `Makefile` and `.github/workflows/publish.yml`.

## Findings

### Passing checks
- Command handlers consistently use `RunE`, which keeps error handling explicit and exit behavior centralized (`cmd/add.go:15-34`, `cmd/root.go:35-39`).
- Errors are wrapped with context in the internal packages (`internal/catalog/catalog.go:61-70`, `internal/project/project.go:37-38,71-72,84-89`, `internal/version/version.go:21-23,28-35,41-42`).
- Exported types and functions have doc comments in the internal packages (`internal/catalog/catalog.go:19-27,45-58,81-109,153-154`, `internal/project/project.go:14-30,64-81,94-106`).
- Tests use `t.TempDir()` rather than shared directories (`cmd/update_test.go:20-27`, `internal/project/project_test.go:11-17`).

### Issues found
| Finding | Severity | Location | Recommendation |
|---|---|---|---|
| Semver parsing and comparison are duplicated across two packages | Medium | `internal/catalog/catalog.go:129-150`, `internal/version/version.go:71-105` | Extract one shared internal semver helper and test it directly. |
| The update check runs synchronously in `PersistentPreRun`, which makes command startup depend on network I/O | Medium | `cmd/root.go:25-31`, `internal/version/version.go:19-45` | Run the check asynchronously or cache it so command execution stays fast and predictable. |
| `updateConfirm` is mutable package-level state rather than command-local flag state | Low | `cmd/update.go:14-48` | Read the flag value from the Cobra command instead of storing it globally. |
| The codebase mixes `os.IsNotExist` with `errors.Is(err, fs.ErrNotExist)` | Low | `internal/project/project.go:34,48,68,97`, `internal/catalog/catalog.go:62,86` | Standardize on `errors.Is` for consistency. |

## Verdict
The code is still clean and recognizably Go-like: small packages, direct data flow, and good error context. The biggest idiomatic issue is duplicated semver logic, not a broad style problem.

## Changes since last assessment
- No major code-shape change in the Go packages.
- The repo improved around CI and docs rather than core Go structure.

## Remaining recommendations
| Finding | Severity | Recommendation |
|---|---|---|
| Duplicated semver helpers | Medium | Consolidate version parsing/comparison into one package. |
| Blocking update check | Medium | Remove network I/O from the hot startup path. |
| Global flag state | Low | Resolve Cobra flags inside the command handler. |
