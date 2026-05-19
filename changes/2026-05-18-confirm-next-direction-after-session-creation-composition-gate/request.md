# Request

## Original Request

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

## Clarified Requirement

After the session creation composition gate, recommend the next bounded milestone direction and record the selection before crossing the next ask-first boundary.

## User-Visible Outcome

The work queue closes the blocked confirmation item after `W-0146`, selects `implement_session_creation_composition`, and opens a bounded implementation milestone.

## Non-Goals

- Do not wire session-validated identity into ordinary protected route policy.
- Do not change WebSocket handshake authentication.
- Do not add transport credential carriers.
- Do not add Protobuf session messages, login response fields, or envelope changes.
- Do not implement logout, refresh, cleanup, token rotation, or active-connection invalidation.
- Do not add reconnect, resume, duplicate replacement, or connection epoch behavior.
- Do not add direct Nakama/Pitaya public API compatibility.
- Do not expand into presence, rooms, parties, match runtime, groups, or broader game backend behavior.

## Unknowns

- Route-policy use of session-validated identity remains deferred.
- Logout/revocation active-connection behavior remains deferred.
- Reconnect/resume/duplicate replacement behavior remains deferred.
- Client-visible session carrier and protocol response shape remain deferred.

## Acceptance Criteria

- [x] The selected direction is recorded as `implement_session_creation_composition`.
- [x] `M-075/W-0147` is closed.
- [x] A bounded implementation milestone/work item exists for the selected direction.
- [x] Nakama/Pitaya reference rationale is recorded.
- [x] Route policy, WebSocket handshake, Protobuf carrier changes, logout, reconnect, direct compatibility, and broader game backend scope remain out of scope.
