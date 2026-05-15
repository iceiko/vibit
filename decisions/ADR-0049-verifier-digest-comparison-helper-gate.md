# ADR-0049: Verifier Digest Comparison Helper Gate

Status: Accepted
Date: 2026-05-15
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-15-define-verifier-digest-comparison-helper-gate/`

Related conversations:

- `conversations/2026-05-15-verifier-digest-comparison-helper-gate.md`

Related artifacts:

- `docs/verifier-digest-comparison-helper-gate.md`
- `docs/verifier-digest-comparison-helper-gate.zh-CN.md`
- `docs/verifier-digest-computation-comparison-boundary.md`
- `docs/verifier-digest-helper-implementation-gate.md`
- `docs/token-credential-verifier-algorithm-redaction-boundary.md`
- `docs/authentication-service-implementation-readiness-gate.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`

## Context

The repository now has explicit verifier key set validation, environment key loading, token and credential material generation helpers, and verifier digest computation helpers. Verifier digest comparison remains deferred.

The next useful preparation step is to define the exact comparison helper implementation gate before adding comparison code. Without this gate, a future agent could compare verifier material with `bytes.Equal`, string equality, database-only equality, or the wrong digest class, or could combine comparison with authentication service behavior, repository calls, token validation, protocol carriers, or startup wiring.

## Decision

Define the verifier digest comparison helper gate before adding code.

Future implementation remains application-owned under:

```text
runtime/internal/app/authentication
```

The future helper files are:

```text
runtime/internal/app/authentication/verifier_comparison.go
runtime/internal/app/authentication/verifier_comparison_test.go
```

`verifier_digest.go` remains computation-only. The future comparison helper should compare a `ComputedDigest` verifier digest to stored verifier digest bytes, use `crypto/hmac.Equal` as the preferred constant-time primitive, accept `crypto/subtle.ConstantTimeCompare` only if the implementation preserves the same fail-closed length and redaction posture, return a redacted match or mismatch result, and reject missing, malformed, or wrong-class input without exposing digest material.

The future helper must not compare raw credential material, raw access-token material, encoded proof text, lookup digests, key ids, account ids, credential ids, token ids, session ids, or protocol metadata. It must not load keys, compute digests, select records, call repositories, inspect record lifecycle state, parse bearer proofs, issue login responses, validate tokens, revoke tokens, refresh tokens, touch protocol carriers, wire startup, add dependencies, or implement production authentication behavior.

This decision does not implement Go code, verifier comparison, digest computation, authentication service behavior, login execution, token validation, logout execution, refresh behavior, cleanup jobs, Protobuf messages, WebSocket proof carriers, startup wiring, repository changes, migrations, dependencies, or production authentication behavior.

## Alternatives Considered

- Implement comparison helpers immediately.
- Add comparison functions to `verifier_digest.go`.
- Compare verifier digests inside future authentication service methods only.
- Compare verifier digests inside `authentication.Repository` or the PostgreSQL adapter.
- Use `bytes.Equal`, `reflect.DeepEqual`, string equality, or database equality for verifier proof comparison.
- Compare lookup digest equality as authentication proof.
- Combine comparison with token validation or login execution.

## Rationale

Verifier comparison is small, but it is security-sensitive and easy for agents to get subtly wrong. A gate makes the next code slice narrow: compare verifier digest bytes only, use a constant-time primitive, preserve class checks, return redacted results, and stop.

Keeping comparison in a separate file from digest computation makes each helper boundary easier to inspect. Computation answers "what digest should this material produce?" Comparison answers "does this computed verifier digest match the stored verifier digest?" Authentication service behavior remains responsible for selecting records, checking lifecycle state, and mapping internal results to public failures later.

Nakama reinforces that game backends need server-side credential and token validation. Pitaya reinforces that realtime route handlers should receive identity context after validation, not own validation internals. vibit adapts those references by keeping comparison application-owned and separately gated from transport, protocol, repository, and service orchestration.

## Agent Reasoning Summary

The repository is ready to define the comparison helper implementation shape, but not to add comparison behavior in this change. The gate protects the project from future agents accidentally turning digest comparison into authentication service logic or using non-constant-time equality.

## Decision Weights

```yaml
decision_weights:
  constant_time_security: high
  agent_context: high
  helper_boundary_clarity: high
  redaction_safety: high
  implementation_control: high
  testability: high
  dependency_minimization: high
  future_service_composability: high
  operations_completeness: medium
confidence: high
```

## Consequences

- `docs/verifier-digest-comparison-helper-gate.md` becomes the standard for the future comparison helper implementation slice.
- `runtime.verifier_digest_comparison_helper_gate` becomes the repository check rule for this gate.
- The next implementation work can add `verifier_comparison.go` and focused tests without choosing service wiring or protocol carriers.
- Verifier comparison code remains unimplemented in this change.
- Runtime authentication behavior remains unimplemented.

## Reversal Conditions

Revisit this decision if:

- A security review requires a different constant-time comparison primitive or result model.
- The first access-token posture changes from opaque bearer material to a signed or structured claim token.
- The first device credential posture changes from server-issued high-entropy material to client-generated or password-like material.
- Production operations require an external cryptography dependency before local helper code exists.
- A future package layout decision moves application authentication helpers out of `runtime/internal/app/authentication`.

## Follow-Up

- Implement the verifier digest comparison helpers as a narrow code slice.
- Keep authentication service behavior, Protobuf messages, WebSocket proof carriers, startup wiring, repositories, migrations, dependencies, and production behavior behind later bounded work items.
