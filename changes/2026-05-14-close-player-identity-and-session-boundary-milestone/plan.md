# Plan

1. Review M-003 completion criteria against standards, manifests, runtime handoff code, inventory permission handoff, and repository checks.
2. Update `docs/player-identity-session-boundary.md` and the Simplified Chinese translation from "next sequence" language to completed boundary milestone language.
3. Update `.arch/runtime.yaml` to record the completed boundary status and next confirmation gate.
4. Mark M-003 and W-0029 completed in `.arch/work-items.yaml`.
5. Do not create a `next_ready` work item because the next meaningful direction requires maintainer confirmation.
6. Run repository verification and record the intentional `check work` warning if no next-ready item remains.

## Files To Edit

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `docs/player-identity-session-boundary.md`
- `docs/player-identity-session-boundary.zh-CN.md`
- `changes/2026-05-14-close-player-identity-and-session-boundary-milestone/*`

## Generated Artifacts

None.

## Handwritten Logic

None.

## Verification Commands

- `cd runtime && go test ./...`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change close-player-identity-and-session-boundary-milestone --json`
- `node tools/vibit check all --json`
- `git diff --check`
