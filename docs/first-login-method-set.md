# First Login Method Set

Status: Draft v0.1
Last updated: 2026-05-14
Scope: Ratified first production login-method set for M-013
Depends on: `docs/first-login-method-candidates.md`
Canonical decision: `ADR-0025`

The paired Simplified Chinese translation is `docs/first-login-method-set.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This document ratifies vibit's first login-method set before runtime authentication implementation begins.

It selects the first production-oriented player login method, records rejected alternatives, and defines the gates that must exist before any handler, credential table, token parser, session store, Protobuf change, WebSocket handshake authentication, runtime player route, or WebSocket route is added.

This document does not implement authentication.

## 2. Ratified Set

The first login-method set is:

```yaml
first_login_method_set:
  - device_credential_login
```

No other login method is part of the first set.

## 3. Selected Method

### `device_credential_login`

Definition:

```text
A player client proves possession of a high-entropy installation credential.
```

The credential is secret proof material. It is not a raw operating-system device ID, advertising ID, model identifier, connection ID, player ID, session ID, or other public metadata.

Ratified posture:

```yaml
method: device_credential_login
actor_kind_after_success: player
production_classification: production_capable_after_required_gates
bootstrap_only: false
local_development_only: false
creates_player_account: allowed_after_account_creation_policy
authenticates_existing_account: allowed_after_credential_lookup_policy
links_existing_account: deferred
recovers_account: deferred
upgrades_anonymous_account: deferred
requires_major_dependency_before_contracts: false
requires_credential_storage_before_implementation: true
requires_external_identity_linking: false
requires_password_hashing: false
requires_oauth_or_oidc: false
requires_provider_sdk: false
requires_protobuf_envelope_change: false
requires_websocket_handshake_authentication: false
confidence: high
```

The first implementation must keep the WebSocket transport credential-neutral. Credential proof must be handled by a future application-owned authentication boundary after protocol decoding and before production-sensitive domain dispatch.

Token format, access-token behavior, refresh behavior, token carrier behavior, runtime session persistence, and WebSocket connection binding remain separate decisions in W-0067 through W-0071.

## 4. Public Rationale

`device_credential_login` is the smallest production-minded first login method that fits vibit's current goals.

It gives a low-friction game entry path similar to the device-style capability found in mature game backends such as Nakama, but it adapts that capability into vibit terms:

- The proof is high-entropy secret material, not public metadata.
- The player account lifecycle remains separate from credential records.
- The WebSocket transport remains credential-neutral.
- The Protobuf `Session` metadata remains metadata-only.
- Runtime handlers do not receive authenticated identity until a validator produces a machine-readable result.
- Later login families remain possible without forcing provider dependencies or password workflows into the first slice.

This also preserves the Pitaya-inspired separation between realtime connection/session vocabulary and handler request context. vibit should eventually bind validated identity into request context, but the binding must come from application-owned validation, not from WebSocket acceptors or route handlers.

## 5. Rejected Alternatives

| Alternative | Decision | Reason |
| --- | --- | --- |
| `guest_anonymous_login` | Deferred. | It is useful for onboarding, but first production player-owned state should not depend on anonymous authority until expiration, upgrade, abuse, recovery, and permission limits are ratified. |
| `custom_id_login` | Deferred. | It is only safe with a trusted issuer. That requires service-auth, subject namespace, collision, linking, replay, and audit boundaries first. |
| `email_password_login` | Deferred. | It requires password hashing, recovery, reset flows, rate limiting, breach posture, and support workflows that are too broad for the first authentication slice. |
| `external_provider_login` | Deferred. | It is important for later platform identity, but provider validation, issuer/audience rules, account linking, conflicts, outages, and dependencies make it too broad now. |
| `service_authentication` | Deferred to a separate track. | It is not a player login method and should not be mixed into the first player-authentication slice. |
| Metadata-only `player_id` or `session_id` | Rejected. | Metadata-only values are not proof and must not satisfy production permissions. |
| Direct Nakama API compatibility | Rejected for now. | Nakama remains a capability reference, not a governing API shape. |
| Pitaya session binding as implementation API | Rejected for now. | Pitaya remains vocabulary input; vibit owns its contracts, manifests, and runtime handoff shape. |

## 6. Decision Weights

```yaml
decision_weights:
  production_safety: medium
  game_onboarding_ergonomics: high
  agent_context: high
  dependency_load: low
  storage_complexity: medium
  abuse_and_recovery_load: medium
  reversibility: high
  long_term_maintainability: high
confidence: high
```

The main reason to accept medium storage and abuse complexity is that the complexity stays local, contract-checkable, and deferrable behind explicit gates.

## 7. Required Gates Before Implementation

Before `device_credential_login` can be implemented, the repository must have:

- Token format and proof carrier posture ratified by W-0067 and W-0068.
- Token lifecycle and storage implications defined by W-0069.
- Authentication command, response, error, permission, and audit surfaces defined by W-0070.
- Credential, token, and session schema gates defined by W-0071.
- Repository checks for selected login/token boundaries added by W-0072.
- A semantic login contract for the selected method.
- A credential storage boundary that keeps credentials out of `player_accounts` and `player_account_events`.
- A credential hash and lookup rule.
- A player account creation or lookup policy.
- Redaction rules for credentials, tokens, logs, traces, errors, change specs, conversation logs, and tests.
- Focused tests for success, missing proof, malformed proof, invalid proof, replay or collision behavior, redaction, and boundary ownership.

Implementation remains deferred until those gates are complete or explicitly superseded by a later accepted decision.

## 8. Non-Authorization

This ratification does not authorize:

- Runtime authentication code.
- Login handlers.
- Token parsing, signing, validation, refresh, revocation, rotation, replay handling, or storage.
- Credential tables.
- External identity tables.
- Token tables.
- Session tables.
- Migrations.
- Password hashing, JWT, OAuth, OIDC, provider SDK, Redis-like, cryptography, key-management, or major authentication dependencies.
- Protobuf envelope changes.
- WebSocket handshake authentication.
- Runtime player account handlers.
- WebSocket routes.
- Treating metadata-only `player_id`, `session_id`, or `connection_id` as proof.

## 9. Known Gaps

The following questions are intentionally left for later M-013 work:

- Whether the initial installation credential is client-generated, server-issued, or both.
- Whether first login creates a player account by default or requires explicit create intent.
- Whether credential rotation is included in the first implementation.
- Whether account recovery exists in the first implementation.
- Whether rate limiting needs a separate store or can start with process-local limits.
- Which token format and token carrier are selected.
- Whether refresh tokens exist in the first production slice.
- Whether runtime sessions are persisted or derived from token validation.

These gaps are not blockers for ratifying the login-method set. They are blockers for implementation.

## 10. Reference Alignment

### Nakama

Nakama remains the capability reference for broad game authentication coverage, including device-style, email, custom identifier, provider login, sessions, refresh, and logout.

vibit adopts the low-friction device-style capability as the first player login direction, but does not copy Nakama's API shape or token/session semantics.

### Pitaya

Pitaya remains a vocabulary reference for session binding, handler context, routing, and realtime server structure.

vibit adapts the session-context idea into application-owned request identity. The selected login method must not place authentication inside WebSocket acceptors, transport handlers, protocol adapters, or domain handlers.

## 11. Follow-Up

Next work:

```text
W-0067 Compare token format and carrier options
```

W-0067 must compare token formats and carrier postures without assuming that selecting `device_credential_login` automatically selects JWT, opaque tokens, refresh tokens, Protobuf envelope changes, WebSocket handshake authentication, or session persistence.
