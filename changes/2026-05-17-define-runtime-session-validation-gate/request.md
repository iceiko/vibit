# Request

## Original Request

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

## Clarified Requirement

Define a gate-only standard for future runtime session validation after the PostgreSQL session adapter implementation, while preserving strict deferrals for runtime behavior.

## User-Visible Outcome

Maintainers and agents can see exactly how future persisted runtime session validation should be implemented and what must remain out of scope until separately authorized.

## Non-Goals

- Do not implement runtime session validation.
- Do not set `RequestIdentity.SessionValidated` true.
- Do not create sessions at login or `BindConnection`.
- Do not change route policy.
- Do not add WebSocket handshake authentication or transport credential parsing.
- Do not add Protobuf session messages, envelope changes, or generated output.
- Do not add logout/revocation active-connection behavior.
- Do not add reconnect/epoch behavior.
- Do not add dependencies or direct Nakama/Pitaya API compatibility.

## Unknowns

- The exact future validator source shape remains deferred.
- Whether validation updates `last_seen_at` remains deferred.
- Whether session-validated identity can satisfy protected routes remains deferred.

## Acceptance Criteria

- [x] English and Simplified Chinese gate documents exist.
- [x] ADR records the decision.
- [x] Work queue records `M-070/W-0142`.
- [x] Repository check rule exists.
- [x] Runtime behavior remains unchanged.
