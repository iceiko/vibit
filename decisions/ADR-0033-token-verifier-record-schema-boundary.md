# ADR-0033: Token Verifier Record Schema Boundary

Status: Accepted
Date: 2026-05-14
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-14-define-token-verifier-record-schema-boundary/`

Related conversations:

- `conversations/2026-05-14-token-verifier-record-schema-boundary.md`

Related artifacts:

- `docs/token-verifier-record-schema-boundary.md`
- `docs/token-verifier-record-schema-boundary.zh-CN.md`
- `docs/credential-token-session-schema-gates.md`
- `docs/token-lifecycle-storage-implications.md`
- `docs/selected-login-token-boundary-checks.md`
- `docs/credential-record-schema-boundary.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`

## Context

M-014 exists to ratify credential and token verifier record schema boundaries before migrations, repositories, adapters, runtime lookup, handlers, routes, generated authentication shapes, Protobuf messages, WebSocket proof carriers, or authentication behavior are added.

W-0074 ratified the credential record schema boundary for `device_credential_login`. The selected first token posture is an opaque high-entropy access token with login-command issuance, explicit request proof payloads, no refresh token in the first posture, and PostgreSQL as the default durable target.

W-0075 turns the token verifier record gate into a concrete schema boundary while preserving implementation deferral.

## Decision

Ratify the token verifier record schema boundary for the first opaque access-token posture.

The future logical table is:

```text
authentication_access_tokens
```

The owner is:

```text
runtime.authentication
```

The boundary ratifies:

- One token verifier record represents one opaque access token.
- `token_record_id` is the log-safe token identifier.
- Raw access-token text is never stored.
- `token_lookup_digest` and `token_verifier_digest` are required non-plaintext verifier materials and are not log-safe.
- Token verifier algorithms must be versioned.
- Token lifecycle states are `active`, `expired`, and `revoked`.
- `active` tokens become invalid after `expires_at` even if expiration is computed rather than eagerly materialized.
- The first actor kind is `player`.
- `player_id` is immutable on token verifier records.
- `credential_record_id` linkage is required for the first posture so future login rotation can revoke previous active tokens for the same credential.
- Tokens cannot move between players or credentials.
- Refresh-token storage, runtime session binding, WebSocket connection binding, Redis-like token/session stores, and service/admin token families remain deferred.
- The first retention recommendation is seven days for expired and revoked token verifier records.
- Cleanup is required before production token storage is enabled, but cleanup implementation remains deferred.
- Replay-sensitive failure classes must distinguish missing, malformed, unsupported, invalid, expired, revoked, disabled-actor, and validator-unavailable cases without leaking whether player, credential, or token records exist beyond the ratified error surface.

This decision does not add migration source, tables, repository interfaces, PostgreSQL adapters, runtime token issuance, runtime token validation, logout, refresh, cleanup jobs, generated authentication shapes, Protobuf messages, WebSocket proof carriers, WebSocket handshake authentication, authentication dependencies, or production authentication behavior.

## Alternatives Considered

- Add token verifier migration source immediately.
- Store raw access tokens and compare plaintext.
- Use JWT as the first token format.
- Use a Redis-like token/session store before PostgreSQL token schema gates are complete.
- Store token state in `player_accounts`.
- Treat current Protobuf `Session` fields as token proof.
- Bind the first token verifier schema to WebSocket connection state.
- Add refresh-token storage in the first posture.
- Copy Nakama session token and refresh API shape.
- Treat Pitaya handler session binding as token verifier persistence.

## Rationale

Opaque access tokens require server-side verifier storage if revocation, logout, rotation, and replay-sensitive failure behavior are production requirements. The verifier record must be explicit before migrations so later agents do not hide token state in player lifecycle tables, credential records, WebSocket session state, or transport handlers.

The future logical table name is specific to access tokens because the first selected token family is specific. A generic session/token table would be more flexible, but it would also invite agents to mix refresh tokens, runtime sessions, service tokens, WebSocket bindings, and audit state before those families are ratified.

The credential linkage is required because the selected login posture rotates access tokens on successful login. It is still only linkage; it is not proof and does not authorize credential lookup or token validation.

Nakama remains useful as proof that mature game backends need token expiration, refresh, revocation, and logout capability, but vibit keeps refresh tokens and public API compatibility deferred. Pitaya remains useful for session and handler vocabulary, but token verifier records must not become transport sessions or handler context state.

## Agent Reasoning Summary

The correct next step is to ratify token verifier record semantics without creating storage. This gives later agents enough structure to plan migrations, repositories, adapters, redaction tests, and token validation behavior safely while keeping runtime authentication blocked until a later implementation milestone explicitly authorizes it.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  security_boundary_clarity: high
  schema_first_discipline: high
  player_account_separation: high
  credential_token_linkage_clarity: high
  implementation_deferral: high
  game_backend_reference_alignment: medium
  future_refresh_session_flexibility: medium
  long_term_maintainability: high
confidence: high
```

## Consequences

- The token verifier record schema boundary is ratified with no schema added.
- Future migration work has a stable logical table target.
- Player account lifecycle tables remain token-free.
- Credential record schema and token verifier schema boundaries are both ratified for M-014.
- Authentication schema migration planning becomes the next active work.
- Token verifier migrations, repository interfaces, PostgreSQL adapters, runtime token issuance, validation, logout, cleanup, generated output, Protobuf changes, WebSocket changes, authentication dependencies, and authentication behavior remain deferred.

## Reversal Conditions

Revisit this decision if:

- A security review requires a Redis-like or distributed revocation store before PostgreSQL token verifier storage.
- The maintainer ratifies refresh tokens in the first authentication implementation.
- A future session persistence decision turns access-token validation into a WebSocket-bound session model.
- A future compatibility ADR explicitly adopts a Nakama-like session token API surface.
- A future distributed runtime decision changes the durable token verifier target away from PostgreSQL.
- A future service/admin authentication family requires a shared token record design before player access-token storage is created.

## Follow-Up

- Plan the authentication schema migration queue in W-0076.
- Keep token verifier migration sources, repository interfaces, PostgreSQL adapters, runtime token validation, logout behavior, cleanup jobs, generated authentication output, Protobuf changes, WebSocket changes, authentication dependencies, and authentication implementation behind later gates.
