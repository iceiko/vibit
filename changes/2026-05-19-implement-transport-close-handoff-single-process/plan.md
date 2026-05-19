# Plan

## Files To Create

- `runtime/internal/platform/transport/ws/close_handoff.go`
- `runtime/internal/platform/transport/ws/close_handoff_test.go`
- `decisions/ADR-0081-transport-close-handoff-single-process-implementation.md`
- `conversations/2026-05-19-transport-close-handoff-single-process-implementation.md`

## Files To Edit

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
- `tools/vibit`

## Generated Artifacts

- None.

## Handwritten Logic

- Add `CloseHandoffRequest`, `CloseHandoffResult`, and redacted outcome vocabulary.
- Add a single-process accepted socket table inside WebSocket transport.
- Register accepted sockets by server-observed connection id and epoch.
- Unregister sockets when the server connection loop exits.
- Close matched active sockets through concrete transport mechanics without selecting close codes or reason text.

## Tests

- Unit tests for close requested, missing socket, stale epoch, already closed, close failed, epoch-safe unregister, and credential-neutral request/result shapes.
- Live WebSocket test proving the handoff closes an accepted socket.

## Verification Commands

- `go test ./internal/platform/transport/ws`
- `cd runtime && go test ./...`
- `node -c tools/vibit`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change implement-transport-close-handoff-single-process --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

No data migration is required. Rollback removes the transport-owned close handoff files and the server socket-table registration.
