# Impact: Pitaya-Aligned Session Binding Kick Disconnect Session Data Source-First Map

This change adds a source-first repository inspection map for Pitaya-aligned session binding, kick/disconnect, and session data vocabulary.

## Runtime Impact

No runtime behavior is added.

The map reports current vibit session-lifecycle-adjacent concepts:

- first-message connection binding route;
- request-level access-token validation and request identity handoff;
- request identity and connection metadata;
- absence of a general session data scope;
- application-owned active connection registry;
- logout service token revocation with unchanged transport close policy;
- absence of kick policy or route;
- close handoff and connection registry cleanup;
- existing transport close reason mapping;
- server-owned presence lifecycle snapshot.

It maps those concepts to future session binding, kick/disconnect, and session data vocabulary without changing dispatch, transport, protocol, domain, persistence, startup, metrics, tracing, or distributed runtime behavior.

## Tooling Impact

`tools/vibit` gains:

- `node tools/vibit inspect pitaya-session-lifecycle --json`;
- `runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_source_first_map`;
- repository checks for W-0269 artifacts, inspection output markers, deferrals, redaction flags, and W-0270 next-ready state.

## Documentation And Memory Impact

The change records `ADR-0177`, a conversation memory entry, work item completion, runtime/reference/convention/contract manifest markers, and continuation docs pointing to W-0270.

## Explicit Non-Impact

- No runtime behavior.
- No protocol shape.
- No generated output.
- No persistence.
- No dependencies.
- No authentication/session behavior change.
- No session binding behavior, kick/disconnect behavior, server-initiated disconnect behavior, server-initiated kick behavior, session data behavior, session data persistence, session unbind behavior, close reason behavior change, connection-session handoff behavior, or presence-session handoff behavior.
- No acceptor behavior, TCP acceptor, WebSocket behavior change, or connection lifecycle behavior change.
- No route handler implementation, handler routing, pipeline, middleware, serializer, forwarding, or backend route targeting behavior.
- No metrics endpoint, tracing pipeline, distributed runtime, service discovery, RPC, remote call, frontend/backend role, cluster-safe session routing, group, broadcast, hosted, SDK, release, or direct compatibility surface.
