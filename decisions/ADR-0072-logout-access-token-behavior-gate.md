# ADR-0072: Logout Access Token Behavior Gate

Status: Accepted
Date: 2026-05-18
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-logout-revocation-active-connection-gate/`
- `changes/2026-05-18-define-logout-access-token-behavior-gate/`

Related conversations:

- `conversations/2026-05-18-logout-access-token-behavior-gate.md`

Related artifacts:

- `docs/logout-access-token-behavior-gate.md`
- `docs/logout-access-token-behavior-gate.zh-CN.md`
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

`W-0154` defined the logout/revocation active-connection gate. That gate separated presented-token logout, runtime session revocation, and active socket invalidation as future decisions, and recommended `presented_access_token_only` as the first future logout scope.

The work queue reached `M-083/W-0155`, a confirmation gate. The maintainer asked the agent to recommend the next ten steps and continue, with Nakama and Pitaya as key references.

The safest next direction is to define the exact behavior boundary for future `LogoutAccessToken` execution before writing revocation code. vibit already has access-token validation, token repository mutation vocabulary, durable token records, and a fail-closed `LogoutAccessToken` service method. It still lacks the rules for proof validation, verifier comparison, transaction ordering, public error collapse, and session/socket deferral when revoking the presented token.

## Decision

Select:

```text
define_logout_access_token_behavior_gate
```

Create a gate-only standard:

```text
docs/logout-access-token-behavior-gate.md
docs/logout-access-token-behavior-gate.zh-CN.md
```

The gate defines the future first logout posture:

- Revoke only the verified presented opaque access token.
- Reject missing or malformed token proof before opening a unit of work.
- Use lookup digest and verifier digest helpers before revocation.
- Require token kind, active status, non-expiry, expected audience, supported verifier metadata, and constant-time verifier comparison before revocation.
- Return logout success only after unit-of-work commit.
- Collapse public failures to missing, malformed, invalid, or unavailable token errors.
- Keep runtime session revocation, active connection invalidation, connection registry, WebSocket close policy, Protobuf logout route, protocol session carriers, refresh, logout-all, admin revocation, cleanup jobs, dependencies, and direct Nakama/Pitaya API compatibility deferred.

This ADR does not implement `LogoutAccessToken`, revoke tokens, revoke runtime sessions, close WebSocket connections, add connection registries, add kick/disconnect behavior, add reconnect or epoch behavior, add Protobuf logout routes, add protocol session carriers, change the existing Protobuf envelope, change WebSocket handshake authentication, add transport credential carriers, add cleanup jobs, add dependencies, add memory durable session behavior, or adopt direct Nakama/Pitaya public API compatibility.

## Alternatives Considered

- Implement `LogoutAccessToken` immediately.
- Treat logout as idempotent success when the token is already revoked.
- Revoke the linked runtime session in the first logout transaction.
- Close active WebSocket connections from logout execution.
- Add a Protobuf logout route before the service behavior is defined.
- Support logout-all sessions as the first posture.
- Add admin revocation as part of player logout.
- Copy Nakama or Pitaya public APIs directly.

## Rationale

Nakama shows that authenticated session material needs explicit lifecycle behavior: logout, refresh, expiration, revocation, and realtime implications all affect gameplay access. vibit adapts that pressure by requiring revoked access tokens to stop authorizing future protected requests while keeping refresh, logout-all, admin revocation, and realtime socket invalidation as separate surfaces.

Pitaya shows the value of keeping session/connection infrastructure separate from handler logic. vibit adapts that by placing logout behavior in the application authentication service and keeping WebSocket transport, protocol adapters, connection registries, and handler packages out of credential parsing and revocation decisions.

The gate is necessary because logout execution crosses cryptographic proof validation, repository mutation, transaction outcome, public error semantics, and future session/socket lifecycle. Defining the behavior before implementation reduces the chance of a hidden logout-all, hidden socket close, or repository-owned credential decision.

## Agent Reasoning Summary

After the active-connection revocation gate, the next practical step is to make the smallest executable logout behavior precise. The selected scope is presented-token logout only. That turns the existing semantic `LogoutAccessToken` contract into an implementable future slice without prematurely taking on connection registry, socket close, protocol route, reconnect, or broader game-server module behavior.

This preserves Nakama's session lifecycle pressure and Pitaya's separation of connection/session/handler concerns while remaining vibit-native and contract-first.

## Decision Weights

```yaml
decision_weights:
  security_and_revocation_correctness: high
  transaction_boundary_clarity: high
  public_error_redaction: high
  nakama_pitaya_alignment: high
  transport_protocol_app_separation: high
  implementation_scope_control: high
  immediate_user_visible_feature_surface: low
  direct_nakama_pitaya_api_compatibility: low
confidence: high
```

## Consequences

- `runtime.logout_access_token_behavior_gate` becomes the repository check rule for this boundary.
- Future logout implementation should be `implement_logout_access_token_behavior`.
- The first implementation must revoke only the verified presented token record and return success only after commit.
- Runtime session revocation and active WebSocket connection invalidation remain later gates.
- Protobuf logout route exposure remains a later gate.
- The work queue blocks again after this gate at `M-085/W-0157`.

## Reversal Conditions

Revisit this decision if a future ADR selects idempotent logout success for already revoked tokens, requires linked runtime session revocation in the first logout transaction, adopts handshake-level authentication as the primary proof carrier, introduces a connection registry before logout execution, requires logout-all sessions as the first public posture, or adopts a direct Nakama or Pitaya compatibility surface.

## Follow-Up

- Implement `LogoutAccessToken` behavior in `runtime/internal/app/authentication` after a bounded implementation slice is selected.
- Define active connection registry behavior before targeting open sockets.
- Define transport close policy before using custom WebSocket close codes or reason text.
- Define protocol logout route mapping before exposing logout to clients.
- Define reconnect and epoch behavior before duplicate replacement or resume behavior.
