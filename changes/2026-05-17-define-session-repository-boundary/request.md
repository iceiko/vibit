# Request

## Original Request

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya
```

## Clarified Requirement

Define a gate-only session repository boundary after the runtime sessions migration source. The boundary must describe the future storage-neutral repository owner, capability vocabulary, data boundary, authentication/session separation, WebSocket/protocol deferrals, and Nakama/Pitaya reference mapping.

## User-Visible Outcome

Maintainers and agents can read a clear standard before adding Go session repository interfaces, PostgreSQL session adapters, or runtime session validation.

## Non-Goals

- Adding `runtime/internal/app/session` or Go repository interfaces.
- Adding PostgreSQL session adapter behavior.
- Changing migrations or adding session tables.
- Creating, validating, revoking, or cleaning up runtime sessions.
- Setting `RequestIdentity.SessionValidated` true.
- Changing route policy.
- Changing WebSocket handshake authentication or transport credential carriers.
- Adding Protobuf session messages or changing the existing Protobuf envelope.
- Adding logout/revocation active-connection invalidation.
- Adding reconnect, resume, duplicate replacement, or durable epoch behavior.
- Adding dependencies.
- Adding direct Nakama/Pitaya API compatibility.

## Acceptance Criteria

- [x] English and Simplified Chinese standards define the session repository boundary.
- [x] `ADR-0061` records the decision.
- [x] The future repository owner candidate is `runtime/internal/app/session`.
- [x] The future PostgreSQL adapter owner remains `runtime/internal/platform/persistence/postgres`.
- [x] Candidate repository capabilities are listed without implementing them.
- [x] Token validation, request identity, WebSocket transport, Protobuf, route policy, logout/revocation, reconnect, dependencies, memory durable session behavior, and direct Nakama/Pitaya compatibility remain deferred.
- [x] Repository checks cover the boundary.
