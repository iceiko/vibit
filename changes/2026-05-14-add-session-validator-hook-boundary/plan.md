# Plan

1. Add `SessionValidator` and `SessionValidatorFunc` in `runtime/internal/app`.
2. Add `MetadataOnlySessionValidator` as the default non-authenticating validator.
3. Add `SessionValidatingDispatcher` as a wrapper around an existing `RouteDispatcher`.
4. Make the wrapper normalize request identity before calling the inner dispatcher.
5. Return an application error when an injected validator rejects a request.
6. Compose default session validation in runtime HTTP handler wiring without changing protocol or transport behavior.
7. Add focused Go tests for the validator and wrapper.
8. Update architecture manifests and runtime/player-session standards to record the hook status.
9. Run Go and repository verification.

## Files To Edit

- `runtime/internal/app/handoff.go`
- `runtime/internal/app/session_validation.go`
- `runtime/internal/app/session_validation_test.go`
- `runtime/cmd/vibit-server/main.go`
- `runtime/cmd/vibit-server/main_test.go`
- `.arch/runtime.yaml`
- `docs/player-identity-session-boundary.md`
- `docs/player-identity-session-boundary.zh-CN.md`
- `docs/runtime-protocol-adapter.md`
- `docs/runtime-protocol-adapter.zh-CN.md`
- `.arch/work-items.yaml`

## Generated Artifacts

None.

## Verification Commands

- `cd runtime && go test ./...`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change add-session-validator-hook-boundary --json`
- `node tools/vibit check all --json`
- `git diff --check`
