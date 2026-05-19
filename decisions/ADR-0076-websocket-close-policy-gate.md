# ADR-0076: WebSocket Close Policy Gate

Status: Accepted
Date: 2026-05-18
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-active-connection-registry-implementation/`
- `changes/2026-05-18-define-websocket-close-policy-gate/`

Related conversations:

- `conversations/2026-05-18-websocket-close-policy-gate.md`

Related artifacts:

- `docs/websocket-close-policy-gate.md`
- `docs/websocket-close-policy-gate.zh-CN.md`
- `docs/active-connection-registry-gate.md`
- `decisions/ADR-0075-active-connection-registry-single-process-implementation.md`
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

`ADR-0075` implemented the first application-owned, single-process active connection registry. That registry can represent server-observed connection ids, connection epochs, validated player linkage, optional runtime session ids, optional access-token record ids, and policy-neutral invalidation markers.

The registry deliberately does not close WebSocket connections, kick users, choose close codes or close reason text, revoke runtime sessions, expose protocol logout routes, add reconnect behavior, or wire transport handoffs.

The work queue reached `M-091/W-0163`, a confirmation gate. The maintainer asked the agent, in Chinese, to recommend the next ten steps and continue, with Nakama and Pitaya as key references.

The next highest-value boundary is WebSocket close policy. Without it, future logout, token revocation, runtime session revocation, admin disconnect, duplicate connection policy, or operational drain could accidentally close sockets from the wrong layer or make hidden failure semantics.

## Decision

Select:

```text
define_websocket_close_policy_gate
```

Create a gate-only standard:

```text
docs/websocket-close-policy-gate.md
docs/websocket-close-policy-gate.zh-CN.md
```

The gate defines future ownership, close intent vocabulary, target vocabulary, reason classes, redaction requirements, close failure questions, registry interaction rules, logout/session revocation boundaries, WebSocket/protocol deferrals, test expectations, and Nakama/Pitaya reference mapping.

The first future close policy posture is conservative:

- WebSocket close policy is application-owned under `runtime/internal/app`.
- The active connection registry remains a target model under `runtime/internal/app/connection`, not the policy owner.
- WebSocket transport may own only a future narrow close handoff after application policy emits a redacted close intent.
- Authentication service behavior may revoke tokens, but it must not directly close sockets.
- Protocol adapters and domain modules must not hide close, kick, or disconnect behavior inside route handling.
- Registry invalidation and concrete socket close remain separate until a later implementation gate explicitly wires them.

Nakama guides the lifecycle distinction: its session documentation separates token logout from open socket disconnection, and its server function reference exposes session disconnection as a distinct server-side operation. vibit adapts that lesson by requiring explicit close policy instead of making logout implicitly close sockets.

Pitaya guides the layering distinction: its documentation separates acceptors, handler service lifecycle, sessions, and protocol-agnostic session operations such as kicking connected users. vibit adapts that lesson by keeping the close decision in application policy and limiting transport to future concrete close mechanics.

This ADR does not implement WebSocket close behavior, add transport close handoff code, choose close codes or reason strings, add kick/disconnect behavior, change logout behavior, revoke runtime sessions, replace duplicate connections, add reconnect or epoch behavior, add Protobuf logout routes, add protocol session carriers, change the existing Protobuf envelope, change WebSocket handshake authentication, add transport credential carriers, add durable or distributed registry storage, add dependencies, broaden game backend modules, or adopt direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement logout-triggered socket close directly from `LogoutAccessToken`.
- Let the active connection registry close sockets when records are invalidated.
- Put close code and reason mapping directly in WebSocket transport.
- Add protocol logout route before close policy exists.
- Add reconnect and connection epoch behavior before close semantics are defined.
- Add admin kick/disconnect behavior as an operations feature first.
- Copy Nakama or Pitaya public APIs directly.

## Rationale

Nakama shows that game backends need predictable lifecycle semantics around session tokens, realtime sockets, logout, and server-directed disconnect. The important lesson is not API compatibility; it is that token invalidation and socket disconnection are distinct lifecycle actions with different player-facing effects.

Pitaya shows that accepting a connection, representing a session, running a route handler, and managing connection lifecycle are distinct responsibilities. The important lesson is not to copy Pitaya's session API; it is to keep business policy out of low-level network acceptors and avoid hidden route-handler side effects.

vibit needs a gate before implementation because the active connection registry now makes targeting possible, but targeting is not policy. A close policy must answer who chooses the target, which reason class applies, whether a close is silent or visible, how failures affect the caller, and which layer performs the concrete transport close.

## Agent Reasoning Summary

After the active connection registry implementation, the next practical issue is not exposing a public logout route. The server first needs a durable rule that says a registry target cannot be converted into a socket close unless application policy explicitly authorizes it.

This keeps future logout, revocation, reconnect, admin operations, presence, rooms, and match runtime work from smuggling close semantics into transport or domain modules.

## Decision Weights

```yaml
decision_weights:
  realtime_lifecycle_clarity: high
  nakama_pitaya_alignment: high
  transport_protocol_app_separation: high
  future_logout_revocation_correctness: high
  close_failure_semantics_clarity: high
  implementation_scope_control: high
  immediate_socket_close_feature_surface: low
  direct_nakama_pitaya_api_compatibility: low
confidence: high
```

## Consequences

- `runtime.websocket_close_policy_gate` becomes the repository check rule for this boundary.
- Future concrete WebSocket close behavior requires a separate implementation gate.
- Future logout-triggered socket close, session revocation close, duplicate replacement, reconnect/epoch behavior, admin kick/disconnect, protocol close messages, close codes, and close reasons remain deferred.
- `LogoutAccessToken` remains token-record scoped and does not close sockets.
- The active connection registry remains policy-neutral.
- WebSocket transport remains credential-neutral and policy-neutral.
- The work queue blocks again after the gate at `M-093/W-0165`.

## Reversal Conditions

Revisit this decision if a future ADR chooses transport-owned authentication state, makes handshake-level authentication the sole lifecycle owner, adopts direct Nakama or Pitaya compatibility, introduces distributed runtime before single-process close policy is proven, or selects a public protocol close-message model that requires a different application/transport handoff.

## Follow-Up

- Implement a single-process WebSocket close policy after explicit confirmation.
- Define protocol logout route mapping before exposing logout to clients.
- Define reconnect and connection epoch behavior before duplicate replacement or resume behavior.
- Define protocol session carriers before clients receive or carry runtime session ids.
- Define operations/admin surfaces before adding administrative kick or disconnect behavior.
