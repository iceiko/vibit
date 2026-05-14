# Request

## Original Request

The maintainer selected option `1` from the next direction gate:

```text
1
```

The option corresponds to:

```text
ratify_authentication_and_token_session_validation_design
```

## Clarified Requirement

Complete `W-0057 Define authentication and token session validation design standard`.

Define a design standard for authentication proof, token behavior, session validation, credential storage boundaries, external identity linking boundaries, session persistence boundaries, request identity trust, Protobuf envelope interaction, and WebSocket handshake interaction before implementation.

Map Nakama and Pitaya reference concepts into vibit terms as adopted, adapted, deferred, or rejected. Preserve current metadata-only validator behavior as non-authenticated. Plan the remaining M-011 work queue without adding runtime authentication code, token parsing, credential storage, session persistence, Protobuf envelope changes, WebSocket handshake changes, runtime player handlers, or WebSocket routes.

## User-Visible Outcome

The repository should have a clear standard that future agents must read before implementing authentication, token/session validation, credentials, external identities, session persistence, or handshake behavior.

`node tools/vibit inspect next --json` should advance from `W-0057` to the next bounded M-011 work item.

## Non-Goals

- Do not implement production authentication.
- Do not choose a concrete login method.
- Do not choose JWT, opaque tokens, refresh tokens, signing, expiration, revocation, rotation, or token storage behavior.
- Do not add credential storage, password hashing, cryptography, OAuth, OIDC, external identity, or session-store dependencies.
- Do not add authentication runtime code, token parsing, credential lookup, external identity linking, or session persistence.
- Do not change Protobuf envelope behavior.
- Do not change WebSocket handshake authentication behavior.
- Do not add runtime player account handlers.
- Do not add WebSocket routes.
- Do not declare metadata-only `player_id` or `session_id` sufficient for production permissions.
- Do not copy Nakama or Pitaya public API shape.

## Acceptance Criteria

- [x] Authentication/token/session validation standard exists in English.
- [x] Simplified Chinese translation exists.
- [x] Durable ADR records the boundary decision.
- [x] Nakama account/auth/session token/refresh/realtime socket concepts are mapped into vibit terms.
- [x] Pitaya session binding/handler/frontend/backend/realtime session concepts are mapped into vibit terms.
- [x] Metadata-only identity remains non-authenticated.
- [x] Remaining M-011 work queue is planned with one next-ready item.
- [x] No runtime authentication code, token parsing, credential storage, session persistence, Protobuf envelope change, WebSocket handshake change, runtime player handler, or WebSocket route is added.
- [x] Verification is recorded.
