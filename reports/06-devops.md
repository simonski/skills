# DevOps

**Score: 74/100** (was 64)

## What is being assessed
DevOps covers build reproducibility, CI/CD correctness, release packaging, secret handling, and whether the pipeline catches likely regressions before publish. Good looks like reliable builds, strong verification, and automation that mirrors the real release path.

## Methodology
Read `Makefile`, the GitHub Actions workflow, release packaging, and supporting docs. Cross-checked the declared workflow against a local test and coverage run.

## Findings

### Passing checks
- The Makefile exposes clear test, lint, setup, release, and publish targets (`Makefile:17-24,55-68,89-115`).
- CI uses `go-version-file: go.mod`, so the runner follows the module's declared Go version (`.github/workflows/publish.yml:12-16,45-48`).
- The test job now verifies modules, runs `go vet`, executes coverage, and runs `govulncheck` before publish (`.github/workflows/publish.yml:17-28`).
- Release automation builds archives for all target platforms and generates an SBOM (`.github/workflows/publish.yml:50-56,73-86`).
- Secrets setup for the tap token is documented (`docs/SECRETS.md:7-28`).

### Issues found
| Finding | Severity | Location | Recommendation |
|---|---|---|---|
| The repo's current coverage is below the workflow's configured 40% threshold, so the pipeline is stricter than the codebase currently satisfies | High | `.github/workflows/publish.yml:20-24` | Raise coverage in `cmd/` immediately or lower the threshold temporarily to reflect the current baseline. |
| `make build` mutates `VERSION`, which makes local builds non-reproducible | Medium | `Makefile:41-47` | Split version bumping into a separate target and keep `build` purely compilational. |
| CI does not run `staticcheck` even though local linting expects it when installed | Medium | `Makefile:65-68`, `.github/workflows/publish.yml:17-28` | Add `staticcheck ./...` to the test job. |
| Publish has no rollback path if the Homebrew tap update fails after the release is created | Medium | `.github/workflows/publish.yml:67-114` | Add cleanup or rollback steps so publish is atomic. |

## Verdict
DevOps is much better than the previous baseline because the repo now has real verification and SBOM generation. The remaining weakness is that release automation is still a bit brittle, and the current codebase does not yet meet its own coverage gate.

## Changes since last assessment
- `go mod verify`, `go vet`, coverage gating, and `govulncheck` were added to CI (`.github/workflows/publish.yml:17-28`).
- Release SBOM generation was added (`.github/workflows/publish.yml:53-56`).
- `make build` still bumps `VERSION`, so the local build experience remains awkward (`Makefile:41-47`).

## Remaining recommendations
| Finding | Severity | Recommendation |
|---|---|---|
| Coverage gate mismatch | High | Bring `cmd` coverage up before the next release or temporarily align the threshold to reality. |
| Non-reproducible local build target | Medium | Separate `bump` from `build`. |
| Missing `staticcheck` in CI | Medium | Run the same lint stack in CI that contributors are asked to run locally. |
| Non-atomic publish | Medium | Add rollback and smoke-test steps around release publication. |
