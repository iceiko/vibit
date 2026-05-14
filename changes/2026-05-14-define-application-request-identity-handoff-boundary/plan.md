# Plan

1. Add application-owned request identity vocabulary in `runtime/internal/app`.
2. Attach `RequestIdentity` to `RouteRequest` and `ApplicationResult`.
3. Add helpers for metadata-only session identity and future validated player identity handoff.
4. Make Protobuf envelope decoding convert existing session metadata into metadata-only application identity.
5. Make dispatch populate metadata-only identity for callers that still provide only `Session`.
6. Preserve identity in application, bootstrap, and inventory application results.
7. Keep frame connection metadata synchronized with metadata-only identity.
8. Add focused Go tests for metadata normalization, dispatcher propagation, envelope conversion, and frame metadata refresh.
9. Mark `W-0024` completed and create the next conservative work item if verification passes.

## Files To Edit

- `runtime/internal/app/handoff.go`
- `runtime/internal/app/dispatch.go`
- `runtime/internal/app/handoff_test.go`
- `runtime/internal/app/dispatch_test.go`
- `runtime/internal/app/bootstrap/inventory.go`
- `runtime/internal/modules/inventory/inventory.go`
- `runtime/internal/platform/protocol/protobuf/envelope.go`
- `runtime/internal/platform/protocol/protobuf/envelope_test.go`
- `runtime/internal/platform/protocol/protobuf/frame_handler.go`
- `runtime/internal/platform/protocol/protobuf/frame_handler_test.go`
- `docs/runtime-protocol-adapter.md`
- `docs/runtime-protocol-adapter.zh-CN.md`
- `docs/player-identity-session-boundary.md`
- `docs/player-identity-session-boundary.zh-CN.md`
- `.arch/runtime.yaml`
- `.arch/protocol.yaml`
- `.arch/conventions.yaml`
- `.arch/work-items.yaml`

## Generated Artifacts

None.

## Verification Commands

- `cd runtime && go test ./...`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change define-application-request-identity-handoff-boundary --json`
- `node tools/vibit check all --json`
- `git diff --check`
