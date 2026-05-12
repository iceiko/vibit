# Plan

## Files To Create

- `conversations/2026-05-12-check-contracts.md`
- `changes/2026-05-12-add-check-contracts/request.md`
- `changes/2026-05-12-add-check-contracts/spec.yaml`
- `changes/2026-05-12-add-check-contracts/impact.md`
- `changes/2026-05-12-add-check-contracts/plan.md`
- `changes/2026-05-12-add-check-contracts/checklist.md`
- `changes/2026-05-12-add-check-contracts/verification.md`

## Files To Edit

- `tools/vibit`
- `rules/check-rules.json`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`

## Generated Artifacts

None.

## Handwritten Logic

Add CLI check logic inside `tools/vibit`.

## Tests To Add Or Update

No separate test suite exists yet. Verification is performed by running the CLI command in text and JSON modes.

## Verification Commands

- `node tools/vibit check contracts`
- `node tools/vibit check contracts --json`
- `node tools/vibit check all --json`
- `node tools/vibit inspect rule contract.registry_declared`
- Secret scan for GitHub token forms.
- `git diff --check`

## Rollback Notes

If contract source format changes, update or replace this check with the new format-aware validator before changing contract files.
