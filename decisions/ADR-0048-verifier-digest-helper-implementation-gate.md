# ADR-0048: Verifier Digest Helper Implementation Gate

Status: Accepted
Date: 2026-05-15
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-15-define-verifier-digest-helper-implementation-gate/`

Related conversations:

- `conversations/2026-05-15-verifier-digest-helper-implementation-gate.md`

Related artifacts:

- `docs/verifier-digest-helper-implementation-gate.md`
- `docs/verifier-digest-helper-implementation-gate.zh-CN.md`
- `docs/verifier-digest-computation-comparison-boundary.md`
- `docs/token-credential-verifier-algorithm-redaction-boundary.md`
- `docs/token-credential-material-generation-implementation-gate.md`
- `docs/authentication-service-implementation-readiness-gate.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`

## Context

The repository now has validated verifier key set loading, environment key loaders, and raw material generation helpers. The verifier digest computation and comparison boundary is already defined, but implementation is still deferred.

The next useful preparation step is to define the exact helper implementation gate before adding digest computation code. Without this gate, a future agent could compute HMAC digests in service orchestration, repositories, transport handlers, protocol adapters, migrations, or fixtures, or could combine digest computation with verifier comparison or authentication behavior.

## Decision

Define the verifier digest helper implementation gate before adding code.

Future implementation remains application-owned under:

```text
runtime/internal/app/authentication
```

The future helper files are:

```text
runtime/internal/app/authentication/verifier_digest.go
runtime/internal/app/authentication/verifier_digest_test.go
```

The future helper slice should build a deterministic canonical byte input with the registered version header, null separator, length-prefixed purpose label, and length-prefixed raw material; compute HMAC-SHA-256 with the matching logical key from an already-validated `VerifierKeySet`; return a copied 32-byte digest; expose typed or sentinel errors; and keep all digest bytes, key values, and raw material out of error text and string formatting.

The future helper should accept an already-validated `VerifierKeySet` and decoded raw material bytes. It may use Go standard library `crypto/hmac`, `crypto/sha256`, and `encoding/binary` after a later implementation work item authorizes code. No external cryptography, password-hashing, JWT, OAuth, OIDC, provider, KMS, cloud secret-manager, or operations dependency is required for the first helper implementation.

This decision does not implement Go code, HMAC computation, digest helpers, verifier comparison, authentication service behavior, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, startup wiring, repository changes, migrations, dependencies, or production authentication behavior.

## Alternatives Considered

- Implement digest helpers immediately.
- Put digest computation inside future authentication service methods.
- Put digest computation inside `authentication.Repository` or the PostgreSQL adapter.
- Combine digest computation with verifier comparison in one helper.
- Use an external HMAC or cryptography dependency before local helper code needs one.
- Compute digests without a canonical versioned input.
- Use raw string concatenation instead of length-prefixed binary encoding for canonical input.

## Rationale

HMAC-SHA-256 digest computation is small enough to implement locally, but the canonical input construction is security-sensitive and must be unambiguous. A gate makes the future helper narrow: build deterministic input, compute HMAC, copy digest bytes, and return.

The canonical input uses a versioned ASCII header, null separator, and explicit length-prefixed fields to prevent domain-separation ambiguity. This design follows the same rationale documented in `ADR-0043` and `docs/verifier-digest-computation-comparison-boundary.md`.

Separating digest computation from verifier comparison keeps each slice independently verifiable. The comparison slice can later add `crypto/hmac.Equal` or `crypto/subtle.ConstantTimeCompare` after the digest helper is validated.

Nakama reinforces that game backends need server-side credential and token verification capability. Pitaya reinforces that realtime route handlers should receive validated identity context rather than own digest computation. vibit adapts those references by keeping digest computation application-owned and separately gated from transport, protocol, repository, comparison, and service behavior.

## Agent Reasoning Summary

The repository is ready to define the implementation shape for verifier digest helpers, but not to add computation behavior in this change. The gate keeps the next implementation slice small enough for agents to execute and verify without crossing into verifier comparison or authentication service behavior.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  secret_redaction_safety: high
  implementation_control: high
  canonical_input_determinism: high
  testability: high
  dependency_minimization: high
  future_service_composability: high
  operations_completeness: medium
confidence: high
```

## Consequences

- `docs/verifier-digest-helper-implementation-gate.md` becomes the standard for the future digest helper implementation slice.
- `runtime.verifier_digest_helper_implementation_gate` becomes the repository check rule for this gate.
- The next implementation work can add `verifier_digest.go` and focused tests without choosing service wiring or protocol carriers.
- Verifier digest computation code remains unimplemented in this change.
- Verifier comparison remains unimplemented.
- Runtime authentication behavior remains unimplemented.

## Reversal Conditions

Revisit this decision if:

- A security review requires a different HMAC algorithm, hash function, or canonical input shape.
- The first access-token posture changes from opaque bearer material to a signed or structured claim token.
- The first device credential posture changes from server-issued high-entropy material to client-generated or password-like material.
- Production operations require an external cryptography dependency before local helper code exists.
- A future package layout decision moves application authentication helpers out of `runtime/internal/app/authentication`.

## Follow-Up

- Implement the verifier digest computation helpers as a narrow code slice.
- Keep verifier comparison, authentication service behavior, Protobuf messages, WebSocket proof carriers, startup wiring, repositories, migrations, dependencies, and production behavior behind later bounded work items.
