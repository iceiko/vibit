# ADR-0054: Runtime Authentication Startup Composition Gate

Status: Accepted
Date: 2026-05-17
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-17-define-runtime-authentication-startup-composition-gate/`

Related conversations:

- `conversations/2026-05-17-runtime-authentication-startup-composition-gate.md`

Related artifacts:

- `docs/runtime-authentication-startup-composition-gate.md`
- `docs/runtime-authentication-startup-composition-gate.zh-CN.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

The runtime already has service-local device credential login, service-local access-token validation, and request-level Protobuf payload-wrapper route protection. The Protobuf frame handler supports a `RouteProtector`, but process startup has not wired the existing authentication service into the PostgreSQL runtime path.

The maintainer explicitly selected:

```text
wire_runtime_authentication_startup_composition
```

The maintainer also emphasized Nakama and Pitaya as important references. Nakama shows the capability expectation that clients authenticate and then use server/realtime features through session/token context. Pitaya shows the architectural expectation that acceptors, sessions, routing, and handlers remain separated.

## Decision

Define and accept the runtime authentication startup composition gate.

The first implementation may wire authentication startup composition only in `runtime/cmd/vibit-server` and only for the explicit PostgreSQL runtime path selected by:

```text
VIBIT_RUNTIME_STORE=postgres
```

The composition may use:

- Existing environment verifier key loading.
- Existing `authentication.Service`.
- Existing PostgreSQL unit-of-work repository capabilities.
- `crypto/rand.Reader` for access-token material generation through the existing service dependency.
- A startup-owned system clock.
- A startup-owned standard-library token record id generator.
- Default access-token lifetime `1h`, optionally overridden by `VIBIT_AUTH_ACCESS_TOKEN_TTL`.
- Default token audience `vibit_gameplay_runtime_requests`, optionally overridden by `VIBIT_AUTH_TOKEN_AUDIENCE`.
- Existing `authentication.NewRouteAccessTokenValidator`.
- Existing `app.NewRouteProtector`.
- Existing Protobuf frame handler route-protector injection.

The default memory runtime path remains a bootstrap path without durable authentication repository capability.

This decision does not add WebSocket handshake authentication, session persistence, authentication command Protobuf messages, login route registration, repository interface changes, PostgreSQL adapter changes, migrations, dependencies, logout, refresh, cleanup, token rotation, token validation audit mutation, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Wire authentication route protection into all runtime stores, including memory.
- Parse access-token proof from WebSocket handshake headers, query strings, cookies, or subprotocols.
- Add a new session persistence model before startup composition.
- Register authentication command routes before startup validation is wired.
- Defer startup composition and move next to logout/refresh/cleanup.
- Adopt Nakama or Pitaya public API compatibility directly.

## Rationale

The PostgreSQL path already has the repository capabilities required by `authentication.Service`. The memory path does not, so making memory fail closed for protected authentication would break the bootstrap developer path without adding durable auth value.

Startup composition is the right next step after route protection because protected routes are ineffective in the process server until the route protector is actually injected. The composition remains narrow: it wires existing pieces and does not invent new auth lifecycle behavior.

Nakama guides the product capability expectation around authenticated sessions/tokens before gameplay and realtime use. Pitaya guides the separation between connection acceptors, session/context handling, routing, and handlers. vibit adapts both by wiring authentication at startup into application/protocol composition while keeping WebSocket transport credential-neutral.

## Agent Reasoning Summary

After request-level route protection exists, the process server still needs a concrete composition point before protected gameplay routes can validate access-token proof. The PostgreSQL runtime path is the only current path with the repository capabilities required by the existing authentication service, so startup composition should be narrow, explicit, and fail-closed there while leaving the memory path as bootstrap-only behavior. This mirrors Nakama's authenticate-before-gameplay capability sequencing and Pitaya's separation of connection, session/context, route, and handler concerns without adopting either project's public API surface.

## Decision Weights

```yaml
decision_weights:
  process_server_usability: high
  authentication_boundary_preservation: high
  transport_credential_neutrality: high
  postgres_repository_capability_fit: high
  memory_bootstrap_preservation: medium
  session_persistence_deferral: high
  implementation_cost: medium
confidence: high
```

## Consequences

- `docs/runtime-authentication-startup-composition-gate.md` becomes the standard for startup composition.
- `runtime.authentication_startup_composition_gate` becomes the repository check rule for the gate and its first implementation slice.
- `runtime/cmd/vibit-server` becomes the only first write area for startup composition implementation.
- The PostgreSQL runtime path must fail closed when authentication verifier key configuration is missing or invalid.
- The memory runtime path remains unauthenticated bootstrap behavior.
- Session persistence, WebSocket handshake authentication, authentication command routes, repository changes, migrations, dependencies, logout, refresh, cleanup, and token rotation remain deferred.

## Reversal Conditions

Revisit this decision if:

- The project ratifies session persistence before request-level startup composition is used in production.
- The memory runtime path grows a storage-neutral authentication repository.
- The first public client authentication route must be exposed before protected gameplay routes are practical.
- A future operations gate requires a different secret-loading or key-rotation posture.
- A later compatibility ADR adopts a direct Nakama or Pitaya API surface.

## Follow-Up

- Implement the bounded startup composition slice in `runtime/cmd/vibit-server`.
- Add focused tests for configuration parsing, fail-closed verifier key loading, route-protector injection, and memory bootstrap preservation.
