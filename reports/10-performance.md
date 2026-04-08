# Performance

**Score: 82/100**

## What is being assessed
N+1 I/O patterns, unbounded resource usage, goroutine safety, build context memory, and startup latency for a CLI tool.

## Methodology
Traced all I/O paths. Checked for loops that perform file reads. Measured conceptual startup cost of embed.FS. Assessed HTTP client usage and timeout configuration.

## Findings

### Passing checks
- **Embedded catalog**: all skill content compiled into binary via `//go:embed` — zero disk reads on startup; catalog access is purely in-memory
- **HTTP timeout**: `LatestRelease` uses a 5-second client timeout — no unbounded blocking — internal/version/version.go:18
- **No goroutine leaks**: no goroutines spawned in application code
- **Small memory footprint**: catalog is ~8 skills × ~10KB each ≈ 80KB embedded data — negligible
- **No N+1 catalog queries**: `catalog.All()` iterates skill IDs once then calls `GetVersion` per skill — each is an in-memory FS read, effectively free
- **project.List** reads directory once then reads each .md file — linear in number of installed skills; fine at expected scale (<20 skills)

### Issues found
| Finding | Severity | Location | Recommendation |
|---------|----------|----------|----------------|
| `checkForUpdates()` blocks the main command path with a synchronous HTTP call (up to 5s) | High | cmd/root.go:48 | Run in a background goroutine; print result after command output; use a channel with a short select timeout (500ms) |
| `catalog.All()` calls `skillIDs()` then `Get()` per skill — two FS traversals | Low | internal/catalog/catalog.go:28 | Single-pass: read dir once, parse version from filename, read each file once. Current code is correct but could be simplified |
| No caching of `LatestRelease` response — every command invocation hits the network | Medium | internal/version/version.go, cmd/root.go | Cache result in `~/.cache/skills/latest-version` with a 24-hour TTL |

## Verdict
Performance is excellent for the intended use case. The only material issue is the synchronous update check which adds up to 5 seconds of latency to every command on a slow network. Moving it to a goroutine with a short timeout would make all commands feel instant.

## Changes since last assessment
First assessment.

## Remaining recommendations
| Finding | Severity | Recommendation |
|---------|----------|----------------|
| Async update check | High | Goroutine with 500ms timeout; cache result 24h |
| Cache latest-version response | Medium | Write to ~/.cache/skills/latest-version with TTL |
