# Impact Analysis

## Affected Modules

- `runtime`: Work queue direction and runtime manifests record the session creation composition implementation slice.
- `authentication`: The application authentication service becomes the selected composition point for login-time durable session creation, while the session repository remains owned by `runtime/internal/app/session`.

## Module Ownership Impact

No ownership changes are implemented by this direction-selection change. Future implementation remains application-owned under `runtime/internal/app/authentication`, with `runtime/internal/app/session.Repository` as the storage-neutral session lifecycle boundary.

## Public Contract Impact

No command, query, event, error, permission, Protobuf, or generated contract shape changes.

## Data And Migration Impact

No migration or data behavior changes in the direction-selection step.

## Test Impact

No runtime tests are required for this direction-selection step.

## Documentation Impact

The work queue, ADR, conversation memory, manifests, guides, rule catalog, and checks record the selected implementation direction.

## Compatibility Risks

No runtime compatibility risk. The change records the selected direction while preserving explicit deferrals for route policy, WebSocket handshake authentication, Protobuf session carriers, logout, reconnect, and direct Nakama/Pitaya compatibility.
