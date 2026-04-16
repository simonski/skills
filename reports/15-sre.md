# SRE

**Score: 58/100** (was 52)

## What is being assessed
SRE review checks operational safety for shipping and supporting the tool: release reliability, rollback readiness, graceful degradation, and basic runbook quality. Good looks like reproducible releases, safe failure modes, and a documented recovery path when publish steps go wrong.

## Methodology
Reviewed the GitHub Actions workflow, Makefile publish path, release packaging, and operational docs. Focused on failure modes in publish rather than service uptime, since this is a CLI.

## Findings

### Passing checks
- Publish is gated on the test job via `needs: test` (`.github/workflows/publish.yml:30-33`).
- The workflow avoids recursive publish loops by skipping release commits (`.github/workflows/publish.yml:34-36`).
- Network update failures degrade gracefully at runtime instead of breaking the command (`cmd/root.go:56-60`).
- The Homebrew formula has a post-install test (`packaging/homebrew/skills.rb.template:31-33`).

### Issues found
| Finding | Severity | Location | Recommendation |
|---|---|---|---|
| Publish is not atomic: the workflow creates the GitHub release before the Homebrew tap update finishes | High | `.github/workflows/publish.yml:67-114` | Add rollback/cleanup steps if tap publication fails after tag or release creation. |
| There is no operational runbook for release rollback or failed publish recovery | Medium | `docs/OPERATIONS.md` (missing) | Add an operations doc covering retry, cleanup, and rollback steps. |
| The workflow has no post-release smoke test for the published Homebrew formula | Medium | `.github/workflows/publish.yml:88-114` | Install from the tap on a fresh runner after update and run a simple command. |

## Verdict
Operationally, this repo is fine for a small CLI but still too optimistic about publish success. The release path needs one more pass for rollback and smoke testing before it is truly low-maintenance.

## Changes since last assessment
- SBOM generation improved release hygiene (`.github/workflows/publish.yml:53-56`).
- The main publish-sequencing risks remain unchanged.

## Remaining recommendations
| Finding | Severity | Recommendation |
|---|---|---|
| Non-atomic publish | High | Add rollback and cleanup steps around tag/release creation. |
| Missing operations doc | Medium | Create `docs/OPERATIONS.md` with recovery guidance. |
| No smoke test | Medium | Verify the brewed binary after tap publication. |
