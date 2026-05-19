# Request

## Original Request

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

## Clarified Requirement

Implement the selected `implement_session_creation_composition` slice by making successful device-credential login create a durable runtime session in the same application unit of work as access-token storage.

## User-Visible Outcome

`AuthenticateWithDeviceCredential` now creates an active runtime session row linked to the stored access-token record before returning the one-time client-visible access token after commit.

## Non-Goals

- Do not expose `session_id` through Protobuf login responses.
- Do not change the existing Protobuf envelope.
- Do not add WebSocket handshake authentication.
- Do not add transport credential carriers.
- Do not wire session validation into ordinary route policy.
- Do not change `ValidateAccessToken` or set `SessionValidated` true there.
- Do not implement logout, refresh, cleanup, token rotation, or active-connection invalidation.
- Do not add reconnect, resume, duplicate connection replacement, or connection epoch behavior.
- Do not add dependencies or direct Nakama/Pitaya public API compatibility.

## Acceptance Criteria

- [x] Login-created runtime sessions are created through `session.Repository.CreateRuntimeSession`.
- [x] Token storage and session creation happen in the same application unit of work.
- [x] Runtime sessions are linked to `access_token_record_id` without storing raw proof or digest material.
- [x] The first session lifetime aligns to the access-token lifetime.
- [x] Session creation failure prevents successful token/session return.
- [x] Commit failure prevents successful token/session return.
- [x] Access-token validation and route policy remain unchanged.
- [x] WebSocket transport, Protobuf envelope, generated output, logout, reconnect, and direct Nakama/Pitaya compatibility remain unchanged.
