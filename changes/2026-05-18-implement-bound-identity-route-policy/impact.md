# Impact Analysis

## Affected Modules

- `runtime`: Adds application-owned route policy classification and checks.
- `authentication`: Continues to provide access-token validation but does not own route policy.

## Module Ownership Impact

Route policy remains in `runtime/internal/app`. WebSocket transport, Protobuf adapters, authentication service, session repository, and persistence adapters do not gain policy ownership.

## Public Contract Impact

No public commands, queries, events, permissions, Protobuf sources, generated output, WebSocket handshake behavior, or database schemas are added or changed.

## Runtime Behavior Impact

Default ordinary protected route behavior is unchanged: request-level authenticated request wrappers still carry access-token proof and route validation still clears `SessionValidated` before domain dispatch.

Explicit route policy configuration can now classify a route as bound-connection, session-validated, or bound-session, but no default production domain route is reclassified by this change.

## Test Impact

Focused Go tests cover default classification, explicit policy classification, metadata-only rejection, bound route success/failure, session route success/failure, bound-session identity agreement, and unchanged request-token behavior.

## Reference Alignment

Nakama informs the route policy goal: authenticated session material can gate gameplay access, but lifecycle concerns such as logout and refresh stay separate. Pitaya informs the layering: acceptor/transport, session context, route policy, and handlers stay separated.

## Compatibility Risks

The main risk is accidentally treating a bound connection or persisted session as global proof. The implementation requires explicit route classification and keeps ordinary protected routes on request-token proof.
