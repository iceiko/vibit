# ADR-0071: Logout Revocation Active Connection Gate

Status: Accepted
Date: 2026-05-18
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-bound-identity-route-policy-implementation/`
- `changes/2026-05-18-define-logout-revocation-active-connection-gate/`

Related conversations:

- `conversations/2026-05-18-logout-revocation-active-connection-gate.md`

Related artifacts:

- `docs/logout-revocation-active-connection-gate.md`
- `docs/logout-revocation-active-connection-gate.zh-CN.md`
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

`W-0152` implemented explicit bound identity route policy. vibit can now distinguish ordinary request-token routes, explicit bound-connection routes, explicit session-validated routes, and explicit bound-session routes.

The work queue reached `M-081/W-0153`, a confirmation gate. The maintainer asked the agent to recommend the next ten steps and continue, with Nakama and Pitaya as key references.

At this point, the highest-risk missing boundary is logout and revocation behavior for active WebSocket connections. vibit has access-token validation, durable token records, login-created runtime sessions, first-message connection binding, durable session validation, and explicit route policy families, but it still has no policy for whether token logout or session revocation closes already-open sockets.

## Decision

Select:

```text
define_logout_revocation_active_connection_gate
```

Create a gate-only standard:

```text
docs/logout-revocation-active-connection-gate.md
docs/logout-revocation-active-connection-gate.zh-CN.md
```

The gate defines future ownership, policy questions, first recommended posture, error/redaction rules, existing runtime relationships, WebSocket/Protobuf deferrals, future test expectations, and Nakama/Pitaya reference mapping.

The recommended future first posture is conservative:

- Presented-token logout should be the first logout scope unless a later ADR broadens it.
- Runtime session revocation remains a separate policy choice.
- Active-connection invalidation requires an explicit connection registry and close policy before implementation.
- WebSocket transport must not own authentication state.
- Ordinary protected routes remain request-token protected by default.
- Reconnect and epoch behavior remain deferred.

This ADR does not implement `LogoutAccessToken`, revoke tokens, revoke runtime sessions, close WebSocket connections, add an active connection registry, add kick/disconnect behavior, add reconnect or epoch behavior, add Protobuf logout routes, add protocol session carriers, change the existing Protobuf envelope, change WebSocket handshake authentication, add transport credential carriers, add cleanup jobs, add dependencies, change route classification, add memory durable session behavior, or adopt direct Nakama/Pitaya public API compatibility.

## Alternatives Considered

- Implement `LogoutAccessToken` immediately.
- Revoke the linked runtime session whenever a token is logged out.
- Close active WebSocket connections immediately on token logout.
- Add an active connection registry before defining policy.
- Treat first-message bound identity as enough to find and close sockets.
- Reclassify ordinary routes to bound-session policy before logout/revocation behavior is defined.
- Add a Protobuf logout route in the same slice.
- Copy Nakama or Pitaya public APIs directly.

## Rationale

Nakama shows that authentication session material has lifecycle consequences: expiration, refresh, logout, revocation, and realtime socket behavior all matter for gameplay access. vibit adapts that by requiring explicit policy for token revocation, runtime session revocation, and active socket invalidation instead of letting storage or transport code decide implicitly.

Pitaya shows the value of separating acceptors, sessions, route handlers, and connection management. vibit adapts that by keeping future logout/revocation policy in the application layer and limiting any future transport participation to a narrow socket-close handoff.

The gate is needed before implementation because active socket invalidation crosses storage, application policy, route authorization, transport side effects, and future reconnect behavior. Combining all of that into a logout implementation would make the boundary too broad and easy to regress.

## Agent Reasoning Summary

After route policy implementation, the next practical question is no longer whether a route can distinguish identity sources; it is how the server reacts when one of those identity sources becomes invalid while a connection is still open.

Choosing a gate-only logout/revocation active-connection boundary keeps the work aligned with Nakama and Pitaya without copying either system. It also prevents hidden WebSocket-side authentication state, metadata-only socket targeting, and silent best-effort revocation behavior.

## Decision Weights

```yaml
decision_weights:
  security_and_revocation_correctness: high
  active_connection_lifecycle_clarity: high
  nakama_pitaya_alignment: high
  transport_protocol_app_separation: high
  implementation_scope_control: high
  immediate_user_visible_feature_surface: low
  direct_nakama_pitaya_api_compatibility: low
confidence: high
```

## Consequences

- `runtime.logout_revocation_active_connection_gate` becomes the repository check rule for this boundary.
- Future logout implementation must answer whether token revocation also revokes runtime sessions.
- Future active-connection invalidation requires an explicit connection registry and transport handoff policy.
- WebSocket transport remains credential-neutral and does not own authentication state.
- Protobuf logout routes, close policy, session carriers, reconnect, and epoch behavior remain separate gates.
- The work queue blocks again after the gate at `M-083/W-0155`.

## Reversal Conditions

Revisit this decision if a future ADR selects handshake-level authentication as the primary identity source, adopts a direct Nakama or Pitaya compatibility surface, requires account-wide revocation as the first logout posture, introduces a cluster session registry before the single-process registry boundary, or chooses a transport-owned authentication model.

## Follow-Up

- Define presented-token logout execution before implementing `LogoutAccessToken`.
- Define an active connection registry before targeting open sockets.
- Define transport close policy before using custom WebSocket close codes or reason text.
- Define reconnect and epoch behavior before duplicate replacement or resume behavior.
- Define protocol logout/session carriers before exposing logout commands, session ids, resume tokens, or connection epochs to clients.
