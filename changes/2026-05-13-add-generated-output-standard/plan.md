# Plan

## Files To Create

- `docs/generated-output.md`
- `docs/generated-output.zh-CN.md`
- `decisions/ADR-0017-generated-output-standard.md`
- `conversations/2026-05-13-generated-output-standard.md`
- `changes/2026-05-13-add-generated-output-standard/`

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
- `proto/README.md`
- `proto/README.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Generated Artifacts

- None.

## Handwritten Logic

Extend `node tools/vibit check generated` with checks for `runtime/internal/generated/proto/`.

## Tests

No dedicated test framework exists for the CLI yet. Verification is through `node --check tools/vibit` and the repository CLI checks.

## Verification Commands

- `node --check tools/vibit`
- `node tools/vibit check generated`
- `node tools/vibit check generated --json`
- `node tools/vibit check architecture`
- `node tools/vibit check schemas`
- `node tools/vibit check memory`
- `node tools/vibit check protocol`
- `node tools/vibit check runtime`
- `node tools/vibit check change add-generated-output-standard`
- `node tools/vibit check all --json`
- `git diff --check`
- secret scan for GitHub tokens outside ignored local files

## Rollback Or Migration Notes

Rollback can remove the generated-output standard and checks because no generated runtime files are added.
