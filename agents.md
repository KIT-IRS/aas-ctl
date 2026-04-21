# AGENTS.md - aas-ctl repository guide for coding agents

This document summarizes how this repository works, based on `README.md` and the implementation in `cmd/`, `utils/`, and `config/`.

## 1) Project at a glance

- Name: `aas-ctl`
- Type: Cobra-based Go CLI for browsing and editing Asset Administration Shell (AAS) repositories.
- Entry point: `main.go` -> `cmd.Execute()`.
- Current version constant: `0.9.0` (`cmd/version.go`).
- Core dependencies:
  - `github.com/spf13/cobra`
  - `github.com/aas-core-works/aas-core3.0-golang`
  - `github.com/fatih/color`

## 2) Main responsibilities of the CLI

The tool is designed to:

- list shells and submodels
- show shell/submodel/submodel-element content
- search shells by submodel/element/value criteria
- discover nested structures interactively
- run raw REST operations (`get`, `post`, `put`, `patch`)
- manage repository profiles (`config create`, `config list`, `config select`, `config show`)

## 3) Command map (implemented behavior)

### Root

- `aas-ctl` (root command)
- Aliases and shortcuts:
  - `aas-ctl list` is the shell list command (same implementation as `aas list`)

### AAS commands

- `aas-ctl aas list`
- `aas-ctl aas show <identifier>`

Persistent output flags on `aas`:

- `--id`
- `--url`
- `--json`

These are mutually exclusive.

### Submodel commands

- `aas-ctl sm list [--aas <identifier>]`
- `aas-ctl sm show <identifier> [--aas <identifier>] [--elementId <idShort> | --elementIdx <idx>] [--value]`

Persistent flags on `sm`:

- `--aas <identifier>`
- `--id`
- `--url`
- `--json`

`sm show` local flags:

- `--elementId`
- `--elementIdx`
- `--value`

`--elementId` and `--elementIdx` are mutually exclusive.

### Generic show

- `aas-ctl show <identifier>`

Behavior from code (`cmd/show.go`):

1. Try AAS lookup first.
2. If not found as AAS, try Submodel lookup.

So this command auto-detects identifiable type but prefers AAS on ambiguous identifiers.

### Search

- `aas-ctl search [--sm <idShort>] [--elementId <idShort> | --elementIdx <idx>] [--value <value>]`

Output flags also supported on search:

- `--id`, `--url`, `--json` (mutually exclusive)

Filter semantics are AND-based across provided criteria.

### Discovery

- `aas-ctl discover <arg1> [arg2 ...]`

Flags:

- `--url`
- `--json`

These are mutually exclusive.

Discovery resolves a path from an identifiable (shell or submodel) through nested referables.

### Raw REST commands

- `aas-ctl get <url>`
- `aas-ctl post <url> [<body>]`
- `aas-ctl put <url> [<body>]`
- `aas-ctl patch <url> <value>`

For `post`/`put`, if body arg is omitted the command reads stdin.

For `patch`, implementation detail from `utils/patch.go` and `cmd/patch.go`:

- `/$value` is auto-appended if missing.
- value argument is JSON-quoted before sending.

### Config commands

- `aas-ctl config create <name>`
- `aas-ctl config list`
- `aas-ctl config select <name>`
- `aas-ctl config show`

## 4) Config model and profile handling

Config file path:

- `~/.aas/config.json` (resolved via `os.UserHomeDir()` in `config/config.go`)

Config JSON model:

- `activeProfile`
- `profiles[]`

Profile fields:

- `name`
- `url`
- `ports.discovery`
- `ports.registry`
- `ports.sm-registry`
- `ports.repository`
- `ports.sm-repository`
- `ports.concept-descriptions`

Derived endpoints are built in `config/profiles.go`:

- discovery: `<url>:<discovery>/lookup/shells`
- registry: `<url>:<registry>/shell-descriptors`
- sm-registry: `<url>:<sm-registry>/submodel-descriptors`
- repository: `<url>:<repository>/shells`
- sm-repository: `<url>:<sm-repository>/submodels`
- concept-descriptions: `<url>:<concept-descriptions>/concept-descriptions`

Important implementation detail:

- `Config.Save()` writes directly to `~/.aas/config.json` and does not create parent directories.
- If `~/.aas/` does not exist, save operations can fail.

## 5) Data access and resolution behavior

From `utils/get.go`, `utils/query.go`, `utils/discovery.go`:

- Identifier resolution typically tries global ID first, then `idShort`.
- For `idShort` collisions, first match wins (shells/submodels).
- `sm show --aas <shellIdentifier>` resolves submodel by `idShort` within that shell.
- Submodel elements can be addressed by `idShort` or index.
- Nested discovery supports traversal into collections and lists.

## 6) Output behavior

Output formatting logic is centralized in `utils/print.go`.

Modes:

- default human-readable format
- `--id` only ID
- `--url` endpoint URL
- `--json` JSON serialized object
- `--value` (for selected SME types)

Implemented value formatting for:

- `Property`
- `MultiLanguageProperty`
- `Range`
- `SubmodelElementList`
- `SubmodelElementCollection`

Unknown/unsupported element types produce a fallback "formatting is not implemented" message.

## 7) Error and HTTP handling

- HTTP calls use `utils.HttpRequest`.
- Non-2xx responses become `HTTPError` with status and URL.
- Many command handlers use `log.Fatal`, so failures terminate with non-zero exit and message.
- Custom errors exist for missing IDs, missing idShorts, and missing submodels in shell context.

## 8) Agent usage guidance

If you are an AI coding or operations agent using this repository/tool:

1. Prefer `aas-ctl` commands over custom REST calls when interacting with AAS data.
2. Check active profile first when target environment is unclear (`aas-ctl config list`).
3. Do read-only commands before write commands when identifiers are uncertain:
   - `list`, `show`, `discover`, `search`
4. For submodel elements with non-unique names, disambiguate with shell context (`sm ... --aas`).
5. For write operations:
   - use `patch` for value-only changes
   - use `put` for full object replacement
   - use `post` for creation
6. When scripting, prefer JSON mode (`--json`) for machine parsing.

## 9) Practical caveats for agents

- `README.md` states Go `1.24+`, while `go.mod` currently says `go 1.23.5`.
- `config select` does not hard-fail with `log.Fatal` on unknown profile; it prints a message.
- Discovery output may show indexed entries; when passing indices as args, use numeric values.
- `show` command auto-detection order is AAS first, then Submodel.

## 10) Repository pointers

- CLI commands: `cmd/*.go`
- Config model: `config/config.go`, `config/profiles.go`
- Retrieval/query/discovery: `utils/get.go`, `utils/query.go`, `utils/discovery.go`
- Output formatting: `utils/print.go`
- HTTP + REST wrappers: `utils/http.go`, `utils/get.go`, `utils/post.go`, `utils/put.go`, `utils/patch.go`
- Example agent assets: `agentic-usage/skills/aas-ctl/SKILL.md`, `agentic-usage/agents/aas.agent.md`
