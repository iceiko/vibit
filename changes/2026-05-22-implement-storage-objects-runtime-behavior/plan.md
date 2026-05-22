# Plan

## Files To Create

- `runtime/internal/app/storage/service.go`
- `runtime/internal/app/storage/service_test.go`
- `decisions/ADR-0117-storage-objects-runtime-behavior-implementation.md`
- `conversations/2026-05-22-storage-objects-runtime-behavior-implementation.md`
- `changes/2026-05-22-implement-storage-objects-runtime-behavior/request.md`
- `changes/2026-05-22-implement-storage-objects-runtime-behavior/impact.md`
- `changes/2026-05-22-implement-storage-objects-runtime-behavior/plan.md`
- `changes/2026-05-22-implement-storage-objects-runtime-behavior/checklist.md`
- `changes/2026-05-22-implement-storage-objects-runtime-behavior/spec.yaml`
- `changes/2026-05-22-implement-storage-objects-runtime-behavior/verification.md`

## Files To Edit

- `tools/vibit`
- `rules/check-rules.json`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `docs/v0.1-alpha-goal.md`
- `docs/v0.1-alpha-goal.zh-CN.md`
- `docs/alpha-developer-flow.md`
- `docs/alpha-developer-flow.zh-CN.md`
- `docs/alpha-acceptance-checklist.md`
- `docs/alpha-acceptance-checklist.zh-CN.md`
- `docs/product-maturity-milestones.md`
- `docs/product-maturity-milestones.zh-CN.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `docs/nakama-pitaya-product-parity-roadmap.zh-CN.md`

## Generated Artifacts

None.

## Handwritten Logic

- Add an application service for storage object get/list/put/delete.
- Validate request identity before repository access.
- Validate collection, key, JSON object values, value size, expected versions, and list pagination.
- Use unit-of-work storage repository capability.
- Map storage repository conflicts to stable redacted public error codes.

## Tests

- Add focused fake-repository service tests under `runtime/internal/app/storage`.
- Run package tests and full runtime tests.

## Verification Commands

- `cd runtime && go test ./internal/app/storage`
- `cd runtime && go test ./...`
- `node -c tools/vibit`
- `node tools/vibit inspect next --json`
- `node tools/vibit inspect rule runtime.storage_objects_runtime_behavior_implementation`
- `node tools/vibit check change implement-storage-objects-runtime-behavior --json`
- `node tools/vibit check module storage --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

No migration rollback is needed because this slice changes no schema and adds no generated output or protocol route.
