# ADR-0080: Transport Close Handoff Gate

Status: Accepted
Date: 2026-05-18
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-protocol-logout-route-implementation/`
- `changes/2026-05-18-define-transport-close-handoff-gate/`

Related conversations:

- `conversations/2026-05-18-protocol-logout-route-next-direction.md`
- `conversations/2026-05-18-transport-close-handoff-gate.md`

Related artifacts:

- `docs/transport-close-handoff-gate.md`
- `docs/transport-close-handoff-gate.zh-CN.md`
- `docs/websocket-close-policy-gate.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Context

The Nakama/Pitaya-class product roadmap keeps runtime lifecycle closure as the near-term phase. `W-0170` exposed protocol logout without closing sockets. `W-0166` implemented an application-owned close policy that can invalidate active bound registry records and emit redacted close intents, but it deliberately uses `mark_invalidated_only` and has no concrete WebSocket socket close handoff.

This leaves a clear gap: the system can decide that a connection should no longer be valid, but WebSocket transport has no narrow policy-neutral way to close the concrete socket. Future presence, chat, matchmaking, match runtime, operations, and distributed runtime work should not depend on this ambiguity.

Nakama informs the product pressure: logout, session lifecycle, realtime socket disconnect, and server-directed disconnect are related but distinct lifecycle behaviors. Pitaya informs the architecture pressure: acceptors, sessions, handlers, and kick/disconnect connection-management mechanics should be distinct surfaces.

## Decision

Select:

```text
define_transport_close_handoff_gate
```

Create the gate-only standard:

```text
docs/transport-close-handoff-gate.md
docs/transport-close-handoff-gate.zh-CN.md
```

Define the repository check rule:

```text
runtime.transport_close_handoff_gate
```

The gate establishes that a future transport close handoff must be narrow and policy-neutral:

- Application close policy remains the owner of close decisions.
- Active connection registry remains the owner of server-observed connection records and lifecycle markers.
- WebSocket transport may later own only concrete socket close mechanics.
- The first handoff target is server-observed `connection_id + connection_epoch`.
- Transport must not close by player id, runtime session id, access-token record id, route identity, request identity, envelope metadata, headers, cookies, query strings, subprotocol values, or remote address.
- Connection epoch must prevent stale close intents from closing a later socket.
- Close code mapping, close reason text, logout-triggered close, runtime session revocation, reconnect/epoch behavior, protocol session carriers, operations/admin disconnect, dependencies, and direct Nakama/Pitaya API compatibility remain deferred.

This ADR does not implement concrete WebSocket close handoff.

## Alternatives Considered

- Implement socket close directly in WebSocket transport without a gate.
- Let active connection registry close sockets whenever records are invalidated.
- Let protocol logout close the current socket after successful token revocation.
- Close sockets by player id or session id inside transport.
- Add reconnect and duplicate replacement before concrete close handoff.
- Add presence lifecycle before concrete close handoff.
- Copy Nakama or Pitaya disconnect APIs directly.

## Rationale

The application layer has the context to decide why a socket should close. Transport has the concrete socket handle. Conflating those responsibilities would make authentication, logout, session revocation, reconnect, and operations behavior harder for agents to reason about.

The first target must be `connection_id + epoch` because that is the smallest server-owned target that can identify one concrete accepted socket and reject stale intents. Higher-level targets such as player id, runtime session id, and token record id are valid application-policy inputs, but they should be resolved to concrete connection/epoch targets before transport sees them.

This preserves the Nakama/Pitaya lessons without copying public APIs. It gives vibit a route to server-directed close behavior while keeping transport credential-neutral and policy-neutral.

## Agent Reasoning Summary

After logout route exposure, the next useful lifecycle step is not presence, chat, groups, or matchmaking. Those features depend on predictable connection close and reconnect behavior. The handoff must be defined before it is implemented because a careless implementation would put policy in the transport layer or use client metadata as authority.

## Decision Weights

```yaml
decision_weights:
  lifecycle_closure: high
  transport_policy_separation: high
  stale_connection_safety: high
  nakama_pitaya_alignment: high
  agent_context_reduction: high
  implementation_scope_control: high
  immediate_socket_close_feature_surface: low
  close_reason_text_now: low
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `runtime.transport_close_handoff_gate` becomes the repository check rule for this boundary.
- Future concrete socket close implementation requires a later bounded implementation slice.
- `LogoutAccessToken` remains token-record scoped and does not close sockets.
- WebSocket transport remains credential-neutral and policy-neutral.
- Reconnect/epoch, protocol session carriers, presence lifecycle, operations/admin disconnect, and distributed close behavior remain later gates.

## Reversal Conditions

Revisit this decision if a later ADR selects a different connection identity model, if WebSocket transport becomes stateful through a ratified connection manager with stronger ownership, if direct Nakama/Pitaya API compatibility is explicitly selected, or if distributed runtime becomes a prerequisite before single-process transport close handoff.

## Follow-Up

- Implement a bounded single-process transport close handoff after this gate.
- Define reconnect and connection epoch behavior after concrete close mechanics are available.
- Define protocol session carriers after close/reconnect semantics are stable.
- Define presence lifecycle after connection lifecycle behavior is less ambiguous.
