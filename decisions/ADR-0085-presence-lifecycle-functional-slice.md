# ADR-0085: Presence Lifecycle Functional Slice

Status: Accepted
Date: 2026-05-20
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-20-define-presence-lifecycle-functional-slice/`

Related conversations:

- `conversations/2026-05-20-presence-lifecycle-functional-slice.md`

Related artifacts:

- `runtime/internal/app/connection/registry.go`
- `runtime/internal/platform/transport/ws/server.go`
- `runtime/cmd/vibit-server/main.go`
- `runtime/cmd/vibit-server/connection_lifecycle.go`
- `runtime/cmd/vibit-server/connection_binding_registry.go`
- `runtime/internal/app/connection/registry_test.go`
- `runtime/internal/platform/transport/ws/server_test.go`
- `runtime/cmd/vibit-server/connection_lifecycle_test.go`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `rules/check-rules.json`

## Context

`ADR-0084` made successful login responses carry the server-created runtime session id through existing `Envelope.Session` metadata. The next product-parity gap was presence lifecycle.

Nakama-class servers expose presence/status as a basic online-service capability. Pitaya-class Go game server architecture keeps acceptors, sessions, route handlers, and connection lifecycle responsibilities separated. vibit already had a single-process active connection registry, first-message connection binding, close policy, transport close handoff, and connection epoch progression, but those pieces were not wired into a first presence primitive.

`ADR-0082` reduced confirmation-gate density for non-security functional slices. Presence lifecycle is a Tier 2 functional slice: it should move forward through a bounded change spec plus implementation, not another pure confirmation milestone.

## Decision

Select:

```text
define_presence_lifecycle_functional_slice
```

as a Tier 2 functional slice and implement the smallest checkable server-owned presence lifecycle behavior directly.

The implementation:

- adds `PresenceForPlayer` to the application-owned connection registry,
- derives `online` from active bound connection records,
- derives `offline` when no active bound player connection exists,
- adds a credential-neutral WebSocket `ConnectionLifecycleObserver`,
- wires PostgreSQL startup so WebSocket open/close registers and closes server-observed connection records,
- wraps successful first-message connection binding so validated player identity is bound into the registry,
- keeps `RequestIdentity` token-record and runtime-session validation expansion deferred.

## Boundaries

This ADR keeps these boundaries:

- Presence lifecycle state is application-owned under `runtime/internal/app/connection`.
- WebSocket transport owns only server-observed open/close lifecycle observation.
- Startup composition adapts transport lifecycle and validated binding into the registry.
- Authentication service remains token lifecycle service logic and does not own presence.
- WebSocket transport remains credential-neutral.

This decision does not add Protobuf presence messages, generated output, protocol presence query routes, presence subscriptions, presence broadcasts, chat, friends, groups, parties, matchmaking, match runtime, operations/admin behavior, reconnect tokens, resume tokens, logout-triggered socket close, runtime session revocation, durable/distributed presence, dependencies, broad product modules, or direct Nakama/Pitaya API compatibility.

The current startup binding records player/connection presence. It does not claim full runtime-session or access-token-record presence linkage from access-token validation because the current `RequestIdentity` does not carry those values.

## Nakama And Pitaya Mapping

Nakama informs the product pressure: player presence must be grounded in authenticated realtime socket lifecycle state before chat, notifications, parties, matchmaking, and match runtime can build on it.

Pitaya informs the layering pressure: acceptor lifecycle, session/binding context, handler routing, and connection management should remain separate surfaces.

vibit adapts those lessons by deriving presence from its application-owned active connection registry while keeping WebSocket transport credential-neutral and protocol shape unchanged.

## Alternatives Considered

- Add a Protobuf `GetPresence` query immediately.
- Add presence subscriptions or broadcasts immediately.
- Add durable or distributed presence storage.
- Put presence directly in WebSocket transport.
- Make authentication service own active connection presence.
- Use client-supplied `player_id` or `session_id` metadata as presence authority.
- Jump directly to chat, friends, parties, matchmaking, or match runtime.

## Rationale

The smallest useful presence lifecycle primitive is not a public route; it is reliable server-owned state that says whether a validated player has active bound connections.

Without registry-backed presence, a protocol query would either return fake data or trust client metadata. With this slice, the server can track open, bound, and closed connection lifecycle state in normal PostgreSQL startup composition while preserving current security boundaries.

This is enough to unblock a later protected presence query slice and keeps the project moving toward Nakama/Pitaya-class coverage without opening broad social or multiplayer modules.

## Agent Reasoning Summary

The maintainer asked for faster progress and fewer confirmation gates. The useful compromise is a small production-shaped lifecycle primitive: wire the existing registry to real WebSocket lifecycle and validated connection binding, then expose the resulting presence snapshot in a future protocol query.

## Decision Weights

```yaml
decision_weights:
  development_velocity: high
  lifecycle_closure: high
  product_parity_progress: high
  transport_auth_separation: high
  protocol_shape_stability: high
  durable_distributed_scope: low
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `runtime.presence_lifecycle_functional_slice` becomes the repository check rule for this slice.
- The runtime can derive player online/offline state from active bound connection records.
- WebSocket open and close lifecycle can update the active connection registry through startup composition.
- Successful first-message binding can bind validated player identity into the registry.
- Future protocol presence query work can return real server-owned presence snapshots instead of placeholder state.

## Reversal Conditions

Revisit this decision if future distributed runtime work replaces single-process registry semantics, if presence subscriptions require a different source of truth, or if direct Nakama/Pitaya API compatibility is explicitly selected and requires a different presence model.

## Follow-Up

- Add a protected protocol presence query functional slice.
- Keep subscriptions, broadcasts, chat, social modules, matchmaking, match runtime, operations/admin behavior, durable/distributed presence, reconnect/resume tokens, logout-triggered close, runtime session revocation, and direct compatibility behind explicit future work.

