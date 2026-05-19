# Impact Analysis

## Affected Modules

- Runtime application session boundary.
- Authentication module manifest and guide references, because authentication remains a token-linkage participant only.

## Module Ownership Impact

Runtime owns the future session repository interface under `runtime/internal/app/session`. Authentication still does not own runtime sessions, session validation, WebSocket sessions, route policy, or direct Nakama/Pitaya compatibility.

## Public Contract Impact

No commands, queries, events, permissions, errors, Protobuf messages, or envelope fields are changed.

## Data And Migration Impact

No migration is added. The existing `runtime_sessions` migration remains the only durable session SQL source.

## Test Impact

The selected implementation slice requires focused Go unit tests for storage-neutral value types, closed status vocabulary, required-field normalization, UTC time normalization, lifecycle transition mutation normalization, listing query bounds, and no secret material fields.

## Documentation Impact

The work queue, ADRs, change specs, manifests, AGENTS guides, and check rules must record the selected direction.

## Compatibility Risks

Runtime behavior remains unchanged. The main risk is accidentally turning the interface slice into adapter or validation behavior; the follow-up check rule guards that boundary.
