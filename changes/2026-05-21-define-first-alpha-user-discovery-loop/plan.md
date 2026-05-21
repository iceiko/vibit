# Plan

## Files To Create

- `docs/first-alpha-user-discovery-loop.md`
- `docs/first-alpha-user-discovery-loop.zh-CN.md`
- `decisions/ADR-0104-first-alpha-user-discovery-loop.md`
- `conversations/2026-05-21-first-alpha-user-discovery-loop.md`
- `changes/2026-05-21-define-first-alpha-user-discovery-loop/spec.yaml`
- `changes/2026-05-21-define-first-alpha-user-discovery-loop/request.md`
- `changes/2026-05-21-define-first-alpha-user-discovery-loop/impact.md`
- `changes/2026-05-21-define-first-alpha-user-discovery-loop/plan.md`
- `changes/2026-05-21-define-first-alpha-user-discovery-loop/checklist.md`
- `changes/2026-05-21-define-first-alpha-user-discovery-loop/verification.md`

## Files To Edit

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `README.md`
- `README.zh-CN.md`
- `docs/v0.1-alpha-goal.md`
- `docs/v0.1-alpha-goal.zh-CN.md`
- `docs/alpha-developer-flow.md`
- `docs/alpha-developer-flow.zh-CN.md`
- `docs/alpha-acceptance-checklist.md`
- `docs/alpha-acceptance-checklist.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Generated Artifacts

None.

## Handwritten Logic

Add a `tools/vibit` repository check for `runtime.first_alpha_user_discovery_loop`.

## Tests

No Go tests are needed. Run repository checks and syntax checks.

## Verification Commands

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check change define-first-alpha-user-discovery-loop --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

## Rollback Or Migration Notes

Rollback is documentation-only: revert the W-0196 docs, ADR, change records, manifest updates, and check-rule additions. No database or runtime rollback is required.
