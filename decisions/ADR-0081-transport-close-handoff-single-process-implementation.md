# ADR-0081: Transport Close Handoff Single Process Implementation

Status: Accepted
Date: 2026-05-19
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-transport-close-handoff-gate/`
- `changes/2026-05-19-implement-transport-close-handoff-single-process/`

Related conversations:

- `conversations/2026-05-18-transport-close-handoff-next-direction.md`
- `conversations/2026-05-19-transport-close-handoff-single-process-implementation.md`

Related artifacts:

- `docs/transport-close-handoff-gate.md`
- `decisions/ADR-0080-transport-close-handoff-gate.md`
- `runtime/internal/platform/transport/ws/close_handoff.go`
- `runtime/internal/platform/transport/ws/close_handoff_test.go`
- `runtime/internal/platform/transport/ws/server.go`
- `runtime/internal/platform/transport/ws/server_test.go`
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

`ADR-0080` defined the transport close handoff gate. The gate selected server-observed `connection_id + connection_epoch` as the first concrete WebSocket handoff target and kept close decisions application-owned under `runtime/internal/app/connection`.

Before this slice, the application close policy could invalidate active registry records and emit redacted intents, but the WebSocket transport did not retain a narrow concrete handle for closing an accepted socket. Protocol logout remained token-record scoped and did not close sockets.

The next bounded step was to add only the single-process transport mechanics needed for a future application-owned policy to request concrete close by server-owned connection metadata.

## Decision

Select:

```text
implement_transport_close_handoff_single_process
```

Add:

```text
runtime/internal/platform/transport/ws/close_handoff.go
runtime/internal/platform/transport/ws/close_handoff_test.go
```

Update:

```text
runtime/internal/platform/transport/ws/server.go
runtime/internal/platform/transport/ws/server_test.go
```

The implementation is owned by WebSocket transport under `runtime/internal/platform/transport/ws`.

`Server` now keeps a single-process in-memory table of accepted WebSocket sockets keyed by server-observed connection id. Each registered socket records the server-observed connection epoch. The public handoff is:

```text
RequestClose(context.Context, CloseHandoffRequest) CloseHandoffResult
```

The request contains only:

- `ConnectionID`
- `ConnectionEpoch`
- `RequestedAt`

The result contains only:

- `ConnectionID`
- `ConnectionEpoch`
- `Outcome`
- `ClosedAt`

The accepted outcomes are:

- `close_requested`
- `socket_not_found`
- `epoch_mismatch`
- `already_closed`
- `close_failed`

The concrete close operation uses the existing WebSocket connection handle and `CloseNow`. This deliberately avoids choosing WebSocket close codes or close reason text in this slice.

## Boundaries

The implementation keeps these boundaries:

- Application close policy remains owned by `runtime/internal/app/connection`.
- Active connection registry lifecycle markers remain owned by `runtime/internal/app/connection`.
- Authentication service behavior remains token lifecycle scoped and does not call WebSocket transport.
- Protocol adapters do not close sockets directly.
- The Protobuf envelope and authentication/logout Protobuf shapes are unchanged.
- The WebSocket transport remains credential-neutral and does not read Authorization headers, cookies, query-string tokens, subprotocol authentication material, player ids, runtime session ids, or token record ids for close authority.

This ADR does not add close code mapping, close reason text, logout-triggered socket close, runtime session revocation, reconnect or resume behavior, protocol session carriers, presence, operations/admin disconnect, dependencies, durable/distributed close handoff, or direct Nakama/Pitaya API compatibility.

## Nakama And Pitaya Mapping

Nakama informs the product need for explicit realtime socket lifecycle behavior that remains distinct from token logout and session lifecycle. vibit adapts that lesson by adding a concrete server-directed close handoff while keeping logout behavior uncoupled.

Pitaya informs the architecture split between acceptors, sessions, route handlers, and connection management. vibit adapts that lesson by keeping concrete socket mechanics in WebSocket transport while keeping application policy outside the transport package.

This decision does not copy either system's public API.

## Alternatives Considered

- Let `LogoutAccessToken` call WebSocket transport directly.
- Let the active connection registry store concrete socket handles.
- Close sockets by player id, runtime session id, access-token record id, route identity, request identity, headers, cookies, query strings, subprotocol values, or remote address.
- Select WebSocket close codes and player-visible close reason text now.
- Add reconnect and duplicate replacement behavior in the same slice.
- Add durable or distributed close handoff before single-process mechanics are proven.
- Copy Nakama or Pitaya disconnect APIs directly.

## Rationale

Connection id and epoch are the smallest server-owned target that can identify an accepted socket and protect against stale close intents. Higher-level application targets can be resolved by close policy before transport sees the handoff request.

`CloseNow` is intentionally used for the first concrete action because this slice must not select close codes or reason text. Future gates can map reason classes to close codes, player-facing text, protocol system messages, or retry semantics.

The in-memory table is intentionally single-process. It matches the current modular monolith runtime and keeps distributed routing, durable registries, operations/admin disconnect, and cluster behavior behind later decisions.

## Agent Reasoning Summary

After the gate, the useful next step was not reconnect, protocol session carriers, or presence. The system first needed a concrete transport-owned way to close an accepted socket after application policy resolves a trusted target.

The implementation stays narrow because transport has the socket handle but not the policy context. It accepts only connection id and epoch, returns redacted outcomes, and leaves all higher-level lifecycle choices to later gates.

## Decision Weights

```yaml
decision_weights:
  lifecycle_closure: high
  transport_policy_separation: high
  stale_connection_safety: high
  credential_neutrality: high
  nakama_pitaya_alignment: high
  testability_without_distributed_runtime: high
  close_code_mapping_now: low
  logout_coupling_now: low
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `runtime.transport_close_handoff_single_process_implementation` becomes the repository check rule for this slice.
- WebSocket transport has a narrow in-process concrete socket close handoff.
- Stale epoch, missing socket, already closed socket, and close failure outcomes are redacted and policy-neutral.
- Application close policy and active connection registry remain the owners of close decisions and registry lifecycle markers.
- The work queue moves to a new confirmation gate before reconnect/epoch behavior, protocol session carriers, logout-triggered socket close, close code mapping, close reason text, operations/admin disconnect, or broad product modules are selected.

## Reversal Conditions

Revisit this decision if a future ADR selects a different server-owned connection identity model, if direct Nakama/Pitaya API compatibility becomes an explicit goal, if distributed runtime routing becomes a prerequisite before single-process close handoff is useful, or if a future WebSocket close-code policy requires a different transport handoff shape.

## Follow-Up

- Confirm the next lifecycle-closure direction after this implementation.
- Define reconnect and connection epoch behavior before duplicate replacement or resume behavior.
- Define protocol session carriers after close/reconnect semantics are stable.
- Define close code mapping and close reason text only through a separate policy gate.
