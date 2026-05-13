# Plan

## Files To Create

- `runtime/internal/platform/protocol/protobuf/error_envelope.go`
- `runtime/internal/platform/protocol/protobuf/error_envelope_test.go`
- `changes/2026-05-13-add-application-error-envelope-mapping/`

## Files To Edit

- `runtime/internal/platform/protocol/protobuf/inventory_bridge.go`
- `runtime/internal/platform/protocol/protobuf/inventory_bridge_test.go`
- `.arch/runtime.yaml`
- `docs/runtime-protocol-adapter.md`
- `docs/runtime-protocol-adapter.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`

## Generated Artifacts

None.

Generated Protobuf files under `runtime/internal/generated/proto/` must not be edited.

## Handwritten Logic

- Add `BuildErrorEnvelopeFromApplicationResult`.
- Add `BuildErrorEnvelopeFromApplicationError`.
- Route successful application results through existing payload mapping.
- Route errored application results through the new error envelope mapping.

## Tests

- Unit tests for application error envelope mapping.
- Integration test for an inventory permission error becoming a protocol error envelope.
- Regression test that successful inventory bridge mapping still returns response envelopes.

## Verification Commands

- `cd runtime && go test ./...`
- `cd runtime && go vet ./...`
- `node tools/vibit check runtime --json`
- `node tools/vibit check change add-application-error-envelope-mapping --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan for GitHub tokens in tracked and unignored files

## Rollback Or Migration Notes

Rollback can remove the new error envelope files and revert protocol adapter documentation/status updates. No public contract, generated output, or database migration rollback is involved.
