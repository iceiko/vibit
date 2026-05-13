# Impact Analysis

## Affected Modules

- `runtime`: adds process startup wiring and a narrow server assembly package.
- `inventory`: supplies the first non-persistent runtime repository and bootstrap permission policy used only for the in-memory request loop.

## Module Ownership Impact

No ownership changes.

`runtime/cmd/vibit-server/` owns process startup and lifecycle only. Runtime assembly lives outside `cmd` so tests can exercise wiring without launching a long-running process. Inventory business behavior remains in the inventory module.

## Public Contract Impact

No public command, query, event, error, permission, or Protobuf schema changes.

The `/v1/ws` endpoint was already planned by `ADR-0015`; this change mounts it for the first runtime process.

## Data And Migration Impact

No durable data or migration impact. The first process wiring uses an in-memory inventory repository until the PostgreSQL boundary is defined.

## Test Impact

Adds a narrow WebSocket integration test that sends a binary Protobuf envelope through `/v1/ws` and decodes the Protobuf response envelope.

## Documentation Impact

Updates runtime manifests and runtime docs to reflect process wiring and the manual run path.

## Compatibility Risks

The endpoint path becomes implemented for the first runtime process. Changing it later after clients exist will require a compatibility-sensitive change spec and ADR.
