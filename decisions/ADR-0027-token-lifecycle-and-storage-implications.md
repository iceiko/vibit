# ADR-0027: Token Lifecycle And Storage Implications

Status: Accepted
Date: 2026-05-14
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-14-define-token-lifecycle-and-storage-implications/`

Related conversations:

- `conversations/2026-05-14-token-lifecycle-storage-implications.md`

Related artifacts:

- `docs/token-lifecycle-storage-implications.md`
- `docs/token-lifecycle-storage-implications.zh-CN.md`
- `docs/first-token-format-proof-carrier-posture.md`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`

## Context

ADR-0026 ratified opaque high-entropy access tokens as the first token format, with login command response issuance and explicit request proof payloads as the first carriers.

The next risk is that a future implementation agent may treat that ratification as permission to add raw token storage, refresh tokens, session token vocabulary, session tables, Redis-like dependencies, or ad hoc lifecycle behavior. W-0069 exists to define lifecycle and storage implications before contracts, schema gates, checks, and implementation.

## Decision

Ratify the first token lifecycle posture as:

```yaml
token_kind: access_token
format: opaque_high_entropy_token
minimum_entropy_bits: 256
token_text_encoding: url_safe_unpadded_base64_or_equivalent
token_ttl: 1h
refresh_token: not_in_first_implementation
renewal_method: reauthenticate_with_selected_login_method
revocation_required: true
rotation_required: true
replay_control_required: true
logout_required: true
cleanup_required: true
audit_required: true
raw_token_storage: forbidden
verifier_storage_required: true
token_storage_default_target: postgresql_schema_gate
session_storage_required_for_first_posture: false
external_identity_storage_required_for_first_posture: false
credential_storage_required_for_device_credential_login: true
implementation_authorized: false
```

The first access-token TTL is one hour. Refresh tokens are not included in the first implementation posture. Renewal happens by reauthenticating with the selected login method.

Token verifier storage is required before opaque token validation can be implemented. The future default durable target is PostgreSQL, but no table, migration, repository, or adapter is added by this decision.

Credential storage is required before `device_credential_login` can be implemented. External identity storage and session storage are not required for the first posture.

Player account lifecycle tables must remain credential-free, token-free, external-identity-free, and session-free.

## Alternatives Considered

- Very short access-token TTL, such as five or fifteen minutes.
- Long access-token TTL, such as twenty-four hours or longer.
- Adding refresh tokens in the first implementation.
- Treating the opaque token as a session token.
- Storing raw token strings.
- Using a Redis-like store for the first token verifier storage.
- Requiring persisted runtime sessions for the first token posture.
- Letting player account lifecycle tables carry credential or token state.

## Rationale

A one-hour access-token TTL is a balanced first default. It limits stolen-token lifetime while remaining practical for early gameplay loops. Because `device_credential_login` is the selected login method, renewal can happen through reauthentication without introducing refresh-token rotation and replay complexity immediately.

Opaque tokens require server-side verification. PostgreSQL is already the ratified authoritative durable store, so it is the correct default target for future token verifier schema gates. Redis-like storage may become useful later, but it is not justified before distributed runtime or high-scale session routing decisions exist.

Keeping session storage out of the first posture avoids conflating access tokens with persisted runtime sessions. This also preserves the current WebSocket transport and Protobuf envelope boundaries.

## Agent Reasoning Summary

The lifecycle posture keeps the first implementation production-minded without expanding into JWT, refresh-token, Redis-like, or session-persistence work. Future agents get concrete defaults for entropy, TTL, renewal, revocation, rotation, redaction, and storage gates while still being blocked from writing tables or runtime behavior prematurely.

## Decision Weights

```yaml
decision_weights:
  production_safety: high
  game_ergonomics: medium
  agent_context: high
  dependency_load: low
  storage_complexity: medium
  protocol_stability: high
  operations_clarity: high
  reversibility: high
  long_term_maintainability: high
confidence: high
```

## Consequences

- Future first token implementation must use at least 256 bits of entropy.
- Future first access tokens expire after one hour unless a later decision supersedes the default.
- Refresh tokens remain deferred.
- Token verifier storage is required before implementation, but schema work is deferred to W-0071.
- Token storage should target PostgreSQL first unless a later storage decision supersedes it.
- Runtime session persistence remains deferred.
- Player account lifecycle storage must not absorb credential, token, external identity, or session state.
- Future checks must prevent shortcut implementations where statically possible.

## Reversal Conditions

Revisit this decision if:

- A security review requires a shorter or longer default TTL.
- The first client experience cannot tolerate reauthentication-based renewal.
- A later session persistence milestone ratifies a session token model.
- A distributed runtime decision justifies Redis-like token or session storage.
- A compatibility goal with Nakama, Pitaya, or another framework explicitly changes lifecycle semantics.

## Follow-Up

- Define authentication contract, error, permission, and audit surfaces.
- Define credential, token, and session schema gates.
- Add repository checks for selected login/token boundaries.
- Keep authentication implementation deferred until those gates are complete.
