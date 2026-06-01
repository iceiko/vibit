# Change Request: Implement Pitaya-Aligned Session Binding Kick Disconnect Session Data Source-First Map

Implement `W-0269 Implement Pitaya-aligned session binding, kick/disconnect, and session data source-first map`.

The maintainer asked to continue toward Pitaya. The narrow request is to make the session binding, kick/disconnect, and session data vocabulary from `ADR-0176` inspectable through a source-first repository inspection map.

## Scope

- Add `node tools/vibit inspect pitaya-session-lifecycle --json`.
- Register `runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_source_first_map`.
- Record `ADR-0177` and the W-0269 change artifacts.
- Open `W-0270 Select next Pitaya-aligned direction after session binding, kick/disconnect, and session data map` as next-ready.

## Non-Goals

- No session binding behavior.
- No kick/disconnect behavior.
- No server-initiated disconnect behavior.
- No server-initiated kick behavior.
- No session data behavior.
- No session data persistence.
- No session unbind behavior.
- No close reason behavior change.
- No connection-session handoff behavior.
- No presence-session handoff behavior.
- No acceptor behavior.
- No TCP acceptors.
- No WebSocket behavior changes.
- No connection lifecycle behavior changes.
- No route handler implementation.
- No handler routing behavior.
- No handler pipeline behavior.
- No pipeline middleware behavior.
- No backend route targeting.
- No protocol messages or routes.
- No Protobuf source.
- No generated output.
- No repository interfaces, PostgreSQL adapters, migrations, or dependencies.
- No metrics endpoints, tracing pipelines, service discovery, RPC, remote calls, frontend/backend roles, cluster-safe session routing, or distributed runtime behavior.
- No hosted deployment, SDK publication, release artifact, or direct Nakama/Pitaya API compatibility.
