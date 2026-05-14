# ADR-0041: Secret Configuration And Verifier Key Loading Boundary

Status: Accepted
Date: 2026-05-15
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-15-define-secret-configuration-verifier-key-loading-boundary/`

Related conversations:

- `conversations/2026-05-15-secret-configuration-verifier-key-loading-boundary.md`

Related artifacts:

- `docs/secret-configuration-verifier-key-loading-boundary.md`
- `docs/secret-configuration-verifier-key-loading-boundary.zh-CN.md`
- `docs/token-credential-verifier-algorithm-redaction-boundary.md`
- `docs/application-authentication-service-interface-boundary.md`
- `docs/runtime-authentication-implementation-boundary.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`

## Context

The first verifier algorithm and redaction boundary is defined. It selects `vibit_hmac_sha256_v1` for high-entropy device credentials and opaque access tokens, while preserving deferral for runtime authentication behavior.

The next risk is that future code may load verifier keys in an ad hoc place, reuse one key across lookup and verifier digests, commit default keys, disclose key identifiers in public artifacts, silently fall back to weak local secrets, or introduce a major KMS or provider dependency before the first local authentication slice needs it.

## Decision

Define the secret configuration and verifier key loading boundary before implementation.

The first local secret posture is:

- Future key loading is application-owned under `runtime/internal/app`, with an optional child package such as `runtime/internal/app/authentication`.
- The first local implementation may use process environment configuration or explicit runtime secret input after a later code gate.
- No external KMS or cloud secret-manager dependency is required for the first local posture.
- Any external secret manager, KMS provider, cloud SDK, or operations integration requires a separate dependency adoption record and operations boundary.
- Four logical server-side keys are required: credential lookup, credential verifier, token lookup, and token verifier.
- `verifier_key_id` identifies a logical key set, is stored internally with verifier records, and is not log-safe by default.
- Production behavior must fail closed for missing, malformed, too-short, duplicate, or incomplete key configuration.
- Development and tests may use explicit test-only fixtures or local ignored environment files, but committed production-like secret values and default production keys are forbidden.

This decision does not implement secret loading, environment parsing, token generation, credential generation, verifier comparison, login execution, token validation, logout execution, cleanup jobs, application authentication service code, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository interface changes, migration schema changes, KMS integration, secret-manager integration, or production authentication behavior.

## Alternatives Considered

- Require external KMS before any local verifier behavior.
- Add committed development defaults.
- Use one server-side key for all digest classes.
- Treat `verifier_key_id` as log-safe.
- Defer secret configuration until service code is written.
- Implement environment parsing immediately.

## Rationale

Process environment or explicit runtime secret input is sufficient for the first local posture and keeps the dependency surface small while preserving a clear path to production secret-management work later.

Four logical keys preserve the lookup/verifier and credential/token separation already established by the verifier algorithm boundary. That separation gives future agents a concrete rule to follow and prevents accidental all-purpose key reuse.

Fail-closed production behavior avoids a dangerous default path where an authentication service appears to run but cannot verify securely.

Deferring KMS and external secret managers is not a rejection of those systems. It keeps the first implementation local, inspectable, and dependency-light until production operations requirements are ratified.

Nakama reinforces that account authentication and token lifecycle must be production-capable. Pitaya reinforces that request/session identity should remain outside low-level transport code. vibit adapts those references by keeping secret configuration application-owned, explicit, redacted, and separately gated.

## Agent Reasoning Summary

The repository is ready to define secret configuration posture but not ready to add key loading code. A process-environment-first local posture is narrow enough for Termux and local development, while the ADR requirement for external secret managers prevents hidden operations dependencies from entering without review.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  security_redaction_clarity: high
  dependency_minimization: high
  implementation_deferral: high
  repository_boundary_clarity: high
  local_development_practicality: high
  production_secret_management_path: medium
  operations_complexity: low
  long_term_maintainability: high
confidence: high
```

## Consequences

- `docs/secret-configuration-verifier-key-loading-boundary.md` becomes the standard for future verifier key configuration and loading posture.
- `runtime.secret_configuration_verifier_key_loading_boundary` becomes the repository check rule for this boundary.
- Future local implementation can use process environment configuration or explicit runtime secret input without an external dependency adoption record.
- Future external secret-manager or KMS integration requires a dependency adoption record and operations boundary.
- Future code must keep verifier key values, encoded key values, environment variable values, and full concrete key identifiers out of public artifacts.
- Runtime authentication behavior remains unimplemented.

## Reversal Conditions

Revisit this decision if:

- A security review requires KMS or a secret manager before any verifier behavior is implemented.
- Production deployment requirements make process environment configuration unacceptable even for a first local posture.
- The first credential posture changes to a low-entropy or password-like credential family.
- Operational key rotation needs a different key hierarchy or key identifier model.
- A future protocol decision changes from opaque high-entropy access tokens to signed or structured claim tokens.

## Follow-Up

- Define token and credential material generation before implementation.
- Define verifier digest computation and constant-time comparison before implementation.
- Keep login execution, token validation, logout execution, cleanup, Protobuf messages, and WebSocket proof carriers behind later bounded work items.
- Add redaction tests in the first code slice that loads keys, computes digests, validates proofs, or produces authentication errors.
