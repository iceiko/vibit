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

## Manual Verification Path

1. Start the server.
2. Connect a WebSocket client to `ws://127.0.0.1:8080/v1/ws`.
3. Send a binary Protobuf `Envelope` for `inventory.GrantItem` or `inventory.GetInventory`.
4. Confirm the response is a binary Protobuf `Envelope` with the same `request_id`.

Text WebSocket messages are rejected by the transport adapter. JSON is not accepted on this endpoint.

## Current Runtime Assumptions

- The runtime uses an in-memory inventory repository.
- Inventory bootstrap permissions allow grant and read operations.
- Authentication and session validation are not implemented yet.
- PostgreSQL persistence is not wired yet. The persistence boundary is defined in `docs/postgresql-persistence-boundary.md`.
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
