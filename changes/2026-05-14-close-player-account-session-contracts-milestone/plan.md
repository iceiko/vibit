# Plan

1. Confirm `W-0031` through `W-0035` satisfy the `M-005` completion criteria.
2. Update `docs/player-account-session-contracts.md` and the Simplified Chinese translation with milestone completion status.
3. Update `.arch/runtime.yaml` and `.arch/reference.yaml` so their active phase/status fields no longer imply that `M-005` is still in progress.
4. Mark `M-005` completed in `.arch/work-items.yaml`.
5. Add `W-0036` as the completed closure work item.
6. Add a blocked `W-0037` confirmation gate for the next milestone direction.
7. Run repository and runtime verification.

## Files To Edit

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `docs/player-account-session-contracts.md`
- `docs/player-account-session-contracts.zh-CN.md`
- `changes/2026-05-14-add-runtime-contract-inspection/*`
- `changes/2026-05-14-close-player-account-session-contracts-milestone/*`

## Generated Artifacts

None.

## Handwritten Logic

None.

## Verification Commands

- `node tools/vibit inspect work --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change add-runtime-contract-inspection --json`
- `node tools/vibit check change close-player-account-session-contracts-milestone --json`
- `node tools/vibit check contracts --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check all --json`
- `cd runtime && go test ./...`
- `cd runtime && go vet ./...`
- `git diff --check`
