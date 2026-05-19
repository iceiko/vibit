# ADR-0075: Active Connection Registry Single Process Implementation

Status: Accepted
Date: 2026-05-18
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-active-connection-registry-gate/`
- `changes/2026-05-18-implement-active-connection-registry-single-process/`

Related conversations:

- `conversations/2026-05-18-active-connection-registry-single-process-implementation.md`

Related artifacts:

- `runtime/internal/app/connection/registry.go`
- `runtime/internal/app/connection/registry_test.go`
- `docs/active-connection-registry-gate.md`
- `decisions/ADR-0074-active-connection-registry-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/authentication/AGENTS.md`
- `modules/authentication/AGENTS.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0074` defined a gate-only boundary for future active connection registry behavior. That gate established the first posture: application-owned, single-process, in-memory, non-durable registry state under `runtime/internal/app/connection`.

The work queue reached `M-089/W-0161`, a confirmation gate. The maintainer asked the agent, in Chinese, to recommend the next ten steps and continue, with Nakama and Pitaya as important game-server references.

The missing implementation was now narrow and concrete: vibit needed a safe server-owned model for active connection targeting before future logout socket invalidation, kick/disconnect, duplicate replacement, reconnect/epoch behavior, presence, rooms, or match runtime can be designed.

## Decision

Select:

```text
implement_active_connection_registry_single_process
```

Implement the first active connection registry in:

```text
runtime/internal/app/connection/registry.go
runtime/internal/app/connection/registry_test.go
```

The implementation adds an application-owned `InMemoryRegistry` with these capabilities:

- `RegisterOpenConnection`
- `BindConnectionIdentity`
- `MarkConnectionClosed`
- `FindConnectionByID`
- `ListConnectionsByPlayerID`
- `ListConnectionsByRuntimeSessionID`
- `ListConnectionsByAccessTokenRecordID`
- `MarkConnectionInvalidated`

The registry records server-observed connection id, connection epoch, active or terminal state, validated player linkage, optional runtime session id, optional access-token record id, and lifecycle timestamps. It rejects duplicate active records for the same connection id and epoch, rejects binding without validated player identity, excludes closed and invalidated records from active target lists, returns record copies, and keeps raw proof material out of record shape.

This ADR does not wire the registry into WebSocket transport, Protobuf adapters, startup composition, route policy, logout, runtime session validation, or authentication service behavior.

It does not close WebSocket connections, add kick/disconnect behavior, revoke runtime sessions, replace duplicate connections automatically, add reconnect or epoch protocol behavior, add Protobuf logout routes, add protocol session carriers, change the existing Protobuf envelope, change WebSocket handshake authentication, add transport credential carriers, add durable or distributed registry storage, add cleanup jobs, add dependencies, add memory durable session behavior, broaden game backend modules, or adopt direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Keep the registry gate only and defer implementation again.
- Wire the registry directly into WebSocket open/close handling in the same slice.
- Make logout close active sockets immediately after token revocation.
- Treat first-message `BindConnection` state as a sufficient registry.
- Add a distributed registry or Redis-like dependency immediately.
- Add reconnect, duplicate replacement, or epoch protocol behavior in the same slice.
- Copy Nakama or Pitaya public APIs directly.

## Rationale

Nakama shows that authenticated session material and realtime sockets need coordinated lifecycle state. vibit adapts that lesson by adding explicit active connection records that can later support lifecycle policy, without copying Nakama session APIs or socket behavior.

Pitaya shows the value of separating acceptors, sessions, route handlers, and connection management. vibit adapts that lesson by keeping registry state in the application layer and keeping WebSocket transport credential-neutral and policy-neutral.

The single-process in-memory posture is intentionally conservative. It is enough for the first modular monolith runtime and future policy tests, while avoiding hidden distributed assumptions before cluster routing, server-to-server RPC, service discovery, or Redis-like storage are ratified.

## Agent Reasoning Summary

After service-level logout and a registry gate, the next useful step is not a public logout route or socket close behavior. The server first needs a safe local targeting primitive for active connections.

The implementation keeps the registry as state, not policy. It can mark an active record invalidated, but it does not decide whether that invalidation should close the socket, send a message, fail future routes, or wait for reconnect handling.

## Decision Weights

```yaml
decision_weights:
  active_connection_targeting_foundation: high
  nakama_pitaya_alignment: high
  transport_protocol_app_separation: high
  future_logout_revocation_correctness: high
  reconnect_and_presence_readiness: medium
  implementation_scope_control: high
  immediate_protocol_surface: low
  direct_nakama_pitaya_api_compatibility: low
confidence: high
```

## Consequences

- `runtime.active_connection_registry_single_process_implementation` becomes the repository check rule for this slice.
- The runtime now has a local application-owned active connection registry primitive.
- Future logout/revocation, WebSocket close, kick/disconnect, duplicate replacement, reconnect/epoch, presence, rooms, parties, and match runtime work can target server-owned registry records instead of client metadata.
- WebSocket transport remains credential-neutral and does not own authentication state.
- Protobuf sources, generated output, migrations, dependencies, and startup wiring remain unchanged.
- The work queue blocks again after this implementation at `M-091/W-0163`.

## Reversal Conditions

Revisit this decision if a future ADR chooses transport-owned authentication state, requires handshake-level authentication as the primary registry binding source, adopts direct Nakama/Pitaya API compatibility, introduces distributed runtime before single-process behavior is stabilized, or changes connection identity from server-observed id and epoch to a different correlation model.

## Follow-Up

- Define WebSocket close policy before sending close codes, close reasons, kick messages, or disconnect behavior.
- Define protocol logout route mapping before exposing logout to clients.
- Define reconnect and connection epoch behavior before duplicate replacement or resume behavior.
- Define protocol session carriers before clients receive or carry runtime session ids.
