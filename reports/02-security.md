# Security

**Score: 48/100** (was 72)

## What is being assessed
Security here covers filesystem safety, secrets handling, outbound network use, dependency hygiene, and release pipeline hardening for a local CLI. Good looks like strict input validation, bounded network behavior, safe secret handling, and no dangerous file operations reachable from user input.

## Methodology
Read the command layer and project filesystem package, then checked CI and release automation for verification steps and secret handling. Files read included `cmd/rm.go`, `cmd/add.go`, `cmd/update.go`, `internal/project/project.go`, `internal/version/version.go`, `.github/workflows/publish.yml`, and `.gitignore`.

## Findings

### Passing checks
- Project installs use explicit permissions: directories at `0o755` and files at `0o644` (`internal/project/project.go:83-89`).
- Version checks use a 5-second HTTP timeout and optional `GITHUB_TOKEN` auth instead of an unbounded client (`internal/version/version.go:19-27`).
- CI now verifies modules, runs `go vet`, and executes `govulncheck` before publish (`.github/workflows/publish.yml:17-28`).
- Release artifacts and Windows binaries are ignored locally, reducing accidental commits of build outputs (`.gitignore:1-4`).

### Issues found
| Finding | Severity | Location | Recommendation |
|---|---|---|---|
| Raw `rm` input can escape the `.skills/` directory because `id` is joined directly into a filesystem path | Critical | `cmd/rm.go:25-37`, `internal/project/project.go:95-103` | Add a shared skill-ID validator such as `^[a-z0-9][a-z0-9-]*$` and reject `.` / `..` segments before any project operation. |
| `project.Get` and `project.SkillPath` also trust caller-supplied IDs, so read paths can traverse outside the project root | High | `internal/project/project.go:25-28`, `internal/project/project.go:65-78` | Enforce the same validation inside `internal/project` so every caller gets the invariant automatically. |
| Release publishing clones the tap repository with the token embedded in the URL | Medium | `.github/workflows/publish.yml:88-99` | Switch to a credential helper or `GIT_ASKPASS` pattern so the token is not passed as a process argument. |

## Verdict
The repo has improved supply-chain checks, but the core local filesystem boundary is still weak. A single unchecked skill ID in `rm` is enough to make this the highest-priority security issue in the codebase.

## Changes since last assessment
- CI hardening improved with `go mod verify`, `go vet`, and `govulncheck` (`.github/workflows/publish.yml:17-28`).
- The previous report's assumption that path traversal was prevented was incorrect; the current code still trusts raw IDs in project path joins (`internal/project/project.go:25-28,95-103`).
- SBOM generation was added to releases (`.github/workflows/publish.yml:53-56`).

## Remaining recommendations
| Finding | Severity | Recommendation |
|---|---|---|
| Path traversal via skill ID | Critical | Validate IDs centrally in `internal/project` and cover invalid IDs with tests. |
| Secret-in-URL release step | Medium | Replace tokenized clone URLs with non-argument credential handling. |
| Silent network failure handling | Low | Log or surface update-check failures behind a debug flag instead of swallowing all errors unconditionally. |
