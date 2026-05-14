# ADR-0042: Token And Credential Material Generation Boundary

Status: Accepted
Date: 2026-05-15
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-15-define-token-credential-material-generation-boundary/`

Related conversations:

- `conversations/2026-05-15-token-credential-material-generation-boundary.md`

Related artifacts:

- `docs/token-credential-material-generation-boundary.md`
- `docs/token-credential-material-generation-boundary.zh-CN.md`
- `docs/token-credential-verifier-algorithm-redaction-boundary.md`
- `docs/secret-configuration-verifier-key-loading-boundary.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`

## Context

The first verifier algorithm posture and secret configuration/key loading posture are now defined. Future authentication service code still lacks a boundary for generating raw device credential and raw access-token material.

Without a generation boundary, a future agent may generate low-entropy strings, embed claims in opaque tokens, store raw token or credential material, use client-supplied metadata as a credential, expose one-time secrets in logs, or add a major dependency before the first local posture needs it.

## Decision

Define token and credential material generation as a future application-owned boundary under `runtime/internal/app`, optionally in an application-owned child package such as `runtime/internal/app/authentication`.

The first posture is server-issued and application-generated for both device credentials and access tokens. Generated raw material must contain at least 256 bits of entropy, with 32 random bytes as the first raw byte shape. Text presentation is URL-safe unpadded Base64 or equivalent.

Raw device credentials and raw access tokens are one-time client-visible secrets. They must not be stored by repositories, PostgreSQL adapters, transport adapters, protocol adapters, domain modules, generated output, logs, traces, metrics labels, public errors, audit-safe facts, change specs, ADRs, documentation examples, or conversation logs.

Future first-posture implementation may use Go standard library `crypto/rand` and `encoding/base64` after a later code gate. No external randomness, cryptography, JWT, OAuth, OIDC, provider, KMS, cloud secret-manager, or operations dependency is required for this generation posture.

This decision does not implement token generation, credential generation, secret loading, verifier digest computation, verifier comparison, login execution, token validation, logout execution, cleanup jobs, application authentication service code, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository interface changes, migration schema changes, or production authentication behavior.

## Alternatives Considered

- Accept client-generated installation credentials in the first posture.
- Generate shorter local development tokens first.
- Use JWT or signed claim tokens for the first access token.
- Store raw token or credential material temporarily for debugging.
- Add an external randomness or secret-generation dependency.
- Defer generation posture until service code is written.

## Rationale

Server-issued application generation lets the first implementation guarantee entropy and encoding without trusting client platform randomness, device identifiers, or metadata. It also keeps generation out of transport and repository code.

Opaque high-entropy tokens should remain opaque. Embedding player ids, timestamps, claims, or route data would change the token model and pull in validation and signing concerns too early.

The Go standard library is enough for secure random bytes and text encoding in the first local posture. More advanced production operations can still adopt external dependencies later through explicit gates.

Nakama reinforces that game backends need token issuance and credential lifecycle capability. Pitaya reinforces that request identity should be bound in application/session context rather than transport acceptors. vibit adapts those references by making raw secret generation application-owned, one-time, redacted, and separately gated.

## Agent Reasoning Summary

The repository is ready to define material generation posture but not ready to add generation helpers. A server-issued first posture gives future agents concrete constraints for entropy, encoding, storage, response presentation, and repository handoff while keeping runtime behavior deferred.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  entropy_guarantee: high
  redaction_clarity: high
  dependency_minimization: high
  implementation_deferral: high
  client_complexity: medium
  future_rotation_support: medium
  long_term_maintainability: high
confidence: high
```

## Consequences

- `docs/token-credential-material-generation-boundary.md` becomes the standard for future raw token and credential material generation.
- `runtime.token_credential_material_generation_boundary` becomes the repository check rule for this boundary.
- The first credential source is server-issued application generation.
- The first access-token source is server-issued application generation.
- Future first-posture generation code can use Go standard library randomness and Base64 encoding without an external dependency adoption record.
- Runtime authentication behavior remains unimplemented.

## Reversal Conditions

Revisit this decision if:

- A client platform integration requires client-generated credential enrollment for the first posture.
- A security review requires a different raw byte size or encoding.
- A future token format changes from opaque high-entropy bearer token to a signed or structured claim token.
- Production operations require an external randomness or secret-generation dependency before local generation code exists.
- The first credential posture changes from high-entropy installation credential to password-like or low-entropy human input.

## Follow-Up

- Define verifier digest computation and constant-time comparison before implementation.
- Keep token generation helpers, credential generation helpers, login execution, token validation, logout execution, cleanup, Protobuf messages, and WebSocket proof carriers behind later bounded work items.
- Add generation redaction tests in the first code slice that produces raw credential or access-token material.
