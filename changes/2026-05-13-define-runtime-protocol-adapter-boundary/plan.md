# Plan

## Files To Create

- `docs/runtime-protocol-adapter.md`
- `docs/runtime-protocol-adapter.zh-CN.md`
- `decisions/ADR-0018-runtime-protocol-adapter-boundary.md`
- `conversations/2026-05-13-runtime-protocol-adapter-boundary.md`
- `changes/2026-05-13-define-runtime-protocol-adapter-boundary/`

## Files To Edit

- `.arch/README.md`
- `.arch/README.zh-CN.md`
- `.arch/conventions.yaml`
- `.arch/protocol.yaml`
- `.arch/runtime.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `docs/game-protocol.md`
- `docs/game-protocol.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Generated Artifacts

- None.

## Handwritten Logic

Extend `node tools/vibit check runtime` to verify runtime protocol adapter boundary references before Go implementation starts.

## Tests

No dedicated CLI test suite exists yet. Verification is through `node --check tools/vibit` and repository checks.

## Verification Commands

- `go version`
- `buf --version`
- `protoc --version`
- `node --check tools/vibit`
- `node tools/vibit check runtime`
- `node tools/vibit check runtime --json`
- `node tools/vibit check architecture`
- `node tools/vibit check schemas`
- `node tools/vibit check memory`
- `node tools/vibit check protocol`
- `node tools/vibit check generated`
- `node tools/vibit check change define-runtime-protocol-adapter-boundary`
- `node tools/vibit check all --json`
- `git diff --check`
- secret scan for GitHub tokens outside ignored local files

## Rollback Or Migration Notes

Rollback can remove the standard, ADR, conversation log, manifest references, and runtime check extensions because no Go runtime source files are added.
