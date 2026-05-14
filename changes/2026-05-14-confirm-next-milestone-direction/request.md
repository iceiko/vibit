# Request

## Original Request

The maintainer selected the next milestone direction:

```text
选1 注意参考nakama pitaya
```

This means the next direction is `ratify_player_account_and_session_contracts`, with explicit attention to Nakama and Pitaya as active references.

## Clarified Requirement

Close the `M-004 Next Direction Confirmation Gate` by recording the maintainer's selected direction and creating the next bounded milestone for player account and session contract ratification.

The new milestone must start with contract standards, not production authentication implementation. Nakama should inform the account, user, authentication, and session capability vocabulary. Pitaya should inform session binding, handler routing, and frontend/backend session vocabulary. Neither reference should become vibit's governing API shape without a separate ADR.

## User-Visible Outcome

Future `continue` / `继续` requests can advance one ready work item into player account and session contract ratification.

## Non-Goals

- Do not implement production authentication.
- Do not choose a token format.
- Do not choose credential storage, password behavior, guest login, device login, social login, or external identity provider behavior.
- Do not add player account database schema or migrations.
- Do not add session persistence.
- Do not change the Protobuf envelope.
- Do not change the WebSocket handshake.
- Do not copy Nakama or Pitaya public APIs directly.
- Do not introduce a major external framework dependency.
- Do not declare metadata-only identity sufficient for production permissions.

## Reference Intake

- Nakama authentication and session documentation inform the capability split between account authentication, user identity, session tokens, refresh, and realtime socket use.
- Pitaya overview, API, and feature documentation inform the distinction between sessions, handlers, routes, frontend servers, backend servers, and session state visibility.
- vibit keeps the Agent-Native requirement stronger than external API compatibility: every adopted capability must become explicit, contract-first, generated where appropriate, and machine-checkable.

## Acceptance Criteria

- `W-0030` is completed with a change trace.
- `M-004` is completed with a concise completion summary.
- A new active milestone records the selected direction.
- Exactly one `next_ready` work item exists for the new milestone.
- The next work item starts with account/session contract standards, not auth implementation.
- Ask-first boundaries remain for authentication scheme, token format, credentials, player account schema, session persistence, Protobuf envelope, WebSocket handshake, major dependencies, and production permission behavior.
