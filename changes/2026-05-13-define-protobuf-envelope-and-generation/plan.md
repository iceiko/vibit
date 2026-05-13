# Plan

## Files To Create

- `buf.yaml`
- `buf.gen.yaml`
- `proto/vibit/protocol/v1/envelope.proto`
- `proto/vibit/inventory/v1/inventory.proto`
- `decisions/ADR-0016-protobuf-envelope-and-generation.md`
- `conversations/2026-05-13-protobuf-envelope-and-generation.md`
- `changes/2026-05-13-define-protobuf-envelope-and-generation/`

## Files To Edit

- `.arch/README.md`
- `.arch/README.zh-CN.md`
- `.arch/contracts.yaml`
- `.arch/conventions.yaml`
- `.arch/protocol.yaml`
- `.arch/runtime.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `docs/game-protocol.md`
- `docs/game-protocol.zh-CN.md`
- `proto/README.md`
- `proto/README.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Generated Artifacts

None.

Generated Go Protobuf output is intentionally not created in this change.

## Handwritten Logic

No server runtime business logic.

Handwritten tooling logic is added to `tools/vibit` to validate the first protocol source files and Buf generation configuration.

## Tests

No runtime tests.

Repository checks and protocol checks are the verification path for this standards/protocol-source change.

## Verification Commands

- `node --check tools/vibit`
- `node tools/vibit check protocol`
- `node tools/vibit check protocol --json`
- `node tools/vibit check architecture`
- `node tools/vibit check schemas`
- `node tools/vibit check memory`
- `node tools/vibit check contracts`
- `node tools/vibit check generated`
- `node tools/vibit check runtime`
- `node tools/vibit check all --json`
- `node tools/vibit check change define-protobuf-envelope-and-generation`
- `git diff --check`
- Secret scan excluding `.git`, `.vibit.local.env`, and `node_modules`

## Rollback Or Migration Notes

Rollback by removing the newly added Protobuf source files, Buf configs, ADR, conversation log, change spec, and reverting manifest/documentation/tooling updates.

No data migration is required.
