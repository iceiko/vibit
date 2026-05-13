# Plan

## Files To Create

- `runtime/internal/app/bootstrap/`
- `runtime/internal/platform/process/`
- `runtime/cmd/vibit-server/main.go`
- `changes/2026-05-13-add-minimal-websocket-process-wiring/*`

## Files To Edit

- `.arch/runtime.yaml`
- `.arch/protocol.yaml`
- `.arch/work-items.yaml`
- `docs/runtime-protocol-adapter.md`
- `docs/runtime-protocol-adapter.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`

## Generated Artifacts

- None.

## Handwritten Logic

- Add in-memory inventory repository and bootstrap policy helpers.
- Add process assembly that registers inventory routes, composes Protobuf frame handling with the WebSocket transport, and mounts `/v1/ws`.
- Add `cmd/vibit-server/main.go` as a thin process entrypoint.

## Tests

- Add an integration test for `/v1/ws` using a test server and a binary Protobuf request.

## Verification Commands

```bash
cd runtime && go test ./...
cd runtime && go vet ./...
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change add-minimal-websocket-process-wiring --json
node tools/vibit check all --json
git diff --check
```

## Rollback Or Migration Notes

Rollback is a normal code revert. No persisted data is involved.
