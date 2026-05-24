# Plan

## Files To Create

- `docs/agent-native-feature-request-test-workflow.md`
- `docs/agent-native-feature-request-test-workflow.zh-CN.md`
- `decisions/ADR-0128-agent-native-feature-request-test-workflow.md`
- `conversations/2026-05-24-agent-native-feature-request-test-workflow.md`
- `changes/2026-05-24-define-agent-native-feature-request-test-workflow/request.md`
- `changes/2026-05-24-define-agent-native-feature-request-test-workflow/spec.yaml`
- `changes/2026-05-24-define-agent-native-feature-request-test-workflow/impact.md`
- `changes/2026-05-24-define-agent-native-feature-request-test-workflow/plan.md`
- `changes/2026-05-24-define-agent-native-feature-request-test-workflow/checklist.md`
- `changes/2026-05-24-define-agent-native-feature-request-test-workflow/verification.md`

## Files To Edit

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`
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
- `tools/vibit`
- `rules/check-rules.json`

## Generated Artifacts

None.

## Handwritten Logic

Only `tools/vibit` check logic is changed. No server runtime behavior is changed.

## Verification Commands

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.agent_native_feature_request_test_workflow
node tools/vibit check change define-agent-native-feature-request-test-workflow --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

## Rollback Notes

This is a docs/manifests/checks slice. Rollback would remove the standard, ADR, change records, check rule, and W-0221 next-ready queue update. No data or runtime migration rollback is required.

