# ADR-0043: Verifier Digest Computation And Comparison Boundary

Status: Accepted
Date: 2026-05-15
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-15-define-verifier-digest-computation-comparison-boundary/`

Related conversations:

- `conversations/2026-05-15-verifier-digest-computation-comparison-boundary.md`

Related artifacts:

- `docs/verifier-digest-computation-comparison-boundary.md`
- `docs/verifier-digest-computation-comparison-boundary.zh-CN.md`
- `docs/token-credential-verifier-algorithm-redaction-boundary.md`
- `docs/secret-configuration-verifier-key-loading-boundary.md`
- `docs/token-credential-material-generation-boundary.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`

## Context

The first verifier algorithm family, secret configuration posture, and raw material generation posture are defined. Future authentication service code still lacks a precise boundary for canonical digest input construction, lookup digest handoff, key-set selection during rotation, and constant-time verifier digest comparison.

Without this boundary, a future agent may compute HMAC inputs inconsistently, use lookup digest equality as authentication proof, compare verifier material with non-constant-time primitives, expose lookup misses differently from verifier mismatches, or move digest logic into repositories and protocol adapters.

## Decision

Define verifier digest computation and comparison as a future application-owned boundary under `runtime/internal/app`, optionally in an application-owned child package such as `runtime/internal/app/authentication`.

The first posture keeps `vibit_hmac_sha256_v1`, version `1`, 32-byte raw digest output, four separated digest classes, four separated logical keys, and the registered purpose labels from ADR-0040.

Future digest helpers must build canonical input as:

```text
ascii("vibit.auth.verifier.input.v1")
|| 0x00
|| uint16be(len(purpose_label))
|| ascii(purpose_label)
|| uint16be(len(raw_material))
|| raw_material
```

Lookup digest database equality is allowed only for candidate record selection. It is not authentication proof. Future validation must still use lifecycle checks and constant-time comparison of computed verifier digest bytes against stored verifier digest bytes.

Future validation must select candidate lookup keys from active and accepted previous key sets because opaque proof text does not carry `verifier_key_id`. After a candidate record is selected, the stored `verifier_key_id` selects the verifier key set for computing the verifier digest. Missing records, verifier mismatches, unknown key ids, unsupported algorithm versions, malformed proof, and expired or revoked tokens collapse to the same public invalid-proof class unless a later semantic error standard explicitly allows more detail.

Future first-posture implementation may use Go standard library `crypto/hmac`, `crypto/sha256`, and `crypto/subtle` after a later code gate. No external cryptography, password-hashing, JWT, OAuth, OIDC, provider, KMS, cloud secret-manager, or operations dependency is required for this digest computation and comparison posture.

This decision does not implement verifier digest computation, verifier comparison, token generation, credential generation, secret loading, login execution, token validation, logout execution, cleanup jobs, application authentication service code, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository interface changes, migration schema changes, or production authentication behavior.

## Alternatives Considered

- Leave canonical digest input shape for implementation time.
- Use raw material directly as the HMAC input without a version header or length prefixes.
- Use database lookup digest equality as complete proof.
- Store a token key id or credential key id inside the opaque proof text.
- Compare digest bytes with `bytes.Equal`.
- Use a JWT or signed claim token to avoid server-side token verifier records.
- Add an external cryptography dependency before local HMAC helpers need one.

## Rationale

Canonical byte construction is a small rule with high long-term value. It lets future agents implement HMAC helpers without inventing incompatible input framing or silently changing digest behavior.

Length-prefixing the purpose label and raw material makes the input unambiguous and versioned while still staying simple enough for Go standard library implementation.

Opaque bearer tokens should not carry key ids or claims in the first posture. Computing lookup digest candidates across active and accepted previous key sets keeps token text opaque while preserving key rotation.

Database equality on lookup digests is useful for indexed record selection, but it is not proof. Keeping verifier digest comparison application-owned and constant-time prevents repositories from becoming hidden authentication engines.

Nakama reinforces the need for robust account/token validation semantics. Pitaya reinforces that request identity should be bound in application/session context rather than transport acceptors. vibit adapts those references by making digest computation, comparison, and public failure collapse explicit and agent-checkable.

## Agent Reasoning Summary

The repository is ready to define digest computation and comparison posture but not ready to add HMAC helpers. A boundary now gives future agents exact byte input, key-selection, comparison, failure, redaction, and dependency constraints before the first authentication service implementation slice.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  digest_determinism: high
  rotation_support: high
  redaction_clarity: high
  dependency_minimization: high
  implementation_deferral: high
  repository_boundary_integrity: high
  long_term_maintainability: high
confidence: high
```

## Consequences

- `docs/verifier-digest-computation-comparison-boundary.md` becomes the standard for future digest computation and verifier comparison.
- `runtime.verifier_digest_computation_comparison_boundary` becomes the repository check rule for this boundary.
- Future first-posture digest helpers can use Go standard library HMAC-SHA-256 and constant-time comparison primitives without an external dependency adoption record.
- Future validation may need repository boundary work if batch lookup by multiple digest candidates is selected.
- Runtime authentication behavior remains unimplemented.

## Reversal Conditions

Revisit this decision if:

- A security review requires a different canonical input shape.
- A future token format changes from opaque high-entropy bearer token to a structured claim token.
- Production operations require a different key-selection strategy for rotation.
- A future repository boundary proves that candidate lookup across accepted key sets is impractical.
- A new verifier algorithm replaces `vibit_hmac_sha256_v1`.

## Follow-Up

- Define an authentication service implementation readiness gate before adding service code.
- Keep verifier digest helpers, verifier comparison helpers, token generation, credential generation, secret loading, login execution, token validation, logout execution, cleanup, Protobuf messages, and WebSocket proof carriers behind later bounded work items.
- Add canonical input and constant-time comparison tests in the first code slice that computes verifier digests.
