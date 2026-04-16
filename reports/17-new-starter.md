# New Starter

**Score: 74/100** (was 58)

## What is being assessed
This category measures how quickly a new engineer can understand the repo, run checks, and make a safe first change. Good looks like a clear reading order, accurate setup docs, and explicit warnings about common pitfalls.

## Methodology
Walked the repo as a new contributor starting from the README, then followed supporting docs into contribution, architecture, and release information. Compared that onboarding path to the actual build and test targets.

## Findings

### Passing checks
- A new contributor can find install, usage, and build/test commands from the README (`README.md:8-38,108-130`).
- Contribution workflow is now documented, including branch naming and commit style (`CONTRIBUTING.md:16-48`).
- The Makefile exposes setup, test, and lint targets clearly (`Makefile:17-24,55-68`).
- Architecture and secrets docs now exist and reduce hidden knowledge (`docs/ARCHITECTURE.md:1-99`, `docs/SECRETS.md:7-28`).
- `CONTRIBUTING.md` explicitly warns that `make build` bumps `VERSION` and suggests a plain `go build` for local work (`CONTRIBUTING.md:12-15`).

### Issues found
| Finding | Severity | Location | Recommendation |
|---|---|---|---|
| There is still no dedicated onboarding guide that ties the docs together in a clear reading order | Medium | `docs/ONBOARDING.md` (missing) | Add `docs/ONBOARDING.md` with reading order, setup, testing, and release gotchas. |
| README examples for installed files are inaccurate, which is exactly the kind of mismatch that trips up first-time contributors | Medium | `README.md:50-61`, `internal/project/project.go:25-28` | Fix the README to match the real `.skills/<id>/SKILL.md` layout. |
| The changelog does not reflect the current repo version, which weakens release-context discovery for new contributors | Low | `CHANGELOG.md:10-26`, `VERSION:1` | Update the changelog through `0.1.6`. |

## Verdict
Onboarding is materially better than before because the repo now has contribution, architecture, and secrets docs. The last missing piece is a single curated onboarding document that gives newcomers a reliable order of operations.

## Changes since last assessment
- `CONTRIBUTING.md`, `docs/ARCHITECTURE.md`, and `docs/SECRETS.md` were added.
- The `make build` gotcha is now documented in contribution guidance.
- A dedicated onboarding guide still has not been created.

## Remaining recommendations
| Finding | Severity | Recommendation |
|---|---|---|
| Missing onboarding guide | Medium | Add `docs/ONBOARDING.md` and link it from README and CONTRIBUTING. |
| README implementation drift | Medium | Correct the installed-file path examples. |
| Stale release context | Low | Keep `CHANGELOG.md` current with `VERSION`. |
