# Plan

## Files To Create

- `changes/2026-05-14-close-player-account-postgresql-persistence-milestone/spec.yaml`
- `changes/2026-05-14-close-player-account-postgresql-persistence-milestone/request.md`
- `changes/2026-05-14-close-player-account-postgresql-persistence-milestone/impact.md`
- `changes/2026-05-14-close-player-account-postgresql-persistence-milestone/plan.md`
- `changes/2026-05-14-close-player-account-postgresql-persistence-milestone/checklist.md`
- `changes/2026-05-14-close-player-account-postgresql-persistence-milestone/verification.md`

## Files To Edit

- `.arch/work-items.yaml`
- `.arch/contracts.yaml`
- `modules/player/AGENTS.md`
- `modules/player/AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`

## Generated Artifacts

None.

## Handwritten Logic

None.

## Tests

No new tests are required for this workflow closeout. Re-run the focused PostgreSQL adapter tests because the closeout depends on their completed behavior.

## Verification Commands

- `cd runtime && go test ./internal/platform/persistence/postgres`
- `cd runtime && go test ./...`
- `node tools/vibit check runtime --json`
- `node tools/vibit check module player --json`
- `node tools/vibit check contracts --json`
- `node tools/vibit check migrations --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change close-player-account-postgresql-persistence-milestone --json`
- `node tools/vibit check all --json`
- `node tools/vibit inspect next --json`
- `git diff --check`

## Rollback Or Migration Notes

No data rollback is needed because this change does not add or modify migrations or runtime behavior. If the closeout is premature, revert only the workflow closeout and confirmation gate metadata, not the already completed persistence adapter work.
