# Impact Analysis

## Affected Modules

- Runtime application session boundary.
- Runtime PostgreSQL platform persistence boundary.
- Authentication module manifest and guide references, because authentication remains a token-linkage participant only.

## Module Ownership Impact

The selected future gate keeps `runtime/internal/app/session` as the storage-neutral repository owner and `runtime/internal/platform/persistence/postgres` as the future adapter owner. Authentication still does not own runtime sessions, session validation, WebSocket sessions, route policy, or direct Nakama/Pitaya compatibility.

## Public Contract Impact

No commands, queries, events, permissions, errors, Protobuf messages, or envelope fields are changed.

## Data And Migration Impact

No migration is added. The existing `runtime_sessions` migration remains the durable SQL source the future adapter must target.

## Test Impact

The selected gate will define required fake-executor PostgreSQL adapter tests, SQL shape checks, row mapping checks, error mapping checks, and transaction-boundary checks. It does not add those adapter tests yet.

## Documentation Impact

The work queue, manifests, AGENTS guides, check rules, change specs, ADR, and conversation memory must record the selected direction.

## Compatibility Risks

Runtime behavior remains unchanged. The main risk is accidentally implementing the adapter or runtime validation during a gate-only step; the follow-up check rule guards that boundary.
