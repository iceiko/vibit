# Plan

1. Capture the selected Nakama-style user requirement in `request.md`.
2. Write `spec.yaml` with workflow phases, Nakama mapping, acceptance criteria, test plan, implementation boundaries, verification, and memory updates.
3. Add `ADR-0129` to record the pilot decision.
4. Add a conversation log preserving maintainer intent and the selected presence/status pilot.
5. Update `.arch/work-items.yaml`:
   - mark `M-149/W-0221` completed;
   - open `M-150/W-0222` as next-ready.
6. Update runtime/reference/contracts/conventions/modules manifests and relevant guides so agents see `W-0222` as the next continuation step.
7. Register `runtime.nakama_aligned_feature_request_workflow_pilot` in `rules/check-rules.json`.
8. Update `tools/vibit` expected-next helpers and runtime check logic.
9. Run verification commands and record the results.

## Files To Create

- `changes/2026-05-24-pilot-nakama-aligned-feature-request-workflow/request.md`
- `changes/2026-05-24-pilot-nakama-aligned-feature-request-workflow/spec.yaml`
- `changes/2026-05-24-pilot-nakama-aligned-feature-request-workflow/impact.md`
- `changes/2026-05-24-pilot-nakama-aligned-feature-request-workflow/plan.md`
- `changes/2026-05-24-pilot-nakama-aligned-feature-request-workflow/checklist.md`
- `changes/2026-05-24-pilot-nakama-aligned-feature-request-workflow/verification.md`
- `decisions/ADR-0129-nakama-aligned-presence-status-workflow-pilot.md`
- `conversations/2026-05-24-nakama-presence-status-workflow-pilot.md`

## Files To Edit

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/contracts.yaml`
- `.arch/conventions.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
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

## Generated Output

No generated output.

## Handwritten Logic

No runtime handwritten behavior in this pilot.

`tools/vibit` receives check and next-ready helper updates only.

## Tests And Verification

Run:

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.nakama_aligned_feature_request_workflow_pilot
node tools/vibit check change pilot-nakama-aligned-feature-request-workflow --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

## Rollback Notes

This is a planning/checks slice. Rollback would remove the pilot records, ADR-0129, rule registration, and W-0222 next-ready update. No data, runtime, generated output, or migration rollback is required.
