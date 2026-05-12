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
  runtime.yaml
  contracts.yaml
  dependencies.yaml
```

This is the first draft. The files describe expected shape before implementation code exists.

`runtime.yaml` records the runtime readiness decisions for the first Go server runtime direction. It points to the Agent Decision Records that govern the first language, server instance model, contract boundary, client protocol, wire format, persistence direction, dependency adoption, and proof slice.

`contracts.yaml` registers public command, query, event, error, and permission contract source files. Contract files live under `contracts/` and are semantic source artifacts, not generated output. Protobuf wire schemas are planned under `proto/` and must align with these semantic contracts.

`dependencies.yaml` records foundational dependency decision slots. It identifies dependency categories that need adoption records before implementation imports or requires concrete packages.

## Expected Future Files

```text
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
5. Read `.arch/runtime.yaml` before changing or creating runtime implementation code.
6. Read `.arch/contracts.yaml` before adding or changing public contracts.
7. Read `.arch/dependencies.yaml` before adding foundational dependencies.
8. Read the affected module's `module.yaml`, when it exists.
9. Update manifests before implementation when public architecture changes.

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
