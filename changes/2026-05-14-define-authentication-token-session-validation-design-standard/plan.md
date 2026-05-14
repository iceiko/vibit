# Plan

## Files To Create

- `docs/authentication-token-session-validation.md`
- `docs/authentication-token-session-validation.zh-CN.md`
- `decisions/ADR-0023-authentication-token-session-validation-design-boundary.md`
- `conversations/2026-05-14-authentication-token-session-validation-design-standard.md`
- `changes/2026-05-14-define-authentication-token-session-validation-design-standard/`

## Files To Edit

- `.arch/conventions.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`

## Generated Artifacts

None.

## Handwritten Logic

None.

## Tests To Add Or Update

No Go tests are required because runtime code is unchanged.

## Verification Commands

```bash
node tools/vibit check architecture --json
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check change define-authentication-token-session-validation-design-standard --json
node tools/vibit check all --json
node tools/vibit inspect next --json
git diff --check
```

## Rollback Notes

Because this is a standard-only change, rollback means removing the new standard/ADR/change/conversation artifacts and restoring the M-011 work queue to `W-0057`.
