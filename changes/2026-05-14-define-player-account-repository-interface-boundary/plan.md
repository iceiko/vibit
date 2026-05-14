# Plan

## Files To Create

- `runtime/internal/modules/player/repository.go`
- `runtime/internal/modules/player/repository_test.go`
- `changes/2026-05-14-define-player-account-repository-interface-boundary/`

## Files To Edit

- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `.arch/work-items.yaml`
- `modules/player/module.yaml`
- `modules/player/AGENTS.md`
- `modules/player/AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `docs/postgresql-persistence-boundary.md`
- `docs/postgresql-persistence-boundary.zh-CN.md`
- `tools/vibit`

## Generated Artifacts

- None.

## Handwritten Logic

Add storage-neutral player account repository interface definitions and small validation helpers. Keep the package standard-library only.

## Tests

Add focused tests under `runtime/internal/modules/player/`.

## Verification Commands

- `go test ./internal/modules/player`
- `node tools/vibit check runtime --json`
- `node tools/vibit check module player --json`
- `node tools/vibit check contracts --json`
- `node tools/vibit check migrations --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change define-player-account-repository-interface-boundary --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

No database rollback is needed. This change adds an interface boundary only.
