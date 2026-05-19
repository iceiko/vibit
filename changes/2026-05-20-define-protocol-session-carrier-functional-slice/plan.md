# Plan

## Files To Create

- `decisions/ADR-0084-protocol-session-carrier-functional-slice.md`
- `conversations/2026-05-20-protocol-session-carrier-functional-slice.md`

## Files To Edit

- `runtime/internal/platform/protocol/protobuf/inventory_bridge.go`
- `runtime/internal/platform/protocol/protobuf/inventory_bridge_test.go`
- `runtime/internal/platform/protocol/protobuf/authentication_bridge.go`
- `runtime/internal/platform/protocol/protobuf/authentication_bridge_test.go`
- `runtime/internal/platform/protocol/protobuf/authenticated_request_test.go`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/authentication/AGENTS.md`
- `modules/authentication/AGENTS.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Generated Artifacts

- None.

## Handwritten Logic

- Add `SessionCarrierFromApplicationResult`.
- Add login-result-specific session carrier extraction for `AuthenticationResult`.
- Keep metadata-only identity unchanged.
- Let validated identity supply response session metadata when application results already carry validated identity.

## Tests

- Test login response envelope carries the login-created runtime session id and player id.
- Test validated identity can provide response session metadata.
- Test metadata-only identity does not upgrade response session metadata.

## Verification Commands

- `cd runtime && go test ./internal/platform/protocol/protobuf`
- `cd runtime && go test ./...`
- `node -c tools/vibit`
- `node tools/vibit check change define-protocol-session-carrier-functional-slice --json`
- `node tools/vibit inspect next`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

No data migration is required. Rollback removes the Protobuf response session-carrier helper, associated tests, and manifest/check-rule updates.
