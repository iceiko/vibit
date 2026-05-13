# Plan

## Files To Create

- `runtime/internal/platform/protocol/protobuf/*.go`
- `runtime/internal/platform/protocol/protobuf/*_test.go`
- `runtime/internal/generated/proto/**/*.pb.go`

## Files To Edit

- `runtime/go.mod`
- `.arch/runtime.yaml`
- `.arch/README.md`
- `.arch/README.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`

## Generated Artifacts

- Generated Go Protobuf output under `runtime/internal/generated/proto/`.

## Handwritten Logic

- Define narrow runtime protocol handoff helpers for route requests and outbound messages.
- Add small tests that prove the generated Protobuf output and runtime adapter helpers compile and round-trip correctly.

## Tests

- Go unit tests for the protocol adapter helpers.
- Go unit tests for generated Protobuf message round-trips.

## Verification Commands

- `buf generate`
- `go mod tidy`
- `node --check tools/vibit`
- `go test ./...` from `runtime/`
- `node tools/vibit check generated`
- `node tools/vibit check runtime`
- `node tools/vibit check all --json`
- `git diff --check`
- secret scan for GitHub tokens outside ignored local files

## Rollback Or Migration Notes

Rollback can remove the generated Go Protobuf files, the runtime protocol adapter helpers, and the added Go dependency from `runtime/go.mod` while leaving the protocol source files and architecture decisions intact.
