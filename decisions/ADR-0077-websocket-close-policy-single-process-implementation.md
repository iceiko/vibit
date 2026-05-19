# ADR-0077: WebSocket Close Policy Single Process Implementation

Status: Accepted
Date: 2026-05-18
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-websocket-close-policy-gate/`
- `changes/2026-05-18-implement-websocket-close-policy-single-process/`

Related conversations:

- `conversations/2026-05-18-websocket-close-policy-single-process-implementation.md`

Related artifacts:

- `docs/websocket-close-policy-gate.md`
- `decisions/ADR-0076-websocket-close-policy-gate.md`
- `runtime/internal/app/connection/registry.go`
- `runtime/internal/app/connection/registry_test.go`
- `runtime/internal/app/connection/close_policy.go`
- `runtime/internal/app/connection/close_policy_test.go`
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

`ADR-0076` defined the WebSocket close policy gate. The gate established that future close policy is application-owned, that the active connection registry is target state rather than policy, and that WebSocket transport may only perform a later narrow concrete close handoff after application policy emits a redacted close intent.

The work queue reached `M-093/W-0165`, a confirmation point. The maintainer asked, in Chinese, for the next recommended ten steps and explicitly emphasized Nakama and Pitaya as reference baselines.

The next bounded step is to implement the first single-process close policy primitive. That primitive is useful now because the active connection registry can already represent active bound connections and policy-neutral invalidation. The implementation must still avoid concrete WebSocket close mechanics.

## Decision

Select:

```text
implement_websocket_close_policy_single_process
```

Add:

```text
runtime/internal/app/connection/close_policy.go
runtime/internal/app/connection/close_policy_test.go
```

The implementation is application-owned under `runtime/internal/app/connection`.

The policy supports these target kinds:

- `connection_id_and_epoch`
- `player_id`
- `runtime_session_id`
- `access_token_record_id`

The policy resolves targets only through server-owned active bound registry records. It does not synthesize close intents from client metadata alone. Missing, unbound, closed, invalidated, or otherwise inactive records produce no close intent.

The policy emits redacted `CloseIntent` values with internal reason class, target kind, transport action, retryability, public visibility, outcome, and registry identity linkage. The only transport action selected in this implementation is:

```text
mark_invalidated_only
```

Matched records are marked invalidated through `MarkConnectionInvalidated`. That invalidation is an application lifecycle marker, not a concrete socket close.

The policy accepts the internal reason classes ratified by the gate:

- `token_revoked`
- `logout_presented_token`
- `session_revoked`
- `duplicate_connection_policy`
- `server_shutdown_or_drain`
- `policy_violation`
- `administrative_action`
- `protocol_error`
- `idle_timeout`
- `unknown_internal`

This ADR does not add concrete WebSocket close handoff, close codes, close reason text, kick/disconnect public API behavior, logout-triggered protocol route behavior, runtime session revocation close behavior, duplicate replacement, reconnect or resume behavior, protocol session carriers, generated Protobuf output, WebSocket handshake authentication, transport credential carriers, durable or distributed registry storage, cleanup jobs, operations/admin disconnect, dependencies, broader game backend behavior, or direct Nakama/Pitaya API compatibility.

## Nakama And Pitaya Mapping

Nakama informs the lifecycle model: token/session lifecycle and realtime socket lifecycle need explicit server-side semantics, and server-directed disconnect is distinct from token invalidation. vibit adapts this by producing close intents and invalidation markers without making logout automatically close sockets.

Pitaya informs the layering model: acceptors, sessions, route handlers, and connection management are separate surfaces. vibit adapts this by keeping close policy in the application layer and preserving concrete transport close mechanics for a future narrow handoff.

This decision does not copy either system's public API.

## Alternatives Considered

- Wire `LogoutAccessToken` directly to active socket close.
- Let the registry close sockets when invalidation is recorded.
- Put close policy and close codes in WebSocket transport.
- Expose a protocol logout route before close intent behavior exists.
- Implement a transport close handoff in the same slice.
- Add reconnect and duplicate connection replacement before close policy exists.
- Copy Nakama or Pitaya session disconnect APIs directly.

## Rationale

The registry made active connection targeting possible, but targeting is not policy. A close policy primitive gives future logout, token revocation, session revocation, duplicate connection handling, operations drain, and admin disconnect work a shared application-owned vocabulary.

The `mark_invalidated_only` first posture is intentionally conservative. It proves the policy and state transition semantics while preserving the harder transport decision for a later gate. That keeps the first implementation small, testable, and compatible with the existing WebSocket transport boundary.

## Agent Reasoning Summary

After `ADR-0076`, the most useful next step is not a user-facing logout route. The server first needs an internal application primitive that can turn trusted registry targets into redacted lifecycle intents without leaking token/session material or letting transport decide application policy.

This follows Nakama's product pressure for predictable realtime lifecycle behavior and Pitaya's architectural separation between connection management and route handlers.

## Decision Weights

```yaml
decision_weights:
  application_owned_lifecycle_policy: high
  nakama_pitaya_alignment: high
  transport_protocol_authentication_separation: high
  future_logout_revocation_composability: high
  redaction_and_agent_verifiability: high
  immediate_concrete_socket_close_feature_surface: low
  direct_nakama_pitaya_api_compatibility: low
confidence: high
```

## Consequences

- `runtime.websocket_close_policy_single_process_implementation` becomes the repository check rule for this slice.
- `runtime/internal/app/connection/close_policy.go` owns the first close policy primitive.
- Focused tests verify registry-backed target resolution, invalidation, redacted intents, skipped targets, and deferrals.
- WebSocket transport remains policy-neutral and does not close sockets for this slice.
- Authentication service behavior remains token lifecycle scoped and does not own socket close policy.
- The work queue blocks again after this implementation at `M-095/W-0167`.

## Reversal Conditions

Revisit this decision if a future ADR selects handshake-level authentication as the sole lifecycle owner, adopts direct Nakama/Pitaya public API compatibility, introduces distributed runtime connection routing before single-process close policy is proven, or chooses a protocol close-message model that requires different application intent fields.

## Follow-Up

- Define protocol logout route behavior before exposing logout to clients.
- Define reconnect and connection epoch behavior before duplicate replacement or resume behavior.
- Define protocol session carriers before clients receive or carry runtime session ids.
- Define concrete transport close handoff before close intents close sockets.
- Define operations/admin surfaces before administrative disconnect behavior is exposed.
