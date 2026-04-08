# Security

**Score: 72/100**

## What is being assessed
Authentication, access control, data protection, input validation, path traversal prevention, rate limiting, and dependency vulnerability management for a CLI tool.

## Methodology
Read all source files in cmd/ and internal/. Checked file-write paths, input parsing, HTTP client usage, dependency list in go.mod/go.sum, and CI pipeline for vulnerability scanning.

## Findings

### Passing checks
- Path traversal prevented: `project.SkillPath` uses `filepath.Join(root, ".skills", id+".md")` — `id` comes from catalog keys (alphanumeric with hyphens) so no traversal possible. (internal/project/project.go:26)
- Skill IDs validated: `parseSkillArg` rejects empty IDs and empty version strings. (cmd/add.go:84)
- No credential storage: the binary stores no secrets, tokens, or passwords.
- File permissions: skill files written with `0o644`, skills directory created with `0o755`. (internal/project/project.go:80,83)
- Dependency count is minimal (cobra, fatih/color, plus transitive). Attack surface is small.
- HTTP client uses 5-second timeout for GitHub release API calls — no unbounded blocking. (internal/version/version.go:18)
- No `exec.Command` usage — no shell injection surface.

### Issues found
| Finding | Severity | Location | Recommendation |
|---------|----------|----------|----------------|
| No `govulncheck` in CI | Medium | .github/workflows/publish.yml | Add `go install golang.org/x/vuln/cmd/govulncheck@latest && govulncheck ./...` as a CI step before publish |
| Unauthenticated GitHub API call — susceptible to rate limiting (60 req/hr per IP) | Low | internal/version/version.go:17 | Accept a `GITHUB_TOKEN` env var and pass it as a Bearer header when present |
| `go.sum` not checked for tampering in CI | Low | .github/workflows/publish.yml | Add `-goflags=-mod=readonly` or `go mod verify` step in CI to detect tampered module cache |
| No `.gitignore` entry for `dist/` binary artifacts on Windows (`.exe`) | Low | .gitignore | Confirm `dist/` glob covers all artifacts |

## Verdict
The attack surface is tiny for a local CLI tool. There are no auth systems, no network servers, and no user-supplied data that reaches dangerous sinks. Main gaps are the absence of automated vulnerability scanning in CI and unauthenticated GitHub API calls that can be rate-limited. Both are straightforward to fix.

## Changes since last assessment
First assessment.

## Remaining recommendations
| Finding | Severity | Recommendation |
|---------|----------|----------------|
| Add govulncheck to CI | Medium | `govulncheck ./...` step before publish job |
| Add go mod verify to CI | Low | `go mod verify` step in test job |
| Pass GITHUB_TOKEN to release API check | Low | Read `os.Getenv("GITHUB_TOKEN")`, set `Authorization: Bearer` header |
