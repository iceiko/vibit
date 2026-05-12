# vibit

vibit is an open-source agent-native server framework for building backends that AI coding agents can understand, extend, verify, and maintain from first principles.

Status: constitutional design phase

## What This Project Means By Agent-Native

Agent-native does not primarily mean that the server has AI features.

It means the server architecture is designed so AI coding agents can work inside it reliably:

- The architecture is explicit instead of tribal knowledge.
- Module ownership is declared instead of guessed.
- Public behavior is contract-first.
- Repeatable structure is generated.
- Business rules are tested as invariants.
- Cross-module communication is bounded.
- Change workflow is documented and verifiable.
- Documentation is written for both humans and agents.

AI gameplay features such as NPC agents, memory, model routing, tool calling, and simulations may become extensions later. They are not the foundation.

## Why This Exists

Many existing server codebases were built for human maintainers with local context, long memory, and implicit team conventions. AI coding agents can help in those codebases, but they often lose force when architecture rules are hidden, module boundaries are weak, tests are incomplete, or public contracts are unclear.

vibit starts from a different premise:

> The next generation of long-lived server software should be designed so agents can safely understand, modify, verify, and extend it.

The goal is not to make agents magically smarter. The goal is to make the codebase more legible, bounded, generated, contract-driven, and testable.

## Current Documents

- `CONSTITUTION.md`: canonical project constitution
- `CONSTITUTION.zh-CN.md`: Simplified Chinese translation
- `AGENTS.md`: repository-level operating guide for coding agents
- `AGENTS.zh-CN.md`: Simplified Chinese translation
- `.arch/README.md`: machine-readable architecture manifest entry point
- `.arch/modules.yaml`: first draft module registry manifest
- `.arch/conventions.yaml`: first draft repository convention manifest
- `.arch/runtime.yaml`: runtime readiness manifest for the first Go server runtime direction
- `.arch/contracts.yaml`: contract registry for public command, query, event, error, and permission source files
- `.arch/dependencies.yaml`: dependency adoption registry for foundational dependency decision slots
- `docs/module-manifest.md`: module manifest standard
- `docs/module-manifest.zh-CN.md`: Simplified Chinese translation
- `docs/change-spec.md`: change spec standard
- `docs/change-spec.zh-CN.md`: Simplified Chinese translation
- `changes/_template/`: reusable change spec template
- `docs/conversation-log.md`: conversation log standard
- `docs/conversation-log.zh-CN.md`: Simplified Chinese translation
- `conversations/`: maintainer-agent project memory
- `docs/agent-decision-record.md`: Agent Decision Record standard
- `docs/agent-decision-record.zh-CN.md`: Simplified Chinese translation
- `decisions/`: durable decision rationale
- `docs/schema-validation.md`: schema validation standard
- `docs/schema-validation.zh-CN.md`: Simplified Chinese translation
- `docs/dependency-adoption.md`: dependency adoption standard
- `docs/dependency-adoption.zh-CN.md`: Simplified Chinese translation
- `schema/`: JSON Schema files for machine-checkable standards
- `rules/`: rule catalogs for machine-readable check metadata

English documents are canonical. Simplified Chinese translations are maintained for human readers and early project discussion.

## Intended Direction

vibit should evolve toward:

- Architecture manifests under `.arch/`
- A first Go server runtime governed by `.arch/runtime.yaml` and Agent Decision Records
- WebSocket as the first gameplay/client protocol
- Protobuf as the first client/server wire message format
- PostgreSQL as the first authoritative durable relational store
- `github.com/coder/websocket` as the first WebSocket platform adapter dependency
- `google.golang.org/protobuf`, `protoc-gen-go`, and Buf CLI as the first Protobuf tooling stack
- `github.com/jackc/pgx/v5` as the first PostgreSQL driver behind platform persistence adapters
- `github.com/pressly/goose/v3` as the first SQL-first migration tooling
- A first Go module planned at `runtime/go.mod` with module path `github.com/iceiko/vibit/runtime`
- Go runtime package boundaries under `runtime/cmd/vibit-server/`, `runtime/internal/app/`, `runtime/internal/platform/`, `runtime/internal/modules/`, and `runtime/internal/generated/`
- Protobuf source files under `proto/vibit/<module>/v1/`, with generated Go Protobuf output under `runtime/internal/generated/proto/`
- SQL-first PostgreSQL migration source files under `runtime/migrations/postgres/`
- S3-compatible object storage as a planned large-object storage abstraction, with MinIO as the preferred local/self-hosted candidate pending dependency adoption
- Module manifests at `modules/<module>/module.yaml`, following `docs/module-manifest.md`
- Module-level agent guides at `modules/<module>/AGENTS.md`
- Contract-first commands, queries, events, errors, permissions, and migrations
- Contract source files under `contracts/`, registered by `.arch/contracts.yaml`
- Foundational dependency decisions registered by `.arch/dependencies.yaml`
- Generated scaffolds for repeatable framework structure
- Architecture checks that verify dependency, contract, event, and generated-file rules
- Change specs under `changes/<date>-<change-id>/`, following `docs/change-spec.md`
- Conversation logs under `conversations/`, following `docs/conversation-log.md`
- Agent Decision Records under `decisions/`, following `docs/agent-decision-record.md`
- Schema validation under `schema/`, following `docs/schema-validation.md`
- Rule catalogs under `rules/`, starting with `rules/check-rules.json`

The first serious prototype should prove one claim:

> Given a new backend requirement, an AI coding agent can identify the affected module, update the correct contracts, generate the correct structure, implement the behavior, add tests, run verification, and update documentation without damaging unrelated architecture.

## CLI Prototype

The first executable standard lives at:

```bash
tools/vibit
```

Initial commands:

```bash
node tools/vibit --help
node tools/vibit check all
node tools/vibit check all --json
node tools/vibit check schemas
node tools/vibit check schemas --json
node tools/vibit check memory
node tools/vibit check memory --json
node tools/vibit check contracts
node tools/vibit check contracts --json
node tools/vibit check generated
node tools/vibit check generated --json
node tools/vibit check runtime
node tools/vibit check runtime --json
node tools/vibit inspect module inventory
node tools/vibit inspect boundary --from inventory --to player
node tools/vibit inspect contract --module inventory --type command --id GrantItem
node tools/vibit inspect change bootstrap-vibit-cli
node tools/vibit inspect memory
node tools/vibit inspect rule check.subcheck
node tools/vibit inspect rules --category check
node tools/vibit check architecture
node tools/vibit check architecture --json
node tools/vibit check change bootstrap-vibit-cli
node tools/vibit check change bootstrap-vibit-cli --json
node tools/vibit check module inventory
node tools/vibit check module inventory --json
node tools/vibit generate module <module>
```

The CLI currently uses Node.js standard-library APIs only. It is a prototype for architecture checks, inspection, and generators. It is not the server runtime and does not determine the server runtime language.

Use `--json` when an agent needs machine-readable check results during intake, verification, or handoff. Human-readable text output remains the default.

Each JSON check result item includes a stable `rule_id` and an `artifact` value so agents can route failures without parsing prose. `check all --json` is a compact overview; run the specific failing check with `--json` to get full result details.

Use `node tools/vibit check memory` to verify required conversation log and Agent Decision Record structure.

Use `node tools/vibit check contracts` to verify that `.arch/contracts.yaml` and registered contract source files are consistent.

Use `node tools/vibit check generated` to verify that module-declared generated files exist and include generated, source, and generator trace markers.

Use `node tools/vibit check runtime` for server runtime verification. Before the Go runtime exists, this check reports that runtime implementation has not started. After Go runtime work begins, it should run the Go runtime test path.

Use `node tools/vibit inspect contract --module <module> --type <type> --id <id>` to inspect one registered command, query, event, error catalog, or permission catalog as JSON during agent intake.

Use `node tools/vibit inspect change <change-id>` to inspect a change spec directory and its verification metadata without manually opening every file.

Use `node tools/vibit inspect memory` to list change specs, conversation logs, and Agent Decision Records as a machine-readable project memory index.

Rule metadata for check output lives in `rules/check-rules.json`.

Use `node tools/vibit inspect rule <rule-id>` to inspect one rule without parsing the full catalog.

Use `node tools/vibit inspect rules` or `node tools/vibit inspect rules --category <category>` to discover available rules.

The first server runtime direction is Go, using a modular monolith single-process server model. WebSocket is the first gameplay/client protocol, and Protobuf is the first client/server wire format. Semantic business contracts remain in vibit manifests and contract source files; Protobuf owns wire schema shape. See `.arch/runtime.yaml`, `decisions/ADR-0008-go-server-runtime-language.md`, and `decisions/ADR-0009-websocket-protobuf-client-protocol.md`.

PostgreSQL is the first authoritative durable relational store for runtime state. S3-compatible object storage is planned for large artifacts such as replays, snapshots, exports, binary assets, and diagnostic archives. MinIO is the preferred local/self-hosted candidate for that S3-compatible role, but it is not a mandatory runtime dependency until a concrete use case and dependency adoption record justify it. Domain modules must use vibit-owned storage interfaces rather than depending directly on database drivers or object-storage clients. See `decisions/ADR-0011-postgresql-and-object-storage-persistence.md`.

The first accepted foundational runtime dependencies are recorded in `decisions/ADR-0013-first-go-runtime-dependencies.md` and `.arch/dependencies.yaml`. They are accepted only for platform adapters and generation tooling, not for direct use inside domain modules. S3 client tooling, MinIO deployment, observability, and external Go test framework adoption remain deferred until concrete runtime needs justify them.

The first Go runtime layout is recorded in `decisions/ADR-0014-go-runtime-layout-and-boundaries.md`. Runtime code has not started yet, but future Go files should follow these boundaries:

- `runtime/cmd/vibit-server/`: process startup, configuration wiring, and lifecycle.
- `runtime/internal/app/`: command/query dispatch, application composition, and transaction orchestration.
- `runtime/internal/platform/`: WebSocket, Protobuf, PostgreSQL, migration, event, and transaction platform adapters.
- `runtime/internal/modules/<module>/`: handwritten domain module logic.
- `runtime/internal/generated/`: generated Go contract and Protobuf output.

State-changing commands should run inside an application-owned unit of work before repository mutation and domain-event recording. Event publication outside the transaction is deferred until vibit adopts an explicit event delivery or outbox standard.

## Early Reference Domain

A small game backend is the recommended first demonstration domain because it naturally contains state, permissions, events, consistency rules, and long-lived modules.

Suggested modules:

- Player accounts
- Inventory
- Currency
- Rewards
- Tasks or quests
- Match sessions

The first backend slice should emphasize maintainability and agent workflow over feature count. It should still be treated as the beginning of a long-maintained system, not as disposable demo code.

## Governance

Project decisions are governed by `CONSTITUTION.md`.

Before changing constitutional principles, ratifying the name, introducing a major architectural pattern, or making a breaking standard change, read the constitution and record the motivation, alternatives, compatibility impact, and migration path.

## Name

`vibit` is the product name.

The intended category phrase is:

```text
agent-native server framework
```

Before final ratification, the name should be checked against major public registries and platforms for obvious conflicts.
