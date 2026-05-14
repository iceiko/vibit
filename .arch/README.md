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
  reference.yaml
  work-items.yaml
```

This is the first draft. The files describe expected shape before implementation code exists.

`runtime.yaml` records the runtime readiness decisions for the first Go server runtime direction. It points to the Agent Decision Records that govern the first language, server instance model, contract boundary, client protocol, wire format, persistence direction, dependency adoption, and proof slice.

`protocol.yaml` records the first game protocol framework. It defines the WebSocket-framed Protobuf envelope, structured routing fields, session identity model, game target scopes, server-authoritative message rules, compatibility expectations, and implementation boundaries. The human-readable standard is `docs/game-protocol.md`, and the governing decisions are `ADR-0015` and `ADR-0016`.

`ADR-0014` records the first Go runtime package layout and boundary rules. The planned Go module root is `runtime/`, with process startup under `runtime/cmd/vibit-server/`, application orchestration under `runtime/internal/app/`, platform adapters under `runtime/internal/platform/`, handwritten domain runtime logic under `runtime/internal/modules/<module>/`, generated Go outputs under `runtime/internal/generated/`, SQL-first PostgreSQL migrations under `runtime/migrations/postgres/`, and Protobuf source files under repository-root `proto/`.

`contracts.yaml` registers public command, query, event, error, and permission contract source files. Contract files live under `contracts/` and are semantic source artifacts, not generated output. The first Protobuf wire schemas live under `proto/` and must align with these semantic contracts.

`buf.yaml` and `buf.gen.yaml` configure Protobuf source discovery, linting, breaking checks, and planned Go generation output. They are root-level generation configuration rather than architecture manifests, but `.arch/protocol.yaml` and `.arch/runtime.yaml` point to them because agents must read them before protocol generation.

`docs/generated-output.md` records the generated output standard. `ADR-0017` governs generated output traceability, the `runtime/internal/generated/proto/` ownership rule, and the requirement that generated Go Protobuf files use `protoc-gen-go` markers and source traces back to `proto/`.

`docs/runtime-protocol-adapter.md` records the runtime protocol adapter boundary standard. `ADR-0018` governs the first handoff between WebSocket transport, Protobuf protocol adaptation, application dispatch, generated code, and domain modules.

`dependencies.yaml` records foundational dependency decision slots. It identifies dependency categories that need adoption records before implementation imports or requires concrete packages.

`reference.yaml` records the active reference baseline for game server capability planning. It links `docs/reference-game-server-alignment.md` and `ADR-0019`. Nakama is the primary reference for broad game backend product capability surface. Pitaya is the primary reference for Go game server framework architecture vocabulary. These references guide planning; they do not override vibit's constitution, ADRs, manifests, generated boundaries, or verification commands.

`docs/player-identity-session-boundary.md` records the active boundary standard for player identity, player accounts, authentication, runtime sessions, transport connection metadata, and request identity context. `ADR-0021` governs this boundary. It keeps `player_id` as planned player-domain identity, not authenticated proof from a client envelope, until session validation exists.

`work-items.yaml` records the active work continuation queue. It links `docs/workflow.md` and defines milestones, work items, dependencies, completion traces, and the `next_ready` item that gives maintainer continuation requests a deterministic meaning.

The first accepted Go runtime dependencies are recorded by `ADR-0013`:

- `github.com/coder/websocket` for the platform WebSocket transport adapter.
- `google.golang.org/protobuf`, `protoc-gen-go`, and Buf CLI for Protobuf runtime, generation, linting, breaking checks, formatting, and orchestration.
- `github.com/jackc/pgx/v5` for PostgreSQL platform persistence adapters.
- `github.com/pressly/goose/v3` for SQL-first migration tooling.

S3 client tooling, MinIO deployment, observability, and external Go test framework adoption remain deferred until concrete runtime needs require them.

The first Go runtime skeleton exists under `runtime/`, and the first narrow runtime handoff slice now contains generated Go Protobuf output plus typed application handoff structures. Server business behavior, WebSocket transport, PostgreSQL persistence, migrations, and full application dispatch have not started yet. Agents should update the relevant manifests before implementation and keep third-party transport, protocol, persistence, and migration dependencies inside their declared owner packages. Generated Go Protobuf output lives under `runtime/internal/generated/proto/`; do not create or edit generated Go Protobuf files by hand. Runtime protocol handoff rules are defined in `docs/runtime-protocol-adapter.md`. Run `node tools/vibit check generated` after generated output changes and `node tools/vibit check runtime` after runtime boundary changes.

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
7. Read `docs/runtime-protocol-adapter.md` before adding WebSocket transport code, Protobuf runtime adapter code, application dispatch code, or domain runtime handlers.
8. Read `.arch/contracts.yaml` before adding or changing public contracts.
9. Read `.arch/dependencies.yaml` before adding foundational dependencies.
10. Read `.arch/reference.yaml` and `docs/reference-game-server-alignment.md` before adding new game server capability families, runtime subsystems, social/realtime features, matchmaking, match runtime, cluster/RPC work, or operational surfaces.
11. Read `docs/player-identity-session-boundary.md` before adding player, account, authentication, session, permission, or request identity behavior.
12. Read `.arch/work-items.yaml` and `docs/workflow.md` before interpreting "continue" or multi-step continuation requests.
13. Read the affected module's `module.yaml`, when it exists.
14. Update manifests before implementation when public architecture changes.

If a manifest is missing information needed for a safe change, update the manifest or document the gap.

Use `ADR-0012` for the decision authority boundary. After maintainer authorization, agents may professionally evaluate technical sub-decisions inside ratified directions, but must still ask before changing product direction, constitutional principles, runtime language, primary protocol direction, persistence direction, major architecture patterns, module ownership, breaking contracts, validation or permission strength, licensing-risk acceptance, hosting, cost, operations, or vendor-lock-in commitments.

## Verification Direction

These manifests should eventually power checks similar to:

```bash
vibit check architecture
vibit check module <module>
vibit check contracts
vibit check protocol
vibit check work
vibit check change <change-id>
```

Until those commands exist, agents must record architecture verification as not available.
