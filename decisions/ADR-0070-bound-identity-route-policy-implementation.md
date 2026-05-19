# ADR-0070: Bound Identity Route Policy Implementation

Status: Accepted
Date: 2026-05-18
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-bound-identity-route-policy-gate/`
- `changes/2026-05-18-implement-bound-identity-route-policy/`

Related conversations:

- `conversations/2026-05-18-bound-identity-route-policy-implementation.md`

Related artifacts:

- `runtime/internal/app/route_authentication.go`
- `runtime/internal/app/route_authentication_test.go`
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

`W-0150` defined the bound identity route policy gate. The gate established that route policy should be application-owned, route-scoped, fail-closed, and explicit before vibit lets request-token identity, bound connection identity, or session-validated identity authorize different route families.

The work queue reached `M-079/W-0151`, a confirmation gate. The maintainer asked the agent to recommend the next ten steps and continue, with Nakama and Pitaya as key reference baselines.

At this point, vibit already had request-level access-token route protection, first-message connection binding, durable runtime session validation, and login-time durable runtime session creation. The missing implementation was route-policy vocabulary and route-scoped classification.

## Decision

Select:

```text
implement_bound_identity_route_policy
```

Implement explicit route policy vocabulary in:

```text
runtime/internal/app/route_authentication.go
runtime/internal/app/route_authentication_test.go
```

The implementation adds:

- `RouteProtectionRequirement`
- `RouteProtectionRouteRequirement`
- `RouteProtectionPublic`
- `RouteProtectionRequestTokenRequired`
- `RouteProtectionBoundConnectionRequired`
- `RouteProtectionSessionValidatedRequired`
- `RouteProtectionBoundSessionRequired`
- `RouteProtectionPolicy.RequirementFor`

The default protected route behavior remains request_token_required. `runtime.authentication.AuthenticateWithDeviceCredential` remains the explicit public route. Ordinary protected domain routes still require request-level access-token proof through the existing authenticated request wrapper.

Explicitly classified bound-connection routes may accept only a validated player identity bound to the observed connection id and epoch. Explicitly classified session-validated routes may accept only an already session-validated player identity with a matching session id. Explicitly classified bound-session routes require bound and session-validated identity sources to agree on actor/player and available session/connection metadata.

Metadata-only identity does not satisfy any protected policy family.

This ADR does not change WebSocket handshake authentication, add transport credential carriers, change Protobuf sources, change the existing Protobuf envelope, expose session ids through Protobuf, add protocol session carriers, wire persistent session validation into frame handling, add a connection binding registry, change startup route registration, remove per-request token proof from ordinary protected routes, implement logout/revocation active-connection invalidation, implement reconnect/epoch behavior, add cleanup jobs, add dependencies, add memory durable session behavior, or adopt direct Nakama/Pitaya public API compatibility.

## Alternatives Considered

- Keep the gate only and defer implementation again.
- Make first-message `BindConnection` sufficient for all ordinary protected routes.
- Make durable `session_id` metadata sufficient for protected routes.
- Require both bound connection identity and durable session validation for every protected route.
- Wire persistent session validation into the Protobuf frame handler immediately.
- Expose session ids in login responses before a protocol carrier gate.
- Implement logout/revocation active-connection behavior in the same slice.
- Copy Nakama or Pitaya public APIs directly.

## Rationale

Nakama shows that authenticated session material can become central to gameplay access, but session lifetime, refresh, logout, and active socket behavior are separate lifecycle concerns. vibit adapts this by adding explicit policy families without treating all ordinary routes as session-bound yet.

Pitaya shows a useful separation between acceptors, sessions, and route handlers. vibit adapts this by keeping policy inside `runtime/internal/app`, keeping WebSocket transport credential-neutral, and requiring handlers to receive normalized identity rather than parsing credentials or session metadata.

The conservative posture is intentional: it creates the vocabulary and fail-closed policy mechanics needed for later game-server features while preserving the already verified request-token route behavior.

## Agent Reasoning Summary

After the gate, implementing the route policy vocabulary was lower risk and higher leverage than jumping to logout, reconnect, protocol carriers, or broader modules. Without this implementation, future work would keep debating what identity source a route is allowed to trust.

The implementation is deliberately route-scoped. It lets future routes opt into bound/session policy families by explicit configuration, while inventory and other ordinary protected routes continue using the request-level access-token proof path.

## Decision Weights

```yaml
decision_weights:
  route_authorization_correctness: high
  nakama_pitaya_alignment: high
  transport_protocol_app_separation: high
  ordinary_route_behavior_stability: high
  future_logout_reconnect_readiness: medium
  immediate_protocol_feature_surface: low
  direct_nakama_pitaya_api_compatibility: low
confidence: high
```

## Consequences

- `RouteProtector` can classify routes into public, request-token, bound-connection, session-validated, and bound-session policies.
- Default protected routes remain request-token protected.
- Bound/session identity can satisfy only explicitly classified routes.
- Metadata-only identity remains rejected for protected policy families.
- `SessionValidated=true` is accepted only from a provided validated identity; route policy does not assert it directly.
- Focused Go tests cover the new policy behavior and unchanged request-token behavior.
- `runtime.bound_identity_route_policy_implementation` becomes the repository check rule for this slice.
- The work queue blocks again after implementation at `M-081/W-0153`.

## Reversal Conditions

Revisit this decision if a future ADR chooses handshake-level authentication as the primary route identity source, adopts direct Nakama or Pitaya public API compatibility, requires connection-bound identity to replace request proof globally, changes the runtime session proof carrier posture, or introduces a production route registry that supersedes the application policy type.

## Follow-Up

- Define logout/revocation active-connection behavior before revocation closes or invalidates WebSocket connections.
- Define reconnect and connection epoch behavior before duplicate replacement or resume behavior.
- Define protocol session carriers before clients receive or carry session ids.
- Explicitly classify any future production route before allowing it to rely on bound or session-validated identity instead of request-token proof.
