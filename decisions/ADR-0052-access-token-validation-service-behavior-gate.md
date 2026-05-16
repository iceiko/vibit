# ADR-0052: Access Token Validation Service Behavior Gate

Status: Accepted
Date: 2026-05-16
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-16-define-access-token-validation-service-behavior-gate/`

Related conversations:

- `conversations/2026-05-16-access-token-validation-service-behavior-gate.md`

Related artifacts:

- `docs/access-token-validation-service-behavior-gate.md`
- `docs/access-token-validation-service-behavior-gate.zh-CN.md`
- `docs/authentication-service-behavior-implementation-gate.md`
- `runtime/internal/app/authentication/service.go`
- `runtime/internal/app/authentication/service_test.go`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`

## Context

`AuthenticateWithDeviceCredential` now executes a bounded device credential login flow and can issue opaque access tokens through the application authentication service. `ValidateAccessToken` still fails closed and does not produce validated request identity.

The next behavior slice is security-sensitive because it turns opaque token proof into `RequestIdentity` for domain dispatch. It must not be mixed with protocol carrier selection, WebSocket handshake authentication, route protection, session persistence, repository interface changes, migrations, startup wiring, or external authentication dependencies.

## Decision

Define the access-token validation service behavior gate before adding validation execution.

Future validation behavior remains application-owned under:

```text
runtime/internal/app/authentication
```

The future code target remains:

```text
runtime/internal/app/authentication/service.go
runtime/internal/app/authentication/service_test.go
```

The first future validation behavior must use the existing skeleton method:

```text
ValidateAccessToken
```

The future implementation must treat `AccessToken` as opaque, high-entropy, URL-safe unpadded Base64 material already extracted into the service request. It must reject missing or malformed proof before any unit-of-work or repository call, compute the token lookup digest before lookup, use authentication and player repositories only through the unit-of-work boundary, check token lifecycle and audience, compare token verifier digest before producing identity, require an active player account, and construct a validated player `RequestIdentity` only after proof and account state are accepted.

The first public disclosure posture collapses lookup miss, wrong token posture, expired token, revoked token, wrong audience, verifier mismatch, and inactive player account to `AUTHENTICATION_TOKEN_INVALID`. More specific public expired, revoked, or account-disabled disclosure requires a later explicit decision.

This decision does not implement access-token validation, protocol carriers, WebSocket handshake authentication, route protection, session persistence, logout, refresh, cleanup jobs, token validation audit mutation, repository changes, migrations, startup wiring, generated files, external dependencies, or production authentication behavior.

## Alternatives Considered

- Implement access-token validation immediately.
- Parse `Bearer` proof strings inside the service before a protocol carrier decision.
- Let WebSocket transport or Protobuf handlers validate tokens directly.
- Let the authentication repository compare raw token proof or construct `RequestIdentity`.
- Publicly expose expired, revoked, or disabled-account distinctions in the first validation behavior.
- Update token validation audit timestamps in the first service slice.
- Introduce session persistence or route protection in the same step.

## Rationale

The project now has enough helper and repository structure to define validation precisely, but executing validation changes request identity trust. That trust boundary should be machine-checkable before code changes.

Keeping validation in the application service follows vibit's architecture: transport and protocol code carry decoded request data, repositories store digest records, and domain modules consume validated identity. The service owns proof validation and public error collapse.

Nakama reinforces token/session validation as a core backend capability. Pitaya reinforces that realtime route handlers should receive context after validation rather than own proof parsing. vibit adapts those capabilities through an explicit agent-native service gate.

## Agent Reasoning Summary

The next useful step is a gate that makes future token validation reviewable: proof parsing, digest lookup, verifier comparison, token lifecycle checks, player account state, request identity handoff, public error collapse, redaction, and tests.

## Decision Weights

```yaml
decision_weights:
  request_identity_trust_boundary: high
  security_boundary_clarity: high
  agent_context: high
  validation_sequence_explicitness: high
  public_error_disclosure_control: high
  protocol_separation: high
  repository_boundary_control: high
  session_persistence_deferral: high
  immediate_product_behavior: low
confidence: high
```

## Consequences

- `docs/access-token-validation-service-behavior-gate.md` becomes the standard for future access-token validation implementation.
- `runtime.access_token_validation_service_behavior_gate` becomes the repository check rule for this gate.
- The next work item can implement validation behavior only inside the existing service boundary.
- Protocol carriers, WebSocket handshake authentication, route protection, session persistence, logout, refresh, cleanup, repository changes, migration changes, dependencies, and production behavior remain deferred.

## Reversal Conditions

Revisit this decision if:

- The first token posture changes away from opaque access-token material.
- A future protocol carrier decision requires a different service request shape.
- The transaction boundary stops exposing repository factories through concrete unit-of-work capabilities.
- The project decides to publicly distinguish expired, revoked, or disabled-account token failures in first validation behavior.
- Session persistence becomes mandatory before route-level validation can be trusted.

## Follow-Up

- Implement access-token validation as a bounded service behavior slice only after a future work item authorizes it.
- Keep protocol carriers, startup wiring, repository changes, migrations, dependencies, session persistence, and production behavior behind later bounded work items.
