# SDLC Assessment Summary

**Date:** 2026-04-16  
**Project:** `skills` — Go CLI for managing AI-agent skill definitions  
**Overall Score:** 72/100

## Score table

| # | Category | Current | Previous | Delta |
|---|---|---:|---:|---:|
| 01 | OpenAPI | N/A | N/A | — |
| 02 | Security | 48 | 72 | -24 |
| 03 | InfoSec / Cyber | 55 | 78 | -23 |
| 04 | Idiomatic Go | 82 | 81 | +1 |
| 05 | Idiomatic JavaScript | N/A | N/A | — |
| 06 | DevOps | 74 | 64 | +10 |
| 07 | QA | 45 | 55 | -10 |
| 08 | Tech Lead | 78 | 79 | -1 |
| 09 | Architect | 87 | 85 | +2 |
| 10 | Performance | 78 | 82 | -4 |
| 11 | Database | N/A | N/A | — |
| 12 | Tech Writer | 80 | 68 | +12 |
| 13 | Product Owner | 78 | 74 | +4 |
| 14 | Compliance | 92 | 88 | +4 |
| 15 | SRE | 58 | 52 | +6 |
| 16 | UX Review | 74 | N/A | new |
| 17 | New Starter | 74 | 58 | +16 |

*Scored categories: 14. Overall = 1003 / 14 = 71.6, rounded to **72/100**.*

## Score distribution

| Band | Categories |
|---|---|
| 90+ | Compliance |
| 80-89 | Idiomatic Go, Architect, Tech Writer |
| 70-79 | DevOps, Tech Lead, Performance, Product Owner, UX Review, New Starter |
| 60-69 | None |
| 50-59 | InfoSec / Cyber, SRE |
| Below 50 | Security, QA |

## What changed since last assessment

- CI is materially stronger: `go mod verify`, `go vet`, coverage gating, and `govulncheck` are now in the test job (`.github/workflows/publish.yml:17-28`).
- Release automation now generates and publishes a CycloneDX SBOM (`.github/workflows/publish.yml:53-56,78-84`).
- Contributor docs improved: `CONTRIBUTING.md`, `CHANGELOG.md`, `docs/ARCHITECTURE.md`, and `docs/SECRETS.md` now exist.
- The old reports overstated path safety. `project.Remove` still trusts raw IDs and can escape `.skills/` with `..` segments (`internal/project/project.go:95-103`), which drags Security and InfoSec down sharply.
- Quality remains uneven: local total coverage is **37.8%** and `cmd` coverage is **13.9%**, below the workflow's configured 40% bar.

## Cumulative improvement

| Category | Original | Current |
|---|---:|---:|
| DevOps | 64 | 74 |
| Tech Writer | 68 | 80 |
| Compliance | 88 | 92 |
| New Starter | 58 | 74 |
| Security | 72 | 48 |
| QA | 55 | 45 |

## Key metrics

| Metric | Value |
|---|---|
| Current version | `0.1.6` |
| Go source files | 18 |
| Test files | 4 |
| Test functions | 27 |
| Overall coverage | 37.8% |
| `cmd` coverage | 13.9% |
| `internal/catalog` coverage | 76.3% |
| `internal/project` coverage | 85.1% |
| `internal/version` coverage | 45.8% |
| Catalog skills | 10 |
| GitHub workflows | 1 |
| CI jobs | 2 |
| SBOM generation | Yes |
| Largest Go file | `cmd/init.go` (308 lines) |
| Database migrations | N/A |
| Database indexes | N/A |

## Remaining action items

| Priority | Category | Location | Action |
|---|---|---|---|
| P0 | Security | `cmd/rm.go:25-37`, `internal/project/project.go:95-103` | Validate skill IDs against an allowlist before any path join or filesystem mutation; reject `.` and `..` segments explicitly. |
| P1 | QA | `cmd/add.go:36-75`, `cmd/init.go:93-308`, `cmd/ls.go:29-66`, `cmd/rm.go:25-37`, `cmd/search.go:24-55`, `cmd/version.go:21-40`, `cmd/versions.go:25-38` | Add tests for the untested command runners and raise total coverage above the configured threshold. |
| P1 | DevOps | `.github/workflows/publish.yml:20-24` | Align the coverage threshold with current reality or, preferably, increase real coverage before the next release. |
| P1 | Performance | `cmd/root.go:25-31`, `internal/version/version.go:19-45` | Move the update check off the hot path and cache release lookups. |
| P2 | Tech Writer | `README.md:50-61`, `README.md:97-106`, `CHANGELOG.md:10-26` | Fix README drift for install paths, refresh the catalog table, and update the changelog past `0.1.3`. |
| P2 | SRE | `.github/workflows/publish.yml:67-114` | Make publish atomic with rollback/cleanup steps and add a post-release smoke test. |
| P2 | New Starter | `docs/ONBOARDING.md` (missing) | Add a single onboarding guide with reading order, setup flow, and common gotchas. |
