# ADR-0014: Go Runtime Layout And Boundaries

Status: Accepted
Date: 2026-05-13
Decision Makers: Agent
Related changes:

- `changes/2026-05-13-define-go-runtime-layout/`

Related conversations:

- `conversations/2026-05-13-go-runtime-layout-and-boundaries.md`

Related artifacts:

- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `decisions/ADR-0004-minimal-server-instance-model.md`
- `decisions/ADR-0005-contract-and-generation-boundary.md`
- `decisions/ADR-0008-go-server-runtime-language.md`
- `decisions/ADR-0009-websocket-protobuf-client-protocol.md`
- `decisions/ADR-0011-postgresql-and-object-storage-persistence.md`
- `decisions/ADR-0013-first-go-runtime-dependencies.md`

## Context

The first runtime direction is Go, WebSocket, Protobuf, PostgreSQL, and a modular monolith single-process server. The first accepted dependencies are already recorded in ADR-0013.

Before adding `go.mod` or implementation code, vibit needs a concrete package layout that agents can follow without inferring architecture from partial files.

The layout must preserve these constraints:

- Domain modules must not import transport, protocol, persistence, migration, object-storage, or framework dependencies directly.
- Transport handlers must not hide business logic.
- Generated files must be traceable and immutable to non-system agents.
- Public contracts must precede handwritten behavior.
- Persistence must use vibit-owned repository and transaction interfaces.

## Decision

Use `runtime/` as the first Go module root.

The first runtime layout is:

```text
runtime/
  go.mod
  cmd/vibit-server/
  internal/app/
  internal/platform/transport/ws/
  internal/platform/protocol/protobuf/
  internal/platform/persistence/postgres/
  internal/platform/migrations/
  internal/platform/events/
  internal/platform/tx/
  internal/modules/<module>/
  internal/generated/contracts/
  internal/generated/proto/
  migrations/postgres/
proto/
  vibit/<module>/v1/
```

Package ownership:

- `runtime/cmd/vibit-server/` owns process startup, configuration wiring, and lifecycle.
- `runtime/internal/app/` owns command/query dispatch, application service composition, and transaction orchestration.
- `runtime/internal/platform/transport/ws/` owns `github.com/coder/websocket` and WebSocket connection adaptation.
- `runtime/internal/platform/protocol/protobuf/` owns Protobuf framing, envelope encode/decode, and conversion to generated message types.
- `runtime/internal/platform/persistence/postgres/` owns `github.com/jackc/pgx/v5` and PostgreSQL adapter implementation.
- `runtime/internal/platform/migrations/` owns `github.com/pressly/goose/v3` invocation and migration validation.
- `runtime/internal/platform/events/` owns event recording and publication mechanisms.
- `runtime/internal/platform/tx/` owns transaction boundary interfaces and unit-of-work execution.
- `runtime/internal/modules/<module>/` owns handwritten domain behavior for that module.
- `runtime/internal/generated/contracts/` owns generated Go contract shapes derived from vibit semantic contracts.
- `runtime/internal/generated/proto/` owns generated Go Protobuf files derived from `proto/`.
- `runtime/migrations/postgres/` owns SQL-first PostgreSQL migration source files.
- `proto/` owns Protobuf source files and should be paired with root `buf.yaml` and `buf.gen.yaml` when Protobuf generation starts.

The first Go module path should be:

```text
github.com/iceiko/vibit/runtime
```

This keeps the runtime implementation explicitly separate from repository-level architecture tooling.

## Alternatives Considered

- Place `go.mod` at repository root.
- Use separate Go modules for runtime, generated code, modules, and tools.
- Put generated code under repository-root `generated/`.
- Put `.proto` files under `runtime/proto/`.
- Put domain modules under repository-root `modules/<module>/runtime/`.
- Put platform adapters directly under each module.
- Delay package layout until implementation starts.

## Rationale

Putting `go.mod` under `runtime/` keeps the Go server implementation separate from the existing Node.js repository tooling and architecture manifests. It also avoids implying that every repository-level artifact belongs to the Go module.

Using Go `internal` packages gives the project a native boundary tool. It prevents external imports from depending on internal runtime implementation and gives future architecture checks clear paths to inspect.

Keeping `.proto` source files at repository-root `proto/` makes wire schema a public source artifact alongside `contracts/`, while generated Go Protobuf code stays inside `runtime/internal/generated/proto/`.

Keeping platform adapters under `runtime/internal/platform/` prevents third-party dependencies from leaking into domain modules.

Keeping domain modules under `runtime/internal/modules/<module>/` makes the first runtime slice concrete while preserving the source-of-truth module manifests under repository-root `modules/<module>/module.yaml`.

## Agent Reasoning Summary

The layout should make the allowed edit path obvious: contracts and manifests are source, generated packages are regenerated, platform adapters contain third-party dependencies, and domain modules contain handwritten behavior behind vibit-owned interfaces.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  human_ergonomics: high
  implementation_cost: medium
  reversibility: medium
  long_term_maintainability: high
confidence: high
```

## Consequences

- The first `go.mod` should be created at `runtime/go.mod`.
- The first server command should be created under `runtime/cmd/vibit-server/`.
- Domain runtime code should live under `runtime/internal/modules/inventory/` for the first proof slice.
- WebSocket, Protobuf runtime, pgx, and goose imports are forbidden outside their declared platform or generated owners.
- `.proto` source files should be created under `proto/vibit/<module>/v1/`.
- Generated Go Protobuf output should go under `runtime/internal/generated/proto/`.
- SQL migrations should go under `runtime/migrations/postgres/`.
- Initial runtime tests should use Go standard-library `testing`.

## Transaction Boundary

State-changing commands must enter through application dispatch and run inside a transaction boundary owned by `runtime/internal/platform/tx/`.

The first transaction model:

- Command handler receives a vibit-owned context and module repository interfaces.
- Application layer opens a unit of work.
- Repository mutations happen inside that unit of work.
- Domain events produced by the command are recorded inside the same unit of work.
- Event publication outside the transaction is deferred until an explicit event delivery or outbox decision is made.
- Query handlers should not mutate state and do not require a write transaction by default.

Outbox storage is deferred until the first cross-boundary durable event delivery requirement appears.

## Reversal Conditions

Revisit this decision if:

- The first Go module under `runtime/` creates excessive import friction.
- Generated code needs to be shared by non-runtime Go modules sooner than expected.
- The modular monolith shape requires a clearer public package boundary than `internal` allows.
- Protobuf source location under root `proto/` conflicts with Buf or release packaging requirements.
- Transaction orchestration needs a stronger outbox or event-log model immediately.

## Follow-Up

- Create `runtime/go.mod` and skeletal package directories according to this layout.
- Add import-boundary checks once Go files exist.
- Add `proto/`, `buf.yaml`, and `buf.gen.yaml` when Protobuf generation starts.
- Add SQL migration naming rules before the first migration file.
- Decide whether the first repository implementation starts with a fake, PostgreSQL, or both.
