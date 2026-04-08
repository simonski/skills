# Compliance

**Score: 88/100**

## What is being assessed
GDPR (data collection, retention, right to erasure), audit trail integrity, cookie consent, data processing documentation, and license compliance.

## Methodology
Read all source files for personal data handling. Checked go.mod for license-incompatible dependencies. Reviewed version check network call. Checked for cookie/session usage.

## Findings

### Passing checks
- **No personal data collected**: the binary collects no user data, no telemetry, no analytics
- **No cookies or sessions**: CLI tool — no HTTP server, no session storage
- **No data sent to third parties**: the only outbound network call is a GET to `api.github.com/repos/simonski/skills/releases/latest` — no user data included in the request
- **MIT license**: project uses MIT license (stated in Homebrew formula) — permissive, compatible with all dependencies
- **All dependencies are permissive**: cobra (Apache 2.0), fatih/color (MIT), spf13/pflag (BSD-3), mattn/go-colorable (MIT), mattn/go-isatty (MIT)
- **No credentials stored**: the binary stores no personal data locally
- **No right-to-erasure concern**: `skills rm <id>` removes the installed file; `rm -rf .skills/` removes all local data. No remote storage.
- **GDPR Article 5 — data minimisation**: satisfied — zero personal data processed

### Issues found
| Finding | Severity | Location | Recommendation |
|---------|----------|----------|----------------|
| No SBOM generated or published with releases | Low | .github/workflows/publish.yml | Generate with cyclonedx-gomod and attach to GitHub releases for supply chain transparency |
| License not stated in repository root (only in Homebrew formula template) | Low | / | Add a LICENSE file to the repository root |
| go.sum present but not explicitly verified in CI | Low | .github/workflows/publish.yml | Add `go mod verify` step |

## Verdict
Compliance posture is strong. The tool processes no personal data, collects no telemetry, uses only permissive open-source dependencies, and has a clean data footprint. The main gaps are cosmetic: a missing LICENSE file at the repo root and absent SBOM generation.

## Changes since last assessment
First assessment.

## Remaining recommendations
| Finding | Severity | Recommendation |
|---------|----------|----------------|
| Add LICENSE file | Low | Create LICENSE (MIT) at repository root |
| Generate SBOM on release | Low | cyclonedx-gomod in publish workflow |
| go mod verify in CI | Low | Step in test job |
