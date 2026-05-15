# ADR-0045: Local Verifier Key Configuration Loading Gate

Status: Accepted
Date: 2026-05-15
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-15-define-local-verifier-key-configuration-loading-gate/`

Related conversations:

- `conversations/2026-05-15-local-verifier-key-configuration-loading-gate.md`

Related artifacts:

- `docs/local-verifier-key-configuration-loading-gate.md`
- `docs/local-verifier-key-configuration-loading-gate.zh-CN.md`
- `docs/secret-configuration-verifier-key-loading-boundary.md`
- `docs/verifier-digest-computation-comparison-boundary.md`
- `docs/authentication-service-implementation-readiness-gate.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`

## Context

The authentication service implementation readiness gate names local verifier key configuration loading as the first recommended implementation gate. The previous secret configuration boundary allowed process environment configuration or explicit runtime secret input, but it intentionally did not choose the first implementation sequence.

The next risk is implementing environment parsing and key validation in the same code slice. That would make the invariant-bearing validator harder to test, couple the first implementation to process environment behavior, and create a path for secret values or key identifiers to leak through errors or logs.

## Decision

Define the local verifier key configuration loading gate before adding code.

The first future implementation slice should validate explicit in-memory verifier key configuration with already-decoded bytes. It should not parse environment variables, decode Base64 text, read local files, wire process startup, or integrate KMS or cloud secret managers.

The future owner package is:

```text
runtime/internal/app/authentication
```

The first future implementation should use files such as:

```text
runtime/internal/app/authentication/verifier_key_config.go
runtime/internal/app/authentication/verifier_key_config_test.go
```

The first validator must require a key-set id and four distinct logical key byte slices, each at least 32 bytes. It must copy inputs, avoid exposing mutable internal slices, reject duplicate logical key bytes, reject all-zero keys, reject obvious repeated single-byte keys, and keep errors redacted.

Environment loading is deferred to a follow-up gate that should call this explicit validator.

This decision does not implement Go code, environment parsing, secret loading, token generation, credential generation, digest helpers, verifier comparison, authentication service behavior, protocol carriers, repository changes, migrations, dependencies, or production authentication behavior.

## Alternatives Considered

- Implement process environment loading first.
- Implement environment loading and validation in one slice.
- Accept Base64-encoded strings as the first core validator input.
- Place key configuration loading in `runtime/cmd/vibit-server`.
- Place key configuration loading in PostgreSQL adapters or the authentication module repository.
- Require KMS or a cloud secret manager before any local key validation code.

## Rationale

Explicit in-memory validation is the smallest useful first slice. It directly captures the security invariants without forcing tests to manipulate process environment variables, local files, or shell encoding.

Environment loading is an adapter, not the invariant-bearing core. Sequencing it second makes future environment, file, KMS, or secret-manager loaders reuse the same validator instead of duplicating validation rules.

Keeping the code under `runtime/internal/app/authentication` preserves the previous application-owned authentication boundary. Transport, protocol, persistence, generated code, and domain modules should not hold verifier key material.

Nakama reinforces the need for trustworthy server-side authentication secret material. Pitaya reinforces that validated identity should arrive at handlers after proof validation, while transport/session mechanics should stay separate from authentication secrets. vibit adapts those lessons by validating verifier key configuration in the application layer before digest or service behavior exists.

## Agent Reasoning Summary

The repository is ready to define the first local key configuration code gate but not to implement it in the same change. Starting with an explicit in-memory validator keeps the next code slice narrow, testable on Termux without external services, and reusable by later environment or operations loaders.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  testability: high
  security_redaction_clarity: high
  implementation_control: high
  dependency_minimization: high
  local_development_practicality: high
  future_loader_reuse: high
  production_operations_completeness: medium
confidence: high
```

## Consequences

- `docs/local-verifier-key-configuration-loading-gate.md` becomes the standard for the first local verifier key configuration loading implementation gate.
- `runtime.local_verifier_key_configuration_loading_gate` becomes the repository check rule for this gate.
- The next implementation slice can be a focused explicit in-memory validator with unit tests.
- Environment variable parsing remains deferred to a follow-up gate.
- Runtime authentication behavior remains unimplemented.

## Reversal Conditions

Revisit this decision if:

- A future implementation proves that explicit in-memory validation cannot be separated cleanly from environment loading.
- Production deployment requirements require an external secret manager before local validation code exists.
- A Go packaging constraint makes `runtime/internal/app/authentication` unsuitable for application-owned authentication helpers.
- The verifier key model changes away from four separated logical keys.

## Follow-Up

- Implement the explicit in-memory verifier key set validator as a narrow code slice.
- Add focused tests for copying, immutability, missing values, short values, duplicate values, weak repeated values, and redaction.
- Define the environment variable loader gate after the explicit validator exists.
- Keep digest helpers, material generation, authentication service behavior, Protobuf messages, WebSocket proof carriers, repositories, migrations, dependencies, and production behavior behind later bounded work items.
