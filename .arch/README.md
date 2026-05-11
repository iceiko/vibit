# Architecture Manifests

Status: Draft v0.1  
Last updated: 2026-05-12  
Scope: Machine-readable architecture entry point for vibit

This directory contains architecture manifests for agents, humans, generators, and future verification commands.

The manifests are not decorative documentation. They are intended to become executable architecture context.

## Purpose

The `.arch/` directory should answer the questions an agent must resolve before changing code:

- What modules exist?
- What does each module own?
- Which dependencies are allowed?
- Which contracts define public behavior?
- Which events, commands, queries, errors, and permissions are registered?
- Which files are generated?
- Which tests or checks prove architecture rules?

## Current Files

```text
.arch/
  README.md
  README.zh-CN.md
  modules.yaml
  conventions.yaml
```

This is the first draft. The files describe expected shape before implementation code exists.

## Expected Future Files

```text
.arch/dependencies.yaml
.arch/commands.yaml
.arch/queries.yaml
.arch/events.yaml
.arch/errors.yaml
.arch/permissions.yaml
.arch/test-matrix.yaml
.arch/generation.yaml
```

The project should add these when the first prototype needs them.

## Agent Rules

Before changing implementation code, agents should:

1. Read `CONSTITUTION.md`.
2. Read `AGENTS.md`.
3. Read `.arch/modules.yaml`.
4. Read `.arch/conventions.yaml`.
5. Read the affected module's `module.yaml`, when it exists.
6. Update manifests before implementation when public architecture changes.

If a manifest is missing information needed for a safe change, update the manifest or document the gap.

## Verification Direction

These manifests should eventually power checks similar to:

```bash
vibit check architecture
vibit check module <module>
vibit check contracts
vibit check change <change-id>
```

Until those commands exist, agents must record architecture verification as not available.
