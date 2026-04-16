# Architect

**Score: 87/100** (was 85)

## What is being assessed
Architecture review checks package boundaries, dependency direction, responsibilities, and whether the design matches the problem size. Good looks like an acyclic dependency graph, clear ownership, and infrastructure complexity proportional to the product.

## Methodology
Read the entrypoint, root command, internal packages, and architecture docs. Compared the actual package dependency direction to the documented DAG.

## Findings

### Passing checks
- The dependency graph is clean and acyclic: `cmd` depends on `internal/catalog`, `internal/project`, and `internal/version`, while internal packages do not import each other (`docs/ARCHITECTURE.md:40-50`, `cmd/root.go:42-53`).
- The repo keeps one responsibility per internal package: embedded catalog, local project storage, and release checking (`CLAUDE.md:24-35`, `docs/ARCHITECTURE.md:23-37`).
- The embedded catalog design is simple and appropriate for a CLI distribution model (`internal/catalog/catalog.go:16-17`, `docs/ARCHITECTURE.md:52-75`).
- Local project state is scoped to `.skills/` rather than spread across user home directories or global config (`internal/project/project.go:12-28`).

### Issues found
| Finding | Severity | Location | Recommendation |
|---|---|---|---|
| The project package exposes path construction and mutation APIs without enforcing the skill-ID invariant | Medium | `internal/project/project.go:25-28,81-103` | Make the storage package the trust boundary and reject invalid IDs internally. |
| Version retrieval is hard-wired to the live GitHub endpoint | Low | `internal/version/version.go:12-21` | Introduce an injectable client or source abstraction for better testability. |
| Version comparison logic is duplicated across packages instead of owned by a single domain helper | Low | `internal/catalog/catalog.go:129-150`, `internal/version/version.go:71-105` | Consolidate semver behavior into one internal module. |

## Verdict
The architecture is a good fit for the size of the tool: small, acyclic, and easy to trace. The main structural gap is that some safety invariants live in callers instead of the storage boundary that should own them.

## Changes since last assessment
- Supporting docs now match the package layout much better (`docs/ARCHITECTURE.md:1-100`).
- No architectural bloat was introduced; growth stayed within the same simple command/internal model.

## Remaining recommendations
| Finding | Severity | Recommendation |
|---|---|---|
| Project package trusts callers too much | Medium | Validate skill IDs inside `internal/project`. |
| Live-bound version source | Low | Allow tests to stub or replace the release source. |
| Duplicated version rules | Low | Centralize semver behavior. |
