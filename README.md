# skills

Agentic skills manager for your project.

`skills` is a single-binary Go CLI that manages AI-agent skill definitions (SKILL.md files) in your project.
A growing catalog of reusable skills is embedded in the binary — install only the ones you need.

## Installation

```bash
brew install simonski/tap/skills
```

Or build from source:

```bash
make build          # produces dist/skills
make install        # installs to $GOPATH/bin
```

## Usage

```
skills                    Show usage
skills ls                 List catalog skills and show installation status
skills add <skill-id>     Add a skill to the current project
skills rm  <skill-id>     Remove a skill from the current project
skills search <term>      Search the catalog for matching skills
skills version            Print version and check for updates
```

### `skills ls`

Lists all skills in the catalog and indicates each skill's status in the current project:

| Colour | Status |
|--------|--------|
| 🟢 Green  | **INSTALLED** — latest version is installed |
| 🟡 Yellow | **UPDATE AVAILABLE** — an older version is installed |
| 🔴 Red    | **NOT INSTALLED** — skill is not present in this project |

### `skills add <skill-id>`

Installs (or updates) a skill into `.skills/<skill-id>.md` in the current directory.

```bash
skills add go
skills add security
```

### `skills rm <skill-id>`

Removes `.skills/<skill-id>.md` from the current project.

```bash
skills rm go
```

### `skills search <term>`

Searches skill IDs, descriptions, and content for the given term (case-insensitive, any word matches).

```bash
skills search python
skills search "docker container"
```

## Catalog

| ID | Description |
|----|-------------|
| `api-design` | REST API design best practices |
| `code-review` | Code review best practices |
| `docker` | Docker and container best practices |
| `git` | Git conventional commits and branching practices |
| `go` | Go programming best practices |
| `python` | Python programming best practices |
| `security` | Security best practices |
| `testing` | Software testing best practices |

## Building

```bash
make           # show help / usage
make build     # build binary into dist/
make test      # run tests
make lint      # run go vet + staticcheck
make release   # cross-compile for all platforms into dist/
make clean     # remove dist/
```

## Version checking

The binary automatically checks for a newer GitHub release on each run and prints a notice if one is available.
You can also check manually with `skills version`.
