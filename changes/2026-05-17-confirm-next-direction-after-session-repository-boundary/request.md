# Request

## Original Request

The maintainer asked the agent to recommend the next ten steps and then continue those ten steps, with Nakama and Pitaya as important references.

## Clarified Requirement

Choose the next bounded milestone after the session repository boundary. The selected direction is `implement_session_repository_interface`.

## User-Visible Outcome

Maintainers and agents can see that the next work item is a storage-neutral Go session repository interface implementation under `runtime/internal/app/session`.

## Non-Goals

- Implementing a PostgreSQL session adapter.
- Adding SQL query behavior for `runtime_sessions`.
- Creating, validating, revoking, expiring, or cleaning up runtime sessions at runtime.
- Setting `RequestIdentity.SessionValidated` true.
- Changing WebSocket handshake authentication, transport credential carriers, Protobuf session messages, the existing envelope, route policy, logout active-connection behavior, reconnect/epoch behavior, dependencies, memory durable session behavior, or direct Nakama/Pitaya API compatibility.

## Unknowns

- The first PostgreSQL adapter gate remains unselected.
- Runtime session validation behavior remains unselected.
- Logout/revocation active-connection and reconnect/epoch behavior remain unselected.

## Acceptance Criteria

- [x] The selected direction is recorded as `implement_session_repository_interface`.
- [x] `M-063/W-0135` is completed.
- [x] `M-064/W-0136` is created as the next bounded implementation slice.
- [x] Ask-first boundaries remain preserved for adapters, runtime behavior, route policy, transport, protocol, dependencies, and direct Nakama/Pitaya API compatibility.
