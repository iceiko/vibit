# ADR-0069: Bound Identity Route Policy Gate

Status: Accepted
Date: 2026-05-18
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-session-creation-composition-implementation/`
- `changes/2026-05-18-define-bound-identity-route-policy-gate/`

Related conversations:

- `conversations/2026-05-18-bound-identity-route-policy-gate.md`

Related artifacts:

- `docs/bound-identity-route-policy-gate.md`
- `docs/bound-identity-route-policy-gate.zh-CN.md`
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

`W-0148` implemented session creation composition. Successful device-credential login now stores the access-token verifier record and creates one active durable runtime session in the same unit of work.

The work queue reached `M-077/W-0149`, a confirmation gate. The maintainer asked the agent to recommend the next ten steps and continue, with Nakama and Pitaya as key reference baselines.

At this point, vibit has request-level access-token route protection, first-message connection binding, durable session persistence, runtime session validation, and login-created durable sessions. The missing boundary is how route policy should use, combine, or deliberately avoid using request proof, bound connection identity, and session-validated identity.

## Decision

Select:

```text
define_bound_identity_route_policy_gate
```

Define the gate-only standard:

```text
docs/bound-identity-route-policy-gate.md
docs/bound-identity-route-policy-gate.zh-CN.md
```

The gate defines future application-owned route policy families:

- `public`
- `request_token_required`
- `bound_connection_required`
- `session_validated_required`
- `bound_session_required`

It recommends that the first implementation remain conservative: ordinary protected domain routes should continue to require request-level access-token proof, while bound connection identity and session-validated identity become explicit route-scoped policy families only after a later implementation slice.

This ADR does not change production route authorization behavior, use bound identity for ordinary protected routes, use session-validated identity for ordinary protected routes, remove per-request access-token proof, change WebSocket handshake authentication, add transport credential carriers, expose session ids through Protobuf, change the existing Protobuf envelope, implement logout/revocation active-connection behavior, implement reconnect/epoch behavior, add generated output, add dependencies, add memory durable session behavior, or adopt direct Nakama/Pitaya public API compatibility.

## Alternatives Considered

- Implement bound identity route policy immediately.
- Make first-message `BindConnection` sufficient for all ordinary protected routes.
- Make durable `session_id` sufficient for protected routes.
- Require both bound connection identity and session validation for all routes immediately.
- Expose session ids in login responses before route policy is defined.
- Move identity policy into WebSocket transport or Protobuf adapters.
- Prioritize logout/revocation active-connection behavior before route policy.
- Prioritize reconnect/epoch behavior before route policy.
- Copy Nakama or Pitaya public APIs directly.

## Rationale

Nakama shows that authenticated session material eventually becomes central to gameplay API access, but also that session lifetime, logout, refresh, and socket behavior have separate lifecycle pressure. vibit adapts this by defining route-scoped policy families instead of treating every identity source as interchangeable proof.

Pitaya shows a useful separation between acceptors, sessions, and route handlers. vibit adapts this by keeping route policy in `runtime/internal/app`, keeping WebSocket transport credential-neutral, and requiring domain handlers to receive normalized identity context rather than parsing transport or protocol credentials.

The gate is deliberately narrower than a full session subsystem: it defines how future route policy should reason about identity sources without changing route behavior yet.

## Agent Reasoning Summary

After login-time session creation, the highest-leverage next boundary is route policy. Implementing logout, reconnect, protocol session carriers, or broader social/realtime modules before route policy would leave the core question unanswered: what identity source may authorize which route.

The recommended route-policy gate lets vibit absorb mature game server lessons while preserving agent-native maintainability: route policy is explicit, application-owned, route-scoped, redacted, and checkable.

## Decision Weights

```yaml
decision_weights:
  route_authorization_correctness: high
  nakama_pitaya_alignment: high
  transport_protocol_app_separation: high
  future_logout_reconnect_readiness: medium
  immediate_behavior_change: low
  direct_nakama_pitaya_api_compatibility: low
confidence: high
```

## Consequences

- `docs/bound-identity-route-policy-gate.md` becomes the canonical gate for future bound/session route policy.
- The Chinese translation is maintained in the same change.
- `runtime.bound_identity_route_policy_gate` becomes the repository check rule for this gate.
- The work queue advances from `M-077/W-0149` to a completed gate milestone and then blocks again at the next direction confirmation gate.
- Future implementation may add explicit route policy families and tests, but that requires a separate bounded work item.
- Existing WebSocket, Protobuf, route protection, connection binding, runtime session validation, session creation, logout, reconnect, and direct compatibility behavior remain unchanged.

## Reversal Conditions

Revisit this decision if a future ADR selects handshake-level authentication as the primary route identity source, adopts direct Nakama or Pitaya public API compatibility, requires connection-bound identity to replace request proof globally, or changes the session/proof carrier posture.

## Follow-Up

- Implement the bound identity route-policy family slice after a separate implementation gate or work item.
- Define logout/revocation active-connection behavior before revocation closes or invalidates WebSocket connections.
- Define reconnect and connection epoch behavior before duplicate replacement or resume behavior.
- Define protocol session carriers before clients carry or receive session ids.
