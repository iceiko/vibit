# Request

## Original Request

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

## Clarified Requirement

Choose the next bounded milestone direction after the bound identity route policy gate.

## User-Visible Outcome

The selected direction is `implement_bound_identity_route_policy`.

## Non-Goals

- Do not select logout/revocation active-connection behavior in this step.
- Do not select reconnect/epoch behavior in this step.
- Do not select protocol session carriers in this step.
- Do not select operations hardening, memory durable session behavior, direct Nakama/Pitaya API compatibility, or broader game backend modules in this step.

## Acceptance Criteria

- [x] The selected direction is recorded.
- [x] The next milestone and work item are created.
- [x] Ask-first boundaries remain explicit.
- [x] Nakama and Pitaya reference lessons are considered without copying public APIs.
