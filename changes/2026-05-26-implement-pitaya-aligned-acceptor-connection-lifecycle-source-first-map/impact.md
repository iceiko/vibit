# Impact: Pitaya-Aligned Acceptor And Connection Lifecycle Source-First Map

This change adds a source-first repository inspection map for Pitaya-aligned acceptor and connection lifecycle vocabulary.

## Runtime Impact

No runtime behavior is added.

The map reports current vibit acceptor and connection-lifecycle-adjacent concepts:

- WebSocket accept loop;
- server-observed connection id metadata;
- server-observed connection epoch metadata;
- first-message connection binding route;
- application-owned active connection registry;
- transport close to application policy handoff;
- server-owned presence lifecycle snapshot.

It maps those concepts to future acceptor and connection lifecycle vocabulary without changing dispatch, transport, protocol, domain, persistence, startup, metrics, tracing, or distributed runtime behavior.

## Tooling Impact

`tools/vibit` gains:

- `node tools/vibit inspect pitaya-acceptor-connection --json`;
- `runtime.pitaya_aligned_acceptor_connection_lifecycle_source_first_map`;
- repository checks for W-0266 artifacts, inspection output markers, deferrals, redaction flags, and W-0267 next-ready state.

## Documentation And Memory Impact

The change records `ADR-0174`, a conversation memory entry, work item completion, runtime/reference/convention/contract manifest markers, and continuation docs pointing to W-0267.

## Explicit Non-Impact

- No runtime behavior.
- No protocol shape.
- No generated output.
- No persistence.
- No dependencies.
- No authentication/session behavior change.
- No acceptor behavior, TCP acceptor, WebSocket behavior change, connection lifecycle behavior change, session binding behavior, kick/disconnect behavior, or concrete socket close behavior change.
- No route handler implementation, handler routing, pipeline, middleware, serializer, forwarding, or backend route targeting behavior.
- No metrics endpoint, tracing pipeline, distributed runtime, service discovery, RPC, remote call, frontend/backend role, cluster-safe session routing, group, broadcast, hosted, SDK, release, or direct compatibility surface.
