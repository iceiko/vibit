# ADR-0040: Token And Credential Verifier Algorithm Redaction Boundary

Status: Accepted
Date: 2026-05-15
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-15-define-token-credential-verifier-algorithm-redaction-boundary/`

Related conversations:

- `conversations/2026-05-15-token-credential-verifier-algorithm-redaction-boundary.md`

Related artifacts:

- `docs/token-credential-verifier-algorithm-redaction-boundary.md`
- `docs/token-credential-verifier-algorithm-redaction-boundary.zh-CN.md`
- `docs/credential-record-schema-boundary.md`
- `docs/token-verifier-record-schema-boundary.md`
- `docs/application-authentication-service-interface-boundary.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`

## Context

The selected first authentication posture uses high-entropy device credentials and opaque high-entropy access tokens. Credential and token verifier record schemas already contain lookup digest, verifier digest, verifier algorithm, verifier version, and verifier key id fields. Application authentication service interfaces are also defined as a future application-owned boundary.

The next risk is that a future agent may implement verifier comparison or token generation with ad hoc hashing, raw token storage, plaintext comparison, leaky errors, or an unnecessary dependency such as JWT, OAuth, OIDC, bcrypt, Argon2, KMS, or a provider SDK. Before code appears, agents need a specific verifier posture, digest classification model, constant-time comparison expectation, dependency posture, and redaction test standard.

## Decision

Define `vibit_hmac_sha256_v1` as the first planned verifier algorithm family for high-entropy device credential and opaque access-token verification.

The planned digest classes are:

- `credential_lookup_digest`
- `credential_verifier_digest`
- `token_lookup_digest`
- `token_verifier_digest`

Each digest class uses HMAC-SHA-256 with a dedicated purpose label and server-side key. Lookup digests are secret-adjacent index material. Verifier digests are secret verifier material and require constant-time comparison.

Raw access-token material and raw device credential material must have at least 256 bits of entropy. Token presentation remains URL-safe unpadded Base64 or equivalent. Raw credential material, raw token material, lookup digests, verifier digests, verifier keys, peppers, and full verifier key identifiers are not log-safe.

Future first-posture implementation may use Go standard library cryptographic primitives: `crypto/hmac`, `crypto/sha256`, `crypto/subtle`, `crypto/rand`, and `encoding/base64`. No external cryptography, password-hashing, JWT, OAuth, OIDC, KMS, provider, or Redis-like token/session dependency is required for this first high-entropy posture.

This decision does not implement token generation, credential generation, verifier comparison, login execution, token validation, logout execution, cleanup jobs, application authentication service code, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository interface changes, migration schema changes, secret configuration loading, or production authentication behavior.

## Alternatives Considered

- Keep the algorithm unspecified until service code is added.
- Store SHA-256 digests without HMAC.
- Use bcrypt or Argon2 for the first device credential.
- Use JWT or signed claim tokens instead of opaque tokens.
- Add KMS or cloud secret-manager integration before local verifier behavior.
- Let the repository or PostgreSQL adapter compute verifier digests.
- Ask the maintainer before choosing a standard-library HMAC posture.

## Rationale

HMAC-SHA-256 with high-entropy inputs is a simple, standard, inspectable posture that fits Go's standard library and avoids adding unnecessary dependencies before behavior exists.

Plain SHA-256 would make digest material more exposed if raw token or credential material leaks elsewhere. HMAC gives the server a secret verifier boundary while preserving efficient lookup digest indexing.

bcrypt and Argon2 are correct tools for password-like or low-entropy human input, but the first selected posture is high-entropy installation credentials and high-entropy opaque access tokens. If password-like input is added later, it needs its own credential boundary and dependency adoption record.

JWT, OAuth, OIDC, KMS, provider SDKs, and Redis-like stores may be useful in later product directions, but they would widen the dependency and operations surface before vibit has the smallest application-owned authentication slice.

Nakama reinforces the need for token lifecycle and account authentication capability. Pitaya reinforces identity context separation in request handling. vibit adapts those lessons by keeping verifier logic application-owned, redacted, versioned, and explicitly gated.

## Agent Reasoning Summary

The repository is ready to define exact verifier posture but not ready to add authentication behavior. Choosing a standard-library HMAC-SHA-256 posture is narrow enough to guide future code and checks while avoiding a major external dependency decision.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  security_redaction_clarity: high
  dependency_minimization: high
  implementation_deferral: high
  repository_boundary_clarity: high
  protocol_stability: high
  future_rotation_support: medium
  operations_complexity: low
  long_term_maintainability: high
confidence: high
```

## Consequences

- `docs/token-credential-verifier-algorithm-redaction-boundary.md` becomes the standard for future verifier algorithm and redaction-test planning.
- `runtime.token_credential_verifier_algorithm_redaction_boundary` becomes the repository check rule for this boundary.
- Future first-posture verifier code can use Go standard library cryptographic primitives without an external dependency adoption record.
- Future code must keep lookup digest, verifier digest, raw proof, server-side keys, and verifier key identifiers out of public artifacts.
- Runtime authentication behavior remains unimplemented.
- The next preparation gate can define secret configuration and verifier key loading before code, or move to a bounded implementation gate if the necessary configuration posture exists.

## Reversal Conditions

Revisit this decision if:

- A security review requires a different algorithm or key hierarchy.
- The first credential posture changes from high-entropy installation credential to password-like or low-entropy human input.
- Operational key rotation requires an external KMS before any verifier behavior is added.
- A future protocol decision changes token format from opaque high-entropy bearer token to a signed or structured claim token.
- The Go standard library primitives become insufficient for the selected production security posture.

## Follow-Up

- Define secret configuration and verifier key loading boundary before implementing verifier code.
- Keep token generation, credential generation, verifier comparison, login execution, token validation, logout execution, cleanup, Protobuf messages, and WebSocket proof carriers behind later bounded work items.
- Add redaction tests in the first code slice that handles credential proof, access-token proof, verifier digests, or authentication errors.
