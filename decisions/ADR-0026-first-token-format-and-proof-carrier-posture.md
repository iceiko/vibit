# ADR-0026: First Token Format And Proof Carrier Posture

Status: Accepted
Date: 2026-05-14
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-14-compare-token-format-and-carrier-options/`
- `changes/2026-05-14-ratify-first-token-format-and-proof-carrier-posture/`

Related conversations:

- `conversations/2026-05-14-token-format-carrier-comparison.md`
- `conversations/2026-05-14-first-token-format-proof-carrier-ratification.md`

Related artifacts:

- `docs/token-format-carrier-options.md`
- `docs/token-format-carrier-options.zh-CN.md`
- `docs/first-token-format-proof-carrier-posture.md`
- `docs/first-token-format-proof-carrier-posture.zh-CN.md`
- `docs/first-login-method-set.md`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`

## Context

M-013 exists to ratify first login methods, token format, proof carrier behavior, lifecycle semantics, schema gates, and repository checks before authentication implementation.

W-0066 ratified `device_credential_login` as the first login-method set. W-0067 compared token formats and proof carriers, recommending opaque high-entropy access tokens issued through a login command response and presented on later authenticated requests through explicit request proof payload fields.

The repository still has no production authentication implementation. Current `Session.session_id`, `Session.player_id`, `Session.connection_id`, and `Session.connection_epoch` fields remain metadata-only.

## Decision

Ratify the first token format and proof carrier posture as:

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

The selected token is a client-presented access token. It is an opaque high-entropy bearer secret with no client-readable claims.

The future issuer is an application-owned authentication boundary after successful credential and account policy validation. The future verifier is an application-owned token validator before production-sensitive domain dispatch.

The token subject is a player account identifier after selected credential and account policy success. The first audience is vibit gameplay runtime requests.

The token must have finite expiration, raw value redaction, and non-plaintext server-side verifier storage if storage is used. Exact expiration, refresh, revocation, rotation, replay, logout, cleanup, audit, and storage policy remain W-0069 work.

This decision does not implement authentication, add token storage, add credential storage, add session persistence, change the Protobuf envelope, change WebSocket handshake authentication, add runtime player handlers, or add WebSocket routes.

## Alternatives Considered

- Signed structured token.
- External provider token as vibit access token.
- Plain session ID as secret.
- Current Protobuf `Session` metadata as proof.
- Protobuf envelope extension.
- WebSocket handshake carrier.
- First system-message binding.
- Refresh token in the first posture.

## Rationale

Opaque high-entropy access tokens are the narrowest production-capable first token format for vibit.

They keep the first implementation storage-backed and explicit, avoid key-management and JWT dependency decisions, make redaction straightforward, and support revocation and logout through future server-side verifier policy. They also reduce the risk that future agents hide authorization facts in token claims.

The selected carriers preserve vibit's current runtime boundaries. Login command response issuance keeps token issuance in semantic application behavior after credential validation. Explicit request proof payloads are verbose, but they keep proof visible in future contracts without changing the Protobuf envelope or WebSocket handshake.

Nakama remains useful for session token, refresh, expiration, and logout capability vocabulary. Pitaya remains useful for session binding and handler context vocabulary. vibit adapts those concepts into application-owned validation before domain dispatch without copying either public API shape.

## Agent Reasoning Summary

The selected posture gives future agents a small and inspectable path: define lifecycle, contracts, schemas, and checks first, then implement a storage-backed opaque access token validator through the application boundary. It avoids premature JWT/key-management work and blocks the dangerous shortcut of treating existing metadata fields as proof.

## Decision Weights

```yaml
decision_weights:
  production_safety: high
  agent_context: high
  dependency_load: low
  protocol_stability: high
  storage_complexity: medium
  revocation_ergonomics: high
  human_ergonomics: medium
  reversibility: high
  long_term_maintainability: high
confidence: high
```

## Consequences

- Future first authentication implementation work must target opaque high-entropy access tokens unless a later ADR supersedes this decision.
- Future agents must not introduce JWT, signed structured tokens, external provider tokens, refresh tokens, or session token vocabulary by convenience.
- Raw token values must be redacted and must not be stored in plaintext.
- The Protobuf envelope remains unchanged.
- Existing Protobuf session metadata remains metadata-only and unauthenticated.
- WebSocket transport remains credential-neutral.
- Token lifecycle and storage implications must be defined next before implementation.

## Reversal Conditions

Revisit this decision if:

- A security review shows opaque bearer access tokens are unacceptable for the first production slice.
- A later protocol decision ratifies an envelope-level proof carrier before implementation.
- A future compatibility goal with Nakama, Pitaya, or another framework requires a different token posture.
- Operational constraints prove storage-backed token verification is too expensive without a signed-token model.
- Session persistence becomes the next ratified direction and justifies replacing access-token vocabulary with a more precise session token model.

## Follow-Up

- Define token lifecycle and storage implications.
- Define authentication contract, error, and permission surfaces.
- Define credential, token, and session schema gates.
- Add repository checks for selected login/token boundaries.
- Keep authentication implementation deferred until those gates are complete.
