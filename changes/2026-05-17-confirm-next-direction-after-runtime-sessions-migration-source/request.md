# Request

## Original Request

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya
```

## Clarified Requirement

Choose the next bounded milestone after the runtime sessions migration source. The selected direction is `define_session_repository_boundary`.

The direction should use Nakama as a reference for first-class authenticated session lifecycle lookup and management, and Pitaya as a reference for keeping session context separate from acceptors, routing, and transport plumbing.

## User-Visible Outcome

Maintainers and agents can see that the next work item is a gate-only session repository boundary before Go repository code, PostgreSQL adapter behavior, or runtime session validation.

## Non-Goals

- Adding Go session repository interfaces in the confirmation change.
- Adding PostgreSQL session adapters.
- Creating, validating, revoking, or cleaning up runtime sessions.
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

- Concrete session repository API shape remains deferred to the boundary and later implementation gate.
- Whether session creation belongs to login, BindConnection, or a separate command remains deferred.
- Whether persisted session identity can satisfy route policy remains deferred.

## Acceptance Criteria

- [x] The selected direction is recorded as `define_session_repository_boundary`.
- [x] `M-061/W-0133` is closed.
- [x] `M-062/W-0134` is created as a bounded gate-only slice.
- [x] Ask-first boundaries are preserved for repository implementation, adapters, runtime behavior, route policy, logout/revocation, reconnect, operations, direct Nakama/Pitaya compatibility, and broader game backend behavior.
