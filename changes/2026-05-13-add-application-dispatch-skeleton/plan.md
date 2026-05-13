# Plan

## Files To Create

- `runtime/internal/app/dispatch.go`
- `runtime/internal/app/dispatch_test.go`
- `conversations/2026-05-13-application-dispatch-skeleton.md`

## Files To Edit

- `.arch/runtime.yaml`
- `README.md`
- `README.zh-CN.md`
- `docs/runtime-protocol-adapter.md`
- `docs/runtime-protocol-adapter.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Generated Artifacts

None.

## Handwritten Logic

- Add a pure application dispatcher with explicit route registration.
- Add stable application dispatch errors.
- Add a handler interface and handler function adapter.
- Preserve request metadata in application results.
- Add runtime layer import boundary checks for app/domain package direction.

## Tests

- Go unit tests for dispatcher behavior.
- Runtime check verification for layer-boundary rules through `node tools/vibit check runtime`.

## Verification Commands

- `node --check tools/vibit`
- `go test ./...` from `runtime/`
- `node tools/vibit check runtime --json`
- `node tools/vibit check change add-application-dispatch-skeleton --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan for GitHub tokens outside ignored local files

## Rollback Or Migration Notes

Rollback can remove the dispatcher files, remove the new runtime layer-boundary check, and restore documentation to the previous runtime handoff slice state. No database or public protocol migration is involved.
