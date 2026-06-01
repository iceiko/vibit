# Request

Define `W-0268`, a gate-only Pitaya-aligned session binding, kick/disconnect, and session data boundary gate after the acceptor and connection lifecycle direction selection.

The change must:

- Accept `ADR-0176`.
- Define the English standard and Simplified Chinese translation.
- Register `runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_boundary_gate`.
- Map current first-message binding, runtime session validation, request identity/session metadata, active connection registry, logout/close handoff, and presence lifecycle surfaces to future session lifecycle vocabulary.
- Open `W-0269 Implement Pitaya-aligned session binding, kick/disconnect, and session data source-first map` as next-ready.

No session binding behavior, kick/disconnect behavior, session data behavior, session data persistence, acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, route handler implementation, backend route targeting, protocol messages or routes, Protobuf sources, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, metrics/tracing behavior, hosted surfaces, SDKs, or direct Nakama/Pitaya API compatibility may be added.
