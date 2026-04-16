# InfoSec / Cyber

**Score: 55/100** (was 78)

## What is being assessed
This category treats the repo like an attack surface: local filesystem inputs, outbound network calls, release supply chain, and any route from untrusted input to a dangerous sink. Good looks like explicit trust boundaries, strong input allowlists, and no easy privilege escalation paths.

## Methodology
Built a threat model around the CLI's three real surfaces: user-supplied skill IDs, local `.skills/` filesystem state, and outbound release publishing. Read `cmd/*.go`, `internal/project/project.go`, `internal/version/version.go`, and `.github/workflows/publish.yml`, and cross-checked for command execution, SSRF, and HTML/JS sinks.

## Findings

### Passing checks
- There is no SQL layer, template renderer, or HTTP server, which removes SQLi, XSS, and CSRF from the live threat surface (`main.go:1-7`, `cmd/root.go:16-53`).
- The outbound release check targets one hard-coded GitHub API URL rather than a user-controlled destination (`internal/version/version.go:13-21`).
- Catalog content is embedded into the binary at build time, reducing runtime tampering risk (`internal/catalog/catalog.go:16-17`).
- No command execution path exists in the application code; the CLI performs file and HTTP operations only (`cmd/root.go:16-53`, `internal/version/version.go:19-45`).

### Issues found
| Finding | Severity | Location | Recommendation |
|---|---|---|---|
| The local filesystem trust boundary is broken: `skills rm ..` resolves outside `.skills/` and reaches `os.RemoveAll` | Critical | `cmd/rm.go:25-37`, `internal/project/project.go:95-103` | Treat skill IDs as hostile input and validate before path construction; add regression tests for `.`, `..`, and slash-containing IDs. |
| Project read paths also rely on unchecked IDs, so hostile input can probe outside the project root | High | `internal/project/project.go:25-28`, `internal/project/project.go:65-78` | Normalize and validate IDs at the package boundary instead of relying on callers. |
| Release provenance stops at an SBOM; there is no artifact signing or provenance attestation | Medium | `.github/workflows/publish.yml:53-86` | Add signed checksums or provenance attestation for published binaries. |

## Verdict
The remote attack surface is small, but the local attack surface is more dangerous than the earlier reports suggested. The unchecked path join in project operations is the kind of low-level issue that can turn routine CLI use into destructive behavior.

## Changes since last assessment
- Supply-chain visibility improved via SBOM generation (`.github/workflows/publish.yml:53-56`).
- The main regression is analytical rather than new code: the old report misclassified path traversal as solved, but the current implementation still permits it (`internal/project/project.go:95-103`).

## Remaining recommendations
| Finding | Severity | Recommendation |
|---|---|---|
| Local path escape via skill IDs | Critical | Lock the project package behind a strict allowlist and add negative tests. |
| Missing binary provenance | Medium | Publish signed checksums or provenance alongside release assets. |
| Release credential exposure risk | Medium | Remove tokenized clone URLs from the publish path. |
