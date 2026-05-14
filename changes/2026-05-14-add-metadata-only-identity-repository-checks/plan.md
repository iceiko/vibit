# Plan

1. Add `runtime.identity_boundary` to the known rule IDs and rule catalog.
2. Extend `node tools/vibit check runtime` with metadata-only identity boundary checks.
3. Check WebSocket transport imports against domain, player, inventory, generated Protobuf, and Protobuf runtime packages.
4. Check domain module imports against WebSocket, generated Protobuf, Protobuf runtime, known auth/token/credential/password dependencies.
5. Check that player runtime implementation, player Protobuf source roots, and player/account migrations remain absent until ratified.
6. Check that `modules/player/module.yaml` still declares boundary-only public API and ratification markers.
7. Document the new check in the player identity/session standard and runtime manifest.
8. Record deferred checks that are too brittle for the current repository shape.
9. Run repository verification.

## Files To Edit

- `tools/vibit`
- `rules/check-rules.json`
- `.arch/runtime.yaml`
- `docs/player-identity-session-boundary.md`
- `docs/player-identity-session-boundary.zh-CN.md`
- `.arch/work-items.yaml`
- `changes/2026-05-14-add-metadata-only-identity-repository-checks/*`

## Generated Artifacts

None.

## Verification Commands

- `node -c tools/vibit`
- `node -e "JSON.parse(require('fs').readFileSync('rules/check-rules.json','utf8'))"`
- `node tools/vibit inspect rule runtime.identity_boundary`
- `node tools/vibit check runtime --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change add-metadata-only-identity-repository-checks --json`
- `node tools/vibit check all --json`
- `git diff --check`
