# Infosec / Cyber

**Score: 78/100**

## What is being assessed
Threat model covering all attack surfaces: code injection, path traversal, supply chain, credential leakage, and privilege escalation. Assessed from a paranoid posture.

## Methodology
Mapped every input entry point, external call, and file-system write. Traced data flow from CLI argument to disk. Checked for non-alphanumeric character handling in all inputs.

## Findings

### Passing checks
- **No SQL injection surface** — no database used.
- **No XSS surface** — no HTML output or web server.
- **No CSRF surface** — no stateful HTTP sessions.
- **No SSRF surface** — only outbound call is to a hardcoded GitHub API URL (`api.github.com/repos/simonski/skills/releases/latest`). URL is not user-controlled. (internal/version/version.go:14)
- **No command injection** — `exec.Command` is never called. Grep confirms zero occurrences.
- **Path traversal prevented** — skill ID is constrained to catalog keys (alphanumeric + hyphens); `filepath.Join` is used throughout; the `.skills/` prefix is always prepended. (internal/project/project.go:22,26)
- **No credential disclosure** — no secrets written to disk or logs. TAP_TOKEN only appears in the workflow YAML as a secret reference, never echoed.
- **Embedded FS** — catalog content is compiled into the binary via `//go:embed`; runtime catalog files cannot be tampered with by unprivileged users. (internal/catalog/catalog.go:14)
- **Minimal privilege** — the binary only writes to `.skills/` within the working directory. No root operations, no system-wide writes.

### Issues found
| Finding | Severity | Location | Recommendation |
|---------|----------|----------|----------------|
| Skill ID not validated against allowlist before filesystem write | Medium | cmd/add.go:35, internal/project/project.go:79 | Validate `id` matches `^[a-z0-9][a-z0-9-]*$` before any file operation; `parseSkillArg` currently only checks for empty string |
| Supply chain: no SLSA provenance or SBOM generated | Medium | .github/workflows/publish.yml | Add `goreleaser` or `cyclonedx-gomod` to generate SBOM; attach to GitHub release |
| GitHub release API response is not verified (no signature/checksum of tag) | Low | internal/version/version.go:29 | Document this is advisory-only (update check); acceptable for this use case |
| `make publish` clones tap repo over HTTPS with PAT in URL — PAT visible in process list | Low | Makefile:publish target | Use git credential helper or `GIT_ASKPASS` instead of embedding token in clone URL |
| No integrity check on installed skill files after `skills add` | Low | cmd/add.go, internal/project/project.go | Consider storing SHA256 in front matter and verifying on `skills ls` |

## Verdict
The threat surface is genuinely small. The most meaningful risk is the absence of an ID allowlist before filesystem writes — a carefully crafted skill ID (e.g. containing `../`) could escape the `.skills/` directory. This should be patched. Supply chain hygiene (SBOM, SLSA) is missing but common for projects at this stage.

## Changes since last assessment
First assessment.

## Remaining recommendations
| Finding | Severity | Recommendation |
|---------|----------|----------------|
| Validate skill ID format | Medium | Regex `^[a-z0-9][a-z0-9-]*$` in parseSkillArg and/or project.Install |
| Generate SBOM on release | Medium | Add cyclonedx-gomod to publish workflow |
| Sanitise PAT from clone URL | Low | Use GIT_ASKPASS or a netrc credential helper |
