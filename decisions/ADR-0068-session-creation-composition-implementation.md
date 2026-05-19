# ADR-0068: Session Creation Composition Implementation

Status: Accepted
Date: 2026-05-18
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-session-creation-composition-gate/`
- `changes/2026-05-18-implement-session-creation-composition/`

Related conversations:

- `conversations/2026-05-18-session-creation-composition-implementation.md`

Related artifacts:

- `runtime/internal/app/authentication/service.go`
- `runtime/internal/app/authentication/service_test.go`
- `runtime/cmd/vibit-server/main.go`
- `runtime/cmd/vibit-server/main_test.go`
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

`W-0146` defined the session creation composition gate after durable `runtime_sessions`, the storage-neutral `session.Repository`, the PostgreSQL adapter, runtime session validation, device-credential login, access-token validation, and startup authentication composition were already present.

The work queue reached `M-075/W-0147`, a confirmation gate. The maintainer asked the agent to recommend the next ten steps and continue, with Nakama and Pitaya as key reference baselines.

The implementation gap was concrete: login could issue and store an opaque access token, and runtime session validation could validate a persisted session, but no production login path created the durable session row that later session-aware behavior can validate.

## Decision

Select:

```text
implement_session_creation_composition
```

Implement login-time durable session creation in:

```text
runtime/internal/app/authentication/service.go
runtime/internal/app/authentication/service_test.go
runtime/cmd/vibit-server/main.go
runtime/cmd/vibit-server/main_test.go
```

`AuthenticateWithDeviceCredential` now composes, in one application unit of work:

1. Device credential proof validation.
2. Credential verifier comparison.
3. Active player account validation.
4. Opaque access-token material generation.
5. Access-token lookup/verifier digest storage.
6. Server-owned runtime session id generation.
7. `session.Repository.CreateRuntimeSession`.
8. Commit.
9. Successful result return after commit.

The runtime session is created as an active player session, links to the stored `access_token_record_id`, initializes `last_seen_at` to `issued_at`, and uses the same first-posture expiration as the access token.

This ADR does not expose session ids through Protobuf login responses, change the existing Protobuf envelope, add Protobuf session messages, add generated output, change WebSocket handshake authentication, add transport credential carriers, wire session validation into protected route policy, change `ValidateAccessToken`, set `RequestIdentity.SessionValidated = true` in login or token validation, implement logout/revocation active-connection invalidation, add refresh, cleanup, token rotation, reconnect/epoch behavior, add dependencies, add memory durable session behavior, or adopt direct Nakama/Pitaya public API compatibility.

## Alternatives Considered

- Keep login token-only and defer session creation again.
- Create sessions during `BindConnection`.
- Create sessions in WebSocket transport or Protobuf protocol adapters.
- Expose the new `session_id` in the Protobuf login response immediately.
- Use session creation as ordinary route authorization immediately.
- Implement logout/revocation active-connection behavior in the same slice.
- Copy Nakama or Pitaya session APIs directly.

## Rationale

Nakama treats authentication as the entry point to a session lifecycle with expiration, refresh, logout, and management implications. vibit adapts that lesson by making login create durable session state transactionally with token storage, so later route policy, logout, and management work has a real lifecycle object to reason about.

Pitaya keeps acceptors, handler routing, and session context separated. vibit adapts that lesson by keeping durable session creation in application composition rather than transport, protocol, or persistence adapters.

The result is intentionally narrower than a full Nakama-like session surface: it creates the durable row, keeps the token/session linkage private, and preserves separate gates for protocol carriers, route policy, logout, reconnect, and operations.

## Agent Reasoning Summary

After the gate, the highest-leverage next implementation was to make the existing login path create the missing durable session row. Without this, session validation existed but had no normal production creation path. Route policy would have been premature because it would depend on session data that login did not create.

The implementation keeps the transaction boundary clear: token storage and session creation succeed or fail together, and client-visible token material is returned only after commit. Failure paths collapse to existing redacted authentication-unavailable public errors.

## Decision Weights

```yaml
decision_weights:
  lifecycle_correctness: high
  unit_of_work_atomicity: high
  nakama_pitaya_alignment: high
  transport_protocol_app_separation: high
  future_route_policy_readiness: high
  immediate_protocol_feature_surface: low
  direct_nakama_pitaya_api_compatibility: low
confidence: high
```

## Consequences

- `AuthenticateWithDeviceCredential` creates a durable runtime session after token storage and before commit.
- `ServiceDependencies` requires a `SessionIDGenerator`.
- The login unit-of-work capability includes `NewSessionRepository`.
- Startup composition provides a high-entropy random runtime session id generator.
- `AuthenticationResult` carries application-owned session metadata, but protocol login responses are unchanged.
- Focused tests cover session creation success, session repository failure, session id failure, session creation failure, commit failure, and unchanged access-token validation.
- `runtime.session_creation_composition_implementation` becomes the repository check rule for this slice.
- The work queue blocks again after implementation at the next confirmation gate.

## Reversal Conditions

Revisit this decision if a future ADR chooses handshake-level authentication as the primary session creation point, makes sessions independent of access-token lifetime, adopts direct Nakama/Pitaya public API compatibility, changes the runtime session repository contract, or requires protocol clients to receive session ids during login.

## Follow-Up

- Define bound/session identity route policy before protected routes can rely on session-validated identity.
- Define logout/revocation active-connection behavior before revocation closes or invalidates WebSocket connections.
- Define reconnect and connection epoch behavior before duplicate replacement or resume behavior.
- Define a protocol/session carrier gate before exposing session ids to clients.
