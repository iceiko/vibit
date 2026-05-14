# Request

## Original Request

The maintainer asked the agent to continue ten steps according to professional judgment, stopping only when a truly necessary decision required human input.

## Clarified Requirement

Advance `W-0067` by comparing token format and carrier options after `device_credential_login` was ratified as the first login-method set.

## User-Visible Outcome

The repository gains a comparison document recommending an opaque high-entropy token, login-command response issuance, and explicit request proof payload posture for W-0068.

## Non-Goals

- Do not implement token parsing or validation.
- Do not add token tables, session tables, credential tables, or migrations.
- Do not add JWT, signing, cryptography, OAuth, OIDC, provider SDK, key-management, Redis-like, or session-store dependencies.
- Do not change Protobuf envelope behavior.
- Do not use current `Session.session_id` metadata as proof.
- Do not change WebSocket handshake authentication behavior.
- Do not add runtime player handlers or WebSocket routes.

## Unknowns

- Whether the first opaque token is named access token, session token, or another vibit term.
- Whether token storage uses PostgreSQL immediately.
- Whether refresh tokens are included in the first production implementation.
- Whether request proof payloads are repeated per contract or generated later.

## Acceptance Criteria

- Compare opaque high-entropy tokens, signed structured tokens, external provider tokens, and plain session-ID-as-secret posture.
- Compare login command response, explicit request payload, Protobuf envelope extension, current Session metadata, WebSocket handshake carriers, and first system-message carriers.
- Recommend a first token format and carrier posture for W-0068.
- Preserve Protobuf, WebSocket handshake, and runtime implementation boundaries.
