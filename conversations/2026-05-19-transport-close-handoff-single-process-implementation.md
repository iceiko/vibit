# Conversation: Transport Close Handoff Single Process Implementation

Date: 2026-05-19
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0081-transport-close-handoff-single-process-implementation.md`
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-transport-close-handoff-gate/`
- `changes/2026-05-19-implement-transport-close-handoff-single-process/`

Related artifacts:

- `runtime/internal/platform/transport/ws/close_handoff.go`
- `runtime/internal/platform/transport/ws/close_handoff_test.go`
- `runtime/internal/platform/transport/ws/server.go`
- `runtime/internal/platform/transport/ws/server_test.go`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/authentication/AGENTS.md`
- `modules/authentication/AGENTS.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Context

The repository was at `M-102/W-0174`, the implementation slice selected after `ADR-0080` defined the transport close handoff gate.

Completed prerequisites included:

- Application-owned active connection registry.
- Application-owned close policy with `mark_invalidated_only`.
- Client-facing protocol logout route that revokes the presented access-token record without closing sockets.
- Gate-only transport close handoff standard selecting server-observed `connection_id + connection_epoch` as the first concrete target.

The missing implementation gap was concrete WebSocket socket close mechanics that remained transport-owned and policy-neutral.

## Maintainer Narrative

The maintainer asked in Chinese:

```text
检查目录中的项目，继续推进，注意及时提交和推送。
```

The maintainer also noted that the GitHub credential existed in a non-committed local file:

```text
ssh key之前是存在了一个不提交的文件里，你自己找一下。
```

The local ignored credential file was used only for pushing commits and was not committed or recorded.

## Agent Response Summary

The agent first pushed the prior completed commit, then advanced `W-0174`.

The implementation added a narrow WebSocket transport handoff:

```text
runtime/internal/platform/transport/ws.CloseHandoffRequest
runtime/internal/platform/transport/ws.CloseHandoffResult
runtime/internal/platform/transport/ws.Server.RequestClose
```

The handoff targets only server-observed connection id and epoch. Accepted sockets are tracked in a single-process in-memory socket table inside WebSocket transport. The handoff closes the matched concrete socket with transport mechanics while preserving close code and close reason deferrals.

## Decisions

- Complete `M-102/W-0174`.
- Accept `ADR-0081`.
- Add `runtime.transport_close_handoff_single_process_implementation` as the repository check rule for the implementation state.
- Move the work queue to `M-103/W-0175`, a next-direction confirmation gate after concrete close handoff implementation.
- Recommend `define_reconnect_connection_epoch_gate` as the next lifecycle-closure direction, without selecting it in this implementation slice.

## Artifacts

- `runtime/internal/platform/transport/ws/close_handoff.go`
- `runtime/internal/platform/transport/ws/close_handoff_test.go`
- `runtime/internal/platform/transport/ws/server.go`
- `runtime/internal/platform/transport/ws/server_test.go`
- `changes/2026-05-19-implement-transport-close-handoff-single-process/`
- `decisions/ADR-0081-transport-close-handoff-single-process-implementation.md`
- `conversations/2026-05-19-transport-close-handoff-single-process-implementation.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/authentication/AGENTS.md`
- `modules/authentication/AGENTS.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Open Questions

- Reconnect, resume, duplicate replacement, and connection epoch behavior remain deferred.
- Protocol session carriers remain deferred.
- Close code mapping and close reason text remain deferred.
- Logout-triggered socket close remains deferred.
- Runtime session revocation remains deferred.
- Operations/admin disconnect remains deferred.
- Durable or distributed close handoff remains deferred.

## Follow-Up

Block at the next confirmation gate before choosing reconnect/epoch behavior, protocol session carriers, logout-triggered socket close, close code mapping, close reason text, runtime session revocation, operations/admin disconnect, direct Nakama/Pitaya API compatibility, or broader game backend planning.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, database credentials, close reason text, remote addresses, headers, or GitHub tokens are recorded in this conversation log.
