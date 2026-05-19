# Request

## Original Request

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya
```

## Clarified Requirement

Implement the selected migration-source-only slice by adding the first PostgreSQL `runtime_sessions` migration source. The migration must follow the prior schema gate and remain limited to durable session lifecycle storage.

## User-Visible Outcome

Maintainers and agents can inspect and apply a concrete SQL migration source for the first durable runtime session lifecycle table.

## Non-Goals

- Adding session repository interfaces.
- Adding PostgreSQL session adapters.
- Creating, validating, revoking, or cleaning up sessions at runtime.
- Setting `RequestIdentity.SessionValidated` true.
- Creating sessions in login or BindConnection.
- Changing WebSocket handshake authentication or parsing transport credential carriers.
- Changing the existing Protobuf envelope.
- Adding Protobuf session messages or generated output.
- Creating `runtime_session_connections` or connection registry storage.
- Making bound connection identity or persisted session identity satisfy ordinary protected route policy.
- Adding logout/revocation active-connection invalidation.
- Adding reconnect, resume, duplicate replacement, or durable epoch behavior.
- Adding dependencies.
- Adding memory durable session behavior.
- Adding direct Nakama/Pitaya API compatibility.

## Unknowns

- Session repository API shape remains deferred.
- Runtime session creation trigger remains deferred.
- Session validation and route-policy semantics remain deferred.
- Logout/revocation active-connection behavior remains deferred.
- Reconnect and connection epoch behavior remain deferred.

## Acceptance Criteria

- [x] `runtime/migrations/postgres/000005_create_runtime_sessions.sql` exists.
- [x] The migration declares goose Up and Down markers.
- [x] The migration declares `-- Module: runtime.session`.
- [x] The migration creates `runtime_sessions`.
- [x] Required lifecycle columns are present.
- [x] Optional revocation and access-token record linkage columns are present.
- [x] Raw token, raw credential, token digest, and credential digest columns are absent.
- [x] No `runtime_session_connections` table is created.
- [x] No Go repository, adapter, runtime behavior, Protobuf, generated output, or dependency change is added by this slice.
