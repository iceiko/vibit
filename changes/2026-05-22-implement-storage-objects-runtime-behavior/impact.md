# Impact Analysis

## Affected Modules

- `runtime`
- `storage`
- `workflow`
- `reference`

## Module Ownership Impact

Application-owned storage object runtime behavior is added under `runtime/internal/app/storage`.

The storage module continues to own storage-neutral value types and repository vocabulary through `runtime/internal/modules/storage`. The PostgreSQL platform package continues to own SQL mapping. Protocol, transport, generated output, startup wiring, authentication, session validation, and migrations remain unchanged.

## Public Contract Impact

No public commands, queries, events, permissions, Protobuf messages, WebSocket routes, or generated contract outputs are added.

The implementation adds application-internal public error codes for later handler mapping:

- `STORAGE_OBJECT_INVALID_REQUEST`
- `STORAGE_OBJECT_NOT_FOUND`
- `STORAGE_OBJECT_ALREADY_EXISTS`
- `STORAGE_OBJECT_VERSION_MISMATCH`
- `STORAGE_OBJECT_UNAVAILABLE`
- `STORAGE_OBJECT_FORBIDDEN`

## Data And Migration Impact

No migrations are added or changed. The service references the existing `storage_objects` logical table only through `runtime/internal/modules/storage.Repository` and unit-of-work repository handoff.

## Test Impact

Adds focused Go tests for:

- dependency validation;
- metadata-only identity rejection before repository access;
- validated player owner derivation;
- get/list/put/delete behavior;
- value copying and size validation;
- expected-version update/delete behavior;
- redacted conflict mapping;
- unit-of-work storage repository handoff.

Default tests do not require live PostgreSQL.

## Documentation Impact

Adds:

- `ADR-0117`
- conversation log
- change spec files
- manifest, guide, and check-rule updates

No new public standard document is required because `ADR-0116` is the ratified implementation gate.

## Compatibility Risks

No public API compatibility risk is introduced because this slice does not expose protocol routes or generated output.

The main behavioral risks are internal service semantics for put/create-or-replace and expected-version handling. These are covered by focused tests and remain adjustable before protocol exposure.
