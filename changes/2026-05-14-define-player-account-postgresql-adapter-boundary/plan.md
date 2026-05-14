# Plan

## Files To Create

- `changes/2026-05-14-define-player-account-postgresql-adapter-boundary/`

## Files To Edit

- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `.arch/work-items.yaml`
- `modules/player/module.yaml`
- `modules/player/AGENTS.md`
- `modules/player/AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `docs/postgresql-persistence-boundary.md`
- `docs/postgresql-persistence-boundary.zh-CN.md`
- `tools/vibit`

## Generated Artifacts

- None.

## Handwritten Logic

- No Go runtime adapter implementation.
- `tools/vibit` runtime identity checks are updated so the adapter boundary is machine-verifiable.

## Tests

- No new Go tests are required because no Go adapter implementation is added.
- Existing runtime and architecture checks must verify the new boundary.

## Verification Commands

```bash
node tools/vibit check runtime --json
node tools/vibit check module player --json
node tools/vibit check contracts --json
node tools/vibit check migrations --json
node tools/vibit check work --json
node tools/vibit check change define-player-account-postgresql-adapter-boundary --json
node tools/vibit check all --json
node tools/vibit inspect next --json
git diff --check
```

## Rollback Or Migration Notes

No database rollback is needed. This change defines an adapter boundary only and does not change migration source or runtime behavior.
