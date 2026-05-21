# vibit

vibit is an open-source agent-native server framework for building backends that AI coding agents can understand, extend, verify, and maintain from first principles.

Status: pre-alpha, building toward `v0.1 alpha`

The short-term target is a first developer-usable `v0.1 alpha`: a single-node, PostgreSQL-backed, WebSocket + Protobuf game backend runtime that real developers can run locally, inspect, and use as a contribution starting point.

The long-term target is an AI-era Nakama or Pitaya: the same class of serious game/backend server capability, adapted around vibit's agent-native maintainability model. This does not mean direct Nakama or Pitaya API compatibility.

## What Exists Today

The repository has moved beyond a pure design phase. Current implemented foundation includes:

- Agent-readable governance: `CONSTITUTION.md`, `AGENTS.md`, change specs, ADRs, conversation logs, and machine-readable architecture manifests.
- A Go runtime under `runtime/`.
- A WebSocket gameplay endpoint at `/v1/ws`.
- A Protobuf envelope and generated Go Protobuf output.
- PostgreSQL migration sources and platform persistence adapters.
- An inventory proof slice.
- Player account persistence boundaries and adapters.
- Device credential login service behavior and protocol route.
- Opaque access-token validation for protected routes.
- Logout service behavior and protocol route.
- Runtime session persistence and response session metadata.
- First-message WebSocket connection binding.
- Single-process active connection lifecycle, close handoff, reconnect epoch handling, and presence lifecycle snapshot.
- `tools/vibit` checks and inspection commands for agents and humans.

This is not yet a finished alpha. The authenticated gameplay end-to-end path is now proven in a focused Go test, the runbook and request-loop script exist, and the runtime exposes a small health/readiness/version/config surface. The most important remaining missing piece is an alpha acceptance checklist.

## Try It Locally

Prerequisites for the current development workflow:

- Go, for the runtime.
- Node.js, for `tools/vibit`.
- PostgreSQL, when using the persistent runtime path.
- Buf and Protobuf tooling, when regenerating Protobuf output.

Run the repository checks:

```bash
node tools/vibit check all
```

Run the Go runtime tests:

```bash
cd runtime
go test ./...
```

Start the bootstrap in-memory runtime:

```bash
cd runtime
go run ./cmd/vibit-server
```

The runtime listens on `:8080` by default and mounts:

```text
/v1/ws
/healthz
/readyz
/version
/configz
```

The endpoint expects binary WebSocket messages containing `vibit.protocol.v1.Envelope` Protobuf bytes. JSON is not accepted on this gameplay endpoint.

`/healthz`, `/readyz`, `/version`, and `/configz` are small JSON troubleshooting endpoints. `/configz` reports only redacted runtime posture, not verifier keys, raw credentials, raw tokens, DSNs, digests, headers, cookies, query strings, subprotocol values, remote addresses, or concrete transport metadata.

The PostgreSQL runtime path is more complete, but it requires migrations, `VIBIT_POSTGRES_DSN`, and authentication verifier key environment variables. See `docs/runtime-runbook.md` for the current operational notes. The runbook is part of the v0.1 alpha hardening path and should be treated as development documentation, not a polished release guide.

## Next Target: v0.1 Alpha

The durable target is recorded in `docs/v0.1-alpha-goal.md`.

`v0.1 alpha` should let a technically capable developer:

- clone the repo,
- prepare local config without committing secrets,
- apply or verify required PostgreSQL migrations,
- create or obtain a first device credential,
- login through the protocol route,
- receive an opaque access token and runtime session metadata,
- bind a WebSocket connection,
- call a protected inventory route,
- query presence,
- logout,
- run checks,
- and identify a concrete next contribution.

Recommended next sequence:

1. Complete `W-0178`: protected presence protocol query. Completed.
2. Define and add first local onboarding/device credential issuance. Completed.
3. Select and prove onboarding -> login -> bind connection -> protected inventory -> presence query -> logout end to end. Completed.
4. Refresh the runtime runbook around the actual alpha path. Completed.
5. Add a minimal example client or request-loop script. Completed.
6. Add health/readiness/version/config surfaces. Completed.
7. Add an alpha acceptance checklist or check.

Run the minimal local alpha request loop with:

```bash
examples/local-alpha-request-loop.sh
```

It wraps the focused authenticated gameplay E2E proof and does not print raw credentials, raw access tokens, verifier keys, DSNs, digests, or transport metadata.

## Continue Development

If you are an agent or contributor continuing the project, start with:

```bash
node tools/vibit inspect next
node tools/vibit check work --json
```

At the time of this README update, the next ready item is:

```text
W-0188 Add alpha acceptance checklist
```

Use `.arch/work-items.yaml` as the source of truth for continuation. A request such as `continue` or `继续推进` means: advance one `next_ready` work item unless blocked by an ask-first boundary or verification failure.

## Agent-Native Means

Agent-native does not primarily mean that the server has AI features.

It means the codebase is designed so AI coding agents can work inside it reliably:

- architecture rules are explicit,
- module ownership is declared,
- public behavior is contract-first,
- generated structure is traceable,
- business rules are tested,
- change workflow is bounded,
- repository state is checkable,
- and project memory is stored in durable artifacts instead of chat history.

AI gameplay features such as NPC agents, memory, model routing, tool calling, and simulations may become extensions later. They are not the foundation.

## Nakama And Pitaya Target

Nakama and Pitaya are active reference baselines for capability planning.

- Nakama guides the broad game backend surface: accounts, sessions, storage, social systems, chat, groups, parties, leaderboards, tournaments, matchmaking, realtime multiplayer, authoritative matches, operations, SDKs, and examples.
- Pitaya guides Go game server architecture vocabulary: acceptors, sessions, routes, handlers, remotes/RPC, frontend/backend roles, groups, broadcast, serializers, and cluster service discovery.

vibit should cover the same class of common capability over time while preserving its own architecture: explicit manifests, contracts, generation, tests, ADRs, repository checks, and bounded agent workflow.

See:

- `docs/reference-game-server-alignment.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `docs/v0.1-alpha-goal.md`

## Project Map

Important entry points:

- `CONSTITUTION.md`: canonical project constitution.
- `AGENTS.md`: repository-level operating guide for coding agents.
- `.arch/README.md`: architecture manifest entry point.
- `.arch/work-items.yaml`: continuation queue.
- `.arch/runtime.yaml`: runtime readiness and implementation state.
- `.arch/reference.yaml`: Nakama/Pitaya reference and product parity planning.
- `docs/v0.1-alpha-goal.md`: short-term alpha and long-term product target.
- `docs/runtime-runbook.md`: current runtime startup and verification notes.
- `docs/nakama-pitaya-product-parity-roadmap.md`: long-term capability roadmap.
- `changes/`: concrete change specs and verification records.
- `conversations/`: durable maintainer-agent project memory.
- `decisions/`: Agent Decision Records.
- `runtime/`: first Go reference runtime.
- `proto/`: Protobuf protocol source files.
- `tools/vibit`: architecture check, inspection, and generator CLI.

English documents are canonical. Simplified Chinese translations are maintained for human readers and early project discussion.

## CLI

The current CLI is:

```bash
node tools/vibit --help
node tools/vibit inspect next
node tools/vibit inspect work
node tools/vibit inspect reference
node tools/vibit check all
node tools/vibit check all --json
node tools/vibit check runtime
node tools/vibit check work
node tools/vibit check memory
node tools/vibit check schemas
```

Use `--json` when an agent needs machine-readable check results during intake, verification, or handoff.

## Governance

Project decisions are governed by `CONSTITUTION.md`.

Before changing constitutional principles, introducing a major architectural pattern, changing public protocol shape, adding dependencies, or shifting release direction, record the motivation, alternatives, compatibility impact, and migration path in a change spec and ADR when appropriate.

`vibit` is the product name. The intended category phrase is:

```text
agent-native server framework
```
