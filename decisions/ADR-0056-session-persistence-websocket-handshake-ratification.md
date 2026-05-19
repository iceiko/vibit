# ADR-0056: Session Persistence And WebSocket Handshake Ratification

Status: Accepted
Date: 2026-05-17
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-authentication-login-route/`
- `changes/2026-05-17-ratify-session-persistence-websocket-handshake-authentication/`

Related conversations:

- `conversations/2026-05-17-authentication-login-route-next-direction-and-session-handshake-ratification.md`

Related artifacts:

- `docs/session-persistence-websocket-handshake-ratification.md`
- `docs/session-persistence-websocket-handshake-ratification.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

The runtime now has a public device credential login command route, an opaque access-token response, request-level route protection through a Protobuf authenticated payload wrapper, access-token validation in the application authentication service, and PostgreSQL startup composition for those pieces.

The work queue reached `M-049/W-0121`, which asks for the next major direction after the public login route. The candidate directions included session persistence and WebSocket handshake authentication, logout/refresh/cleanup behavior, token rotation, memory durable authentication support, broader game backend modules, and operations tooling.

The maintainer asked the agent to recommend the next ten steps and continue according to that recommendation, while keeping Nakama and Pitaya as important references.

## Decision

Select `ratify_session_persistence_and_websocket_handshake_authentication` as the next direction and accept the first ratified session/handshake posture.

The selected current path remains:

- Request-level opaque access-token validation.
- The existing `vibit.authentication.v1.AuthenticatedRequest` Protobuf payload wrapper as the proof carrier for protected requests.
- The existing public login route as the token issuance path.
- `RequestIdentity.SessionValidated == false` until later session validation exists.
- WebSocket transport remains credential-neutral.
- The existing Protobuf envelope remains unchanged.

The preferred future connection-level gate is first-message protocol binding owned by protocol/application layers, not immediate WebSocket handshake credential parsing.

The preferred first durable target for future session persistence is PostgreSQL, but this ADR does not add a session table, migration, repository interface, adapter, cleanup job, or dependency.

## Alternatives Considered

- Implement WebSocket handshake authentication immediately.
- Put access tokens into HTTP `Authorization`, Bearer, cookies, query strings, or WebSocket subprotocols.
- Add session fields or token proof fields to the existing Protobuf envelope.
- Add a runtime session table before first-message binding behavior is specified.
- Implement logout, refresh, cleanup, token rotation, or active connection invalidation first.
- Expand into social, chat, presence, matchmaking, or other game backend modules before session/connection lifecycle is ratified.

## Rationale

Nakama guides the capability sequence: clients authenticate, receive token or session material, and then use normal gameplay or realtime socket features. vibit has now reached the point where the next architectural question is how validated identity attaches to ongoing connections and future sessions.

Pitaya guides the architecture vocabulary: acceptors, sessions, route handlers, groups, and cluster roles are separate concerns. vibit adapts this by keeping WebSocket transport narrow and credential-neutral, and by planning connection binding through a protocol/application gate instead of hiding authentication inside the transport adapter.

Request-level validation is already implemented and verifiable. Keeping it as the current path avoids premature coupling to browser-specific carrier behavior, pre-envelope error semantics, or transport-owned auth. First-message binding is a better next gate because it can use the existing WebSocket Protobuf request loop while still defining explicit connection-bound state later.

## Agent Reasoning Summary

After public login route exposure, the highest-value next step is not another token lifecycle feature or a new game module. It is ratifying the identity-to-connection/session model that later logout, reconnect, presence, rooms, parties, matches, and groups will depend on. The gate deliberately stops short of implementation so future agents cannot accidentally add token carriers, session tables, or handshake auth in the wrong layer.

## Decision Weights

```yaml
decision_weights:
  game_server_capability_sequence: high
  transport_boundary_preservation: high
  compatibility_with_existing_request_level_validation: high
  future_reconnect_presence_room_match_foundation: high
  implementation_risk_reduction: high
  immediate_feature_delivery: medium
  reversibility: medium
confidence: high
```

## Consequences

- `docs/session-persistence-websocket-handshake-ratification.md` becomes the standard for the first session/handshake posture after login.
- `runtime.session_persistence_websocket_handshake_ratification` becomes the repository check rule.
- The work queue moves past `M-049/W-0121` and records `M-050/W-0122` as completed.
- The next direction is blocked again at `M-051/W-0123` before implementing first-message binding, session schema, handshake authentication, logout/revocation, reconnect/epoch behavior, operations posture, or broader game backend modules.
- Existing runtime behavior is not changed by this ADR.

## Reversal Conditions

Revisit this decision if:

- A future client compatibility requirement requires handshake-level authentication before first-message binding.
- A future non-WebSocket transport changes the tradeoff for request-level validation.
- A later ADR adopts direct Nakama or Pitaya public API compatibility.
- Operational load from unauthenticated WebSocket connections becomes a primary near-term constraint.
- Session persistence must be implemented before any connection binding semantics can be useful.

## Follow-Up

- Define a first-message connection binding protocol gate if connection-bound identity is selected next.
- Define a PostgreSQL session persistence schema gate before adding any session table, repository, migration, or cleanup job.
- Define logout/revocation and active-connection invalidation before logout behavior affects bound connections.
