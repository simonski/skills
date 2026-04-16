# Performance

**Score: 78/100** (was 82)

## What is being assessed
Performance review checks startup latency, repeated work, scaling behavior, and whether any code path can grow unexpectedly with user input. Good looks like fast command startup, bounded I/O, and no unnecessary repeated reads or network work.

## Methodology
Reviewed startup flow, search/list implementations, embedded catalog access, and release lookups. Focused on user-visible latency rather than micro-optimizations.

## Findings

### Passing checks
- The catalog is embedded, so catalog reads avoid disk and network I/O (`internal/catalog/catalog.go:16-17,57-79`).
- Search is a simple in-memory scan over a 10-skill catalog, which is acceptable at current scale (`cmd/search.go:24-53`).
- Version lookups use a bounded 5-second timeout instead of an unbounded client (`internal/version/version.go:19-21`).

### Issues found
| Finding | Severity | Location | Recommendation |
|---|---|---|---|
| Every command except `version` blocks on a network update check before doing useful work | High | `cmd/root.go:25-31`, `internal/version/version.go:19-45` | Move the release check off the hot path and/or cache results with a TTL. |
| `ls` calls `project.Get` once per catalog skill instead of loading installed skills once | Medium | `cmd/ls.go:47-63`, `internal/project/project.go:65-78` | Read installed skills in one pass and build an ID-to-version map before rendering the table. |
| Release checks are never cached, so repeated local usage causes repeated GitHub API traffic | Medium | `internal/version/version.go:19-45` | Cache the latest successful release value locally for a short interval. |

## Verdict
Current scale hides most performance issues, but startup latency is still unnecessarily coupled to the network. The biggest win would be eliminating synchronous release lookups from routine commands.

## Changes since last assessment
- No caching or async work was added around update checks.
- Catalog growth remained small enough that the linear scans are still acceptable.

## Remaining recommendations
| Finding | Severity | Recommendation |
|---|---|---|
| Blocking update check | High | Make release checks async or cached. |
| Repeated per-skill filesystem reads in `ls` | Medium | Load installed state once per command. |
| No release-check cache | Medium | Add a small local cache with expiry. |
