# ADR-0079: Protocol Logout Route Gate

Status: Accepted
Date: 2026-05-18
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-18-ratify-nakama-pitaya-product-parity-roadmap/`
- `changes/2026-05-18-define-protocol-logout-route-gate/`

Related conversations:

- `conversations/2026-05-18-protocol-logout-route-gate.md`

Related artifacts:

- `docs/protocol-logout-route-gate.md`
- `docs/protocol-logout-route-gate.zh-CN.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
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

`ADR-0078` ratified the Nakama/Pitaya-class product parity roadmap and selected lifecycle closure as the near-term execution focus. The next concrete direction in that roadmap is:

```text
define_protocol_logout_route_gate
```

The runtime already has service-level `LogoutAccessToken` behavior from `ADR-0073`. It revokes only the verified presented opaque access-token record, returns success only after unit-of-work commit, and keeps runtime session revocation and active socket behavior deferred.

The remaining client-facing gap is that clients cannot request logout through the WebSocket Protobuf route surface. That route must be defined carefully because logout crosses proof carriers, route protection policy, application handler registration, protocol payload mapping, transaction bypass, error redaction, and future socket/session lifecycle boundaries.

Nakama informs this boundary because its server runtime has separate session logout and session disconnect operations. Pitaya informs this boundary because its architecture separates acceptors, sessions, routes, handlers, and connection management. vibit should expose logout without collapsing those surfaces.

## Decision

Select:

```text
define_protocol_logout_route_gate
```

Create a gate-only standard:

```text
docs/protocol-logout-route-gate.md
docs/protocol-logout-route-gate.zh-CN.md
```

Define the repository check rule:

```text
runtime.protocol_logout_route_gate
```

The gate authorizes a later bounded implementation slice to expose:

```text
runtime.authentication.LogoutAccessToken
```

as an explicit command route using future Protobuf messages:

```text
vibit.authentication.v1.LogoutAccessTokenRequest
vibit.authentication.v1.LogoutAccessTokenResponse
```

The first proof carrier posture is `access_token_in_logout_request_payload`, and the first route-protection posture is `explicit_service_validated_token_lifecycle_route`. The route should not use the `AuthenticatedRequest` wrapper because the logout service must receive the exact presented access token to revoke it after verifier comparison.

The future implementation must keep the existing Protobuf envelope unchanged, keep WebSocket transport credential-neutral, call only the existing application authentication service, bypass the outer transactional dispatcher, preserve public error collapse, and avoid socket close/session revocation side effects.

This ADR does not add Protobuf logout messages, generated output, route registration, transport behavior, close codes, close reason text, active connection invalidation, runtime session revocation, reconnect/epoch behavior, protocol session carriers, dependencies, broader product modules, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Expose logout through the existing `AuthenticatedRequest` wrapper and omit token proof from the logout payload.
- Make logout a normal protected gameplay route.
- Add logout route implementation immediately without a gate.
- Close the active WebSocket socket after logout succeeds.
- Revoke the linked runtime session after logout succeeds.
- Treat already revoked tokens as idempotent logout success.
- Copy Nakama's session API shape directly.
- Copy Pitaya's session kick model directly.

## Rationale

The logout route is not a normal protected route. Its purpose is to revoke the proof it receives. If the route-protection wrapper consumes the proof before dispatch, the application service cannot revoke that exact presented token without adding a hidden side channel. Carrying the access token in `LogoutAccessTokenRequest` keeps the behavior explicit and lets the existing service perform the same validation sequence used by service-level logout.

The decision also preserves the lifecycle separation established by earlier ADRs. Logout invalidates token material; it does not automatically close sockets, revoke sessions, replace duplicate connections, or decide reconnect behavior. Those behaviors are still important for Nakama/Pitaya-class parity, but they need separate gates because their failure and player-facing semantics are different.

## Agent Reasoning Summary

After product parity ratification, the best next step is not presence, chat, social, or match runtime. The lowest unstable shared surface is still lifecycle: clients need a logout route, but route exposure must not erase the separation between proof validation, route handling, socket close, session revocation, and reconnect.

## Decision Weights

```yaml
decision_weights:
  lifecycle_closure: high
  client_visible_logout_capability: high
  token_proof_redaction: high
  route_policy_clarity: high
  nakama_pitaya_alignment: high
  transport_protocol_app_separation: high
  implementation_scope_control: high
  direct_api_compatibility: low
  immediate_socket_close_behavior: low
confidence: high
```

## Consequences

- `runtime.protocol_logout_route_gate` becomes the repository check rule for this boundary.
- `docs/protocol-logout-route-gate.md` becomes the standard for future logout route implementation.
- The next ready work item is `W-0170`, a bounded implementation slice for the protocol logout route.
- The future implementation may add logout Protobuf messages and generated output through the normal Buf path.
- The future implementation must not close sockets, revoke runtime sessions, add reconnect behavior, or copy Nakama/Pitaya public APIs.

## Reversal Conditions

Revisit this decision if a later ADR selects handshake-level authentication as the only accepted proof carrier, removes service-level presented-token logout, adopts direct Nakama/Pitaya API compatibility, or decides that logout must be idempotent success for already revoked tokens.

## Follow-Up

- Implement the bounded protocol logout route slice.
- Define concrete transport close handoff after route exposure.
- Define reconnect and connection epoch behavior.
- Define protocol session carrier behavior.
- Continue to presence lifecycle only after logout/close/reconnect/session-carrier semantics are stable.
