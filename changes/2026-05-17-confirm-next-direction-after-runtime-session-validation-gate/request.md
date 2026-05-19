# Request

## Original Request

```text
你建议后面十步应该做什么？然后按照你的建议继续做十步。重点参考Nakama Pitaya，注意吸收他们设计中的优点。
```

## Clarified Requirement

After the runtime session validation gate, recommend the next bounded milestone direction and record the selection before crossing the next ask-first boundary.

## User-Visible Outcome

The work queue closes the blocked confirmation item after `W-0142`, selects `implement_runtime_session_validation`, and opens the bounded implementation milestone.

## Non-Goals

- Do not create runtime sessions at login or `BindConnection`.
- Do not change route policy to accept session-validated or bound identity.
- Do not change WebSocket handshake authentication.
- Do not add transport credential carriers.
- Do not add Protobuf session messages or envelope changes.
- Do not add logout/revocation active-connection behavior.
- Do not add reconnect or connection epoch behavior.
- Do not add direct Nakama/Pitaya API compatibility.

## Unknowns

- Session creation composition remains deferred.
- Route-policy use of session-validated identity remains deferred.
- Logout/revocation active-connection behavior remains deferred.
- Reconnect and duplicate replacement behavior remain deferred.

## Acceptance Criteria

- [x] The selected direction is recorded as `implement_runtime_session_validation`.
- [x] `M-071/W-0143` is closed.
- [x] A bounded implementation milestone/work item exists for the selected direction.
- [x] Nakama/Pitaya reference rationale is recorded.
- [x] Session creation, route policy, logout, reconnect, protocol, transport, dependencies, and direct compatibility remain out of scope.
