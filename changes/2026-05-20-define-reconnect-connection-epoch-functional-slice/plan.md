# Plan

## Files To Create

- `decisions/ADR-0083-reconnect-connection-epoch-functional-slice.md`
- `conversations/2026-05-20-reconnect-connection-epoch-functional-slice.md`

## Files To Edit

- `runtime/internal/app/connection/registry.go`
- `runtime/internal/app/connection/registry_test.go`
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

- Add `StateSuperseded`.
- Add `ErrorCodeConnectionEpochStale`.
- Add `SupersededAt` and `SupersededByEpoch` registry record metadata.
- Make `RegisterOpenConnection` supersede earlier active epochs when a newer epoch arrives.
- Make `RegisterOpenConnection` reject stale or repeated epochs after a newer epoch exists.
- Preserve copy-on-read behavior for new pointer fields.

## Tests

- Registry test for newer epoch superseding earlier active epoch.
- Registry test for stale epoch rejection.
- Registry test for superseded lifecycle inspection and active-list exclusion.
- Existing close policy tests continue to prove only active bound records are targeted.

## Verification Commands

- `cd runtime && go test ./internal/app/connection`
- `cd runtime && go test ./...`
- `node -c tools/vibit`
- `node tools/vibit check change define-reconnect-connection-epoch-functional-slice --json`
- `node tools/vibit inspect next`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

No data migration is required. Rollback removes the superseded registry state, stale epoch error, and associated tests and manifest/check-rule updates.
