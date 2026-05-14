# Request

## Original Request

The maintainer responded to the blocked next-direction gate:

```text
按照你的建议和判断推进10步，除非有非常必要的，需要我决策的，再停下来问。
```

## Clarified Requirement

Interpret the maintainer's message as authorization to choose the agent's recommended direction from `W-0063`:

```text
login_method_and_token_format_ratification
```

Close `M-012 Next Direction Confirmation Gate`, mark `W-0063` completed with the selected direction, and create a bounded next milestone for login method plus token format ratification.

The new milestone must begin with a ratification standard and reference mapping, not implementation. It must not add production authentication code, token parsing, credential storage, external identity linking, session persistence, Protobuf envelope changes, WebSocket handshake authentication, runtime player handlers, or WebSocket routes before the required artifacts are explicit.

## User-Visible Outcome

`node tools/vibit inspect next --json` should no longer show `W-0063` as blocked. The work queue should move to `M-013 Login Method And Token Format Ratification` with `W-0064` as the next ready work item.

## Non-Goals

- Do not implement production authentication.
- Do not choose concrete login methods in this direction-gate closeout.
- Do not choose concrete token format in this direction-gate closeout.
- Do not add token parsing, signing, refresh, expiration, revocation, rotation, storage, or carrier behavior.
- Do not add credential storage, password hashing, OAuth, OIDC, provider SDK, cryptography, key-management, or external identity-provider dependencies.
- Do not add credential tables, external identity tables, token tables, session tables, or migrations.
- Do not change Protobuf envelope behavior.
- Do not change WebSocket handshake authentication behavior.
- Do not add runtime player account command/query handlers.
- Do not add WebSocket routes for player account contracts.
- Do not copy Nakama or Pitaya public API shapes.
- Do not declare metadata-only `player_id` or `session_id` sufficient for production permissions.

## Acceptance Criteria

- [x] `W-0063` is marked completed with `selected_direction: login_method_and_token_format_ratification`.
- [x] `M-012` is marked completed.
- [x] `M-013` is added as the active milestone.
- [x] `W-0064` is created as the first `next_ready` work item.
- [x] Manifests record the selected direction.
- [x] A conversation log records the maintainer authorization.
- [x] Verification is recorded.
