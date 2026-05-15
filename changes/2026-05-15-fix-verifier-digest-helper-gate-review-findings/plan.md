# Plan

## Files To Create

- `conversations/2026-05-15-verifier-digest-helper-implementation-gate.md`
- `changes/2026-05-15-fix-verifier-digest-helper-gate-review-findings/`

## Files To Edit

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `AGENTS.zh-CN.md`
- `changes/2026-05-15-define-verifier-digest-helper-implementation-gate/verification.md`
- `tools/vibit`

## Generated Artifacts

None.

## Handwritten Logic

Add a focused `checkRuntimeVerifierDigestHelperImplementationGate` function to `tools/vibit` and call it from `checkRuntime`.

## Tests

No Go tests are required. Validate with repository checks.

## Verification Commands

- `node -c tools/vibit`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change fix-verifier-digest-helper-gate-review-findings --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

Rollback is a normal git revert. No runtime data or migration state is affected.
