# SRE

**Score: 52/100**

## What is being assessed
Observability (metrics, structured logging, tracing), alerting readiness, runbook existence, incident response documentation, backup/restore procedures, capacity planning, SLA/SLO definition, and graceful degradation patterns.

## Methodology
Searched for Prometheus instrumentation, slog/zerolog usage, distributed tracing, health check endpoints, and operational documentation.

## Findings

### Passing checks
- **Graceful degradation on network failure**: `checkForUpdates()` silently ignores errors — users never see a failure if GitHub API is unavailable — cmd/root.go:48
- **Timeout on external calls**: 5-second timeout on GitHub API — prevents indefinite hangs — internal/version/version.go:18
- **Release process is automated**: GitHub Actions workflow handles tagging, release creation, and Homebrew tap update — .github/workflows/publish.yml
- **No stateful service to monitor**: the binary is a local CLI tool — no server processes, no daemons, no databases to observe
- **Homebrew formula includes test stanza**: `assert_match version.to_s, shell_output("#{bin}/skills version")` — verifies binary executes correctly post-install — packaging/homebrew/skills.rb.template

### Issues found
| Finding | Severity | Location | Recommendation |
|---------|----------|----------|----------------|
| No CHANGELOG.md — incidents and regressions cannot be traced to a specific release | Medium | / | Create CHANGELOG.md per Keep a Changelog format |
| No rollback procedure documented | Medium | / | Document: to roll back, `brew install simonski/tap/skills@<version>` or download previous GitHub release |
| CI publish failure leaves partial state (tag pushed but release not created, or release created but tap not updated) | High | .github/workflows/publish.yml | Add atomic rollback: if tap update fails, delete the GitHub release and tag; use job-level `if: failure()` steps |
| No smoke test after Homebrew tap update (verify formula installs and runs) | Medium | .github/workflows/publish.yml | Add a post-publish verification step: `brew tap simonski/tap && brew install skills && skills version` on macOS runner |
| No monitoring of GitHub release download counts or install metrics | Low | — | Optional: GitHub Insights provides release download counts natively |
| Update check cache miss on every invocation adds network dependency to every command | Medium | cmd/root.go | Cache latest-version with 24h TTL in ~/.cache/skills/ to eliminate repeated network calls |

## Verdict
As a local CLI tool, traditional SRE concerns (uptime, latency SLOs, pager alerts) do not apply. The relevant SRE concerns are release reliability and rollback capability. The main gap is that the publish pipeline is not atomic — a failure mid-publish can leave an inconsistent state with a GitHub release but an un-updated Homebrew formula. Adding rollback steps to the workflow is the priority fix.

## Changes since last assessment
First assessment.

## Remaining recommendations
| Finding | Severity | Recommendation |
|---------|----------|----------------|
| Atomic publish pipeline | High | Add cleanup steps on publish job failure |
| Document rollback procedure | Medium | Add to README or docs/OPERATIONS.md |
| Post-publish smoke test | Medium | macOS runner step after tap update |
| Cache update check | Medium | 24h TTL file in ~/.cache/skills/ |
