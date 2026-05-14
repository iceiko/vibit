# Request

## Original Request

The maintainer accepted the professional recommendation to continue from the post-tooling direction gate.

## Clarified Requirement

Close `M-008 Next Direction Confirmation Gate`, record the selected direction as `ratify_player_account_postgresql_persistence_schema`, and create a bounded next milestone for player account PostgreSQL persistence schema ratification.

## User-Visible Outcome

`node tools/vibit inspect next --json` should no longer show `W-0049` as blocked. The work queue should move to `M-009 Player Account PostgreSQL Persistence Schema` with a bounded first work item.

## Non-Goals

- Do not implement production authentication.
- Do not choose login methods.
- Do not choose token format, signing, refresh, expiration, or revocation behavior.
- Do not add credential storage, password hashing, cryptography, OAuth, OIDC, or external auth dependencies.
- Do not add session persistence.
- Do not change the Protobuf envelope.
- Do not change WebSocket handshake authentication behavior.
- Do not add runtime player account handlers or WebSocket routes.
- Do not copy Nakama or Pitaya public API shapes.
- Do not declare metadata-only `player_id` or `session_id` sufficient for production permissions.

## Acceptance Criteria

- [x] `W-0049` is marked completed with the selected direction.
- [x] `M-008` is marked completed.
- [x] `M-009` is added as the active milestone.
- [x] The next milestone has a bounded `next_ready` work item.
- [x] Ask-first boundaries remain explicit.
- [x] Verification is recorded.
