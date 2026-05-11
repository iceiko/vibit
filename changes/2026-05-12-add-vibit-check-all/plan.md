# Plan

## Files To Create

- `conversations/2026-05-12-check-all-and-open-server-decisions.md`

## Files To Edit

- `tools/vibit`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `changes/2026-05-12-add-vibit-check-all/*`

## Generated Artifacts

- None.

## Handwritten Logic

- Add `check all` command dispatch.
- Add discovery for change specs under `changes/`.
- Add lightweight discovery for registered modules from `.arch/modules.yaml`.
- Run subchecks and aggregate failures.

## Tests

Manual CLI verification:

- `node tools/vibit --help`
- `node tools/vibit check all`
- `node tools/vibit check architecture`
- `node tools/vibit check change add-vibit-check-all`
- `node tools/vibit check module inventory`

## Verification Commands

- `node tools/vibit --help`
- `node tools/vibit check all`
- `node tools/vibit check architecture`
- `node tools/vibit check change add-vibit-check-all`
- `node tools/vibit check module inventory`
- `rg -n "ghp_[A-Za-z0-9]|github_pat_[A-Za-z0-9]" .`

## Rollback Or Migration Notes

Remove the `check all` command and documentation references if aggregate checks prove premature.
