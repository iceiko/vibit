# Request

## Original Request

```text
按照你的建议继续推进。
```

## Clarified Requirement

Define a gate-only PostgreSQL session persistence schema standard after first-message connection binding implementation. The gate must plan the future first durable runtime session table and explicitly defer SQL migration source, repositories, adapters, runtime behavior, route policy, logout/revocation, reconnect behavior, transport credential carriers, dependencies, and direct Nakama/Pitaya compatibility.

## User-Visible Outcome

Maintainers and agents can inspect the future session schema boundary before any session persistence implementation is added.

## Non-Goals

- Adding `runtime_sessions` migration source.
- Adding `runtime_session_connections` or connection registry storage.
- Adding repository interfaces or PostgreSQL adapters.
- Creating, validating, revoking, or cleaning up sessions at runtime.
- Setting `RequestIdentity.SessionValidated` true.
- Changing WebSocket handshake authentication or parsing transport credential carriers.
- Changing the existing Protobuf envelope.
- Adding Protobuf messages or generated output.
- Making bound connection identity satisfy ordinary protected route policy.
- Adding logout/revocation active-connection invalidation.
- Adding reconnect, resume, duplicate replacement, or durable epoch behavior.
- Adding dependencies.
- Adding direct Nakama/Pitaya API compatibility.

## Unknowns

- Exact SQL constraints and indexes remain deferred to the migration-source work item.
- Session repository API shape remains deferred.
- Whether session creation is tied to login, BindConnection, or a separate command remains deferred.
- Whether bound identity can satisfy route policy remains deferred.

## Acceptance Criteria

- [x] English and Simplified Chinese gate standards exist.
- [x] ADR-0059 records the decision.
- [x] The future `runtime_sessions` logical table candidate is described.
- [x] The future migration source candidate is recorded without creating it.
- [x] Runtime, contract, convention, reference, module, AGENTS, work, and rule artifacts reference the gate.
- [x] Repository check rule exists.
- [x] Deferrals are explicit.
