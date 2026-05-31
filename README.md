# vibit

vibit is an open-source **agent-native server framework** for building backends that AI coding agents and human developers can understand, extend, verify, and maintain from first principles.

Latest source alpha: `v0.1.0-alpha.1`

This is an early alpha for developers who want to inspect the architecture, run the first authenticated gameplay loop locally, and help shape an AI-maintainable backend framework. It is not a production server distribution yet.

## Why Try It

Most backend frameworks were designed before AI coding agents became part of the engineering loop. They may be powerful, but their architecture rules, ownership boundaries, change history, and verification paths are often scattered across code, docs, issues, and maintainer memory.

vibit is testing a stricter model:

- architecture rules live in machine-readable manifests;
- changes are bounded by work items, specs, ADRs, and verification records;
- public behavior is contract-first;
- generated output is traceable;
- module ownership is explicit;
- agents can inspect the next safe task with `tools/vibit`;
- humans can audit why a change exists without reading old chats.

The first domain focus is game/backend servers. The long-term target is an AI-era Nakama-class open-source backend framework, adapted around agent-native maintainability rather than direct API compatibility. Pitaya is deferred as a future architecture reference for distributed runtime concerns.

The product purpose is AI-native development and AI-native testing: a user states a backend requirement, and AI helps produce the bounded spec, acceptance criteria, tests, implementation, verification, and durable project records.

## Try The Alpha

Fastest source checkout path:

```bash
git clone https://github.com/iceiko/vibit.git
cd vibit
node tools/vibit check all
cd runtime && go test ./...
cd .. && examples/local-alpha-request-loop.sh
```

What this proves today:

- repository architecture checks run through `tools/vibit`;
- Go runtime tests pass;
- the local alpha request loop exercises the authenticated gameplay path;
- the script avoids printing raw credentials, raw access tokens, verifier keys, DSNs, digests, or transport metadata.

The packaged prototype-ready local path is:

- `docs/prototype-ready-local-development-path-package.md`
- `examples/README.md`
- `examples/local.prototype.env.example`

Prerequisites:

- Go;
- Node.js;
- PostgreSQL for the persistent runtime path;
- Buf and Protobuf tooling only when regenerating Protobuf output.

## What Exists

The repository has moved beyond a pure design phase. Current implemented foundation includes:

- Go runtime under `runtime/`;
- WebSocket gameplay endpoint at `/v1/ws`;
- Protobuf envelope and generated Go Protobuf output;
- PostgreSQL migration sources and platform persistence adapters;
- inventory proof slice;
- player account persistence boundaries and adapters;
- local onboarding/device credential issuance for development;
- device credential login service behavior and protocol route;
- opaque access-token validation for protected routes;
- runtime session persistence and response session metadata;
- first-message WebSocket connection binding;
- protected inventory and protected presence query path;
- logout service behavior and protocol route;
- single-process active connection lifecycle, close handoff, reconnect epoch handling, and presence lifecycle snapshot;
- health, readiness, version, and redacted config endpoints;
- `tools/vibit` checks and inspection commands for agents and humans.

The runtime listens on `:8080` by default and mounts:

```text
/v1/ws
/healthz
/readyz
/version
/configz
```

`/v1/ws` expects binary WebSocket messages containing `vibit.protocol.v1.Envelope` Protobuf bytes. JSON is not accepted on this gameplay endpoint.

`/healthz`, `/readyz`, `/version`, and `/configz` are small JSON troubleshooting endpoints. `/configz` reports only redacted runtime posture, not verifier keys, raw credentials, raw tokens, DSNs, digests, headers, cookies, query strings, subprotocol values, remote addresses, or concrete transport metadata.

## Who This Is For

Try vibit now if you:

- build or operate game/backend servers;
- have used or evaluated Nakama, Colyseus, Pomelo, Agones, Pitaya, or custom Go backends;
- want AI coding agents to make safer changes in a serious backend codebase;
- care about explicit architecture, contracts, generated structure, tests, and durable decision records;
- want to help define the first useful agent-native server framework shape before it hardens.

This alpha is not yet for production deployment, plug-and-play SDK use, hosted operations, or teams looking for packaged binaries and containers.

## Current Limits

`v0.1.0-alpha.1` is source-first:

- no release binaries;
- no packages;
- no container images;
- no checksum files;
- no provenance or signing artifacts;
- no hosted deployment;
- no install script;
- no SDK package;
- no direct Nakama/Pitaya API compatibility promise;
- no Pitaya-style cluster/RPC/frontend-backend work before a later ADR reactivates it.

The PostgreSQL runtime path is the most complete local path, but it still expects development setup: migrations, `VIBIT_POSTGRES_DSN`, and authentication verifier key environment variables. See `docs/runtime-runbook.md`.
For the current packaged local path, see `docs/prototype-ready-local-development-path-package.md`.

## Release Notes

Release notes for the source alpha live at:

- `docs/releases/v0.1.0-alpha.1.md`
- `docs/releases/v0.1.0-alpha.1.zh-CN.md`

The release authorization record is:

- `docs/release-execution-final-authorization.md`
- `decisions/ADR-0103-release-execution-final-authorization.md`

The first alpha user discovery loop is:

- `docs/first-alpha-user-discovery-loop.md`
- `decisions/ADR-0104-first-alpha-user-discovery-loop.md`

The first alpha feedback intake surface and product maturity milestones are:

- `.github/ISSUE_TEMPLATE/alpha-feedback.yml`
- `docs/first-alpha-feedback-intake-surfaces.md`
- `docs/product-maturity-milestones.md`
- `decisions/ADR-0105-first-alpha-feedback-intake-and-product-maturity-milestones.md`

The prototype-ready foundation execution plan is:

- `docs/prototype-ready-foundation-execution-plan.md`
- `decisions/ADR-0106-prototype-ready-foundation-execution-plan.md`

The prototype-ready local development path gate is:

- `docs/prototype-ready-local-development-path-gate.md`
- `decisions/ADR-0107-prototype-ready-local-development-path-gate.md`

The prototype-ready local development path package is:

- `docs/prototype-ready-local-development-path-package.md`
- `decisions/ADR-0108-prototype-ready-local-development-path-package.md`

The storage objects behavior gate is:

- `docs/storage-objects-behavior-gate.md`
- `decisions/ADR-0109-storage-objects-behavior-gate.md`

The storage objects persistence schema gate is:

- `docs/storage-objects-persistence-schema-gate.md`
- `decisions/ADR-0110-storage-objects-persistence-schema-gate.md`

The storage objects migration source is:

- `runtime/migrations/postgres/000006_create_storage_objects.sql`
- `decisions/ADR-0111-storage-objects-migration-source.md`

The storage objects repository boundary is:

- `docs/storage-objects-repository-boundary.md`
- `decisions/ADR-0112-storage-objects-repository-boundary.md`

## Continue Development

If you are an agent or contributor continuing the project, start with:

```bash
node tools/vibit inspect next
node tools/vibit check work --json
```

The current next work item is:

```text
W-0250 Define Pitaya-aligned server-to-server RPC boundary gate
```

Use `.arch/work-items.yaml` as the source of truth for continuation. `W-0231` through `W-0242` completed the friends relationship lifecycle, persistence, repository, PostgreSQL adapter, runtime behavior, protocol route, and local proof path. `W-0243` through `W-0249` completed the post-friends route capability selection, minimum operations inspection gate, source-first operations inspection, Pitaya-aligned distributed runtime vocabulary gate, `node tools/vibit inspect pitaya-vocabulary --json`, the frontend/backend role boundary gate, and `node tools/vibit inspect pitaya-roles --json`. `ADR-0157` registered `runtime.pitaya_aligned_frontend_backend_role_source_first_map` and opened `W-0250 Define Pitaya-aligned server-to-server RPC boundary gate`. The next direction is `define_pitaya_aligned_server_to_server_rpc_boundary_gate`. A request such as `continue` or `继续推进` means: advance one `next_ready` work item unless blocked by an ask-first boundary, verification failure, or required maintainer confirmation.

Recent friends relationship trace decisions remain recorded in `decisions/ADR-0141-friends-relationship-migration-source.md`, `decisions/ADR-0142-friends-relationship-repository-boundary.md`, `decisions/ADR-0143-friends-relationship-repository-interface-implementation.md`, `decisions/ADR-0144-friends-relationship-postgresql-adapter-gate.md`, `decisions/ADR-0145-friends-relationship-postgresql-adapter-implementation.md`, `decisions/ADR-0146-friends-relationship-runtime-behavior-gate.md`, `decisions/ADR-0147-friends-relationship-runtime-behavior-implementation.md`, `decisions/ADR-0148-friends-relationship-protocol-route-gate.md`, `decisions/ADR-0149-friends-relationship-protocol-route-implementation.md`, and `decisions/ADR-0150-friends-relationship-protocol-route-local-proof.md`.

## Agent-Native Means

Agent-native does not primarily mean that the server has AI features.

It means the codebase is designed so AI coding agents can work inside it reliably and can turn user requirements into tested backend changes:

- architecture rules are explicit;
- module ownership is declared;
- public behavior is contract-first;
- generated structure is traceable;
- user requirements become bounded specs;
- acceptance criteria and test plans are written before non-trivial implementation;
- business rules are tested;
- change workflow is bounded;
- repository state is checkable;
- project memory is stored in durable artifacts instead of chat history.

AI gameplay features such as NPC agents, memory, model routing, tool calling, and simulations may become extensions later. They are not the foundation.

## Nakama-First Target

Nakama is the active primary product reference for capability planning.

- Nakama guides the broad game backend surface: accounts, sessions, storage, social systems, chat, groups, parties, leaderboards, tournaments, matchmaking, realtime multiplayer, authoritative matches, operations, SDKs, and examples.
- Pitaya is deferred as a future architecture reference for acceptors, sessions, routes, handlers, remotes/RPC, frontend/backend roles, groups, broadcast, serializers, and cluster service discovery.

vibit should cover the same class of common Nakama-style capability over time while preserving its own architecture: explicit manifests, contracts, generation, tests, ADRs, repository checks, bounded agent workflow, and AI-native requirement-to-test-to-implementation flow.

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
- `.arch/reference.yaml`: Nakama-first reference and product parity planning.
- `docs/v0.1-alpha-goal.md`: short-term alpha and long-term product target.
- `docs/alpha-developer-flow.md`: packaged local alpha developer journey.
- `docs/alpha-acceptance-checklist.md`: local v0.1 alpha acceptance checklist.
- `docs/runtime-runbook.md`: current runtime startup and verification notes.
- `docs/release-publishing-decision-gate.md`: release publishing decision boundary.
- `docs/release-execution-preparation-gate.md`: release execution preparation boundary.
- `docs/release-execution-authorization-gate.md`: release execution authorization criteria.
- `docs/release-execution-final-authorization.md`: final release authorization record.
- `docs/first-alpha-user-discovery-loop.md`: first alpha user discovery loop.
- `docs/first-alpha-feedback-intake-surfaces.md`: first alpha feedback intake surface.
- `docs/product-maturity-milestones.md`: source alpha, prototype-ready, production-candidate, and product-class milestones.
- `docs/prototype-ready-local-development-path-package.md`: repeatable source-first local development path.
- `docs/storage-objects-behavior-gate.md`: first player-owned small storage objects behavior gate.
- `docs/releases/v0.1.0-alpha.1.md`: alpha release notes.
- `docs/nakama-pitaya-product-parity-roadmap.md`: long-term Nakama-first capability roadmap.
- `examples/README.md`: local examples and redacted configuration template guidance.
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
