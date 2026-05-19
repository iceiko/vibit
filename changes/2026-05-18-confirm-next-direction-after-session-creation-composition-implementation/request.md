# Request

## Original Request

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

## Clarified Requirement

Choose the next bounded milestone direction after session creation composition implementation, using Nakama and Pitaya as active reference baselines.

## Selected Direction

```text
define_bound_identity_route_policy_gate
```

## Rationale

Login now creates durable runtime sessions, and vibit already has request-level token validation, first-message connection binding, and runtime session validation. The next missing design boundary is how ordinary route policy should use request-token identity, bound connection identity, and session-validated identity.

## Non-Goals

- Do not implement route-policy use of bound identity in this confirmation step.
- Do not implement route-policy use of session-validated identity in this confirmation step.
- Do not change WebSocket handshake authentication.
- Do not add transport credential carriers.
- Do not expose session ids through Protobuf.
- Do not change the existing Protobuf envelope.
- Do not implement logout/revocation active-connection behavior.
- Do not implement reconnect or connection epoch behavior.
- Do not add direct Nakama/Pitaya public API compatibility.

## Acceptance Criteria

- [x] The selected direction is recorded as `define_bound_identity_route_policy_gate`.
- [x] The selection references Nakama session/gameplay-access pressure and Pitaya session/handler separation.
- [x] A new bounded milestone and work item are created.
- [x] Ask-first boundaries remain preserved for implementation, logout, reconnect, protocol carriers, operations, memory durable sessions, direct compatibility, and broader game backend modules.
