# Plan

## Files To Create

- `docs/session-postgresql-adapter-gate.md`
- `docs/session-postgresql-adapter-gate.zh-CN.md`
- `decisions/ADR-0063-session-postgresql-adapter-gate.md`
- `conversations/2026-05-17-session-postgresql-adapter-gate.md`
- `changes/2026-05-17-define-session-postgresql-adapter-gate/*`

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

No runtime logic is added. The change is a gate-only standard plus repository checks and manifests.

## Tests

No Go tests are added because no Go adapter code is added. Verification uses repository checks.

## Verification Commands

- `node -c tools/vibit`
- `go test ./...`
- `node tools/vibit check migrations --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check module authentication --json`
- `node tools/vibit check work --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check change confirm-next-direction-after-session-repository-interface --json`
- `node tools/vibit check change define-session-postgresql-adapter-gate --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

Reversal means deleting `docs/session-postgresql-adapter-gate.md`, the paired translation, ADR-0063, W-0138 manifest/check entries, and reopening `M-065/W-0137`. No data migration rollback is required.
