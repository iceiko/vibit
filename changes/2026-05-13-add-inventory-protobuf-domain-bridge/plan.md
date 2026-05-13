# Plan

## Files To Create

- `runtime/internal/platform/protocol/protobuf/inventory_bridge.go`
- `runtime/internal/platform/protocol/protobuf/inventory_bridge_test.go`
- `changes/2026-05-13-add-inventory-protobuf-domain-bridge/`

## Files To Edit

- `.arch/runtime.yaml`
- `docs/runtime-protocol-adapter.md`
- `docs/runtime-protocol-adapter.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/inventory/AGENTS.md`
- `modules/inventory/AGENTS.zh-CN.md`
- `modules/inventory/module.yaml`

## Generated Artifacts

None.

Generated Protobuf files under `runtime/internal/generated/proto/` must not be edited.

## Handwritten Logic

- Add an explicit inventory protocol bridge under `runtime/internal/platform/protocol/protobuf/`.
- Keep generic envelope decoding available for tests and low-level protocol use.
- Add a bridge-aware envelope-to-route helper for application dispatch.
- Add application result and event payload mapping helpers for inventory responses and events.

## Tests

- Unit tests for request payload bridge behavior.
- Unit tests for response and event payload bridge behavior.
- Integration test that builds a Protobuf envelope, adapts it to a domain route request, dispatches through inventory handlers, and builds a Protobuf response envelope.

## Verification Commands

- `cd runtime && go test ./...`
- `cd runtime && go vet ./...`
- `node tools/vibit check runtime --json`
- `node tools/vibit check module inventory --json`
- `node tools/vibit check change add-inventory-protobuf-domain-bridge --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan for GitHub tokens in tracked and unignored files

## Rollback Or Migration Notes

Rollback can remove the bridge files and revert documentation status updates. No public contract, generated output, or database migration rollback is involved.
