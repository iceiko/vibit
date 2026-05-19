# Impact Analysis

## Affected Modules

- Runtime PostgreSQL persistence boundary.
- Runtime application session boundary.
- Authentication module manifest and guide references, because authentication remains token-record-linkage only.

## Module Ownership Impact

The future PostgreSQL session adapter is platform persistence-owned under `runtime/internal/platform/persistence/postgres`. The storage-neutral repository interface remains owned by `runtime/internal/app/session`. Authentication still does not own runtime sessions, adapter behavior, validation policy, WebSocket sessions, route policy, or direct Nakama/Pitaya compatibility.

## Public Contract Impact

No public commands, queries, events, errors, permissions, Protobuf messages, or envelope fields are changed.

## Data And Migration Impact

No migration is added or changed. The gate references the existing `runtime_sessions` migration as the table the future adapter must target.

## Test Impact

No Go tests are added by this gate-only change. The standard defines future adapter test requirements for interface conformance, SQL shape, normalization, row mapping, error mapping, listing bounds, transaction neutrality, and redaction.

## Documentation Impact

`ADR-0063`, change specs, conversation memory, manifests, AGENTS guides, and repository checks are updated.

## Compatibility Risks

Runtime behavior is unchanged. The gate reduces future compatibility risk by preventing a future adapter from silently owning token validation, request identity construction, route policy, WebSocket behavior, reconnect state, or direct Nakama/Pitaya APIs.
