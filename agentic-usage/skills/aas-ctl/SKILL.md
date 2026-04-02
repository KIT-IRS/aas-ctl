---
name: aas-ctl
description: "Interact with Asset Administration Shells (AAS) via the aas-ctl CLI. Use when: querying AAS repositories, listing shells or submodels, showing submodel elements, searching shells by criteria, creating or updating submodels, navigating AAS structures, configuring AAS profiles, or performing REST operations (GET/POST/PUT/PATCH) on AAS endpoints."
argument-hint: "Describe what AAS data you want to access or modify"
---

# AAS-CTL — Asset Administration Shell CLI

A command-line tool for browsing and interacting with Asset Administration Shell (AAS) repositories. Enables efficient access to shells, submodels, and submodel elements.

## When to Use

- List, show, or search Asset Administration Shells and Submodels
- Navigate AAS structures interactively via discovery
- Create, update, or patch submodels and submodel elements
- Configure connections to different AAS repository environments
- Perform raw REST operations against AAS API endpoints

## Prerequisites

- `aas-ctl` must be installed and available in PATH
- An active profile must be configured pointing to a running AAS repository

## Command Reference

### Configuration

```bash
# List all profiles (active profile is highlighted)
aas-ctl config list

# Create a new profile
aas-ctl config create <profile-name>

# Switch active profile
aas-ctl config select <profile-name>

# Show config file location (~/.aas/config.json)
aas-ctl config show
```

A profile stores: name, base URL, and ports for Discovery (8084), Registry (8082), SmRegistry (8083), Repository (8081), SmRepository (8081), ConceptDescriptions (8081).

### Listing

```bash
# List all shells
aas-ctl aas list

# List all submodels
aas-ctl sm list

# List submodels of a specific shell
aas-ctl sm list --aas <Identifier>

# Shortcut: list everything (shells)
aas-ctl list
```

### Showing Details

```bash
# Show a shell and its submodels
aas-ctl aas show <Identifier>

# Show a submodel and its elements
aas-ctl sm show <Identifier>

# Show a specific submodel element by IDShort
aas-ctl sm show <SmIdentifier> --elementId <ElementIDShort>

# Show a specific submodel element by index
aas-ctl sm show <SmIdentifier> --elementIdx <Index>

# Show only the value of an element
aas-ctl sm show <SmIdentifier> --elementId <ElementIDShort> --value

# Auto-detect shell or submodel
aas-ctl show <Identifier>
```

Identifiers can be a global ID (e.g., `urn:example:aas#001`) or an IDShort (e.g., `MyShell`).

### Searching / Filtering

```bash
# Find shells containing a submodel with a given IDShort
aas-ctl search --sm <SubmodelIDShort>

# Chain filters (AND logic): submodel + element + value
aas-ctl search --sm <SmIDShort> --elementId <ElementIDShort> --value <Value>

# Return only IDs of matching shells
aas-ctl search --sm <SmIDShort> --id
```

### Discovery (Interactive Navigation)

```bash
# Start from a shell — lists its submodels
aas-ctl discover <ShellIdentifier>

# Navigate deeper — list submodel elements
aas-ctl discover <ShellIdentifier> <SubmodelIdentifier>

# Navigate to a specific element
aas-ctl discover <ShellIdentifier> <SubmodelIdentifier> <ElementIDShort>
```

### Raw REST Operations

```bash
# GET any endpoint
aas-ctl get <endpoint>

# POST data (or pipe from stdin)
aas-ctl post <endpoint> '<json-data>'
cat data.json | aas-ctl post <endpoint>

# PUT data (or pipe from stdin)
aas-ctl put <endpoint> '<json-data>'

# PATCH data
aas-ctl patch <endpoint> '<json-data>'
```

### Output Flags (mutually exclusive)

| Flag | Effect |
|------|--------|
| `--id` | Print only the global ID |
| `--url` | Print only the URL |
| `--json` | Print full JSON representation |

## Output Formats

**Default format:**
```
urn:example:aas#001<AssetAdministrationShell> MyShell
```

**Verbose/show format (Submodel):**
```
urn:example:sm#001<Submodel> DataSM
 0  Temperature<Property>: 35.5<xs:string>
 1  Pressure<MultiLanguageProperty>: 1.2<en>
 2  Range<Range>: 0.0-100.0<xs:float>
```

## Typical AI Agent Workflows

### Read a specific value from an AAS
```bash
aas-ctl sm show <SubmodelIDShort> --elementId <ElementIDShort> --value
```

### Find all shells matching criteria
```bash
aas-ctl search --sm <SubmodelIDShort> --elementId <ElementIDShort> --value <TargetValue> --id
```

### Get full JSON for programmatic processing
```bash
aas-ctl sm show <SubmodelIdentifier> --json
```

### Update a submodel element value
```bash
aas-ctl patch <endpoint> '<new-value>'
```
