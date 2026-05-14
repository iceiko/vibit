# Token Format And Carrier Option Comparison

Status: Draft v0.1
Last updated: 2026-05-14
Scope: Token format and proof carrier comparison for M-013
Depends on: `docs/first-login-method-set.md`

The paired Simplified Chinese translation is `docs/token-format-carrier-options.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This document compares token format and proof carrier options after `device_credential_login` was ratified as the first login-method set.

It recommends a first token format and proof carrier posture for W-0068.

It does not implement token parsing, signing, validation, refresh, revocation, rotation, replay handling, token storage, session persistence, Protobuf envelope changes, WebSocket handshake authentication, runtime player handlers, or WebSocket routes.

## 2. Evaluation Criteria

Each option is evaluated against:

- Production safety.
- Revocation and logout ergonomics.
- Redaction safety.
- Agent-readable implementation shape.
- Dependency load.
- Storage and migration implications.
- Fit with `device_credential_login`.
- Fit with WebSocket-first gameplay.
- Nakama session capability alignment.
- Pitaya session/context vocabulary alignment.
- Reversibility.

## 3. Token Format Summary

| Format | Recommendation | Reason |
| --- | --- | --- |
| `opaque_high_entropy_token` | Recommend as first token format. | It keeps validation explicit and storage-backed, is easy to redact, supports revocation/logout cleanly, avoids key-management and JWT dependencies, and gives agents a narrow lookup-based implementation path. |
| `signed_structured_token` | Defer. | Useful later, but signing, key rotation, issuer/audience/claim drift, revocation, and replay controls are too broad for the first slice. |
| `external_provider_token` | Defer. | It belongs to future provider login and external identity linking, not the first device credential path. |
| `plain_session_id_as_secret` | Reject for first token format. | It blurs identifier and proof vocabulary and risks turning existing metadata-only `session_id` into authority. |

## 4. Proof Carrier Summary

| Carrier | Recommendation | Reason |
| --- | --- | --- |
| `login_command_response_token` | Recommend for issuing the first token after successful login. | It keeps the credential exchange and token issuance in an application/protocol command flow, not in the WebSocket handshake or envelope metadata. |
| `explicit_request_proof_payload` | Recommend as first request proof posture for authenticated routes. | It keeps proof visible to semantic contracts without changing the current Protobuf envelope. It is verbose but agent-readable and reversible. |
| `first_system_message_binding` | Defer. | Useful for future connection binding, but requires system-message contracts, timeout behavior, reconnect rules, and connection state. |
| `protobuf_envelope_extension` | Defer. | It may become cleaner later, but it is a compatibility-sensitive protocol change. |
| Current `Session.session_id` metadata | Reject for proof. | It is metadata-only today and must not become proof by reinterpretation. |
| WebSocket handshake carrier | Defer. | It can reject unauthenticated connections early, but it risks putting authentication into transport and requires separate browser/non-browser carrier analysis. |
| WebSocket subprotocol, cookie, or query parameter | Defer. | Each requires explicit transport risk analysis and should not be selected by convenience. |

## 5. Recommended First Posture

Recommended first token format:

```text
opaque_high_entropy_token
```

Recommended first proof carrier posture:

```yaml
token_issuance_carrier: login_command_response_token
request_proof_carrier: explicit_request_proof_payload
protobuf_envelope_change: false
websocket_handshake_authentication_change: false
current_session_metadata_as_proof: false
first_system_message_binding: deferred
```

This recommendation means:

- A future successful `device_credential_login` command may issue an opaque high-entropy access token.
- The raw token value must be presented only by clients and redacted everywhere else.
- Server-side storage, if used, must store a lookup-safe hash or equivalent non-plaintext verifier, not raw token strings.
- Authenticated gameplay requests should carry proof through future explicit semantic contract fields until a later protocol decision ratifies a cleaner carrier.
- Existing Protobuf `Session.session_id`, `Session.player_id`, `Session.connection_id`, and `Session.connection_epoch` remain metadata-only.
- WebSocket transport remains credential-neutral.

W-0068 must ratify or adjust this recommendation with issuer, verifier, subject, audience, expiration, refresh, revocation, rotation, replay, redaction, and storage implications.

## 6. Token Format Details

### Opaque High-Entropy Token

Position:

```text
recommend_first
```

Benefits:

- No signing or key-management dependency required for the first slice.
- Server-side validation can be explicit and contract-checkable.
- Revocation, logout, and forced invalidation are straightforward if token records exist.
- Redaction is simple because the token has no client-inspectable claims.
- Fits `device_credential_login`: login proves credential possession; token proves a subsequent session or access grant.
- Keeps claim evolution out of client-visible token contents.

Risks:

- Requires token lookup storage or an equivalent verifier.
- Requires careful token hashing and indexing rules.
- Requires cleanup and expiration logic.
- Adds database or session-store work before implementation.
- Every validation may hit storage unless caching or session binding is later ratified.

Required artifacts before implementation:

- Token issuance contract.
- Token validation contract.
- Token storage or verifier schema gate.
- Hash and lookup rule.
- Expiration rule.
- Revocation/logout rule.
- Redaction rule.
- Replay and collision tests.
- Repository checks preventing plaintext token storage and metadata-only proof shortcuts.

Recommended status:

```yaml
format: opaque_high_entropy_token
recommended_for_first_posture: true
requires_signing_dependency: false
requires_key_management: false
requires_token_storage_or_verifier: true
revocation_fit: high
redaction_fit: high
agent_clarity: high
confidence: high
```

### Signed Structured Token

Position:

```text
defer
```

Benefits:

- Can validate without storage lookup.
- Can carry issuer, audience, subject, expiration, and claims.
- Familiar to many backend teams.

Risks:

- Requires signing dependency and key-management posture.
- Revocation is harder without a server-side denylist or session store.
- Claims can drift from server truth.
- Agents may be tempted to add authorization facts into token claims instead of module-owned permission logic.
- Key rotation, algorithm agility, clock skew, audience validation, and replay posture must be explicit.

Recommendation:

```text
Defer until key management, revocation, and claim ownership are ratified.
```

### External Provider Token

Position:

```text
defer
```

Benefits:

- Useful for future platform, social, OAuth, or OIDC login.
- Can reuse provider-issued proof when provider validation is ratified.

Risks:

- Not vibit's own session token.
- Provider issuer, audience, key, metadata, outage, and refresh semantics differ.
- Requires external identity linking and provider dependency decisions.

Recommendation:

```text
Defer until external provider login and external identity linking are ratified.
```

### Plain Session ID As Secret

Position:

```text
reject_for_first
```

Benefits:

- Simple vocabulary.
- May look similar to common session-cookie systems.

Risks:

- Easy to confuse identifier with proof.
- Conflicts with current metadata-only `Session.session_id`.
- Encourages future agents to treat existing envelope fields as authenticated.
- Less clear than explicitly naming an opaque access token and a later logical session.

Recommendation:

```text
Do not use plain session ID as the first token format.
```

A future runtime session may have a high-entropy session identifier, but that must be ratified as session persistence, not as a reinterpretation of current metadata.

## 7. Carrier Details

### Login Command Response Token

Position:

```text
recommend_for_issuance
```

The first token should be issued by a future semantic login command response after credential validation succeeds.

Benefits:

- Keeps transport credential-neutral.
- Avoids Protobuf envelope changes.
- Keeps authentication result tied to explicit command contracts.
- Easy for agents to test as request/response behavior.

Risks:

- Requires future login command and response contracts.
- Requires redaction and logging rules around response payloads.
- Does not by itself bind a WebSocket connection.

### Explicit Request Proof Payload

Position:

```text
recommend_for_first_authenticated_routes
```

Authenticated commands or queries should carry token proof in explicit contract-owned payload fields until a later carrier is ratified.

Benefits:

- No envelope or handshake change required.
- Clear in semantic contracts.
- Works across transports.
- Avoids treating metadata-only fields as authority.
- Reversible if a later envelope or session-binding carrier is ratified.

Risks:

- More verbose than envelope-level or connection-bound proof.
- Requires every authenticated contract to declare proof semantics or share generated proof wrapper conventions later.
- May increase token exposure if redaction is not enforced.

### First System Message Binding

Position:

```text
defer
```

This may become a strong WebSocket-first model later, especially for realtime gameplay. It requires system-message contracts, binding state, timeout behavior, reconnect rules, and connection lifecycle tests.

### Protobuf Envelope Extension

Position:

```text
defer
```

An envelope-level proof field may eventually reduce per-contract repetition, but it changes the public wire schema and generated output. It requires a protocol change spec and compatibility review.

### Current Session Metadata

Position:

```text
reject_as_proof
```

Current `Session.session_id`, `Session.player_id`, `Session.connection_id`, and `Session.connection_epoch` fields remain metadata-only. They may be copied, normalized, or logged as metadata only where allowed. They must not authorize player-owned behavior.

### WebSocket Handshake Carrier

Position:

```text
defer
```

Handshake proof may be useful later for early rejection, but it touches the transport/process boundary and must be ratified through a dedicated handshake decision. It must not be added while comparing token formats.

## 8. Comparative Matrix

Scores are qualitative:

```text
high = favorable
medium = manageable with gates
low = unfavorable for first slice
```

| Token format | Production safety | Revocation | Redaction | Dependency load | Storage complexity | Agent clarity | First-slice fit |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `opaque_high_entropy_token` | High | High | High | High | Medium | High | High |
| `signed_structured_token` | Medium | Low | Medium | Low | Low to medium | Medium | Medium |
| `external_provider_token` | Medium | Medium | Low | Low | Medium | Low | Low |
| `plain_session_id_as_secret` | Low | Medium | Medium | High | Medium | Low | Low |

| Carrier | Boundary clarity | Protocol compatibility | WebSocket fit | Agent clarity | First-slice fit |
| --- | --- | --- | --- | --- | --- |
| `login_command_response_token` | High | High | High | High | High |
| `explicit_request_proof_payload` | High | High | Medium | High | High |
| `first_system_message_binding` | Medium | Medium | High | Medium | Medium |
| `protobuf_envelope_extension` | Medium | Low | High | Medium | Medium |
| Current `Session.session_id` metadata | Low | High | Medium | Low | Low |
| WebSocket handshake carrier | Medium | Medium | High | Medium | Medium |

## 9. Reference Mapping

### Nakama

Nakama commonly returns session tokens after authentication and supports refresh-token based continuation.

vibit adapts this as:

- Adopted concept: authentication produces a server-recognized token or session proof.
- Recommended first token format: opaque high-entropy token, not Nakama-compatible token format.
- Deferred concept: refresh token lifecycle.
- Deferred concept: realtime socket bound to authenticated session.
- Rejected for now: direct Nakama API compatibility.

### Pitaya

Pitaya's relevant input is session context and handler binding.

vibit adapts this as:

- Future handlers receive normalized request identity after validation.
- Proof carrier must produce `RequestIdentity` before domain dispatch.
- Connection binding and session object behavior remain deferred.
- Authentication must not be hidden in transport acceptors or route handlers.

## 10. Recommended W-0068 Ratification Packet

W-0068 should ratify or adjust:

```yaml
first_token_format: opaque_high_entropy_token
token_issuance_carrier: login_command_response_token
request_proof_carrier: explicit_request_proof_payload
access_token: selected
refresh_token: deferred
session_token_vocabulary: deferred_until_session_persistence
protobuf_envelope_change: false
websocket_handshake_authentication_change: false
current_session_metadata_as_proof: false
requires_before_implementation:
  - token_issuance_contract
  - token_validation_contract
  - token_hash_lookup_rule
  - token_storage_or_verifier_schema_gate
  - expiration_rule
  - revocation_logout_rule
  - redaction_rules
  - repository_checks
  - focused_tests
does_not_authorize:
  - token_parser_code
  - token_tables
  - session_tables
  - credential_tables
  - jwt_or_signing_dependency
  - protobuf_envelope_change
  - websocket_handshake_authentication
  - runtime_player_handlers
  - websocket_routes
```

## 11. Open Questions For Later Work

- Whether the first opaque token is called an access token, session token, or another vibit term in public contracts.
- Whether token storage uses PostgreSQL immediately or a later session store.
- Whether refresh tokens exist in the first production implementation.
- Whether expiration is short enough to avoid refresh in the first implementation.
- Whether revocation and logout are mandatory before the first runtime login implementation.
- Whether request proof payload fields should be repeated in every authenticated command or generated through a shared contract wrapper later.

These are not blockers for W-0067. They should be answered by W-0068 through W-0071.
