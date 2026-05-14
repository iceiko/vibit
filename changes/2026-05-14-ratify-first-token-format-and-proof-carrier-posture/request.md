# Request

## Original Request

The maintainer asked the agent to continue ten steps according to professional judgment, stopping only when a truly necessary decision required human input.

## Clarified Requirement

Advance `W-0068` by ratifying the first token format and proof carrier posture recommended by `W-0067`.

## User-Visible Outcome

The repository gains a durable ratification document and ADR selecting an opaque high-entropy access token, login command response issuance, and explicit request proof payload posture while preserving implementation deferral.

## Non-Goals

- Do not implement token generation, parsing, signing, validation, refresh, revocation, rotation, replay handling, or storage.
- Do not add credential tables, external identity tables, token tables, session tables, or migrations.
- Do not add password hashing, JWT, OAuth, OIDC, provider SDK, cryptography, key-management, Redis-like, or other major dependencies.
- Do not change Protobuf envelope behavior.
- Do not change WebSocket handshake authentication behavior.
- Do not add first system-message authentication.
- Do not add runtime player account handlers or WebSocket routes.
- Do not treat metadata-only `player_id`, `session_id`, `connection_id`, or `connection_epoch` as proof.

## Unknowns

- Exact access-token expiration duration.
- Refresh token inclusion.
- Revocation, rotation, replay, logout, cleanup, and audit policy.
- Token verifier storage schema.
- Future login response and authenticated request proof contract shapes.

## Acceptance Criteria

- Ratify a first access-token format.
- Ratify token issuance and request proof carrier posture.
- Record issuer, verifier, subject, audience, expiration, refresh, revocation, rotation, replay, redaction, and storage implications.
- State unchanged Protobuf envelope, current Session metadata, WebSocket handshake, and first system-message behavior.
- Preserve implementation deferral.
