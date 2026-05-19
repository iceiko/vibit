# Impact Analysis

## Affected Modules

- `runtime`: Future route policy remains application-owned under `runtime/internal/app`.
- `authentication`: Authentication remains proof/token/session-creation composition and does not own route policy.

## Module Ownership Impact

No Go ownership changes are implemented. The gate records the future ownership boundary:

- Route policy: `runtime/internal/app`
- Authentication service: `runtime/internal/app/authentication`
- Session repository: `runtime/internal/app/session`
- Protocol adapters: `runtime/internal/platform/protocol/protobuf`
- WebSocket transport: `runtime/internal/platform/transport/ws`

## Public Contract Impact

No public command, query, event, permission, error, Protobuf source, generated output, or WebSocket behavior changes.

## Data And Migration Impact

No migrations are added or changed.

## Test Impact

No Go tests are required for a gate-only documentation boundary. The future implementation must test public route behavior, request-token policy, bound identity policy, session-validated policy, identity mismatch fail-closed behavior, redaction, WebSocket neutrality, and unchanged Protobuf shape.

## Reference Alignment

Nakama informs the relationship between authenticated session material and gameplay access. vibit adapts this into explicit route policy families rather than direct API compatibility.

Pitaya informs separation between acceptors, sessions, and handlers. vibit adapts this by keeping route policy in application code and keeping handlers free of transport credential parsing.

## Compatibility Risks

The main risk is accidentally treating bound connection identity or session id metadata as proof. The gate explicitly forbids that until a later route-scoped implementation slice authorizes it.
