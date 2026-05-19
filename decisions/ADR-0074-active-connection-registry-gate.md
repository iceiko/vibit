# ADR-0074: Active Connection Registry Gate

Status: Accepted
Date: 2026-05-18
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-logout-access-token-behavior-implementation/`
- `changes/2026-05-18-define-active-connection-registry-gate/`

Related conversations:

- `conversations/2026-05-18-active-connection-registry-gate.md`

Related artifacts:

- `docs/active-connection-registry-gate.md`
- `docs/active-connection-registry-gate.zh-CN.md`
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

`ADR-0073` implemented presented access-token logout behavior. That behavior revokes the verified presented opaque access-token record, but it deliberately does not revoke runtime sessions, close active WebSocket connections, add a registry, or expose a protocol logout route.

The work queue reached `M-087/W-0159`, a confirmation gate. The maintainer asked the agent, in Chinese, to recommend the next ten steps and continue, with Nakama and Pitaya as key references.

The next highest-value boundary is the active connection registry. Without it, future logout, revocation, kick/disconnect, duplicate replacement, reconnect, and route policy decisions would have no safe server-owned model for active sockets.

## Decision

Select:

```text
define_active_connection_registry_gate
```

Create a gate-only standard:

```text
docs/active-connection-registry-gate.md
docs/active-connection-registry-gate.zh-CN.md
```

The gate defines future ownership, first single-process in-memory posture, candidate registry record vocabulary, candidate capabilities, redaction rules, WebSocket/Protobuf deferrals, future test expectations, and Nakama/Pitaya reference mapping.

The first future registry posture is conservative:

- Registry behavior is application-owned, with candidate package `runtime/internal/app/connection`.
- The first registry is single-process, in-memory, and non-durable.
- Registry records model server-observed active connection state and validated identity linkage, not client proof.
- WebSocket transport remains credential-neutral and does not own authentication state.
- Registry lookup, registry invalidation, and concrete socket close are separate future policy choices.
- Cluster routing, distributed kick/disconnect, reconnect/epoch behavior, protocol logout routes, protocol session carriers, and direct Nakama/Pitaya API compatibility remain deferred.

This ADR does not implement a registry, close WebSocket connections, add kick/disconnect behavior, revoke runtime sessions, change route policy, add reconnect or epoch behavior, add Protobuf logout routes, add protocol session carriers, change the existing Protobuf envelope, change WebSocket handshake authentication, add transport credential carriers, add durable connection storage, add dependencies, or adopt direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Expose a protocol logout route before defining active connection targeting.
- Implement best-effort WebSocket close directly from `LogoutAccessToken`.
- Treat first-message connection binding state as a complete active connection registry.
- Add reconnect and epoch behavior before registry ownership is defined.
- Add a distributed registry or Redis-like dependency immediately.
- Copy Nakama or Pitaya public APIs directly.

## Rationale

Nakama shows that authenticated session material and realtime socket lifecycle need coordinated policy. vibit adapts that by making active connections first-class application runtime state before logout or revocation can target open sockets.

Pitaya shows the value of separating acceptors, sessions, route handlers, and connection management. vibit adapts that by keeping registry state in the application layer and limiting transport participation to future narrow lifecycle handoffs.

A gate is needed before implementation because a registry sits between connection binding, route policy, logout/revocation, WebSocket transport lifecycle, future reconnect behavior, and future operations tooling. Defining the boundary first prevents hidden transport-owned authentication state and metadata-only socket targeting.

## Agent Reasoning Summary

After service-level logout implementation, the next practical issue is not the public logout route. It is whether the server can safely identify active sockets affected by future lifecycle changes. A registry boundary is the prerequisite for targeted invalidation, duplicate replacement, reconnect, presence, and multiplayer routing.

Choosing a gate-only registry boundary keeps the work aligned with Nakama and Pitaya without importing their APIs or distributed assumptions too early.

## Decision Weights

```yaml
decision_weights:
  active_connection_lifecycle_clarity: high
  nakama_pitaya_alignment: high
  transport_protocol_app_separation: high
  future_logout_revocation_correctness: high
  reconnect_and_duplicate_replacement_readiness: high
  implementation_scope_control: high
  immediate_user_visible_feature_surface: low
  direct_nakama_pitaya_api_compatibility: low
confidence: high
```

## Consequences

- `runtime.active_connection_registry_gate` becomes the repository check rule for this boundary.
- Future active socket targeting requires a registry implementation gate.
- Future WebSocket close/kick/disconnect behavior requires a separate close policy gate or implementation authorization.
- `LogoutAccessToken` remains token-record scoped and does not close sockets.
- WebSocket transport remains credential-neutral and does not own authentication state.
- The work queue blocks again after the gate at `M-089/W-0161`.

## Reversal Conditions

Revisit this decision if a future ADR chooses handshake-level authentication as the primary identity source, requires transport-owned auth state, adopts direct Nakama or Pitaya compatibility, introduces distributed runtime before single-process registry behavior, or selects a protocol-first logout route that does not need active socket targeting.

## Follow-Up

- Implement a single-process in-memory active connection registry after explicit confirmation.
- Define WebSocket close policy before sending custom close reasons or close codes.
- Define reconnect and epoch behavior before duplicate replacement or resume behavior.
- Define protocol logout/session carriers before exposing logout, session ids, resume tokens, or connection epochs to clients.
