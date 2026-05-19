# ADR-0053: Access Token Protocol Carrier And Route Protection Gate

Status: Accepted
Date: 2026-05-17
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-17-define-access-token-protocol-carrier-route-protection-gate/`

Related conversations:

- `conversations/2026-05-17-access-token-protocol-carrier-route-protection-gate.md`

Related artifacts:

- `docs/access-token-protocol-carrier-route-protection-gate.md`
- `docs/access-token-protocol-carrier-route-protection-gate.zh-CN.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`AuthenticateWithDeviceCredential` and `ValidateAccessToken` now execute inside the application authentication service. The service can issue opaque access tokens and turn valid token proof into application-owned `RequestIdentity`.

The next boundary is exposing token proof to client protocol traffic and protecting normal gameplay routes. That boundary is security-sensitive because it can accidentally put authentication behavior into WebSocket transport, Protobuf envelope metadata, route fields, domain handlers, repository code, startup wiring, or session persistence.

## Decision

Define the access-token protocol carrier and route-protection gate before implementation.

The first selected future posture is request-level validation with an explicit Protobuf payload wrapper for protected routes. The existing envelope remains unchanged. WebSocket handshake authentication remains deferred. Session persistence remains deferred.

The planned wrapper candidate is:

```text
vibit.authentication.v1.AuthenticatedRequest
```

The future wrapper carries an opaque access token plus the original inner payload type and bytes. The outer envelope still routes to the original domain command or query. Protocol adapter code may extract the access-token field only as a narrow handoff to application-owned validation. Application route policy must require validated player identity for protected routes before domain dispatch. Domain handlers must consume `RequestIdentity` and must not parse tokens.

This decision does not add `.proto` files, generated Protobuf output, route-protection code, protocol adapter behavior, startup wiring, WebSocket handshake authentication, session persistence, repository changes, migrations, logout, refresh, cleanup, dependencies, or production authentication behavior.

## Alternatives Considered

- Put access tokens in the existing Protobuf `Session` metadata fields.
- Add access-token fields directly to the existing `Envelope`.
- Parse HTTP `Authorization: Bearer` values in WebSocket transport.
- Use cookies, query strings, or WebSocket subprotocols for the first posture.
- Trust metadata-only `player_id` or `session_id` for protected gameplay routes.
- Add session persistence or WebSocket handshake authentication before request-level route protection.
- Let each domain route define and parse its own token field.

## Rationale

The payload-wrapper posture keeps the current envelope stable and avoids treating metadata as proof. It also keeps WebSocket transport credential-neutral, while still creating a concrete path to protect gameplay routes before domain handlers run.

Request-level validation fits the existing `SessionValidator` and application dispatch shape, but the final implementation must still add explicit route policy and protocol adapter tests. The gate keeps implementation from drifting into ad hoc proof parsing or implicit route defaults.

Nakama guides the need for token/session validation before gameplay requests. Pitaya guides the separation between connection handling and route handler identity context. vibit adapts both through an explicit application-owned validation and route-policy boundary.

## Agent Reasoning Summary

After service-local login and token validation, the highest-value next step is to expose proof to client protocol traffic in a controlled way. A gate is needed first because protocol carrier selection, route protection, and startup composition touch long-lived boundaries.

## Decision Weights

```yaml
decision_weights:
  request_identity_trust_boundary: high
  protocol_stability: high
  transport_credential_neutrality: high
  route_protection_clarity: high
  session_persistence_deferral: high
  implementation_cost: medium
  reversibility: medium
confidence: high
```

## Consequences

- `docs/access-token-protocol-carrier-route-protection-gate.md` becomes the standard for future protocol carrier and route-protection implementation.
- `runtime.access_token_protocol_carrier_route_protection_gate` becomes the repository check rule for this gate.
- The next implementation work can target a bounded payload-wrapper plus route-policy slice.
- WebSocket handshake authentication, session persistence, envelope changes, startup wiring, repository changes, migrations, logout, refresh, cleanup, and dependencies remain deferred.

## Reversal Conditions

Revisit this decision if:

- The project selects handshake-level authentication before request-level route protection.
- The existing envelope must version before protected gameplay routes can be carried.
- A future client compatibility requirement makes payload wrapper ergonomics unacceptable.
- Session persistence becomes mandatory before route-level request validation can be trusted.
- A later ADR chooses a different proof carrier for browser or non-browser clients.

## Follow-Up

- Implement the bounded payload-wrapper and route-policy slice only after a future work item authorizes the exact `.proto`, generated output, protocol adapter, application policy, tests, and startup deferrals.
