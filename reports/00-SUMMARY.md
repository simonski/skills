# SDLC Assessment Summary

**Date:** 2025-07-14  
**Project:** `skills` — Go CLI for managing AI-agent skill definitions  
**Overall Score: 73/100**

---

## Score Table

| # | Category | Score | Notes |
|---|----------|-------|-------|
| 01 | OpenAPI | N/A | No HTTP surface — pure CLI |
| 02 | Security | 72/100 | Input validation gap, no govulncheck in CI |
| 03 | InfoSec / Cyber | 78/100 | Path traversal risk on skill IDs, no allowlist |
| 04 | Idiomatic Go | 81/100 | Clean idioms; duplicate semver logic; blocking HTTP |
| 05 | Idiomatic JavaScript | N/A | No JavaScript — pure CLI |
| 06 | DevOps | 64/100 | No lint, no govulncheck, no go mod verify in CI |
| 07 | QA | 55/100 | `cmd` package at 21.5% coverage — critical gap |
| 08 | Tech Lead | 79/100 | Duplicate semver, oversized runInit |
| 09 | Architect | 85/100 | Clean DAG; missing interfaces for testability |
| 10 | Performance | 82/100 | Blocking `checkForUpdates()` on every command |
| 11 | Database | N/A | No database |
| 12 | Tech Writer | 68/100 | No CHANGELOG, CONTRIBUTING, ONBOARDING, docs/ |
| 13 | Product Owner | 74/100 | Undocumented commands, UX polish gaps |
| 14 | Compliance | 88/100 | No LICENSE file, no SBOM |
| 15 | SRE | 52/100 | Non-atomic publish; no rollback docs or smoke test |
| 16 | UX Review | N/A | CLI UX solid; minor gaps (NO_COLOR, search status) |
| 17 | New Starter | 58/100 | `make build` version-bump gotcha; no ONBOARDING.md |

*N/A categories excluded from overall average. Scored categories: 02–04, 06–10, 12–15, 17 (12 categories).*  
*Overall: (72+78+81+64+55+79+85+82+68+74+88+52+58) / 13 = **73/100***

---

## Score Distribution

| Band | Categories |
|------|-----------|
| 80–89 | Idiomatic Go (81), Architect (85), Performance (82), Compliance (88) |
| 70–79 | Security (72), InfoSec (78), Tech Lead (79), Product Owner (74) |
| 60–69 | DevOps (64), Tech Writer (68) |
| 50–59 | QA (55), SRE (52), New Starter (58) |

---

## Key Metrics

| Metric | Value |
|--------|-------|
| Go source files | ~15 |
| Test files | 6 |
| Overall test coverage | 43% |
| `cmd` package coverage | 21.5% |
| `internal/catalog` coverage | 76.3% |
| `internal/project` coverage | 86.2% |
| CI jobs | 2 (test, publish) |
| GitHub Action workflows | 1 |
| Skills in catalog | ~10 |
| Open tickets | 0 |

---

## What Changed Since Last Assessment

*(First assessment — no prior baseline.)*

- `cmd/update.go` added: `skills update` command with dry-run default and `-y` flag
- `.github/workflows/publish.yml` added: test + publish CI/CD pipeline
- `reports/` directory created with 17 category reports

---

## Prioritised Action Items

### 🔴 High Priority

| Finding | Category | File | Recommendation |
|---------|----------|------|----------------|
| Skill ID not validated against allowlist | Security / InfoSec | `cmd/add.go` | Reject IDs not matching `^[a-z0-9][a-z0-9-]*$` before any filesystem write |
| `cmd` package at 21.5% coverage | QA | `cmd/*.go` | Add tests for `runAdd`, `runRM`, `runLS`, `runSearch`, `parseSkillArg` |
| No lint or govulncheck in CI | DevOps | `.github/workflows/publish.yml` | Add `golangci-lint run` and `govulncheck ./...` to test job |
| No `go mod verify` in CI | DevOps | `.github/workflows/publish.yml` | Add `go mod verify` step before build |
| `checkForUpdates()` blocks every command | Performance | `cmd/root.go` | Move to goroutine with 2s timeout; cache result in `~/.skills/.version-check` |
| `make build` silently bumps VERSION | DevOps / New Starter | `Makefile` | Separate `make bump` from `make build`; document the distinction |
| No ONBOARDING.md | New Starter | (missing) | Create `docs/ONBOARDING.md` — template in report 17 |
| TAP_TOKEN secret undocumented | DevOps | `.github/workflows/publish.yml` | Document required secret in README or docs/ |

### 🟡 Medium Priority

| Finding | Category | File | Recommendation |
|---------|----------|------|----------------|
| Duplicate semver logic | Tech Lead / Go | `internal/catalog/catalog.go:131`, `internal/version/version.go:62` | Extract `internal/semver` package; delete duplicates |
| No CHANGELOG.md | Tech Writer | (missing) | Create `CHANGELOG.md` following Keep a Changelog format |
| No CONTRIBUTING.md | Tech Writer | (missing) | Create `CONTRIBUTING.md` with branch, commit, PR conventions |
| Non-atomic publish pipeline | SRE | `.github/workflows/publish.yml` | Add smoke test step; document rollback procedure |
| No HTTP client injection | Architect | `internal/version/version.go` | Accept `*http.Client` parameter for testability |
| `runInit` too large | Tech Lead | `cmd/init.go` | Split into smaller named functions |
| No LICENSE file | Compliance | (missing) | Add `LICENSE` (MIT or Apache-2.0) at repo root |

### 🟢 Low Priority

| Finding | Category | File | Recommendation |
|---------|----------|------|----------------|
| `NO_COLOR` not respected | UX | `cmd/*.go` | Check `os.Getenv("NO_COLOR")` before emitting ANSI codes |
| `skills search` doesn't show install status | UX / PO | `cmd/search.go` | Add `[installed]` marker in search output |
| `skills rm` lacks `-y` flag | UX | `cmd/rm.go` | Add `--yes` dry-run guard consistent with `update` |
| No SBOM | Compliance | (missing) | Add `make sbom` target using `syft` or `cyclonedx-gomod` |
| Error sentinel inconsistency | Tech Lead | `cmd/*.go` | Define package-level sentinel errors; replace ad-hoc `fmt.Errorf` |

---

## Verdict

The `skills` project is well-structured with clean Go idioms, a sensible package layout, and good internal package coverage. The three areas that most need investment are:

1. **Test coverage** — the `cmd` package is the user-facing surface and is largely untested.
2. **CI hardening** — adding lint and vulnerability scanning would catch regressions automatically.
3. **Documentation** — a new starter currently has no single document explaining how to work on the project.

Addressing the 8 high-priority items would move the overall score from **73 → ~83/100**.
