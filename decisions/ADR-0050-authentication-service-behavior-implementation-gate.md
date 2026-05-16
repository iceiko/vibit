# ADR-0050: Authentication Service Behavior Implementation Gate

Status: Accepted
Date: 2026-05-16
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-15-define-authentication-service-behavior-implementation-gate/`

Related conversations:

- `conversations/2026-05-16-authentication-service-behavior-implementation-gate.md`

Related artifacts:

- `docs/authentication-service-behavior-implementation-gate.md`
- `docs/authentication-service-behavior-implementation-gate.zh-CN.md`
- `docs/runtime-authentication-implementation-boundary.md`
- `docs/application-authentication-service-interface-boundary.md`
- `docs/token-credential-verifier-algorithm-redaction-boundary.md`
- `docs/secret-configuration-verifier-key-loading-boundary.md`
- `docs/token-credential-material-generation-boundary.md`
- `docs/verifier-digest-computation-comparison-boundary.md`
- `docs/authentication-service-implementation-readiness-gate.md`
- `docs/verifier-digest-comparison-helper-gate.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`

## Context

The repository has completed the authentication helper chain needed before service orchestration can be designed precisely:

- In-memory verifier key set validation.
- Process environment verifier key loading.
- Token and credential material generation helpers.
- Lookup and verifier digest computation helpers.
- Constant-time verifier digest comparison helpers.

The next risk is not missing a helper. The risk is service behavior becoming an unbounded place where future agents mix login execution, token validation, repository calls, protocol carriers, startup wiring, and public failure decisions without a narrow standard.

## Decision

Define the authentication service behavior implementation gate before adding service behavior code.

Future service behavior remains application-owned under:

```text
runtime/internal/app/authentication
```

The future service files are:

```text
runtime/internal/app/authentication/service.go
runtime/internal/app/authentication/service_test.go
```

The next authorized implementation should be a skeleton-only service slice unless a later work item explicitly authorizes login execution or token validation execution. The skeleton may define typed dependencies, request/result vocabulary, internal redacted failure classes, and fail-closed `not implemented` behavior. It must not call repositories, issue tokens, validate tokens, revoke tokens, refresh tokens, clean up tokens, expose Protobuf or WebSocket carriers, change repositories or migrations, wire startup, add dependencies, or implement production authentication behavior.

The future real service flow must compose existing helpers in this order:

- Parse already-decoded proof input at the application service boundary.
- Compute lookup digests before repository lookup.
- Use `authentication.Repository` only through the application unit-of-work boundary.
- Compute verifier digests with the correct verifier key context.
- Compare verifier digests with `CompareCredentialVerifierDigest` or `CompareTokenVerifierDigest`.
- Collapse lookup miss, wrong key/version, mismatch, and unusable proof states to redacted public invalid-proof errors.
- Convert validated access-token proof into `RequestIdentity` before production-sensitive domain dispatch.

This decision does not implement Go service code, login execution, token validation, logout, refresh, cleanup jobs, protocol carriers, repository changes, migrations, startup wiring, dependencies, or production authentication behavior.

## Alternatives Considered

- Implement login and token validation immediately.
- Add service behavior to existing helper files.
- Let the PostgreSQL adapter decide authentication outcomes.
- Let WebSocket or Protobuf handlers parse and validate proof directly.
- Expose token expiration and revocation distinctions in the first behavior by default.
- Create Protobuf authentication messages before application service behavior is bounded.

## Rationale

The helper chain is now ready, but authentication service behavior is the first place where multiple security-sensitive boundaries meet. A gate prevents future agents from turning a small service slice into a cross-repository feature.

Separating the service skeleton from real login and token validation gives agents a narrow next code step: define the shape and fail closed. That is useful because it can be tested without PostgreSQL, protocol carriers, or startup wiring.

Nakama reinforces the need for server-side account authentication, token issuance, and token validation. Pitaya reinforces that realtime handlers should receive validated identity context rather than own proof validation. vibit adapts both by keeping proof validation application-owned and keeping transport/protocol layers as carriers only.

## Agent Reasoning Summary

The project is ready to define the authentication service behavior gate, but not to expose authentication behavior. The correct next step is to make the orchestration plan, failure collapse posture, repository handoff, redaction requirements, and file boundaries inspectable before code.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  security_boundary_clarity: high
  public_error_redaction: high
  repository_boundary_control: high
  protocol_transport_separation: high
  future_testability: high
  dependency_minimization: high
  implementation_speed: medium
  immediate_product_behavior: low
confidence: high
```

## Consequences

- `docs/authentication-service-behavior-implementation-gate.md` becomes the standard for future authentication service behavior work.
- `runtime.authentication_service_behavior_implementation_gate` becomes the repository check rule for this gate.
- The next work item can add a skeleton-only service shape without login execution or token validation.
- Authentication service behavior remains unimplemented in this change.
- Protobuf authentication messages, WebSocket proof carriers, repository changes, migrations, startup wiring, dependencies, and production behavior remain deferred.

## Reversal Conditions

Revisit this decision if:

- A future security review requires token expired/revoked states to be publicly distinguished in first behavior.
- A future protocol carrier decision requires a different application service request shape.
- The application unit-of-work boundary changes.
- The authentication repository boundary moves out of `runtime/internal/modules/authentication`.
- The first access-token posture changes from opaque bearer material to a signed or structured claim token.

## Follow-Up

- Add a skeleton-only authentication service behavior implementation slice if the next work item authorizes it.
- Keep login execution, token validation execution, protocol carriers, startup wiring, repository changes, migrations, dependencies, and production behavior behind later bounded work items.
