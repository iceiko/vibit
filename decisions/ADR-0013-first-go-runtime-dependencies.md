# ADR-0013: First Go Runtime Dependencies

Status: Accepted
Date: 2026-05-13
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-13-adopt-first-runtime-dependencies/`

Related conversations:

- `conversations/2026-05-13-technical-decision-authority-and-runtime-dependencies.md`

Related artifacts:

- `.arch/dependencies.yaml`
- `.arch/runtime.yaml`
- `docs/dependency-adoption.md`
- `decisions/ADR-0010-foundational-dependency-adoption.md`
- `decisions/ADR-0011-postgresql-and-object-storage-persistence.md`
- `decisions/ADR-0012-agent-technical-decision-authority.md`

## Context

The first server runtime direction is Go. The first gameplay/client protocol is WebSocket. The first wire message format is Protobuf. PostgreSQL is the first authoritative durable relational store.

Before implementation starts, vibit needs enough dependency decisions to keep the first runtime slice coherent. These dependencies must not leak into domain modules. Domain modules should depend on vibit-owned contracts, repositories, and service interfaces while platform adapters and generation tooling own third-party packages.

External facts checked on 2026-05-13:

- `github.com/coder/websocket`: ISC license, active repository, latest release `v1.8.14` published 2025-09-06, repository updated 2026-05-11, about 5.1k stars. The README describes it as a minimal idiomatic Go WebSocket library with context support, zero dependencies, Autobahn test compliance, concurrent writes, close handshake, ping API, and Wasm support.
- `github.com/gorilla/websocket`: BSD-2-Clause license, about 24.7k stars, mature and widely used, stable API, latest release `v1.5.3` published 2024-06-14.
- `google.golang.org/protobuf` / `github.com/protocolbuffers/protobuf-go`: BSD-3-Clause license, latest release `v1.36.11` published 2025-12-12, official Go support for Protobuf runtime and `protoc-gen-go`.
- `bufbuild/buf`: Apache-2.0 license, latest release `v1.69.0` published 2026-04-29, about 11.1k stars. The README documents linting, breaking-change detection, generation, formatting, and JSON-capable CLI output.
- `github.com/jackc/pgx/v5`: MIT license, active repository, about 13.8k stars. The README describes it as a pure Go PostgreSQL driver and toolkit with PostgreSQL-specific features, connection pool support, semantic versioning for stable public API, and `v5` as latest stable major version.
- `github.com/pressly/goose/v3`: MIT license, latest release `v3.27.1` published 2026-04-24, about 10.6k stars. The README describes it as a database migration CLI and library supporting incremental SQL changes and Go functions.
- `golang-migrate/migrate`: about 18.5k stars and widely used, but broader and less focused for vibit's SQL-first PostgreSQL path.
- `ariga/atlas`: Apache-2.0 license and strong schema-as-code tooling, but heavier and more opinionated than vibit needs before the first persistent slice.

Source links:

- https://github.com/coder/websocket
- https://github.com/gorilla/websocket
- https://github.com/protocolbuffers/protobuf-go
- https://github.com/bufbuild/buf
- https://github.com/jackc/pgx
- https://github.com/pressly/goose
- https://github.com/golang-migrate/migrate
- https://github.com/ariga/atlas

## Decision

Accept these first foundational dependencies:

```yaml
websocket_server: github.com/coder/websocket
protobuf_tooling:
  runtime: google.golang.org/protobuf
  generator: google.golang.org/protobuf/cmd/protoc-gen-go
  orchestrator: buf.build CLI
postgresql_driver: github.com/jackc/pgx/v5
migration_tooling: github.com/pressly/goose/v3
test_framework: Go standard library testing first; no external test framework accepted yet.
```

Keep these categories deferred:

```yaml
s3_client: deferred
minio_server: deferred
observability: deferred
```

`github.com/coder/websocket` is accepted over `github.com/gorilla/websocket` for the first platform adapter because its API is more context-oriented, minimal, dependency-light, and agent-readable. Gorilla remains a credible replacement candidate because it is more widely starred and mature.

`goose` is accepted as SQL-first migration tooling. vibit should prefer SQL migration files for schema changes. Go migrations are allowed only for platform migration tooling when a change spec explains why SQL is insufficient; they must not become a hidden place for domain business logic.

## Evaluation

### WebSocket Server

Accepted dependency: `github.com/coder/websocket`.

Evaluation:

- Maintenance activity: active repository with recent updates and a 2025 release.
- License compatibility: ISC, compatible with vibit's MIT project license.
- API stability: Go module with documented API and examples.
- Production adoption signals: lower star count than Gorilla, but credible adoption and active maintenance.
- Security and supply chain: zero runtime dependencies reduces transitive dependency risk.
- Operational fit: fits Go `context.Context`, close handling, pings, and WebSocket connection lifecycle.
- Agent readability: minimal API and context-oriented flow are easier for agents to follow than larger event-driven networking packages.
- Testability: platform transport adapter can be tested through adapter-level connection tests and protocol fixtures.
- Generated-code compatibility: no generated code dependency; it should only consume generated dispatch helpers through a vibit adapter.

### Protobuf Tooling

Accepted dependencies:

- `google.golang.org/protobuf`.
- `google.golang.org/protobuf/cmd/protoc-gen-go`.
- Buf CLI from `github.com/bufbuild/buf`.

Evaluation:

- Maintenance activity: official Go Protobuf runtime and generator remain active; Buf has active releases.
- License compatibility: Go Protobuf uses BSD-3-Clause; Buf uses Apache-2.0.
- API stability: official Protobuf Go module and generated code are expected to be stable, with documented compatibility caveats.
- Production adoption signals: Protobuf is a widely adopted wire format; Buf is widely used for Protobuf linting and breaking checks.
- Security and supply chain: generation should be pinned later through tool versions or reproducible install scripts.
- Operational fit: directly supports the ratified WebSocket + Protobuf client protocol.
- Agent readability: Buf lint, breaking checks, formatting, and machine-friendly output give agents clearer feedback than raw compiler failures alone.
- Testability: generated outputs can be checked for source traces and contract-to-proto alignment.
- Generated-code compatibility: this tooling directly shapes generated protocol packages and must be declared by generated files.

### PostgreSQL Driver

Accepted dependency: `github.com/jackc/pgx/v5`.

Evaluation:

- Maintenance activity: active repository with recent updates.
- License compatibility: MIT.
- API stability: `v5` is the current stable major version and follows semantic versioning for documented public APIs.
- Production adoption signals: mature Go PostgreSQL driver and toolkit with broad ecosystem usage.
- Security and supply chain: direct dependency in platform persistence adapters only; future version updates should be verified through integration tests.
- Operational fit: PostgreSQL-specific features match vibit's accepted PostgreSQL persistence direction.
- Agent readability: explicit APIs and pooling support are easier to reason about than hiding PostgreSQL behavior behind generic database portability too early.
- Testability: persistence adapters can be covered by repository tests and migration/integration tests.
- Generated-code compatibility: no direct generated output requirement, but generated repositories or query helpers must stay behind vibit-owned interfaces.

### Migration Tooling

Accepted dependency: `github.com/pressly/goose/v3`.

Evaluation:

- Maintenance activity: active project with recent release.
- License compatibility: MIT.
- API stability: mature CLI and Go module.
- Production adoption signals: widely used migration tool with SQL and Go migration support.
- Security and supply chain: migration execution is operationally sensitive, so future use must pin versions and run validation before applying migrations.
- Operational fit: SQL-first migrations fit PostgreSQL and keep schema changes inspectable for agents.
- Agent readability: versioned SQL files are simple, durable, and easy for agents to diff.
- Testability: migrations can be validated with goose commands and later exercised against disposable PostgreSQL instances.
- Generated-code compatibility: migration files are source artifacts unless a future generator explicitly owns them.

## Boundary

Allowed direct use:

- `github.com/coder/websocket`: platform transport adapter packages only.
- `google.golang.org/protobuf`: generated protocol packages, protocol serialization helpers, and generation tooling only.
- `protoc-gen-go`: generation tooling only.
- Buf CLI: repository tooling, generation, formatting, linting, and breaking checks only.
- `github.com/jackc/pgx/v5`: platform persistence adapter packages only.
- `github.com/pressly/goose/v3`: platform migration tooling only.

Forbidden direct use:

- Domain modules.
- Domain command handlers.
- Domain query handlers.
- Domain policies.
- Domain invariants.
- Transport handlers as a place to hide business logic.
- Protobuf generated files as handwritten extension points.
- Go migrations as a place to hide domain business logic.

Vibit-owned abstraction:

- Transport adapters convert WebSocket Protobuf frames into vibit commands and queries.
- Generated contract and protocol packages expose typed message shapes.
- Domain modules depend on vibit-owned command, query, event, repository, and transaction interfaces.
- Persistence adapters implement vibit-owned repository interfaces using `pgx/v5`.
- Migration tooling owns SQL migration apply, rollback, status, and validation commands.

## Replacement Path

Replacement path by dependency:

- `github.com/coder/websocket`: replace inside platform transport adapter with `github.com/gorilla/websocket` or another WebSocket implementation without changing domain modules.
- `google.golang.org/protobuf` and `protoc-gen-go`: replacing the official Go Protobuf stack would require a new protocol generation ADR and generated-code migration plan.
- Buf CLI: raw `protoc`, a custom vibit generator, or another Protobuf build tool could replace Buf if linting, breaking checks, formatting, and generation orchestration remain covered.
- `github.com/jackc/pgx/v5`: replace inside platform persistence adapters with `database/sql` plus a PostgreSQL driver if portability becomes more valuable than PostgreSQL-specific behavior.
- `github.com/pressly/goose/v3`: replace with `golang-migrate/migrate`, Atlas, or a custom vibit migration runner if schema drift, declarative schema state, or stronger verification becomes necessary.

## Verification

Current verification for this decision is documentation and manifest based:

```bash
node tools/vibit check architecture
node tools/vibit check memory
node tools/vibit check all
```

Future verification once Go runtime code exists:

```bash
go test ./...
go vet ./...
buf lint
buf breaking
goose validate
node tools/vibit check all
```

Future architecture checks should verify:

- Domain packages do not import accepted foundational dependencies directly.
- Platform adapters are the only direct owners of WebSocket and PostgreSQL dependencies.
- Protobuf generated files declare generator and source trace.
- Migration files follow the accepted SQL-first layout and naming convention.
- Go migrations are absent unless a change spec explicitly allows one.

## Alternatives Considered

- Use `github.com/gorilla/websocket` as the first WebSocket library.
- Use `golang.org/x/net/websocket`.
- Hand-roll WebSocket protocol handling.
- Use raw `protoc` and `protoc-gen-go` without Buf.
- Use Buf as the only Protobuf tool without explicitly adopting the official Go runtime and generator.
- Use `database/sql` with `lib/pq`.
- Use `github.com/jackc/pgx/v5` directly behind platform adapters.
- Use `golang-migrate/migrate` for migrations.
- Use `ariga/atlas` for schema-as-code migrations.
- Avoid a migration tool and write a custom runner immediately.
- Adopt an external Go test framework before runtime code exists.

## Rationale

`github.com/coder/websocket` fits vibit's first Go transport adapter because its context-first API and minimal dependency shape are easier for agents to reason about. The lower star count compared with Gorilla is acceptable because WebSocket usage will be isolated behind a vibit platform adapter and replacement path is straightforward.

`google.golang.org/protobuf` and `protoc-gen-go` are the official Go Protobuf implementation and generator. They should be explicit because generated Go code and wire serialization depend on them.

Buf is accepted because vibit needs Protobuf linting, breaking checks, generation orchestration, formatting, and machine-friendly output. These properties directly support agent-native contract verification. The Buf Schema Registry is not required by this decision.

`pgx/v5` is a strong fit because vibit has already chosen PostgreSQL as the first authoritative store. It gives PostgreSQL-specific behavior, pooling, and stable Go APIs without forcing domain modules to import driver packages.

`goose/v3` is a pragmatic migration tool for the first persistent slice because it supports simple SQL migrations, has a CLI and library, and is lighter than schema-as-code systems that would shape the architecture too early.

The standard Go `testing` package is enough for first runtime tests. Additional helpers can be adopted later through the dependency adoption process when a real testing pain appears.

## Agent Reasoning Summary

Accept only dependencies that unblock the first Go runtime path and can be contained behind adapters or generation tooling. Defer dependencies whose use case is not concrete yet. Prefer choices that reduce agent ambiguity through explicit contracts, generated output, verification, and narrow ownership.

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

- `.arch/dependencies.yaml` should mark WebSocket, Protobuf, PostgreSQL driver, and migration tooling as accepted.
- Runtime implementation may later add these dependencies only in the allowed owners declared by `.arch/dependencies.yaml`.
- Domain modules must not directly import `github.com/coder/websocket`, `google.golang.org/protobuf`, `github.com/jackc/pgx/v5`, or `github.com/pressly/goose/v3`.
- Generated Protobuf files must trace to `.proto` sources and generation tooling.
- Migration files should be SQL-first and owned by platform migration tooling.
- S3/MinIO, observability, and external test framework adoption remain open until concrete runtime needs appear.

## Reversal Conditions

Revisit the WebSocket decision if `github.com/coder/websocket` maintenance slows significantly, if the API proves unsuitable for vibit's connection lifecycle, or if Gorilla's ecosystem maturity becomes more valuable than the context-first API.

Revisit Buf adoption if its CLI introduces unacceptable supply-chain, licensing, or operational friction, or if vibit's own manifest-to-Protobuf generator makes Buf unnecessary for local checks.

Revisit `pgx/v5` if a future persistence adapter needs database portability over PostgreSQL-specific features.

Revisit `goose/v3` if migration verification requires richer schema diffing, declarative schema state, or stronger drift detection than SQL-first migrations can provide.

## Follow-Up

- Define the Go package layout before creating `go.mod`.
- Define platform adapter package names for WebSocket transport and PostgreSQL persistence.
- Define `proto/`, `buf.yaml`, and `buf.gen.yaml` layout before generating Protobuf output.
- Define migration directory layout and file naming before adding the first database migration.
- Add import-boundary checks once Go runtime code exists.
