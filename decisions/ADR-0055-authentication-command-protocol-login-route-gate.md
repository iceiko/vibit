# ADR-0055: Authentication Command Protocol And Login Route Gate

Status: Accepted
Date: 2026-05-17
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-runtime-authentication-startup-composition/`
- `changes/2026-05-17-define-authentication-command-protocol-login-route-gate/`

Related conversations:

- `conversations/2026-05-17-runtime-authentication-startup-composition-next-direction.md`
- `conversations/2026-05-17-authentication-command-protocol-login-route-gate.md`

Related artifacts:

- `docs/authentication-command-protocol-login-route-gate.md`
- `docs/authentication-command-protocol-login-route-gate.zh-CN.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

The runtime already has service-local device credential login, service-local access-token validation, request-level route protection, and PostgreSQL startup composition for route protection.

The remaining client-facing gap is that `runtime.authentication.AuthenticateWithDeviceCredential` is declared as the explicit public route in application route policy, but the runtime does not yet have a Protobuf request/response message, protocol bridge, application route handler registration, or startup composition that exposes the existing service method to clients.

The maintainer asked for a recommendation and then authorized the recommended next ten steps. The selected direction is:

```text
add_authentication_command_protocol_messages_and_login_route_registration
```

## Decision

Define and accept the authentication command protocol and login route gate.

The next implementation may add only:

- `proto/vibit/authentication/v1/authentication.proto`
- generated Go Protobuf output from that source
- protocol bridge behavior for `AuthenticateWithDeviceCredential`
- application-owned route handler registration for the explicit public login route
- PostgreSQL startup composition that registers the login route only when the authentication service is composed
- a transaction-wrapper bypass for the authentication route so the authentication service owns its own unit-of-work
- focused tests and check-rule coverage

The existing Protobuf envelope remains unchanged. WebSocket transport remains credential-neutral. Session persistence and WebSocket handshake authentication remain deferred.

This decision does not add HTTP `Authorization`, Bearer, cookie, query-string, or WebSocket subprotocol credential carriers. It does not add logout, refresh, cleanup, token rotation, repository interface changes, PostgreSQL adapter changes, migrations, dependencies, memory-store durable authentication behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Move directly to session persistence and WebSocket handshake authentication.
- Add logout, refresh, cleanup, or token rotation before exposing login.
- Parse login credentials from WebSocket handshake headers or query strings.
- Put credential proof into the existing Protobuf `Session` metadata.
- Register all authentication commands at once.
- Let the protocol adapter call authentication repositories directly.
- Keep the login route unavailable and require out-of-band token seeding.

## Rationale

The runtime now validates tokens on protected routes, but clients need an in-band way to obtain the access token first. Exposing the already-implemented login service through an explicit public command route closes that loop without selecting session persistence or handshake authentication prematurely.

The slice remains smaller than session persistence because it uses existing service, repository, startup, route policy, and error-envelope boundaries. It also makes the next system behavior easier to verify: unauthenticated clients can call only the public login route, then use the returned opaque token in the existing request-level wrapper for protected gameplay routes.

Nakama guides the capability sequence: authenticate first, receive token/session material, then use normal server or realtime features. Pitaya guides the layering: transport acceptors, session/context, routing, and handlers stay separate. vibit adapts both through explicit Protobuf payloads and application-owned route handlers while keeping WebSocket transport credential-neutral.

## Agent Reasoning Summary

After startup composition, the highest-value next step is not deeper session machinery; it is making the existing login service reachable through the same route/protocol pipeline that already protects gameplay requests. This creates a minimal authenticate-then-gameplay loop while preserving the later option to ratify session persistence, handshake authentication, logout, refresh, cleanup, and operations posture separately.

## Decision Weights

```yaml
decision_weights:
  client_authentication_loop_closure: high
  authentication_boundary_preservation: high
  transport_credential_neutrality: high
  protocol_contract_clarity: high
  session_persistence_deferral: high
  implementation_cost: medium
  reversibility: medium
confidence: high
```

## Consequences

- `docs/authentication-command-protocol-login-route-gate.md` becomes the standard for public login route protocol and application registration.
- `runtime.authentication_command_protocol_login_route_gate` becomes the repository check rule for the gate and its first implementation slice.
- The next implementation work can expose `AuthenticateWithDeviceCredential` through Protobuf and application dispatch.
- `TransactionalDispatcher` needs a narrow bypass for the public authentication route so the authentication service owns its own unit-of-work.
- Session persistence, WebSocket handshake authentication, logout, refresh, cleanup, token rotation, repository changes, migrations, dependencies, and direct Nakama/Pitaya API compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- The project chooses handshake-level authentication before public login command exposure.
- A future protocol compatibility requirement rejects Protobuf login command payloads.
- The authentication service unit-of-work boundary changes and no longer needs route-level transaction bypass.
- A later ADR adopts direct Nakama or Pitaya public API compatibility.
- The memory runtime path gains a durable authentication repository.

## Follow-Up

- Implement the bounded Protobuf login command and route-registration slice after `W-0120` is active.
- Keep the first implementation limited to device credential login and the PostgreSQL runtime path.
