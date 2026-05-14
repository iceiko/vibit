# Token Lifecycle And Storage Implications

Status: Draft v0.1
Last updated: 2026-05-14
Scope: Token issuance, expiration, refresh, revocation, rotation, replay, logout, cleanup, redaction, audit, and storage implications for the first token posture
Depends on: `docs/first-token-format-proof-carrier-posture.md`
Canonical decision: `ADR-0027`

The paired Simplified Chinese translation is `docs/token-lifecycle-storage-implications.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This document defines the lifecycle and storage implications of vibit's first ratified token posture:

```yaml
first_access_token_format: opaque_high_entropy_token
token_issuance_carrier: login_command_response_token
request_proof_carrier: explicit_request_proof_payload
```

It turns the W-0068 token posture into implementation gates for future work.

This document does not implement token generation, parsing, validation, refresh, revocation, rotation, replay handling, storage, cleanup jobs, audit events, migrations, Protobuf changes, WebSocket handshake authentication, runtime player handlers, or WebSocket routes.

## 2. Lifecycle Summary

The first lifecycle posture is:

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

The `1h` access-token TTL is the first production-minded default for the initial implementation queue. It is short enough to limit stolen-token lifetime and long enough to avoid constant relogin for early gameplay loops. A later operations decision may tune this value through configuration, but the default must remain finite.

## 3. Issuance

A future successful `device_credential_login` command may issue an access token only after all of these conditions are true:

- The credential proof is valid.
- The player account creation or lookup policy succeeds.
- The account lifecycle state allows login.
- The token verifier storage boundary exists.
- Redaction and audit rules exist.
- The login command contract and response contract exist.

Issuance implications:

```yaml
issuer_owner: future_application_owned_authentication_boundary
issuer_layer: application_after_protocol_decode
transport_issuer: forbidden
protobuf_adapter_issuer: forbidden
domain_module_issuer: forbidden
player_account_repository_issuer: forbidden
generated_file_issuer: forbidden
```

The token response is a one-time client-visible secret presentation. Server logs, traces, errors, test snapshots, conversation logs, and change specs must not store the raw token.

## 4. Token Shape

The raw access token must be generated from cryptographically secure randomness with at least 256 bits of entropy.

The first acceptable text encoding is URL-safe unpadded Base64 or an equivalent encoding that avoids control characters, whitespace, path separators, query delimiters, and visually ambiguous formatting.

Rules:

- Token values are case-sensitive.
- Token values are bearer secrets.
- Token values must not contain client-readable claims.
- Token values must not embed `player_id`, `session_id`, provider subject, route name, timestamp, permission, or account lifecycle state.
- Token values must not be accepted from URL query parameters in the first posture.
- Token values must not be copied into Protobuf `Session` metadata fields.

## 5. Expiration

The first access-token TTL is:

```yaml
access_token_ttl: 1h
```

Rules:

- Tokens must have `issued_at` and `expires_at` semantics in the future verifier schema.
- Expiration is evaluated by the future application-owned token validator.
- Expired tokens produce a distinct expired-proof failure class.
- Expired tokens must not be silently treated as missing proof.
- Expired token records may be retained temporarily for audit, replay detection, or abuse analysis.

The exact retention period for expired token records remains a schema gate in W-0071. The first lifecycle recommendation is:

```yaml
expired_token_retention_recommendation: 7d
```

This is a recommendation, not a migration.

## 6. Refresh And Renewal

Refresh tokens are not part of the first implementation posture.

The first renewal method is:

```yaml
renewal_method: reauthenticate_with_device_credential_login
```

This means a client obtains a new access token by performing the selected login method again.

Rules:

- Do not add refresh-token contracts in W-0069.
- Do not add refresh-token storage in W-0069.
- Do not call the first access token a session token.
- Do not use current `Session.session_id` as refresh or renewal proof.
- A future refresh token requires its own rotation, revocation, replay, storage, cleanup, redaction, error, permission, and test gates.

## 7. Revocation And Logout

Revocation is required for the first opaque-token implementation.

Required future statuses:

```yaml
token_statuses:
  - active
  - expired
  - revoked
```

Future schema work may add more statuses, but those three are the minimum.

Logout semantics:

```yaml
logout_scope_first_posture: presented_access_token
logout_all_sessions: deferred
admin_revocation: deferred_to_permission_surface
forced_account_revocation: deferred_to_account_policy_and_audit_surface
```

The first logout behavior should revoke the presented access token. It must not revoke every token for a player, credential, device, or account unless a later contract and permission decision explicitly grants that behavior.

Revocation implications:

- A revoked token must fail validation distinctly from malformed or expired proof.
- Revocation must take effect before production-sensitive domain dispatch.
- Revocation must be visible to audit tooling when audit storage is ratified.
- Runtime memory-only revocation is insufficient for a production opaque-token posture unless the whole implementation is explicitly local-only.

## 8. Rotation

Rotation is required for new issuance.

First posture:

```yaml
rotation_on_successful_login: required
previous_token_for_same_credential: revoke_when_schema_supports_credential_token_linkage
automatic_background_rotation: deferred
refresh_rotation: deferred
```

The first implementation should rotate access tokens on successful login. Once schema gates define the relationship between credential proof and token verifier records, successful login should revoke previous active access tokens for the same credential installation unless a later decision explicitly supports concurrent active tokens for the same credential.

This does not require session persistence.

## 9. Replay Controls

Opaque access tokens are bearer secrets. Replay cannot be eliminated by token format alone.

Required first controls:

- High-entropy token generation.
- Finite TTL.
- Non-plaintext verifier storage.
- Revocation and logout.
- Rotation on successful login.
- Token redaction.
- Token carrier exclusion from route names, request IDs, target IDs, logs, URL query parameters, and Protobuf `Session` metadata.
- Future tests for replay-sensitive failure classes and stolen-token behavior within the defined model.

Deferred controls:

- Per-request nonce.
- Token binding to a WebSocket connection.
- Token binding to device fingerprint, IP address, TLS session, or first system message.
- Distributed replay cache.
- Redis-like token/session store.

These deferred controls require future architecture decisions because they affect protocol shape, state model, distributed runtime, and operational dependencies.

## 10. Cleanup

Cleanup is required before production token storage is enabled.

The first cleanup posture:

```yaml
cleanup_required: true
cleanup_owner: future_authentication_or_token_storage_boundary
cleanup_target: expired_and_revoked_token_verifier_records
cleanup_trigger_first_posture: explicit_maintenance_command_or_scheduled_runtime_job_deferred
default_retention_recommendation: 7d
```

No cleanup job is added by this document.

Future cleanup work must define:

- Whether cleanup runs through a CLI command, scheduled process, admin operation, or maintenance worker.
- Whether cleanup is safe to run concurrently.
- How cleanup is audited.
- Whether cleanup is required in default verification.
- How local development avoids destructive surprises.

## 11. Redaction

Raw token values are secret material.

Required redaction rule:

```yaml
redact_raw_tokens_in:
  - logs
  - errors
  - traces
  - metrics_labels
  - test_snapshots
  - migration_fixtures
  - conversation_logs
  - change_specs
  - documentation_examples
  - panic_or_recovery_output
```

Allowed references:

- Stable token record identifiers that are not derived from raw token text.
- Short redacted fingerprints, if a future standard defines the fingerprint algorithm.
- Hash/verifier values only in controlled database tests when raw tokens are absent and fixture policy is explicit.

Forbidden references:

- Raw token value.
- Token prefix long enough to be brute-force useful.
- Token embedded inside a URL.
- Token copied into `player_id`, `session_id`, `connection_id`, `request_id`, `target_id`, route name, or error message.

## 12. Audit Implications

Audit is required as a future capability, but no audit event is added by this document.

Future authentication audit surface should cover:

- Token issued.
- Token validation failed.
- Token expired.
- Token revoked.
- Token logout requested.
- Token rotated.
- Token cleanup executed.
- Credential mismatch or account state blocked issuance.

Audit events must not contain raw token values.

W-0070 must define public contract, error, permission, and audit surfaces before runtime implementation. W-0071 must define storage gates before any audit persistence exists.

## 13. Storage Implications

### Credential Storage

Credential storage is required before `device_credential_login` can be implemented.

Status:

```yaml
credential_storage_required: true
credential_storage_added_now: false
credential_storage_schema_gate: W-0071
```

Credential storage remains separate from player account lifecycle storage.

### Token Verifier Storage

Token verifier storage is required before opaque access-token validation can be implemented.

Status:

```yaml
token_verifier_storage_required: true
token_verifier_storage_added_now: false
token_verifier_schema_gate: W-0071
default_store_target: PostgreSQL
redis_like_store_selected: false
```

The first durable storage target should be PostgreSQL because PostgreSQL is already the ratified authoritative durable store. A Redis-like store remains deferred until dependency adoption and distributed runtime needs justify it.

The future token verifier schema must support:

- Non-plaintext token verifier.
- Token status.
- Subject actor.
- Audience.
- Issued-at timestamp.
- Expires-at timestamp.
- Revoked-at timestamp, when revoked.
- Rotation lineage or replacement relationship, if ratified.
- Credential-token linkage if one-active-token-per-credential is enforced.
- Audit-safe record identifiers.

No table, migration, repository, or adapter is added now.

### External Identity Storage

External identity storage is not required for the first `device_credential_login` posture.

Status:

```yaml
external_identity_storage_required_for_first_posture: false
external_identity_storage_added_now: false
```

Provider login and identity linking remain deferred.

### Session Storage

Session storage is not required for the first posture.

Status:

```yaml
session_storage_required_for_first_posture: false
session_storage_added_now: false
session_token_vocabulary: deferred_until_session_persistence
websocket_connection_binding: deferred
```

The first posture can validate access-token proof per authenticated request without persisted runtime sessions. A later session persistence milestone may choose a session store, connection binding, first-message authentication, handshake authentication, or hybrid model.

### Player Account Lifecycle Storage

Player account lifecycle storage must remain credential-free, token-free, external-identity-free, and session-free.

The current lifecycle tables remain:

```text
player_accounts
player_account_events
```

Forbidden by this lifecycle standard:

- Adding credential columns to `player_accounts`.
- Adding token columns to `player_accounts`.
- Adding provider subject columns to `player_accounts`.
- Adding session or WebSocket state columns to `player_accounts`.
- Adding raw token, token verifier, credential, provider subject, session, or WebSocket state rows to `player_account_events`.

## 14. Reference Alignment

### Nakama

Nakama remains the capability reference for access/session token issuance, refresh, expiration, logout, and revocation vocabulary.

vibit adapts those lifecycle dimensions but keeps refresh tokens and session token vocabulary deferred. The first vibit posture is a storage-backed opaque access token renewed by the selected login method, not a direct copy of Nakama's token/session behavior.

### Pitaya

Pitaya remains the vocabulary reference for session context and connection/session binding.

vibit defers connection-bound session behavior. The first lifecycle posture keeps token validation application-owned and does not place token state inside WebSocket acceptors or route handlers.

## 15. Required Future Gates

Before implementation, future work must provide:

- W-0070 contract, error, permission, and audit surfaces.
- W-0071 credential, token verifier, and optional session schema gates.
- W-0072 repository checks for forbidden shortcuts.
- A token verifier repository boundary.
- A credential lookup boundary.
- Redaction tests.
- Expiration tests.
- Revocation tests.
- Logout tests.
- Rotation tests.
- Replay-sensitive tests.
- Cleanup tests or explicit cleanup deferral with production classification.

## 16. Non-Authorization

This document does not authorize:

- Token generation code.
- Token parsing or validation code.
- Token storage tables.
- Credential storage tables.
- External identity tables.
- Session tables.
- Migrations.
- Cleanup jobs.
- Audit event implementation.
- Refresh tokens.
- JWT, signing, key-management, OAuth, OIDC, provider SDK, Redis-like, or password-hashing dependencies.
- Protobuf envelope changes.
- WebSocket handshake authentication.
- First system-message binding.
- Runtime player handlers.
- WebSocket routes.

## 17. Follow-Up

Next work:

```text
W-0070 Define authentication contract error permission surfaces
```

W-0070 must define the planned semantic contract, error, permission, and audit surfaces needed by the selected login and token posture before implementation.
