# ADR-0036: Runtime Authentication Implementation Boundary

Status: Accepted
Date: 2026-05-14
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-14-define-runtime-authentication-implementation-boundary/`

Related conversations:

- `conversations/2026-05-14-runtime-authentication-implementation-boundary.md`

Related artifacts:

- `docs/runtime-authentication-implementation-boundary.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`
- `runtime/internal/modules/authentication/repository.go`
- `runtime/internal/platform/persistence/postgres/authentication_repository.go`

## Context

The authentication persistence path is now concrete enough to tempt a future agent into implementing login or token validation by convenience.

The repository already has semantic authentication contracts, selected `device_credential_login`, opaque access-token posture, credential and token verifier migrations, a storage-neutral `authentication.Repository`, and a PostgreSQL adapter. What it still does not have is production runtime authentication behavior.

The next architectural risk is not missing code. The risk is unbounded authentication code appearing in transport, protocol adapters, repositories, generated files, or domain modules.

## Decision

Define the runtime authentication implementation boundary before adding runtime behavior.

Future runtime authentication is application-owned under `runtime/internal/app`. The application boundary will orchestrate login, token validation, logout, error mapping, repository usage, and request identity handoff after separate implementation gates authorize those behaviors.

The authentication module remains the storage-neutral repository interface owner. The PostgreSQL adapter remains persistence-only. Protocol adapters may only translate ratified wire messages. WebSocket transport remains credential-neutral unless a future carrier decision changes that.

Token generation, verifier comparison, login execution, access-token validation, logout execution, cleanup jobs, generated authentication shapes, Protobuf messages, WebSocket proof carriers, and authentication dependencies remain separate gates.

## Alternatives Considered

- Start `AuthenticateWithDeviceCredential` implementation immediately after the PostgreSQL adapter.
- Put token validation in the Protobuf adapter.
- Put token proof in WebSocket handshake headers as the first implementation.
- Let the authentication repository compare verifiers.
- Generate authentication contract shapes before defining the runtime boundary.
- Ask for a new maintainer decision before any non-behavioral planning step.

## Rationale

Authentication has high drift risk. The code can appear to work while silently weakening module boundaries, request identity trust, error redaction, or transport neutrality.

An explicit implementation boundary gives future agents a smaller and safer intake path. They can see which layer owns each responsibility and which gates remain closed before they touch runtime behavior.

Nakama and Pitaya are useful references for capability coverage and vocabulary, but vibit needs stricter ownership rules because its primary goal is agent-native maintainability.

## Agent Reasoning Summary

The safest next step is to define the runtime implementation boundary, not to write authentication behavior. The persistence adapter makes later runtime work possible, but it does not change where authentication proof should live or how request identity becomes trustworthy.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  security_boundary_clarity: high
  implementation_cost: low
  reversibility: medium
  long_term_maintainability: high
  similarity_to_nakama_capability_coverage: medium
  similarity_to_pitaya_session_vocabulary: medium
  protocol_stability: high
confidence: high
```

## Consequences

- `docs/runtime-authentication-implementation-boundary.md` becomes the standard for first runtime authentication implementation planning.
- Runtime authentication behavior remains unimplemented.
- Application ownership is defined before service code appears.
- Repository usage remains through `authentication.Repository`.
- Generated authentication shapes, Protobuf authentication messages, WebSocket proof carriers, and authentication dependencies remain deferred.
- Future work items can advance in smaller gates without asking the maintainer for routine technical sequencing.

## Reversal Conditions

Revisit this decision if:

- A future protocol decision requires authentication proof before protocol decoding.
- The application boundary cannot validate proof before domain dispatch without breaking the runtime request loop.
- A security review requires verifier algorithm ratification before any application service boundary can be named.
- The maintainer chooses WebSocket handshake authentication as the first proof carrier.
- Nakama or Pitaya alignment reveals a required capability that cannot fit the current ownership split.

## Follow-Up

- Add or refine runtime authentication implementation boundary checks.
- Close `M-016` only after the boundary is inspectable in standards, manifests, guides, and work items.
- Keep runtime login, token generation, token validation, verifier comparison, logout execution, refresh, cleanup, Protobuf, WebSocket, generated output, runtime session persistence, audit persistence, and major authentication dependencies behind later gates.
