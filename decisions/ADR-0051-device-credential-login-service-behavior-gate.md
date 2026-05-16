# ADR-0051: Device Credential Login Service Behavior Gate

Status: Accepted
Date: 2026-05-16
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-16-define-device-credential-login-service-behavior-gate/`

Related conversations:

- `conversations/2026-05-16-device-credential-login-service-behavior-gate.md`

Related artifacts:

- `docs/device-credential-login-service-behavior-gate.md`
- `docs/device-credential-login-service-behavior-gate.zh-CN.md`
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

The authentication service skeleton now exists under `runtime/internal/app/authentication`. It reserves the service method vocabulary and fails closed, but does not execute login, validate access tokens, issue tokens, call repositories, expose protocol carriers, wire startup, or add production authentication behavior.

The next concrete behavior slice will be device credential login. That slice is security-sensitive because it joins proof parsing, lookup digest computation, repository lookup, verifier comparison, player account state, access-token generation, token record storage, public error mapping, and redaction.

## Decision

Define the device credential login service behavior gate before adding login execution.

Future login behavior remains application-owned under:

```text
runtime/internal/app/authentication
```

The future code target remains:

```text
runtime/internal/app/authentication/service.go
runtime/internal/app/authentication/service_test.go
```

The first future login behavior must use the existing skeleton method:

```text
AuthenticateWithDeviceCredential
```

The future implementation must treat `CredentialProof` as server-issued, high-entropy, URL-safe unpadded Base64 device credential material. It must reject missing or malformed proof before any unit-of-work or repository call, compute the credential lookup digest before lookup, use the authentication repository and player account repository only through the unit-of-work boundary, compare the credential verifier digest before token generation, require an active credential and active player account, generate opaque access-token material only after proof acceptance, store token digests only, and return the raw access-token text only once after the unit of work succeeds.

This decision does not implement login execution, token issuance, token validation, logout, refresh, cleanup jobs, protocol carriers, repository changes, migrations, startup wiring, generated files, external dependencies, or production authentication behavior.

## Alternatives Considered

- Implement device credential login immediately.
- Let the authentication repository compare proof and decide login results.
- Let the PostgreSQL adapter issue tokens.
- Parse `Bearer` proof strings inside the service before a protocol carrier decision.
- Change the global `tx.UnitOfWork` interface before the service needs it.
- Expose lookup miss, verifier mismatch, disabled credential, or missing player as distinct public login failures.
- Choose a production token lifetime or external id package in this gate.

## Rationale

The helper chain is ready, but real login is the first slice that creates an actual authenticated token. The project should force the behavior sequence to be explicit before implementation.

Keeping proof validation application-owned follows vibit's architecture and keeps transport, protocol, repository, and persistence adapters narrow. Using a local unit-of-work capability interface lets the login service access existing repositories without widening the global transaction interface prematurely.

Nakama reinforces server-side account authentication and token issuance as core game backend capabilities. Pitaya reinforces that realtime handlers should receive identity context after validation rather than own proof validation. vibit adapts those capabilities through agent-native boundaries instead of copying public APIs.

## Agent Reasoning Summary

The next useful step is not more skeleton code. It is a gate that makes the future real login path machine-checkable: proof pre-validation, helper order, repository handoff, player account check, token issuance, public error collapse, redaction, and tests.

## Decision Weights

```yaml
decision_weights:
  security_boundary_clarity: high
  agent_context: high
  login_sequence_explicitness: high
  redaction_safety: high
  repository_boundary_control: high
  future_testability: high
  protocol_separation: high
  dependency_minimization: high
  immediate_product_behavior: low
confidence: high
```

## Consequences

- `docs/device-credential-login-service-behavior-gate.md` becomes the standard for future device credential login implementation.
- `runtime.device_credential_login_service_behavior_gate` becomes the repository check rule for this gate.
- The next work item can implement login behavior only inside the existing service boundary.
- Access-token validation, logout, refresh, cleanup, protocol carriers, startup wiring, repository changes, migration changes, dependencies, and production behavior remain deferred.

## Reversal Conditions

Revisit this decision if:

- The first login method changes away from `device_credential_login`.
- The first token posture changes away from opaque access-token material.
- A future protocol carrier decision requires a different service request shape.
- The transaction boundary stops exposing repository factories through concrete unit-of-work capabilities.
- A security review requires more specific public failure disclosure or stricter collapse.

## Follow-Up

- Implement device credential login as a bounded service behavior slice only after a future work item authorizes it.
- Keep access-token validation, protocol carriers, startup wiring, repository changes, migrations, dependencies, and production behavior behind later bounded work items.
