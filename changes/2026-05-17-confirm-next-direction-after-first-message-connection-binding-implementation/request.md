# Request

## Original Request

```text
按照你的建议继续推进。
```

## Clarified Requirement

After first-message connection binding implementation, choose the next milestone direction according to the agent's recommendation. The selected direction is `define_postgres_session_persistence_schema_gate`.

## User-Visible Outcome

The work queue records that PostgreSQL session persistence schema gating is the next selected direction. No SQL migration, repository, runtime behavior, WebSocket handshake authentication, route-policy bound identity, logout/revocation, reconnect behavior, dependency, or direct Nakama/Pitaya compatibility is added by this direction-only change.

## Non-Goals

- Implementing session persistence.
- Creating `runtime_sessions` SQL migration source.
- Adding session repository interfaces or PostgreSQL adapters.
- Making `BindConnection` create or validate durable sessions.
- Setting `RequestIdentity.SessionValidated` true.
- Adding WebSocket handshake authentication or transport credential carriers.
- Adding logout/revocation active-connection invalidation.
- Adding reconnect, resume, duplicate replacement, or durable epoch behavior.
- Adding route-policy use of bound connection identity.
- Adding direct Nakama/Pitaya API compatibility.

## Unknowns

- The exact future SQL migration is still deferred.
- Session repository shape is still deferred.
- Runtime session validation semantics are still deferred.
- Active connection invalidation and reconnect behavior are still deferred.

## Acceptance Criteria

- [x] The selected direction is recorded as `define_postgres_session_persistence_schema_gate`.
- [x] Nakama/Pitaya reference basis is recorded.
- [x] Work queue advances from `M-057/W-0129` to the gate milestone.
- [x] Deferrals remain explicit.
