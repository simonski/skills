# Compliance

**Score: 92/100** (was 88)

## What is being assessed
Compliance for this repo means licensing, software inventory, data handling, and whether the release process supports downstream review. Good looks like clear licensing, minimal data collection, dependency visibility, and documented secret usage.

## Methodology
Reviewed the license, dependency manifest, secrets docs, version-check behavior, and release workflow. Focused on what the binary stores, transmits, and publishes.

## Findings

### Passing checks
- The repository has an MIT license at the root (`LICENSE:1-21`).
- The module dependency set is small and clearly declared (`go.mod:1-16`).
- The only outbound network request is a GitHub release lookup, and it does not send project data (`internal/version/version.go:13-27`).
- Secrets for the publish path are documented (`docs/SECRETS.md:7-28`).
- Release automation now generates and publishes a CycloneDX SBOM (`.github/workflows/publish.yml:53-56,78-84`).

### Issues found
| Finding | Severity | Location | Recommendation |
|---|---|---|---|
| Release assets are not signed and no provenance attestation is published | Medium | `.github/workflows/publish.yml:73-86` | Add signed checksums or provenance attestations for release artifacts. |
| The repo does not explicitly document its local-data footprint for privacy-conscious users | Low | `README.md:50-61`, `internal/project/project.go:25-28` | Add one short note explaining that state is stored only in `.skills/<id>/SKILL.md` under the project root. |

## Verdict
Compliance posture is strong for a local CLI: clear license, small dependency surface, documented secrets, and an SBOM in the release flow. The next maturity step is artifact provenance rather than basic legal or data-handling cleanup.

## Changes since last assessment
- The repo gained a root license file and SBOM publication in CI.
- Secret setup is now documented in `docs/SECRETS.md`.

## Remaining recommendations
| Finding | Severity | Recommendation |
|---|---|---|
| Unsigned release artifacts | Medium | Publish signed checksums or provenance for binaries. |
| Implicit local-data model | Low | Document exactly what state the tool writes and where. |
