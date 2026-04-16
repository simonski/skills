# Product Owner

**Score: 78/100** (was 74)

## What is being assessed
Product-owner review checks whether the implemented command set matches the repo's stated goals and whether important user journeys are complete and understandable. Good looks like discoverable workflows, accurate docs, and low-friction paths for common tasks.

## Methodology
Compared the stated product in `README.md` with the command implementations in `cmd/`. Focused on the main user journeys: list, search, add, update, inspect, initialize, and remove.

## Findings

### Passing checks
- The README and root command expose a coherent command set for listing, adding, removing, searching, updating, initializing, and inspecting skills (`README.md:24-37`, `cmd/root.go:42-53`).
- Version pinning is supported for `add` (`cmd/add.go:16-29,84-96`).
- `update` is dry-run by default, which is a good safety choice for bulk changes (`cmd/update.go:16-33,95-97,140-143`).
- `init` can detect agent-related files and guide setup interactively (`cmd/init.go:18-68,109-191`).

### Issues found
| Finding | Severity | Location | Recommendation |
|---|---|---|---|
| The README's installed-file path is wrong, which makes a core user journey inaccurate | Medium | `README.md:50-61`, `internal/project/project.go:25-28` | Update product docs so users know exactly what lands in their project. |
| `search` does not show whether a matching skill is already installed, unlike `ls` | Low | `cmd/search.go:49-53`, `cmd/ls.go:44-63` | Add an install-status column or marker to search results. |
| `rm` performs destructive deletion with no dry-run or confirmation path | Low | `cmd/rm.go:13-37` | Add `--yes` or a confirmation prompt consistent with `update`. |

## Verdict
The product scope is solid and the command set fits the stated goal of managing embedded skills in a project. The remaining gaps are mostly around user trust and clarity rather than missing features.

## Changes since last assessment
- `update`, `init`, `get`, and `versions` are now fully documented in the README (`README.md:24-37`).
- Product docs still lag the implementation on the actual installed file layout.

## Remaining recommendations
| Finding | Severity | Recommendation |
|---|---|---|
| Incorrect install-path docs | Medium | Fix README path examples and descriptions. |
| Search lacks status context | Low | Show whether a search hit is installed or updatable. |
| Remove flow is too sharp-edged | Low | Add a confirmation or explicit apply flag. |
