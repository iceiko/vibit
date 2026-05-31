# Plan

## Files To Create

- `docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md`
- `docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.zh-CN.md`
- `decisions/ADR-0154-pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md`
- `conversations/2026-05-26-pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md`
- `changes/2026-05-26-define-pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate/`

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

Add repository check coverage to `tools/vibit` for the W-0246 gate artifacts and deferrals.

## TDD

1. Run `node tools/vibit inspect rule runtime.pitaya_aligned_distributed_runtime_vocabulary_reactivation_gate`.
2. Confirm it fails with `Unknown rule_id`.
3. Add the check rule and repository validation.
4. Run the focused and full verification commands.

## Verification Commands

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_distributed_runtime_vocabulary_reactivation_gate
node tools/vibit check change define-pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

## Rollback Or Migration Notes

Rollback is source-only: remove the gate standard, ADR/change/memory artifacts, rule registration, and manifest/doc pointer updates. No data migration is involved.
