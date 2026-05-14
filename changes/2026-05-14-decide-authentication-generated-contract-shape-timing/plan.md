# Plan

## Files To Create

- `docs/authentication-generated-contract-shape-timing.md`
- `docs/authentication-generated-contract-shape-timing.zh-CN.md`
- `decisions/ADR-0038-authentication-generated-contract-shape-timing.md`
- `conversations/2026-05-14-authentication-generated-contract-shape-timing.md`

## Files To Edit

- `docs/generated-output.md`
- `docs/generated-output.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/authentication/module.yaml`
- `modules/authentication/AGENTS.md`
- `modules/authentication/AGENTS.zh-CN.md`

## Generated Artifacts

None. Generated authentication contract shape files remain deferred.

## Handwritten Logic

None. This change does not modify Go runtime code or generator code.

## Tests

No Go tests are required for this timing decision.

## Verification Commands

- `node tools/vibit check contracts --json`
- `node tools/vibit inspect generated --json`
- `node tools/vibit check generated --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change decide-authentication-generated-contract-shape-timing --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

Rollback by removing the timing standard, ADR, manifest references, and work queue update before any generated authentication shape files are produced.

