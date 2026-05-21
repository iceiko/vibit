# Plan

## Files To Create

- `changes/2026-05-21-refresh-runtime-runbook-for-alpha/request.md`
- `changes/2026-05-21-refresh-runtime-runbook-for-alpha/spec.yaml`
- `changes/2026-05-21-refresh-runtime-runbook-for-alpha/impact.md`
- `changes/2026-05-21-refresh-runtime-runbook-for-alpha/plan.md`
- `changes/2026-05-21-refresh-runtime-runbook-for-alpha/checklist.md`
- `changes/2026-05-21-refresh-runtime-runbook-for-alpha/verification.md`
- `decisions/ADR-0093-runtime-runbook-alpha-path-refresh.md`
- `conversations/2026-05-21-runtime-runbook-alpha-path-refresh.md`

## Files To Edit

- `docs/runtime-runbook.md`
- `docs/runtime-runbook.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `docs/v0.1-alpha-goal.md`
- `docs/v0.1-alpha-goal.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Generated Artifacts

None.

## Handwritten Logic

No runtime handwritten logic changes.

Tooling changes are limited to repository check coverage for the runbook refresh state.

## Tests

No new tests.

The existing E2E proof remains the referenced behavioral verification.

## Verification Commands

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check change refresh-runtime-runbook-for-alpha --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

## Rollback Or Migration Notes

No migration or data rollback is needed. Reverting this slice restores the previous runbook and queue/check metadata.
