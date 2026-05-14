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
- `../docs/postgresql-verification-environment.md`, before live PostgreSQL verification work
- `../docs/authentication-token-session-validation.md`, before authentication, token, credential, external identity, session persistence, request identity trust, WebSocket handshake, player handler, or player route work
- `../docs/runtime-runbook.md`
- `../decisions/ADR-0014-go-runtime-layout-and-boundaries.md`
- `../decisions/ADR-0018-runtime-protocol-adapter-boundary.md`
- `../decisions/ADR-0020-postgresql-persistence-boundary.md`, before persistence work
- `../decisions/ADR-0022-player-account-postgresql-schema-boundary.md`, before player account persistence work
- `../decisions/ADR-0023-authentication-token-session-validation-design-boundary.md`, before authentication/session design or implementation work
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

The first PostgreSQL inventory repository adapter is `internal/platform/persistence/postgres/inventory_repository.go`. Construct it with `NewInventoryRepositoryForUnitOfWork` and pass an executor supplied by the application-owned unit of work, such as a `pgx.Tx` or compatible test executor. The adapter must not call `BEGIN`, `COMMIT`, or `ROLLBACK`; transaction lifetime belongs to `internal/platform/tx` and `internal/app`.

PostgreSQL configuration is owned by `internal/platform/persistence/postgres/config.go`. It reads `VIBIT_POSTGRES_DSN`, `VIBIT_POSTGRES_MAX_CONNS`, and `VIBIT_POSTGRES_MIN_CONNS`, builds pgx pool configuration, and must not require a live PostgreSQL server in normal unit tests. Connection strings and credentials must come from environment or explicit runtime input and must not be committed.

The pgx-backed transaction runner is `internal/platform/persistence/postgres/runner.go`. It implements `internal/platform/tx.Runner` while keeping pgx transaction handles inside the PostgreSQL platform package. It commits successful command units of work, rolls back failed callback units of work, and exposes package-owned helpers such as `UnitOfWork.NewInventoryRepository` for future persistent composition. Do not import the PostgreSQL runner from `internal/app/` or domain modules; persistent runtime wiring must happen in an approved composition boundary.

`GrantItemMutation` carries `event_id`, `occurred_at`, and `reason` so the PostgreSQL adapter can persist `inventory_item_grants` in the same executor path as the item quantity update.

The first inventory migration source is `migrations/postgres/000001_create_inventory_state.sql`. It creates `inventory_accounts`, `inventory_items`, and `inventory_item_grants`. Run `node ../tools/vibit check migrations` when migration sources or migration guidance change. Migration status and apply behavior are covered by the opt-in live durable inventory request-loop verification when `VIBIT_POSTGRES_TEST_DSN` is set.

The ratified player account PostgreSQL schema boundary is documented in `../docs/postgresql-persistence-boundary.md` and `../decisions/ADR-0022-player-account-postgresql-schema-boundary.md`. The first player account migration source is `migrations/postgres/000002_create_player_account_state.sql`. That migration creates only `player_accounts` and `player_account_events` lifecycle state. It must not add credentials, password hashes, external identity links, access tokens, refresh tokens, runtime session rows, WebSocket connection state, request identity validation results, inventory state, or permission grants.

The player account repository interface boundary is `internal/modules/player/repository.go`. It is storage-neutral domain code and may define account lifecycle structs, `Repository.CreatePlayerAccount`, `Repository.GetPlayerAccount`, and durable mutation metadata for persistence adapters. The PostgreSQL adapter is `internal/platform/persistence/postgres/player_account_repository.go`, with focused tests at `internal/platform/persistence/postgres/player_account_repository_test.go`. It uses `NewPlayerAccountRepositoryForUnitOfWork(executor)`, implements `player.Repository`, receives its executor from the application-owned unit of work, and must not call `BEGIN`, `COMMIT`, or `ROLLBACK`. `UnitOfWork.NewPlayerAccountRepository` is a PostgreSQL package helper and must not expose pgx to application or domain packages.

The player account PostgreSQL adapter does not authorize runtime handlers, WebSocket routes, authentication, token behavior, credential storage, external identity linking, or session persistence. The adapter may only write `player_accounts`, write `player_account_events` for `PlayerAccountCreated`, and read current lifecycle rows from `player_accounts` until a later change ratifies more behavior.

The authentication, token, and session validation design boundary is documented in `../docs/authentication-token-session-validation.md` and `../decisions/ADR-0023-authentication-token-session-validation-design-boundary.md`. It separates authentication proof, login methods, tokens, credentials, external identity links, runtime sessions, request identity, WebSocket handshake authentication, player account lifecycle, transport connection metadata, and Protobuf envelope metadata. The current `MetadataOnlySessionValidator` is a non-authenticated bootstrap path. Do not treat metadata-only `player_id`, `session_id`, or `connection_id` as production proof, and do not add authentication runtime code, token parsing, credential lookup, external identity linking, session persistence, Protobuf envelope authentication changes, WebSocket handshake authentication, runtime player handlers, or WebSocket routes until separately ratified. `runtime.authentication_token_session_boundary` is the repository check rule for this boundary.

The first explicit PostgreSQL migration runner is `internal/platform/migrations/postgres.go`. It owns `github.com/pressly/goose/v3`, accepts a caller-supplied `*sql.DB` and migration source filesystem or directory, lists SQL migration sources, reports structured status, and applies pending migrations only when explicitly invoked. Do not wire it into normal `cmd/vibit-server` startup without a change spec.

Live PostgreSQL verification is governed by `../docs/postgresql-verification-environment.md`. It is opt-in through `VIBIT_POSTGRES_TEST_DSN`; normal unit tests, `node ../tools/vibit check runtime`, and default repository checks must not require a running PostgreSQL server. When a live PostgreSQL check is skipped because no disposable DSN is available, record that explicitly.

## 6. Generated Files

Generated files are immutable to non-system agents.

If generated output is wrong, change the source contract, schema, template, or generator. Do not hand-edit generated files unless a change spec or Agent Decision Record explicitly grants `generated_file_override`.

Go Protobuf generated output under `internal/generated/proto/` must be produced from `../proto/` sources through the accepted Buf and `protoc-gen-go` path. Files under that root must be generated `*.pb.go` files with the `protoc-gen-go` marker and source trace, or temporary `.gitkeep` placeholders while generation has not run.

Do not place handwritten runtime code under `internal/generated/proto/` or `internal/generated/contracts/`.

## 7. Current State

This runtime workspace now has the first generated Protobuf output, the first narrow runtime handoff slice, the first WebSocket transport adapter, a small application dispatch skeleton for command and query routes, the first transaction boundary skeleton, the first inventory repository/policy/handler runtime boundary with a command-safe mutation lock, the first PostgreSQL configuration parser, the first pgx-backed transaction runner adapter, the first PostgreSQL inventory repository adapter, the first inventory Protobuf/domain payload bridge, the first application-error-to-Protobuf-error-envelope mapper, the first frame-to-Protobuf-to-application composition adapter, a package-local request-loop test fixture for Protobuf command/query tests, minimal process wiring that mounts `/v1/ws`, an explicit PostgreSQL inventory runtime composition path, and an opt-in live PostgreSQL durable inventory request-loop verification test.

The workspace has a documented PostgreSQL persistence boundary, transaction skeleton, PostgreSQL configuration parser, pgx-backed transaction runner, first inventory migration source, first explicit migration apply/status runner, first PostgreSQL repository adapter, explicit runtime store selection, a ratified player account PostgreSQL lifecycle schema boundary with its first migration source, storage-neutral repository interface, focused PostgreSQL adapter implementation, PostgreSQL unit-of-work factory helper, and a ratified authentication/token/session validation design boundary without runtime authentication implementation, plus a live verification test that skips unless `VIBIT_POSTGRES_TEST_DSN` is set. `VIBIT_RUNTIME_STORE=memory` remains the default. `VIBIT_RUNTIME_STORE=postgres` enables PostgreSQL-backed inventory composition when `VIBIT_POSTGRES_DSN` is provided. The workspace still does not implement generated route registration, generated protocol bridge creation, production authentication/session validation, runtime player account handlers, WebSocket player routes, automatic startup migrations, or catalog-driven error retryability yet.

The first manual process run path is:

```bash
cd runtime
go run ./cmd/vibit-server
```

The first explicit persistent process run path is:

```bash
cd runtime
VIBIT_RUNTIME_STORE=postgres VIBIT_POSTGRES_DSN='postgres://user:pass@127.0.0.1:5432/vibit?sslmode=disable' go run ./cmd/vibit-server
```

Migrations are not applied automatically during normal server startup.

The first opt-in live durable inventory verification command is:

```bash
cd runtime
VIBIT_POSTGRES_TEST_DSN='postgres://user:pass@127.0.0.1:5432/vibit_test?sslmode=disable' VIBIT_POSTGRES_TEST_ALLOW_DESTRUCTIVE=1 go test ./internal/platform/protocol/protobuf -run TestPostgresPersistentInventoryRequestLoop -v
```

If `VIBIT_POSTGRES_TEST_DSN` is unset, this test skips and records that live PostgreSQL verification was unavailable. The first live execution has passed on local Termux PostgreSQL 18.2.

## 8. Verification

Run repository verification from the repository root:

```bash
node tools/vibit check runtime
node tools/vibit check generated
node tools/vibit check migrations
node tools/vibit check postgres-env
node tools/vibit check all
```

When Go source files exist and the local Go toolchain is available, runtime verification should include:

```bash
go test ./...
go vet ./...
```

Do not claim Go test verification when the Go toolchain is unavailable or tests were not run.
