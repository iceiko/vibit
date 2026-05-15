# ADR-0046: Environment Verifier Key Loader Gate

Status: Accepted
Date: 2026-05-15
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-15-define-environment-verifier-key-loader-gate/`

Related conversations:

- `conversations/2026-05-15-environment-verifier-key-loader-gate.md`

Related artifacts:

- `docs/environment-verifier-key-loader-gate.md`
- `docs/environment-verifier-key-loader-gate.zh-CN.md`
- `docs/local-verifier-key-configuration-loading-gate.md`
- `runtime/internal/app/authentication/verifier_key_config.go`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`

## Context

The explicit in-memory verifier key set validator now exists under `runtime/internal/app/authentication`. It accepts already-decoded bytes and enforces the key-set invariants without parsing process environment variables or decoding text.

The next useful preparation step is to define the future process environment loader. Without a gate, a future agent could duplicate validation rules, accept ambiguous key encodings, leak environment variable values through errors, or wire loader behavior into startup while implementing parsing.

## Decision

Define the environment verifier key loader gate before adding code.

The future loader is an application-owned adapter under:

```text
runtime/internal/app/authentication
```

The future loader should use files such as:

```text
runtime/internal/app/authentication/verifier_key_env.go
runtime/internal/app/authentication/verifier_key_env_test.go
```

The future process environment contract uses:

```text
VIBIT_AUTH_VERIFIER_KEY_SET_ID
VIBIT_AUTH_CREDENTIAL_LOOKUP_KEY
VIBIT_AUTH_CREDENTIAL_VERIFIER_KEY
VIBIT_AUTH_TOKEN_LOOKUP_KEY
VIBIT_AUTH_TOKEN_VERIFIER_KEY
```

The four key variables must be Base64 text decoded into bytes and then passed to `NewVerifierKeySet`. The key set id is a required string and remains not log-safe by default.

The future loader must call the explicit in-memory validator instead of duplicating key-set validation logic. It may have a testable lookup-function entry point and a tiny process environment adapter after a later implementation work item authorizes code.

This decision does not implement Go code, environment parsing, Base64 decoding, startup wiring, local secret files, `.env` behavior, CLI secret input, KMS, cloud secret managers, token generation, credential generation, digest helpers, verifier comparison, authentication service behavior, protocol carriers, repository changes, migrations, dependencies, or production authentication behavior.

## Alternatives Considered

- Implement environment parsing immediately.
- Combine environment parsing and startup wiring.
- Let the future loader duplicate key validation rules.
- Accept raw key text, hex text, JSON blobs, or partial key sets.
- Add dotenv support in the first environment loader.
- Require KMS or a cloud secret manager before process environment loading is defined.

## Rationale

Environment loading is an adapter around the invariant-bearing validator. Keeping the loader second and gate-only first preserves a clean handoff: text source to decoded bytes to `NewVerifierKeySet`.

Process environment variables are practical for the first local posture on Termux and simple development environments, but their values are easy to leak through errors or shell history. A gate lets the project declare exact variable names, decoding policy, redaction expectations, and tests before code exists.

Nakama and Pitaya remain references for server-side authentication capability and realtime identity handoff. vibit adapts those ideas by keeping verifier key loading application-owned and away from transport, protocol, persistence, generated code, and domain modules.

## Agent Reasoning Summary

The repository is ready to define the environment loader contract because the in-memory validator is implemented and verified. It is not yet useful to wire startup or authentication behavior. A gate keeps future implementation narrow, testable, and redaction-focused.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  redaction_safety: high
  implementation_control: high
  testability: high
  dependency_minimization: high
  local_development_practicality: high
  operations_completeness: medium
confidence: high
```

## Consequences

- `docs/environment-verifier-key-loader-gate.md` becomes the standard for future process environment verifier key loader implementation.
- `runtime.environment_verifier_key_loader_gate` becomes the repository check rule for this gate.
- The next implementation slice can add a focused environment loader without choosing startup wiring.
- Process environment parsing and Base64 decoding remain unimplemented in this change.
- Runtime authentication behavior remains unimplemented.

## Reversal Conditions

Revisit this decision if:

- The project rejects process environment configuration as a first local posture.
- A future operations standard requires external secret-manager integration before any environment loader exists.
- The verifier key model changes away from four separated logical keys.
- A future Go packaging constraint makes `runtime/internal/app/authentication` unsuitable for loader helpers.

## Follow-Up

- Implement the environment verifier key loader as a narrow code slice.
- Keep startup wiring behind a separate composition gate.
- Keep `.env`, local secret files, CLI secret input, KMS, cloud secret managers, digest helpers, authentication service behavior, Protobuf messages, WebSocket proof carriers, repositories, migrations, dependencies, and production behavior behind later bounded work items.
