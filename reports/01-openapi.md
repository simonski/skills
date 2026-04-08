# OpenAPI

**Score: N/A**

## What is being assessed
OpenAPI spec coverage, route/operationId drift, schema correctness, and generated code integrity.

## Methodology
Searched for any .yaml/.json OpenAPI specs, route definitions, and generated code stubs across the entire codebase.

## Findings

### Passing checks
- N/A — this is a CLI tool with no HTTP API surface.

### Issues found
| Finding | Severity | Location | Recommendation |
|---------|----------|----------|----------------|
| No OpenAPI spec exists | N/A | — | Not applicable for a CLI tool |

## Verdict
This category is not applicable. The project is a pure CLI binary with no HTTP API, no routes, and no generated server code. Score is not assigned.

## Changes since last assessment
First assessment.

## Remaining recommendations
None applicable.
