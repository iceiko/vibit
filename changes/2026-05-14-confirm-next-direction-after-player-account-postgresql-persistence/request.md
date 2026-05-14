# Request

## Original Request

The maintainer selected option 1 from the next-direction confirmation gate:

```text
1
```

The option corresponds to:

```text
ratify_authentication_and_token_session_validation_design
```

## Clarified Requirement

Close `M-010 Next Direction Confirmation Gate`, record the selected direction as `ratify_authentication_and_token_session_validation_design`, and create a bounded next milestone for authentication plus token/session validation design.

The new milestone must start with design standards and reference mapping, not implementation. It must not choose concrete login methods, token format, signing behavior, refresh behavior, credential storage, external identity linking, session persistence, Protobuf envelope changes, or WebSocket handshake authentication behavior before those decisions are explicitly ratified inside the new milestone.

## User-Visible Outcome

`node tools/vibit inspect next --json` should no longer show `W-0056` as blocked. The work queue should move to `M-011 Authentication And Token Session Validation Design` with a bounded first work item.

## Non-Goals

- Do not implement production authentication.
- Do not choose login methods.
- Do not choose JWT, opaque token, or another concrete token format.
- Do not choose token signing, refresh, expiration, revocation, or storage behavior.
- Do not add credential storage, password hashing, OAuth, OIDC, or external identity-provider dependencies.
- Do not add external identity linking tables or behavior.
- Do not add session persistence.
- Do not change Protobuf envelope behavior.
- Do not change WebSocket handshake authentication behavior.
- Do not add runtime player account command/query handlers.
- Do not add WebSocket routes for player account contracts.
- Do not copy Nakama or Pitaya public API shapes.
- Do not declare metadata-only `player_id` or `session_id` sufficient for production permissions.

## Acceptance Criteria

- [x] `W-0056` is marked completed with the selected direction.
- [x] `M-010` is marked completed.
- [x] `M-011` is added as the active milestone.
- [x] The new milestone has exactly one `next_ready` first work item.
- [x] The first work item is a design standard step, not implementation.
- [x] Ask-first boundaries remain explicit for token format, login methods, credential storage, external identity linking, session persistence, Protobuf envelope behavior, WebSocket handshake behavior, runtime player handlers, and WebSocket routes.
- [x] Verification is recorded.
