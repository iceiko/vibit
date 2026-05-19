# Request

## Original Request

The maintainer asked the agent to recommend the next ten steps and then continue those ten steps, with Nakama and Pitaya as important references.

## Clarified Requirement

Choose the next bounded milestone after the storage-neutral session repository interface implementation. The selected direction is `define_session_postgresql_adapter_gate`.

## User-Visible Outcome

Maintainers and agents can see that the next work item is a gate-only PostgreSQL session adapter standard. The gate prepares a future adapter implementation without adding SQL execution behavior yet.

## Non-Goals

- Implementing a PostgreSQL session adapter in this direction-confirmation change.
- Adding `runtime/internal/platform/persistence/postgres/*session*` Go files.
- Adding unit-of-work factory wiring for session repositories.
- Creating, validating, revoking, expiring, or cleaning up runtime sessions at runtime.
- Setting `RequestIdentity.SessionValidated` true.
- Changing WebSocket handshake authentication, transport credential carriers, Protobuf session messages, the existing envelope, route policy, logout active-connection behavior, reconnect/epoch behavior, dependencies, memory durable session behavior, or direct Nakama/Pitaya API compatibility.

## Unknowns

- The concrete PostgreSQL adapter implementation remains a future work item after the gate.
- Runtime session validation behavior remains unselected.
- Whether session creation belongs to login, BindConnection, or a later command remains unselected.
- Logout/revocation active-connection behavior and reconnect/epoch behavior remain unselected.

## Acceptance Criteria

- [x] The selected direction is recorded as `define_session_postgresql_adapter_gate`.
- [x] `M-065/W-0137` is completed.
- [x] `M-066/W-0138` is created as a gate-only PostgreSQL session adapter milestone.
- [x] Ask-first boundaries remain preserved for adapter implementation, runtime validation, route policy, transport, protocol, dependencies, active connection behavior, reconnect behavior, and direct Nakama/Pitaya API compatibility.
