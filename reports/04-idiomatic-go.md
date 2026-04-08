# Idiomatic Go

**Score: 81/100**

## What is being assessed
Error handling patterns, context propagation, concurrency safety, package organisation, interface design, naming conventions, build tooling (Makefile, CI), linting configuration, and test patterns — assessed against effective Go and community standards.

## Methodology
Read all .go files in cmd/ and internal/. Checked error wrapping, use of context, goroutine usage, package naming, interface definitions, variable naming, Makefile targets, and CI workflow.

## Findings

### Passing checks
- All errors are wrapped with `%w` and include call-site context (e.g. `fmt.Errorf("reading skill %s@%s: %w", ...)`) — internal/catalog/catalog.go:62
- `errors.Is(err, fs.ErrNotExist)` used correctly for file-not-found checks — internal/catalog/catalog.go:63
- Package names are lowercase single words: `catalog`, `project`, `version`, `cmd` — correct
- `cobra.Command.RunE` used throughout (returns errors rather than calling os.Exit directly) — all cmd/*.go
- `os.IsNotExist(err)` used in project package alongside the newer `errors.Is` — minor inconsistency but not wrong
- Exported functions all have doc comments starting with the function name — internal/catalog/catalog.go, internal/project/project.go
- No use of `panic` anywhere in application code
- `t.TempDir()` used in all project tests — correct cleanup pattern — internal/project/project_test.go
- `embed.FS` used correctly with `//go:embed` directive — internal/catalog/catalog.go:14
- `bufio.Scanner` used for line-by-line parsing rather than `strings.Split` — internal/catalog/catalog.go:158

### Issues found
| Finding | Severity | Location | Recommendation |
|---------|----------|----------|----------------|
| Global mutable `updateConfirm bool` var used for flag state | Low | cmd/update.go:15 | Acceptable cobra pattern but note it is not safe for parallel command execution in tests; use `cmd.Flags().GetBool("yes")` inside RunE instead |
| `checkForUpdates()` makes a live HTTP call on every command invocation (except `version`) | Medium | cmd/root.go:26-34 | Run in a goroutine with a context so it never blocks the main command path; or cache result in a temp file with a TTL |
| `os.IsNotExist` mixed with `errors.Is(err, fs.ErrNotExist)` | Low | internal/project/project.go:35 vs internal/catalog/catalog.go:63 | Standardise on `errors.Is(err, fs.ErrNotExist)` throughout |
| No `.golangci.yml` linting configuration — `make lint` only runs `go vet` + `staticcheck` if installed | Low | Makefile:lint | Add `.golangci.yml` with at minimum `errcheck`, `govet`, `staticcheck`, `gosimple` |
| CI has no lint step | Medium | .github/workflows/publish.yml | Add `go vet ./...` and `staticcheck ./...` as a step in the test job |
| `compareVersions` and `parseVersion` in catalog.go duplicate logic from `semverGT`/`splitSemver` in version.go | Medium | internal/catalog/catalog.go:131, internal/version/version.go:62 | Extract a shared `internal/semver` package used by both |
| `make build` bumps the patch version on every invocation — breaks reproducibility in local dev | Medium | Makefile:build | Separate `make bump` (version increment) from `make build` (compile only); CI should call `make bump && make release` |

## Verdict
The code is clean, well-structured, and follows Go idioms closely. The main actionable issues are: the duplicate semver logic between two packages, the blocking HTTP update check on the critical path, and the absence of a lint step in CI. None are severe, and the overall code quality is high for a project of this size.

## Changes since last assessment
First assessment.

## Remaining recommendations
| Finding | Severity | Recommendation |
|---------|----------|----------------|
| Async update check | Medium | Wrap checkForUpdates in a goroutine; print result at end of command |
| Deduplicate semver logic | Medium | Create internal/semver package |
| Separate version bump from build | Medium | Add `make bump` target; make `make build` idempotent |
| Add .golangci.yml | Low | Configure errcheck, govet, staticcheck, gosimple |
| Add lint step to CI | Medium | go vet + staticcheck in test job |
