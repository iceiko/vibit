# Plan

## Files To Create

- `runtime/internal/platform/protocol/protobuf/request_loop_fixture_test.go`
- `changes/2026-05-13-add-request-loop-test-fixture/*`

## Files To Edit

- `runtime/internal/platform/protocol/protobuf/frame_handler_test.go`
- `runtime/internal/platform/protocol/protobuf/inventory_bridge_test.go`
- `.arch/runtime.yaml`
- `.arch/work-items.yaml`
- `docs/runtime-protocol-adapter.md`
- `docs/runtime-protocol-adapter.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`

## Generated Artifacts

- None.

## Handwritten Logic

- Add a package-local test fixture that builds an in-memory inventory dispatcher with existing inventory runtime helpers.
- Add shared helpers for Protobuf envelope construction, frame payload marshaling, response envelope unmarshaling, and route request conversion.
- Refactor existing Protobuf protocol adapter tests to use the fixture instead of duplicated local fake implementations.

## Tests

- Update existing frame handler and inventory bridge tests to use the fixture.
- Preserve existing request-loop behavior assertions.

## Verification Commands

```bash
cd runtime && go test ./...
cd runtime && go vet ./...
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change add-request-loop-test-fixture --json
node tools/vibit check all --json
git diff --check
```

## Rollback Or Migration Notes

Rollback is a normal code revert. No generated output or persisted data is involved.
