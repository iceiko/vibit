# ADR-0035: Authentication PostgreSQL Adapter Implementation Gate

Status: Accepted
Date: 2026-05-14
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-14-close-credential-token-verifier-schema-ratification-milestone/`

Related artifacts:

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `docs/authentication-schema-migration-queue.md`
- `docs/postgresql-persistence-boundary.md`
- `modules/authentication/module.yaml`
- `runtime/internal/modules/authentication/repository.go`

## Context

`M-014` ratified the credential and token verifier schema path for the selected `device_credential_login` and opaque access-token posture.

The repository now has:

- Credential record schema boundary.
- Token verifier record schema boundary.
- SQL-first PostgreSQL migration sources for `authentication_device_credentials` and `authentication_access_tokens`.
- Static migration checks.
- A storage-neutral `authentication.Repository` interface.
- A declared authentication PostgreSQL adapter boundary.

The next risk is that an agent treats the existence of token and credential persistence as permission to implement production authentication, login, token validation, WebSocket carriers, Protobuf messages, or token generation.

## Decision

Close `M-014` and open `M-015 Authentication PostgreSQL Adapter Implementation`.

`M-015` authorizes only persistence-adapter implementation for the already-ratified authentication repository interface:

```text
runtime/internal/platform/persistence/postgres/authentication_repository.go
runtime/internal/platform/persistence/postgres/authentication_repository_test.go
runtime/internal/platform/persistence/postgres.UnitOfWork.NewAuthenticationRepository
```

The adapter may insert, read, and update the ratified credential and token verifier records through a caller-supplied executor. It must not generate tokens, create credential material, compare verifiers, parse bearer tokens, validate access tokens, execute login, execute logout semantics beyond a repository update method, refresh tokens, run cleanup jobs, know WebSocket behavior, know Protobuf behavior, or add authentication dependencies.

The first work item in `M-015` refines repository checks before implementation so the check suite can distinguish allowed persistence-adapter vocabulary from forbidden runtime authentication behavior.

## Alternatives Considered

- Start runtime authentication immediately after schema ratification.
- Ask for a new maintainer confirmation gate before adapter implementation.
- Implement adapter and runtime login together.
- Keep all authentication work blocked until Protobuf and WebSocket proof-carrier decisions are ratified.
- Implement token generation before the persistence adapter.

## Rationale

The PostgreSQL adapter is the smallest useful next implementation step because the schema, repository interface, and adapter boundary are already ratified.

Implementing the adapter does not require deciding protocol shape, WebSocket proof carriers, raw token generation, verifier algorithms, handler routing, or production authentication behavior. It can be tested with fake executors and verified locally.

Refining checks before writing adapter code keeps the project agent-native: future agents receive machine-readable permission for exactly one new kind of code without weakening broader authentication boundaries.

## Agent Reasoning Summary

The next useful step is persistence implementation, not authentication behavior. Persistence-adapter code can make the ratified schema usable while preserving the stronger product decisions for later milestones.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  bounded_implementation: high
  security_boundary_clarity: high
  similarity_to_existing_player_persistence_pattern: high
  local_verifiability: high
  protocol_stability: high
  delivery_speed: medium
  future_authentication_readiness: high
confidence: high
```

## Consequences

- `M-014` is completed.
- `M-015` becomes the active milestone.
- `W-0083` becomes the next ready work item.
- Authentication adapter implementation can proceed only after checks are refined.
- Runtime authentication remains deferred.
- WebSocket and Protobuf behavior remain unchanged.
- Major authentication dependencies remain deferred.

## Reversal Conditions

Revisit this decision if:

- The repository interface is found to encode runtime authentication behavior rather than persistence.
- The migration schema cannot support the repository interface without a breaking schema change.
- A security review requires verifier algorithm ratification before any persistence adapter exists.
- The maintainer chooses to prioritize protocol proof carriers or generated authentication shapes before persistence adapter implementation.

## Follow-Up

- Advance `W-0083`.
- Then implement `W-0084` if checks pass.
- Keep login, token generation, token validation, logout execution, refresh, cleanup jobs, Protobuf, WebSocket, generated output, external identity, session persistence, audit persistence, and major authentication dependencies behind later gates.
