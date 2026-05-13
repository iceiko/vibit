# Go Runtime Agent Guide

Status: Draft v0.1
Last updated: 2026-05-13
Scope: `runtime/` Go server runtime workspace
Canonical source: `../CONSTITUTION.md`, `../AGENTS.md`, and `../decisions/ADR-0014-go-runtime-layout-and-boundaries.md`

This guide applies to the first Go server runtime implementation.

The paired Simplified Chinese translation is `runtime/AGENTS.zh-CN.md`. The English file is authoritative.

## 1. Purpose

`runtime/` is the first Go module for vibit's server runtime.

The Go module path is:

```text
github.com/iceiko/vibit/runtime
```

The runtime exists to prove vibit's core claim through a small, long-lived backend slice:

```text
requirement -> spec -> contract -> generated shape -> handwritten logic -> tests -> verification -> docs
```

Do not treat this workspace as a disposable demo.

## 2. Required Reading

Before changing files under `runtime/`, read:

- `../CONSTITUTION.md`
- `../AGENTS.md`
- `../.arch/runtime.yaml`
- `../.arch/dependencies.yaml`
- `../.arch/contracts.yaml`
- `../docs/generated-output.md`
- `../docs/runtime-protocol-adapter.md`
- `../docs/postgresql-persistence-boundary.md`, before persistence work
- `../docs/runtime-runbook.md`
- `../decisions/ADR-0014-go-runtime-layout-and-boundaries.md`
- `../decisions/ADR-0018-runtime-protocol-adapter-boundary.md`
- `../decisions/ADR-0020-postgresql-persistence-boundary.md`, before persistence work
- The affected module manifest, such as `../modules/inventory/module.yaml`
- The relevant change spec under `../changes/`

## 3. Package Ownership

Use these package boundaries:

- `cmd/vibit-server/`: process startup, configuration wiring, and lifecycle.
- `internal/app/`: command/query dispatch, application composition, and transaction orchestration.
- `internal/platform/transport/ws/`: WebSocket transport adapter and `github.com/coder/websocket` ownership.
- `internal/platform/protocol/protobuf/`: Protobuf framing, envelope conversion, and wire message adaptation.
- `internal/platform/persistence/postgres/`: PostgreSQL adapter implementation and `github.com/jackc/pgx/v5` ownership.
- `internal/platform/migrations/`: migration tooling invocation and validation.
- `internal/platform/events/`: event recording and publication mechanisms.
- `internal/platform/tx/`: transaction boundary and unit-of-work interfaces.
- `internal/modules/<module>/`: handwritten domain module runtime logic.
- `internal/generated/contracts/`: generated Go contract shapes.
- `internal/generated/proto/`: generated Go Protobuf files.
- `migrations/postgres/`: SQL-first PostgreSQL migration sources.

## 4. Dependency Rules

Domain modules must not import third-party transport, protocol, persistence, migration, object-storage, or framework dependencies directly.

Allowed owner packages:

- `github.com/coder/websocket`: `internal/platform/transport/ws/` only.
- `google.golang.org/protobuf`: generated protocol packages and protocol adapter packages only.
- `github.com/jackc/pgx/v5`: `internal/platform/persistence/postgres/` only.
- `github.com/pressly/goose/v3`: `internal/platform/migrations/` only.

Do not add new foundational dependencies without checking `../.arch/dependencies.yaml` and creating the required adoption record.

## 5. Runtime Boundary Rules

Runtime protocol handoff must follow `../docs/runtime-protocol-adapter.md`.

WebSocket transport reads and writes frames. Protobuf protocol adaptation decodes and encodes envelopes. Application dispatch routes commands and queries. Domain modules enforce invariants. Generated packages provide shapes only.

WebSocket transport handlers pass opaque frame bytes to injected protocol/application composition. They do not adapt requests into commands or queries directly, and they must not hide business logic.

State-changing commands should enter through `internal/app/` and run inside an application-owned unit of work. Repository mutations and domain event recording should happen inside that same unit of work.

The current transaction skeleton is `internal/platform/tx.Runner`, `internal/platform/tx.UnitOfWork`, and `internal/app.TransactionalDispatcher`. Application code may import this transaction boundary package, but it must not import persistence, migration, protocol, or transport platform adapters. Query routes should pass through without a write unit of work by default.

Query handlers should not mutate state and do not require a write transaction by default.

Event publication outside the transaction remains deferred until vibit adopts an explicit event delivery or outbox standard.

PostgreSQL persistence work must follow `../docs/postgresql-persistence-boundary.md`. Repository interfaces stay module-owned, `pgx` stays under `internal/platform/persistence/postgres/`, `goose` stays under `internal/platform/migrations/`, and SQL migration sources stay under `migrations/postgres/`.

For the first durable inventory implementation, `GrantItem` must use a transaction-bound repository and call `LockInventoryForMutation` before reading current items and applying capacity-sensitive mutations. The returned `MutationLock` is a locked aggregate view, not a transaction owner. Repositories must not silently open independent write transactions for command flows.

The first inventory migration source is `migrations/postgres/000001_create_inventory_state.sql`. It creates `inventory_accounts`, `inventory_items`, and `inventory_item_grants`. Run `node ../tools/vibit check migrations` when migration sources or migration guidance change. Migration apply/rollback verification remains pending until migration tooling can run against a disposable PostgreSQL environment.

## 6. Generated Files

Generated files are immutable to non-system agents.

If generated output is wrong, change the source contract, schema, template, or generator. Do not hand-edit generated files unless a change spec or Agent Decision Record explicitly grants `generated_file_override`.

Go Protobuf generated output under `internal/generated/proto/` must be produced from `../proto/` sources through the accepted Buf and `protoc-gen-go` path. Files under that root must be generated `*.pb.go` files with the `protoc-gen-go` marker and source trace, or temporary `.gitkeep` placeholders while generation has not run.

Do not place handwritten runtime code under `internal/generated/proto/` or `internal/generated/contracts/`.

## 7. Current State

This runtime workspace now has the first generated Protobuf output, the first narrow runtime handoff slice, the first WebSocket transport adapter, a small application dispatch skeleton for command and query routes, the first transaction boundary skeleton, the first inventory repository/policy/handler runtime boundary with a command-safe mutation lock, the first inventory Protobuf/domain payload bridge, the first application-error-to-Protobuf-error-envelope mapper, the first frame-to-Protobuf-to-application composition adapter, a package-local request-loop test fixture for Protobuf command/query tests, and minimal process wiring that mounts `/v1/ws`.

The workspace has a documented PostgreSQL persistence boundary, transaction skeleton, and first inventory migration source, but it still does not implement PostgreSQL repository adapters, migration apply/rollback tooling, persistent runtime wiring, generated route registration, generated protocol bridge creation, authentication/session validation, or catalog-driven error retryability yet.

The first manual process run path is:

```bash
cd runtime
go run ./cmd/vibit-server
```

## 8. Verification

Run repository verification from the repository root:

```bash
node tools/vibit check runtime
node tools/vibit check generated
node tools/vibit check migrations
node tools/vibit check all
```

When Go source files exist and the local Go toolchain is available, runtime verification should include:

```bash
go test ./...
go vet ./...
```

Do not claim Go test verification when the Go toolchain is unavailable or tests were not run.
