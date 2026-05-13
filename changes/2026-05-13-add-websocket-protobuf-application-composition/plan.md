# Plan

## Files To Create

- `runtime/internal/platform/protocol/protobuf/frame_handler.go`
- `runtime/internal/platform/protocol/protobuf/frame_handler_test.go`
- `changes/2026-05-13-add-websocket-protobuf-application-composition/*`

## Files To Edit

- `.arch/runtime.yaml`
- `.arch/work-items.yaml`
- `docs/runtime-protocol-adapter.md`
- `docs/runtime-protocol-adapter.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`

## Generated Artifacts

- None.

## Handwritten Logic

Add a small `FrameHandler` in the Protobuf protocol adapter package. It accepts copied frame payload bytes and transport metadata, decodes a Protobuf envelope, maps it to an application route request, dispatches it through an injected dispatcher interface, maps application results or application errors back into Protobuf envelopes, and returns encoded frame bytes.

## Tests

Add focused Go tests for the composition adapter.

## Verification Commands

```bash
cd runtime && go test ./...
cd runtime && go vet ./...
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change add-websocket-protobuf-application-composition --json
node tools/vibit check all --json
git diff --check
```

## Rollback Or Migration Notes

Rollback is a normal code revert. No data migration is involved.
