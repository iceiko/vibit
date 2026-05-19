# Request

## Original Request

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

## Clarified Requirement

Define the conservative `define_websocket_close_policy_gate` boundary after the active connection registry implementation.

## User-Visible Outcome

The repository now has a gate-only WebSocket close policy standard:

```text
docs/websocket-close-policy-gate.md
docs/websocket-close-policy-gate.zh-CN.md
```

The gate says future close policy is application-owned, transport owns only a future narrow concrete close handoff, the active connection registry is target state rather than policy, and logout/session revocation do not close sockets until a later implementation gate explicitly selects that behavior.

## Non-Goals

- Do not implement WebSocket close behavior.
- Do not add transport close handoff code.
- Do not add close codes, close reasons, kick messages, disconnect commands, or protocol close messages.
- Do not change logout behavior or runtime session revocation.
- Do not change active connection registry behavior.
- Do not add duplicate replacement, reconnect, resume, durable epoch behavior, or protocol session carriers.
- Do not change the existing Protobuf envelope, add Protobuf logout routes, add generated output, change WebSocket handshake authentication, add transport credential carriers, or add dependencies.
- Do not add durable/distributed registry storage, cleanup jobs, operations dependencies, memory durable session behavior, broader game backend modules, or direct Nakama/Pitaya API compatibility.

## Acceptance Criteria

- [x] WebSocket close policy ownership is documented.
- [x] Transport close handoff is deferred and kept narrow.
- [x] Registry invalidation and concrete socket close remain separate.
- [x] Close intent, target, reason class, retryability, and visibility vocabulary are recorded as future choices.
- [x] Close code and reason text selection remains deferred.
- [x] Logout and runtime session revocation behavior remains unchanged.
- [x] WebSocket transport stays credential-neutral.
- [x] Existing Protobuf sources and generated output stay unchanged.
- [x] Nakama and Pitaya lessons are mapped into vibit-native design without public API compatibility claims.
- [x] Repository checks enforce the boundary.
