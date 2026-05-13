# Plan

## Files To Create

- `runtime/internal/platform/transport/ws/server.go`
- `runtime/internal/platform/transport/ws/server_test.go`
- `changes/2026-05-13-add-first-websocket-transport-adapter/`

## Files To Edit

- `runtime/go.mod`
- `runtime/go.sum`
- `.arch/runtime.yaml`
- `docs/runtime-protocol-adapter.md`
- `docs/runtime-protocol-adapter.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`

## Generated Artifacts

None.

## Handwritten Logic

- Define `Frame`, `FrameHandler`, and `Server` types in the WebSocket transport package.
- Accept WebSocket connections through `github.com/coder/websocket`.
- Read binary frames, copy payload bytes, and pass them to the injected handler.
- Write handler responses as binary frames.
- Reject unsupported frame types.

## Tests

- Use `httptest.Server` and `github.com/coder/websocket` clients.
- Test binary frame echo through the handler.
- Test text frame rejection.
- Test connection metadata passed to the handler.

## Verification Commands

- `cd runtime && go mod tidy`
- `cd runtime && go test ./...`
- `cd runtime && go vet ./...`
- `node tools/vibit check runtime --json`
- `node tools/vibit check change add-first-websocket-transport-adapter --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan for GitHub tokens in tracked and unignored files

## Rollback Or Migration Notes

Rollback can remove the WebSocket transport files and dependency update. No public contract, generated output, or database migration rollback is involved.
