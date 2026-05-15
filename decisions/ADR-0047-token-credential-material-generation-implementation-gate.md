# ADR-0047: Token And Credential Material Generation Implementation Gate

Status: Accepted
Date: 2026-05-15
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-15-define-token-credential-material-generation-implementation-gate/`

Related conversations:

- `conversations/2026-05-15-token-credential-material-generation-implementation-gate.md`

Related artifacts:

- `docs/token-credential-material-generation-implementation-gate.md`
- `docs/token-credential-material-generation-implementation-gate.zh-CN.md`
- `docs/token-credential-material-generation-boundary.md`
- `docs/verifier-digest-computation-comparison-boundary.md`
- `docs/authentication-service-implementation-readiness-gate.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`

## Context

The repository now has a validated local verifier key set shape and a narrow process environment verifier key loader. The raw token and credential material generation boundary is already defined, but implementation is still deferred.

The next useful preparation step is to define the exact helper implementation gate before adding generation code. Without this gate, a future agent could generate raw secrets in service orchestration, repositories, transport handlers, protocol adapters, migrations, or fixtures, or could combine generation with digest computation, verifier comparison, or authentication behavior.

## Decision

Define the token and credential material generation implementation gate before adding code.

Future implementation remains application-owned under:

```text
runtime/internal/app/authentication
```

The future helper files are:

```text
runtime/internal/app/authentication/material_generation.go
runtime/internal/app/authentication/material_generation_test.go
```

The future helper slice should generate 32 bytes of raw material, encode it with URL-safe unpadded Base64, preserve a distinct material kind for device credential versus access token, copy raw bytes on return, fail closed on missing or failing randomness, reject all-zero and repeated-single-byte material, and expose only redacted errors.

The future helper should accept an explicit `io.Reader` for testability and may use Go standard library `crypto/rand`, `encoding/base64`, and `io` after a later implementation work item authorizes code. No external randomness, cryptography, JWT, OAuth, OIDC, provider, KMS, cloud secret-manager, operations, or password-hashing dependency is required for the first helper implementation.

This decision does not implement Go code, token generation, credential generation, digest computation, verifier comparison, authentication service behavior, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, startup wiring, repository changes, migrations, dependencies, or production authentication behavior.

## Alternatives Considered

- Implement generation helpers immediately.
- Put generation inside future authentication service methods.
- Put generation inside `authentication.Repository` or the PostgreSQL adapter.
- Generate JWT-like structured token text.
- Store raw material temporarily for debugging.
- Use an external randomness or token-generation dependency before local helper code needs one.
- Combine generation with verifier digest computation in one helper.

## Rationale

Raw token and credential generation is small enough to implement locally, but secret material is easy to mishandle. A gate makes the future helper narrow: obtain bytes, validate minimal generated shape, encode for one-time presentation, and return a value object.

The explicit `io.Reader` seam gives tests deterministic control while keeping production randomness separate. URL-safe unpadded Base64 gives shell- and protocol-friendly presentation text without making the token a claim container.

Nakama reinforces that game backends need server-side credential and token lifecycle capability. Pitaya reinforces that realtime route handlers should receive validated identity context rather than own proof generation. vibit adapts those references by keeping raw material generation application-owned and separately gated from transport, protocol, repository, digest, and service behavior.

## Agent Reasoning Summary

The repository is ready to define the implementation shape for material generation helpers, but not to add generation behavior in this change. The gate keeps the next implementation slice small enough for agents to execute and verify without crossing into authentication service behavior.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  secret_redaction_safety: high
  implementation_control: high
  testability: high
  dependency_minimization: high
  future_service_composability: high
  operations_completeness: medium
confidence: high
```

## Consequences

- `docs/token-credential-material-generation-implementation-gate.md` becomes the standard for the future helper implementation slice.
- `runtime.token_credential_material_generation_implementation_gate` becomes the repository check rule for this gate.
- The next implementation work can add `material_generation.go` and focused tests without choosing service wiring or protocol carriers.
- Token and credential generation code remains unimplemented in this change.
- Runtime authentication behavior remains unimplemented.

## Reversal Conditions

Revisit this decision if:

- A security review requires a different raw byte size or encoding.
- The first access-token posture changes from opaque bearer material to a signed or structured claim token.
- The first device credential posture changes from server-issued high-entropy material to client-generated or password-like material.
- Production operations require an external randomness dependency before local helper code exists.
- A future package layout decision moves application authentication helpers out of `runtime/internal/app/authentication`.

## Follow-Up

- Implement the token and credential material generation helpers as a narrow code slice.
- Keep digest computation, verifier comparison, authentication service behavior, Protobuf messages, WebSocket proof carriers, startup wiring, repositories, migrations, dependencies, and production behavior behind later bounded work items.
