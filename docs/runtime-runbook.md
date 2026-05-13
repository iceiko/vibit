# Runtime Runbook

Status: Draft v0.1
Last updated: 2026-05-13
Scope: First Go runtime process startup and manual verification

This runbook records how to start the first vibit Go runtime process.

The paired Simplified Chinese translation is `docs/runtime-runbook.zh-CN.md`. The English file is authoritative.

## Current Runtime Surface

The current runtime process mounts one gameplay WebSocket endpoint:

```text
/v1/ws
```

The endpoint expects binary WebSocket messages containing `vibit.protocol.v1.Envelope` Protobuf bytes.

The first mounted request loop is:

```text
WebSocket binary frame
-> Protobuf envelope
-> application dispatch
-> inventory command or query handler
-> Protobuf response envelope
-> WebSocket binary frame
```

## Start The Server

From the Go runtime module:

```bash
cd runtime
go run ./cmd/vibit-server
```

The default listen address is:

```text
:8080
```

Override it with:

```bash
VIBIT_ADDR=:9090 go run ./cmd/vibit-server
```

The default runtime store is in memory:

```text
VIBIT_RUNTIME_STORE=memory
```

To start the explicit PostgreSQL-backed inventory composition path, provide both the store selector and a PostgreSQL DSN:

```bash
VIBIT_RUNTIME_STORE=postgres VIBIT_POSTGRES_DSN='postgres://user:pass@127.0.0.1:5432/vibit?sslmode=disable' go run ./cmd/vibit-server
```

Optional PostgreSQL pool settings:

```text
VIBIT_POSTGRES_MAX_CONNS
VIBIT_POSTGRES_MIN_CONNS
```

Normal server startup does not apply migrations. Apply or verify migrations explicitly before using the PostgreSQL store path against a fresh database.

## Manual Verification Path

1. Start the server.
2. Connect a WebSocket client to `ws://127.0.0.1:8080/v1/ws`.
3. Send a binary Protobuf `Envelope` for `inventory.GrantItem` or `inventory.GetInventory`.
4. Confirm the response is a binary Protobuf `Envelope` with the same `request_id`.

Text WebSocket messages are rejected by the transport adapter. JSON is not accepted on this endpoint.

## Current Runtime Assumptions

- The runtime uses an in-memory inventory repository by default.
- `VIBIT_RUNTIME_STORE=postgres` enables the explicit PostgreSQL inventory composition path.
- Inventory bootstrap permissions allow grant and read operations.
- Authentication and session validation are not implemented yet.
- PostgreSQL persistence is wired for inventory runtime composition only when explicitly selected. The persistence boundary is defined in `docs/postgresql-persistence-boundary.md`.
- PostgreSQL migrations are not applied automatically during normal server startup.
- Optional live PostgreSQL verification is defined in `docs/postgresql-verification-environment.md`; it requires `VIBIT_POSTGRES_TEST_DSN` and is not part of default server startup.
- Generated route registration is not implemented yet.

These are bootstrap assumptions for the first request loop, not long-term production policy.

## Verification Commands

Run from the repository root unless noted:

```bash
cd runtime && go test ./...
cd runtime && go vet ./...
node tools/vibit check runtime
node tools/vibit check postgres-env
node tools/vibit check all
```

`node tools/vibit check postgres-env` is a static standards check. It does not connect to PostgreSQL. Live PostgreSQL verification remains opt-in through `VIBIT_POSTGRES_TEST_DSN`.

Run the current live durable inventory verification against a disposable PostgreSQL database with:

```bash
cd runtime && VIBIT_POSTGRES_TEST_DSN='postgres://user:pass@127.0.0.1:5432/vibit_test?sslmode=disable' VIBIT_POSTGRES_TEST_ALLOW_DESTRUCTIVE=1 go test ./internal/platform/protocol/protobuf -run TestPostgresPersistentInventoryRequestLoop -v
```

This test applies the inventory migration explicitly and verifies the WebSocket Protobuf `GrantItem` then `GetInventory` request loop through the PostgreSQL-backed runtime composition. If `VIBIT_POSTGRES_TEST_DSN` is unset, the test skips and records that live PostgreSQL verification was unavailable.

The test uses `drop_schema` cleanup semantics by default. Other cleanup modes are intentionally skipped for this test because migration apply must be verified from a clean schema.
