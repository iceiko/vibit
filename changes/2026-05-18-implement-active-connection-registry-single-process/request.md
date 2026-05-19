# Request

## Original Request

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

## Clarified Requirement

Implement the conservative `implement_active_connection_registry_single_process` slice authorized after the active connection registry gate.

## User-Visible Outcome

The runtime now has an application-owned in-memory active connection registry under:

```text
runtime/internal/app/connection
```

The registry can:

- Register server-observed open connections by connection id and epoch.
- Bind validated player identity linkage, with optional runtime session id and access-token record id.
- Mark active records closed or invalidated.
- Find lifecycle records by connection id and epoch.
- List active bound records by player id, runtime session id, or access-token record id.

## Non-Goals

- Do not close WebSocket connections.
- Do not add kick or disconnect behavior.
- Do not revoke runtime sessions.
- Do not add duplicate replacement behavior.
- Do not add reconnect, resume, connection-epoch protocol behavior, or durable epoch behavior.
- Do not add Protobuf logout routes or protocol session carriers.
- Do not change the existing envelope.
- Do not add WebSocket handshake authentication or transport credential carriers.
- Do not add durable/distributed registry storage, cleanup jobs, dependencies, memory durable session behavior, broader game backend modules, or direct Nakama/Pitaya API compatibility.

## Acceptance Criteria

- [x] Registry is application-owned under `runtime/internal/app/connection`.
- [x] Registry is single-process, in-memory, and non-durable.
- [x] Registration requires server-observed connection id and epoch.
- [x] Binding requires validated player identity linkage.
- [x] Metadata-only targeting is rejected.
- [x] Closed and invalidated records are excluded from active target lists.
- [x] Returned records and slices are copies.
- [x] Raw token or credential material, lookup digests, verifier digests, verifier key ids, transport headers, cookies, query strings, subprotocol values, remote addresses, and inner payload bytes are not stored.
- [x] WebSocket, Protobuf, PostgreSQL, generated output, dependencies, and direct compatibility remain unchanged.
- [x] Tests and repository checks pass.
