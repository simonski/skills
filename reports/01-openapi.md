# OpenAPI

**Score: N/A** (was N/A)

## What is being assessed
OpenAPI coverage checks whether a project exposes an HTTP API with a maintained spec, matching routes, stable operation IDs, and trustworthy generated code. Good looks like zero drift between the spec and the runtime surface.

## Methodology
Read the CLI entrypoint and command wiring, then searched for API specs, HTTP routers, handlers, and generated server code in the repository. Files checked included `main.go`, `cmd/root.go`, and `internal/version/version.go`.

## Findings

### Passing checks
- The binary entrypoint only calls the Cobra command tree; there is no server bootstrap path (`main.go:1-7`).
- Root wiring registers CLI subcommands only; no HTTP routes, middleware, or handlers are present (`cmd/root.go:16-53`).
- The only HTTP usage is an outbound GET to GitHub releases for version checks, not an inbound API surface (`internal/version/version.go:13-45`).

### Issues found
| Finding | Severity | Location | Recommendation |
|---|---|---|---|
| No OpenAPI surface exists in this repository | N/A | `main.go:1-7`, `cmd/root.go:16-53` | Keep this category N/A unless the project adds an HTTP API. |

## Verdict
This category is not applicable. `skills` is a local CLI, not a network service, so there is no spec-to-route drift to assess.

## Changes since last assessment
- No HTTP API has been introduced; the project remains CLI-only.

## Remaining recommendations
| Finding | Severity | Recommendation |
|---|---|---|
| Future API work would need a spec-first workflow | N/A | Reassess this category only if the project adds HTTP endpoints or generated API code. |
