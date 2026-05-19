# Request

## Original Request

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

## Clarified Requirement

Define the gate-only boundary for future route-policy use of request-token identity, bound connection identity, and session-validated identity after session creation composition implementation.

## User-Visible Outcome

The repository now has `docs/bound-identity-route-policy-gate.md` and its Simplified Chinese translation. The gate explains how future route policy should classify routes and compose identity sources before implementation.

## Non-Goals

- Do not implement production use of bound identity for ordinary protected routes.
- Do not implement production use of session-validated identity for ordinary protected routes.
- Do not remove per-request access-token proof from ordinary protected routes.
- Do not change WebSocket handshake authentication.
- Do not add transport credential carriers.
- Do not add Protobuf session messages or generated output.
- Do not expose session ids in login responses.
- Do not change the existing Protobuf envelope.
- Do not implement logout, refresh, cleanup, token rotation, active-connection invalidation, reconnect, duplicate replacement, presence, rooms, parties, match runtime, group behavior, or broader game backend behavior.
- Do not add memory durable session behavior or direct Nakama/Pitaya public API compatibility.

## Acceptance Criteria

- [x] The gate defines future route policy ownership under `runtime/internal/app`.
- [x] The gate defines candidate policy families for public, request-token, bound-connection, session-validated, and bound-session routes.
- [x] The gate preserves current route behavior.
- [x] The gate maps Nakama and Pitaya concepts into vibit terms.
- [x] The gate preserves WebSocket, Protobuf, logout, reconnect, operations, memory durable session, and direct compatibility deferrals.
- [x] The gate is reflected in ADR, manifests, guides, check rules, and work queue.
