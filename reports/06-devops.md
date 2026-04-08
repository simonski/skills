# DevOps

**Score: 64/100**

## What is being assessed
Build pipeline quality (Makefile), CI/CD configuration, release process, secrets management, Docker usage, version management, and operational readiness.

## Methodology
Read Makefile, .github/workflows/publish.yml, go.mod, VERSION file, .gitignore, and packaging/homebrew/skills.rb.template.

## Findings

### Passing checks
- Makefile has clearly named, documented targets: build, test, lint, release, publish, clean, install — Makefile:help
- Cross-compilation covers darwin/amd64, darwin/arm64, linux/amd64, linux/arm64, windows/amd64 — Makefile:release
- Release archives use correct formats (.tar.gz for unix, .zip for windows) — Makefile:release
- GitHub Actions workflow correctly separates test and publish jobs with `needs: test` — .github/workflows/publish.yml
- Loop prevention: publish job skips commits starting with "chore: release" — .github/workflows/publish.yml:13
- Go build uses `-trimpath` flag (removes local paths from binary) — Makefile:GOFLAGS
- ldflags injects version at build time — Makefile:build
- Homebrew formula template uses SHA256 checksums for all platform archives — packaging/homebrew/skills.rb.template
- VERSION file is single source of truth for version — Makefile, cmd/root.go
- `.gitignore` covers dist/ and common Go artifacts

### Issues found
| Finding | Severity | Location | Recommendation |
|---------|----------|----------|----------------|
| No lint step in CI — `go vet` and `staticcheck` only run via `make lint` locally | High | .github/workflows/publish.yml | Add `run: go vet ./...` to test job; add staticcheck |
| No vulnerability scanning in CI (`govulncheck`) | High | .github/workflows/publish.yml | Add `go install golang.org/x/vuln/cmd/govulncheck@latest && govulncheck ./...` |
| `make build` auto-bumps patch version — running it during development creates unintended version commits | Medium | Makefile:build | Split into `make bump` and `make build`; `make build` should compile only |
| CI does not run `go mod verify` | Medium | .github/workflows/publish.yml | Add `go mod verify` to detect tampered module cache |
| No coverage threshold enforced in CI | Medium | .github/workflows/publish.yml | Add `go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out` with a minimum threshold check |
| Publish job uses `ubuntu-latest` — Go version comes from go.mod but ubuntu-latest runner Go may lag | Low | .github/workflows/publish.yml | Pin `go-version-file: go.mod` (already done) — acceptable |
| No Docker image or container build | Low | — | Not required for a CLI tool distributed via Homebrew |
| No RELEASE.md or changelog automation | Low | — | Add a CHANGELOG.md or use `gh release` with auto-generated notes |
| TAP_TOKEN required but not documented in repo | Medium | .github/workflows/publish.yml | Add docs/SECRETS.md or README section documenting required secrets |

## Verdict
The CI/CD foundation is solid: test-before-publish gating, cross-platform release builds, and automated Homebrew tap updates are all present. The critical gaps are the absence of linting and vulnerability scanning in CI, which mean quality regressions and known CVEs could ship undetected. These are one-line additions.

## Changes since last assessment
First assessment.

## Remaining recommendations
| Finding | Severity | Recommendation |
|---------|----------|----------------|
| Add lint + govulncheck to CI | High | go vet, staticcheck, govulncheck in test job |
| Add go mod verify to CI | Medium | Step in test job |
| Separate make bump from make build | Medium | Refactor Makefile |
| Enforce coverage threshold in CI | Medium | Fail if total coverage drops below 60% |
| Document required secrets | Medium | Add SECRETS.md or README section |
