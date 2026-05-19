# Request

## Original Request

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya
```

## Clarified Requirement

Choose the next bounded milestone after the PostgreSQL session persistence schema gate. The selected direction is `implement_runtime_sessions_migration_source`.

The direction should use Nakama as a reference for first-class authenticated session lifecycle records and Pitaya as a reference for keeping session context separate from transport acceptors and handlers.

## User-Visible Outcome

Maintainers and agents can see that the next work item is a migration-source-only implementation slice for `runtime_sessions`.

## Non-Goals

- Adding session repository interfaces or PostgreSQL adapters in the confirmation change.
- Creating runtime session validation behavior.
- Setting `RequestIdentity.SessionValidated` true.
- Changing route policy.
- Changing WebSocket handshake authentication or parsing transport credential carriers.
- Changing the existing Protobuf envelope.
- Creating durable connection registry storage.
- Adding logout/revocation active-connection invalidation.
- Adding reconnect, resume, duplicate replacement, or durable epoch behavior.
- Adding dependencies.
- Adding direct Nakama/Pitaya API compatibility.

## Unknowns

- Session repository API shape remains deferred.
- Whether session creation belongs to login, BindConnection, or a separate command remains deferred.
- Whether persisted session identity can satisfy route policy remains deferred.

## Acceptance Criteria

- [x] The selected direction is recorded as `implement_runtime_sessions_migration_source`.
- [x] `M-059/W-0131` is closed.
- [x] `M-060/W-0132` is created as a bounded migration-source-only slice.
- [x] Ask-first boundaries are preserved for repositories, adapters, runtime behavior, route policy, logout/revocation, reconnect, operations, direct Nakama/Pitaya compatibility, and broader game backend behavior.
