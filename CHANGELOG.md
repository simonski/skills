# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.3] - 2025-01-01

### Added
- `skills update` command — shows what would be updated (dry-run by default); use `-y` to apply
- `skills init` interactive wizard — auto-detects installed coding agents and pre-selects relevant skills
- GitHub Actions CI/CD pipeline — test + publish jobs with Homebrew tap auto-update
- Versioned folder-per-skill catalog structure (`skills/<id>/<version>.md`)
- `skills versions <id>` command to list all available versions of a skill
- `skills get <id>` command to print full skill content to stdout

### Changed
- Catalog restructured to support multiple versions per skill
- `skills add <id>@<version>` now supports version-pinning

### Fixed
- Semver comparison handles multi-segment versions correctly

## [0.1.1] - 2024-12-01

### Added
- Initial public release
- `skills ls` — list catalog skills with colour-coded installation status
- `skills add <skill-id>` — install a skill into `.skills/`
- `skills rm <skill-id>` — remove an installed skill
- `skills search <term>` — search catalog by ID, description, and content
- `skills version` — print version and check for newer releases
- Homebrew tap distribution via `simonski/homebrew-tap`
- Cross-platform release builds (darwin/amd64, darwin/arm64, linux/amd64, linux/arm64, windows/amd64)
- Embedded skill catalog (go, python, security, docker, git, testing, api-design, code-review)
- Automatic update check on every command invocation

[Unreleased]: https://github.com/simonski/skills/compare/v0.1.3...HEAD
[0.1.3]: https://github.com/simonski/skills/compare/v0.1.1...v0.1.3
[0.1.1]: https://github.com/simonski/skills/releases/tag/v0.1.1
