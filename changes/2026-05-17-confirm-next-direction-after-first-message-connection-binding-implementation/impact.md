# Impact Analysis

## Affected Modules

- Runtime workflow and manifests.
- Authentication module guidance, because authentication remains storage-neutral and must not own runtime sessions.

## Module Ownership Impact

This direction confirmation does not change runtime ownership by itself. It selects a future gate to define runtime-owned session persistence schema boundaries. The player module continues to own player accounts, the authentication module continues to own credential/token verifier boundaries, and WebSocket transport continues to own only connection plumbing.

## Public Contract Impact

No command, query, event, error, permission, Protobuf, or envelope shape changes are made by this confirmation step.

## Data And Migration Impact

No SQL migration, session table, repository, adapter, cleanup job, or data migration is added. The next gate may define a future `runtime_sessions` schema candidate but still must not implement it unless a later work item authorizes it.

## Test Impact

No Go tests are required for this direction-only change. Repository checks must verify the work queue and manifests.

## Documentation Impact

The work queue, manifests, AGENTS guides, change spec, and conversation log record the selected direction and the deferrals.

## Compatibility Risks

There is no runtime compatibility impact because no runtime behavior changes. WebSocket handshake authentication, route-policy bound identity, logout/revocation, reconnect behavior, and direct Nakama/Pitaya compatibility remain deferred.
