# Plan

## Files To Create

- `runtime/internal/app/session/repository.go`
- `runtime/internal/app/session/repository_test.go`
- `decisions/ADR-0062-session-repository-interface-implementation.md`
- `conversations/2026-05-17-session-repository-interface-implementation.md`
- `changes/2026-05-17-implement-session-repository-interface/*`

## Files To Edit

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/authentication/module.yaml`
- `modules/authentication/AGENTS.md`
- `modules/authentication/AGENTS.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Generated Artifacts

- None.

## Handwritten Logic

Add storage-neutral Go value types, repository interface methods, and normalization helpers. The code must not import PostgreSQL, WebSocket, generated Protobuf, or platform adapters.

## Tests

Add focused package tests under `runtime/internal/app/session`.

## Verification Commands

- `node -c tools/vibit`
- `go test ./...`
- `node tools/vibit check migrations --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check module authentication --json`
- `node tools/vibit check work --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check change confirm-next-direction-after-session-repository-boundary --json`
- `node tools/vibit check change implement-session-repository-interface --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

Reversal means deleting `runtime/internal/app/session`, removing ADR-0062 and the W-0136 manifest/check entries, and reopening `M-063/W-0135`. No data migration rollback is required.
