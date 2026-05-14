# First Token Format And Proof Carrier Posture

Status: Draft v0.1
Last updated: 2026-05-14
Scope: Ratified first token format and proof carrier posture for M-013
Depends on: `docs/token-format-carrier-options.md`
Canonical decision: `ADR-0026`

The paired Simplified Chinese translation is `docs/first-token-format-proof-carrier-posture.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This document ratifies vibit's first token format and proof carrier posture after `device_credential_login` was selected as the first login-method set.

It selects the first access-token format, issuance carrier, request proof carrier, and non-authorized boundaries before authentication implementation begins.

It does not implement token generation, parsing, signing, validation, refresh, revocation, rotation, replay handling, storage, session persistence, Protobuf envelope changes, WebSocket handshake authentication, runtime player handlers, or WebSocket routes.

## 2. Ratified Posture

The first token posture is:

```yaml
first_access_token_format: opaque_high_entropy_token
token_issuance_carrier: login_command_response_token
request_proof_carrier: explicit_request_proof_payload
refresh_token: deferred
session_token_vocabulary: deferred_until_session_persistence
protobuf_envelope_change: false
websocket_handshake_authentication_change: false
current_session_metadata_as_proof: false
first_system_message_binding: deferred
implementation_authorized: false
```

The selected token is an access token. It is not a refresh token, persisted runtime session identifier, WebSocket connection identifier, Protobuf `Session.session_id`, player ID, or external provider token.

## 3. Access Token Format

### `opaque_high_entropy_token`

Definition:

```text
A client-presented bearer secret with high entropy and no client-readable claims.
```

Ratified posture:

```yaml
format: opaque_high_entropy_token
token_kind: access_token
issuer: future_application_owned_authentication_boundary
verifier: future_application_owned_token_validator_before_domain_dispatch
subject: player_account_id_after_credential_and_account_policy_success
audience: vibit_gameplay_runtime_requests
expiration: required_finite_expiration_exact_policy_deferred_to_W_0069
refresh: deferred
revocation: required_capability_policy_deferred_to_W_0069
rotation: required_for_new_issuance_policy_deferred_to_W_0069
replay_posture: bearer_token_risk_must_be_controlled_by_lifecycle_policy
redaction: raw_token_secret_redacted_everywhere_except_client_presentation
storage: lookup_safe_hash_or_equivalent_non_plaintext_verifier_required_before_implementation
requires_signing_dependency: false
requires_key_management: false
requires_protobuf_envelope_change: false
requires_websocket_handshake_authentication: false
confidence: high
```

The token has no client-inspectable claims. Authorization facts must not be hidden inside token contents. Module-owned permission logic remains separate from token validation.

## 4. Issuer And Verifier

The future issuer is an application-owned authentication boundary reached after protocol decoding and after `device_credential_login` credential validation succeeds.

The future verifier is an application-owned token validator that runs before production-sensitive domain dispatch.

The following layers must not issue or verify the token:

- WebSocket transport adapters.
- Protobuf envelope adapters.
- Domain modules such as inventory.
- Player account persistence repositories.
- Generated Protobuf files or generated contract shape files.

This preserves the runtime boundary used elsewhere in vibit: transport and protocol move bytes and envelopes, application dispatch owns validation handoff, and domain modules receive already-normalized request identity.

## 5. Subject And Audience

The token subject is a player account identifier only after the selected login method and account creation or lookup policy succeed.

The initial audience is:

```text
vibit gameplay runtime requests
```

The first token posture does not authorize:

- Service-to-service authority.
- Admin authority.
- External provider sessions.
- Cross-game or cross-project audience sharing.
- Long-lived offline credentials.

Those require later contracts and decisions.

## 6. Expiration, Refresh, Revocation, Rotation, And Replay

The first access token must have a finite expiration. The exact expiration duration is deferred to W-0069.

Refresh tokens are deferred. A future refresh token must not be added without rotation, revocation, replay, storage, cleanup, redaction, and error semantics.

Revocation is required as a capability because opaque tokens are storage-backed or verifier-backed. The exact revocation model is deferred to W-0069.

Rotation is required for new issuance and replacement policy, but exact rotation triggers are deferred to W-0069.

Opaque access tokens are bearer secrets. A stolen token can be replayed until expiration or revocation unless later binding is ratified. W-0069 must define replay controls before implementation.

This ratification does not bind tokens to WebSocket connections, first system messages, device fingerprints, IP addresses, or current Protobuf session metadata.

## 7. Redaction And Storage

Raw token values are secret material.

Rules:

- Raw token values must be presented only by clients.
- Raw token values must not be stored in plaintext.
- Server-side storage, if used, must store a lookup-safe hash or equivalent non-plaintext verifier.
- Logs, errors, traces, conversation logs, change specs, tests, and documentation must redact token values.
- Token values must not appear in route names, request IDs, target IDs, player IDs, session IDs, connection IDs, or migration fixtures.

W-0069 and W-0071 must define the concrete lifecycle and schema gates before storage exists.

## 8. Proof Carrier

### Issuance Carrier

The first token should be issued by a future successful login command response:

```yaml
token_issuance_carrier: login_command_response_token
```

This means a future semantic `device_credential_login` command may return an access token after credential validation and account policy succeed.

This does not authorize the command contract yet. W-0070 must define the semantic command, response, errors, permissions, and audit surface before runtime implementation.

### Request Proof Carrier

The first request proof carrier for authenticated routes is:

```yaml
request_proof_carrier: explicit_request_proof_payload
```

Authenticated commands or queries may carry access-token proof in explicit contract-owned payload fields until a later protocol decision ratifies a cleaner carrier.

This is intentionally verbose. The benefit is that future agents can see the proof requirement in semantic contracts without changing the current Protobuf envelope or WebSocket handshake.

## 9. Unchanged Protocol And Transport Behavior

This ratification does not change the Protobuf envelope.

Current Protobuf session metadata remains metadata-only:

- `Session.session_id`
- `Session.player_id`
- `Session.connection_id`
- `Session.connection_epoch`

These fields are not authenticated proof and must not satisfy production permissions by reinterpretation.

This ratification does not change WebSocket handshake authentication. The WebSocket transport remains credential-neutral.

This ratification does not add a first system-message binding step. Connection binding, reconnect behavior, and session persistence remain future decisions.

## 10. Rejected And Deferred Options

| Option | Status | Reason |
| --- | --- | --- |
| Signed structured token | Deferred. | Requires signing dependency, key management, claim ownership, key rotation, revocation, replay, and clock-skew decisions. |
| External provider token as vibit access token | Deferred. | Belongs to provider login and external identity linking, not first device credential login. |
| Plain session ID as secret | Rejected for first posture. | Blurs identifier and proof vocabulary and risks reinterpreting metadata-only session fields. |
| Current Protobuf `Session` metadata as proof | Rejected. | Existing fields are metadata-only and must not become authority by convenience. |
| Protobuf envelope extension | Deferred. | It is a compatibility-sensitive wire decision. |
| WebSocket handshake carrier | Deferred. | It risks putting authentication into transport and needs browser/non-browser carrier analysis. |
| First system-message binding | Deferred. | Requires system-message contracts, timeout behavior, reconnect rules, and connection state. |
| Refresh token | Deferred. | Requires lifecycle, rotation, revocation, replay, redaction, storage, and error semantics. |

## 11. Reference Alignment

### Nakama

Nakama remains the capability reference for session token, refresh token, expiration, and logout vocabulary.

vibit adapts the idea that a successful authentication exchange can issue a client-presented token, but it does not copy Nakama's token format, public API, refresh behavior, or realtime socket binding semantics.

### Pitaya

Pitaya remains a vocabulary reference for session binding and handler context.

vibit adapts the session-context idea by requiring a future application-owned validator to produce request identity before domain dispatch. It does not put token validation inside WebSocket acceptors or route handlers.

## 12. Required Gates Before Implementation

Before this posture can be implemented, the repository must have:

- Token lifecycle and storage implications defined by W-0069.
- Authentication contract, error, permission, and audit surfaces defined by W-0070.
- Credential, token, and session schema gates defined by W-0071.
- Repository checks for selected login/token boundaries added by W-0072.
- A semantic login command and response contract.
- A semantic authenticated-request proof contract shape.
- A redaction rule for credentials and tokens.
- A token verifier storage or equivalent lookup boundary if validation requires durable lookup.
- Focused tests for success, missing proof, malformed proof, invalid proof, expired proof, revoked proof, replay/collision behavior, redaction, and layer ownership.

## 13. Non-Authorization

This ratification does not authorize:

- Runtime authentication code.
- Login handlers.
- Token generation, parsing, signing, validation, refresh, revocation, rotation, replay handling, or storage.
- Credential tables.
- External identity tables.
- Token tables.
- Session tables.
- Migrations.
- Password hashing, JWT, OAuth, OIDC, provider SDK, Redis-like, cryptography, key-management, or major authentication dependencies.
- Protobuf envelope changes.
- WebSocket handshake authentication.
- First system-message authentication.
- Runtime player account handlers.
- WebSocket routes.
- Treating metadata-only `player_id`, `session_id`, `connection_id`, or `connection_epoch` as proof.

## 14. Follow-Up

Next work:

```text
W-0069 Define token lifecycle and storage implications
```

W-0069 must define issuance, expiration, revocation, rotation, replay, cleanup, logout, storage, audit, and redaction implications before implementation.
