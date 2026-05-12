# Plan

## Files To Edit

- `tools/vibit`
- `rules/check-rules.json`
- `.arch/README.md`
- `.arch/README.zh-CN.md`
- `.arch/conventions.yaml`
- `.arch/runtime.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `proto/README.md`
- `proto/README.zh-CN.md`

## Files To Create

- `.arch/protocol.yaml`
- `docs/game-protocol.md`
- `docs/game-protocol.zh-CN.md`
- `decisions/ADR-0015-game-protocol-framework.md`
- `conversations/2026-05-13-game-protocol-framework.md`
- `changes/2026-05-13-define-game-protocol-framework/`

## Generated Artifacts

None.

## Handwritten Logic

Extend `node tools/vibit check protocol` to validate the protocol framework manifest before `.proto` files exist.

The check should verify:

- `.arch/protocol.yaml` exists.
- It declares `schema_version: 0.1`, `project: vibit`, and `game_protocol_framework: ADR-0015`.
- It records WebSocket, Protobuf, envelope, message kinds, session model, target model, authority model, error model, compatibility, implementation boundaries, and verification fields.

## Tests To Add Or Update

No separate test framework exists for `tools/vibit` yet.

The change is verified through CLI check commands.

## Verification Plan

Run:

```bash
node --check tools/vibit
node tools/vibit check protocol
node tools/vibit check protocol --json
node tools/vibit check architecture
node tools/vibit check schemas
node tools/vibit check memory
node tools/vibit check contracts
node tools/vibit check generated
node tools/vibit check runtime
node tools/vibit check all --json
node tools/vibit check change define-game-protocol-framework
git diff --check
rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" --hidden --glob '!.git/**' --glob '!.vibit.local.env' --glob '!node_modules/**' .
```

## Rollback

Remove `.arch/protocol.yaml`, the game protocol standard documents, ADR-0015, the protocol manifest check additions, and references to the new standard if the protocol framework proves too early before first runtime implementation.
