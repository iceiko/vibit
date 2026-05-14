# Request

## Original Request

The maintainer asked the agent to continue ten steps according to professional judgment, stopping only when a truly necessary decision required human input.

## Clarified Requirement

Advance `W-0066` by ratifying the first login-method set recommended by `W-0065`.

## User-Visible Outcome

The repository gains a durable ratification document and ADR selecting `device_credential_login` as the first login-method set, with rejected alternatives, decision weights, known gaps, and implementation gates.

## Non-Goals

- Do not implement any login method.
- Do not add credential tables, external identity tables, token tables, session tables, or migrations.
- Do not add password hashing, JWT, OAuth, OIDC, provider SDK, cryptography, key-management, Redis-like, or other major dependencies.
- Do not change Protobuf envelope behavior.
- Do not change WebSocket handshake authentication behavior.
- Do not add runtime player account handlers or WebSocket routes.
- Do not treat metadata-only `player_id`, `session_id`, or `connection_id` as proof.

## Unknowns

- Whether installation credentials are client-generated, server-issued, or both.
- Whether account creation is implicit or explicit.
- Which token format and carrier will be selected.
- Whether refresh tokens or persisted runtime sessions are in the first implementation.

## Acceptance Criteria

- Ratify a first login-method set.
- Record public rationale, rejected alternatives, decision weights, confidence, and known gaps.
- State production, bootstrap, and local-development classification.
- State account creation, account linking, and existing-account authentication posture.
- Preserve implementation deferral.
