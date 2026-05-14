# ADR-0039: Application Authentication Service Interface Boundary

Status: Accepted
Date: 2026-05-15
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-15-define-application-authentication-service-interface-boundary/`

Related conversations:

- `conversations/2026-05-15-application-authentication-service-interface-boundary.md`

Related artifacts:

- `docs/application-authentication-service-interface-boundary.md`
- `docs/application-authentication-service-interface-boundary.zh-CN.md`
- `docs/runtime-authentication-implementation-boundary.md`
- `docs/authentication-generated-contract-shape-timing.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`

## Context

Runtime authentication now has semantic contracts, selected login/token posture, PostgreSQL schema sources, a storage-neutral repository interface, a PostgreSQL adapter, an application-owned implementation boundary, and metadata-only generated authentication contract shapes.

The next risk is that a future agent may jump from generated shape files or repository code directly into login execution, token validation, or transport behavior. Before behavior appears, agents need a stable application service-interface boundary that maps generated contract vocabulary to future request/result DTOs, redaction expectations, unit-of-work usage, repository usage, error mapping, permission mapping, audit handoff, and request identity handoff.

## Decision

Define the application authentication service interface boundary before adding service code or runtime authentication behavior.

Future authentication service interfaces are application-owned under `runtime/internal/app`. They may later live directly under that package or under an application-owned child package such as `runtime/internal/app/authentication`, but no code package is created by this decision.

Generated authentication contract shapes inform service-level request and result vocabulary, but they remain metadata-only generated output. Transport handlers, Protobuf adapters, domain modules, repositories, and PostgreSQL adapters must not import generated authentication shapes as runtime behavior owners.

Future service behavior may use `authentication.Repository` only through the application unit-of-work boundary. The repository remains storage-neutral, and the PostgreSQL adapter remains persistence-only.

Token generation, verifier comparison, login execution, access-token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository interface changes, and migration schema changes remain deferred.

## Alternatives Considered

- Implement Go authentication service interfaces immediately.
- Implement `AuthenticateWithDeviceCredential` before defining the service boundary.
- Put token validation in the Protobuf adapter.
- Let WebSocket transport own the first proof carrier and service entrypoint.
- Let `authentication.Repository` own verifier comparison.
- Let generated authentication contract shapes become runtime registries.
- Ask the maintainer for a decision before this boundary-only sequencing step.

## Rationale

The application service boundary is the safest bridge between generated semantic contract metadata and future handwritten authentication behavior. It gives future agents a narrow target without creating behavior early.

Authentication service code has a high chance of drifting into token generation, verifier comparison, transport carriers, public wire shape, or persistence decisions if those responsibilities are not split before code appears.

Nakama and Pitaya reinforce the need for authentication, session, request identity, and handler context concepts in production game backends. vibit adapts those capabilities into stricter agent-native ownership boundaries rather than copying their APIs.

## Agent Reasoning Summary

The repository is ready to define the future service-interface shape, but it is not yet ready to implement authentication behavior. A boundary-only step keeps the project self-bootstrapping and gives later agents a stable intake surface for service code, tests, and protocol work.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  generated_contract_traceability: high
  security_redaction_clarity: high
  request_identity_handoff_clarity: high
  implementation_deferral: high
  protocol_stability: high
  repository_boundary_clarity: high
  tooling_cost: low
  long_term_maintainability: high
confidence: high
```

## Consequences

- `docs/application-authentication-service-interface-boundary.md` becomes the standard for future runtime authentication service-interface planning.
- `runtime.application_authentication_service_interface_boundary` becomes the repository check rule for this boundary.
- The application owner remains `runtime/internal/app`.
- Generated authentication shapes remain metadata-only and immutable.
- Runtime authentication behavior remains unimplemented.
- The next preparation gate can define verifier algorithms and redaction test expectations before service code is added.

## Reversal Conditions

Revisit this decision if:

- A future protocol decision requires authentication proof before protocol decoding.
- The application layer cannot obtain `authentication.Repository` through the unit-of-work boundary without weakening transaction ownership.
- A security review requires verifier algorithm selection before any service interface vocabulary can be named.
- Generated authentication contract shapes prove unsuitable for guiding service DTO vocabulary.
- The maintainer chooses WebSocket handshake authentication as the first proof carrier.

## Follow-Up

- Define token and credential verifier algorithm posture and redaction test expectations before service code.
- Keep Go service implementation behind a later bounded work item.
- Keep Protobuf authentication messages and WebSocket proof carriers behind separate gates.
