# Plan

## Files To Create

- `docs/reference-game-server-alignment.md`
- `docs/reference-game-server-alignment.zh-CN.md`
- `.arch/reference.yaml`
- `decisions/ADR-0019-nakama-and-pitaya-reference-baseline.md`
- `conversations/2026-05-13-nakama-and-pitaya-reference-alignment.md`

## Files To Edit

- `.arch/README.md`
- `.arch/README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `.arch/conventions.yaml`
- `tools/vibit`

## Generated Artifacts

None.

## Handwritten Logic

Small architecture check extension in `tools/vibit` to require the reference manifest, standard, ADR, and entry-point references.

## Tests

No runtime tests are needed for this standards-only change.

## Verification Commands

- `node tools/vibit check architecture --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check change align-with-nakama-and-pitaya --json`
- `node tools/vibit check all --json`
- `git diff --check`
- secret scan for GitHub tokens outside ignored local files

## Rollback Or Migration Notes

Rollback can remove the new standard, ADR, reference manifest, and intake references. No code, public protocol, or database migration is involved.
