# ADR-0073: Logout Access Token Behavior Implementation

Status: Accepted
Date: 2026-05-18
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-logout-access-token-behavior-gate/`
- `changes/2026-05-18-implement-logout-access-token-behavior/`

Related conversations:

- `conversations/2026-05-18-logout-access-token-behavior-implementation.md`

Related artifacts:

- `runtime/internal/app/authentication/service.go`
- `runtime/internal/app/authentication/service_test.go`
- `docs/logout-access-token-behavior-gate.md`
- `decisions/ADR-0072-logout-access-token-behavior-gate.md`
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

`ADR-0072` defined the future first logout posture: revoke only the verified presented opaque access-token record through the application authentication service. The work queue then blocked at `M-085/W-0157`, requiring a maintainer direction before code could change.

The maintainer asked the agent in Chinese to recommend the next ten steps and continue, with Nakama and Pitaya as key references.

The smallest useful implementation now is the service-level `LogoutAccessToken` behavior. The repository already has token lookup and revocation mutation vocabulary, and the application package already has proof decoding, lookup digest computation, verifier digest computation, constant-time comparison, token posture checks, redacted service errors, and unit-of-work composition.

## Decision

Select:

```text
implement_logout_access_token_behavior
```

Implement `Service.LogoutAccessToken` in `runtime/internal/app/authentication` with this sequence:

- Reject missing or malformed access-token proof before opening a unit of work.
- Decode the access token as Base64URL unpadded 32-byte opaque material.
- Compute token lookup digest before repository lookup.
- Open the application unit of work and obtain the authentication repository.
- Find the token record by lookup digest.
- Require the existing active access-token posture.
- Compute and compare the token verifier digest before revocation.
- Revoke the presented token record with reason `logout_presented_access_token`.
- Return success only after unit-of-work commit.

Public invalid-token behavior stays collapsed for lookup miss, revoked, expired, wrong kind, wrong audience, unsupported verifier metadata, unknown key id, and verifier mismatch. Repository and commit failures return token-store-unavailable behavior without claiming revocation success.

This ADR does not revoke runtime sessions, close WebSocket connections, add a connection registry, add kick/disconnect behavior, add reconnect or epoch behavior, add Protobuf logout routes, add protocol session carriers, change the existing Protobuf envelope, change WebSocket handshake authentication, add transport credential carriers, add refresh, add logout-all, add admin revocation, add cleanup jobs, add dependencies, add memory durable session behavior, broaden game backend modules, or adopt direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Keep `LogoutAccessToken` fail-closed.
- Return idempotent success for already revoked tokens.
- Revoke runtime sessions together with the presented token.
- Close active WebSocket connections from the logout service.
- Add a protocol logout route in the same slice.
- Add logout-all sessions or admin revocation.
- Copy Nakama or Pitaya public logout APIs directly.

## Rationale

Nakama shows that authenticated token/session material has lifecycle state and that revoked material must not authorize future gameplay requests. vibit adapts that by making token revocation explicit and transactional, while keeping refresh, logout-all, admin revocation, and realtime socket invalidation as separate surfaces.

Pitaya shows the value of separating connection/session infrastructure from handler logic. vibit adapts that by keeping logout in the application authentication service and leaving WebSocket transport, protocol adapters, session repositories, and connection management out of this slice.

The presented-token-only implementation is the smallest behavior that converts the existing logout contract from fail-closed placeholder to useful lifecycle behavior without hiding broader session or socket decisions.

## Agent Reasoning Summary

The previous gate already answered the hard ordering questions for service-level logout. Implementing exactly that sequence improves security posture and future gameplay correctness while keeping the design bounded. It also creates a concrete base for later protocol logout routes or active connection invalidation gates.

## Decision Weights

```yaml
decision_weights:
  revocation_correctness: high
  transaction_boundary_clarity: high
  public_error_redaction: high
  nakama_pitaya_alignment: high
  transport_protocol_app_separation: high
  implementation_scope_control: high
  immediate_protocol_surface: low
  direct_nakama_pitaya_api_compatibility: low
confidence: high
```

## Consequences

- `runtime.logout_access_token_behavior_implementation` becomes the repository check rule for this slice.
- `LogoutAccessToken` now revokes verified presented opaque access-token records.
- Existing request-level access-token validation will reject revoked tokens on later protected requests when validation is invoked.
- Runtime session revocation and active WebSocket connection invalidation remain later gates.
- Protobuf logout route exposure remains a later gate.
- The work queue blocks again after this implementation at `M-087/W-0159`.

## Reversal Conditions

Revisit this decision if a future ADR chooses idempotent success for already revoked tokens, requires linked runtime session revocation in the same transaction, adopts handshake-level authentication as the primary proof carrier, introduces a connection registry before logout execution, requires logout-all sessions as the first public posture, or adopts direct Nakama/Pitaya API compatibility.

## Follow-Up

- Define active connection registry behavior before targeting open sockets.
- Define protocol logout route mapping before exposing logout to clients.
- Define runtime session revocation policy before invalidating bound-session routes after logout.
- Define reconnect and epoch behavior before duplicate replacement or resume behavior.
