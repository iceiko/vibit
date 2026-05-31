# Plan

## Files To Create

- `decisions/ADR-0155-pitaya-aligned-distributed-runtime-vocabulary-source-first-map.md`
- `conversations/2026-05-26-pitaya-aligned-distributed-runtime-vocabulary-source-first-map.md`
- `changes/2026-05-26-implement-pitaya-aligned-distributed-runtime-vocabulary-source-first-map/`

## Files To Edit

- `tools/vibit`
- `rules/check-rules.json`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `modules/friends/module.yaml`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`
- `modules/friends/AGENTS.md`
- `modules/friends/AGENTS.zh-CN.md`
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

Add a `tools/vibit inspect pitaya-vocabulary` repository inspection command and runtime check coverage for the W-0247 artifacts, output markers, redaction posture, and deferrals.

## TDD

1. Run `node tools/vibit inspect rule runtime.pitaya_aligned_distributed_runtime_vocabulary_source_first_map`.
2. Confirm it fails with `Unknown rule_id`.
3. Run `node tools/vibit inspect pitaya-vocabulary --json`.
4. Confirm it fails with `Unknown command`.
5. Add the inspection command, rule, change artifacts, and repository validation.
6. Run the focused and full verification commands.

## Verification Commands

```bash
node -c tools/vibit
node tools/vibit inspect pitaya-vocabulary --json
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_distributed_runtime_vocabulary_source_first_map
node tools/vibit check change implement-pitaya-aligned-distributed-runtime-vocabulary-source-first-map --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

## Rollback Or Migration Notes

Rollback is source-only: remove the inspection command, ADR/change/memory artifacts, rule registration, and manifest/doc pointer updates. No data migration is involved.
