# Request

## Original Request

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

## Clarified Requirement

Implement the conservative bound identity route policy slice authorized by the previous gate.

## User-Visible Outcome

`RouteProtector` now supports explicit route policy families:

- `public`
- `request_token_required`
- `bound_connection_required`
- `session_validated_required`
- `bound_session_required`

The default protected route posture remains request-token proof.

## Non-Goals

- Do not reclassify ordinary production domain routes away from request-token proof.
- Do not add WebSocket handshake authentication.
- Do not add transport credential carriers.
- Do not add Protobuf session carriers or change the existing envelope.
- Do not expose session ids to clients.
- Do not wire persistent session validation into frame handling.
- Do not add connection binding registries.
- Do not add logout/revocation active-connection invalidation.
- Do not add reconnect, resume, duplicate replacement, cleanup jobs, dependencies, memory durable session behavior, direct Nakama/Pitaya API compatibility, or broader game backend behavior.

## Acceptance Criteria

- [x] Route policy vocabulary is explicit.
- [x] Default protected route policy remains `request_token_required`.
- [x] Public authentication route remains explicit.
- [x] Metadata-only identity is rejected for protected policy families.
- [x] Bound-connection policy requires matching validated bound identity.
- [x] Session-validated policy requires already session-validated identity.
- [x] Bound-session policy requires identity source agreement.
- [x] WebSocket, Protobuf, logout, reconnect, dependencies, generated output, and direct compatibility remain unchanged.
- [x] Tests and repository checks pass.
