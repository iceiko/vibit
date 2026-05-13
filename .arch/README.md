# Architecture Manifests

Status: Draft v0.1  
Last updated: 2026-05-13
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
  protocol.yaml
  runtime.yaml
  contracts.yaml
  dependencies.yaml
```

This is the first draft. The files describe expected shape before implementation code exists.

`runtime.yaml` records the runtime readiness decisions for the first Go server runtime direction. It points to the Agent Decision Records that govern the first language, server instance model, contract boundary, client protocol, wire format, persistence direction, dependency adoption, and proof slice.

`protocol.yaml` records the first game protocol framework. It defines the WebSocket-framed Protobuf envelope, structured routing fields, session identity model, game target scopes, server-authoritative message rules, compatibility expectations, and implementation boundaries. The human-readable standard is `docs/game-protocol.md`, and the governing decisions are `ADR-0015` and `ADR-0016`.

`ADR-0014` records the first Go runtime package layout and boundary rules. The planned Go module root is `runtime/`, with process startup under `runtime/cmd/vibit-server/`, application orchestration under `runtime/internal/app/`, platform adapters under `runtime/internal/platform/`, handwritten domain runtime logic under `runtime/internal/modules/<module>/`, generated Go outputs under `runtime/internal/generated/`, SQL-first PostgreSQL migrations under `runtime/migrations/postgres/`, and Protobuf source files under repository-root `proto/`.

`contracts.yaml` registers public command, query, event, error, and permission contract source files. Contract files live under `contracts/` and are semantic source artifacts, not generated output. The first Protobuf wire schemas live under `proto/` and must align with these semantic contracts.

`buf.yaml` and `buf.gen.yaml` configure Protobuf source discovery, linting, breaking checks, and planned Go generation output. They are root-level generation configuration rather than architecture manifests, but `.arch/protocol.yaml` and `.arch/runtime.yaml` point to them because agents must read them before protocol generation.

`dependencies.yaml` records foundational dependency decision slots. It identifies dependency categories that need adoption records before implementation imports or requires concrete packages.

The first accepted Go runtime dependencies are recorded by `ADR-0013`:

- `github.com/coder/websocket` for the platform WebSocket transport adapter.
- `google.golang.org/protobuf`, `protoc-gen-go`, and Buf CLI for Protobuf runtime, generation, linting, breaking checks, formatting, and orchestration.
- `github.com/jackc/pgx/v5` for PostgreSQL platform persistence adapters.
- `github.com/pressly/goose/v3` for SQL-first migration tooling.

S3 client tooling, MinIO deployment, observability, and external Go test framework adoption remain deferred until concrete runtime needs require them.

The first Go runtime skeleton exists under `runtime/`, but server business implementation has not started. Agents should update the relevant manifests before implementation and keep third-party transport, protocol, persistence, and migration dependencies inside their declared owner packages. Protobuf generation output remains planned under `runtime/internal/generated/proto/`; do not create or edit generated Go Protobuf files by hand.

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
6. Read `.arch/protocol.yaml` before adding or changing `.proto` files, WebSocket protocol handlers, generated protocol output, or client/server protocol rules.
7. Read `.arch/contracts.yaml` before adding or changing public contracts.
8. Read `.arch/dependencies.yaml` before adding foundational dependencies.
9. Read the affected module's `module.yaml`, when it exists.
10. Update manifests before implementation when public architecture changes.

If a manifest is missing information needed for a safe change, update the manifest or document the gap.

Use `ADR-0012` for the decision authority boundary. After maintainer authorization, agents may professionally evaluate technical sub-decisions inside ratified directions, but must still ask before changing product direction, constitutional principles, runtime language, primary protocol direction, persistence direction, major architecture patterns, module ownership, breaking contracts, validation or permission strength, licensing-risk acceptance, hosting, cost, operations, or vendor-lock-in commitments.

## Verification Direction

These manifests should eventually power checks similar to:

```bash
vibit check architecture
vibit check module <module>
vibit check contracts
vibit check protocol
vibit check change <change-id>
```

Until those commands exist, agents must record architecture verification as not available.
