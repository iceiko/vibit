# Agent Operating Guide

Status: Draft v0.1  
Last updated: 2026-05-13
Scope: Repository-level operating instructions for coding agents  
Canonical source: `CONSTITUTION.md`

This guide turns the constitution into working rules for agents. It does not replace the constitution. When this guide and `CONSTITUTION.md` conflict, follow `CONSTITUTION.md` and update this guide.

The paired Simplified Chinese translation is `AGENTS.zh-CN.md`. The English file is authoritative.

## 1. Project Identity

Working name:

```text
vibit
```

Category:

```text
Agent-Native Server Framework
```

Positioning:

```text
vibit is an open-source agent-native server framework for building backends that AI coding agents can understand, extend, verify, and maintain from first principles.
```

In this repository, "AI-native" means agent-native maintainability first. It does not primarily mean adding AI gameplay features or AI product features.

## 2. Required Reading

Before making a non-trivial change, read:

- `CONSTITUTION.md`
- This file
- The relevant architecture manifests under `.arch/`, when they exist
- The relevant module manifest at `modules/<module>/module.yaml`, when it exists
- The relevant module guide at `modules/<module>/AGENTS.md`, when it exists
- The relevant change spec under `changes/`, when the change has one

If an expected artifact does not exist yet, do not invent hidden assumptions. Either create the missing artifact as part of the change or record that it is not yet available.

## 3. Current Repository State

This repository is currently in the constitutional and standards-design phase.

Existing foundation:

- `CONSTITUTION.md`
- `CONSTITUTION.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `.arch/README.md`
- `.arch/modules.yaml`
- `.arch/conventions.yaml`
- `.arch/protocol.yaml`
- `.arch/runtime.yaml`
- `.arch/contracts.yaml`
- `.arch/dependencies.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `docs/module-manifest.md`
- `docs/module-manifest.zh-CN.md`
- `docs/change-spec.md`
- `docs/change-spec.zh-CN.md`
- `changes/_template/`
- `docs/conversation-log.md`
- `docs/conversation-log.zh-CN.md`
- `conversations/`
- `docs/agent-decision-record.md`
- `docs/agent-decision-record.zh-CN.md`
- `decisions/`
- `docs/schema-validation.md`
- `docs/schema-validation.zh-CN.md`
- `docs/dependency-adoption.md`
- `docs/dependency-adoption.zh-CN.md`
- `docs/game-protocol.md`
- `docs/game-protocol.zh-CN.md`
- `docs/generated-output.md`
- `docs/generated-output.zh-CN.md`
- `docs/runtime-protocol-adapter.md`
- `docs/runtime-protocol-adapter.zh-CN.md`
- `docs/reference-game-server-alignment.md`
- `docs/reference-game-server-alignment.zh-CN.md`
- `docs/workflow.md`
- `docs/workflow.zh-CN.md`
- `schema/`
- `rules/`

Framework implementation code, generators, modules, and verification commands may not exist yet. When they do not exist, document that verification is not available instead of pretending that it ran.

Runtime readiness decisions currently point to Go as the first server runtime implementation language, WebSocket as the first gameplay/client protocol, Protobuf as the first wire message format, a modular monolith single-process server model, contract-first commands/queries/events/errors/permissions, and `inventory` as the preferred first proof slice. Read `.arch/runtime.yaml`, `ADR-0004` through `ADR-0010`, and note that `ADR-0003` is superseded before creating runtime implementation code.

The first game protocol framework is recorded in `.arch/protocol.yaml`, `docs/game-protocol.md`, `ADR-0015`, and `ADR-0016`. It defines a WebSocket-framed Protobuf envelope with explicit `kind`, `module`, and `name` routing fields, session metadata, game target scopes, server-authoritative message rules, error mapping, and compatibility expectations. Read it before adding `.proto` files, WebSocket protocol handlers, generated protocol output, or client/server protocol rules.

The first protocol source files are `proto/vibit/protocol/v1/envelope.proto` and `proto/vibit/inventory/v1/inventory.proto`. Buf configuration lives at `buf.yaml` and `buf.gen.yaml`. `ADR-0016` records the envelope and generation configuration decision. Generated Go Protobuf output remains planned under `runtime/internal/generated/proto/`; do not create or edit generated Go Protobuf files by hand.

The generated output standard is `docs/generated-output.md`, with `docs/generated-output.zh-CN.md` as the paired Simplified Chinese translation. `ADR-0017` records the generated output decision. Read these before adding generated files, generated output checks, or generator behavior. Go Protobuf output under `runtime/internal/generated/proto/` must be `*.pb.go`, must contain the `protoc-gen-go` generated-code marker, and must trace to an existing `.proto` source.

The runtime protocol adapter boundary standard is `docs/runtime-protocol-adapter.md`, with `docs/runtime-protocol-adapter.zh-CN.md` as the paired Simplified Chinese translation. `ADR-0018` records the boundary decision. Read these before adding WebSocket transport code, Protobuf runtime adapter code, application dispatch code, or domain runtime handlers.

The active game server reference alignment standard is `docs/reference-game-server-alignment.md`, with `docs/reference-game-server-alignment.zh-CN.md` as the paired Simplified Chinese translation. `ADR-0019` records Nakama and Pitaya as active reference baselines. Read `.arch/reference.yaml` and the standard before adding new game server capability families, runtime subsystems, social/realtime features, matchmaking, match runtime, cluster/RPC work, or operational surfaces. Nakama and Pitaya guide capability planning; they do not override vibit's constitution, ADRs, manifests, generated boundaries, or verification commands.

The work continuation standard is `docs/workflow.md`, with `docs/workflow.zh-CN.md` as the paired Simplified Chinese translation. The machine-readable work queue is `.arch/work-items.yaml`. When the maintainer says "continue" or "继续", interpret that as advancing one `next_ready` work item unless blocked or confirmation is required. When the maintainer asks to continue multiple steps, advance up to that many work items in order, stopping at blockers, verification failures, ask-first boundaries, or maintainer redirection.

Current executable tooling:

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
node tools/vibit check protocol
node tools/vibit check protocol --json
node tools/vibit check generated
node tools/vibit check generated --json
node tools/vibit check runtime
node tools/vibit check runtime --json
node tools/vibit check work
node tools/vibit check work --json
node tools/vibit inspect module <module>
node tools/vibit inspect boundary --from <module> --to <module>
node tools/vibit inspect contract --module <module> --type <type> --id <id>
node tools/vibit inspect change <change-id>
node tools/vibit inspect work
node tools/vibit inspect work --json
node tools/vibit inspect memory
node tools/vibit inspect rule <rule-id>
node tools/vibit inspect rules
node tools/vibit inspect rules --category <category>
node tools/vibit check architecture
node tools/vibit check architecture --json
node tools/vibit check change <change-id>
node tools/vibit check change <change-id> --json
node tools/vibit check module <module>
node tools/vibit check module <module> --json
node tools/vibit generate module <module>
```

Use `node tools/vibit check all` as the default repository verification command when CLI tooling is available.

The current CLI is Node.js standard-library tooling only. Do not treat CLI implementation language as the server runtime language.

Use `--json` when an agent needs machine-readable check results during intake, verification, or handoff.

Every JSON check result item should include `rule_id` and `artifact`. Treat `check all --json` as a compact overview, then run the specific failing check with `--json` for full detail.

Use `node tools/vibit check memory` when conversation logs or Agent Decision Records are added or changed.

Use `node tools/vibit check contracts` when contract source files or `.arch/contracts.yaml` are added or changed.

Use `node tools/vibit check protocol` before creating or changing `.proto` files, generated Protobuf output, or protocol generation rules. Missing `.proto` files may pass while protocol sources are still planned, but once a `.proto` file exists it must align with registered command, query, and event contracts.

Use `node tools/vibit check generated` when generated files, module manifest `generated` declarations, generated output standards, or Go Protobuf generated output are added or changed.

Use `node tools/vibit check runtime` when runtime module behavior, runtime adapter boundaries, runtime guidance, or tests are added or changed. Before the Go runtime exists, this check should pass as not applicable because runtime implementation has not started. After `runtime/go.mod` exists but before Go source files exist, this check should verify the ADR-0014 skeleton and the ADR-0018 runtime protocol adapter boundary, and pass without running `go test`. Once Go source files exist, runtime checks require Go test files and a local Go toolchain.

Use `node tools/vibit inspect work` before interpreting a continuation request. Use `node tools/vibit check work` when `.arch/work-items.yaml`, workflow docs, or work item state changes. The default continuation unit is one work item.

Use `node tools/vibit inspect contract --module <module> --type <type> --id <id>` during intake when an agent needs one contract's registry entry, source summary, module manifest declaration, and consistency status as JSON.

Use `node tools/vibit inspect change <change-id>` during intake or handoff when a change spec exists and an agent needs a structured summary of its files, metadata, affected modules, and verification state.

Use `node tools/vibit inspect memory` when an agent needs a structured index of change specs, conversation logs, and Agent Decision Records before choosing which artifacts to read in full.

Use `rules/check-rules.json` to interpret check result `rule_id` values.

Use `node tools/vibit inspect rule <rule-id>` when only one rule's metadata is needed.

Use `node tools/vibit inspect rules --category <category>` to discover rules by category.

Use `.arch/runtime.yaml` as the machine-readable intake point for runtime readiness. It links the ADRs that govern language, server instance model, contract and generation boundary, client protocol, dependency adoption, and first proof slice.

Use `.arch/protocol.yaml` as the machine-readable intake point for game protocol framework decisions. It links `ADR-0015` and defines the first WebSocket Protobuf envelope, route fields, session model, target scopes, authority rules, error model, and first inventory slice protocol scope.

Use `ADR-0016`, `ADR-0017`, `buf.yaml`, `buf.gen.yaml`, `proto/README.md`, and `docs/generated-output.md` before changing the protocol envelope, inventory Protobuf source, Buf generation configuration, generated output checks, or generated Go Protobuf output path.

Use `ADR-0018` and `docs/runtime-protocol-adapter.md` before changing runtime code that sits between WebSocket transport, Protobuf protocol adaptation, application dispatch, generated code, and domain modules.

Use `.arch/dependencies.yaml` as the machine-readable intake point before adding foundational dependencies. Use `docs/dependency-adoption.md` and `docs/_templates/dependency-adoption.md` for adoption records.

Use `.arch/reference.yaml` as the machine-readable intake point for Nakama/Pitaya reference alignment. Nakama is the primary reference for broad game backend product capability surface. Pitaya is the primary reference for Go game server framework architecture vocabulary. Preserve vibit's Agent-Native constraints when adapting reference patterns, and record why a reference pattern is adopted, adapted, or rejected.

Use `.arch/work-items.yaml` as the machine-readable intake point for continuation. Work item IDs such as `W-0007` are execution steps; ADR IDs remain architectural decisions; change spec IDs remain concrete execution records; Git hashes remain repository snapshots; versions remain release identifiers.

Use `ADR-0014` before changing Go runtime files. The first Go module lives at `runtime/go.mod` with module path `github.com/iceiko/vibit/runtime`. Keep process startup under `runtime/cmd/vibit-server/`, application dispatch and composition under `runtime/internal/app/`, platform adapters under `runtime/internal/platform/`, handwritten domain module logic under `runtime/internal/modules/<module>/`, generated Go contract shapes under `runtime/internal/generated/contracts/`, generated Go Protobuf output under `runtime/internal/generated/proto/`, SQL-first PostgreSQL migrations under `runtime/migrations/postgres/`, and Protobuf source files under repository-root `proto/vibit/<module>/v1/`.

## 4. Documentation Rules

English is the canonical documentation language.

Every public-facing document should have:

- An English source document
- A Simplified Chinese human-readable translation

Naming examples:

```text
CONSTITUTION.md
CONSTITUTION.zh-CN.md
AGENTS.md
AGENTS.zh-CN.md
docs/<name>.md
docs/<name>.zh-CN.md
.arch/README.md
.arch/README.zh-CN.md
```

Rules:

- Update the Chinese translation in the same change when the English source changes materially.
- If the translation cannot be updated in the same change, mark it clearly as out of date.
- Keep machine-readable identifiers in English.
- Use English for code identifiers, module names, commands, events, permissions, and errors unless a strong domain reason exists.
- Preserve meaning in translation. Do not force literal word-by-word translation when it reduces clarity.

## 5. Standard Change Workflow

For every non-trivial feature, bug fix, migration, refactor, or standard change:

1. Clarify the requirement.
2. Identify affected modules and contracts.
3. Write or update the change spec when the change is large enough to need durable context.
4. Update schemas, manifests, or contracts before implementation when public behavior changes.
5. Generate repeatable structure when generators exist.
6. Implement only inside the declared boundary.
7. Add or update focused tests.
8. Run relevant verification commands.
9. Update documentation and translations.
10. Record what was verified and what was not verified.

For early design-only changes, steps involving code, tests, generators, and verification may be not applicable. Say that explicitly.

## 6. Architecture Rules

Prefer designs that:

- Give agents less ambiguous context
- Create stronger module boundaries
- Make behavior easier to verify
- Make contracts explicit
- Reduce hidden coupling
- Support code generation
- Remain practical for human developers

Do not rely on maintainer memory for architecture rules. If a rule matters, it should eventually be represented in a document, schema, manifest, test, generator, or architecture check.

## 7. Module Rules

When modules exist, each module should declare:

- What it owns
- What it does not own
- Public commands
- Public queries
- Published events
- Subscribed events
- Allowed dependencies
- Forbidden dependencies
- Invariants
- Required tests
- Generated files
- Handwritten extension points

Other modules must not reach into a module's internals directly. Cross-module communication should happen through commands, queries, events, public module APIs, or generated clients.

Use `docs/module-manifest.md` as the source standard for `modules/<module>/module.yaml`.

Use `docs/change-spec.md` as the source standard for `changes/<date>-<change-id>/`.

Use `docs/conversation-log.md` as the source standard for `conversations/`.

When the maintainer introduces product intent, rejects an interpretation, names a concept, or makes an architectural decision, preserve that context in a conversation log. Redact secrets before committing.

Use `docs/agent-decision-record.md` as the source standard for `decisions/`.

When a decision affects long-term architecture, generated file conventions, module ownership, or a rejected plausible alternative, create or update an Agent Decision Record. Keep rationale concise and public; do not store hidden chain-of-thought.

Generated files are immutable to non-system agents. If generated output is wrong, change the source schema, template, or generator unless a change spec or decision record explicitly grants a `generated_file_override`.

For the server runtime, Go is the first implementation language. WebSocket is the first gameplay/client protocol. Protobuf is the first wire message format. PostgreSQL is the first authoritative durable relational store. S3-compatible object storage is a planned object-storage abstraction, with MinIO as the preferred local/self-hosted candidate pending a dependency adoption record. Domain modules must not depend directly on third-party transport, protocol, persistence, object-storage, or framework libraries; platform adapters own those dependencies behind vibit-owned interfaces.

Accepted first Go runtime dependencies are recorded in `ADR-0013` and `.arch/dependencies.yaml`:

- `github.com/coder/websocket` for the platform WebSocket transport adapter.
- `google.golang.org/protobuf` and `google.golang.org/protobuf/cmd/protoc-gen-go` for Go Protobuf runtime and generation.
- Buf CLI for Protobuf linting, breaking checks, formatting, and generation orchestration.
- `github.com/jackc/pgx/v5` for PostgreSQL platform persistence adapters.
- `github.com/pressly/goose/v3` for SQL-first migration tooling.
- Go standard-library `testing` first; no external test framework is adopted yet.

Direct imports or invocations of accepted dependencies are allowed only in their declared owner layers. Domain runtime logic and domain modules must use vibit-owned interfaces, generated contracts, repositories, and adapters.

Goose migrations should be SQL-first. Go migrations require a change spec explaining why SQL is insufficient and must not hide domain business logic.

When adding Go runtime code, follow the ADR-0014 package boundary:

- `runtime/cmd/vibit-server/` owns startup, configuration wiring, and process lifecycle.
- `runtime/internal/app/` owns command/query dispatch, application service composition, and transaction orchestration.
- `runtime/internal/platform/transport/ws/` owns `github.com/coder/websocket`.
- `runtime/internal/platform/protocol/protobuf/` owns Protobuf framing and envelope conversion.
- `runtime/internal/platform/persistence/postgres/` owns `github.com/jackc/pgx/v5`.
- `runtime/internal/platform/migrations/` owns `github.com/pressly/goose/v3` invocation and migration validation.
- `runtime/internal/platform/events/` owns event recording and publication mechanisms.
- `runtime/internal/platform/tx/` owns unit-of-work and transaction boundary interfaces.
- `runtime/internal/modules/<module>/` owns handwritten domain behavior only.

Runtime protocol handoff must follow `docs/runtime-protocol-adapter.md`: WebSocket transport reads and writes frames, the Protobuf adapter converts envelopes and payloads, application dispatch routes commands and queries, domain modules enforce invariants, and generated packages provide shapes only.

State-changing commands should enter through application dispatch and run inside an application-owned unit of work. Domain events produced by a command should be recorded in the same unit of work. Query handlers should not mutate state and do not require a write transaction by default. Event publication outside the transaction remains deferred until an explicit event delivery or outbox decision exists.

Before adding persistence implementation, agents must declare or update the relevant repository interfaces, migration expectations, transaction boundaries, and storage verification path. Do not add PostgreSQL drivers, migration tools, S3 SDKs, or MinIO clients without a change spec or adoption record that follows `ADR-0010` and `ADR-0011`.

Before adding foundational dependencies, agents must check `.arch/dependencies.yaml`. Do not change a dependency slot to `accepted` until an adoption record documents the problem solved, license, maintenance activity, abstraction boundary, allowed owners, forbidden owners, replacement path, and verification path.

Before adding new game server capability families or runtime subsystems, agents must check `.arch/reference.yaml` and `docs/reference-game-server-alignment.md`. Map the proposal to the relevant Nakama/Pitaya capability family, then keep the implementation sequence aligned with vibit's contract-first, manifest-first, generated, and checkable architecture. Do not copy external APIs without an explicit compatibility ADR. Do not add Pitaya-style cluster/RPC/service-discovery work before the modular monolith proof slice is stable.

Use `ADR-0012` for decision authority. After explicit maintainer authorization, agents may professionally evaluate and decide technical sub-decisions inside an already ratified direction. Still ask the maintainer before changing constitutional principles, product direction, runtime language, primary protocol direction, persistence direction, major architecture patterns, module ownership, breaking contracts, validation or permission strength, licensing-risk acceptance, hosting, cost, operations, or vendor-lock-in commitments.

Use `docs/schema-validation.md` as the source standard for `schema/`.

When changing the shape of module manifests, change specs, Agent Decision Records, or tool JSON output, update the paired schema file and run `node tools/vibit check schemas`.

## 8. Contract Rules

Public behavior should be specified before implementation.

Contract-bearing artifacts may include:

- API schemas
- Command schemas
- Query schemas
- Event schemas
- Error catalogs
- Permission catalogs
- Database migration schemas
- Generated clients

Rules:

- Public contracts must be declared before use.
- Compatibility-sensitive contracts must be versioned.
- Breaking changes must be explicit.
- Generated output must be traceable to source schema.
- Do not hand-edit generated contract output unless the generator itself is being changed.

## 9. Testing And Verification

Testing is part of the architecture, not a finishing step.

When implementation code exists, relevant verification may include:

- Unit tests
- Contract tests
- Invariant tests
- Integration tests
- Migration tests
- Replay tests
- Architecture checks
- Generator checks
- Documentation consistency checks

This repository does not yet define final verification commands. Until it does, record verification as one of:

```text
Verified: <commands or checks run>
Not verified: <reason>
Not applicable: <reason>
```

Never claim that a change is verified when verification was not run.

## 10. Ask First

Ask the human maintainer before:

- Changing constitutional principles
- Ratifying or replacing the project name
- Redefining module ownership
- Introducing a new architectural pattern
- Making breaking API, command, query, or event changes
- Changing generated file conventions
- Removing tests
- Weakening validation or permission checks
- Moving data ownership between modules
- Accepting meaningful licensing-risk, hosting, cost, operations, or vendor-lock-in commitments
- Changing server runtime language, primary protocol direction, persistence direction, or core project thesis
- Adding a major external framework dependency

## 11. Never

Never:

- Treat AI gameplay features as the foundation of this project
- Bypass module boundaries for convenience
- Hide business logic in transport handlers
- Add unregistered public events
- Add unregistered permissions
- Add untyped cross-module payloads
- Make broad repository edits without a declared boundary
- Hand-edit generated files without documenting why
- Leave an English public document materially changed while its Chinese translation silently falls behind
- Claim verification was run when it was not

## 12. When Adding New Standards

New standards should explain:

- The problem being solved
- The rule being introduced
- The reason the rule helps agents
- The impact on humans
- The expected artifacts
- The verification path
- The migration path from existing work

Prefer a small standard that can be enforced over a broad statement that cannot be checked.

## 13. When Adding Implementation Code

Do not start by scattering framework code across the repository.

Start from the smallest complete slice that proves the core claim:

```text
requirement -> spec -> contract -> generated shape -> handwritten logic -> tests -> verification -> docs
```

A good first implementation target should include a small but complete backend domain, such as player accounts, inventory, currency, rewards, quests, or match sessions.

## 14. Bootstrapping Control

Self-bootstrapping is useful only while it improves the path to a working server framework.

Before adding a new standard, inspect command, check command, schema, generator, or workflow rule, confirm that it directly supports at least one of:

- The next runtime vertical slice
- A concrete module boundary
- A public contract or generated shape
- A test or verification path
- Agent context reduction for an expected implementation task

If the benefit is mainly that the tooling becomes more complete, defer it.

When the repository already has enough tooling to attempt a small end-to-end backend capability, prefer runtime readiness work over additional meta-tooling, then implement the runtime slice.

Runtime readiness should answer only the decisions needed to make the first slice coherent:

- Implementation language and package layout
- Minimal server instance model
- First module and capability boundary
- Contract format
- Generated versus handwritten file boundary
- Minimum test and verification strategy
- Persistence and migration assumptions

Do not rush into implementation when these choices are still ambiguous. Also do not extend readiness work after it stops changing how the first slice will be built, verified, or maintained.

Record exceptions in a change spec or Agent Decision Record.

## 15. Handoff Requirements

At the end of a change, leave enough context for the next agent or human to continue.

Record:

- What changed
- Why it changed
- Which files changed
- Which contracts or standards changed
- What was verified
- What was not verified
- Which open questions remain

If the work is incomplete, state the next concrete action.
